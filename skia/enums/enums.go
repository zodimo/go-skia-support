package enums

// PaintCap represents the stroke cap style
type PaintCap uint8

const (
	PaintCapButt    PaintCap = 0 // no stroke extension
	PaintCapRound   PaintCap = 1 // adds circle
	PaintCapSquare  PaintCap = 2 // adds square
	PaintCapLast    PaintCap = PaintCapSquare
	PaintCapDefault PaintCap = PaintCapButt
)

// PaintCapCount is the number of valid cap values
const PaintCapCount = PaintCapLast + 1

// PaintJoin represents the stroke join style
type PaintJoin uint8

const (
	PaintJoinMiter   PaintJoin = 0 // extends to miter limit
	PaintJoinRound   PaintJoin = 1 // adds circle
	PaintJoinBevel   PaintJoin = 2 // connects outside edges
	PaintJoinLast    PaintJoin = PaintJoinBevel
	PaintJoinDefault PaintJoin = PaintJoinMiter
)

// PaintStyle represents the paint style (fill, stroke, or both)
type PaintStyle uint8

const (
	PaintStyleFill          PaintStyle = 0 // set to fill geometry
	PaintStyleStroke        PaintStyle = 1 // set to stroke geometry
	PaintStyleStrokeAndFill PaintStyle = 2 // sets to stroke and fill geometry
)

// Corner represents a corner of a rounded rectangle
type Corner int

const (
	CornerUpperLeft  Corner = 0 // index of top-left corner radii
	CornerUpperRight Corner = 1 // index of top-right corner radii
	CornerLowerRight Corner = 2 // index of bottom-right corner radii
	CornerLowerLeft  Corner = 3 // index of bottom-left corner radii
)

// RRectType represents the type of a rounded rectangle
type RRectType int

const (
	RRectTypeEmpty     RRectType = 0 // zero width or height
	RRectTypeRect      RRectType = 1 // non-zero width and height, and zeroed radii
	RRectTypeOval      RRectType = 2 // non-zero width and height filled with radii
	RRectTypeSimple    RRectType = 3 // non-zero width and height with equal radii
	RRectTypeNinePatch RRectType = 4 // non-zero width and height with axis-aligned radii
	RRectTypeComplex   RRectType = 5 // non-zero width and height with arbitrary radii
)

// MatrixType represents the type mask for matrix classification
type MatrixType uint8

const (
	// MatrixTypeIdentity represents an identity matrix
	MatrixTypeIdentity MatrixType = 0
	// MatrixTypeTranslate represents a translation matrix
	MatrixTypeTranslate MatrixType = 0x01
	// MatrixTypeScale represents a scale matrix
	MatrixTypeScale MatrixType = 0x02
	// MatrixTypeAffine represents an affine matrix (skew or rotate)
	MatrixTypeAffine MatrixType = 0x04
	// MatrixTypePerspective represents a perspective matrix
	MatrixTypePerspective MatrixType = 0x08
)

// PathFillType represents the fill rule for paths
type PathFillType uint8

const (
	PathFillTypeWinding        PathFillType = 0
	PathFillTypeEvenOdd        PathFillType = 1
	PathFillTypeInverseWinding PathFillType = 2
	PathFillTypeInverseEvenOdd PathFillType = 3
	PathFillTypeDefault        PathFillType = PathFillTypeWinding
)

// PathConvexity represents the convexity type of a path
type PathConvexity uint8

const (
	PathConvexityConvexCW         PathConvexity = 0
	PathConvexityConvexCCW        PathConvexity = 1
	PathConvexityConvexDegenerate PathConvexity = 2
	PathConvexityConcave          PathConvexity = 3
	PathConvexityUnknown          PathConvexity = 4
)

// PathVerb represents a path verb
type PathVerb uint8

const (
	PathVerbMove  PathVerb = 0
	PathVerbLine  PathVerb = 1
	PathVerbQuad  PathVerb = 2
	PathVerbConic PathVerb = 3
	PathVerbCubic PathVerb = 4
	PathVerbClose PathVerb = 5
)

// PathDirection represents the direction for adding closed contours
type PathDirection uint8

const (
	PathDirectionCW      PathDirection = 0
	PathDirectionCCW     PathDirection = 1
	PathDirectionDefault PathDirection = PathDirectionCW
)

// AddPathMode represents how paths are added together
type AddPathMode uint8

const (
	AddPathModeAppend AddPathMode = 0
	AddPathModeExtend AddPathMode = 1
)

// PathFirstDirection represents the first direction of a path
type PathFirstDirection uint8

const (
	PathFirstDirectionCW      PathFirstDirection = 0
	PathFirstDirectionCCW     PathFirstDirection = 1
	PathFirstDirectionUnknown PathFirstDirection = 2
)

// DirChange represents a direction change in the path
type DirChange uint8

const (
	DirChangeUnknown   DirChange = 0
	DirChangeLeft      DirChange = 1
	DirChangeRight     DirChange = 2
	DirChangeStraight  DirChange = 3
	DirChangeBackwards DirChange = 4
	DirChangeInvalid   DirChange = 5
)
