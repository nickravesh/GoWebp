package converter

import (
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
)

type Settings struct {
	Quality        float32
	MaxWorkers     int
	AllowOverwrite bool
}

func NewSettings() *Settings {
	return &Settings{
		Quality:        85,
		MaxWorkers:     4,
		AllowOverwrite: false,
	}
}

type Job struct {
	SrcPath string
	DstPath string
	RelPath string
}

type Result struct {
	Job   Job
	Error error
}

func FindImages(paths []string) ([]Job, error) {
	var jobs []Job
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			err := filepath.Walk(p, func(path string, fi os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !fi.IsDir() && isImageFile(path) {
					rel, _ := filepath.Rel(p, path)
					jobs = append(jobs, Job{SrcPath: path, RelPath: rel})
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else if isImageFile(p) {
			jobs = append(jobs, Job{SrcPath: p, RelPath: filepath.Base(p)})
		}
	}
	return jobs, nil
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func ConvertImage(job Job, outputDir string, settings *Settings) error {
	dstPath := filepath.Join(outputDir, strings.TrimSuffix(job.RelPath, filepath.Ext(job.RelPath))+".webp")
	job.DstPath = dstPath

	if !settings.AllowOverwrite {
		if _, err := os.Stat(dstPath); err == nil {
			return fmt.Errorf("file exists: %s", dstPath)
		}
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	in, err := os.Open(job.SrcPath)
	if err != nil {
		return err
	}
	defer in.Close()

	img, _, err := image.Decode(in)
	if err != nil {
		return err
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	opts := &webp.Options{Quality: settings.Quality}
	if err := webp.Encode(out, img, opts); err != nil {
		return err
	}
	return nil
}

func WorkerPool(jobs []Job, outputDir string, settings *Settings, logWriter io.Writer, progress func(int, int)) {
	numJobs := len(jobs)
	jobCh := make(chan Job)
	resultCh := make(chan Result)
	var wg sync.WaitGroup

	// Workers
	for i := 0; i < settings.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				err := ConvertImage(job, outputDir, settings)
				resultCh <- Result{Job: job, Error: err}
			}
		}()
	}

	// Dispatcher
	go func() {
		for _, job := range jobs {
			jobCh <- job
		}
		close(jobCh)
	}()

	// Collector
	go func() {
		success, fail := 0, 0
		for i := 0; i < numJobs; i++ {
			res := <-resultCh
			if res.Error != nil {
				fail++
				fmt.Fprintf(logWriter, "FAIL: %s: %v\n", res.Job.SrcPath, res.Error)
			} else {
				success++
				fmt.Fprintf(logWriter, "OK: %s\n", res.Job.SrcPath)
			}
			progress(success+fail, numJobs)
		}
		close(resultCh)
	}()

	wg.Wait()
}
