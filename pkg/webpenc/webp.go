//go:build !cgo || !webpenc

package webpenc

import (
	"errors"
	"image"
	"io"

	"github.com/cdillond/imgconv/pkg/utils"
)

const MAX_ENCODE_TYPE utils.FileType = 3

func EncodeWebP(w io.Writer, img image.Image, opt WebPOptions) error {
	// this code should be unreachable
	return errors.New("webp encoding is not enabled; review docs at github.com/cdillond/imgconv for details")
}
