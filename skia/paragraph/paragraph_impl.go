package paragraph

import (
	"math"
	"sync"

	"github.com/zodimo/go-skia-support/skia/interfaces"
)

// InternalState represents the layout progress of the paragraph.
//
// Ported from: skia-source/modules/skparagraph/src/ParagraphImpl.h
type InternalState int

const (
	// StateUnknown is the initial state.
	StateUnknown InternalState = iota
	// StateIndexed means text properties have been computed.
	StateIndexed
	// StateShaped means text has been shaped into runs.
	StateShaped
	// StateLineBroken means text has been broken into lines.
	StateLineBroken
	// StateFormatted means lines have been formatted/justified.
	StateFormatted
)

// ParagraphImpl is the concrete implementation of the Paragraph interface.
// It orchestrates all other components (Run, OneLineShaper, TextLine, TextWrapper)
// to perform layout and painting.
//
// Ported from: skia-source/modules/skparagraph/src/ParagraphImpl.cpp
type ParagraphImpl struct {
	// Input
	text           string
	paragraphStyle ParagraphStyle
	textStyles     []Block       // styled blocks
	placeholders   []Placeholder // placeholder elements
	fontCollection *FontCollection
	unicode        interfaces.SkUnicode

	// Internal state
	state                     InternalState
	runs                      []*Run
	clusters                  []*Cluster
	clustersIndexFromCodeUnit []int
	codeUnitProperties        []int // unicode flags per code unit
	bidiRegions               []BidiRegion
	lines                     []*TextLine
	words                     []int // word boundary positions

	// UTF mappings (for query methods)
	utf8IndexForUTF16Index []int
	utf16IndexForUTF8Index []int
	utf16MappingOnce       sync.Once

	// Metrics
	width                float32
	height               float32
	maxIntrinsicWidth    float32
	minIntrinsicWidth    float32
	longestLine          float32
	alphabeticBaseline   float32
	ideographicBaseline  float32
	exceededMaxLines     bool
	unresolvedGlyphs     int
	unresolvedCodepoints map[rune]struct{}

	// Caching
	oldWidth             float32
	oldHeight            float32
	maxWidthWithTrailing float32
	emptyMetrics         InternalLineMetrics
	strutMetrics         InternalLineMetrics

	// Text properties
	hasLineBreaks        bool
	hasWhitespacesInside bool
	trailingSpaces       int
}

// NewParagraphImpl creates a new ParagraphImpl with the given parameters.
func NewParagraphImpl(
	text string,
	style ParagraphStyle,
	blocks []Block,
	placeholders []Placeholder,
	fontCollection *FontCollection,
	unicode interfaces.SkUnicode,
) *ParagraphImpl {
	// Ensure we have at least one placeholder (the "fake" placeholder)
	if len(placeholders) == 0 {
		placeholders = []Placeholder{NewPlaceholderDefault()}
	}

	return &ParagraphImpl{
		text:                 text,
		paragraphStyle:       style,
		textStyles:           blocks,
		placeholders:         placeholders,
		fontCollection:       fontCollection,
		unicode:              unicode,
		state:                StateUnknown,
		runs:                 make([]*Run, 0),
		clusters:             make([]*Cluster, 0),
		lines:                make([]*TextLine, 0),
		unresolvedCodepoints: make(map[rune]struct{}),
		emptyMetrics:         NewInternalLineMetrics(),
		strutMetrics:         NewInternalLineMetrics(),
		trailingSpaces:       len(text),
	}
}

// --- TextLineOwner interface implementation ---

// Styles returns all styled blocks.
func (p *ParagraphImpl) Styles() []Block {
	return p.textStyles
}

// Run returns a run by index.
func (p *ParagraphImpl) Run(index int) *Run {
	if index >= 0 && index < len(p.runs) {
		return p.runs[index]
	}
	return nil
}

// Cluster returns a cluster by index.
func (p *ParagraphImpl) Cluster(index int) *Cluster {
	if index >= 0 && index < len(p.clusters) {
		return p.clusters[index]
	}
	return nil
}

// Block returns a block by index.
func (p *ParagraphImpl) Block(index int) Block {
	if index >= 0 && index < len(p.textStyles) {
		return p.textStyles[index]
	}
	return Block{}
}

// GetUnicode returns the Unicode interface.
func (p *ParagraphImpl) GetUnicode() interfaces.SkUnicode {
	return p.unicode
}

// ParagraphStyle returns the paragraph style.
func (p *ParagraphImpl) ParagraphStyle() ParagraphStyle {
	return p.paragraphStyle
}

// FontCollection returns the font collection.
func (p *ParagraphImpl) FontCollection() *FontCollection {
	return p.fontCollection
}

// GetText returns the text.
func (p *ParagraphImpl) GetText() string {
	return p.text
}

// --- Metrics accessors (Paragraph interface) ---

// GetMaxWidth returns the layout width.
func (p *ParagraphImpl) GetMaxWidth() float32 {
	return p.width
}

// GetHeight returns the total paragraph height.
func (p *ParagraphImpl) GetHeight() float32 {
	return p.height
}

// GetMinIntrinsicWidth returns the narrowest width without breaking words.
func (p *ParagraphImpl) GetMinIntrinsicWidth() float32 {
	return p.minIntrinsicWidth
}

// GetMaxIntrinsicWidth returns the width without any wrapping.
func (p *ParagraphImpl) GetMaxIntrinsicWidth() float32 {
	return p.maxIntrinsicWidth
}

// GetAlphabeticBaseline returns the alphabetic baseline.
func (p *ParagraphImpl) GetAlphabeticBaseline() float32 {
	return p.alphabeticBaseline
}

// GetIdeographicBaseline returns the ideographic baseline.
func (p *ParagraphImpl) GetIdeographicBaseline() float32 {
	return p.ideographicBaseline
}

// GetLongestLine returns the width of the longest line.
func (p *ParagraphImpl) GetLongestLine() float32 {
	return p.longestLine
}

// DidExceedMaxLines returns true if max lines was exceeded.
func (p *ParagraphImpl) DidExceedMaxLines() bool {
	return p.exceededMaxLines
}

// LineNumber returns the number of lines.
func (p *ParagraphImpl) LineNumber() int {
	return len(p.lines)
}

// --- State management ---

// MarkDirty marks the paragraph as needing relayout.
func (p *ParagraphImpl) MarkDirty() {
	if p.state > StateIndexed {
		p.state = StateIndexed
	}
}

// UnresolvedGlyphs returns the count of unresolved glyphs, or -1 if not yet shaped.
func (p *ParagraphImpl) UnresolvedGlyphs() int {
	if p.state < StateShaped {
		return -1
	}
	return p.unresolvedGlyphs
}

// UnresolvedCodepoints returns the set of unresolved codepoints.
func (p *ParagraphImpl) UnresolvedCodepoints() []rune {
	result := make([]rune, 0, len(p.unresolvedCodepoints))
	for r := range p.unresolvedCodepoints {
		result = append(result, r)
	}
	return result
}

// --- Internal helpers ---

// text returns a substring for the given range.
func (p *ParagraphImpl) textRange(tr TextRange) string {
	if tr.Start < 0 || tr.End > len(p.text) {
		return ""
	}
	return p.text[tr.Start:tr.End]
}

// clusterIndex returns the cluster index for a text index.
func (p *ParagraphImpl) clusterIndex(textIdx int) int {
	if textIdx >= 0 && textIdx < len(p.clustersIndexFromCodeUnit) {
		return p.clustersIndexFromCodeUnit[textIdx]
	}
	return -1
}

// runByCluster returns the run for a cluster index.
func (p *ParagraphImpl) runByCluster(clusterIdx int) *Run {
	if c := p.Cluster(clusterIdx); c != nil {
		return p.Run(c.RunIndex())
	}
	return nil
}

// strutEnabled returns true if strut is enabled.
func (p *ParagraphImpl) strutEnabled() bool {
	return p.paragraphStyle.StrutStyle.StrutEnabled
}

// strutForceHeight returns true if strut height is forced.
func (p *ParagraphImpl) strutForceHeight() bool {
	return p.paragraphStyle.StrutStyle.ForceStrutHeight
}

// getEllipsis returns the ellipsis string to use.
func (p *ParagraphImpl) getEllipsis() string {
	return p.paragraphStyle.Ellipsis
}

// codeUnitHasProperty checks if a code unit has the specified property.
func (p *ParagraphImpl) codeUnitHasProperty(index int, property int) bool {
	if index >= 0 && index < len(p.codeUnitProperties) {
		return (p.codeUnitProperties[index] & property) == property
	}
	return false
}

// littleRound rounds a value for Flutter test compatibility.
func littleRound(a float32) float32 {
	val := float64(a)
	if val < 0 {
		val = -val
	}
	if val < 10000 {
		return float32(math.Round(float64(a)*100.0) / 100.0)
	} else if val < 100000 {
		return float32(math.Round(float64(a)*10.0) / 10.0)
	}
	return float32(math.Floor(float64(a)))
}

// --- TextWrapperOwner interface implementation ---

// Clusters returns all clusters.
func (p *ParagraphImpl) Clusters() []*Cluster {
	return p.clusters
}

// StructForceHeight returns whether strut height is forced.
func (p *ParagraphImpl) StructForceHeight() bool {
	return p.strutForceHeight()
}

// GetApplyRoundingHack returns whether to apply rounding hack.
func (p *ParagraphImpl) GetApplyRoundingHack() bool {
	return p.paragraphStyle.ApplyRoundingHack
}

// GetEmptyMetrics returns the empty paragraph metrics.
func (p *ParagraphImpl) GetEmptyMetrics() InternalLineMetrics {
	return p.emptyMetrics
}

// StrutEnabled returns whether strut is enabled.
func (p *ParagraphImpl) StrutEnabled() bool {
	return p.strutEnabled()
}

// StrutMetrics returns the strut metrics.
func (p *ParagraphImpl) StrutMetrics() InternalLineMetrics {
	return p.strutMetrics
}

// Lines returns all lines.
func (p *ParagraphImpl) Lines() []*TextLine {
	return p.lines
}

// Text returns the paragraph text.
func (p *ParagraphImpl) Text() string {
	return p.text
}

// TODO: Interface checks will be added when all methods are implemented
// var _ Paragraph = (*ParagraphImpl)(nil)
var _ TextLineOwner = (*ParagraphImpl)(nil)
var _ TextWrapperOwner = (*ParagraphImpl)(nil)
