<h2>About</h2> 
<p>Imgconv is a CLI tool for basic image manipulation. It can be used to convert jpeg, gif, png, tiff, and webp files to jpeg, gif, png, or tiff files. It can also be used to rescale images. Imgconv is powered mainly by Go's standard image library, and it doesn't require the use of cgo. Encoding webp files is currently disabled by default. (Imgconv's experimental support for webp encoding <i>does</i> require cgo.)</p>

<h2>How to Use</h2>
<p>To begin, install this package using the Go compiler:
<pre>go install github.com/cdillond/imgconv@latest</pre>
</p>
<p>When running imgconv, the following parameters are mandatory:
<pre><code>-mode string [REQUIRED] local, remote, or dir
-to string [REQUIRED] the file format of the output image; gif, jpeg, png, and tiff are supported
-url string [REQUIRED] the url of the source image or, if -mode=dir, the path of the target directory</code></pre>
A complete list of accepted parameters can be found in the Flags section.
</p>
<h2>mode</h2>
<p>The <code>-mode</code> flag accepts three possible values: local, remote, and dir.
<ul>
<li>local - imgconv parses the local file specified by the <code>-url</code> flag.</li>
<li>remote - imgconv downloads the resource specified by the <code>-url</code> flag. Only https and http schemes are allowed.</li>
<li>dir - imgconv parses all files in the local directory specified by the <code>-url</code> flag. By default, all subdirectories and their contents are ignored.</li>
</ul>
</p>
<h2>Flags</h2>
The following flags are accepted:
<table>
<tr><th>Flag</th><th>Type</th><th>Usage</th><th>Default</th></tr>
<tr><td><code>-allowUpsize</code></td><td><code>string</code></td><td>permit image pixel dimensions to increase when resizing</td><td><code>false</code></td></tr>
<tr><td><code>-dstDir</code></td><td><code>string</code></td><td>the path of the destination directory; if not specified, the current working directory will be used</td><td>current working directory</td></tr>
<tr><td><code>-gifNumColors</code></td><td><code>uint</code></td><td>the maximum number of colors in output gif files; accepted values are 1-256</td><td><code>256</code></td></tr>
<tr><td><code>-height</code></td><td><code>int</code></td><td>height of the output image in pixels; does not preserve the proportions of the source image</td><td></td></tr>
<tr><td><code>-interpolator</code></td><td><code>string</code></td><td>the interpolation algorithm used to resample images; options are CatmullRom (low speed/high quality), NearestNeighbor (high speed/low quality), and ApproxBiLinear (medium speed/medium quality)</td><td><code>CatmullRom</code></td></tr>
<tr><td><code>-jpegQual</code></td><td><code>uint</code></td><td>the image quality of output jpeg files; accepted values are 0-100 (low - high)</td><td><code>100</code></td></tr>
<tr><td><code>-maxProcs</code></td><td><code>uint</code></td><td>the maximum number of files that can be processed in parallel in dir mode</td><td><code>10</code></td></tr>
<tr><td><code>-maxSidePixels</code></td><td><code>int</code></td><td>size of the greatest dimension of the output image rectangle in pixels; preserves the proportions of the source image</td><td></td></tr>
<tr><td><code>-minSidePixels</code></td><td><code>int</code></td><td>size of the smallest dimension of the output image rectangle in pixels; preserves the proportions of the source image</td><td></td></tr>
<tr><td><code>-mode</code></td><td><code>string</code></td><td><b>[REQUIRED]</b> local, remote, or dir</td><td></td></tr>
<tr><td><code>-out</code></td><td><code>string</code></td><td> the path of the output file; if not specified, the source file name (with an updated extension) will be used (see docs for exceptions); if the path is absolute, it overrides dstDir, but, otherwise, it is relative to dstDir (if specified) or the current working directory; cannot be used in dir mode</td><td></td></tr>
<tr><td><code>-recursive</code></td><td><code>bool</code></td><td>if <code>true</code> and <code>-mode=dir</code>, imgconv will parse all files in the target directory, including all subdirectories</td><td><code>false</code></td></tr>
<tr><td><code>-scaleToHeight</code></td><td><code>int</code></td><td>size of the output image height in pixels; preserves the proportions of the source image</td><td></td></tr>
<tr><td><code>-scaleToWidth</code></td><td><code>int</code></td><td>size of the output image width in pixels; preserves the proportions of the source image</td><td></td></tr>
<tr><td><code>-to</code></td><td><code>string</code></td><td><b>[REQUIRED]</b> the file format of the output image; gif, jpeg, png, and tiff are supported</td><td></td></tr>
<tr><td><code>-url</code></td><td><code>string</code></td><td><b>[REQUIRED]</b> the url of the source image or, if <code>-mode=dir</code>, the path of the target directory</td><td></td></tr>
<tr><td><code>-webpLossy</code></td><td><code>bool</code></td><td>if <code>true</code>, lossy compression will be used for webp encoding</td><td><code>false</code></td></tr>
<tr><td><code>-webpQual</code></td><td><code>uint</code></td><td>the image quality of output webp files when <code>-webpLossy=true</code>; accepted values are 0-100 (low - high)</td><td><code>100</code></td></tr>
<tr><td><code>-width</code></td><td><code>int</code></td><td>width of the output image in pixels; does not preserve the proportions of the source image</td><td></td></tr>
</table>
<h2>Special cases</h2>
<p>Certain flags cannot be used in all cases. Be aware of the following restrictions:
<ul>
<li><code>-out</code> cannot be used in dir mode.</li> 
<li>If <code>-out</code> is an absolute path, it overrides <code>-dstDir</code>.</li>
<li><code>-maxProcs</code> should only be used in dir mode.</li>
<li><code>-recursive</code> should only be used in dir mode.</li>
<li>Only one of <code>-maxSidePixels</code> and <code>-minSidePixels</code> should be specified at a time. If values for both flags are provided, only <code>-maxSidePixels</code> will be used.
<li>Only one of <code>-scaleToHeight</code> and <code>-scaleToWidth</code> should be specified at a time. If values for both flags are provided, only <code>-scaleToHeight</code> will be used.
<li>At most one of <code>-maxSidePixels</code>, <code>-minSidePixels</code>, <code>-scaleToHeight</code>, and <code>-scaleToWidth</code> should be specified at a time. If multiple values are provided anyway, <code>-scaleToHeight</code> and <code>-scaleToWidth</code> override <code>-maxSidePixels</code> and <code>-minSidePixels</code>.</li>
<li>If at least one of <code>-height</code> and <code>-width</code> is specified, any values provided for <code>-maxSidePixels</code>, <code>-minSidePixels</code>, <code>-scaleToHeight</code>, or <code>-scaleToWidth</code> are ignored.</li>
<li><code>-webpLossy</code> and <code>-webpQual</code> are only available if webp encoding is explicitly enabled at build time.</li>
</ul>
</p>
<h2>Naming procedure</h2>
<p>By default, files will be saved under the same name as the source file, with the appropriate file extension. This behavior can be modified using the <code>-out</code> flag. Important exceptions include:
<ol>
<li>File names that contain periods may be truncated, and may not be parsed correctly.</li>
<li>If an output file name cannot be parsed, a name will be assigned at random.</li>
<li>If an output file name conflicts with an existing file, "_v" and a version number will be appended to the end of the new file name, unless the file name is specified by the <code>-out</code> flag, in which case the new file will replace the existing one.</li>
</ol></p>
<h2>Enabling webp encoding</h2>
<p>Imgconv provides <i>experimental</i> support for webp encoding via bindings to Google's <a href="https://developers.google.com/speed/webp/download">libwebp</a> C library. To use this feature, libwebp must be installed in a standard location and cgo must be enabled. When building imgconv, include <code>webpenc</code> as a build tag. On Debian-based Linux systems, for example, this can be achieved using the following commands:
<pre>$ sudo apt install libwebp-dev<br>$ git clone https://github.com/cdillond/imgconv<br>$ cd imgconv<br>$ go env -w CGO_ENABLED=1<br>$ go build -tags webpenc ./cmd/imgconv</pre>
This solution is suboptimal, and setting it up might be more hassle than it is worth. It has only been tested on Linux and Windows.
</p>
