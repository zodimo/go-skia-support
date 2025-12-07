package interfaces

import (
	"github.com/zodimo/go-skia-support/skia/enums"
)

// SkCanvas represents a drawing surface that can be implemented by any graphics backend.
// This interface follows the "Bring Your Own Graphics Backend" (BYOG) philosophy,
// allowing users to implement their own rendering backends (OpenGL, Vulkan, Metal, CPU, etc.)
// while using the library's helper functions and drawing primitives.
//
// The interface matches the public API of C++ SkCanvas from include/core/SkCanvas.h,
// ensuring familiarity for developers who have used Skia C++.
//
// Usage Example:
//
//	type MyCanvas struct {
//		// Your backend-specific fields
//	}
//
//	func (c *MyCanvas) DrawRect(rect Rect, paint SkPaint) {
//		// Implement rectangle drawing using your backend
//	}
//
//	// Implement all other SkCanvas methods...
//
//	// Now you can use helper functions that accept SkCanvas:
//	func DrawRoundedButton(canvas SkCanvas, rect Rect, radius Scalar, paint SkPaint) {
//		rrect := models.NewRRect(rect, radius, radius, radius, radius)
//		canvas.DrawRRect(rrect, paint)
//	}
//
// Ported from: skia-source/include/core/SkCanvas.h
type SkCanvas interface {
	// Drawing Methods
	// Ported from: skia-source/include/core/SkCanvas.h lines 1287-1584

	// DrawPaint fills the current clip region with the specified paint.
	// The paint's color, blend mode, shader, color filter, and image filter affect drawing.
	// MaskFilter and PathEffect in paint are ignored (have no effect).
	// The paint is applied to the clip region, not transformed by the current matrix.
	// If no clip is set, fills the entire canvas.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:drawPaint() - line 1287
	DrawPaint(paint SkPaint)

	// DrawRect draws a rectangle using the current clip, matrix transformation, and paint properties.
	// PaintStyle determines if rectangle is stroked or filled:
	//   - PaintStyleFill: Fills the rectangle interior
	//   - PaintStyleStroke: Strokes the rectangle outline
	//   - PaintStyleStrokeAndFill: Both fills and strokes
	// If stroked, StrokeWidth describes line thickness, PaintJoin draws corners.
	// Rectangle coordinates are transformed by the current matrix before drawing.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:drawRect() - line 1406
	DrawRect(rect Rect, paint SkPaint)

	// DrawRRect draws a rounded rectangle with up to eight corner radii (four corners, each with x and y radii).
	// PaintStyle determines if rrect is stroked or filled.
	// Each corner can have independent x and y radii. If radii are zero, draws as a regular rectangle.
	// If radii exceed rectangle dimensions, they are scaled down to fit.
	// RRect is transformed by the current matrix before drawing.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:drawRRect() - line 1457
	DrawRRect(rrect RRect, paint SkPaint)

	// DrawDRRect draws a "donut" shape - the area between outer and inner rounded rectangles.
	// outer must contain inner or drawing behavior is undefined.
	// PaintStyle determines if SkRRect is stroked or filled.
	// Both Rects are transformed by the current matrix before drawing.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:drawDRRect() - line 1478
	DrawDRRect(outer RRect, inner RRect, paint SkPaint)

	// DrawOval draws an oval (ellipse) bounded by the specified rectangle.
	// PaintStyle determines if oval is stroked or filled. If stroked, StrokeWidth describes line thickness.
	// Rectangle bounds are transformed by the current matrix before drawing the oval.
	// Empty rectangles (width or height <= 0) result in no drawing.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:drawOval() - line 1443
	DrawOval(oval Rect, paint SkPaint)

	// DrawArc draws an arc that is part of an oval bounded by oval, sweeping from startAngle to startAngle + sweepAngle.
	// Angles are in degrees. Zero degrees places start point at the right middle edge of oval (3 o'clock position).
	// Positive sweepAngle is clockwise; negative is counterclockwise. sweepAngle may exceed 360 degrees.
	// If useCenter is true, draws a wedge including lines from oval center to arc end points.
	// If useCenter is false, draws only the arc between end points.
	// If oval is empty or sweepAngle is zero, nothing is drawn.
	// Oval bounds are transformed by the current matrix before drawing.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:drawArc() - line 1527
	DrawArc(oval Rect, startAngle Scalar, sweepAngle Scalar, useCenter bool, paint SkPaint)

	// DrawPath draws a path containing one or more contours, each of which may be open or closed.
	// If filled: PathFillType determines whether path contour describes inside or outside of fill.
	// If stroked: StrokeWidth describes line thickness, PaintCap describes line ends, PaintJoin describes corners.
	// Path may contain multiple contours (each starting with MoveTo). Each contour may be open or closed.
	// Path geometry is transformed by the current matrix before drawing.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:drawPath() - line 1584
	DrawPath(path SkPath, paint SkPaint)

	// DrawPoints draws an array of points according to the specified mode.
	// PointMode behavior:
	//   - PointModePoints: Draws each point separately. Shape depends on PaintCap:
	//     - PaintCapRound: Circle of diameter StrokeWidth
	//     - PaintCapSquare or PaintCapButt: Square of width/height StrokeWidth
	//   - PointModeLines: Each pair of points draws a line segment. One line per two points.
	//     If count is odd, final point is ignored.
	//   - PointModePolygon: Each adjacent pair draws a line segment. count-1 lines drawn, connecting all points sequentially.
	// PaintStyle is ignored - always treated as stroke. PaintJoin is ignored - elements drawn one at a time.
	// If len(points) < 1, has no effect.
	// Points are transformed by the current matrix before drawing.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:drawPoints() - line 1330
	DrawPoints(mode enums.PointMode, points []Point, paint SkPaint)

	// Clipping Methods
	// Ported from: skia-source/include/core/SkCanvas.h lines 1019-1151

	// ClipRect replaces clip with the intersection or difference of current clip and rect.
	// ClipOp behavior:
	//   - ClipOpIntersect: Clip becomes intersection of current clip and rect (default, most common)
	//   - ClipOpDifference: Clip becomes current clip minus rect
	//   - ClipOpUnion: Clip becomes union of current clip and rect
	//   - ClipOpXOR: Clip becomes XOR of current clip and rect
	//   - ClipOpReverseDifference: Clip becomes rect minus current clip
	// doAntiAlias controls edge quality:
	//   - false: Aliased clip - pixels are fully contained by the clip (faster)
	//   - true: Anti-aliased clip - smooth edges with partial pixel coverage (slower, better quality)
	// rect is transformed by the current matrix before being combined with clip.
	// Clips are cumulative and cannot be expanded (except via Save/Restore).
	//
	// Ported from: skia-source/include/core/SkCanvas.h:clipRect() - line 1019
	ClipRect(rect Rect, op enums.ClipOp, doAntiAlias bool)

	// ClipRRect replaces clip with the intersection or difference of current clip and rounded rectangle.
	// ClipOp and doAntiAlias behavior same as ClipRect.
	// rrect is transformed by the current matrix before being combined with clip.
	// Each corner can have independent radii. Zero radii behave like regular rectangle clipping.
	// Clips are cumulative and restrictive.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:clipRRect() - line 1073
	ClipRRect(rrect RRect, op enums.ClipOp, doAntiAlias bool)

	// ClipPath replaces clip with the intersection or difference of current clip and path.
	// PathFillType determines if path describes area inside or outside contours, and how overlaps are handled:
	//   - PathFillTypeWinding: Non-zero winding rule
	//   - PathFillTypeEvenOdd: Even-odd rule
	//   - Inverse variants describe outside instead of inside
	// ClipOp and doAntiAlias behavior same as ClipRect.
	// path is transformed by the current matrix before being combined with clip.
	// More expensive than rect/RRect clipping due to path complexity.
	// Clips are cumulative and restrictive.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:clipPath() - line 1109
	ClipPath(path SkPath, op enums.ClipOp, doAntiAlias bool)

	// Transformation State Methods
	// Ported from: skia-source/include/core/SkCanvas.h lines 876-988

	// Save saves the current matrix and clip state to a stack. Returns the save count (depth of stack before this save).
	// Initial canvas has save count of 1. Each Save() increments the count.
	// Both matrix transformation and clip region are saved.
	// Return value can be used with RestoreToCount.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:save() - line ~850
	Save() int

	// Restore removes the most recent save state from the stack, restoring matrix and clip to previous values.
	// Does nothing if the stack is empty (save count is 1).
	// Both matrix and clip are restored to their state at the last Save() call.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:restore() - line 876
	Restore()

	// RestoreToCount restores state to the matrix and clip values when Save() returned saveCount.
	// Does nothing if saveCount is greater than current state stack count.
	// Restores to initial values if saveCount is less than or equal to 1.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:restoreToCount() - line 898
	RestoreToCount(saveCount int)

	// GetSaveCount returns the number of saved states (depth of save stack).
	// New canvas has save count of 1.
	// Equals number of Save() calls less number of Restore() calls plus one.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:getSaveCount() - line 886
	GetSaveCount() int

	// Concat replaces current matrix with matrix premultiplied with existing matrix.
	// Mathematical effect: newMatrix = matrix * currentMatrix (matrix applied first, then current transformation).
	// The matrix transformation is applied to geometry first, then the existing matrix transformation.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:concat() - line 987
	Concat(matrix SkMatrix)

	// Translate translates the current matrix by dx along x-axis and dy along y-axis.
	// Mathematical effect: Premultiplies current matrix with a translation matrix.
	// Moves drawing by (dx, dy) before applying existing matrix transformation.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:translate() - line 913
	Translate(dx Scalar, dy Scalar)

	// Scale scales the current matrix by sx on x-axis and sy on y-axis.
	// Mathematical effect: Premultiplies current matrix with a scale matrix.
	// Scales drawing by (sx, sy) before applying existing matrix transformation.
	// Negative values mirror/flip. Zero values collapse dimension.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:scale() - line 928
	Scale(sx Scalar, sy Scalar)

	// Rotate rotates the current matrix by degrees around the origin (0, 0).
	// Mathematical effect: Premultiplies current matrix with a rotation matrix.
	// Positive degrees rotates clockwise (mathematical convention).
	// Rotates drawing by degrees around origin before applying existing matrix transformation.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:rotate() - line 942
	Rotate(degrees Scalar)

	// Skew skews the current matrix by sx on x-axis and sy on y-axis.
	// Mathematical effect: Premultiplies current matrix with a skew matrix.
	// Direction:
	//   - Positive sx: Skews right as y-axis values increase
	//   - Positive sy: Skews down as x-axis values increase
	// Skews drawing by (sx, sy) before applying existing matrix transformation.
	//
	// Ported from: skia-source/include/core/SkCanvas.h:skew() - line 976
	Skew(sx Scalar, sy Scalar)
}
