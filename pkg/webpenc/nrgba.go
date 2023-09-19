//go:build cgo && webpenc

package webpenc

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
	"io"
	"unsafe"
)

func EncodeNRGBA(w io.Writer, nrgba *image.NRGBA, opt WebPOptions) error {
	// guard against a possible panic if len(nrgba.Pix) < 1
	if len(nrgba.Pix) < 1 {
		return errors.New("error encoding webp file; could not convert source image to nrgba")
	}
	nrgba_pixels := (*C.uint8_t)(&nrgba.Pix[0])
	// allocate C memory for the image. the max size of a WebP file is 4 GiB (1<<33 bytes)
	// but that seems like too much; instead, i'll assume that the encoding won't require
	// more memory than the rgba pixels do. i'll add 1<<8 extra bytes as a cushion.
	// i'm not sure what happens if it overflows this... might be bad!
	// the pointers get a little confusing, but output is a uint8_t**
	// it points to a pointer to the memory allocated by C.malloc
	p := (*C.uint8_t)(C.malloc(C.size_t(len(nrgba.Pix) + 256)))
	defer C.free(unsafe.Pointer(p)) // DO NOT FORGET
	output := &p

	var size C.size_t
	var err error
	switch opt.IsLossy {
	case true:
		size, err = C.encodeLossyRGBA(nrgba_pixels,
			C.int(nrgba.Rect.Max.X),
			C.int(nrgba.Rect.Max.Y),
			C.int(nrgba.Stride),
			C.float(float32(opt.Quality)),
			output)
		if err != nil {
			return err
		}
	case false:
		size, err = C.encodeLosslessRGBA(nrgba_pixels,
			C.int(nrgba.Rect.Max.X),
			C.int(nrgba.Rect.Max.Y),
			C.int(nrgba.Stride),
			output)
		if err != nil {
			return err
		}
	}
	b := C.GoBytes(unsafe.Pointer(p), C.int(size))
	_, err = w.Write(b)
	return err
}
