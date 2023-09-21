package main

import (
	"fmt"

	"net/url"
	"strings"
)

func ParseDataUrl(u *url.URL) (string, error) {
	opaque := u.Opaque
	_, after, found := strings.Cut(opaque, ";")
	if !found {
		return *new(string), fmt.Errorf("unable to parse data url")
	}

	before, after, found := strings.Cut(after, ",")
	if !found {
		return *new(string), fmt.Errorf("unable to parse data url image format")
	}
	if before != "base64" {
		return *new(string), fmt.Errorf("unable to parse data url encoding")
	}
	return after, nil
}
