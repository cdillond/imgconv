//go:build cgo && webpenc

package webpenc

/*
	#cgo LDFLAGS: -lwebp
   	#include <webp/encode.h>
   	#include <stdio.h>
   	#include <stdlib.h>
   	#include <errno.h>

uint8_t* encodeLossyYCbCr(float quality, int width, int height, uint8_t* y, uint8_t* u, uint8_t* v, int y_stride, int uv_stride, size_t* size) {

	WebPPicture img;
	WebPConfig cfg;
	WebPMemoryWriter writer;

  	if (!WebPConfigInitInternal(&cfg, 0, quality, WEBP_ENCODER_ABI_VERSION)) {
		return writer.mem;
	}

	if (!WebPPictureInit(&img)) {
    	return writer.mem;
  	}

	img.colorspace = 0;
	img.width = width;
	img.height = height;
	img.y = y;
	img.u = u;
	img.v = v;
	img.y_stride = y_stride;
	img.uv_stride = uv_stride;

	WebPMemoryWriterInit(&writer);

	img.custom_ptr = &writer;
 	img.writer = WebPMemoryWrite;

	if (!WebPEncode(&cfg, &img)) {
		return writer.mem;
	}
	*size = writer.size;
	return writer.mem;
}
*/
import "C"

import (
	"errors"
	"image"
	"io"
	"unsafe"
)

func EncodeYCbCrLossy(w io.Writer, ycbcr *image.YCbCr, opt WebPOptions) error {
	// this only gets used for lossy images
	y := (*C.uint8_t)(&ycbcr.Y[0])
	u := (*C.uint8_t)(&ycbcr.Cb[0])
	v := (*C.uint8_t)(&ycbcr.Cr[0])

	size := (*C.size_t)(C.malloc(C.size_t(1)))
	defer C.free(unsafe.Pointer(size))

	out := C.encodeLossyYCbCr(
		C.float(opt.Quality),
		C.int(ycbcr.Rect.Max.X),
		C.int(ycbcr.Rect.Max.Y),
		y, u, v,
		C.int(ycbcr.YStride),
		C.int(ycbcr.CStride),
		size,
	)
	defer C.free(unsafe.Pointer(out))

	if uint(*size) == 0 {
		return errors.New("could not encode YCbCr to lossy webp")
	}
	b := C.GoBytes(unsafe.Pointer(out), C.int(*size))
	_, err := w.Write(b)
	return err
}
