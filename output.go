package main

import (
	"image"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/cdillond/imgconv/pkg/utils"

	"github.com/google/uuid"
)

func SaveFile(img image.Image, dstPath string, encCfg EncodeCfg) error {
	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	err = Encode(img, f, encCfg)
	if err != nil {
		f.Close()          // can ignore this error
		os.Remove(dstPath) // can ignore the error returned here
		return err
	}
	return f.Close()
}

func GetDstFilePath(dstFileName, dstDir, srcUrl string, isRemote bool, fileType utils.FileType) (string, error) {
	var dstName string
	if dstFileName != "" {
		// TO DO validate dstFileName to avoid possible issues with hidden files, files without extensions, and file names that include multiple periods
		// if -out is an absolute file path, no further action is needed
		if filepath.IsAbs(dstFileName) {
			return dstFileName, nil
		}

		// if -out is not absolute, check if -dstDir is specified
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
		dstFileNameExt = uuid.NewString() + "." + utils.FileTypeToString(fileType)
	} else {
		dstFileNameExt = srcNameExtSlice[0] + "." + utils.FileTypeToString(fileType)
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
