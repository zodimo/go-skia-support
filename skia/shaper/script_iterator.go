package shaper

import (
	"encoding/binary"
	"unicode"
)

const (
	// ScriptCommon is the tag for Common script (Zyyy)
	ScriptCommon = 0x5A797979
	// ScriptInherited is the tag for Inherited script (Zinh)
	ScriptInherited = 0x5A696E68
	// ScriptLatin is the tag for Latin script (Latn)
	ScriptLatin = 0x4C61746E
)

// scriptRun represents a resolved script run
type scriptRun struct {
	script uint32
	end    int
}

type scriptRunIterator struct {
	text        []rune
	totalLength int // in bytes
	byteOffset  int
	runs        []scriptRun
	currentIdx  int
}

// NewScriptRunIterator creates a new ScriptRunIterator.
func NewScriptRunIterator(text string, length int) ScriptRunIterator {
	runes := []rune(text)
	runs := computeScriptRuns(runes, text)
	return &scriptRunIterator{
		text:        runes,
		totalLength: length,
		runs:        runs,
		currentIdx:  0,
	}
}

func (s *scriptRunIterator) CurrentScript() uint32 {
	if s.currentIdx >= len(s.runs) {
		return ScriptCommon
	}
	return s.runs[s.currentIdx].script
}

func (s *scriptRunIterator) EndOfCurrentRun() int {
	if s.currentIdx >= len(s.runs) {
		return s.totalLength
	}
	return s.runs[s.currentIdx].end
}

func (s *scriptRunIterator) Consume() {
	if s.currentIdx < len(s.runs) {
		s.currentIdx++
	}
}

func (s *scriptRunIterator) AtEnd() bool {
	return s.currentIdx >= len(s.runs)
}

// computeScriptRuns segments text into runs of the same script, accounting for Common/Inherited.
func computeScriptRuns(runes []rune, text string) []scriptRun {
	if len(runes) == 0 {
		return nil
	}

	// 1. Map each rune to a script tag
	tags := make([]uint32, len(runes))
	for i, r := range runes {
		tags[i] = getScriptTag(r)
	}

	// 2. Resolve Common/Inherited scripts
	// Forward pass: resolve Common/Inherited to previous script
	currentParamsScript := uint32(ScriptCommon)
	lastScript := currentParamsScript
	for i := 0; i < len(tags); i++ {
		if tags[i] == ScriptCommon || tags[i] == ScriptInherited {
			if lastScript != ScriptCommon && lastScript != ScriptInherited {
				tags[i] = lastScript
			}
		} else {
			lastScript = tags[i]
		}
	}

	// Backward pass
	for i := len(tags) - 1; i >= 0; i-- {
		if tags[i] == ScriptCommon || tags[i] == ScriptInherited {
			if i+1 < len(tags) {
				tags[i] = tags[i+1]
			}
		}
	}

	// 3. Final pass
	for i := 0; i < len(tags); i++ {
		if tags[i] == ScriptCommon || tags[i] == ScriptInherited {
			tags[i] = ScriptLatin
		}
	}

	// 4. Merge identical sequential tags into runs
	var runs []scriptRun
	if len(tags) == 0 {
		return runs
	}

	// Need to map back to byte offsets
	byteOffsets := make([]int, len(runes)+1)
	currByte := 0
	for i, r := range runes {
		byteOffsets[i] = currByte
		currByte += len(string(r))
	}
	byteOffsets[len(runes)] = currByte

	currentTag := tags[0]

	for i := 1; i < len(tags); i++ {
		if tags[i] != currentTag {
			runs = append(runs, scriptRun{
				script: currentTag,
				end:    byteOffsets[i],
			})
			currentTag = tags[i]
		}
	}
	// Final run
	runs = append(runs, scriptRun{
		script: currentTag,
		end:    byteOffsets[len(runes)],
	})

	return runs
}

func getScriptTag(r rune) uint32 {
	// Optimization: check Latin first
	if unicode.Is(unicode.Latin, r) {
		return ScriptLatin
	}
	// Common (Zyyy) check - usually punctuations, numbers
	if unicode.Is(unicode.Common, r) {
		return ScriptCommon
	}
	// Inherited (Zinh) - usually combining marks
	if unicode.Is(unicode.Inherited, r) {
		return ScriptInherited
	}

	// Iterate over Scripts map for others
	// This is slow (O(NumScripts)). Ideally we'd use a trie or interval tree.
	// For "Production Ready" correctness, we use accurate map.
	// For performance, we can optimize later or use x/text/unicode/script if available.

	// Check known scripts table
	for name, table := range unicode.Scripts {
		if unicode.Is(table, r) {
			return scriptNameToTag(name)
		}
	}

	return ScriptCommon
}

// scriptNameToTag converts "Latin" -> "Latn" tag (uint32 big-endian).
// We simplify by using a map for common ones and algorithmic fallback if possible.
func scriptNameToTag(name string) uint32 {
	if tag, ok := scriptMap[name]; ok {
		return tag
	}
	// Fallback? or generate 4-char code?
	// Most script codes are not just first 4 chars.
	// e.g. "Canadian_Aboriginal" -> "Cans"
	return ScriptCommon
}

// Helper to make FourCC
func makeTag(s string) uint32 {
	if len(s) != 4 {
		return 0
	}
	return binary.BigEndian.Uint32([]byte(s))
}

var scriptMap = map[string]uint32{
	"Latin":      ScriptLatin,
	"Greek":      makeTag("Grek"),
	"Cyrillic":   makeTag("Cyrl"),
	"Arabic":     makeTag("Arab"),
	"Hebrew":     makeTag("Hebr"),
	"Han":        makeTag("Hani"),
	"Hiragana":   makeTag("Hira"),
	"Katakana":   makeTag("Kana"),
	"Hangul":     makeTag("Hang"),
	"Thai":       makeTag("Thai"),
	"Devanagari": makeTag("Deva"),
	"Bengali":    makeTag("Beng"),
	"Gurmukhi":   makeTag("Guru"),
	"Gujarati":   makeTag("Gujr"),
	"Oriya":      makeTag("Orya"),
	"Tamil":      makeTag("Taml"),
	"Telugu":     makeTag("Telu"),
	"Kannada":    makeTag("Knda"),
	"Malayalam":  makeTag("Mlym"),
	"Sinhala":    makeTag("Sinh"),
	"Myanmar":    makeTag("Mymr"),
	"Khmer":      makeTag("Khmr"),
	"Lao":        makeTag("Laoo"),
	"Tibetan":    makeTag("Tibt"),
	"Georgian":   makeTag("Geor"),
	"Armenian":   makeTag("Armn"),
	"Braille":    makeTag("Brai"),
	"Common":     ScriptCommon,
	"Inherited":  ScriptInherited,
}
