package adapters

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"sync"

	"github.com/HugoSmits86/nativewebp"
	"github.com/gen2brain/avif"
	"github.com/sergeymakinen/go-ico"
	"github.com/sergeymakinen/go-ico/cur"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	"golang.org/x/image/tiff"

	_ "image/jpeg"
	_ "image/png"

	_ "github.com/gen2brain/avif"
	_ "github.com/gen2brain/heic"
	_ "github.com/sergeymakinen/go-ico"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

var imageSupportedExt sync.Map

func init() {
	importFormats := []string{"png", "jpg", "jpeg", "webp", "bmp", "tiff", "avif", "ico", "cur", "heic", "heif"}
	exportFormats := []string{"png", "jpg", "jpeg", "webp", "bmp", "tiff", "avif", "ico", "cur"}

	for _, src := range importFormats {
		imageSupportedExt.Store(src, exportFormats)
	}
}

type ImageConverter struct{}

func NewImageConverter() ImageConverter {
	return ImageConverter{}
}

func (c ImageConverter) Supports(ctx context.Context, srcExt, dstExt string) bool {
	val, ok := imageSupportedExt.Load(srcExt)
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

func (c ImageConverter) Convert(ctx context.Context, srcPath, dstPath, dstExt string) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}
	file, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("imageConverter open src error: %w", err)
	}
	defer file.Close()

	if err := ctx.Err(); err != nil {
		return err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("imageConverter decode error: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("imageConverter create dst error: %w", err)
	}

	defer func() {
		closeErr := dstFile.Close()
		if err == nil {
			err = closeErr
		}
		if ctx.Err() != nil {
			os.Remove(dstPath)
		}
	}()

	writer := ctxWriter{
		ctx: ctx,
		w:   dstFile,
	}

	switch dstExt {
	case "png":
		return png.Encode(writer, img)
	case "jpg", "jpeg":
		return jpeg.Encode(writer, img, nil)
	case "webp":
		return nativewebp.Encode(writer, img, nil)
	case "bmp":
		return bmp.Encode(writer, img)
	case "tiff":
		return tiff.Encode(writer, img, nil)
	case "avif":
		return avif.Encode(writer, img)
	case "ico", "cur":
		b := img.Bounds()
		if b.Dx() > 256 || b.Dy() > 256 {
			img = resizeImage(img, b.Dx(), b.Dy())
		}
		if dstExt == "ico" {
			return ico.Encode(writer, img)
		}
		return cur.Encode(writer, img)
	}
	return nil
}

func resizeImage(src image.Image, srcWidth, srcHeight int) image.Image {
	var newWidth, newHeight int
	if srcWidth > srcHeight {
		newWidth = 256
		newHeight = (srcHeight * 256) / srcWidth
	} else {
		newHeight = 256
		newWidth = (srcWidth * 256) / srcHeight
	}

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}
