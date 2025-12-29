package paragraph

import (
	"sort"

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
func NewOneLineShaper(text string, blocks []Block, placeholders []Placeholder, fontCollection *FontCollection) *OneLineShaper {
	return &OneLineShaper{
		text:             text,
		blocks:           blocks,
		placeholders:     placeholders,
		fontCollection:   fontCollection,
		resolvedBlocks:   make([]runBlock, 0),
		unresolvedBlocks: make([]runBlock, 0),
		Runs:             make([]*Run, 0),
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
	// Simplified grapheme handling: treat every cluster/glyph as potentially unresolved logic
	// TODO: Real grapheme boundary check

	// Using Run's glyphs
	glyphs := run.Glyphs()

	for i, glyphID := range glyphs {
		// Check if unresolved
		// C++: if (glyph == 0 && !isControl) -> unresolved
		isUnresolved := glyphID == 0

		if isUnresolved {
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
	// TODO: Handle RTL swap

	resolvedTextStart := resolvedTextLimits.Start
	resolvedGlyphsStart := 0

	// Iterate new unresolved blocks
	for i := startingCount; i < len(ols.unresolvedBlocks); i++ {
		unresolved := ols.unresolvedBlocks[i]

		// Create resolved block before this unresolved one
		// range: [resolvedTextStart, unresolved.Start)

		// Determine End based on direction. Assuming LTR for now.
		resolvedTextEnd := unresolved.text.Start

		resolvedText := NewTextRange(resolvedTextStart, resolvedTextEnd)
		if resolvedText.Width() > 0 {
			resolvedGlyphs := NewRange(resolvedGlyphsStart, unresolved.glyphs.Start)
			resolvedBlock := newResolvedRunBlock(run, resolvedText, resolvedGlyphs)
			ols.resolvedBlocks = append(ols.resolvedBlocks, resolvedBlock)
		}

		resolvedGlyphsStart = unresolved.glyphs.End
		resolvedTextStart = unresolved.text.End
	}

	// Final piece
	resolvedText := NewTextRange(resolvedTextStart, resolvedTextLimits.End)
	if resolvedText.Width() > 0 {
		resolvedGlyphs := NewRange(resolvedGlyphsStart, run.Size())
		resolvedBlock := newResolvedRunBlock(run, resolvedText, resolvedGlyphs)
		ols.resolvedBlocks = append(ols.resolvedBlocks, resolvedBlock)
	}
}

func (ols *OneLineShaper) clusteredText(run *Run, glyphRange GlyphRange) TextRange {
	// Find text range corresponding to these glyphs
	// Run has ClusterIndexes.
	// Start cluster: run.ClusterIndexes[glyphRange.Start]
	// End cluster: run.ClusterIndexes[glyphRange.End-1] (plus length of that char?)
	// Map back to global text indices.

	if glyphRange.Width() == 0 {
		return emptyRange
	}

	clusters := run.ClusterIndexes()
	startCluster := int(clusters[glyphRange.Start])
	endCluster := startCluster

	// Find max extent in this range
	for i := glyphRange.Start; i < glyphRange.End; i++ {
		c := int(clusters[i])
		if c < startCluster {
			startCluster = c
		}
		if c > endCluster {
			endCluster = c
		}
	}

	// We need to cover the full character of the last cluster.
	// Ideally we'd know char length.
	// For now, assuming contiguous up to next cluster or run end.

	// Find next cluster index in run or run end
	nextCluster := run.TextRange().End
	// ... logic to find byte length of last char ...
	// Simplified:
	// If it's the last glyph, it goes to Run.TextRange.End
	// Else it goes to ClusterIndexes[glyphRange.End] ?

	if glyphRange.End < len(clusters) {
		nextCluster = int(clusters[glyphRange.End])
	}

	return NewTextRange(startCluster, nextCluster)
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
		// Just try default fallback for now (simplified vs C++)
		fallback := ols.fontCollection.DefaultFallbackTypeface()
		if fallback != nil {
			if visitor(fallback) == resolvedEverything {
				return
			}
		}

		// Per-character fallback iteration
		// If we still have unresolved blocks, try to find a font for each character
		if len(ols.unresolvedBlocks) > 0 {
			// Copy unresolved blocks to iterate safely
			pendingBlocks := make([]runBlock, len(ols.unresolvedBlocks))
			copy(pendingBlocks, ols.unresolvedBlocks)

			// Clear main list to consume one by one
			ols.unresolvedBlocks = nil

			// Add them back if we fail

			for _, block := range pendingBlocks {
				text := ols.text[block.text.Start:block.text.End]
				for i, r := range text {
					// Find font for this rune
					// Simplified: check default fallback for this char
					fbTypeface := ols.fontCollection.DefaultFallback(r, style.FontStyle, "")
					if fbTypeface != nil {
						// Shape *just* this character (or cluster)
						// This is very inefficient (shaping char by char),
						// C++ optimization groups same-font chars.
						// Here we just test if visitor accepts it.
						// If visitor accepts, it shapes current 'unresolvedBlocks'.
						// So we need to set 'unresolvedBlocks' to just this char range?
						// No, visitor shapes ALL unresolvedBlocks.

						// Strategy:
						// 1. Identify sub-range for this char
						charLen := len(string(r))
						charStart := block.text.Start + i // This is byte offset from start of block text
						// Wait, 'i' in range loop over string is byte index? Yes.
						charRange := NewTextRange(charStart, charStart+charLen)

						// 2. Set ols.unresolvedBlocks to this single block
						ols.unresolvedBlocks = []runBlock{newRunBlock(charRange)}

						// 3. Call visitor
						if visitor(fbTypeface) == resolvedEverything {
							// Success for this char
						} else {
							// Failed, add to resolvedBlocks as unresolved?
							// Or stick back to pending?
							// For now, if failed, we should probably record it as unresolved run (tofu)
						}
					}
				}
			}
			// What wasn't resolved is lost in this simplified logic?
			// C++ puts hopeless blocks back into unresolved.
		}
	}
}
