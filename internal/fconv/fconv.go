package fconv

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/hidnt/fconv/internal/models"
	"github.com/hidnt/fconv/internal/service"
	"golang.org/x/sync/semaphore"
)

var (
	bufferLength = 1024
)

type App struct {
	cfg  models.Config
	jobs chan string
	wg   sync.WaitGroup
	srv  service.ConverterService
}

func New(cfg models.Config, srv service.ConverterService) App {
	return App{
		cfg:  cfg,
		jobs: make(chan string, bufferLength),
		srv:  srv,
	}
}

func (a *App) weight(path string, maxWorkers int64) int64 {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))

	switch ext {
	case "png", "jpg", "jpeg", "webp", "bmp", "tiff", "avif", "ico", "cur", "heic", "heif":
		return 1
	default:
		weight := maxWorkers / 2
		if weight < 1 {
			return 1
		}
		return weight
	}
}

func (a *App) dispatcher(ctx context.Context) {
	maxWorkers := int64(runtime.NumCPU())
	sem := semaphore.NewWeighted(maxWorkers)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		for path := range a.jobs {
			weight := a.weight(path, maxWorkers)
			if err := sem.Acquire(ctx, weight); err != nil {
				break
			}

			a.wg.Add(1)
			go func(path string, weight int64) {
				defer a.wg.Done()
				defer sem.Release(weight)
				a.processFile(ctx, path)
			}(path, weight)
		}
	}()
}

func (a *App) processFile(ctx context.Context, srcPath string) {
	rawExt := filepath.Ext(srcPath)
	srcExt := strings.ToLower(strings.TrimPrefix(rawExt, "."))
	basename := strings.TrimSuffix(filepath.Base(srcPath), rawExt)

	var dstBasePath string
	if a.cfg.DstFolder == "" {
		dstBasePath = filepath.Dir(srcPath)
	} else {
		dstBasePath = a.cfg.DstFolder
	}
	dstPath := filepath.Join(dstBasePath, basename+"."+a.cfg.DstExt)

	if _, err := os.Stat(dstPath); err == nil && !a.cfg.Force {
		fmt.Fprintln(os.Stderr, "File already exists: "+dstPath)
		return
	}

	fmt.Println("Converting: " + srcPath)
	if err := a.srv.Convert(ctx, srcPath, dstPath, srcExt, a.cfg.DstExt); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to convert "+srcPath+": "+err.Error())
		return
	}
	if a.cfg.Delete {
		os.Remove(srcPath)
	}
	fmt.Println("Сonversion: " + srcPath + " completed")
}

func (a *App) recWalkDir(ctx context.Context, path string, currentLevel int) {
	if a.cfg.LevelOfRec != -1 && currentLevel >= a.cfg.LevelOfRec {
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ReadDir error: "+path)
	}

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return
		default:
			path := filepath.Join(path, entry.Name())
			if entry.IsDir() {
				a.recWalkDir(ctx, path, currentLevel+1)
			} else {
				a.jobs <- path
			}
		}
	}
}

func (a *App) Fconv(ctx context.Context, paths []string) {
	a.dispatcher(ctx)

	if !a.cfg.NeedRecursion {
		a.cfg.LevelOfRec = 1
	}

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to obtain information about "+path)
			continue
		}
		if stat.IsDir() {
			a.recWalkDir(ctx, path, 0)
			continue
		}

		select {
		case <-ctx.Done():
			close(a.jobs)
			return
		case a.jobs <- path:
		}
	}

	close(a.jobs)
	a.wg.Wait()
}
