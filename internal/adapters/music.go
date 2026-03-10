package adapters

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

var musicSupportedExt = make(map[string][]string, 16)

func init() {
	importFormats := []string{"mp3", "wav", "flac", "aac"}
	exportFormats := []string{"mp3", "wav", "flac", "aac"}

	for _, src := range importFormats {
		musicSupportedExt[src] = exportFormats
	}
}

type MusicConverter struct {
	FFmpegPath string
}

func NewMusicConverter() MusicConverter {
	return MusicConverter{
		FFmpegPath: "ffmpeg",
	}
}

func (c MusicConverter) Supports(ctx context.Context, srcExt, dstExt string) bool {
	allowedExt, ok := musicSupportedExt[srcExt]
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

func (c MusicConverter) Convert(ctx context.Context, srcPath, dstPath, dstExt string) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("musicConverter file doesn't exist: %w", err)
	}

	args := []string{
		"-y",
		"-i", srcPath,
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

		return fmt.Errorf("musicConverter ffmpeg error: %w, stderr: %s", err, stderr.String())
	}

	return nil
}
