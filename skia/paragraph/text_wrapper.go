package paragraph

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/models"
)

// ClusterPos represents a position within a cluster range.
type ClusterPos struct {
	clusterIndex int
	position     int
}

// NewClusterPos creates a new ClusterPos.
func NewClusterPos(clusterIndex, position int) ClusterPos {
	return ClusterPos{clusterIndex: clusterIndex, position: position}
}

// Clean resets the ClusterPos.
func (cp *ClusterPos) Clean() {
	cp.clusterIndex = -1
	cp.position = 0
}

// Move advances or retreats the cluster position.
func (cp *ClusterPos) Move(up bool, owner TextLineOwner) {
	if up {
		cp.clusterIndex++
		cp.position = 0
	} else {
		cp.clusterIndex--
		if cp.clusterIndex >= 0 && owner != nil {
			cluster := owner.Cluster(cp.clusterIndex)
			if cluster != nil {
				cp.position = cluster.EndPos()
			}
		}
	}
}

// TextStretch represents a stretch of text being processed.
type TextStretch struct {
	start                ClusterPos
	end                  ClusterPos
	breakPos             ClusterPos
	metrics              InternalLineMetrics
	width                float32
	widthWithGhostSpaces float32
}

// NewTextStretch creates a new empty TextStretch.
func NewTextStretch() TextStretch {
	return TextStretch{
		start: ClusterPos{clusterIndex: -1},
		end:   ClusterPos{clusterIndex: -1},
	}
}

// NewTextStretchFromClusters creates a TextStretch from cluster range.
func NewTextStretchFromClusters(owner TextLineOwner, startIdx, endIdx int, forceStrut bool) TextStretch {
	ts := TextStretch{
		start:   ClusterPos{clusterIndex: startIdx, position: 0},
		end:     ClusterPos{clusterIndex: endIdx, position: 0},
		metrics: NewInternalLineMetrics(),
	}
	ts.metrics.ForceStrut = forceStrut

	// Set end position
	if endIdx >= 0 && owner != nil {
		endCluster := owner.Cluster(endIdx)
		if endCluster != nil {
			ts.end.position = endCluster.EndPos()
		}
	}

	// Accumulate metrics and width
	for i := startIdx; i <= endIdx; i++ {
		cluster := owner.Cluster(i)
		if cluster == nil {
			continue
		}
		run := cluster.Run()
		if run != nil {
			ts.metrics.AddRun(run)
		}
		if i < endIdx {
			ts.width += cluster.Width()
		}
	}
	ts.widthWithGhostSpaces = ts.width
	return ts
}

// Width returns the width.
func (ts *TextStretch) Width() float32 {
	return ts.width
}

// WidthWithGhostSpaces returns the width including trailing spaces.
func (ts *TextStretch) WidthWithGhostSpaces() float32 {
	return ts.widthWithGhostSpaces
}

// StartClusterIndex returns the start cluster index.
func (ts *TextStretch) StartClusterIndex() int {
	return ts.start.clusterIndex
}

// EndClusterIndex returns the end cluster index.
func (ts *TextStretch) EndClusterIndex() int {
	return ts.end.clusterIndex
}

// BreakClusterIndex returns the break cluster index.
func (ts *TextStretch) BreakClusterIndex() int {
	return ts.breakPos.clusterIndex
}

// Metrics returns the metrics.
func (ts *TextStretch) Metrics() *InternalLineMetrics {
	return &ts.metrics
}

// StartPos returns the start position.
func (ts *TextStretch) StartPos() int {
	return ts.start.position
}

// EndPos returns the end position.
func (ts *TextStretch) EndPos() int {
	return ts.end.position
}

// EndOfCluster returns true if at end of cluster.
func (ts *TextStretch) EndOfCluster(owner TextLineOwner) bool {
	if ts.end.clusterIndex < 0 || owner == nil {
		return false
	}
	cluster := owner.Cluster(ts.end.clusterIndex)
	if cluster == nil {
		return false
	}
	return ts.end.position == cluster.EndPos()
}

// EndOfWord returns true if at end of word.
func (ts *TextStretch) EndOfWord(owner TextLineOwner) bool {
	if !ts.EndOfCluster(owner) {
		return false
	}
	cluster := owner.Cluster(ts.end.clusterIndex)
	if cluster == nil {
		return false
	}
	return cluster.IsHardBreak() || cluster.IsSoftBreak()
}

// Extend extends by another TextStretch.
func (ts *TextStretch) Extend(other *TextStretch) {
	ts.metrics.Add(other.metrics)
	ts.end = other.end
	ts.width += other.width
	other.Clean()
}

// ExtendCluster extends by a single cluster.
func (ts *TextStretch) ExtendCluster(owner TextLineOwner, clusterIdx int) {
	cluster := owner.Cluster(clusterIdx)
	if cluster == nil {
		return
	}

	if ts.start.clusterIndex < 0 {
		ts.start = ClusterPos{clusterIndex: clusterIdx, position: cluster.StartPos()}
	}
	ts.end = ClusterPos{clusterIndex: clusterIdx, position: cluster.EndPos()}

	run := cluster.Run()
	if run != nil && !cluster.IsHardBreak() && !run.IsPlaceholder() {
		ts.metrics.AddRun(run)
	}
	ts.width += cluster.Width()
}

// ExtendClusterWithPos extends to a specific position.
func (ts *TextStretch) ExtendClusterWithPos(owner TextLineOwner, clusterIdx, pos int) {
	ts.end = ClusterPos{clusterIndex: clusterIdx, position: pos}
	cluster := owner.Cluster(clusterIdx)
	if cluster != nil {
		run := cluster.Run()
		if run != nil {
			ts.metrics.AddRun(run)
		}
	}
}

// StartFrom starts from a cluster.
func (ts *TextStretch) StartFrom(owner TextLineOwner, clusterIdx, pos int) {
	ts.start = ClusterPos{clusterIndex: clusterIdx, position: pos}
	ts.end = ClusterPos{clusterIndex: clusterIdx, position: pos}

	cluster := owner.Cluster(clusterIdx)
	if cluster != nil {
		run := cluster.Run()
		if run != nil && !run.IsPlaceholder() {
			ts.metrics.AddRun(run)
		}
	}
	ts.width = 0
}

// SaveBreak saves the current break point.
func (ts *TextStretch) SaveBreak() {
	ts.widthWithGhostSpaces = ts.width
	ts.breakPos = ts.end
}

// RestoreBreak restores to the saved break point.
func (ts *TextStretch) RestoreBreak() {
	ts.width = ts.widthWithGhostSpaces
	ts.end = ts.breakPos
}

// ShiftBreak shifts the break forward.
func (ts *TextStretch) ShiftBreak(owner TextLineOwner) {
	ts.breakPos.Move(true, owner)
}

// Trim trims trailing whitespace.
func (ts *TextStretch) Trim(owner TextLineOwner) {
	if ts.end.clusterIndex < 0 || owner == nil {
		return
	}
	cluster := owner.Cluster(ts.end.clusterIndex)
	if cluster == nil || ts.width <= 0 {
		return
	}
	run := cluster.Run()
	if run != nil && run.IsPlaceholder() {
		return
	}
	ts.width -= cluster.Width() - cluster.TrimmedWidth(ts.end.position)
}

// TrimCluster trims a specific cluster.
func (ts *TextStretch) TrimCluster(owner TextLineOwner, clusterIdx int) {
	if ts.end.clusterIndex != clusterIdx {
		return
	}
	cluster := owner.Cluster(clusterIdx)
	if cluster == nil {
		return
	}

	if ts.end.clusterIndex > ts.start.clusterIndex {
		ts.end.Move(false, owner)
		ts.width -= cluster.Width()
	} else {
		ts.end.position = ts.start.position
		ts.width = 0
	}
}

// Empty returns true if stretch is empty.
func (ts *TextStretch) Empty() bool {
	return ts.start.clusterIndex == ts.end.clusterIndex &&
		ts.start.position == ts.end.position
}

// SetMetrics sets the metrics.
func (ts *TextStretch) SetMetrics(m InternalLineMetrics) {
	ts.metrics = m
}

// Clean resets the stretch.
func (ts *TextStretch) Clean() {
	ts.start.Clean()
	ts.end.Clean()
	ts.breakPos.Clean()
	ts.width = 0
	ts.widthWithGhostSpaces = 0
	ts.metrics.Clean()
}

// TextWrapper wraps text into lines.
//
// Ported from: skia-source/modules/skparagraph/src/TextWrapper.cpp
type TextWrapper struct {
	words    TextStretch
	clusters TextStretch
	clip     TextStretch
	endLine  TextStretch

	lineNumber     int
	tooLongWord    bool
	tooLongCluster bool

	hardLineBreak    bool
	exceededMaxLines bool

	height            float32
	minIntrinsicWidth float32
	maxIntrinsicWidth float32
}

// NewTextWrapper creates a new TextWrapper.
func NewTextWrapper() *TextWrapper {
	return &TextWrapper{
		lineNumber: 1,
	}
}

// Height returns the total height.
func (tw *TextWrapper) Height() float32 {
	return tw.height
}

// MinIntrinsicWidth returns the minimum intrinsic width.
func (tw *TextWrapper) MinIntrinsicWidth() float32 {
	return tw.minIntrinsicWidth
}

// MaxIntrinsicWidth returns the maximum intrinsic width.
func (tw *TextWrapper) MaxIntrinsicWidth() float32 {
	return tw.maxIntrinsicWidth
}

// ExceededMaxLines returns true if max lines were exceeded.
func (tw *TextWrapper) ExceededMaxLines() bool {
	return tw.exceededMaxLines
}

// AddLineToParagraph is the callback for adding a line.
type AddLineToParagraph func(
	textExcludingSpaces TextRange,
	text TextRange,
	textIncludingNewlines TextRange,
	clusters ClusterRange,
	clustersWithGhosts ClusterRange,
	widthWithSpaces float32,
	startClip, endClip int,
	offset, advance models.Point,
	metrics InternalLineMetrics,
	addEllipsis bool,
)

// reset resets the wrapper for a new line.
func (tw *TextWrapper) reset() {
	tw.words.Clean()
	tw.clusters.Clean()
	tw.clip.Clean()
	tw.tooLongCluster = false
	tw.tooLongWord = false
	tw.hardLineBreak = false
}

// LineBreakerWithLittleRounding handles rounding edge cases.
type LineBreakerWithLittleRounding struct {
	lower             float32
	maxWidth          float32
	upper             float32
	applyRoundingHack bool
}

// NewLineBreakerWithLittleRounding creates a new breaker.
func NewLineBreakerWithLittleRounding(maxWidth float32, applyRoundingHack bool) LineBreakerWithLittleRounding {
	return LineBreakerWithLittleRounding{
		lower:             maxWidth - 0.25,
		maxWidth:          maxWidth,
		upper:             maxWidth + 0.25,
		applyRoundingHack: applyRoundingHack,
	}
}

// BreakLine returns true if width exceeds max.
func (lb LineBreakerWithLittleRounding) BreakLine(width float32) bool {
	if width < lb.lower {
		return false
	} else if width > lb.upper {
		return true
	}

	val := float32(math.Abs(float64(width)))
	var roundedWidth float32
	if lb.applyRoundingHack {
		if val < 10000 {
			roundedWidth = float32(math.Round(float64(width*100))) * (1.0 / 100)
		} else if val < 100000 {
			roundedWidth = float32(math.Round(float64(width*10))) * (1.0 / 10)
		} else {
			roundedWidth = float32(math.Floor(float64(width)))
		}
	} else {
		if val < 10000 {
			roundedWidth = float32(math.Floor(float64(width*100))) * (1.0 / 100)
		} else if val < 100000 {
			roundedWidth = float32(math.Floor(float64(width*10))) * (1.0 / 10)
		} else {
			roundedWidth = float32(math.Floor(float64(width)))
		}
	}
	return roundedWidth > lb.maxWidth
}

// TextWrapperOwner is the interface for the paragraph that owns the wrapper.
type TextWrapperOwner interface {
	TextLineOwner
	Clusters() []*Cluster
	StructForceHeight() bool
	GetApplyRoundingHack() bool
	GetEmptyMetrics() InternalLineMetrics
	StrutEnabled() bool
	StrutMetrics() InternalLineMetrics
	Lines() []*TextLine
	Text() string
}

// BreakTextIntoLines breaks the text into lines.
func (tw *TextWrapper) BreakTextIntoLines(
	parent TextWrapperOwner,
	maxWidth float32,
	addLine AddLineToParagraph,
) {
	tw.height = 0
	tw.minIntrinsicWidth = float32(math.SmallestNonzeroFloat32)
	tw.maxIntrinsicWidth = float32(math.SmallestNonzeroFloat32)

	clusters := parent.Clusters()
	if len(clusters) == 0 {
		return
	}

	style := parent.ParagraphStyle()
	maxLines := style.MaxLines
	if maxLines == 0 {
		maxLines = math.MaxInt
	}
	align := style.TextAlign
	unlimitedLines := maxLines == math.MaxInt
	endlessLine := math.IsInf(float64(maxWidth), 1)
	hasEllipsis := style.Ellipsis != ""

	disableFirstAscent := style.TextHeightBehavior&TextHeightBehaviorDisableFirstAscent != 0
	disableLastDescent := style.TextHeightBehavior&TextHeightBehaviorDisableLastDescent != 0
	firstLine := true

	softLineMaxIntrinsicWidth := float32(0)
	tw.endLine = NewTextStretchFromClusters(parent, 0, 0, parent.StructForceHeight())

	endClusterIdx := len(clusters) - 1
	needEllipsis := false

	for tw.endLine.EndClusterIndex() < endClusterIdx {
		tw.lookAhead(parent, maxWidth, endClusterIdx, parent.GetApplyRoundingHack())

		lastLine := (hasEllipsis && unlimitedLines) || tw.lineNumber >= maxLines
		needEllipsis = hasEllipsis && !endlessLine && lastLine

		tw.moveForward(needEllipsis)
		if tw.endLine.EndClusterIndex() < endClusterIdx-1 {
			needEllipsis = needEllipsis && true
		} else {
			needEllipsis = false
		}

		tw.trimEndSpaces(parent, align)

		startLineIdx, pos, widthWithSpaces := tw.trimStartSpaces(parent, endClusterIdx)

		if needEllipsis && !tw.hardLineBreak {
			tw.endLine.RestoreBreak()
			widthWithSpaces = tw.endLine.WidthWithGhostSpaces()
		}

		if tw.endLine.Metrics().IsClean() {
			tw.endLine.SetMetrics(parent.GetEmptyMetrics())
		}

		// Handle placeholder runs
		lastRunIdx := -1
		for i := tw.endLine.StartClusterIndex(); i <= tw.endLine.EndClusterIndex(); i++ {
			cluster := parent.Cluster(i)
			if cluster == nil {
				continue
			}
			run := cluster.Run()
			if run == nil || run.Index() == lastRunIdx {
				continue
			}
			lastRunIdx = run.Index()
			if run.IsPlaceholder() {
				run.UpdateMetrics(tw.endLine.Metrics())
			}
		}

		maxRunMetrics := tw.endLine.metrics
		maxRunMetrics.ForceStrut = false

		// Calculate text ranges
		startCluster := parent.Cluster(tw.endLine.StartClusterIndex())
		endCluster := parent.Cluster(tw.endLine.EndClusterIndex())
		breakCluster := parent.Cluster(tw.endLine.BreakClusterIndex())
		startLineCluster := parent.Cluster(startLineIdx)

		var textExcludingSpaces, text, textIncludingNewlines TextRange
		if startCluster != nil && endCluster != nil {
			textExcludingSpaces = NewTextRange(startCluster.TextRange().Start, endCluster.TextRange().End)
		}
		if startCluster != nil && breakCluster != nil {
			text = NewTextRange(startCluster.TextRange().Start, breakCluster.TextRange().Start)
		}
		if startCluster != nil && startLineCluster != nil {
			textIncludingNewlines = NewTextRange(startCluster.TextRange().Start, startLineCluster.TextRange().Start)
		}

		if startLineIdx >= endClusterIdx {
			textIncludingNewlines.End = len(parent.Text())
			text.End = len(parent.Text())
		}

		clusterRange := NewClusterRange(tw.endLine.StartClusterIndex(), tw.endLine.EndClusterIndex()+1)
		clustersWithGhosts := NewClusterRange(tw.endLine.StartClusterIndex(), startLineIdx)

		if disableFirstAscent && firstLine {
			tw.endLine.metrics.Ascent = tw.endLine.metrics.RawAscent
		}
		if disableLastDescent && (lastLine || (startLineIdx >= endClusterIdx && !tw.hardLineBreak)) {
			tw.endLine.metrics.Descent = tw.endLine.metrics.RawDescent
		}

		if parent.StrutEnabled() {
			strutMetrics := parent.StrutMetrics()
			strutMetrics.UpdateLineMetrics(&tw.endLine.metrics)
		}

		lineHeight := tw.endLine.Metrics().Height()
		firstLine = false

		if tw.endLine.Empty() {
			textExcludingSpaces.End = textExcludingSpaces.Start
			clusterRange.End = clusterRange.Start
		}

		if text.End < textExcludingSpaces.End {
			text.End = textExcludingSpaces.End
		}

		addLine(
			textExcludingSpaces,
			text,
			textIncludingNewlines,
			clusterRange,
			clustersWithGhosts,
			widthWithSpaces,
			tw.endLine.StartPos(),
			tw.endLine.EndPos(),
			models.Point{X: 0, Y: models.Scalar(tw.height)},
			models.Point{X: models.Scalar(tw.endLine.Width()), Y: models.Scalar(lineHeight)},
			tw.endLine.metrics,
			needEllipsis && !tw.hardLineBreak,
		)

		softLineMaxIntrinsicWidth += widthWithSpaces
		if tw.maxIntrinsicWidth < softLineMaxIntrinsicWidth {
			tw.maxIntrinsicWidth = softLineMaxIntrinsicWidth
		}
		if tw.hardLineBreak {
			softLineMaxIntrinsicWidth = 0
		}

		tw.height += lineHeight
		if !tw.hardLineBreak || startLineIdx < endClusterIdx {
			tw.endLine.Clean()
		}
		tw.endLine.StartFrom(parent, startLineIdx, pos)

		if hasEllipsis && unlimitedLines {
			if !tw.hardLineBreak {
				break
			}
		} else if lastLine {
			tw.hardLineBreak = false
			break
		}

		tw.lineNumber++
	}

	// Scan remaining text for metrics
	if tw.endLine.EndClusterIndex() >= 0 {
		lastWordLength := float32(0)
		for i := tw.endLine.EndClusterIndex(); i <= endClusterIdx; i++ {
			tw.exceededMaxLines = true
			cluster := parent.Cluster(i)
			if cluster == nil {
				continue
			}

			if cluster.IsHardBreak() {
				if tw.maxIntrinsicWidth < softLineMaxIntrinsicWidth {
					tw.maxIntrinsicWidth = softLineMaxIntrinsicWidth
				}
				softLineMaxIntrinsicWidth = 0
				if tw.minIntrinsicWidth < lastWordLength {
					tw.minIntrinsicWidth = lastWordLength
				}
				lastWordLength = 0
			} else if cluster.IsWhitespaceBreak() {
				softLineMaxIntrinsicWidth += cluster.Width()
				if tw.minIntrinsicWidth < lastWordLength {
					tw.minIntrinsicWidth = lastWordLength
				}
				lastWordLength = 0
			} else {
				run := cluster.Run()
				if run != nil && run.IsPlaceholder() {
					if tw.minIntrinsicWidth < lastWordLength {
						tw.minIntrinsicWidth = lastWordLength
					}
					softLineMaxIntrinsicWidth += cluster.Width()
					if tw.minIntrinsicWidth < cluster.Width() {
						tw.minIntrinsicWidth = cluster.Width()
					}
					lastWordLength = 0
				} else {
					softLineMaxIntrinsicWidth += cluster.Width()
					lastWordLength += cluster.Width()
				}
			}
		}
		if tw.minIntrinsicWidth < lastWordLength {
			tw.minIntrinsicWidth = lastWordLength
		}
		if tw.maxIntrinsicWidth < softLineMaxIntrinsicWidth {
			tw.maxIntrinsicWidth = softLineMaxIntrinsicWidth
		}
	}

	// Handle trailing hard line break
	if tw.hardLineBreak {
		if disableLastDescent {
			tw.endLine.metrics.Descent = tw.endLine.metrics.RawDescent
		}

		if parent.StrutEnabled() {
			strutMetrics := parent.StrutMetrics()
			strutMetrics.UpdateLineMetrics(&tw.endLine.metrics)
		}

		breakCluster := parent.Cluster(tw.endLine.BreakClusterIndex())
		endCluster := parent.Cluster(tw.endLine.EndClusterIndex())

		var textRange TextRange
		if breakCluster != nil {
			textRange = breakCluster.TextRange()
		}

		clusterRange := NewClusterRange(tw.endLine.BreakClusterIndex(), tw.endLine.EndClusterIndex())

		var textIncludingNewlines TextRange
		if endCluster != nil {
			textIncludingNewlines = endCluster.TextRange()
		}

		addLine(
			textRange,
			textRange,
			textIncludingNewlines,
			clusterRange,
			clusterRange,
			0,
			0, 0,
			models.Point{X: 0, Y: models.Scalar(tw.height)},
			models.Point{X: 0, Y: models.Scalar(tw.endLine.Metrics().Height())},
			tw.endLine.metrics,
			needEllipsis,
		)
		tw.height += tw.endLine.Metrics().Height()
	}

	// Correct line metric styles
	lines := parent.Lines()
	if len(lines) == 0 {
		return
	}
	if disableFirstAscent {
		lines[0].ascentStyle = LineMetricStyleTypographic
	}
	if disableLastDescent {
		lines[len(lines)-1].descentStyle = LineMetricStyleTypographic
	}
}

// lookAhead looks ahead to find break opportunities.
func (tw *TextWrapper) lookAhead(parent TextWrapperOwner, maxWidth float32, endClusterIdx int, applyRoundingHack bool) {
	tw.reset()
	tw.endLine.Metrics().Clean()
	tw.words.StartFrom(parent, tw.endLine.StartClusterIndex(), tw.endLine.StartPos())
	tw.clusters.StartFrom(parent, tw.endLine.StartClusterIndex(), tw.endLine.StartPos())
	tw.clip.StartFrom(parent, tw.endLine.StartClusterIndex(), tw.endLine.StartPos())

	breaker := NewLineBreakerWithLittleRounding(maxWidth, applyRoundingHack)

	for i := tw.endLine.EndClusterIndex(); i <= endClusterIdx; i++ {
		cluster := parent.Cluster(i)
		if cluster == nil {
			continue
		}

		if cluster.IsHardBreak() {
			tw.hardLineBreak = true
			break
		}

		width := tw.words.Width() + tw.clusters.Width() + cluster.Width()
		if breaker.BreakLine(width) {
			if cluster.IsWhitespaceBreak() {
				tw.clusters.ExtendCluster(parent, i)
				trimmedWidth := tw.getClustersTrimmedWidth(parent)
				if tw.minIntrinsicWidth < trimmedWidth {
					tw.minIntrinsicWidth = trimmedWidth
				}
				tw.words.Extend(&tw.clusters)
				continue
			}

			run := cluster.Run()
			if run != nil && run.IsPlaceholder() {
				if !tw.clusters.Empty() {
					trimmedWidth := tw.getClustersTrimmedWidth(parent)
					if tw.minIntrinsicWidth < trimmedWidth {
						tw.minIntrinsicWidth = trimmedWidth
					}
					tw.words.Extend(&tw.clusters)
				}

				if cluster.Width() > maxWidth && tw.words.Empty() {
					tw.clusters.ExtendCluster(parent, i)
					tw.tooLongCluster = true
					tw.tooLongWord = true
				}
				break
			}

			// Check if word is too long
			nextWordLength := tw.clusters.Width()
			for j := i; j <= endClusterIdx; j++ {
				further := parent.Cluster(j)
				if further == nil {
					continue
				}
				if further.IsSoftBreak() || further.IsHardBreak() || further.IsWhitespaceBreak() {
					break
				}
				furtherRun := further.Run()
				if furtherRun != nil && furtherRun.IsPlaceholder() {
					break
				}
				if maxWidth == 0 {
					if nextWordLength < further.Width() {
						nextWordLength = further.Width()
					}
				} else {
					nextWordLength += further.Width()
				}
			}

			if nextWordLength > maxWidth {
				if tw.minIntrinsicWidth < nextWordLength {
					tw.minIntrinsicWidth = nextWordLength
				}
				if tw.clusters.EndPos()-tw.clusters.StartPos() > 1 || tw.words.Empty() {
					tw.tooLongWord = true
				}
			}

			if cluster.Width() > maxWidth {
				tw.clusters.ExtendCluster(parent, i)
				tw.tooLongCluster = true
				tw.tooLongWord = true
			}
			break
		}

		run := cluster.Run()
		if run != nil && run.IsPlaceholder() {
			if !tw.clusters.Empty() {
				trimmedWidth := tw.getClustersTrimmedWidth(parent)
				if tw.minIntrinsicWidth < trimmedWidth {
					tw.minIntrinsicWidth = trimmedWidth
				}
				tw.words.Extend(&tw.clusters)
			}

			if tw.minIntrinsicWidth < cluster.Width() {
				tw.minIntrinsicWidth = cluster.Width()
			}
			tw.words.ExtendCluster(parent, i)
		} else {
			tw.clusters.ExtendCluster(parent, i)

			if tw.clusters.EndOfWord(parent) {
				trimmedWidth := tw.getClustersTrimmedWidth(parent)
				if tw.minIntrinsicWidth < trimmedWidth {
					tw.minIntrinsicWidth = trimmedWidth
				}
				tw.words.Extend(&tw.clusters)
			}
		}

		if cluster.IsHardBreak() {
			tw.hardLineBreak = true
			break
		}
	}
}

// moveForward advances the line.
func (tw *TextWrapper) moveForward(hasEllipsis bool) {
	if !tw.words.Empty() {
		tw.endLine.Extend(&tw.words)
		if !tw.tooLongWord && !hasEllipsis {
			return
		}
	}
	if !tw.clusters.Empty() {
		tw.endLine.Extend(&tw.clusters)
		if !tw.tooLongCluster {
			return
		}
	}

	if !tw.clip.Empty() {
		tw.endLine.Metrics().Add(tw.clip.metrics)
	}
}

// trimEndSpaces trims trailing spaces.
func (tw *TextWrapper) trimEndSpaces(parent TextWrapperOwner, align TextAlign) {
	tw.endLine.SaveBreak()
	for i := tw.endLine.EndClusterIndex(); i >= tw.endLine.StartClusterIndex(); i-- {
		cluster := parent.Cluster(i)
		if cluster == nil || !cluster.IsWhitespaceBreak() {
			break
		}
		tw.endLine.TrimCluster(parent, i)
	}
	tw.endLine.Trim(parent)
}

// trimStartSpaces trims leading spaces for the next line.
func (tw *TextWrapper) trimStartSpaces(parent TextWrapperOwner, endClusterIdx int) (int, int, float32) {
	if tw.hardLineBreak {
		width := tw.endLine.Width()
		i := tw.endLine.EndClusterIndex() + 1
		for i <= tw.endLine.BreakClusterIndex() {
			cluster := parent.Cluster(i)
			if cluster == nil || !cluster.IsWhitespaceBreak() {
				break
			}
			width += cluster.Width()
			i++
		}
		return tw.endLine.BreakClusterIndex() + 1, 0, width
	}

	width := tw.endLine.WidthWithGhostSpaces()
	i := tw.endLine.BreakClusterIndex() + 1
	for i <= endClusterIdx {
		cluster := parent.Cluster(i)
		if cluster == nil || !cluster.IsWhitespaceBreak() {
			break
		}
		width += cluster.Width()
		i++
	}

	breakCluster := parent.Cluster(tw.endLine.BreakClusterIndex())
	if breakCluster != nil && breakCluster.IsWhitespaceBreak() && tw.endLine.BreakClusterIndex() < endClusterIdx {
		tw.endLine.ShiftBreak(parent)
	}

	return i, 0, width
}

// getClustersTrimmedWidth returns the trimmed width of clusters.
func (tw *TextWrapper) getClustersTrimmedWidth(parent TextWrapperOwner) float32 {
	width := float32(0)
	trailingSpaces := true
	for i := tw.clusters.EndClusterIndex(); i >= tw.clusters.StartClusterIndex(); i-- {
		cluster := parent.Cluster(i)
		if cluster == nil {
			continue
		}
		run := cluster.Run()
		if run != nil && run.IsPlaceholder() {
			continue
		}
		if trailingSpaces {
			if !cluster.IsWhitespaceBreak() {
				width += cluster.TrimmedWidth(cluster.EndPos())
				trailingSpaces = false
			}
			continue
		}
		width += cluster.Width()
	}
	return width
}
