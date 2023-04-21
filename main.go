package main

import (
	"flag"
	"fmt"
)

func main() {
	mode := flag.String("mode", "", "[REQUIRED] local, remote, or dir")
	srcUrl := flag.String("url", "", "[REQUIRED] the url of the source image OR path of the target directory (if mode=dir)")
	toFileType := flag.String("to", "", "[REQUIRED] the file format of the output image; gif, jpeg, png, and tiff are supported")
	dstDir := flag.String("dstDir", "", "the path of the destination directory; if not specified, the current working directory will be used")
	dstFileName := flag.String("out", "", "the path of the output file; if not specified, the source file name (with an updated extension) will be used (the destination file will be assigned a random name if the source file name cannot be parsed); if the path is absolute, it overrides dstDir, but, otherwise, it is relative to dstDir (if specified) or the current working directory; local and remote mode only")
	resample := flag.Bool("resample", false, "whether to resample the source image; if true, additional height/width parameter(s) must also be specified")
	maxSidePixels := flag.Int("maxSidePixels", -1, "size of the greatest dimension of the destination image rectangle in pixels when rescaled; preserves the proportions of the source image")
	minSidePixels := flag.Int("minSidePixels", -1, "size of the smallest dimension of the destination image rectangle in pixels when rescaled; preserves the proportions of the source image")
	scaleToHeight := flag.Int("scaleToHeight", -1, "length of the destination image height in pixels; the width of the destination image will be scaled to retain the source image's original proportions")
	scaleToWidth := flag.Int("scaleToWidth", -1, "length of the destination image width in pixels; the height of the destination image will be scaled to retain the source image's original proportions")
	height := flag.Int("height", -1, "height of the output image in pixels; the source image's proportions will not be preserved")
	width := flag.Int("width", -1, "width of the output image in pixels when resized; the source image's proportions will not be preserved")
	allowUpsize := flag.Bool("allowUpsize", false, "permit image pixel dimensions to increase when resizing")
	jpegQuality := flag.Uint("jpegQual", 100, "specify the image quality of output jpeg files; accepted values are 0-100 (low - high)")
	gifNumColors := flag.Uint("gifNumColors", 256, "specify the number of colors in output gif files; accepted values are 1-256")
	interpolator := flag.String("interpolator", "", "specify the interpolation algorithm used to resample images; options are CatmullRom (default, slow/high quality), NearestNeighbor (fast/poor quality), and ApproxBiLinear (medium/medium)")
	recursive := flag.Bool("recursive", false, "if true and dirMode=true, imgconv will parse all files in the target directory, including all subdirectories")
	maxProcs := flag.Uint("maxProcs", 10, "the maximum number of concurrent files that can be processed in dirMode (min: 1, default: 10); higher may be quicker, but can also lead to greater memory consumption")

	flag.Parse()

	dstFormat := StringToFileType(*toFileType)
	if dstFormat > 4 {
		fmt.Println("unsupported output file format")
		return
	}

	encCfg := NewEncodeCfg(
		dstFormat,
		WithJpegQuality(int(*jpegQuality)),
		WithGifNumColors(int(*gifNumColors)),
	)
	rsmplCfg := NewResampleCfg(
		*resample,
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
	if *resample {
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
