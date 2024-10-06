package storer

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/utils"
)

type GCPStorer struct {
	cfg *config.Config
}

func NewGCPStorer(cfg *config.Config) Storer {
	return &GCPStorer{
		cfg: cfg,
	}
}

func (s *GCPStorer) UploadFile(file io.Reader, dest string) (FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.App.Timeout)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient failed: %w", err)
	}
	defer client.Close()

	return s.uploadFile(ctx, client, file, dest)
}

func (s *GCPStorer) UploadFiles(files []io.Reader, dests []string) ([]FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.App.Timeout)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return []FileRes{}, fmt.Errorf("storage.NewClient failed: %w", err)
	}
	defer client.Close()

	length := len(files)
	jobCh := make(chan struct {
		file io.Reader
		dest string
	}, length)
	resCh := make(chan FileRes, length)
	errCh := make(chan error, length)

	// spawn workers
	numWorkers := 5
	for range numWorkers {
		go func() {
			for job := range jobCh {
				select {
				case <-ctx.Done():
					errCh <- fmt.Errorf("upload file cancelled (timeout)")
					return
				default:
					res, err := s.uploadFile(ctx, client, job.file, job.dest)
					if err != nil {
						errCh <- err
						return
					}
					errCh <- nil
					resCh <- res
				}
			}
		}()
	}

	// assign jobs
	for i, file := range files {
		jobCh <- struct {
			file io.Reader
			dest string
		}{file: file, dest: dests[i]}
	}
	close(jobCh)

	// wait for results and error
	var results []FileRes
	for range length {
		err = <-errCh
		if err != nil {
			return results, err
		}

		results = append(results, <-resCh)
	}

	return results, nil
}

func (s *GCPStorer) uploadFile(
	ctx context.Context,
	client *storage.Client,
	file io.Reader,
	dest string,
) (FileRes, error) {
	bucket := s.cfg.App.GcpBucket
	if dest[0] == '/' {
		dest = dest[1:]
	}

	o := client.Bucket(bucket).Object(dest)

	wc := o.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return nil, fmt.Errorf("io.Copy failed: %w", err)
	}
	if err := wc.Close(); err != nil {
		return nil, fmt.Errorf("Writer.Close failed: %w", err)
	}

	fmt.Printf("Blob %v uploaded.\n", dest)
	return utils.NewUrlInfoByDest(s.cfg.App.BasePath, dest), nil
}

func (s *GCPStorer) DeleteFile(dest string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.App.Timeout)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient failed: %w", err)
	}
	defer client.Close()

	return s.deleteFile(ctx, client, dest)
}

func (s *GCPStorer) DeleteFiles(dests []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.App.Timeout)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient failed: %w", err)
	}
	defer client.Close()

	length := len(dests)
	jobCh := make(chan string, length)
	errCh := make(chan error, length)

	// spawn workers
	numWorkers := 5
	for range numWorkers {
		go func() {
			for dest := range jobCh {
				select {
				case <-ctx.Done():
					errCh <- fmt.Errorf("delete file cancelled (timeout)")
					return
				default:
				}

				err := s.deleteFile(ctx, client, dest)
				if err != nil {
					errCh <- err
					return
				}

				errCh <- nil
			}
		}()
	}

	// assign jobs
	for _, dest := range dests {
		jobCh <- dest
	}
	close(jobCh)

	// wait for error
	for range length {
		err = <-errCh
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *GCPStorer) deleteFile(
	ctx context.Context,
	client *storage.Client,
	dest string,
) error {
	if dest[0] == '/' {
		dest = dest[1:]
	}

	bucket := s.cfg.App.GcpBucket
	o := client.Bucket(bucket).Object(dest)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs failed: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete failed: %w", dest, err)
	}

	fmt.Printf("Blob %v deleted.\n", dest)
	return nil
}
