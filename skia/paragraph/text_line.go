package paragraph

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
	"github.com/zodimo/go-skia-support/skia/shaper"
)

// TextLineOwner abstracts the ParagraphImpl for TextLine.
type TextLineOwner interface {
	Styles() []Block
	Run(index int) *Run
	Cluster(index int) *Cluster
	Block(index int) Block
	GetUnicode() interfaces.SkUnicode
	ParagraphStyle() ParagraphStyle
	FontCollection() *FontCollection
	GetText() string
}

// ClipContext contains context for clipping operations.
type ClipContext struct {
	Run            *Run
	Pos            int
	Size           int
	TextShift      float32
	Clip           models.Rect
	TrailingSpaces float32
	ClippingNeeded bool
}

// TextAdjustment specifies how text adjustment should be applied.
type TextAdjustment int

const (
	TextAdjustmentGlyphCluster    TextAdjustment = 1 // All text producing glyphs pointing to the same ClusterIndex
	TextAdjustmentGlyphemeCluster TextAdjustment = 2 // base glyph + all attached diacritics
	TextAdjustmentGrapheme        TextAdjustment = 4 // Text adjusted to graphemes
	TextAdjustmentGraphemeCluster TextAdjustment = 5 // GlyphCluster & Grapheme
)

// TextLine represents a single rendered line.
//
// Ported from: skia-source/modules/skparagraph/src/TextLine.h
type TextLine struct {
	owner                  TextLineOwner
	blockRange             BlockRange
	textExcludingSpaces    TextRange
	text                   TextRange
	textIncludingNewlines  TextRange
	clusterRange           ClusterRange
	ghostClusterRange      ClusterRange
	runsInVisualOrder      []int
	advance                models.Point
	offset                 models.Point
	shift                  float32
	widthWithSpaces        float32
	ellipsis               *Run
	sizes                  InternalLineMetrics
	maxRunMetrics          InternalLineMetrics
	hasBackground          bool
	hasShadows             bool
	hasDecorations         bool
	ascentStyle            LineMetricStyle
	descentStyle           LineMetricStyle
	textBlobCache          []TextBlobRecord
	textBlobCachePopulated bool
}

// NewTextLine creates a new TextLine.
func NewTextLine(
	owner TextLineOwner,
	offset models.Point,
	advance models.Point,
	blocks BlockRange,
	textExcludingSpaces TextRange,
	text TextRange,
	textIncludingNewlines TextRange,
	clusters ClusterRange,
	clustersWithGhosts ClusterRange,
	widthWithSpaces float32,
	sizes InternalLineMetrics,
) *TextLine {
	tl := &TextLine{
		owner:                 owner,
		offset:                offset,
		advance:               advance,
		blockRange:            blocks,
		textExcludingSpaces:   textExcludingSpaces,
		text:                  text,
		textIncludingNewlines: textIncludingNewlines,
		clusterRange:          clusters,
		ghostClusterRange:     clustersWithGhosts,
		widthWithSpaces:       widthWithSpaces,
		sizes:                 sizes,
		ascentStyle:           LineMetricStyleCSS,
		descentStyle:          LineMetricStyleCSS,
	}

	// Reorder visual runs (simplified for now, assumes logical=visual if bidi unimplemented or done in owner)
	// Real implementation needs Bidi reordering
	// Assuming `runsInVisualOrder` needs to be populated.
	// We'll populate it based on clusters.
	if clustersWithGhosts.Width() > 0 {
		start := owner.Cluster(clustersWithGhosts.Start)
		end := owner.Cluster(clustersWithGhosts.End - 1)

		// Collect unique runs in range
		// Using a map to track added runs to preserve order/uniqueness
		runMap := make(map[int]bool)
		for i := start.RunIndex(); i <= end.RunIndex(); i++ {
			if !runMap[i] {
				tl.runsInVisualOrder = append(tl.runsInVisualOrder, i)
				runMap[i] = true
			}
			// Update max run metrics
			run := owner.Run(i)
			tl.maxRunMetrics.AddRun(run)

			// Check flags
			// Need to check style for background/decorations
			// Loop through blocks covering this run??
			// Actually C++ loops through blocks in `fBlockRange`
		}
	}

	// Check styles in block range
	for i := blocks.Start; i < blocks.End; i++ {
		block := owner.Block(i)
		if block.Style.HasBackground {
			tl.hasBackground = true
		}
		if block.Style.Decoration.Type != TextDecorationNone {
			tl.hasDecorations = true
		}
		if len(block.Style.TextShadows) > 0 {
			tl.hasShadows = true
		}
	}

	return tl
}

// Format formats the line based on alignment and width.
func (tl *TextLine) Format(align TextAlign, maxWidth float32) {
	delta := maxWidth - tl.Width()
	if delta <= 0 {
		return
	}

	if align == TextAlignJustify {
		if !tl.isHardBreak() {
			tl.Justify(maxWidth)
		} else if tl.owner.ParagraphStyle().TextDirection == TextDirectionRTL {
			tl.shift = delta
		}
	} else if align == TextAlignRight {
		tl.shift = delta
	} else if align == TextAlignCenter {
		tl.shift = delta / 2
	}
}

// Justify justifies the line to fill the max width.
func (tl *TextLine) Justify(maxWidth float32) {
	// Count whitespace patches
	whitespacePatches := 0
	textLen := float32(0)
	whitespaceLen := float32(0)

	start := tl.clusterRange.Start
	end := tl.clusterRange.End

	for i := start; i < end; i++ {
		cluster := tl.owner.Cluster(i)
		if cluster.IsWhitespaceBreak() {
			whitespacePatches++
			whitespaceLen += cluster.Width()
		} else {
			textLen += cluster.Width()
		}
	}

	if whitespacePatches == 0 {
		return
	}

	step := (maxWidth - textLen) / float32(whitespacePatches)
	totalShift := float32(0)

	for i := start; i < end; i++ {
		cluster := tl.owner.Cluster(i)
		if cluster.IsWhitespaceBreak() {
			totalShift += step - cluster.Width()
			cluster.width = step // update cluster width
		}
		// Shift cluster visually? Run positions?
		// C++ updates fShift or runs?
		// C++ TextLine::justify updates cluster widths and shifts them in paint/scan?
		// No, TextLine::justify updates cluster->space() and sets fWidthWithSpaces.
		// It assumes ScanStyles uses cluster widths.
		// Runs are not moved here, visual offset is calculated during iteration.
	}
}

// isHardBreak returns true if the line ends with a hard break.
func (tl *TextLine) isHardBreak() bool {
	if tl.clusterRange.Width() == 0 {
		return false
	}
	// Check last cluster
	lastCluster := tl.owner.Cluster(tl.clusterRange.End - 1)
	return lastCluster.IsHardBreak()
}

// Width returns the width of the line.
func (tl *TextLine) Width() float32 {
	w := tl.advance.X
	if tl.ellipsis != nil {
		w += tl.ellipsis.Advance().X
	}
	return float32(w)
}

// Height returns the height of the line.
func (tl *TextLine) Height() float32 {
	return float32(tl.advance.Y)
}

// Baseline returns the baseline of the line.
func (tl *TextLine) Baseline() float32 {
	return tl.sizes.Baseline()
}

// ScanStyles iterates styles over the line.
func (tl *TextLine) ScanStyles(styleType StyleType, visitor func(TextRange, TextStyle, ClipContext)) {
	if tl.textExcludingSpaces.Width() == 0 {
		return
	}

	tl.iterateThroughVisualRuns(false, func(run *Run, runOffset float32, textRange TextRange, width *float32) bool {
		*width = tl.iterateThroughSingleRunByStyles(TextAdjustmentGlyphCluster, run, runOffset, textRange, styleType, visitor)
		return true
	})
}

// iterateThroughVisualRuns implements the visitor pattern for runs.
func (tl *TextLine) iterateThroughVisualRuns(includingGhostSpaces bool, visitor func(*Run, float32, TextRange, *float32) bool) {
	currentOffset := float32(0)

	for _, runIndex := range tl.runsInVisualOrder {
		run := tl.owner.Run(runIndex)

		// Calculate text range for this run within the line
		// Intersection of run text range and line text range
		// If includingGhostSpaces, use textIncludingNewlines (or ghost range?)
		// C++ logic:
		lineRange := tl.textExcludingSpaces
		if includingGhostSpaces {
			lineRange = tl.textIncludingNewlines
		}

		runRange := run.TextRange()
		intersection := lineRange.Intersection(runRange)

		if intersection.Width() == 0 {
			continue
		}

		var runWidth float32
		if !visitor(run, currentOffset, intersection, &runWidth) {
			return
		}
		currentOffset += runWidth
	}
}

// iterateThroughSingleRunByStyles iterates styles within a run.
func (tl *TextLine) iterateThroughSingleRunByStyles(
	adj TextAdjustment,
	run *Run,
	runOffset float32,
	textRange TextRange,
	styleType StyleType,
	visitor func(TextRange, TextStyle, ClipContext),
) float32 {
	// Intersection of run text range and line text range is passed as textRange
	// We need to iterate blocks that cover this textRange

	// currentPos tracks progress through textRange (unused in simple iteration)
	// currentPos := textRange.Start
	totalWidth := float32(0)

	for i := tl.blockRange.Start; i < tl.blockRange.End; i++ {
		block := tl.owner.Block(i)
		intersection := block.Range.Intersection(textRange)
		if intersection.Width() == 0 {
			continue
		}

		// If we skipped some text (because blocks are not contiguous? Blocks should cover everything)
		// Assuming blocks are contiguous coverage of line.

		// Measure text
		context := tl.measureTextInsideOneRun(intersection, run, runOffset, 0, false, adj) // simplified args
		visitor(intersection, block.Style, context)                                        // Pass actual style

		// Clip width calculation (Rect has no Width method)
		clipWidth := float32(context.Clip.Right - context.Clip.Left)
		totalWidth += clipWidth
	}
	return totalWidth
}

// measureTextInsideOneRun measures text.
func (tl *TextLine) measureTextInsideOneRun(
	textRange TextRange,
	run *Run,
	runOffsetInLine float32,
	textOffsetInRun float32,
	includeGhostSpaces bool,
	adj TextAdjustment,
) ClipContext {
	startGlyph, endGlyph := run.TextToGlyphRange(textRange)
	if startGlyph == endGlyph {
		// No glyphs in range
		return ClipContext{
			Run:       run,
			Clip:      models.Rect{Left: models.Scalar(runOffsetInLine), Right: models.Scalar(runOffsetInLine)},
			TextShift: runOffsetInLine,
		}
	}

	// Calculate width from positions
	startX := run.PositionX(startGlyph)
	endX := run.PositionX(endGlyph)
	width := endX - startX

	// Handle RTL: if startX > endX, swap?
	// Width should be positive. PositionX accounts for visual order.
	// If LTR: pos(start) < pos(end). width > 0.
	// If RTL: pos(start) > pos(end)?
	// Run.positions are visual. simpler:
	// If run is RTL, startGlyph (lower index) corresponds to... visual left?
	// Wait, internal storage of Run:
	// positions[0] is left-most visual? Or logical?
	// Usually positions are absolute X coordinates.
	// So finding min/max X is safer.
	if width < 0 {
		width = -width
		startX = endX
	}

	return ClipContext{
		Run:            run,
		Pos:            startGlyph,
		Size:           endGlyph - startGlyph,
		Clip:           models.Rect{Left: models.Scalar(runOffsetInLine + startX), Right: models.Scalar(runOffsetInLine + startX + width)},
		TextShift:      runOffsetInLine + startX,               // Is checking relative shift?
		ClippingNeeded: models.Scalar(width) < run.Advance().X, // heuristic
	}
}

// Paint paints the line.
func (tl *TextLine) Paint(painter ParagraphPainter, x, y float32) {
	// Background
	if tl.hasBackground {
		tl.ScanStyles(StyleTypeBackground, func(tr TextRange, ts TextStyle, cc ClipContext) {
			// Paint background
			if ts.HasBackground {
				// rect := cc.Clip.Offset(x + tl.offset.X, y + tl.offset.Y)
				// painter.DrawRect(rect, ts.Background)
			}
		})
	}

	// Shadows
	if tl.hasShadows {
		// ...
	}

	// Foreground (TextBlob)
	tl.ensureTextBlobCachePopulated()
	for _, record := range tl.textBlobCache {
		record.Paint(painter, x, y)
	}

	// Decorations
	if tl.hasDecorations {
		// ...
	}
}

// ensureTextBlobCachePopulated populates the blob cache.
func (tl *TextLine) ensureTextBlobCachePopulated() {
	if tl.textBlobCachePopulated {
		return
	}

	// Logic to build text blobs

	tl.textBlobCachePopulated = true
}

// CreateEllipsis creates the ellipsis run if needed.
func (tl *TextLine) CreateEllipsis(maxWidth float32, ellipsis string, ltr bool) {
	if ellipsis == "" {
		return
	}

	// Logic to find where to cut
	// Iterate backwards from ghost clusters

	// Simplified implementation:
	if tl.Width() <= maxWidth {
		return
	}

	// Iterate backwards through clusters
	// clustersWithGhosts includes trailing spaces
	// We want to remove clusters until we have room for ellipsis

	// Start with current width (excluding ellipsis as it is nil)
	width := tl.advance.X
	var ellipsisRun *Run

	// Range: ghosts + normal clusters
	// C++ iterates fGhostClusterRange which behaves as "all clusters capable of being trimmed"
	start := tl.ghostClusterRange.Start
	end := tl.ghostClusterRange.End

	for i := end - 1; i >= start; i-- {
		cluster := tl.owner.Cluster(i)

		// Shape ellipsis using this cluster's style/font
		// We should cache/optimize this (changes only when run changes)
		// Simplified: re-shape always
		ellipsisRun = tl.shapeEllipsis(ellipsis, cluster)

		if ellipsisRun == nil {
			// Failed to shape, try next?
			continue
		}

		ellipsisWidth := ellipsisRun.Advance().X

		// Check if we fit: current_width - cluster_width + ellipsis <= maxWidth
		// Note: width variable tracks current remaining line width

		if float32(width)+float32(ellipsisWidth) <= maxWidth {
			// Fits!
			tl.ellipsis = ellipsisRun
			tl.advance.X = models.Scalar(width) // Update line width to the reduction

			// Update ranges
			tl.clusterRange.End = i
			tl.ghostClusterRange.End = i
			tl.textExcludingSpaces.End = cluster.TextRange().Start
			tl.text.End = cluster.TextRange().Start
			tl.textIncludingNewlines.End = cluster.TextRange().Start

			return
		}

		// Remove cluster width
		width -= models.Scalar(cluster.Width())
	}

	// Fallback: clear line if ellipsis doesn't fit at all?
	// C++ handles "weird situation"
}

// shapeEllipsis shapes the ellipsis text.
func (tl *TextLine) shapeEllipsis(ellipsis string, cluster *Cluster) *Run {
	handler := &ellipsisRunHandler{
		lineHeight:     tl.sizes.Height(),
		useHalfLeading: false,
		baselineShift:  0,
		ellipsis:       ellipsis,
	}

	var run *Run
	if cluster != nil {
		run = cluster.Run()
		if run != nil {
			handler.useHalfLeading = run.UseHalfLeading()
			handler.baselineShift = run.BaselineShift()
			handler.lineHeight = run.HeightMultiplier() * float32(run.Font().Size())
		}
	}

	shapeWith := func(typeface interfaces.SkTypeface) *Run {
		fontSize := float32(14)
		if run != nil {
			fontSize = float32(run.Font().Size())
		}
		font := impl.NewFontWithTypefaceAndSize(typeface, models.Scalar(fontSize))

		hbShaper := shaper.NewHarfbuzzShaper()
		fontIter := shaper.NewTrivialFontRunIterator(font, len(ellipsis))
		bidiIter := shaper.NewTrivialBiDiRunIterator(0, len(ellipsis))
		scriptIter := shaper.NewTrivialScriptRunIterator(0, len(ellipsis))
		langIter := shaper.NewTrivialLanguageRunIterator("en", len(ellipsis))

		hbShaper.ShapeWithIterators(
			ellipsis,
			fontIter,
			bidiIter,
			scriptIter,
			langIter,
			nil,
			0,
			handler,
		)
		if handler.run != nil {
			handler.run.isEllipsis = true
		}
		return handler.run
	}

	if run != nil {
		r := shapeWith(run.Font().Typeface())
		if r != nil && r.IsResolved() {
			return r
		}
	}

	if fc := tl.owner.FontCollection(); fc != nil && fc.FontFallbackEnabled() {
		// Verify fields on FontStyle. Models.FontStyle contains Weight, Width, Slant.
		// Use models.FontStyleNormal object if it exists or construct one.
		// If models.FontStyleNormal is a function, call it?
		// Check models package usage. Assuming models.FontStyleNormal() for now based on lint.
		r := shapeWith(fc.DefaultFallback('.', models.FontStyleNormal(), "en"))
		if r != nil {
			return r
		}
	}

	return nil
}

type ellipsisRunHandler struct {
	run            *Run
	lineHeight     float32
	useHalfLeading bool
	baselineShift  float32
	ellipsis       string
}

func (h *ellipsisRunHandler) BeginLine()  {}
func (h *ellipsisRunHandler) CommitLine() {}
func (h *ellipsisRunHandler) RunInfo(info shaper.RunInfo) {
	h.run = NewRun(
		info,
		0,
		h.lineHeight/float32(info.Font.Size()),
		h.useHalfLeading,
		h.baselineShift,
		0,
		0,
	)
}
func (h *ellipsisRunHandler) CommitRunInfo() {}
func (h *ellipsisRunHandler) RunBuffer(info shaper.RunInfo) shaper.Buffer {
	if h.run == nil {
		return shaper.Buffer{}
	}
	return h.run.NewRunBuffer()
}
func (h *ellipsisRunHandler) CommitRunBuffer(info shaper.RunInfo) {}

// Helper structs

type TextBlobRecord struct {
	// ...
}

func (r *TextBlobRecord) Paint(painter ParagraphPainter, x, y float32) {
	// ...
}

// GetRectsForRange returns bounding boxes for the given text range.
func (tl *TextLine) GetRectsForRange(textRange TextRange, rectHeightStyle RectHeightStyle, rectWidthStyle RectWidthStyle) []TextBox {
	boxes := make([]TextBox, 0)

	// Check line intersection
	// Use textIncludingNewlines to ensure we cover the whole line range if needed
	lineRange := tl.textIncludingNewlines
	intersection := lineRange.Intersection(textRange)
	if intersection.Width() == 0 {
		// Handle zero-width range (cursor) if it falls exactly at start/end?
		if intersection.Start == intersection.End && (intersection.Start == lineRange.Start || intersection.Start == lineRange.End) {
			// fallback to standard processing to potentially generate a cursor rect
		} else {
			return boxes
		}
	}

	// Iterate runs
	for _, runIndex := range tl.runsInVisualOrder {
		run := tl.owner.Run(runIndex)
		runRange := run.TextRange()
		runIntersection := runRange.Intersection(textRange)

		if runIntersection.Width() == 0 && runIntersection.Start != runIntersection.End {
			continue
		}
		// If intersection is empty but we are at the edge, we might need a cursor rect?
		// For now focus on selection.

		if runIntersection.Start >= runIntersection.End {
			continue
		}

		// Calculate run bounds for this intersection
		// Iterate glyphs to find min/max X
		minX := float32(math.Inf(1))
		maxX := float32(math.Inf(-1))
		// found := -1

		// Access run internals (package private)
		glyphCount := len(run.glyphs)
		for i := 0; i < glyphCount; i++ {
			cluster := int(run.clusterIndexes[i])
			// Check if cluster is within intersection
			// Be careful with cluster mapping (logic is separate from run visual logic)
			// Simply check if cluster index is in range
			if cluster >= runIntersection.Start && cluster < runIntersection.End {
				pos := run.positions[i].X
				nextPos := run.positions[i+1].X // run.positions has glyphCount+1 elements
				// width := nextPos - pos

				// Handle RTL/LTR
				// If run is RTL, pos might be greater than nextPos?
				// Layout usually normalizes positions to be visual L->R or R->L?
				// positions are absolute X offsets.

				left := float32(pos)
				right := float32(nextPos)
				if left > right {
					left, right = right, left
				}

				if left < minX {
					minX = left
				}
				if right > maxX {
					maxX = right
				}
			}
		}

		if minX == float32(math.Inf(1)) {
			// No glyphs matched? Might be whitespace or unmapped
			continue
		}

		// Apply run offset
		minX += float32(run.offset.X)
		maxX += float32(run.offset.X)

		// Calculate vertical bounds based on RectHeightStyle
		top := float32(tl.offset.Y)
		bottom := float32(tl.offset.Y) + tl.Height()

		switch rectHeightStyle {
		case RectHeightStyleMax:
			// Use line height (already set)
		case RectHeightStyleIncludeLineSpacingTop:
			// Simplified to Max for now
		case RectHeightStyleIncludeLineSpacingMiddle:
			// Simplified to Max
		case RectHeightStyleIncludeLineSpacingBottom:
			// Simplified to Max
		case RectHeightStyleStrut:
			// Use strut
		case RectHeightStyleTight:
			// Use run metrics
			baseline := tl.Baseline()
			top = float32(tl.offset.Y) + baseline - run.Ascent()
			bottom = float32(tl.offset.Y) + baseline + run.Descent()
		}

		rect := models.Rect{
			Left:   models.Scalar(minX + float32(tl.offset.X)),
			Top:    models.Scalar(top),
			Right:  models.Scalar(maxX + float32(tl.offset.X)),
			Bottom: models.Scalar(bottom),
		}

		direction := TextDirectionLTR
		if run.bidiLevel%2 != 0 {
			direction = TextDirectionRTL
		}

		boxes = append(boxes, NewTextBox(rect, direction))
	}

	// Merge boxes if RectWidthStyle says so (e.g. Tight)
	// Only merge adjacent boxes?

	return boxes
}
