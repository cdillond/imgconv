package main

import (
	"errors"
	"net/http"
	"net/url"
	"os"

	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"imgconv/pkg/utils"
)

func DecodeLocal(srcUrl string) (image.Image, utils.FileType, error) {
	f, err := os.Open(srcUrl)
	if err != nil {
		return nil, 6, err
	}
	img, format, err := image.Decode(f)
	f.Close()
	return img, utils.StringToFileType(format), err
}

func DecodeRemote(u string) (image.Image, utils.FileType, error) {
	srcUrl, err := url.Parse(u)
	if err != nil {
		return nil, 6, err
	}

	// infer scheme if not provided
	if srcUrl.Scheme == "" {
		schemes := []string{"https", "http"}
		for _, scheme := range schemes {
			srcUrl.Scheme = scheme
			resp, err := http.Get(srcUrl.String())
			if err != nil {
				continue
			}
			img, format, err := image.Decode(resp.Body)
			resp.Body.Close()
			if err == nil {
				return img, utils.StringToFileType(format), err
			}
		}
		return nil, 6, errors.New("could not infer scheme from incomplete url: " + u)
	}
	if srcUrl.Scheme != "https" && srcUrl.Scheme != "http" {
		return nil, 6, errors.New("unsupported url scheme: " + srcUrl.Scheme + ". use https or http instead.")
	}
	resp, err := http.Get(srcUrl.String())
	if err != nil {
		return nil, 6, err
	}
	img, format, err := image.Decode(resp.Body)
	resp.Body.Close()
	return img, utils.StringToFileType(format), err
}
