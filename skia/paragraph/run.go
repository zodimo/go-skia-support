// Package paragraph provides the SkParagraph text layout module types.
//
// This file contains the Run struct which represents a contiguous sequence of glyphs
// that share the same font, style, and script. It is the fundamental building block
// for storing shaped text output from the shaper.
//
// Ported from: skia-source/modules/skparagraph/src/Run.cpp
package paragraph

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
	"github.com/zodimo/go-skia-support/skia/shaper"
)

// --- Type Aliases ---

// RunIndex is an index into a slice of runs.
type RunIndex = int

// ClusterIndex is an index into a slice of clusters.
type ClusterIndex = int

// ClusterRange is a range of cluster indices.
type ClusterRange = Range[int]

// GlyphIndex is an index into a slice of glyphs.
type GlyphIndex = int

// GlyphRange is a range of glyph indices.
type GlyphRange = Range[int]

// --- Constants ---

// EmptyRun is a sentinel value indicating no run.
const EmptyRun = EmptyIndex

// EmptyCluster is a sentinel value indicating no cluster.
const EmptyCluster = EmptyIndex

// EmptyClusters is a sentinel value indicating no cluster range.
var EmptyClusters = EmptyRange

// --- Run Struct ---

// Run represents a contiguous sequence of glyphs that share the same font,
// style, and script. It stores the shaped output from the shaper and provides
// methods for accessing glyph positions, metrics, and cluster mappings.
//
// Ported from: skia-source/modules/skparagraph/src/Run.h
type Run struct {
	// --- Core Data ---
	textRange    TextRange          // UTF-8 text range this run covers
	clusterRange ClusterRange       // cluster indices (populated during layout)
	font         interfaces.SkFont  // font used for this run
	fontMetrics  models.FontMetrics // cached font metrics

	// --- Glyph Data (shaped output) ---
	glyphs         []uint16       // glyph IDs
	positions      []models.Point // glyph positions (len = glyphCount + 1 for trailing advance)
	offsets        []models.Point // glyph offsets (for kerning/adjustment)
	clusterIndexes []uint32       // cluster index (text offset) for each glyph

	// --- Layout Properties ---
	advance      models.Point // total advance of the run (width, height)
	offset       models.Point // offset from line origin
	clusterStart int          // first char index in paragraph
	utf8Range    shaper.Range // original UTF-8 range from shaper

	// --- Style/Script Info ---
	bidiLevel uint8  // BiDi embedding level (even=LTR, odd=RTL)
	script    uint32 // script tag (e.g., 'Latn')
	language  string // language tag (e.g., "en")

	// --- Metrics (computed) ---
	heightMultiplier float32 // height multiplier for line spacing
	useHalfLeading   bool    // whether to use half-leading model
	baselineShift    float32 // vertical baseline adjustment
	correctAscent    float32 // adjusted ascent (with leading/multiplier)
	correctDescent   float32 // adjusted descent
	correctLeading   float32 // adjusted leading

	// --- Run State ---
	index            int  // run index in paragraph
	isEllipsis       bool // whether this is an ellipsis run
	placeholderIndex int  // placeholder index, or MaxInt if not placeholder

	// --- Justification ---
	justificationShifts []models.Point // (current, prev) shifts for justification
}

// NewRun creates a new Run from shaper RunInfo.
//
// Parameters:
//   - info: The RunInfo from the shaper containing font, metrics, and glyph count
//   - firstChar: The first character index in the paragraph
//   - heightMultiplier: Height multiplier for line spacing (0 means use default)
//   - useHalfLeading: Whether to use half-leading model for metrics
//   - baselineShift: Vertical baseline adjustment
//   - index: The run index in the paragraph
//   - offsetX: Initial X offset for the run
func NewRun(
	info shaper.RunInfo,
	firstChar int,
	heightMultiplier float32,
	useHalfLeading bool,
	baselineShift float32,
	index int,
	offsetX float32,
) *Run {
	glyphCount := int(info.GlyphCount)

	r := &Run{
		textRange:    NewTextRange(firstChar+info.Utf8Range.Begin, firstChar+info.Utf8Range.End),
		clusterRange: EmptyClusters,
		font:         info.Font,

		glyphs:         make([]uint16, glyphCount),
		positions:      make([]models.Point, glyphCount+1),
		offsets:        make([]models.Point, glyphCount+1),
		clusterIndexes: make([]uint32, glyphCount+1),

		advance:      models.Point{X: info.Advance.X, Y: info.Advance.Y},
		offset:       models.Point{X: models.Scalar(offsetX), Y: 0},
		clusterStart: firstChar,
		utf8Range:    info.Utf8Range,

		bidiLevel: info.BidiLevel,
		script:    info.Script,
		language:  info.Language,

		heightMultiplier: heightMultiplier,
		useHalfLeading:   useHalfLeading,
		baselineShift:    baselineShift,

		index:            index,
		isEllipsis:       false,
		placeholderIndex: math.MaxInt,
	}

	// Get font metrics
	r.fontMetrics = getFontMetrics(info.Font)

	// Calculate adjusted metrics
	r.calculateMetrics()

	// Set trailing position and cluster index (edge case handling)
	r.positions[glyphCount] = models.Point{
		X: r.offset.X + r.advance.X,
		Y: r.offset.Y + r.advance.Y,
	}
	r.offsets[glyphCount] = models.Point{X: 0, Y: 0}
	if r.LeftToRight() {
		r.clusterIndexes[glyphCount] = uint32(info.Utf8Range.End)
	} else {
		r.clusterIndexes[glyphCount] = uint32(info.Utf8Range.Begin)
	}

	return r
}

// getFontMetrics extracts FontMetrics from an SkFont.
// This is a helper function since SkFont interface may not directly expose GetMetrics.
func getFontMetrics(font interfaces.SkFont) models.FontMetrics {
	// The SkFont interface doesn't have a GetMetrics method yet.
	// We'll use a heuristic based on font size for now.
	// TODO: Add GetMetrics to SkFont interface when available.
	size := font.Size()
	return models.FontMetrics{
		Ascent:  -size * 0.8, // typical ascent ratio
		Descent: size * 0.2,  // typical descent ratio
		Leading: size * 0.05, // typical leading ratio
	}
}

// calculateMetrics computes the correct ascent, descent, and leading
// based on height multiplier and half-leading settings.
//
// Ported from: Run::calculateMetrics() in Run.cpp
func (r *Run) calculateMetrics() {
	// Start with font metrics, applying half-leading split
	r.correctAscent = float32(r.fontMetrics.Ascent) - float32(r.fontMetrics.Leading)*0.5
	r.correctDescent = float32(r.fontMetrics.Descent) + float32(r.fontMetrics.Leading)*0.5
	r.correctLeading = 0

	// If height multiplier is near zero, use default metrics
	if nearlyZero(r.heightMultiplier) {
		return
	}

	// Calculate the target run height
	runHeight := r.heightMultiplier * float32(r.font.Size())
	fontIntrinsicHeight := r.correctDescent - r.correctAscent

	if r.useHalfLeading {
		// Half-leading model: split extra space evenly above and below
		extraLeading := (runHeight - fontIntrinsicHeight) / 2
		r.correctAscent -= extraLeading
		r.correctDescent += extraLeading
	} else {
		// Scale ascent and descent proportionally
		multiplier := runHeight / fontIntrinsicHeight
		r.correctAscent *= multiplier
		r.correctDescent *= multiplier
	}

	// Apply baseline shift
	r.correctAscent += r.baselineShift
	r.correctDescent += r.baselineShift
}

// nearlyZero checks if a scalar value is close to zero.
func nearlyZero(x float32) bool {
	const epsilon = 1e-6
	return x > -epsilon && x < epsilon
}

// --- Accessor Methods ---

// NewRunBuffer returns a shaper.Buffer for the shaper to fill with glyph data.
// This allows the shaper to write directly into the Run's storage.
func (r *Run) NewRunBuffer() shaper.Buffer {
	return shaper.Buffer{
		Glyphs:    r.glyphs,
		Positions: r.positions,
		Clusters:  r.clusterIndexes,
		Point:     r.offset,
	}
}

// PosX returns the X position at the given glyph index.
func (r *Run) PosX(index int) float32 {
	return float32(r.positions[index].X)
}

// PosY returns the Y position at the given glyph index.
func (r *Run) PosY(index int) float32 {
	return float32(r.positions[index].Y)
}

// AddX adds a shift to the X position at the given glyph index.
func (r *Run) AddX(index int, shift float32) {
	r.positions[index].X += models.Scalar(shift)
}

// Size returns the number of glyphs in the run.
func (r *Run) Size() int {
	return len(r.glyphs)
}

// SetWidth sets the advance width of the run.
func (r *Run) SetWidth(width float32) {
	r.advance.X = models.Scalar(width)
}

// SetHeight sets the advance height of the run.
func (r *Run) SetHeight(height float32) {
	r.advance.Y = models.Scalar(height)
}

// Shift moves the run offset by the given delta.
func (r *Run) Shift(shiftX, shiftY float32) {
	r.offset.X += models.Scalar(shiftX)
	r.offset.Y += models.Scalar(shiftY)
}

// Advance returns the total advance of the run.
func (r *Run) Advance() models.Point {
	return models.Point{
		X: r.advance.X,
		Y: r.fontMetrics.Descent - r.fontMetrics.Ascent + r.fontMetrics.Leading,
	}
}

// Offset returns the run offset from line origin.
func (r *Run) Offset() models.Point {
	return r.offset
}

// Ascent returns the font ascent plus baseline shift.
func (r *Run) Ascent() float32 {
	return float32(r.fontMetrics.Ascent) + r.baselineShift
}

// Descent returns the font descent plus baseline shift.
func (r *Run) Descent() float32 {
	return float32(r.fontMetrics.Descent) + r.baselineShift
}

// Leading returns the font leading.
func (r *Run) Leading() float32 {
	return float32(r.fontMetrics.Leading)
}

// CorrectAscent returns the adjusted ascent (with height multiplier applied).
func (r *Run) CorrectAscent() float32 {
	return r.correctAscent + r.baselineShift
}

// CorrectDescent returns the adjusted descent (with height multiplier applied).
func (r *Run) CorrectDescent() float32 {
	return r.correctDescent + r.baselineShift
}

// CorrectLeading returns the adjusted leading.
func (r *Run) CorrectLeading() float32 {
	return r.correctLeading
}

// Font returns the font used for this run.
func (r *Run) Font() interfaces.SkFont {
	return r.font
}

// LeftToRight returns true if the run is left-to-right.
func (r *Run) LeftToRight() bool {
	return r.bidiLevel%2 == 0
}

// TextDirection returns the text direction based on bidi level.
func (r *Run) TextDirection() TextDirection {
	if r.LeftToRight() {
		return TextDirectionLTR
	}
	return TextDirectionRTL
}

// Index returns the run index in the paragraph.
func (r *Run) Index() int {
	return r.index
}

// HeightMultiplier returns the height multiplier.
func (r *Run) HeightMultiplier() float32 {
	return r.heightMultiplier
}

// UseHalfLeading returns whether half-leading model is used.
func (r *Run) UseHalfLeading() bool {
	return r.useHalfLeading
}

// BaselineShift returns the baseline shift.
func (r *Run) BaselineShift() float32 {
	return r.baselineShift
}

// IsPlaceholder returns true if this run is a placeholder.
func (r *Run) IsPlaceholder() bool {
	return r.placeholderIndex != math.MaxInt
}

// IsEllipsis returns true if this run is an ellipsis.
func (r *Run) IsEllipsis() bool {
	return r.isEllipsis
}

// TextRange returns the text range covered by this run.
func (r *Run) TextRange() TextRange {
	return r.textRange
}

// ClusterRange returns the cluster range for this run.
func (r *Run) ClusterRange() ClusterRange {
	return r.clusterRange
}

// SetClusterRange sets the cluster range for this run.
func (r *Run) SetClusterRange(from, to int) {
	r.clusterRange = NewRange(from, to)
}

// ClusterIndex returns the cluster index at the given glyph position.
func (r *Run) ClusterIndex(pos int) int {
	return int(r.clusterIndexes[pos])
}

// GlobalClusterIndex returns the global cluster index at the given glyph position.
func (r *Run) GlobalClusterIndex(pos int) int {
	return r.clusterStart + int(r.clusterIndexes[pos])
}

// PositionX returns the X position at pos, including justification shifts.
func (r *Run) PositionX(pos int) float32 {
	x := r.PosX(pos)
	if len(r.justificationShifts) > pos {
		x += float32(r.justificationShifts[pos].Y) // Y holds cumulative shift
	}
	return x
}

// Glyphs returns the slice of glyph IDs.
func (r *Run) Glyphs() []uint16 {
	return r.glyphs
}

// Positions returns the slice of glyph positions.
func (r *Run) Positions() []models.Point {
	return r.positions
}

// Offsets returns the slice of glyph offsets.
func (r *Run) Offsets() []models.Point {
	return r.offsets
}

// ClusterIndexes returns the slice of cluster indexes.
func (r *Run) ClusterIndexes() []uint32 {
	return r.clusterIndexes
}

// Script returns the script tag for this run.
func (r *Run) Script() uint32 {
	return r.script
}

// Language returns the language tag for this run.
func (r *Run) Language() string {
	return r.language
}

// BidiLevel returns the bidi level for this run.
func (r *Run) BidiLevel() uint8 {
	return r.bidiLevel
}

// --- Justification ---

// ResetJustificationShifts clears the justification shifts.
func (r *Run) ResetJustificationShifts() {
	r.justificationShifts = nil
}

// --- Script Detection ---

// IsCursiveScript returns true if this run uses a cursive script
// (Arabic, Syriac, etc.) where letter spacing should not be applied.
//
// Ported from: Run::isCursiveScript() in Run.cpp
func (r *Run) IsCursiveScript() bool {
	switch r.script {
	case makeFourByteTag('A', 'r', 'a', 'b'): // ARABIC
		return true
	case makeFourByteTag('R', 'o', 'h', 'g'): // HANIFI_ROHINGYA
		return true
	case makeFourByteTag('M', 'a', 'n', 'd'): // MANDAIC
		return true
	case makeFourByteTag('M', 'o', 'n', 'g'): // MONGOLIAN
		return true
	case makeFourByteTag('N', 'k', 'o', 'o'): // NKO
		return true
	case makeFourByteTag('P', 'h', 'a', 'g'): // PHAGS_PA
		return true
	case makeFourByteTag('S', 'y', 'r', 'c'): // SYRIAC
		return true
	default:
		return false
	}
}

// makeFourByteTag creates a four-byte tag from four characters.
func makeFourByteTag(a, b, c, d byte) uint32 {
	return uint32(a)<<24 | uint32(b)<<16 | uint32(c)<<8 | uint32(d)
}

// IsResolved returns true if all glyphs in the run are resolved (non-zero).
func (r *Run) IsResolved() bool {
	for _, glyph := range r.glyphs {
		if glyph == 0 {
			return false
		}
	}
	return true
}

// --- Width Calculation ---

// CalculateWidth calculates the width of a glyph range [start, end).
// If clip is true, clips to the actual advance.
func (r *Run) CalculateWidth(start, end int, clip bool) float32 {
	if start >= end || start < 0 || end > len(r.positions) {
		return 0
	}
	width := float32(r.positions[end].X - r.positions[start].X)
	if clip {
		// Handle justification shifts if present
		if len(r.justificationShifts) > 0 && end > 0 {
			width += float32(r.justificationShifts[end-1].X)
		}
	}
	return width
}

// CalculateHeight calculates the height using the specified metric styles.
func (r *Run) CalculateHeight(ascentStyle, descentStyle LineMetricStyle) float32 {
	var ascent, descent float32

	if ascentStyle == LineMetricStyleTypographic {
		ascent = r.Ascent()
	} else {
		ascent = r.CorrectAscent()
	}

	if descentStyle == LineMetricStyleTypographic {
		descent = r.Descent()
	} else {
		descent = r.CorrectDescent()
	}

	return descent - ascent
}

// Clip returns the bounding rectangle of the run.
func (r *Run) Clip() models.Rect {
	return models.Rect{
		Left:   r.offset.X,
		Top:    r.offset.Y,
		Right:  r.offset.X + r.advance.X,
		Bottom: r.offset.Y + r.advance.Y,
	}
}

// UpdateMetrics updates line metrics based on placeholder style.
// Used for placeholder runs to update the line metrics according to alignment.
func (r *Run) UpdateMetrics(metrics *InternalLineMetrics) {
	if !r.IsPlaceholder() {
		return
	}
	// Placeholder runs affect metrics based on their baseline alignment.
	// For now, use the run's ascent/descent.
	metrics.AddRun(r)
}

// TextToGlyphRange maps a text range to a glyph range.
// Returns start and end indices into the glyphs/positions slice.
func (r *Run) TextToGlyphRange(textRange TextRange) (int, int) {
	if r.Size() == 0 {
		return 0, 0
	}

	startGlyph := -1
	endGlyph := -1

	// Iterate to find the range of glyphs covered by the text range.
	// Note: glyphs are stored in visual order, so cluster indexes may not be monotonic if mixed (unlikely in single run)
	// but for LTR/RTL they should be monotonic.
	// We scan all glyphs to be safe and simple.
	glyphCount := r.Size()
	for i := 0; i < glyphCount; i++ {
		cluster := int(r.clusterIndexes[i])
		if cluster >= textRange.Start && cluster < textRange.End {
			if startGlyph == -1 {
				startGlyph = i
			}
			endGlyph = i
		}
	}

	if startGlyph == -1 {
		return 0, 0
	}

	// endGlyph is inclusive in loop, make it exclusive for return
	return startGlyph, endGlyph + 1
}
