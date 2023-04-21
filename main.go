package main

import (
	"flag"
	"fmt"
)

func main() {
	mode := flag.String("mode", "", "[REQUIRED] local, remote, or dir")
	srcUrl := flag.String("url", "", "[REQUIRED] the url of the source image or, if mode=dir, the path of the target directory")
	toFileType := flag.String("to", "", "[REQUIRED] the file format of the output image; gif, jpeg, png, and tiff are supported")
	dstDir := flag.String("dstDir", "", "the path of the destination directory; if not specified, the current working directory will be used")
	dstFileName := flag.String("out", "", "the path of the output file; if not specified, the source file name (with an updated extension) will be used (see docs for exceptions); if the path is absolute, it overrides dstDir, but, otherwise, it is relative to dstDir (if specified) or the current working directory; cannot be used in dir mode")
	//resample := flag.Bool("resample", false, "whether to resample the source image; if true, additional height/width parameter(s) must also be specified")
	maxSidePixels := flag.Int("maxSidePixels", -1, "size of the greatest dimension of the output image rectangle in pixels; preserves the proportions of the source image")
	minSidePixels := flag.Int("minSidePixels", -1, "size of the smallest dimension of the output image rectangle in pixels; preserves the proportions of the source image")
	scaleToHeight := flag.Int("scaleToHeight", -1, "size of the output image height in pixels; preserves the proportions of the source image")
	scaleToWidth := flag.Int("scaleToWidth", -1, "size of the output image width in pixels; preserves the proportions of the source image")
	height := flag.Int("height", -1, "height of the output image in pixels; does not preserve the proportions of the source image")
	width := flag.Int("width", -1, "width of the output image in pixels; does not preserve the proportions of the source image")
	allowUpsize := flag.Bool("allowUpsize", false, "permit image pixel dimensions to increase when resizing")
	jpegQuality := flag.Uint("jpegQual", 100, "the image quality of output jpeg files; accepted values are 0-100 (low - high)")
	gifNumColors := flag.Uint("gifNumColors", 256, "the maximum number of colors in output gif files; accepted values are 1-256")
	interpolator := flag.String("interpolator", "", "the interpolation algorithm used to resample images; options are CatmullRom (default, low speed/high quality), NearestNeighbor (high speed/low quality), and ApproxBiLinear (medium speed/medium quality)")
	recursive := flag.Bool("recursive", false, "if true and mode=dir, imgconv will parse all files in the target directory, including all subdirectories")
	maxProcs := flag.Uint("maxProcs", 10, "the maximum number of files that can be processed in parallel in dir mode")

	flag.Parse()

	dstFormat := StringToFileType(*toFileType)
	if dstFormat > 3 {
		fmt.Println("unsupported output file format")
		return
	}

	encCfg := NewEncodeCfg(
		dstFormat,
		WithJpegQuality(int(*jpegQuality)),
		WithGifNumColors(int(*gifNumColors)),
	)
	rsmplCfg := NewResampleCfg(
		WithAllowUpsize(*allowUpsize),
		WithMinOrMaxSidePixels(*minSidePixels, *maxSidePixels),
		WithScaleToHeightOrWidth(*scaleToHeight, *scaleToWidth),
		WithHeightAndOrWidth(*height, *width),
		WithAlgorithm(*interpolator),
	)
	var b []byte
	var t FileType
	var err error

	switch *mode {
	case "dir":
		err = ProcessDir(*srcUrl, *dstDir, *maxProcs, *recursive, encCfg, rsmplCfg)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	case "local":
		b, t, err = GetBytesAndFileTypeLocal(*srcUrl)
	case "remote":
		b, t, err = GetBytesAndFileTypeRemote(*srcUrl)
	default:
		fmt.Println("mode flag is required")
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	img, err := BytesToImage(b, t)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if rsmplCfg.IsUsed {
		img = Rescale(img, rsmplCfg)
	}
	dstPath, err := GetDstFilePath(*dstFileName, *dstDir, *srcUrl, *mode == "remote", dstFormat)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = SaveFile(img, dstPath, encCfg)
	if err != nil {
		fmt.Println(err.Error())
	}
}
