package main

import (
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/tiff"
)

type EncodeCfg struct {
	FileType      FileType
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

func NewEncodeCfg(fileType FileType, opts ...EncodeOpt) EncodeCfg {
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

func WithJpegQuality(n int) func(*EncodeCfg) {
	return func(e *EncodeCfg) {
		e.JpegQuality = n
	}
}

func WithGifNumColors(n int) func(*EncodeCfg) {
	return func(e *EncodeCfg) {
		e.GifNumColors = n
	}
}

func WithWebPLossy(l bool) func(*EncodeCfg) {
	return func(e *EncodeCfg) {
		e.WebPLossy = l
	}
}

func WithWebPQual(u uint) func(*EncodeCfg) {
	return func(e *EncodeCfg) {
		e.WebPQuality = u
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
	case GIF:
		return gif.Encode(w, img, &gif.Options{
			NumColors: cfg.GifNumColors,
			Quantizer: cfg.GifQuantizer,
			Drawer:    cfg.GifDrawer})
	case JPEG:
		return jpeg.Encode(w, img, &jpeg.Options{Quality: cfg.JpegQuality})
	case PNG:
		return png.Encode(w, img)
	case TIFF:
		return tiff.Encode(w, img, &tiff.Options{
			Compression: cfg.TiffCompType,
			Predictor:   cfg.TiffPredictor})
	case WEBP:
		return EncodeWebP(w, img, WebPOptions{cfg.WebPLossy, cfg.WebPQuality})
	default:
		return errors.New("unsupported file type")
	}
}

/*
	TODO
	implement webp encoding
*/
