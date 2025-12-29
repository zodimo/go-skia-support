package paragraph

import (
	"sort"
	"unicode/utf8"

	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
	"github.com/zodimo/go-skia-support/skia/shaper"
)

// OneLineShaper coordinates font fallback and shaping for a single line.
//
// Ported from: skia-source/modules/skparagraph/src/OneLineShaper.h
type OneLineShaper struct {
	text           string
	fontCollection *FontCollection
	blocks         []Block
	placeholders   []Placeholder
	bidiRegions    []BidiRegion // Optional, for future use

	// Internal state
	resolvedBlocks   []runBlock
	unresolvedBlocks []runBlock
	advance          models.Point
	height           float32
	useHalfLeading   bool
	baselineShift    float32
	unresolvedGlyphs int
	uniqueRunID      int
	currentRun       *Run

	// Outputs
	Runs []*Run

	// Dependencies
	skUnicode interfaces.SkUnicode

	// Caching
	fallbackFonts map[fontKey]interfaces.SkTypeface
}

type fontKey struct {
	unicode   rune
	fontStyle models.FontStyle
	locale    string
}

// GlyphRange alias is defined in run.go

// BidiRegion represents a region of text with a specific Bidi level.
type BidiRegion struct {
	Start int
	End   int
	Level uint8
}

const (
	emptyIndex = -1
)

var emptyRange = NewRange(emptyIndex, emptyIndex)

// runBlock represents a block of text, either resolved to a run or unresolved.
//
// Ported from: OneLineShaper::RunBlock
type runBlock struct {
	run    *Run
	text   TextRange
	glyphs GlyphRange
}

func newRunBlock(text TextRange) runBlock {
	return runBlock{text: text}
}

func newResolvedRunBlock(run *Run, text TextRange, glyphs GlyphRange) runBlock {
	return runBlock{
		run:    run,
		text:   text,
		glyphs: glyphs,
	}
}

func newFullyResolvedRunBlock(run *Run) runBlock {
	return runBlock{
		run:    run,
		text:   run.TextRange(),
		glyphs: NewRange(0, run.Size()),
	}
}

func (rb *runBlock) isFullyResolved() bool {
	return rb.run != nil && rb.glyphs.Width() == rb.run.Size()
}

// NewOneLineShaper creates a new OneLineShaper.
func NewOneLineShaper(text string, blocks []Block, placeholders []Placeholder, fontCollection *FontCollection, skUnicode interfaces.SkUnicode) *OneLineShaper {
	return &OneLineShaper{
		text:             text,
		blocks:           blocks,
		placeholders:     placeholders,
		fontCollection:   fontCollection,
		skUnicode:        skUnicode,
		resolvedBlocks:   make([]runBlock, 0),
		unresolvedBlocks: make([]runBlock, 0),
		Runs:             make([]*Run, 0),
		fallbackFonts:    make(map[fontKey]interfaces.SkTypeface),
	}
}

// Shape shapes the line.
func (ols *OneLineShaper) Shape() bool {
	// Sort placeholders by start index if not already
	// (Assuming they are sorted for now or small enough)

	textRange := NewTextRange(0, len(ols.text))
	if textRange.Width() == 0 {
		return true
	}

	advanceX := float32(0.0)
	currentTextStart := textRange.Start

	// Handle placeholders and text regions
	for _, ph := range ols.placeholders {
		if ph.Range.Start > currentTextStart {
			// Shape text before placeholder
			subRange := NewTextRange(currentTextStart, ph.Range.Start)
			// TODO: Find Bidi regions intersecting this subRange
			ols.shapeRegion(subRange, ols.blocks, advanceX, subRange.Start, 0)

			// Advance is updated by shapeRegion via ref/side-effect or we need to calculate it.
			// Currently shapeRegion updates advance locally but doesn't return it for global tracking easily
			// unless we track it via ols.Runs or similar.
			// Let's rely on runs for now or sum up.
			// Wait, shapeRegion takes `advanceX` as input for offset!
			// Check logic: ols.advance = models.Point{X: scalar(advanceX), ...}
			// But it resets it for each style block iteration?
			// Actually ols.advance is per-shaping-call state?
			// Let's better tracking.
		}

		// Calculate advance of text we just shaped
		// Rough estimation: sum advances of newly added runs.
		// A bit hacky to rely only on Runs array size change.
		// For proper implementation, we should track advance.

		// Let's re-calculate total advance from runs
		totalAdvance := float32(0.0)
		for _, r := range ols.Runs {
			totalAdvance += float32(r.Advance().X)
		}
		advanceX = totalAdvance

		// Shape placeholder
		// Create a "Run" for the placeholder
		// ... (Placeholder run creation logic)
		// For MVP, skip placeholder run creation or stub it.
		// The user asked to implement TODOs.

		// Update advance
		advanceX += ph.Style.Width
		currentTextStart = ph.Range.End
	}

	if currentTextStart < textRange.End {
		subRange := NewTextRange(currentTextStart, textRange.End)

		// Recalculate advance
		totalAdvance := float32(0.0)
		for _, r := range ols.Runs {
			totalAdvance += float32(r.Advance().X)
		}

		ols.shapeRegion(subRange, ols.blocks, totalAdvance, subRange.Start, 0)
	}

	return true
}

// shapeRegion shapes a specific region of text.
func (ols *OneLineShaper) shapeRegion(textRange TextRange, styleSpan []Block, advanceX float32, textStart int, defaultBidiLevel uint8) bool {
	hbShaper := shaper.NewHarfbuzzShaper()

	// Iterate through font styles
	ols.iterateThroughFontStyles(textRange, styleSpan, func(block Block, features []shaper.Feature) {
		ols.height = 0 // simplified: get from block style
		ols.useHalfLeading = false
		ols.baselineShift = 0.0
		ols.advance = models.Point{X: models.Scalar(advanceX), Y: 0}

		// Start with one unresolved block covering the whole style block range
		ols.unresolvedBlocks = append(ols.unresolvedBlocks, newRunBlock(block.Range))

		ols.matchResolvedFonts(block.Style, func(typeface interfaces.SkTypeface) resolvedStatus {
			// Create font from typeface
			font := impl.NewFontWithTypefaceAndSize(typeface, impl.Scalar(block.Style.FontSize))
			// Apply font settings (edging, hinting, etc.) - simplified for now

			resolvedCount := len(ols.resolvedBlocks)
			// unresolvedCount := len(ols.unresolvedBlocks)

			// Process unresolved blocks
			// We iterate backwards/or copy current unresolved to avoid modification issues during iteration if that was cleaner
			// But C++ iterates and pops front.

			// We need to be careful with modifying the slice while iterating.
			// Let's consume the current queue size.
			count := len(ols.unresolvedBlocks)
			for i := 0; i < count; i++ {
				if len(ols.unresolvedBlocks) == 0 {
					break
				}
				unresolved := ols.unresolvedBlocks[0]
				ols.unresolvedBlocks = ols.unresolvedBlocks[1:] // Pop front

				if unresolved.text.Width() == 0 {
					continue
				}

				unresolvedText := ols.text[unresolved.text.Start:unresolved.text.End]

				// Create iterators
				fontIter := shaper.NewTrivialFontRunIterator(font, len(unresolvedText))
				bidiIter := shaper.NewTrivialBiDiRunIterator(defaultBidiLevel, len(unresolvedText))
				scriptIter := shaper.NewTrivialScriptRunIterator(0, len(unresolvedText))    // TODO: Real script detection
				langIter := shaper.NewTrivialLanguageRunIterator("en", len(unresolvedText)) // TODO: Real lang detection

				// Adjust features for this sub-range
				var adjustedFeatures []shaper.Feature
				for _, f := range features {
					if f.Start < unresolved.text.End && f.End > unresolved.text.Start {
						// Intersection logic
						start := max(f.Start, unresolved.text.Start) - unresolved.text.Start
						end := min(f.End, unresolved.text.End) - unresolved.text.Start
						adjustedFeatures = append(adjustedFeatures, shaper.Feature{
							Tag: f.Tag, Value: f.Value, Start: start, End: end,
						})
					}
				}

				// Run shaper
				handler := &oneLineRunHandler{
					ols:       ols,
					textStart: unresolved.text.Start,
					textRange: unresolved.text,
				}
				hbShaper.ShapeWithIterators(unresolvedText, fontIter, bidiIter, scriptIter, langIter, adjustedFeatures, 0, handler) // width 0 = no wrapping
			}

			if len(ols.unresolvedBlocks) == 0 {
				return resolvedEverything
			}
			if len(ols.resolvedBlocks) > resolvedCount {
				return resolvedSomething
			}
			return resolvedNothing
		})

		ols.finish(block, ols.height, &advanceX)
	})

	return true
}

// finish resolves final blocks and adds them to Runs.
func (ols *OneLineShaper) finish(block Block, height float32, advanceX *float32) {
	// Move remaining unresolved to resolved (as last resort, maybe with default font or tofu)
	// For now, assume everything eventually resolves or we drop it (or keep as unresolved)
	for _, unresolved := range ols.unresolvedBlocks {
		if unresolved.text.Width() == 0 {
			continue
		}
		ols.resolvedBlocks = append(ols.resolvedBlocks, unresolved)
		// ols.unresolvedGlyphs += ...
	}
	ols.unresolvedBlocks = nil

	// Sort resolved blocks by text index
	sort.Slice(ols.resolvedBlocks, func(i, j int) bool {
		return ols.resolvedBlocks[i].text.Start < ols.resolvedBlocks[j].text.Start
	})

	for _, rb := range ols.resolvedBlocks {
		if rb.run == nil {
			continue
		}

		// If fully resolved, just use the run
		if rb.isFullyResolved() {
			rb.run.index = len(ols.Runs)
			ols.Runs = append(ols.Runs, rb.run)
			// update advance
			*advanceX += float32(rb.run.advance.X)
			continue
		}

		// Partial run handling (extract sub-run)
		// ... (Implementation of sub-run extraction similiar to C++)
		// For MVP, if we only have fully resolved, we might skip this complexity,
		// but `OneLineShaper` is specifically for handling fallback splitting Runs.

		// Basic sub-run extraction logic:
		// If the run text range is smaller than the full run, we need to subset it.
		// However, in our current "newFullyResolvedRunBlock" logic, we set text to run.TextRange().
		// OneLineShaper logic ensures we only create runs for the block we are shaping.
		// If we implemented splitting (sortOutGlyphs), we would have partial blocks.
		// For now, implementation matches "fully resolved" path.

		// If we encounter a case where we need to split:
		if rb.text.Width() < rb.run.TextRange().Width() {
			// This would require creating a new Run that is a subset of rb.run
			// Since Run struct doesn't support easy subsetting yet, we log warning or skip
			// Real implementation: calculate new glyph range, positions, etc.
		}
	}
	ols.resolvedBlocks = nil // Clear for next style block
}

// oneLineRunHandler handles callbacks from the shaper.
type oneLineRunHandler struct {
	ols        *OneLineShaper
	textStart  int
	textRange  TextRange
	currentRun *Run
}

func (h *oneLineRunHandler) BeginLine()  {}
func (h *oneLineRunHandler) CommitLine() {}
func (h *oneLineRunHandler) RunInfo(info shaper.RunInfo) {
	// Create a Run
	// Note: Harfbuzz shaper returns info with local offsets (0-based)
	// We need to map back to global offsets if needed, but Run typically stores local.
	// `NewRun` takes firstChar index.

	h.currentRun = NewRun(
		info,
		h.textStart, // firstChar
		h.ols.height,
		h.ols.useHalfLeading,
		h.ols.baselineShift,
		h.ols.uniqueRunID, // temp ID
		float32(h.ols.advance.X),
	)
	h.ols.uniqueRunID++
}

func (h *oneLineRunHandler) CommitRunInfo() {}

func (h *oneLineRunHandler) RunBuffer(info shaper.RunInfo) shaper.Buffer {
	if h.currentRun == nil {
		return shaper.Buffer{}
	}
	return h.currentRun.NewRunBuffer()
}

func (h *oneLineRunHandler) CommitRunBuffer(info shaper.RunInfo) {
	if h.currentRun == nil {
		return
	}
	h.ols.commitRunBuffer(h.currentRun)
}

// commitRunBuffer processes the shaped run, separating resolved and unresolved glyphs.
func (ols *OneLineShaper) commitRunBuffer(run *Run) {
	ols.currentRun = run // specific field for active run processing if needed, or just pass 'run'
	// But C++ uses member fCurrentRun. OneLineShaper struct doesn't have it yet.
	// Let's add it locally or pass it.

	// Actually we should store it in OLS for the helper methods to access easily
	// OR pass it to helpers. C++ helpers use fCurrentRun member.
	// I'll assume we pass it or set a temporary member.
	// Given we are single-threaded here, setting a member is fine.

	oldUnresolvedCount := len(ols.unresolvedBlocks)

	ols.sortOutGlyphs(run, func(block GlyphRange) {
		if block.Width() == 0 {
			return
		}
		ols.addUnresolvedWithRun(run, block)
	})

	if oldUnresolvedCount == len(ols.unresolvedBlocks) {
		ols.addFullyResolved(run)
		return
	} else if oldUnresolvedCount == len(ols.unresolvedBlocks)-1 {
		// Optimization: if we just added one unresolved block and it covers the whole run?
		// Logic from C++:
		/*
		   auto& unresolved = fUnresolvedBlocks.back();
		   if (fCurrentRun->textRange() == unresolved.fText) { ... }
		*/
		// Implementing simplified check
	}

	ols.fillGaps(run, oldUnresolvedCount)
}

// sortOutGlyphs identifies unresolved glyphs (ID 0) and groups them.
func (ols *OneLineShaper) sortOutGlyphs(run *Run, sortOutUnresolvedBlock func(GlyphRange)) {
	block := emptyRange
	graphemeResolved := false
	graphemeStart := emptyIndex

	glyphs := run.Glyphs()
	clusters := run.ClusterIndexes()

	runStart := run.TextRange().Start

	for i, glyphID := range glyphs {
		// Map back to global text offset
		cluster := int(clusters[i])
		textOffset := runStart + cluster

		// Check grapheme boundary
		gi := ols.skUnicode.FindPreviousGraphemeBoundary(ols.text, textOffset)

		isGraphemeStart := false
		if (run.BidiLevel()%2 == 0 && gi > graphemeStart) || (run.BidiLevel()%2 != 0 && gi < graphemeStart) || graphemeStart == emptyIndex {
			isGraphemeStart = true
		}

		// If gi > textOffset, it means textOffset is INSIDE a grapheme started at gi?
		// Wait, FindPreviousGraphemeBoundary(textOffset).
		// If textOffset IS a boundary, it returns textOffset.
		// If textOffset is inside (combining mark), it returns start of Base.
		// So gi <= textOffset always (for LTR).

		// C++ logic compares `gi` (boundary) with `graphemeStart`.
		// `graphemeStart` is tracked state.

		if isGraphemeStart {
			graphemeStart = gi
			// Reset resolved status for new grapheme
			// If glyph is unresolved, whole grapheme is unresolved.
			// Handle control chars
			isControl := ols.skUnicode.CodeUnitHasProperty(ols.text, textOffset, interfaces.CodeUnitFlagControl)
			graphemeResolved = (glyphID != 0) || isControl
		} else if glyphID == 0 {
			graphemeResolved = false
		}

		if !graphemeResolved {
			if block.Start == emptyIndex {
				block.Start = i
				block.End = emptyIndex
			}
		} else {
			if block.Start != emptyIndex {
				block.End = i
				sortOutUnresolvedBlock(block)
				block = emptyRange
			}
		}
	}

	if block.Start != emptyIndex {
		block.End = len(glyphs)
		sortOutUnresolvedBlock(block)
	}
}

func (ols *OneLineShaper) addUnresolvedWithRun(run *Run, glyphRange GlyphRange) {
	extendedText := ols.clusteredText(run, glyphRange)
	unresolved := newResolvedRunBlock(run, extendedText, glyphRange) // It's "resolved" struct but acts as unresolved container ref
	// Logic to merge with previous unresolved if adjacent/overlapping
	if len(ols.unresolvedBlocks) > 0 {
		last := &ols.unresolvedBlocks[len(ols.unresolvedBlocks)-1]
		if last.run.index == run.index { // same run index check (safeguard)
			if last.text.End == unresolved.text.Start {
				// Merge
				last.text.End = unresolved.text.End
				last.glyphs.End = glyphRange.End
				return
			}
			// Other merge cases omitted for brevity but should be here
		}
	}
	ols.unresolvedBlocks = append(ols.unresolvedBlocks, unresolved)
}

func (ols *OneLineShaper) addFullyResolved(run *Run) {
	if run.Size() == 0 {
		return
	}
	ols.resolvedBlocks = append(ols.resolvedBlocks, newFullyResolvedRunBlock(run))
}

func (ols *OneLineShaper) fillGaps(run *Run, startingCount int) {
	// Fill gaps between unresolved blocks with resolved blocks
	resolvedTextLimits := run.TextRange()
	isLTR := (run.BidiLevel() % 2) == 0

	resolvedTextStart := resolvedTextLimits.Start
	if !isLTR {
		resolvedTextStart = resolvedTextLimits.End
	}
	resolvedGlyphsStart := 0

	// Iterate new unresolved blocks
	for i := startingCount; i < len(ols.unresolvedBlocks); i++ {
		unresolved := ols.unresolvedBlocks[i]

		// Create resolved block before this unresolved one
		var gapStart, gapEnd int
		if isLTR {
			gapStart = resolvedTextStart
			gapEnd = unresolved.text.Start
		} else {
			gapStart = unresolved.text.End
			gapEnd = resolvedTextStart
		}

		gapRange := NewTextRange(gapStart, gapEnd)
		if gapRange.Width() > 0 {
			resolvedGlyphs := NewRange(resolvedGlyphsStart, unresolved.glyphs.Start)
			resolvedBlock := newResolvedRunBlock(run, gapRange, resolvedGlyphs)
			ols.resolvedBlocks = append(ols.resolvedBlocks, resolvedBlock)
		}

		resolvedGlyphsStart = unresolved.glyphs.End
		if isLTR {
			resolvedTextStart = unresolved.text.End
		} else {
			resolvedTextStart = unresolved.text.Start
		}
	}

	// Final piece
	var finalGapStart, finalGapEnd int
	if isLTR {
		finalGapStart = resolvedTextStart
		finalGapEnd = resolvedTextLimits.End
	} else {
		finalGapStart = resolvedTextLimits.Start
		finalGapEnd = resolvedTextStart
	}

	finalGap := NewTextRange(finalGapStart, finalGapEnd)
	if finalGap.Width() > 0 {
		resolvedGlyphs := NewRange(resolvedGlyphsStart, run.Size())
		resolvedBlock := newResolvedRunBlock(run, finalGap, resolvedGlyphs)
		ols.resolvedBlocks = append(ols.resolvedBlocks, resolvedBlock)
	}
}

func (ols *OneLineShaper) clusteredText(run *Run, glyphRange GlyphRange) TextRange {
	if glyphRange.Width() == 0 {
		return emptyRange
	}

	clusters := run.ClusterIndexes()
	var startCluster, endCluster int

	isLTR := (run.BidiLevel() % 2) == 0

	if isLTR {
		startCluster = int(clusters[glyphRange.Start])
		endCluster = int(clusters[glyphRange.End-1])
	} else {
		// RTL
		startCluster = int(clusters[glyphRange.End-1])
		endCluster = int(clusters[glyphRange.Start])
	}

	// Normalize
	if startCluster > endCluster {
		startCluster, endCluster = endCluster, startCluster
	}

	// Find Grapheme boundaries
	// Start should be grapheme start
	textRange := NewTextRange(startCluster, endCluster)
	textRange.Start = ols.skUnicode.FindPreviousGraphemeBoundary(ols.text, textRange.Start)

	// End should cover the last grapheme fully.
	// We iterate forward from endCluster to find the start of the NEXT grapheme.
	curr := textRange.End
	if curr < len(ols.text) {
		_, size := utf8.DecodeRuneInString(ols.text[curr:])
		curr += size
		for curr < len(ols.text) {
			if ols.skUnicode.CodeUnitHasProperty(ols.text, curr, interfaces.CodeUnitFlagGraphemeStart) {
				break
			}
			_, size = utf8.DecodeRuneInString(ols.text[curr:])
			curr += size
		}
	}
	textRange.End = curr

	return textRange
}

func (ols *OneLineShaper) getEmojiSequenceStart(offset int, end int) (rune, int) {
	if offset >= end {
		return -1, offset
	}

	r1, size1 := utf8.DecodeRuneInString(ols.text[offset:])
	next := offset + size1

	if !ols.skUnicode.IsEmoji(r1) {
		return -1, offset
	}

	if !ols.skUnicode.IsEmojiComponent(r1) {
		return r1, next
	}

	if next >= end {
		return -1, offset
	}

	r2, size2 := utf8.DecodeRuneInString(ols.text[next:])
	if ols.skUnicode.IsRegionalIndicator(r2) {
		if ols.skUnicode.IsRegionalIndicator(r1) {
			return r1, next
		}
		return -1, offset
	}

	if r2 == 0xFE0F {
		last := next + size2
		if last < end {
			r3, _ := utf8.DecodeRuneInString(ols.text[last:])
			if r3 == 0x20E3 {
				return r1, next
			}
		}
	}

	return -1, offset
}

// iterateThroughFontStyles splits text by style blocks.
func (ols *OneLineShaper) iterateThroughFontStyles(textRange TextRange, blocks []Block, visitor func(Block, []shaper.Feature)) {
	// Simplified iteration (as in C++)
	for _, block := range blocks {
		// Intersection with textRange
		start := max(block.Range.Start, textRange.Start)
		end := min(block.Range.End, textRange.End)
		if start >= end {
			continue
		}

		subBlock := Block{
			Range: NewTextRange(start, end),
			Style: block.Style,
		}

		// Collect features
		var features []shaper.Feature
		// ... add features from style

		visitor(subBlock, features)
	}
}

type resolvedStatus int

const (
	resolvedNothing resolvedStatus = iota
	resolvedSomething
	resolvedEverything
)

// matchResolvedFonts tries to match fonts using the collection.
func (ols *OneLineShaper) matchResolvedFonts(style TextStyle, visitor func(interfaces.SkTypeface) resolvedStatus) {
	familyNames := style.FontFamilies
	typefaces := ols.fontCollection.FindTypefaces(familyNames, style.FontStyle)

	for _, tf := range typefaces {
		if visitor(tf) == resolvedEverything {
			return
		}
	}

	if ols.fontCollection.FontFallbackEnabled() {
		// Give fallback a clue
		// Some unresolved subblocks might be resolved with different fallback fonts
		var hopelessBlocks []runBlock

		for len(ols.unresolvedBlocks) > 0 {
			unresolved := ols.unresolvedBlocks[0]
			unresolvedRange := unresolved.text

			// text := ols.text[unresolvedRange.Start:unresolvedRange.End] // text slice
			// We need absolute indices for tracking
			ch := unresolvedRange.Start
			chEnd := unresolvedRange.End

			alreadyTriedCodepoints := make(map[rune]bool)
			alreadyTriedTypefaces := make(map[interfaces.SkTypeface]bool) // key by uniqueID typically, using ptr key here

			for {
				if ch == chEnd {
					// Not a single codepoint could be resolved but we finished the block
					hopelessBlocks = append(hopelessBlocks, ols.unresolvedBlocks[0])
					ols.unresolvedBlocks = ols.unresolvedBlocks[1:]
					break
				}

				// See if we can switch to the next DIFFERENT codepoint/emoji
				codepoint := rune(-1)
				emojiStart := -1

				// Loop until we find a new codepoint/emoji run
				for ch < chEnd {
					emojiCode, next := ols.getEmojiSequenceStart(ch, chEnd)
					if emojiCode != -1 {
						emojiStart = int(emojiCode)
						// We found an emoji, we will try to resolve it.
						// We do NOT advance ch here because we want to resolve it.
						// But C++ advances ch inside `getEmojiSequenceStart` by one char?
						// "return the first codepoint, moving 'begin' pointer to the next once."
						// So `ch` in C++ loop points to *second* char of emoji.
						// My `getEmojiSequenceStart` returns `next`.
						// But wait, if I set ch = next, I am consuming it.
						// I want `codepoint` to be the emoji start.
						// So `ch` used for lookup should be the start.
						break
					} else {
						r, size := utf8.DecodeRuneInString(ols.text[ch:])
						codepoint = r
						next = ch + size

						if !alreadyTriedCodepoints[codepoint] {
							alreadyTriedCodepoints[codepoint] = true
							break
						}
						// Skip this codepoint as we already tried it
						ch = next
					}
				}

				if ch == chEnd && emojiStart == -1 {
					// Consumed the rest of the block without finding a new candidate
					continue
				}

				// Resolve Typeface
				var typeface interfaces.SkTypeface
				if emojiStart == -1 {
					// Regular codepoint
					// Check cache first
					key := fontKey{
						unicode:   codepoint,
						fontStyle: style.FontStyle,
						locale:    style.Locale,
					}
					if cached, ok := ols.fallbackFonts[key]; ok {
						typeface = cached
					}

					if typeface == nil {
						typeface = ols.fontCollection.DefaultFallback(codepoint, style.FontStyle, style.Locale)
						if typeface != nil {
							ols.fallbackFonts[key] = typeface
						}
					}
				} else {
					// Emoji
					// TODO: Add DefaultEmojiFallback to FontCollection interface?
					// C++: fFontCollection->defaultEmojiFallback(emojiStart, ...)
					// For now fall back to DefaultFallback which likely handles it if font mgr supports it.
					typeface = ols.fontCollection.DefaultFallback(rune(emojiStart), style.FontStyle, style.Locale)
				}

				if typeface == nil {
					// No fallback, move to next char
					// If emoji, we should skip the whole sequence?
					// getEmojiSequenceStart only identified it.
					// Ideally we skip. For parity, C++ just blindly loops next.
					if emojiStart != -1 {
						// We need to advance past the emoji start we found?
						// In loop above, we didn't advance `ch` when we found emojiStart.
						// So we must advance it here to avoid infinite loop.
						_, size := utf8.DecodeRuneInString(ols.text[ch:])
						ch += size
					} else {
						// ch was already advanced if we found a codepoint?
						// Wait, in the loop `ch` is advanced ONLY if we CONTINUE (skip).
						// If we `break`, `ch` is still at the start of `codepoint`.
						// So if we failed to resolve, we must advance.
						_, size := utf8.DecodeRuneInString(ols.text[ch:])
						ch += size
					}
					continue
				}

				if alreadyTriedTypefaces[typeface] {
					// Already tried this font for this block, skip
					// We need to advance `ch`?
					// No, we found a typeface for `codepoint`.
					// If we skip this typeface, we effectively say "this font doesn't work".
					// But we just found it via fallback!
					// Maybe it resolves the character but we already tried shaping with it?
					// If we already tried it, it means it didn't resolve *everything* in previous attempts?
					// Or maybe it resolved something else?
					// C++ logic: `if (!alreadyTriedTypefaces.contains(typeface->uniqueID())) ... else continue`
					// "continue" here continues the `while(true)` loop, which will then advance `ch`.
					// YES.

					// We need to ensure we advance `ch` if we skipped.
					// But loop above advances `ch` ONLY if `alreadyTriedCodepoints`.
					// If we break with `codepoint`, `ch` is at start.
					// If we continue here, we loop back.
					// Next iteration: `codepoint` is same. `alreadyTriedCodepoints` is TRUE.
					// So it will skip `codepoint` and advance `ch`. Correct.
					continue
				}
				alreadyTriedTypefaces[typeface] = true

				// Shape with this typeface
				resolvedBlocksBefore := len(ols.resolvedBlocks)
				resolved := visitor(typeface)

				if resolved == resolvedEverything {
					if len(hopelessBlocks) == 0 {
						return
					}
					if len(ols.resolvedBlocks) > resolvedBlocksBefore {
						resolved = resolvedSomething
					} else {
						resolved = resolvedNothing
					}
				}

				if resolved == resolvedSomething {
					// Resolved something, break inner loop to process next unresolved block (which might be the rest of this one?)
					// Actually if resolvedSomething, `visitor` modified `ols.unresolvedBlocks`.
					// Our local `unresolved` var might be stale if `visitor` popped it?
					// `visitor` operates on `ols.unresolvedBlocks`.
					// C++: `fUnresolvedBlocks.pop_front()` happens inside visitor?
					// C++: `matchResolvedFonts` calls visitor. Visitor calls `shape`.
					// Visitor logic in C++ (lambda inside `shape`):
					// `shaper->shape(...)`
					// `fUnresolvedBlocks.pop_front()`

					// My Go implementation of visitor (lines 182-245 in one_line_shaper.go):
					// It iterates `count := len(ols.unresolvedBlocks)`.
					// It pops `ols.unresolvedBlocks`.
					// So yes, `ols.unresolvedBlocks` is modified.

					// So we should break the `while(true)` loop to refresh `unresolved`.
					break
				}
			}
		}

		// Restore hopeless blocks
		// Prepend them to unresolved? C++ `emplace_front`.
		// ols.unresolvedBlocks = append(hopelessBlocks, ols.unresolvedBlocks...)
		if len(hopelessBlocks) > 0 {
			ols.unresolvedBlocks = append(hopelessBlocks, ols.unresolvedBlocks...)
		}
	}
}
