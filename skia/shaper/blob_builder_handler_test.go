package shaper_test

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/models"
	"github.com/zodimo/go-skia-support/skia/shaper"
)

func TestTextBlobBuilderRunHandler_Interface(t *testing.T) {
	// Verify that TextBlobBuilderRunHandler implements RunHandler
	var _ shaper.RunHandler = (*shaper.TextBlobBuilderRunHandler)(nil)
}

func TestNewTextBlobBuilderRunHandler(t *testing.T) {
	offset := models.Point{X: 10, Y: 20}
	handler := shaper.NewTextBlobBuilderRunHandler("Hello", offset)

	if handler == nil {
		t.Fatal("NewTextBlobBuilderRunHandler returned nil")
	}

	// EndPoint should return the initial offset
	endPoint := handler.EndPoint()
	if endPoint.X != 10 || endPoint.Y != 20 {
		t.Errorf("EndPoint() = %v, want {10, 20}", endPoint)
	}
}

func TestTextBlobBuilderRunHandler_BeginLine(t *testing.T) {
	offset := models.Point{X: 5, Y: 15}
	handler := shaper.NewTextBlobBuilderRunHandler("Test", offset)

	// BeginLine should reset current position to offset
	handler.BeginLine()

	// After beginLine, the endpoint should still be the offset
	// (endpoint tracks offset, not current position)
	endPoint := handler.EndPoint()
	if endPoint.X != 5 || endPoint.Y != 15 {
		t.Errorf("After BeginLine, EndPoint() = %v, want {5, 15}", endPoint)
	}
}

func TestTextBlobBuilderRunHandler_EmptyMakeBlob(t *testing.T) {
	handler := shaper.NewTextBlobBuilderRunHandler("", models.Point{X: 0, Y: 0})

	// MakeBlob without any runs should return nil
	blob := handler.MakeBlob()
	if blob != nil {
		t.Error("MakeBlob() should return nil when no runs added")
	}
}

func TestTextBlobBuilderRunHandler_RunInfoWithNilFont(t *testing.T) {
	handler := shaper.NewTextBlobBuilderRunHandler("Test", models.Point{X: 0, Y: 0})
	handler.BeginLine()

	// RunInfo with nil font should not panic
	info := shaper.RunInfo{
		Font:       nil,
		BidiLevel:  0,
		Advance:    models.Point{X: 100, Y: 0},
		GlyphCount: 4,
		Utf8Range:  shaper.Range{Begin: 0, End: 4},
	}
	handler.RunInfo(info) // Should not panic
}

func TestTextBlobBuilderRunHandler_CommitRunInfo(t *testing.T) {
	handler := shaper.NewTextBlobBuilderRunHandler("Test", models.Point{X: 0, Y: 0})
	handler.BeginLine()
	handler.CommitRunInfo()
	// CommitRunInfo should complete without panic
}

func TestTextBlobBuilderRunHandler_RunBufferWithZeroGlyphs(t *testing.T) {
	handler := shaper.NewTextBlobBuilderRunHandler("Test", models.Point{X: 0, Y: 0})
	handler.BeginLine()

	info := shaper.RunInfo{
		Font:       nil,
		GlyphCount: 0, // Zero glyphs
	}
	buffer := handler.RunBuffer(info)

	// With zero glyphs, buffer should be empty
	if len(buffer.Glyphs) != 0 {
		t.Errorf("RunBuffer with 0 glyphs should return empty glyphs slice, got %d", len(buffer.Glyphs))
	}
}

func TestTextBlobBuilderRunHandler_CommitLine(t *testing.T) {
	offset := models.Point{X: 0, Y: 0}
	handler := shaper.NewTextBlobBuilderRunHandler("Test", offset)

	handler.BeginLine()
	// No runs added, metrics are 0, so offset shouldn't change much
	handler.CommitLine()

	// EndPoint should reflect offset update
	// With no runs, metrics are 0, so Y offset should be unchanged
	endPoint := handler.EndPoint()
	if endPoint.X != 0 {
		t.Errorf("EndPoint().X = %v, want 0", endPoint.X)
	}
}

func TestTextBlobBuilderRunHandler_FullWorkflow(t *testing.T) {
	// Test a complete workflow: beginLine -> runInfo -> commitRunInfo ->
	// (runBuffer -> commitRunBuffer) -> commitLine -> makeBlob
	handler := shaper.NewTextBlobBuilderRunHandler("Hello", models.Point{X: 0, Y: 0})

	// Start line
	handler.BeginLine()

	// Provide run info (nil font for this test since we don't have a mock)
	info := shaper.RunInfo{
		Font:       nil,
		BidiLevel:  0,
		Advance:    models.Point{X: 50, Y: 0},
		GlyphCount: 5,
		Utf8Range:  shaper.Range{Begin: 0, End: 5},
	}
	handler.RunInfo(info)

	// Commit run info
	handler.CommitRunInfo()

	// Get run buffer - will be empty since font is nil and AllocRunPos will fail
	buffer := handler.RunBuffer(info)
	if buffer.Glyphs == nil && info.Font == nil {
		// Expected: with nil font, RunBuffer returns empty buffer
		// This is fine for testing the workflow
	}

	// Commit run buffer
	handler.CommitRunBuffer(info)

	// Commit line
	handler.CommitLine()

	// Make blob - may be nil since we couldn't allocate runs without a font
	_ = handler.MakeBlob()
}
