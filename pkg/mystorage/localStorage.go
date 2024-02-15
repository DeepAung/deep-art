package mystorage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/pkg/utils"
)

type localStorage struct {
	cfg config.IAppConfig
}

func NewLocalStorage(cfg config.IAppConfig) IStorage {
	return &localStorage{
		cfg: cfg,
	}
}

// read-only variable
func (s *localStorage) storagePath() string { return "./public/storage/" }

func (s *localStorage) UploadFiles(files []*multipart.FileHeader, dir string) ([]*FileRes, error) {
	filesInfo := newFilesInfo(files, dir)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	length := len(filesInfo)
	jobsCh := make(chan *fileInfo, length)
	resultsCh := make(chan *FileRes, length)
	errorsCh := make(chan error, length)

	results := make([]*FileRes, length)

	// Start workers
	numWorkers := 5 // TODO:
	for id := 1; id <= numWorkers; id++ {
		go func(id int) {
			s.uploadWorker(ctx, jobsCh, resultsCh, errorsCh, id)
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

func (s *localStorage) DeleteFiles(destinations []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	length := len(destinations)
	jobsCh := make(chan string, length)
	errorsCh := make(chan error, length)

	// Start workers
	numWorkers := 5 // TODO:
	for id := 1; id <= numWorkers; id++ {
		go func(id int) {
			s.deleteWorkers(ctx, jobsCh, errorsCh, id)
		}(id)
	}

	// Send jobs to workers
	for _, dest := range destinations {
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

func (s *localStorage) uploadWorker(
	ctx context.Context,
	jobsCh <-chan *fileInfo,
	resultsCh chan<- *FileRes,
	errorsCh chan<- error,
	id int,
) {
	for job := range jobsCh {
		f, err := job.File.Open()
		if err != nil {
			errorsCh <- err
			return
		}

		b, err := io.ReadAll(f)
		if err != nil {
			errorsCh <- err
			return
		}

		// Upload an object to storage
		filePath := utils.Join(s.storagePath(), job.Destination)
		dirPath := utils.Join(s.storagePath(), job.Dir)

		if err := os.WriteFile(filePath, b, 0777); err != nil {
			if err := os.MkdirAll(dirPath, 0777); err != nil {
				errorsCh <- fmt.Errorf("mkdir \"%s\" failed: %v", dirPath, err)
				return
			}

			if err := os.WriteFile(filePath, b, 0777); err != nil {
				errorsCh <- fmt.Errorf("write file failed: %v", err)
				return
			}
		}

		errorsCh <- nil
		resultsCh <- newLocalFileRes(s.cfg, job.Filename, job.Destination)
	}
}

func (s *localStorage) deleteWorkers(
	ctx context.Context,
	jobsCh <-chan string,
	errorsCh chan<- error,
	id int,
) {
	for dest := range jobsCh {
		filePath := utils.Join(s.storagePath(), dest)
		if err := os.Remove(filePath); err != nil {
			errorsCh <- fmt.Errorf("remove file: \"%s\" failed: %v", dest, err)
			return
		}
		errorsCh <- nil
	}
}
