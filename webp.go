//go:build !cgo || !webpenc

package main

import (
	"errors"
	"image"
	"io"
)

const MAX_ENCODE_TYPE FileType = 3

func EncodeWebP(w io.Writer, img image.Image, opt WebPOptions) error {
	// this code should be unreachable
	return errors.New("webp encoding is not enabled; review docs at github.com/cdillond/imgconv for details")
}
