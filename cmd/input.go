package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"imgconv/pkg/utils"
)

func GetBytesAndFileTypeLocal(u string) ([]byte, utils.FileType, error) {
	var b []byte
	srcPath, err := filepath.Abs(u)
	if err != nil {
		return b, 6, err
	}
	srcExt := filepath.Ext(srcPath)

	f, err := os.Open(srcPath)
	if err != nil {
		return b, 6, err
	}
	defer f.Close()
	b, err = io.ReadAll(f)
	return b, utils.StringToFileType(srcExt), err
}

func GetBytesAndFileTypeRemote(u string) ([]byte, utils.FileType, error) {
	var b []byte
	srcUrl, err := url.Parse(u)
	if err != nil {
		return b, 6, err
	}
	var mimeType string
	// infer scheme if not provided
	if srcUrl.Scheme == "" {
		schemes := []string{"https", "http"}

		for _, scheme := range schemes {
			srcUrl.Scheme = scheme
			resp, err := http.Get(srcUrl.String())
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			mimeType = resp.Header.Get("Content-Type")
			b, err = io.ReadAll(resp.Body)
			return b, utils.StringToFileType(mimeType), err
		}
		return b, 6, errors.New("could not infer scheme from incomplete url: " + u)
	}
	if srcUrl.Scheme != "https" && srcUrl.Scheme != "http" {
		return b, 6, errors.New("unsupported url scheme: " + srcUrl.Scheme + ". use https or http instead.")
	}
	resp, err := http.Get(srcUrl.String())
	if err != nil {
		return b, 6, err
	}
	defer resp.Body.Close()
	mimeType = resp.Header.Get("Content-Type")
	b, err = io.ReadAll(resp.Body)
	return b, utils.StringToFileType(mimeType), err
}
