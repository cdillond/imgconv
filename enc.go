package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/tiff"

	"github.com/cdillond/imgconv/pkg/utils"
	"github.com/cdillond/imgconv/pkg/webpenc"
)

type EncodeCfg struct {
	FileType      utils.FileType
	GifNumColors  int
	GifQuantizer  draw.Quantizer
	GifDrawer     draw.Drawer
	JpegQuality   int
	TiffCompType  tiff.CompressionType
	TiffPredictor bool
	WebPLossy     bool
	WebPQuality   uint
}
type EncodeOpt func(*EncodeCfg)

func NewEncodeCfg(fileType utils.FileType, opts ...EncodeOpt) EncodeCfg {
	cfg := EncodeCfg{
		FileType:      fileType,
		GifNumColors:  256,
		GifQuantizer:  nil,
		GifDrawer:     nil,
		JpegQuality:   100,
		TiffCompType:  0,
		TiffPredictor: false,
		WebPLossy:     false,
		WebPQuality:   100,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// I think these constraints are enforced by the image package anyway, but might as well be safe...
func WithJpegQuality(n int) func(*EncodeCfg) {
	return func(e *EncodeCfg) {
		if n < 0 {
			e.JpegQuality = 0
		} else if n > 100 {
			e.JpegQuality = 100
		} else {
			e.JpegQuality = n
		}

	}
}

func WithGifNumColors(n int) func(*EncodeCfg) {
	return func(e *EncodeCfg) {
		if n < 0 {
			e.GifNumColors = 1
		} else if n > 256 {
			e.GifNumColors = 256
		} else {
			e.GifNumColors = n
		}
	}
}

func WithWebPLossy(l bool) func(*EncodeCfg) {
	return func(e *EncodeCfg) {
		e.WebPLossy = l
	}
}

func WithWebPQual(u uint) func(*EncodeCfg) {
	return func(e *EncodeCfg) {
		if u > 100 {
			e.WebPQuality = 100
		} else {
			e.WebPQuality = u
		}
	}
}

/*
TO DO
func WithGifQuantizer(q draw.Quantizer) func(*EncodeCfg) {
	return func(e *EncodeCfg) {
		e.GifQuantizer = q
	}
}
*/

func Encode(img image.Image, w io.Writer, cfg EncodeCfg) error {
	switch cfg.FileType {
	case utils.GIF:
		return gif.Encode(w, img, &gif.Options{
			NumColors: cfg.GifNumColors,
			Quantizer: cfg.GifQuantizer,
			Drawer:    cfg.GifDrawer})
	case utils.JPEG:
		return jpeg.Encode(w, img, &jpeg.Options{Quality: cfg.JpegQuality})
	case utils.PNG:
		return png.Encode(w, img)
	case utils.TIFF:
		return tiff.Encode(w, img, &tiff.Options{
			Compression: cfg.TiffCompType,
			Predictor:   cfg.TiffPredictor})
	case utils.WEBP:
		return webpenc.EncodeWebP(w, img, webpenc.WebPOptions{IsLossy: cfg.WebPLossy, Quality: cfg.WebPQuality})
	default:
		return fmt.Errorf("unsupported file type")
	}
}

/*
	TODO
	implement webp encoding in Go
*/
