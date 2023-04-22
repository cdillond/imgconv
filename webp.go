//go:build !cgo || !webpenc

package main

import (
	"errors"
)

func EncodeWebp() error {
	return errors.New("webp encoding is not enabled; review docs at github.com/cdillond/imgconv for details")
}
