//go:build cgo && webpenc

package main

/*
	#cgo LDFLAGS: -lwebp
	#include <webp/encode.h>
	#include <stdio.h>

	void encodeImg();
 	void encodeImg()
  	{
		printf("%d\n",WEBP_HINT_PICTURE,);
 		return;
 	}
*/
import "C"

import (
	"errors"
)

func EncodeWebp() error {
	C.encodeImg()
	return errors.New("not implemented yet")
}
