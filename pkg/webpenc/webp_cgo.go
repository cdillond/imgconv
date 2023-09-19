//go:build cgo && webpenc

package webpenc

import (
	"image"
	"image/draw"
	"io"

	"github.com/cdillond/imgconv/pkg/utils"
)

const MAX_ENCODE_TYPE utils.FileType = 4

func EncodeWebP(w io.Writer, img image.Image, opt WebPOptions) error {
	// Lossless webp is NRGBA, lossy webp is (N)YCbCr(A)
	// libwebp's NRGBA interface is simple; its YCbCr interface is less so
	switch v := img.(type) {
	case *image.NRGBA:
		return EncodeNRGBA(w, v, opt)
	case *image.YCbCr:
		if opt.IsLossy && v.SubsampleRatio == image.YCbCrSubsampleRatio420 {
			return EncodeYCbCrLossy(w, v, opt)
		}
		// redraw YCbCr to 	NRGBA for lossless and non 420 subsample lossy
		nrgba := image.NewNRGBA(img.Bounds())
		draw.Draw(nrgba, nrgba.Bounds(), img, img.Bounds().Min, draw.Src)
		return EncodeNRGBA(w, nrgba, opt)
	case *image.NYCbCrA:
		if opt.IsLossy && v.YCbCr.SubsampleRatio == image.YCbCrSubsampleRatio420 {
			return EncodeNYCbCrALossy(w, v, opt)
		}
		// redraw NYCbCrA to NRGBA for lossless and non 420 subsample lossy
		nrgba := image.NewNRGBA(img.Bounds())
		draw.Draw(nrgba, nrgba.Bounds(), img, img.Bounds().Min, draw.Src)
		return EncodeNRGBA(w, nrgba, opt)
	default:
		// it's easiest just to draw the image to an NRGBA, plus NYCbCrA doesn't implement draw.Image
		nrgba := image.NewNRGBA(img.Bounds())
		draw.Draw(nrgba, nrgba.Bounds(), img, img.Bounds().Min, draw.Src)
		return EncodeNRGBA(w, nrgba, opt)
	}
}
