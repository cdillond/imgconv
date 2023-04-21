package main

import (
	"image"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func SaveFile(img image.Image, dstPath string, encCfg EncodeCfg) error {
	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return Encode(img, f, encCfg)
}

func GetDstFilePath(dstFileName, dstDir, srcUrl string, isRemote bool, fileType FileType) (string, error) {
	var dstName string
	if dstFileName != "" {
		// if -out is an absolute file path, no further action is needed
		if filepath.IsAbs(dstFileName) {
			return dstFileName, nil
		}

		// if -out is not aboslute, check if -dstDir is specified
		if dstDir != "" {
			// if so, join -dstDir and -out, and return as an absolute path
			dstName = filepath.Join(dstDir, dstFileName)
			return filepath.Abs(dstName)
		} else {
			// if no, join cwd and -out
			return filepath.Abs(dstFileName)
		}
	}
	// if -out is not specified...
	// parse current file name and concat it with an appropriate extension
	var srcFileNameExt, dstFileNameExt string
	if isRemote {
		u, _ := url.Parse(srcUrl) // this should already have succeeded, so no error this time
		path := u.Path
		srcFileNameExt = filepath.Base(path)
	} else {
		srcFileNameExt = filepath.Base(srcUrl)
	}
	// does not allow periods in file names
	srcNameExtSlice := strings.Split(srcFileNameExt, ".")
	if len(srcNameExtSlice) < 1 {
		// assign random name
		dstFileNameExt = uuid.NewString() + "." + FileTypeToString(fileType)
	} else {
		dstFileNameExt = srcNameExtSlice[0] + "." + FileTypeToString(fileType)
	}

	if dstDir != "" {
		// if -dstDir is specified but -out is not
		dstName = filepath.Join(dstDir, dstFileNameExt)
		return filepath.Abs(dstName)
	} else {
		// if neither -dstDir nor -out is specified
		return filepath.Abs(dstFileNameExt)
	}

}
