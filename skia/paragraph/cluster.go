package paragraph

import (
	"github.com/zodimo/go-skia-support/skia/interfaces"
)

// BreakType represents the type of break associated with a cluster.
type BreakType int

const (
	BreakNone BreakType = iota
	BreakGrapheme
	BreakSoftLine
	BreakHardLine
)

// Cluster represents a cluster of glyphs that should be treated as a unit.
//
// Ported from: skia-source/modules/skparagraph/src/Run.h (Cluster class)
type Cluster struct {
	owner     TextLineOwner // Changed from ParagraphImpl to Owner interface
	runIndex  int
	textRange TextRange

	start             int
	end               int
	width             float32
	height            float32
	halfLetterSpacing float32

	isWhitespaceBreak bool
	isIntraWordBreak  bool
	isHardBreak       bool
	isIdeographic     bool
}

// NewCluster creates a new Cluster.
func NewCluster(owner TextLineOwner, runIndex int, start, end int, textRange TextRange, width, height float32) *Cluster {
	return &Cluster{
		owner:     owner,
		runIndex:  runIndex,
		start:     start,
		end:       end,
		textRange: textRange,
		width:     width,
		height:    height,
	}
}

// SetOwner sets the owner.
func (c *Cluster) SetOwner(owner TextLineOwner) {
	c.owner = owner
}

// Run returns the run associated with this cluster.
func (c *Cluster) Run() *Run {
	if c.owner == nil {
		return nil
	}
	return c.owner.Run(c.runIndex)
}

// Font returns the font used by this cluster.
func (c *Cluster) Font() interfaces.SkFont {
	run := c.Run()
	if run != nil {
		return run.Font()
	}
	return nil
}

// Size returns the number of glyphs in the cluster.
func (c *Cluster) Size() int {
	return c.end - c.start
}

// StartPos returns the starting glyph index.
func (c *Cluster) StartPos() int {
	return c.start
}

// EndPos returns the ending glyph index.
func (c *Cluster) EndPos() int {
	return c.end
}

// Width returns the width of the cluster.
func (c *Cluster) Width() float32 {
	return c.width
}

// Height returns the height of the cluster.
func (c *Cluster) Height() float32 {
	return c.height
}

// TextRange returns the text range covered by this cluster.
func (c *Cluster) TextRange() TextRange {
	return c.textRange
}

// RunIndex returns the run index.
func (c *Cluster) RunIndex() int {
	return c.runIndex
}

// IsWhitespaceBreak returns true if this cluster is a whitespace break.
func (c *Cluster) IsWhitespaceBreak() bool {
	return c.isWhitespaceBreak
}

// IsIntraWordBreak returns true if this cluster is an intra-word break.
func (c *Cluster) IsIntraWordBreak() bool {
	return c.isIntraWordBreak
}

// IsHardBreak returns true if this cluster is a hard line break.
func (c *Cluster) IsHardBreak() bool {
	return c.isHardBreak
}

// IsIdeographic returns true if this cluster is ideographic.
func (c *Cluster) IsIdeographic() bool {
	return c.isIdeographic
}

// SetHalfLetterSpacing sets the half letter spacing.
func (c *Cluster) SetHalfLetterSpacing(spacing float32) {
	c.halfLetterSpacing = spacing
}

// GetHalfLetterSpacing returns the half letter spacing.
func (c *Cluster) GetHalfLetterSpacing() float32 {
	return c.halfLetterSpacing
}

// Space adds width to the cluster.
func (c *Cluster) Space(shift float32) {
	c.width += shift
}

// TrimmedWidth returns the width of the cluster trimmed at the given position.
func (c *Cluster) TrimmedWidth(pos int) float32 {
	// TODO: Implement proper trimming logic based on glyphosate positions
	// This requires access to the Run's glyph positions.
	// C++ implementation:
	// return fRun->positionX(pos) - fRun->positionX(fStart);
	run := c.Run()
	if run == nil {
		return 0
	}
	return run.PositionX(pos) - run.PositionX(c.start)
}

// IsSoftBreak returns true if this cluster is a soft line break.
func (c *Cluster) IsSoftBreak() bool {
	if c.owner == nil {
		return false
	}
	text := c.owner.GetText()
	unicode := c.owner.GetUnicode()
	if unicode == nil {
		return false
	}
	return unicode.CodeUnitHasProperty(text, c.textRange.Start, interfaces.CodeUnitFlagSoftLineBreakBefore)
}

// IsGraphemeBreak returns true if this cluster is a grapheme break.
func (c *Cluster) IsGraphemeBreak() bool {
	if c.owner == nil {
		return false
	}
	text := c.owner.GetText()
	unicode := c.owner.GetUnicode()
	if unicode == nil {
		return false
	}
	return unicode.CodeUnitHasProperty(text, c.textRange.Start, interfaces.CodeUnitFlagGraphemeStart)
}

// Contains returns true if the char index is within this cluster.
func (c *Cluster) Contains(ch int) bool {
	return ch >= c.textRange.Start && ch < c.textRange.End
}

// Belongs returns true if this cluster belongs to the given text range.
func (c *Cluster) Belongs(text TextRange) bool {
	return c.textRange.Start >= text.Start && c.textRange.End <= text.End
}

// StartsIn returns true if this cluster starts in the given text range.
func (c *Cluster) StartsIn(text TextRange) bool {
	return c.textRange.Start >= text.Start && c.textRange.Start < text.End
}

// --- Internal ---

// SetBreakType sets the break properties.
func (c *Cluster) SetBreakType(whiteSpace, intraWord, hardBreak, ideographic bool) {
	c.isWhitespaceBreak = whiteSpace
	c.isIntraWordBreak = intraWord
	c.isHardBreak = hardBreak
	c.isIdeographic = ideographic
}
