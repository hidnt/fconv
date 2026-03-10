package adapters

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

var mus2VidSupportedExt = make(map[string][]string, 16)

func init() {
	importFormats := []string{"mp4", "mkv", "mov", "avi", "wmv", "ogg"}
	exportFormats := []string{"mp3", "wav", "flac", "aac"}

	for _, src := range importFormats {
		mus2VidSupportedExt[src] = exportFormats
	}
}

type Mus2VidConverter struct {
	FFmpegPath string
}

func NewMus2VidConverter() Mus2VidConverter {
	return Mus2VidConverter{
		FFmpegPath: "ffmpeg",
	}
}

func (c Mus2VidConverter) Supports(ctx context.Context, srcExt, dstExt string) bool {
	allowedExt, ok := mus2VidSupportedExt[srcExt]
	if !ok {
		return false
	}

	for _, allowed := range allowedExt {
		if allowed == dstExt {
			return true
		}
	}

	return false
}

func (c Mus2VidConverter) Convert(ctx context.Context, srcPath, dstPath, dstExt string) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("mus2VidConverter file doesn't exist: %w", err)
	}

	args := []string{
		"-y",
		"-i", srcPath,
		"-vn",
		"-q:a", "0",
		dstPath,
	}

	cmd := exec.CommandContext(ctx, c.FFmpegPath, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		_ = os.Remove(dstPath)

		if ctx.Err() != nil {
			return ctx.Err()
		}

		return fmt.Errorf("mus2VidConverter ffmpeg error: %w, stderr: %s", err, stderr.String())
	}

	return nil
}
