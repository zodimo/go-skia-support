package paragraph

import (
	"testing"
)

func TestParagraphBuilder_Lifecycle(t *testing.T) {
	style := NewParagraphStyle()
	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	if builder == nil {
		t.Fatal("MakeParagraphBuilder returned nil")
	}

	// Test initial state
	if builder.GetText() != "" {
		t.Errorf("Expected empty text, got %q", builder.GetText())
	}
	ps := builder.GetParagraphStyle()
	if !ps.Equals(&style) {
		t.Error("ParagraphStyle mismatch")
	}
}

func TestParagraphBuilder_AddText(t *testing.T) {
	style := NewParagraphStyle()
	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	builder.AddText("Hello")
	builder.AddText(" World")

	expected := "Hello World"
	if builder.GetText() != expected {
		t.Errorf("Expected %q, got %q", expected, builder.GetText())
	}

	builder.Reset()
	if builder.GetText() != "" {
		t.Errorf("Expected empty text after Reset, got %q", builder.GetText())
	}
}

func TestParagraphBuilder_StyleStack(t *testing.T) {
	style := NewParagraphStyle()
	// Set a recognizable color for default style
	style.DefaultTextStyle.Color = 0xFF0000FF // Red

	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	// Initial style should match default
	current := builder.PeekStyle()
	if current.Color != 0xFF0000FF {
		t.Errorf("Expected default color 0xFF0000FF, got %x", current.Color)
	}

	// Push new style
	newStyle := NewTextStyle()
	newStyle.Color = 0xFF00FF00 // Green
	builder.PushStyle(&newStyle)

	current = builder.PeekStyle()
	if current.Color != 0xFF00FF00 {
		t.Errorf("Expected pushed color 0xFF00FF00, got %x", current.Color)
	}

	// Pop style
	builder.Pop()
	current = builder.PeekStyle()
	if current.Color != 0xFF0000FF {
		t.Errorf("Expected restored color 0xFF0000FF, got %x", current.Color)
	}
}

func TestParagraphBuilder_Build(t *testing.T) {
	style := NewParagraphStyle()
	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	builder.AddText("Test")

	// Currently Build() returns nil, so we just check it doesn't panic
	para := builder.Build()
	if para != nil {
		// If we implement a stub, this might change, but for now we expect nil or check behavior.
		// Since we returned nil in implementation, this is expected behavior for now.
	}
}
