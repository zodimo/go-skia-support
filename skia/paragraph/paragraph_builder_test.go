package paragraph

import (
	"strings"
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

	para := builder.Build()
	if para == nil {
		t.Fatal("Build() returned nil, expected ParagraphImpl")
	}

	// Verify it's a ParagraphImpl
	_, ok := para.(*ParagraphImpl)
	if !ok {
		t.Error("Build() did not return a *ParagraphImpl")
	}
}

func TestParagraphBuilder_Build_Empty(t *testing.T) {
	style := NewParagraphStyle()
	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	// Build without adding any text
	para := builder.Build()
	if para == nil {
		t.Fatal("Build() returned nil for empty paragraph")
	}
}

// --- Placeholder Tests ---

func TestParagraphBuilder_AddPlaceholder(t *testing.T) {
	style := NewParagraphStyle()
	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	builder.AddText("Hello ")

	// Add a placeholder
	phStyle := NewPlaceholderStyle()
	phStyle.Width = 50
	phStyle.Height = 20
	builder.AddPlaceholder(phStyle)

	builder.AddText(" World")

	// Text should contain the placeholder character (U+FFFC)
	text := builder.GetText()
	if !strings.Contains(text, "\uFFFC") {
		t.Error("Text should contain placeholder marker character U+FFFC")
	}

	// Build and verify
	para := builder.Build()
	if para == nil {
		t.Fatal("Build() returned nil")
	}
}

func TestParagraphBuilder_MultiplePlaceholders(t *testing.T) {
	style := NewParagraphStyle()
	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	// Add text with multiple placeholders
	builder.AddText("Start ")
	builder.AddPlaceholder(PlaceholderStyle{Width: 10, Height: 10})
	builder.AddText(" middle ")
	builder.AddPlaceholder(PlaceholderStyle{Width: 20, Height: 20})
	builder.AddText(" end")

	text := builder.GetText()
	count := strings.Count(text, "\uFFFC")
	if count != 2 {
		t.Errorf("Expected 2 placeholder markers, got %d", count)
	}

	para := builder.Build()
	if para == nil {
		t.Fatal("Build() returned nil")
	}
}

func TestParagraphBuilder_PlaceholderAtStart(t *testing.T) {
	style := NewParagraphStyle()
	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	// Start with a placeholder
	builder.AddPlaceholder(PlaceholderStyle{Width: 30, Height: 30})
	builder.AddText("Text after placeholder")

	text := builder.GetText()
	if !strings.HasPrefix(text, "\uFFFC") {
		t.Error("Text should start with placeholder marker")
	}

	para := builder.Build()
	if para == nil {
		t.Fatal("Build() returned nil")
	}
}

func TestParagraphBuilder_PlaceholderAtEnd(t *testing.T) {
	style := NewParagraphStyle()
	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	builder.AddText("Text before placeholder")
	builder.AddPlaceholder(PlaceholderStyle{Width: 30, Height: 30})

	text := builder.GetText()
	if !strings.HasSuffix(text, "\uFFFC") {
		t.Error("Text should end with placeholder marker")
	}

	para := builder.Build()
	if para == nil {
		t.Fatal("Build() returned nil")
	}
}

func TestParagraphBuilder_PlaceholderOnly(t *testing.T) {
	style := NewParagraphStyle()
	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	// Only placeholders, no text
	builder.AddPlaceholder(PlaceholderStyle{Width: 50, Height: 50})
	builder.AddPlaceholder(PlaceholderStyle{Width: 30, Height: 30})

	text := builder.GetText()
	if text != "\uFFFC\uFFFC" {
		t.Errorf("Expected two placeholder markers, got %q", text)
	}

	para := builder.Build()
	if para == nil {
		t.Fatal("Build() returned nil")
	}
}

func TestParagraphBuilder_ResetClearsPlaceholders(t *testing.T) {
	style := NewParagraphStyle()
	fc := NewFontCollection()
	builder := MakeParagraphBuilder(style, fc)

	builder.AddText("Hello")
	builder.AddPlaceholder(PlaceholderStyle{Width: 10, Height: 10})

	builder.Reset()

	// After reset, text should be empty
	if builder.GetText() != "" {
		t.Error("Text should be empty after Reset")
	}

	// Build after reset should work
	para := builder.Build()
	if para == nil {
		t.Fatal("Build() returned nil after Reset")
	}
}
