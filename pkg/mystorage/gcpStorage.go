package mystorage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"github.com/DeepAung/deep-art/config"
)

type gcpStorage struct {
	cfg config.IAppConfig
}

func NewGCPStorage(cfg config.IAppConfig) IStorage {
	return &gcpStorage{
		cfg: cfg,
	}
}

func (s *gcpStorage) UploadFiles(files []*multipart.FileHeader, dir string) ([]*FileRes, error) {
	filesInfo := newFilesInfo(files, dir)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	length := len(filesInfo)
	jobsCh := make(chan *fileInfo, length)
	resultsCh := make(chan *FileRes, length)
	errorsCh := make(chan error, length)

	results := make([]*FileRes, length)

	// Start workers
	numWorkers := 5 // TODO:
	for id := 1; id <= numWorkers; id++ {
		go func(id int) {
			s.uploadWorker(ctx, client, jobsCh, resultsCh, errorsCh, id)
		}(id)
	}

	// Send jobs to workers
	for _, info := range filesInfo {
		jobsCh <- info
	}
	close(jobsCh)

	// Get results
	for i := 0; i < length; i++ {
		err := <-errorsCh
		if err != nil {
			return nil, err
		}

		results[i] = <-resultsCh
	}

	return results, nil
}

func (s *gcpStorage) DeleteFiles(destinations []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	length := len(destinations)
	fmt.Println("length: ", length)
	jobsCh := make(chan string, length)
	errorsCh := make(chan error, length)

	// Start workers
	numWorkers := 5 // TODO:
	for id := 1; id <= numWorkers; id++ {
		go func(id int) {
			s.deleteWorkers(ctx, client, jobsCh, errorsCh, id)
		}(id)
	}

	// Send jobs to workers
	for _, dest := range destinations {
		fmt.Println("send job: ", dest)
		jobsCh <- dest
	}
	close(jobsCh)

	// Get errors
	for i := 0; i < length; i++ {
		err := <-errorsCh
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *gcpStorage) uploadWorker(
	ctx context.Context,
	client *storage.Client,
	jobsCh <-chan *fileInfo,
	resultsCh chan<- *FileRes,
	errorsCh chan<- error,
	id int,
) {
	for job := range jobsCh {
		fmt.Printf("worker %d start: %v upload to %v.\n", id, job.Filename, job.Destination)

		f, err := job.File.Open()
		if err != nil {
			errorsCh <- err
			return
		}

		bucket := s.cfg.GCPBucket()
		o := client.Bucket(bucket).Object(job.Destination)

		wc := o.NewWriter(ctx)
		if _, err = io.Copy(wc, f); err != nil {
			errorsCh <- fmt.Errorf("io.Copy: %w", err)
			return
		}
		if err := wc.Close(); err != nil {
			errorsCh <- fmt.Errorf("Writer.Close: %w", err)
			return
		}

		if err := s.makePublic(ctx, client, bucket, job.Destination); err != nil {
			errorsCh <- err
			return
		}

		errorsCh <- nil
		resultsCh <- newGCPFileRes(job.Filename, job.Destination, bucket)

		fmt.Printf("worker %d finish: %v upload to %v.\n", id, job.Filename, job.Destination)
	}
}

func (s *gcpStorage) makePublic(
	ctx context.Context,
	client *storage.Client,
	bucket, destination string,
) error {
	acl := client.Bucket(bucket).Object(destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %v", err)
	}

	fmt.Printf("Blob %v is now publicly accessible.\n", destination)
	return nil
}

func (s *gcpStorage) deleteWorkers(
	ctx context.Context,
	client *storage.Client,
	jobsCh <-chan string,
	errorsCh chan<- error,
	id int,
) {
	for dest := range jobsCh {
		fmt.Printf("worker %d start: delete %s.\n", id, dest)

		o := client.Bucket(s.cfg.GCPBucket()).Object(dest)

		attrs, err := o.Attrs(ctx)
		if err != nil {
			errorsCh <- fmt.Errorf("object.Attrs: %v", err)
			return
		}
		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		if err := o.Delete(ctx); err != nil {
			errorsCh <- fmt.Errorf("Object(%q).Delete: %v", dest, err)
			return
		}

		errorsCh <- nil

		fmt.Printf("worker %d finish: delete %s.\n", id, dest)
	}
}
