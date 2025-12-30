package models

import "github.com/zodimo/go-skia-support/skia/enums"

// SamplingOptions describes how to sample an image (filter, mipmap, cubic, etc.)
// Matches C++ SkSamplingOptions
type SamplingOptions struct {
	UseCubic   bool
	CubicB     float32
	CubicC     float32
	FilterMode enums.FilterMode
	MipmapMode enums.MipmapMode
}

func NewSamplingOptions(filter enums.FilterMode) SamplingOptions {
	return SamplingOptions{
		FilterMode: filter,
		MipmapMode: enums.MipmapModeNone,
	}
}

func NewSamplingOptionsMipmap(filter enums.FilterMode, mipmap enums.MipmapMode) SamplingOptions {
	return SamplingOptions{
		FilterMode: filter,
		MipmapMode: mipmap,
	}
}
