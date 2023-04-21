package main

import (
	"image"

	"golang.org/x/image/draw"
)

type ResampleCfg struct {
	IsUsed        bool
	MaxSidePxls   int
	MinSidePxls   int
	ScaleToHeight int
	ScaleToWidth  int
	Width         int
	Height        int
	Algorithm     draw.Interpolator
	AllowUpsize   bool
}

func NewResampleCfg(isUsed bool, opts ...ResampleOpt) ResampleCfg {
	cfg := ResampleCfg{
		IsUsed:        isUsed,
		MaxSidePxls:   -1,
		MinSidePxls:   -1,
		ScaleToHeight: -1,
		ScaleToWidth:  -1,
		Width:         -1,
		Height:        -1,
		Algorithm:     draw.CatmullRom,
		AllowUpsize:   false,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

type ResampleOpt func(*ResampleCfg)

func WithAllowUpsize(allowUpsize bool) func(*ResampleCfg) {
	return func(r *ResampleCfg) {
		r.AllowUpsize = allowUpsize
	}
}

func WithIsUsed(isUsed bool) func(*ResampleCfg) {
	return func(r *ResampleCfg) {
		r.IsUsed = isUsed
	}
}

func WithMinOrMaxSidePixels(min, max int) func(*ResampleCfg) {
	// choose max if both are valid
	if min > 0 && max > 0 {
		return func(r *ResampleCfg) {
			r.MaxSidePxls = max
		}
	}
	if min > 0 {
		return func(r *ResampleCfg) {
			r.MinSidePxls = min
		}
	}
	if max > 0 {
		return func(r *ResampleCfg) {
			r.MaxSidePxls = max
		}
	}
	return func(r *ResampleCfg) {}
}

func WithScaleToHeightOrWidth(h, w int) func(*ResampleCfg) {
	if h > 0 {
		return func(r *ResampleCfg) {
			r.ScaleToHeight = h

			r.MaxSidePxls = -1
			r.MinSidePxls = -1
		}
	}
	if w > 0 {
		return func(r *ResampleCfg) {
			r.ScaleToWidth = w

			r.MaxSidePxls = -1
			r.MinSidePxls = -1
		}
	}
	return func(r *ResampleCfg) {}
}

func WithHeightAndOrWidth(h, w int) func(*ResampleCfg) {
	// overrides MinSidePxls and MaxSidePxls if specified
	if h > 0 || w > 0 {
		return func(r *ResampleCfg) {
			r.MaxSidePxls = -1
			r.MinSidePxls = -1
			r.ScaleToHeight = -1
			r.ScaleToWidth = -1

			r.Height = h
			r.Width = w
		}

	}
	return func(r *ResampleCfg) {}
}

func WithAlgorithm(s string) func(r *ResampleCfg) {
	switch s {
	case "CatmullRom":
		return func(r *ResampleCfg) {
			r.Algorithm = draw.CatmullRom
		}
	case "NearestNeighbor":
		return func(r *ResampleCfg) {
			r.Algorithm = draw.NearestNeighbor
		}
	case "ApproxBiLinear":
		return func(r *ResampleCfg) {
			r.Algorithm = draw.ApproxBiLinear
		}
	default:
		return func(r *ResampleCfg) {}
	}
}

func DstRect(srcRect image.Rectangle, cfg ResampleCfg) image.Rectangle {
	dstRect := srcRect //srcRect.Bounds()
	var done bool
	if cfg.Height > 0 {
		dstRect.Max.Y = cfg.Height
		done = true
	}
	if cfg.Width > 0 {
		dstRect.Max.X = cfg.Width
		done = true
	}
	if done {
		return dstRect
	}

	if cfg.ScaleToHeight > 0 {
		dstRect.Max.X = (srcRect.Max.X * cfg.ScaleToHeight) / srcRect.Max.Y
		dstRect.Max.Y = cfg.ScaleToHeight
		return dstRect
	}

	if cfg.ScaleToWidth > 0 {
		dstRect.Max.X = cfg.ScaleToWidth
		dstRect.Max.Y = (srcRect.Max.Y * cfg.ScaleToWidth) / srcRect.Max.X
		return dstRect
	}

	if cfg.MaxSidePxls > 0 {
		absMax := dstRect.Max.X
		var isMaxY bool
		if dstRect.Max.Y > absMax {
			isMaxY = true
			absMax = dstRect.Max.Y
		}

		if absMax > cfg.MaxSidePxls || cfg.AllowUpsize {
			// resize is needed
			var X, Y int
			if isMaxY {
				X = (srcRect.Max.X * cfg.MaxSidePxls) / srcRect.Max.Y
				Y = cfg.MaxSidePxls
			} else {
				X = cfg.MaxSidePxls
				Y = (srcRect.Max.Y * cfg.MaxSidePxls) / srcRect.Max.X
			}
			dstRect.Max.X = X
			dstRect.Max.Y = Y
		}
		return dstRect
	}

	if cfg.MinSidePxls > 0 {
		// we are assuming that the recangle Min point is always (0,0)
		absMin := srcRect.Max.X
		var isMinY bool
		if absMin > srcRect.Max.Y {
			isMinY = true
			absMin = srcRect.Max.Y
		}
		if absMin < cfg.MinSidePxls || cfg.AllowUpsize {
			var X, Y int
			if isMinY {
				X = (srcRect.Max.X * cfg.MinSidePxls) / srcRect.Max.Y
				Y = cfg.MinSidePxls
			} else {
				X = cfg.MinSidePxls
				Y = (srcRect.Max.Y * cfg.MinSidePxls) / srcRect.Max.X
			}
			dstRect.Max.X = X
			dstRect.Max.Y = Y
		}
	}
	return dstRect
}

func Rescale(src image.Image, cfg ResampleCfg) image.Image {
	dstRect := DstRect(src.Bounds(), cfg)
	dstImg := image.NewNRGBA(dstRect)
	cfg.Algorithm.Scale(dstImg, dstImg.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dstImg
}
