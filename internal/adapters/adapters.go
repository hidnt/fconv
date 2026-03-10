package adapters

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/google/wire"
	"github.com/hidnt/fconv/internal/service"
)

type ctxWriter struct {
	ctx context.Context
	w   io.Writer
}

func (cw ctxWriter) Write(p []byte) (int, error) {
	select {
	case <-cw.ctx.Done():
		return 0, cw.ctx.Err()
	default:
		return cw.w.Write(p)
	}
}

func checkFFmpeg(ffmpegPath string) error {
	cmd := exec.Command(ffmpegPath, "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg not found or not working: %w", err)
	}
	return nil
}

var AdaptersSet = wire.NewSet(
	NewImageConverter,
	NewVideoConverter,
	NewMusicConverter,
	NewMus2VidConverter,
	ProvideConverters,
)

func ProvideConverters(
	a1 ImageConverter,
	a2 VideoConverter,
	a3 MusicConverter,
	a4 Mus2VidConverter,
) []service.Converter {
	return []service.Converter{
		a1,
		a2,
		a3,
		a4,
	}
}
