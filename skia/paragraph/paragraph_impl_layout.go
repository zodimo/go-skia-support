package paragraph

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/models"
	"golang.org/x/text/unicode/bidi"
)

// --- Layout (main entry point) ---

// Layout performs the paragraph layout at the given width.
//
// Ported from: ParagraphImpl::layout (lines 138-242)
func (p *ParagraphImpl) Layout(width float32) {
	// Apply rounding hack for Flutter compatibility
	floorWidth := width
	if p.paragraphStyle.ApplyRoundingHack {
		floorWidth = float32(math.Floor(float64(floorWidth)))
	}

	// Check if we can reuse previous layout results
	if (!math.IsInf(float64(width), 0) || p.longestLine <= floorWidth) &&
		p.state >= StateLineBroken &&
		len(p.lines) == 1 && p.lines[0].ellipsis == nil {
		// Most common case: single line without ellipsis
		p.width = floorWidth
		p.state = StateShaped
	} else if p.state >= StateLineBroken && p.oldWidth != floorWidth {
		// Width changed, need to re-break lines
		p.state = StateShaped
	}
	// else: Nothing changed, reuse previous results

	if p.state < StateShaped {
		// Need to shape the text
		if p.state < StateIndexed {
			// First layout: compute text properties
			if p.computeCodeUnitProperties() {
				p.state = StateIndexed
			}
		}

		// Clear previous runs and clusters
		p.runs = make([]*Run, 0)
		p.clusters = make([]*Cluster, 0)
		p.clustersIndexFromCodeUnit = make([]int, len(p.text)+1)
		for i := range p.clustersIndexFromCodeUnit {
			p.clustersIndexFromCodeUnit[i] = -1
		}

		if !p.shapeTextIntoEndlessLine() {
			// Empty or failed shaping
			p.resetContext()
			p.resolveStrut()
			p.computeEmptyMetrics()
			p.lines = make([]*TextLine, 0)

			// Set metrics for empty paragraph
			p.width = floorWidth
			p.height = p.emptyMetrics.Height()
			if p.strutEnabled() && p.strutForceHeight() {
				p.height = p.strutMetrics.Height()
			}
			p.alphabeticBaseline = p.emptyMetrics.AlphabeticBaseline()
			p.ideographicBaseline = p.emptyMetrics.IdeographicBaseline()
			p.longestLine = -math.MaxFloat32
			p.minIntrinsicWidth = 0
			p.maxIntrinsicWidth = 0
			p.oldWidth = floorWidth
			p.oldHeight = p.height
			return
		}
		p.state = StateShaped
	}

	if p.state == StateShaped {
		p.resetContext()
		p.resolveStrut()
		p.computeEmptyMetrics()
		p.lines = make([]*TextLine, 0)
		p.breakShapedTextIntoLines(floorWidth)
		p.state = StateLineBroken
	}

	if p.state == StateLineBroken {
		p.resetShifts()
		p.formatLines(p.width)
		p.state = StateFormatted
	}

	p.oldWidth = floorWidth
	p.oldHeight = p.height

	// Apply rounding hack for Flutter compatibility
	if p.paragraphStyle.ApplyRoundingHack {
		p.minIntrinsicWidth = littleRound(p.minIntrinsicWidth)
		p.maxIntrinsicWidth = littleRound(p.maxIntrinsicWidth)
	}

	// Flutter-specific: single line or unlimited lines with ellipsis
	maxLines := p.paragraphStyle.MaxLines
	if maxLines == 1 || (maxLines == 0 && p.paragraphStyle.Ellipsis != "") {
		p.minIntrinsicWidth = p.maxIntrinsicWidth
	}

	// Ensure min <= max
	if p.maxIntrinsicWidth < p.minIntrinsicWidth {
		p.maxIntrinsicWidth = p.minIntrinsicWidth
	}
}

// --- Internal layout methods ---

// resetContext resets computed layout values.
func (p *ParagraphImpl) resetContext() {
	p.alphabeticBaseline = 0
	p.height = 0
	p.width = 0
	p.ideographicBaseline = 0
	p.maxIntrinsicWidth = 0
	p.minIntrinsicWidth = 0
	p.longestLine = 0
	p.maxWidthWithTrailing = 0
	p.exceededMaxLines = false
}

// resetShifts resets justification shifts on all runs.
func (p *ParagraphImpl) resetShifts() {
	for _, run := range p.runs {
		run.ResetJustificationShifts()
	}
}

// computeCodeUnitProperties computes bidi and whitespace properties.
func (p *ParagraphImpl) computeCodeUnitProperties() bool {
	// Simplified implementation without full Unicode support
	// Initialize code unit properties based on basic whitespace detection
	p.codeUnitProperties = make([]int, len(p.text))

	// Simple whitespace and line break detection
	p.trailingSpaces = len(p.text)
	firstWhitespace := -1

	for i := 0; i < len(p.text); i++ {
		ch := p.text[i]
		flags := 0

		// Basic whitespace detection
		if ch == ' ' || ch == '\t' {
			flags |= 0x0008 // partOfWhiteSpaceBreak
			if p.trailingSpaces == len(p.text) {
				p.trailingSpaces = i
			}
			if firstWhitespace == -1 {
				firstWhitespace = i
			}
		} else {
			p.trailingSpaces = len(p.text)
		}

		// Hard line break detection
		if ch == '\n' {
			flags |= 0x0010 // hardLineBreakBefore (set on next char)
			p.hasLineBreaks = true
		}

		p.codeUnitProperties[i] = flags
	}

	// Set hard line break flags on characters after \n
	for i := 1; i < len(p.text); i++ {
		if p.text[i-1] == '\n' {
			p.codeUnitProperties[i] |= 0x0010 // hardLineBreakBefore
		}
	}

	if firstWhitespace != -1 && firstWhitespace < p.trailingSpaces {
		p.hasWhitespacesInside = true
	}

	// BiDi Analysis
	// Use golang.org/x/text/unicode/bidi to determine embedding levels
	paragraphDir := bidi.LeftToRight
	if p.paragraphStyle.TextDirection == TextDirectionRTL {
		paragraphDir = bidi.RightToLeft
	}

	// Analyze the text
	fallback := true
	if len(p.text) > 0 {
		var bidiPara bidi.Paragraph
		// Use bidi.DefaultDirection to set the base direction preference
		_, err := bidiPara.SetString(p.text, bidi.DefaultDirection(paragraphDir))

		if err == nil {
			// Get ordering to ensure we have runs
			ordering, err := bidiPara.Order()
			if err == nil && ordering.NumRuns() > 0 {
				run := bidiPara.RunAt(0)
				if len(run.String()) > 0 {
					fallback = false
					p.bidiRegions = make([]BidiRegion, 0)

					for pos := 0; pos < len(p.text); {
						run := bidiPara.RunAt(pos)
						text := run.String()
						length := len(text)
						if length == 0 {
							break
						}

						dir := run.Direction()
						level := uint8(0)
						if dir == bidi.RightToLeft {
							level = 1
						}

						p.bidiRegions = append(p.bidiRegions, BidiRegion{
							Start: pos,
							End:   pos + length,
							Level: level,
						})

						pos += length
					}
				}
			}
		}
	}

	if fallback {
		// Fallback for empty text or failure
		level := uint8(0)
		if p.paragraphStyle.TextDirection == TextDirectionRTL {
			level = 1
		}
		p.bidiRegions = []BidiRegion{{Start: 0, End: len(p.text), Level: level}}
	}

	return true
}

// shapeTextIntoEndlessLine shapes the text using OneLineShaper.
func (p *ParagraphImpl) shapeTextIntoEndlessLine() bool {
	if len(p.text) == 0 {
		return false
	}

	// Clear unresolved tracking
	p.unresolvedCodepoints = make(map[rune]struct{})

	// Create shaper and shape
	shaper := NewOneLineShaper(p.text, p.textStyles, p.placeholders, p.fontCollection, p.unicode)
	result := shaper.Shape()
	p.unresolvedGlyphs = shaper.unresolvedGlyphs

	// Copy runs from shaper
	p.runs = shaper.Runs

	// Build cluster table with spacing
	p.applySpacingAndBuildClusterTable()

	return result
}

// applySpacingAndBuildClusterTable builds clusters with letter/word spacing.
func (p *ParagraphImpl) applySpacingAndBuildClusterTable() {
	// Check if we need to apply any spacing
	letterSpacingStyles := 0
	hasWordSpacing := false
	for _, block := range p.textStyles {
		if block.Range.Width() > 0 {
			if !nearlyZero(block.Style.LetterSpacing) {
				letterSpacingStyles++
			}
			if !nearlyZero(block.Style.WordSpacing) {
				hasWordSpacing = true
			}
		}
	}

	// For simplicity, always build clusters first
	p.buildClusterTable()

	// TODO: Implement full spacing logic as in C++ if needed
	_ = letterSpacingStyles
	_ = hasWordSpacing
}

// buildClusterTable builds the cluster lookup table.
func (p *ParagraphImpl) buildClusterTable() {
	// Count total clusters needed
	clusterCount := 1
	for _, run := range p.runs {
		if run.IsPlaceholder() {
			clusterCount++
		} else {
			clusterCount += run.Size()
		}
	}

	p.clusters = make([]*Cluster, 0, clusterCount)
	p.clustersIndexFromCodeUnit = make([]int, len(p.text)+1)
	for i := range p.clustersIndexFromCodeUnit {
		p.clustersIndexFromCodeUnit[i] = -1
	}

	// Build clusters from runs
	for _, run := range p.runs {
		runStart := len(p.clusters)

		if run.IsPlaceholder() {
			// Placeholder gets one cluster
			tr := run.TextRange()
			for i := tr.Start; i < tr.End; i++ {
				if i < len(p.clustersIndexFromCodeUnit) {
					p.clustersIndexFromCodeUnit[i] = len(p.clusters)
				}
			}
			advance := run.Advance()
			cluster := NewCluster(p, run.Index(), 0, 1, tr, float32(advance.X), float32(advance.Y))
			p.clusters = append(p.clusters, cluster)
		} else {
			// Create clusters from glyphs - one cluster per glyph for simplicity
			// In full implementation, we'd group by text cluster index
			glyphCount := run.Size()
			tr := run.TextRange()

			if glyphCount > 0 {
				// Calculate width per cluster (simplified)
				advance := run.Advance()
				widthPerCluster := float32(advance.X) / float32(glyphCount)
				charsPerCluster := (tr.End - tr.Start) / glyphCount
				if charsPerCluster < 1 {
					charsPerCluster = 1
				}

				for g := 0; g < glyphCount; g++ {
					charStart := tr.Start + g*charsPerCluster
					charEnd := charStart + charsPerCluster
					if charEnd > tr.End {
						charEnd = tr.End
					}
					if g == glyphCount-1 {
						charEnd = tr.End
					}

					for i := charStart; i < charEnd; i++ {
						if i < len(p.clustersIndexFromCodeUnit) {
							p.clustersIndexFromCodeUnit[i] = len(p.clusters)
						}
					}

					clusterRange := NewTextRange(charStart, charEnd)
					cluster := NewCluster(p, run.Index(), g, g+1, clusterRange, widthPerCluster, 0)

					// Set break properties based on text content
					if charStart < len(p.text) {
						ch := p.text[charStart]
						isWhitespace := ch == ' ' || ch == '\t'
						isHardBreak := ch == '\n'
						cluster.SetBreakType(isWhitespace, false, isHardBreak, false)
					}

					p.clusters = append(p.clusters, cluster)
				}
			} else {
				// Empty run: create one empty cluster
				for i := tr.Start; i < tr.End; i++ {
					if i < len(p.clustersIndexFromCodeUnit) {
						p.clustersIndexFromCodeUnit[i] = len(p.clusters)
					}
				}
				cluster := NewCluster(p, run.Index(), 0, 0, tr, 0, 0)
				p.clusters = append(p.clusters, cluster)
			}
		}

		run.SetClusterRange(runStart, len(p.clusters))
		advance := run.Advance()
		p.maxIntrinsicWidth += float32(advance.X)
	}

	// Add end marker
	if len(p.text) < len(p.clustersIndexFromCodeUnit) {
		p.clustersIndexFromCodeUnit[len(p.text)] = len(p.clusters)
	}
	endCluster := NewCluster(p, -1, 0, 0, NewTextRange(len(p.text), len(p.text)), 0, 0)
	p.clusters = append(p.clusters, endCluster)
}

// breakShapedTextIntoLines breaks text into lines using TextWrapper.
func (p *ParagraphImpl) breakShapedTextIntoLines(maxWidth float32) {
	// Short path: single run, no breaks, fits in width
	if !p.hasLineBreaks && !p.hasWhitespacesInside &&
		len(p.placeholders) == 1 && len(p.runs) == 1 {

		run := p.runs[0]
		runAdvance := run.Advance()
		if float32(runAdvance.X) <= maxWidth {
			advanceX := runAdvance.X
			textRange := NewTextRange(0, len(p.text))
			textExcludingSpaces := NewTextRange(0, p.trailingSpaces)

			metrics := NewInternalLineMetrics()
			metrics.ForceStrut = p.strutForceHeight()
			metrics.AddRun(run)

			// Apply text height behavior
			if p.paragraphStyle.TextHeightBehavior&TextHeightBehaviorDisableFirstAscent != 0 {
				metrics.Ascent = metrics.RawAscent
			}
			if p.paragraphStyle.TextHeightBehavior&TextHeightBehaviorDisableLastDescent != 0 {
				metrics.Descent = metrics.RawDescent
			}

			if p.strutEnabled() {
				p.strutMetrics.UpdateLineMetrics(&metrics)
			}

			// Find trailing spaces cluster
			trailingSpaces := len(p.clusters)
			for trailingSpaces > 0 {
				trailingSpaces--
				cluster := p.clusters[trailingSpaces]
				if !cluster.IsWhitespaceBreak() {
					trailingSpaces++
					break
				}
				advanceX -= models.Scalar(cluster.Width())
			}

			advanceY := models.Scalar(metrics.Height())
			clusterRange := NewClusterRange(0, trailingSpaces)
			clusterRangeWithGhosts := NewClusterRange(0, len(p.clusters)-1)

			line := NewTextLine(
				p,
				models.Point{X: 0, Y: 0},
				models.Point{X: advanceX, Y: advanceY},
				NewBlockRange(0, len(p.textStyles)),
				textExcludingSpaces,
				textRange,
				textRange,
				clusterRange,
				clusterRangeWithGhosts,
				float32(runAdvance.X),
				metrics,
			)
			p.lines = append(p.lines, line)

			p.longestLine = float32(advanceX)
			if nearlyZero(float32(advanceX)) {
				p.longestLine = float32(runAdvance.X)
			}
			p.height = float32(advanceY)
			p.width = maxWidth
			p.maxIntrinsicWidth = float32(runAdvance.X)
			p.minIntrinsicWidth = float32(advanceX)
			p.alphabeticBaseline = metrics.AlphabeticBaseline()
			p.ideographicBaseline = metrics.IdeographicBaseline()
			if len(p.lines) > 0 {
				p.alphabeticBaseline = p.lines[0].Baseline()
				p.ideographicBaseline = p.lines[0].sizes.IdeographicBaseline()
			}
			p.exceededMaxLines = false
			return
		}
	}

	// Full line breaking with TextWrapper
	wrapper := NewTextWrapper()
	wrapper.BreakTextIntoLines(
		p,
		maxWidth,
		func(textExcludingSpaces, text, textWithNewlines TextRange,
			clusters, clustersWithGhosts ClusterRange,
			widthWithSpaces float32,
			startPos, endPos int,
			offset, advance models.Point,
			metrics InternalLineMetrics,
			addEllipsis bool) {

			blocks := p.findAllBlocks(textExcludingSpaces)
			line := NewTextLine(
				p,
				offset,
				advance,
				blocks,
				textExcludingSpaces,
				text,
				textWithNewlines,
				clusters,
				clustersWithGhosts,
				widthWithSpaces,
				metrics,
			)
			p.lines = append(p.lines, line)

			if addEllipsis {
				line.CreateEllipsis(maxWidth, p.getEllipsis(), true)
			}

			lineWidth := line.Width()
			if nearlyZero(lineWidth) {
				lineWidth = widthWithSpaces
			}
			if lineWidth > p.longestLine {
				p.longestLine = lineWidth
			}
		},
	)

	p.height = wrapper.Height()
	p.width = maxWidth
	p.maxIntrinsicWidth = wrapper.MaxIntrinsicWidth()
	p.minIntrinsicWidth = wrapper.MinIntrinsicWidth()
	if len(p.lines) > 0 {
		p.alphabeticBaseline = p.lines[0].Baseline()
		p.ideographicBaseline = p.lines[0].sizes.IdeographicBaseline()
	} else {
		p.alphabeticBaseline = p.emptyMetrics.AlphabeticBaseline()
		p.ideographicBaseline = p.emptyMetrics.IdeographicBaseline()
	}
	p.exceededMaxLines = wrapper.ExceededMaxLines()
}

// formatLines formats each line based on alignment.
func (p *ParagraphImpl) formatLines(maxWidth float32) {
	align := p.paragraphStyle.EffectiveAlign()

	// Check if left-aligned
	isLeftAligned := align == TextAlignLeft ||
		(align == TextAlignJustify && p.paragraphStyle.TextDirection == TextDirectionLTR)

	// Clear lines if infinite width and not left-aligned
	if math.IsInf(float64(maxWidth), 0) && !isLeftAligned {
		p.lines = make([]*TextLine, 0)
		return
	}

	for _, line := range p.lines {
		line.Format(align, maxWidth)
	}
}

// resolveStrut computes strut metrics if enabled.
func (p *ParagraphImpl) resolveStrut() {
	strutStyle := p.paragraphStyle.StrutStyle
	if !strutStyle.StrutEnabled || strutStyle.FontSize < 0 {
		return
	}

	// Find typefaces for strut
	typefaces := p.fontCollection.FindTypefaces(
		strutStyle.FontFamilies,
		strutStyle.FontStyle,
	)
	if len(typefaces) == 0 {
		return
	}

	// Get font metrics
	fontSize := strutStyle.FontSize
	strutLeading := float32(0)
	if strutStyle.Leading >= 0 {
		strutLeading = strutStyle.Leading * fontSize
	}

	// Estimate metrics (should use actual font metrics)
	ascent := -fontSize * 0.8
	descent := fontSize * 0.2

	if strutStyle.HeightOverride {
		if strutStyle.HalfLeading {
			occupiedHeight := descent - ascent
			flexibleHeight := strutStyle.Height*fontSize - occupiedHeight
			flexibleHeight /= 2
			p.strutMetrics = NewInternalLineMetricsFromValues(
				ascent-flexibleHeight,
				descent+flexibleHeight,
				strutLeading,
			)
		} else {
			metricsHeight := descent - ascent
			multiplier := float32(1)
			if metricsHeight != 0 {
				multiplier = strutStyle.Height * fontSize / metricsHeight
			}
			p.strutMetrics = NewInternalLineMetricsFromValues(
				ascent*multiplier,
				descent*multiplier,
				strutLeading,
			)
		}
	} else {
		p.strutMetrics = NewInternalLineMetricsFromValues(
			ascent,
			descent,
			strutLeading,
		)
	}
	p.strutMetrics.ForceStrut = strutStyle.ForceStrutHeight
}

// computeEmptyMetrics computes metrics for empty paragraphs.
func (p *ParagraphImpl) computeEmptyMetrics() {
	if len(p.textStyles) == 0 {
		p.emptyMetrics = NewInternalLineMetrics()
		return
	}

	// Use first text style to compute metrics
	style := p.textStyles[0].Style
	fontSize := style.FontSize
	if fontSize <= 0 {
		fontSize = 14 // Default size
	}

	// Estimate metrics
	ascent := -fontSize * 0.8
	descent := fontSize * 0.2
	leading := float32(0)

	p.emptyMetrics = NewInternalLineMetricsFromValues(ascent, descent, leading)
}

// findAllBlocks finds all blocks covering a text range.
func (p *ParagraphImpl) findAllBlocks(textRange TextRange) BlockRange {
	begin := -1
	end := -1

	for i, block := range p.textStyles {
		if block.Range.End <= textRange.Start {
			continue
		}
		if block.Range.Start >= textRange.End {
			break
		}
		if begin == -1 {
			begin = i
		}
		end = i
	}

	if begin == -1 || end == -1 {
		return EmptyRange
	}
	return NewBlockRange(begin, end+1)
}

// addUnresolvedCodepoints adds codepoints from a range to unresolved set.
func (p *ParagraphImpl) addUnresolvedCodepoints(textRange TextRange) {
	text := p.textRange(textRange)
	for _, r := range text {
		p.unresolvedCodepoints[r] = struct{}{}
	}
}
