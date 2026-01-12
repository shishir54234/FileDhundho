package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"sync"
)

type CompressEngine struct {
	workers int
	mu      sync.Mutex
}
type FileJob struct {
	Path    string
	RelPath string
	IsDir   bool
}

type CompressedFile struct {
	RelPath  string
	Data     []byte
	IsDir    bool
	OrigSize int64
}

// maybe some code to zip a file lets see if it works first

func NewCompressEngine(workers int) *CompressEngine {
	if workers <= 0 {
		workers = 4
	}
	return &CompressEngine{workers: workers}
}

func (e *CompressEngine) CompressFileZip(sourcePath string, destZipPath string) error {
	// implement zipping logic here
	file, err := os.Create(destZipPath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	var fileJobs []FileJob

	err = filepath.Walk(sourcePath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		fileJobs = append(fileJobs, FileJob{
			Path:    path,
			RelPath: relPath,
			IsDir:   fileInfo.IsDir(),
		})

		return nil
	})

	if err != nil {
		return err
	}

	jobs := make(chan FileJob, len(fileJobs))
	results := make(chan CompressedFile, len(fileJobs))
	errors := make(chan error, len(fileJobs))

	var wg sync.WaitGroup

	for i := 0; i < e.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				var compressed CompressedFile
				if job.IsDir {
					compressed = CompressedFile{
						RelPath: job.RelPath,
						IsDir:   true,
					}
				} else {
					fileData, err := os.ReadFile(job.Path)
					if err != nil {
						errors <- err
						continue
					}
					compressed = CompressedFile{
						RelPath:  job.RelPath,
						Data:     fileData,
						IsDir:    false,
						OrigSize: int64(len(fileData)),
					}
				}
				results <- compressed
			}
		}()
	}

	for _, job := range fileJobs {
		jobs <- job
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// now we have to take stuff from the results and write it to one place
	for compressed := range results {
		if compressed.IsDir {
			_, err := zipWriter.Create(compressed.RelPath + "/")
			if err != nil {
				return err
			}
		} else {
			writer, err := zipWriter.Create(compressed.RelPath)
			if err != nil {
				return err
			}
			_, err = writer.Write(compressed.Data)
			if err != nil {
				return err
			}
		}
	}
	return nil
} 



