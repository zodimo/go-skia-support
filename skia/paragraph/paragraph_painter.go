package paragraph

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

// DashPathEffect defines a dash effect for stroking paths.
//
// Ported from: skia-source/modules/skparagraph/include/ParagraphPainter.h
type DashPathEffect struct {
	OnLength  float32
	OffLength float32
}

// DecorationStyle defines how to stroke or fill a decoration.
//
// Ported from: skia-source/modules/skparagraph/include/ParagraphPainter.h
type DecorationStyle struct {
	Color          models.Color4f
	StrokeWidth    float32
	DashPathEffect *DashPathEffect // Optional
}

// ParagraphPainter is access to a Skia canvas to draw paragraph contents.
// This interface allows abstracting the drawing commands.
//
// Ported from: skia-source/modules/skparagraph/include/ParagraphPainter.h
type ParagraphPainter interface {
	// DrawTextBlob draws a text blob at the given coordinates.
	DrawTextBlob(blob interfaces.SkTextBlob, x, y float32, paint interfaces.SkPaint)

	// DrawTextShadow draws a shadow for the text blob.
	DrawTextShadow(blob interfaces.SkTextBlob, x, y float32, color models.Color4f, blurSigma float64)

	// DrawRect draws a rectangle with the given paint.
	DrawRect(rect models.Rect, paint interfaces.SkPaint)

	// DrawFilledRect draws a filled rectangle with the given decoration style.
	DrawFilledRect(rect models.Rect, style DecorationStyle)

	// DrawPath draws a path with the given decoration style.
	DrawPath(path interfaces.SkPath, style DecorationStyle)

	// DrawLine draws a line with the given decoration style.
	DrawLine(x0, y0, x1, y1 float32, style DecorationStyle)

	// ClipRect adds a rectangle to the clip.
	ClipRect(rect models.Rect)

	// Translate applies a translation to the current transform.
	Translate(dx, dy float32)

	// Save saves the current state of the canvas.
	Save()

	// Restore restores the state of the canvas.
	Restore()
}
