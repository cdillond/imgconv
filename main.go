package main

import (
	"flag"
	"image"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cdillond/imgconv/pkg/utils"
	"github.com/cdillond/imgconv/pkg/webpenc"
)

func main() {
	mode := flag.String("mode", "", "[REQUIRED] local, remote, or dir")
	srcUrl := flag.String("url", "", "[REQUIRED] the url of the source image or, if -mode=dir, the path of the target directory")
	toFileType := flag.String("to", "", "[REQUIRED] the file format of the output image; gif, jpeg, png, and tiff are supported")
	dstDir := flag.String("dstDir", "", "the path of the destination directory; if not specified, the current working directory will be used")
	dstFileName := flag.String("out", "", "the path of the output file; if not specified, the source file name (with an updated extension) will be used (see docs for exceptions); if the path is absolute, it overrides dstDir, but, otherwise, it is relative to dstDir (if specified) or the current working directory; cannot be used in dir mode")
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
	recursive := flag.Bool("recursive", false, "if true and -mode=dir, imgconv will parse all files in the target directory, including all subdirectories")
	maxProcs := flag.Uint("maxProcs", 10, "the maximum number of files that can be processed in parallel in dir mode")
	webpLossy := flag.Bool("webpLossy", false, "if true, lossy compression will be used for webp encoding")
	webpQuality := flag.Uint("webpQual", 100, "the image quality of output webp files when -webpLossy=true; accepted values are 0-100 (low - high)")

	flag.Parse()

	if *mode == "" || *srcUrl == "" || *toFileType == "" {
		log.Println("mode, url, and to options must be provided")
		flag.Usage()
		os.Exit(1)
	}

	dstFormat := utils.StringToFileType(*toFileType)
	if dstFormat > webpenc.MAX_ENCODE_TYPE { // defined in webp.go and webp_cgo.go
		if dstFormat == utils.UNSUPPORTED {
			log.Fatalln("unsupported output file format")
		}
		if webpenc.MAX_ENCODE_TYPE < utils.WEBP {
			log.Fatalln("webp encoding is not enabled; review the documentation at github.com/cdillond/imgconv for details")
		}
		log.Fatalln("error parsing to file format flag input") // this should not happen
	}

	encCfg := NewEncodeCfg(
		dstFormat,
		WithJpegQuality(int(*jpegQuality)),
		WithGifNumColors(int(*gifNumColors)),
		WithWebPLossy(*webpLossy),
		WithWebPQual(*webpQuality),
	)
	rsmplCfg := NewResampleCfg(
		WithAllowUpsize(*allowUpsize),
		WithRescale(*height, *width, *scaleToHeight, *scaleToWidth, *maxSidePixels, *minSidePixels),
		WithInterpolator(*interpolator),
	)

	var err error
	var img image.Image
	switch *mode {
	case "dir":
		err = ProcessDir(*srcUrl, *dstDir, *maxProcs, *recursive, encCfg, rsmplCfg)
		if err != nil {
			log.Fatalln(err.Error())
		}
		return
	case "local":
		img, _, err = DecodeLocal(*srcUrl)
	case "remote":
		img, _, err = DecodeRemote(*srcUrl)
	default:
		log.Fatalln("a valid mode flag is required")
	}
	if err != nil {
		log.Fatalln(err.Error())
	}

	if rsmplCfg.IsUsed {
		img = Rescale(img, rsmplCfg)
	}
	dstPath, err := GetDstFilePath(*dstFileName, *dstDir, *srcUrl, *mode == "remote", dstFormat)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// check if file path already exists, and adjust name to avoid collisions
	// specifying an output file name via --out overrides this behavior
	if _, err = os.Stat(dstPath); err == nil && *dstFileName == "" {
		var version int
		fdir := filepath.Dir(dstPath)
		fNameExt := filepath.Base(dstPath)
		fNameExtSlice := strings.Split(fNameExt, ".")
		if len(fNameExtSlice) <= 1 {
			// output file names are required to have extensions and should not otherwise include periods
			log.Fatalln("invalid output file path " + dstPath)
		}
		for ; err == nil; version++ {
			dstNameVersionExt := fNameExtSlice[0] + "_v" + strconv.Itoa(version+1) + "." + fNameExtSlice[1]
			dstPath = filepath.Join(fdir, dstNameVersionExt)
			_, err = os.Stat(dstPath)
		}
	}

	err = SaveFile(img, dstPath, encCfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
