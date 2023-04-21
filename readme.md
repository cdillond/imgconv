imgconv is a simple CLI tool for basic image manipulation, written in Go. It can be used to convert images from jpeg, gif, png, tiff, and webp files to jpeg, gif, png, or tiff files. It can also be used to resize/rescale images. 
There are three modes: remote, local, and dir. Remote mode fetches an image resource from a remote url using an http(s) get request. Local mode loads an image from a local path. Dir mode parses all files in a specified directory. By default, it ignores any subdirectories and their contents, but this behavior can be changed by setting the -recursive flag to true.
The following flags are accepted:
```
 -allowUpsize
        permit image pixel dimensions to increase when resizing
  -dstDir string
        the path of the destination directory; if not specified, the current working directory will be used
  -gifNumColors uint
        specify the number of colors in output gif files; accepted values are 1-256 (default 256)
  -height int
        height of the output image in pixels; the source image's proportions will not be preserved (default -1)
  -interpolator string
        specify the interpolation algorithm used to resample images; options are CatmullRom (default, slow/high quality), NearestNeighbor (fast/poor quality), and ApproxBiLinear (medium/medium)
  -jpegQual uint
        specify the image quality of output jpeg files; accepted values are 0-100 (low - high) (default 100)
  -maxProcs uint
        the maximum number of concurrent files that can be processed in dirMode (min: 1, default: 10); higher may be quicker, but can also lead to greater memory consumption (default 10)
  -maxSidePixels int
        size of the greatest dimension of the destination image rectangle in pixels when rescaled; preserves the proportions of the source image (default -1)
  -minSidePixels int
        size of the smallest dimension of the destination image rectangle in pixels when rescaled; preserves the proportions of the source image (default -1)
  -mode string
        [REQUIRED] local, remote, or dir
  -out string
        the path of the output file; if not specified, the source file name (with an updated extension) will be used (the destination file will be assigned a random name if the source file name cannot be parsed); if the path is absolute, it overrides dstDir, but, otherwise, it is relative to dstDir (if specified) or the current working directory; local and remote mode only
  -recursive
        if true and dirMode=true, imgconv will parse all files in the target directory, including all subdirectories
  -resample
        whether to resample the source image; if true, additional height/width parameter(s) must also be specified
  -scaleToHeight int
        length of the destination image height in pixels; the width of the destination image will be scaled to retain the source image's original proportions (default -1)
  -scaleToWidth int
        length of the destination image width in pixels; the height of the destination image will be scaled to retain the source image's original proportions (default -1)
  -to string
        [REQUIRED] the file format of the output image; gif, jpeg, png, and tiff are supported
  -url string
        [REQUIRED] the url of the source image OR path of the target directory (if mode=dir)
  -width int
        width of the output image in pixels when resized; the source image's proportions will not be preserved (default -1)

```