//go:build cgo && webpenc

package main

/*
   #cgo LDFLAGS: -lwebp
   #include <webp/encode.h>
   #include <stdio.h>
   #include <stdlib.h>
   #include <errno.h>

   size_t encodeLosslessRGBA(const uint8_t* rgba, int width, int height, int stride, uint8_t** output) {
       return WebPEncodeLosslessRGBA(rgba, width, height, stride, output);
   };

   size_t encodeLossyRGBA(const uint8_t* rgba, int width, int height, int stride, float quality_factor, uint8_t** output) {
		return WebPEncodeRGBA(rgba, width, height, stride, quality_factor, output);
   };
*/
import "C"

import (
	"errors"
	"image"
	"image/draw"
	"io"
	"unsafe"
)

const MAX_ENCODE_TYPE FileType = 4

func EncodeWebP(w io.Writer, img image.Image, o WebPOptions) error {
	// check if already an NRGBA image (n.b. NRGBA = non-premultiplied alpha RGBA)
	rgba, ok := img.(*image.NRGBA)
	if !ok {
		// if not, draw src image to rgba
		rgba = image.NewNRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	}

	// allocate C memory for the image. the max size of a WebP file is 4 GiB (1<<33 bytes)
	// but that's seems like too much; instead i'll allocate 67,108,864 bytes (1<<26).
	// i'm not sure what happens if it overflows this... might be bad!
	// the pointers get a little confusing, but output is a uint8_t**
	// it points to a pointer to the memory allocated by C.malloc
	p := (*C.uint8_t)(C.malloc(1 << 26))
	defer C.free(unsafe.Pointer(p)) // DO NOT FORGET
	output := &p
	// guard against a possible panic if len(rgba.Pix) < 1
	if len(rgba.Pix) < 1 {
		return errors.New("error encoding webp file; could not convert to rgba")
	}
	rgba_pixels := (*C.uint8_t)(&rgba.Pix[0])

	var s C.size_t
	switch o.Lossy {
	case true:
		s, err = C.encodeLossyRGBA(rgba_pixels,
			C.int(rgba.Rect.Max.X),
			C.int(rgba.Rect.Max.Y),
			C.int(rgba.Stride),
			C.float(float32(o.WebPQuality)),
			output)
		if err != nil {
			return err
		}
	case false:
		s, err = C.encodeLosslessRGBA(rgba_pixels,
			C.int(rgba.Rect.Max.X),
			C.int(rgba.Rect.Max.Y),
			C.int(rgba.Stride),
			output)
		if err != nil {
			return err
		}
	}
	b := C.GoBytes(unsafe.Pointer(p), C.int(s))
	_, err = w.Write(b)
	return err
}
