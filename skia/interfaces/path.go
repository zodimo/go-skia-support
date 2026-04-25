package interfaces

import (
	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

// SkPath represents a path that can be drawn on a canvas.
// It provides methods for creating and manipulating path geometry.
type SkPath interface {
	// FillType returns the fill type used to determine which parts are inside.
	FillType() enums.PathFillType

	// SetFillType sets the fill type used to determine which parts are inside.
	SetFillType(fillType enums.PathFillType)

	// IsInverseFillType returns true if the fill type is inverse.
	IsInverseFillType() bool

	// ToggleInverseFillType toggles between inverse and non-inverse fill types.
	ToggleInverseFillType()

	// Convexity returns the convexity type of the path.
	Convexity() enums.PathConvexity

	// IsConvex returns true if the path is convex.
	IsConvex() bool

	// Reset clears the path, removing all verbs, points, and conic weights.
	Reset()

	// IsEmpty returns true if the path has no verbs.
	IsEmpty() bool

	// IsFinite returns true if all points in the path are finite.
	IsFinite() bool

	// IsLine returns true if the path contains only one line.
	IsLine() bool

	// CountPoints returns the number of points in the path.
	CountPoints() int

	// Point returns the point at the specified index.
	Point(index int) models.Point

	// GetPoints copies all points from the path into the provided slice.
	GetPoints(points []models.Point) int

	// CountVerbs returns the number of verbs in the path.
	CountVerbs() int

	// GetVerbs copies all verbs from the path into the provided slice.
	GetVerbs(verbs []enums.PathVerb) int

	// ConicWeights returns a read-only view of the path's conic weights.
	// Returns a copy of the conic weights slice.
	ConicWeights() []base.Scalar

	// GetLastPoint returns the last point in the path.
	// Returns the point and true if the path contains one or more points,
	// otherwise returns a zero point and false.
	GetLastPoint() (models.Point, bool)

	// Bounds returns the bounding box of the path.
	Bounds() models.Rect

	// UpdateBoundsCache updates the cached bounds of the path.
	UpdateBoundsCache()

	// ComputeTightBounds returns a tight bounding box of the path.
	ComputeTightBounds() models.Rect

	// MoveTo starts a new contour at the specified point.
	MoveTo(x, y base.Scalar)

	// MoveToPoint starts a new contour at the specified point.
	MoveToPoint(p models.Point)

	// LineTo adds a line from the last point to the specified point.
	LineTo(x, y base.Scalar)

	// LineToPoint adds a line from the last point to the specified point.
	LineToPoint(p models.Point)

	// QuadTo adds a quadratic bezier from the last point to the specified point.
	QuadTo(cx, cy, x, y base.Scalar)

	// QuadToPoint adds a quadratic bezier from the last point to the specified point.
	QuadToPoint(c, p models.Point)

	// ConicTo adds a conic bezier from the last point to the specified point.
	ConicTo(cx, cy, x, y base.Scalar, w base.Scalar)

	// ConicToPoint adds a conic bezier from the last point to the specified point.
	ConicToPoint(c, p models.Point, w base.Scalar)

	// CubicTo adds a cubic bezier from the last point to the specified point.
	CubicTo(cx1, cy1, cx2, cy2, x, y base.Scalar)

	// CubicToPoint adds a cubic bezier from the last point to the specified point.
	CubicToPoint(c1, c2, p models.Point)

	// Close closes the current contour.
	Close()

	// AddRect adds a rectangle to the path.
	AddRect(rect models.Rect, dir enums.PathDirection, startIndex uint)

	// AddOval adds an oval to the path.
	AddOval(rect models.Rect, dir enums.PathDirection)

	// AddCircle adds a circle to the path.
	AddCircle(cx, cy, radius base.Scalar, dir enums.PathDirection)

	// AddRRect adds a rounded rectangle to the path.
	AddRRect(rrect models.RRect, dir enums.PathDirection)

	// AddPath adds another path to this path with offset.
	AddPath(path SkPath, dx, dy base.Scalar, addMode enums.AddPathMode)

	// AddPathNoOffset adds another path to this path without offset.
	AddPathNoOffset(path SkPath, addMode enums.AddPathMode)

	// AddPathMatrix adds another path to this path with matrix transformation.
	AddPathMatrix(path SkPath, matrix SkMatrix, addMode enums.AddPathMode)

	// Transform applies a matrix transformation to the path.
	Transform(matrix SkMatrix)

	// Offset translates the path by the specified offset.
	Offset(dx, dy base.Scalar)

	// ArcTo appends arc from oval from startAngle through sweepAngle.
	// Angles are in degrees. Positive sweep is clockwise.
	// If forceMoveTo is true, starts a new contour; otherwise connects to last point.
	// Ported from: SkPath.h arcTo(oval, startAngle, sweepAngle, forceMoveTo)
	ArcTo(oval models.Rect, startAngle, sweepAngle base.Scalar, forceMoveTo bool)

	// ArcToTangent appends arc tangent to line from last point through (x1,y1)
	// to line from (x1,y1) to (x2,y2), with specified radius.
	// Implements HTML Canvas arcTo and PostScript arct.
	// Ported from: SkPath.h arcTo(x1, y1, x2, y2, radius)
	ArcToTangent(x1, y1, x2, y2, radius base.Scalar)

	// ArcToRotated appends SVG-style elliptical arc to (x,y).
	// rx, ry are the ellipse radii; xAxisRotate is the rotation in degrees.
	// Ported from: SkPath.h arcTo(rx, ry, xAxisRotate, largeArc, sweep, x, y)
	ArcToRotated(rx, ry, xAxisRotate base.Scalar, largeArc enums.ArcSize, sweep enums.PathDirection, x, y base.Scalar)

	// RArcTo appends SVG-style elliptical arc relative to current point.
	// Ported from: SkPath.h rArcTo(rx, ry, xAxisRotate, largeArc, sweep, dx, dy)
	RArcTo(rx, ry, xAxisRotate base.Scalar, largeArc enums.ArcSize, sweep enums.PathDirection, dx, dy base.Scalar)

	// AddArc adds arc as a new contour (starts with implicit MoveTo).
	// Ported from: SkPath.h addArc(oval, startAngle, sweepAngle)
	AddArc(oval models.Rect, startAngle, sweepAngle base.Scalar)
}
