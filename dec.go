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
)

func Decode(r io.Reader, t FileType) (image.Image, error) {
	switch t {
	case GIF:
		return gif.Decode(r)
	case JPEG:
		return jpeg.Decode(r)
	case PNG:
		return png.Decode(r)
	case TIFF:
		return tiff.Decode(r)
	case WEBP:
		return webp.Decode(r)
	default:
		var img image.Image
		return img, errors.New("unsupported file type")
	}
}

// tries decoding with likelySrcFmt first, then iterates through all supported file types
func BytesToImage(b []byte, likelySrcFmt FileType) (image.Image, error) {
	if likelySrcFmt < 5 {
		img, err := Decode(bytes.NewReader(b), likelySrcFmt)
		if err == nil {
			return img, nil
		}
	}
	for i := 0; i < 5; i++ {
		img, err := Decode(bytes.NewReader(b), FileType(i))
		if err == nil {
			return img, nil
		}
	}
	return nil, errors.New("unsupported source file format")
}
