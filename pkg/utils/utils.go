package utils

import (
	"strings"
)

type FileType uint

const (
	GIF FileType = iota
	JPEG
	PNG
	TIFF
	WEBP
	UNSUPPORTED
)

func FileTypeToString(f FileType) string {
	switch f {
	case 0:
		return "gif"
	case 1:
		return "jpeg"
	case 2:
		return "png"
	case 3:
		return "tiff"
	case 4:
		return "webp"
	default:
		return "unsupported"
	}
}

// RETURNS utils.UNSUPPORTED IF FILETYPE IS NOT VALID
func StringToFileType(s string) FileType {
	switch strings.ToLower(s) {
	case "gif", "image/gif":
		return GIF
	case "jpeg", "jpg", "image/jpeg":
		return JPEG
	case "png", "image/png":
		return PNG
	case "tiff", "image/tiff":
		return TIFF
	case "webp", "image/webp":
		return WEBP
	default:
		return UNSUPPORTED
	}
}
