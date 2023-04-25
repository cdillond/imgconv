//go:build cgo && webpenc

package webpenc

/*
	#cgo LDFLAGS: -lwebp
   	#include <webp/encode.h>
   	#include <stdio.h>
   	#include <stdlib.h>
   	#include <errno.h>

uint8_t* encodeLossyNYCbCrA(float quality, int width, int height, uint8_t* y, uint8_t* u, uint8_t* v, uint8_t* a, int y_stride, int uv_stride, int a_stride, size_t* size) {
	WebPPicture img;
	WebPConfig cfg;
	WebPMemoryWriter writer;

  	if (!WebPConfigInitInternal(&cfg, 4, quality, WEBP_ENCODER_ABI_VERSION)) {
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
	img.a = a;
	img.a_stride = a_stride;

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

func EncodeNYCbCrALossy(w io.Writer, nycbcra *image.NYCbCrA, opt WebPOptions) error {

	y := (*C.uint8_t)(&nycbcra.YCbCr.Y[0])
	u := (*C.uint8_t)(&nycbcra.YCbCr.Cb[0])
	v := (*C.uint8_t)(&nycbcra.YCbCr.Cr[0])
	a := (*C.uint8_t)(&nycbcra.A[0])

	size := (*C.size_t)(C.malloc(C.size_t(1)))
	defer C.free(unsafe.Pointer(size))

	out := C.encodeLossyNYCbCrA(
		C.float(opt.Quality),
		C.int(nycbcra.YCbCr.Rect.Max.X),
		C.int(nycbcra.YCbCr.Rect.Max.Y),
		y, u, v, a,
		C.int(nycbcra.YCbCr.YStride),
		C.int(nycbcra.YCbCr.CStride),
		C.int(nycbcra.AStride),
		size,
	)
	defer C.free(unsafe.Pointer(out))

	if uint(*size) == 0 {
		return errors.New("could not encode NYCbCrA to lossy webp")
	}
	b := C.GoBytes(unsafe.Pointer(out), C.int(*size))
	_, err := w.Write(b)
	return err
}
