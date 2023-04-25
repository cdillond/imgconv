package main

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"

	"imgconv/pkg/utils"
)

func Decode(r io.Reader, t utils.FileType) (image.Image, error) {
	switch t {
	case utils.GIF:
		return gif.Decode(r)
	case utils.JPEG:
		return jpeg.Decode(r)
	case utils.PNG:
		return png.Decode(r)
	case utils.TIFF:
		return tiff.Decode(r)
	case utils.WEBP:
		return webp.Decode(r)
	default:
		var img image.Image
		return img, errors.New("unsupported file type")
	}
}

// tries decoding with likelySrcFmt first, then iterates through all supported file types
func BytesToImage(b []byte, likelySrcFmt utils.FileType) (image.Image, error) {
	if likelySrcFmt < 5 {
		img, err := Decode(bytes.NewReader(b), likelySrcFmt)
		if err == nil {
			return img, nil
		}
	}
	for i := 0; i < 5; i++ {
		img, err := Decode(bytes.NewReader(b), utils.FileType(i))
		if err == nil {
			return img, nil
		}
	}
	return nil, errors.New("unsupported source file format")
}
