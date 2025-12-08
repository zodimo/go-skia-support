# Epic 2: Text Support - Brownfield Enhancement

**Epic Goal:** Add text rendering capabilities to the Go Skia port by implementing SkTextBlob, SkFont, and SkTypeface interfaces, enabling text drawing through the Canvas API.

---

## Epic Description

### Existing System Context

- **Current Functionality:** Go port of Skia's core types (`SkMatrix`, `SkPaint`, `SkPath`) with implementations complete. Canvas interface supports drawing primitives (rectangles, paths, ovals, arcs) but lacks text rendering capabilities.
- **Technology Stack:** Pure Go, backend-agnostic calculation library following BYOG (Bring Your Own Graphics Backend) philosophy
- **Integration Points:** 
  - Canvas interface at `skia/interfaces/canvas.go` (currently missing text methods)
  - Reference C++ Skia source at `/home/jaco/SecondBrain/1-Projects/GoProjects/Development/skia-source`
  - C++ API reference: `include/core/SkCanvas.h` (`drawTextBlob` method at line ~2008)
  - C++ types: `include/core/SkTextBlob.h`, `include/core/SkFont.h`, `include/core/SkTypeface.h`

### Enhancement Details

**What's being added:** 
- `SkTextBlob` interface and implementation - immutable container for glyphs, positions, and font attributes
- `SkFont` interface and implementation - lightweight object holding typeface reference plus styling (size, scale, skew, edging, hinting)
- `SkTypeface` interface - represents font data (immutable font face)
- **Canvas text drawing methods** (three methods matching C++ API):
  - `DrawTextBlob` - draws pre-constructed text blob at specified position
  - `DrawSimpleText` - draws text directly with encoding, byte length, font, and paint
  - `DrawString` - convenience method for drawing UTF-8 strings (wraps DrawSimpleText)
- Supporting enums: `TextEncoding`, `FontEdging`, `FontHinting`, `FontStyle` (weight, width, slant)
- Helper functions for creating text blobs from strings

**How it integrates:** 
- New interfaces follow existing pattern: `skia/interfaces/textblob.go`, `skia/interfaces/font.go`, `skia/interfaces/typeface.go`
- Implementations in `skia/impl/textblob.go`, `skia/impl/font.go`, `skia/impl/typeface.go`
- Canvas interface extended with three text drawing methods:
  - `DrawTextBlob(blob SkTextBlob, x Scalar, y Scalar, paint SkPaint)` - primary method for pre-constructed blobs
  - `DrawSimpleText(text []byte, encoding enums.TextEncoding, x Scalar, y Scalar, font SkFont, paint SkPaint)` - low-level text drawing
  - `DrawString(str string, x Scalar, y Scalar, font SkFont, paint SkPaint)` - convenience method for UTF-8 strings
- Enums added to `skia/enums/` package
- Follows existing code patterns: interface-based design, backend-agnostic calculations, comprehensive test coverage

**Success criteria:**
- Canvas interface includes all three text drawing methods matching C++ API signatures:
  - `DrawTextBlob` (line ~2008 in C++ API)
  - `DrawSimpleText` (line ~1834 in C++ API)
  - `DrawString` (line ~1861 in C++ API)
- `SkTextBlob` interface matches C++ `SkTextBlob` API (bounds, iteration, creation from strings)
- `SkFont` interface matches C++ `SkFont` API (typeface, size, scale, skew, edging, hinting)
- `SkTypeface` interface provides font data access (basic interface for MVP)
- Text blob creation from UTF-8 strings works correctly
- Direct text drawing via `DrawSimpleText` and `DrawString` works correctly
- Bounds calculation for text blobs is accurate
- Integration with existing Paint interface works (color, blend modes, filters apply to text)
- Test coverage matches existing component standards

---

## Stories

### Story 1: Text Blob Interface and Core Implementation
**Priority:** High  
**Description:** Implement `SkTextBlob` interface and core implementation following Skia C++ API. Includes blob creation from strings, bounds calculation, and basic iteration capabilities.

**Key Deliverables:**
- `SkTextBlob` interface in `skia/interfaces/textblob.go` matching C++ API:
  - `Bounds()` method returning bounding rectangle
  - `MakeFromString(text string, font SkFont)` factory method
  - Basic blob structure supporting single run of glyphs
- `SkTextBlob` implementation in `skia/impl/textblob.go`:
  - Internal storage for glyphs, positions, and font reference
  - Bounds calculation from glyph metrics
  - UTF-8 string to glyph conversion (basic implementation)
- Supporting enums: `TextEncoding` (UTF8, UTF16, UTF32, GlyphID)
- Test file `textblob_test.go` with:
  - Bounds calculation tests
  - String-to-blob conversion tests
  - Empty/invalid input handling

### Story 2: Font and Typeface Interfaces
**Priority:** High  
**Description:** Implement `SkFont` and `SkTypeface` interfaces to provide font configuration and font data access. Font holds typeface reference plus styling attributes.

**Key Deliverables:**
- `SkFont` interface in `skia/interfaces/font.go` matching C++ API:
  - `Typeface()` method returning SkTypeface
  - `Size()`, `ScaleX()`, `SkewX()` getters
  - `Edging()`, `Hinting()` getters
  - Setter methods for all attributes
- `SkFont` implementation in `skia/impl/font.go`:
  - Typeface reference storage
  - Font attribute storage and accessors
  - Default font creation
- `SkTypeface` interface in `skia/interfaces/typeface.go`:
  - Basic interface for MVP (can be extended later)
  - `FontStyle()` method returning font style information
- `SkTypeface` implementation in `skia/impl/typeface.go`:
  - Basic typeface structure (placeholder for system font integration)
  - Default typeface creation
- Supporting enums: `FontEdging` (Alias, AntiAlias, SubpixelAntiAlias), `FontHinting` (None, Slight, Normal, Full), `FontStyle` (weight, width, slant)
- Test files `font_test.go` and `typeface_test.go` with attribute tests

### Story 3: Canvas Text Drawing Integration
**Priority:** High  
**Description:** Add all three text drawing methods to Canvas interface (`DrawTextBlob`, `DrawSimpleText`, `DrawString`) and ensure proper integration with existing Paint system. Text rendering respects canvas transformations, clipping, and paint properties.

**Key Deliverables:**
- **Three text drawing methods** added to `SkCanvas` interface in `skia/interfaces/canvas.go`:
  1. `DrawTextBlob(blob SkTextBlob, x Scalar, y Scalar, paint SkPaint)` - Primary method for pre-constructed text blobs
     - Ported from C++ API line ~2008
     - Documentation matching C++ API comments
  2. `DrawSimpleText(text []byte, encoding enums.TextEncoding, x Scalar, y Scalar, font SkFont, paint SkPaint)` - Low-level text drawing
     - Ported from C++ API line ~1834
     - Supports multiple text encodings (UTF8, UTF16, UTF32, GlyphID)
     - Takes raw byte array with encoding specification
  3. `DrawString(str string, x Scalar, y Scalar, font SkFont, paint SkPaint)` - Convenience method for UTF-8 strings
     - Ported from C++ API line ~1861
     - Wraps `DrawSimpleText` with UTF-8 encoding
     - Most common use case for simple text rendering
- Integration verification:
  - All three methods respect canvas matrix transformations
  - All three methods respect canvas clipping regions
  - Paint properties apply correctly (color, blend mode, filters)
  - Text positioning (x, y offset) works correctly
  - `DrawString` correctly delegates to `DrawSimpleText` with UTF-8 encoding
- Test file `canvas_text_test.go` with:
  - Basic text drawing tests for all three methods
  - Transformation tests (translate, scale, rotate with text)
  - Clipping tests with text
  - Paint property tests (color, blend modes)
  - Encoding tests for `DrawSimpleText` (UTF8, UTF16, UTF32)
  - Verification that `DrawString` produces same results as `DrawSimpleText` with UTF-8
- Documentation updates to reflect new Canvas text rendering capabilities

---

## Compatibility Requirements

- [x] Existing APIs remain unchanged (text is additive feature)
- [x] Canvas interface extension is backward compatible (new method, no breaking changes)
- [x] No changes to existing Paint, Path, or Matrix interfaces
- [x] Text implementation follows existing code patterns and file organization
- [x] Performance impact is minimal (text calculations are pure Go, no backend dependencies)

---

## Risk Mitigation

**Primary Risk:** Text rendering involves complex glyph shaping and font metrics that may require system font access or font data parsing, which could complicate the backend-agnostic design.

**Mitigation:** 
- MVP focuses on basic text blob creation from strings with simple glyph positioning
- Font metrics calculation can use simplified algorithms initially
- Typeface interface is designed to be extensible for system font integration later
- Text blob contains pre-calculated glyph positions, keeping rendering backend-agnostic

**Rollback Plan:** 
- Text interfaces and implementations are isolated in separate files
- Canvas interface addition is non-breaking (new method only)
- Can disable text features via feature flag if needed
- Existing drawing primitives remain unaffected

---

## Definition of Done

- [ ] All three stories completed with acceptance criteria met
- [ ] Canvas interface includes all three text drawing methods (`DrawTextBlob`, `DrawSimpleText`, `DrawString`)
- [ ] `SkTextBlob`, `SkFont`, and `SkTypeface` interfaces match C++ API structure
- [ ] Text blob creation from UTF-8 strings works correctly
- [ ] Bounds calculation for text blobs is accurate
- [ ] Text respects canvas transformations and clipping
- [ ] Paint properties (color, blend modes) apply correctly to text
- [ ] Test coverage matches existing component standards (comprehensive unit tests)
- [ ] Code follows existing patterns and coding standards
- [ ] Documentation updated to reflect text rendering capabilities
- [ ] No regression in existing features (all existing tests pass)

---

## Technical Notes

**C++ Reference Types:**
- `SkTextBlob` - Primary text container (confirmed correct type)
- `SkFont` - Font configuration (typeface + styling)
- `SkTypeface` - Font data representation
- `SkFontMgr` - Font manager (out of scope for MVP, can be added later)

**Key C++ API Methods:**
- `SkTextBlob::MakeFromString(text, font)` - Create blob from string
- `SkTextBlob::bounds()` - Get bounding rectangle
- `SkCanvas::drawTextBlob(blob, x, y, paint)` - Draw pre-constructed text blob (line ~2008)
- `SkCanvas::drawSimpleText(text, byteLength, encoding, x, y, font, paint)` - Low-level text drawing (line ~1834)
- `SkCanvas::drawString(str, x, y, font, paint)` - Convenience method for UTF-8 strings (line ~1861)
- `SkFont` constructor and attribute setters/getters

**Implementation Approach:**
- Start with basic single-run text blobs (one font, one text string)
- Use simplified glyph positioning initially (left-to-right, no complex shaping)
- Font metrics can use approximations for MVP
- Extend to multi-run blobs and advanced features in future epics


