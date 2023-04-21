package main

import (
	"errors"
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
			dstName = dstFileName
		} else {
			// if -out is not aboslute, check if -dstDir is specified
			// if so, join -dstDir and -out, and check if that path is absolute
			// if yes, no further action is needed
			// if no, join cwd and dstName
			// if -dstDir is not specified, join cwd and -out
			if dstDir != "" {
				dstName = filepath.Join(dstDir, dstFileName)
				dstName, err := filepath.Abs(dstName)
				if err != nil {
					return dstName, err
				}
			} else {
				dstName, err := filepath.Abs(dstFileName)
				return dstName, err
			}
		}
	} else {
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
			dstName, err := filepath.Abs(dstName)
			return dstName, err
		} else {
			// if neither -dstDir nor -out is specified
			dstName, err := filepath.Abs(dstFileNameExt)
			return dstName, err
		}
	}
	return dstName, errors.New("uknown problem determing output file name")
}
