# SkShaper Contract & Semantics

This document defines the semantic contract for the `SkShaper::Shape` method and the expected behaviors of its collaborators (`RunIterator`s and `RunHandler`). Any implementation of the `SkShaper` interface (whether wrapping C++ Skia or a pure Go implementation) MUST adhere to these rules.

## The `Shape` Method

The core entry point is:

```go
func (s *Shaper) Shape(
    text string, 
    fontIter FontRunIterator, 
    bidiIter BiDiRunIterator, 
    scriptIter ScriptRunIterator, 
    langIter LanguageRunIterator,
    features []Feature,
    width float32, 
    handler RunHandler,
)
```

### Lifecycle

1.  **Input Validation**: The shaper validates that the text and iterators are consistent (e.g., iterator lengths match text length).
2.  **Run Resolution**: The shaper iterates through the text, finding the intersection of all four iterators (Font, BiDi, Script, Language). Each intersection forms a "Run".
3.  **Shaping**: For each Run:
    *   The corresponding `Font`, `BiDiLevel`, `Script`, and `Language` are resolved.
    *   Glyphs are generated for the text in that run using the resolved Font and features.
4.  **Handling**: The shaper passes the shaped results to the `RunHandler`.

---

## RunIterator Advancement

Iterators partition the text into contiguous ranges sharing a specific property (e.g., "this range is Bold", "this range is Arabic").

*   **Intersection**: The Shaper advances all iterators simultaneously. The next break point is determined by `min(fontIter.EndOfCurrentRun(), bidiIter.EndOfCurrentRun(), ...)`.
*   **Consumption**: Once a run is processed up to a certain index, all iterators are expected to satisfy `CurrentRun().End > current_index`. If an iterator's current run ends exactly at `current_index`, `Consume()` is called to move it to the next run.
*   **Invariant**: At any point `i` in the text, all iterators must yield valid metadata for `text[i]`.

---

## RunHandler Call Sequence

The `RunHandler` acts as a state machine. The Shaper drives it via the following sequence:

### 1. `BeginLine()`
Called once at the very beginning of the `Shape` call.
*   **Purpose**: signals the start of processing.

### 2. `RunInfo(info RunInfo) -> (buffer RunBuffer, err error)`
Called for each resolved Run *before* glyphs are irrevocably written.
*   **Input**: `RunInfo` containing the `Font`, `BidiLevel`, `Advance` (total width), and count of glyphs.
*   **Purpose**: "Here is what I plan to write."
*   **Return**: The Handler must invoke `CommitRunInfo()`.
    *   Historically in Skia C++, `RunInfo` returns nothing or just informs internal state. In Go adaptation, this steps acts as the precursor to buffer allocation.

### 3. `CommitRunInfo()`
*   **Purpose**: Acknowledges the `RunInfo`. The handler may update internal line breaking state here.

### 4. `RunBuffer(info RunInfo) -> Buffer`
Called after `RunInfo` / `CommitRunInfo` for the same run.
*   **Purpose**: Requesting memory to write the actual glyph data.
*   **Return**: A `Buffer` struct containing slices for `Glyphs`, `Positions`, `Offsets`, and `Clusters`.
    *   **Crucial**: The Shaper WILL write directly into these slices. They must be of sufficient length (`info.GlyphCount`).
    *   `Glyphs`: `[]uint16` (Glyph IDs)
    *   `Positions`: `[]Point` (x, y coordinates relative to the line origin)
    *   `Offsets`: `[]Point` (optional per-glyph offsets)
    *   `Clusters`: `[]uint32` (indices into the original UTF-8 text)

### 5. `CommitRunBuffer(info RunInfo)`
Called after the Shaper has filled the requested `Buffer`.
*   **Purpose**: "I am done writing to the buffer you gave me."
*   **Side Effect**: The Handler now owns the data in the buffer and can accumulate it into the final formatted line.

### 6. `CommitLine()`
Called once at the very end of processing.
*   **Purpose**: Finalizes the shape operation. The Handler can now package the result (e.g., into a `TextBlob` or `Paragraph` layout).

---

## Buffer Management Expectations

*   **Allocation Responsibility**: The `RunHandler` is responsible for allocating the backing arrays for Glyphs, Positions, etc.
*   **Safety**: The Shaper assumes the pointers/slices returned by `RunBuffer` are valid for the duration of the write operations immediately following.
*   **Indices**: Note that `Clusters` indices refer to the byte offset in the original UTF-8 string. They must be monotonically increasing or decreasing depending on the BiDi level, but within a single run, the logic is strictly defined by the font's internal shaping logic (e.g., HarfBuzz).
