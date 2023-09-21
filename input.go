package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/cdillond/imgconv/pkg/utils"
)

func DecodeLocal(srcUrl string) (image.Image, utils.FileType, error) {
	f, err := os.Open(srcUrl)
	if err != nil {
		return nil, utils.UNSUPPORTED, err
	}
	img, format, err := image.Decode(f)
	if err == nil {
		return img, utils.StringToFileType(format), f.Close()
	}
	f.Close() // otherwise, ignore this error
	return img, utils.StringToFileType(format), err
}

func DecodeRemote(u string) (image.Image, utils.FileType, error) {
	srcUrl, err := url.Parse(u)
	if err != nil {
		return nil, utils.UNSUPPORTED, err
	}

	// prevent client from following redirects
	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	// infer scheme if not provided
	if srcUrl.Scheme == "" {
		schemes := []string{"https", "http"}
		for _, scheme := range schemes {
			srcUrl.Scheme = scheme

			resp, err := client.Get(srcUrl.String())
			if err != nil {
				continue
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				// ignore responses with non-2XX status codes
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				continue
			}
			img, format, err := image.Decode(resp.Body)
			resp.Body.Close()
			if err == nil {
				return img, utils.StringToFileType(format), err
			}
		}
		return nil, utils.UNSUPPORTED, fmt.Errorf("could not infer scheme from incomplete url: %s", u)
	}
	if srcUrl.Scheme == "data" {
		src, err := ParseDataUrl(srcUrl)
		if err != nil {
			return nil, utils.UNSUPPORTED, err
		}
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(src))
		img, format, err := image.Decode(reader)
		if err != nil {
			return img, utils.UNSUPPORTED, err
		}
		return img, utils.StringToFileType(format), ErrDataURL
	}
	if srcUrl.Scheme != "https" && srcUrl.Scheme != "http" {
		return nil, utils.UNSUPPORTED, fmt.Errorf("unsupported url scheme: %s", srcUrl.Scheme)
	}

	resp, err := client.Get(srcUrl.String())
	if err != nil {
		return nil, utils.UNSUPPORTED, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// ignore responses with non-2XX status codes
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		return nil, utils.UNSUPPORTED, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	img, format, err := image.Decode(resp.Body)
	resp.Body.Close()
	return img, utils.StringToFileType(format), err
}
