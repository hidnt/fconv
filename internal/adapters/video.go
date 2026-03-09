package adapters

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var videoSupportedExt sync.Map

func init() {
	importFormats := []string{"mp4", "mkv", "mov", "avi", "gif", "wmv", "ogg"}
	exportFormats := []string{"mp4", "mkv", "mov", "avi", "gif", "wmv"}

	for _, src := range importFormats {
		videoSupportedExt.Store(src, exportFormats)
	}
}

type VideoConverter struct {
	FFmpegPath string
}

func NewVideoConverter() VideoConverter {
	return VideoConverter{
		FFmpegPath: "ffmpeg",
	}
}

func (c VideoConverter) CheckFFmpeg() error {
	cmd := exec.Command(c.FFmpegPath, "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg not found or not working: %w", err)
	}
	return nil
}

func (c VideoConverter) Supports(ctx context.Context, srcExt, dstExt string) bool {
	val, ok := videoSupportedExt.Load(srcExt)
	if !ok {
		return false
	}

	allowedDsts := val.([]string)
	for _, allowed := range allowedDsts {
		if allowed == dstExt {
			return true
		}
	}

	return false
}

func (c VideoConverter) Convert(ctx context.Context, srcPath, dstPath, dstExt string) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("videoConverter file doesn't exist: %w", err)
	}

	args := []string{
		"-y",
		"-i", srcPath,
	}

	switch dstExt {
	case "gif":
		args = append(args, "-vf", "fps=15,scale=480:-1:flags=lanczos")
	case "mp4":
		args = append(args, "-c:v", "libx264", "-c:a", "aac", "-strict", "experimental")
	}

	args = append(args, dstPath)
	cmd := exec.CommandContext(ctx, c.FFmpegPath, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		_ = os.Remove(dstPath)

		if ctx.Err() != nil {
			return ctx.Err()
		}

		return fmt.Errorf("videoConverter ffmpeg error: %w, stderr: %s", err, stderr.String())
	}

	return nil
}
