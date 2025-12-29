package shaper_test

import (
	"testing"

	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
	"github.com/zodimo/go-skia-support/skia/shaper"
)

// mockRunHandler implements shaper.RunHandler for testing.
type mockRunHandler struct{}

func (m *mockRunHandler) BeginLine()                  {}
func (m *mockRunHandler) RunInfo(info shaper.RunInfo) {}
func (m *mockRunHandler) CommitRunInfo()              {}
func (m *mockRunHandler) RunBuffer(info shaper.RunInfo) shaper.Buffer {
	return shaper.Buffer{}
}
func (m *mockRunHandler) CommitRunBuffer(info shaper.RunInfo) {}
func (m *mockRunHandler) CommitLine()                         {}

// mockShaper implements shaper.Shaper for testing.
type mockShaper struct {
	shapeCalled bool
}

func (m *mockShaper) Shape(text string, font interfaces.SkFont, leftToRight bool, width float32, runHandler shaper.RunHandler, features []shaper.Feature) {
	m.shapeCalled = true
}

func (m *mockShaper) ShapeWithIterators(text string,
	fontIter shaper.FontRunIterator,
	bidiIter shaper.BiDiRunIterator,
	scriptIter shaper.ScriptRunIterator,
	langIter shaper.LanguageRunIterator,
	features []shaper.Feature,
	width float32,
	runHandler shaper.RunHandler) {
	m.shapeCalled = true
}

func TestShaperInterface_Shape(t *testing.T) {
	var s shaper.Shaper = &mockShaper{}
	s.Shape("test", nil, true, 0, nil, nil)

	if !s.(*mockShaper).shapeCalled {
		t.Errorf("Shape was not called")
	}
}

func TestInterfaces(t *testing.T) {
	// Verify that we can instantiate structs
	f := shaper.Feature{Tag: 1, Value: 0, Start: 0, End: 1}
	if f.Tag != 1 {
		t.Errorf("Feature tag mismatch")
	}

	p := models.Point{X: 10, Y: 20}
	b := shaper.Buffer{
		Glyphs:    []uint16{0},
		Positions: []models.Point{p},
		Point:     p,
	}
	if len(b.Glyphs) != 1 {
		t.Errorf("Buffer glyphs mismatch")
	}

	// Verify interface satisfaction
	var rh shaper.RunHandler = &mockRunHandler{}
	var s shaper.Shaper = &mockShaper{}

	// simple call
	s.Shape("test", nil, true, 100.0, rh, nil)
}
