package storer

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"

	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/utils"
)

type localStorer struct {
	cfg *config.Config
}

func NewLocalStorer(cfg *config.Config) *localStorer {
	return &localStorer{
		cfg: cfg,
	}
}

func (s *localStorer) UploadFiles(files []*multipart.FileHeader, dir string) ([]FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.App.Timeout)
	defer cancel()

	length := len(files)
	jobCh := make(chan *multipart.FileHeader, length)
	resCh := make(chan FileRes, length)
	errCh := make(chan error, length)

	// spawn workers
	numWorkers := 5
	for range numWorkers {
		go func() {
			for file := range jobCh {
				select {
				case <-ctx.Done():
					errCh <- fmt.Errorf("upload file cancelled")
					return
				default:
				}

				res, err := s.uploadFile(cancel, file, dir)
				if err != nil {
					errCh <- err
					return
				}

				errCh <- nil
				resCh <- res
			}
		}()
	}

	// assign jobs
	for _, file := range files {
		jobCh <- file
	}
	close(jobCh)

	// wait for results and error
	var results []FileRes
	for range length {
		err := <-errCh
		if err != nil {
			return results, err
		}

		results = append(results, <-resCh)
	}

	return results, nil
}

func (s *localStorer) DeleteFiles(dests []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.App.Timeout)
	defer cancel()

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
					errCh <- fmt.Errorf("delete file cancelled")
					return
				default:
				}

				err := s.deleteFile(cancel, dest)
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
		err := <-errCh
		if err != nil {
			return err
		}
	}

	return nil
}

// res, err := s.uploadFile(ctx, cancel, file, dir)
func (s *localStorer) uploadFile(
	cancel func(),
	file *multipart.FileHeader,
	dir string,
) (FileRes, error) {
	f, err := file.Open()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("file.Open: %w", err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	// Upload an object to storage
	dest := utils.Join(dir, file.Filename)
	filePath := "." + utils.Join(s.cfg.App.BasePath, dest)
	dirPath := "." + utils.Join(s.cfg.App.BasePath, dir)

	if err := os.WriteFile(filePath, b, fs.ModePerm); err != nil {
		if err := os.MkdirAll(dirPath, fs.ModePerm); err != nil {
			return nil, fmt.Errorf("mkdir \"%s\" failed: %v", dirPath, err)
		}

		if err := os.WriteFile(filePath, b, fs.ModePerm); err != nil {
			return nil, fmt.Errorf("write file failed: %v", err)
		}
	}

	return utils.NewUrlInfoByDest(s.cfg.App.BasePath, dest), nil
}

func (s *localStorer) deleteFile(
	cancel func(),
	dest string,
) error {
	filePath := "." + utils.Join(s.cfg.App.BasePath, dest)
	if err := os.Remove(filePath); err != nil {
		cancel()
		return fmt.Errorf("remove file: \"%s\" failed: %v", dest, err)
	}

	return nil
}
