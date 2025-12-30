package models

import "unsafe"

// Pixmap pairs ImageInfo with pixels and rowBytes.
// Matches C++ SkPixmap
type Pixmap struct {
	Info     ImageInfo
	Addr     unsafe.Pointer // using unsafe.Pointer to match raw C++ pointer semantics, though []byte is safer in Go
	RowBytes int
}

// Note: In a pure Go port, we might prefer []byte for pixel storage.
// However, to strictly model the Skia API which often wraps raw memory, unsafe.Pointer allows integration with C/C++ memory.
// For Go-native usage, helper methods using []byte should be added.

func NewPixmap(info ImageInfo, addr unsafe.Pointer, rowBytes int) Pixmap {
	return Pixmap{
		Info:     info,
		Addr:     addr,
		RowBytes: rowBytes,
	}
}

func (p *Pixmap) Reset() {
	p.Info = ImageInfo{}
	p.Addr = nil
	p.RowBytes = 0
}
