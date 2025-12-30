// Conic curve operations for Skia paths.
// Ported from: skia-source/src/core/SkGeometry.cpp (SkConic struct and methods)
// https://github.com/google/skia/blob/main/src/core/SkGeometry.cpp
package geometry

import (
	"math"

	"github.com/zodimo/go-skia-support/skia/enums"
)

// Conic represents a conic curve (weighted quadratic bezier).
// A conic is defined by 3 points and a weight.
// When w=1, it's equivalent to a quadratic bezier.
// Ported from: SkConic in SkGeometry.h
type Conic struct {
	Pts [3]Point
	W   Scalar
}

// MaxConicsForArc is the maximum number of conics needed to represent any arc
const MaxConicsForArc = 5

// NewConic creates a new Conic from 3 points and a weight
func NewConic(p0, p1, p2 Point, w Scalar) Conic {
	c := Conic{
		Pts: [3]Point{p0, p1, p2},
	}
	c.SetW(w)
	return c
}

// SetW sets the weight, guarding against invalid values
func (c *Conic) SetW(w Scalar) {
	if math.IsInf(float64(w), 0) || math.IsNaN(float64(w)) || w <= 0 {
		w = 1
	}
	c.W = w
}

// EvalAt evaluates the conic at parameter t and returns the point
// Ported from: SkConic::evalAt
func (c *Conic) EvalAt(t Scalar) Point {
	// Rational quadratic bezier formula:
	// P(t) = (1-t)² P0 + 2(1-t)t w P1 + t² P2
	//        --------------------------------
	//        (1-t)² + 2(1-t)t w + t²

	t2 := t * t
	u := 1 - t
	u2 := u * u
	tw := 2 * t * u * c.W

	denom := u2 + tw + t2

	return Point{
		X: (u2*c.Pts[0].X + tw*c.Pts[1].X + t2*c.Pts[2].X) / denom,
		Y: (u2*c.Pts[0].Y + tw*c.Pts[1].Y + t2*c.Pts[2].Y) / denom,
	}
}

// EvalTangentAt evaluates the conic tangent at parameter t
// Ported from: SkConic::evalTangentAt
func (c *Conic) EvalTangentAt(t Scalar) Point {
	// Handle degenerate cases
	if (t == 0 && c.Pts[0] == c.Pts[1]) || (t == 1 && c.Pts[1] == c.Pts[2]) {
		return Point{X: c.Pts[2].X - c.Pts[0].X, Y: c.Pts[2].Y - c.Pts[0].Y}
	}

	p20 := Point{X: c.Pts[2].X - c.Pts[0].X, Y: c.Pts[2].Y - c.Pts[0].Y}
	p10 := Point{X: c.Pts[1].X - c.Pts[0].X, Y: c.Pts[1].Y - c.Pts[0].Y}

	// C = w * P10
	cx := c.W * p10.X
	cy := c.W * p10.Y

	// A = w * P20 - P20
	ax := c.W*p20.X - p20.X
	ay := c.W*p20.Y - p20.Y

	// B = P20 - 2 * C
	bx := p20.X - 2*cx
	by := p20.Y - 2*cy

	// Eval At² + Bt + C
	return Point{
		X: (ax*t+bx)*t + cx,
		Y: (ay*t+by)*t + cy,
	}
}

// Chop subdivides the conic at t=0.5 into two conics
// Ported from: SkConic::chop
func (c *Conic) Chop() (Conic, Conic) {
	scale := 1 / (1 + c.W)
	newW := Scalar(math.Sqrt(float64(0.5 + c.W*0.5)))

	wp1 := Point{X: c.W * c.Pts[1].X, Y: c.W * c.Pts[1].Y}

	midX := (c.Pts[0].X + 2*wp1.X + c.Pts[2].X) * scale * 0.5
	midY := (c.Pts[0].Y + 2*wp1.Y + c.Pts[2].Y) * scale * 0.5
	mid := Point{X: midX, Y: midY}

	dst0 := Conic{
		Pts: [3]Point{
			c.Pts[0],
			{X: (c.Pts[0].X + wp1.X) * scale, Y: (c.Pts[0].Y + wp1.Y) * scale},
			mid,
		},
		W: newW,
	}
	dst1 := Conic{
		Pts: [3]Point{
			mid,
			{X: (wp1.X + c.Pts[2].X) * scale, Y: (wp1.Y + c.Pts[2].Y) * scale},
			c.Pts[2],
		},
		W: newW,
	}

	return dst0, dst1
}

// quadrantPts are predefined unit circle points at 0°, 45°, 90°, 135°, etc.
var quadrantPts = [8]Point{
	{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}, {X: -1, Y: 1},
	{X: -1, Y: 0}, {X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1},
}

// BuildUnitArc builds conic arcs for the unit circle from start to stop vectors.
// Returns up to MaxConicsForArc conics.
// Ported from: SkConic::BuildUnitArc
func BuildUnitArc(uStart, uStop Point, dir enums.PathDirection, userMatrix *Matrix) []Conic {
	// Rotate to canonical form: uStart becomes (1, 0)
	// x = dot(uStart, uStop), y = cross(uStart, uStop)
	x := uStart.X*uStop.X + uStart.Y*uStop.Y
	y := uStart.X*uStop.Y - uStart.Y*uStop.X

	absY := Scalar(math.Abs(float64(y)))

	// Check for (effectively) coincident vectors
	if absY <= ScalarNearlyZero && x > 0 &&
		((y >= 0 && dir == enums.PathDirectionCW) || (y <= 0 && dir == enums.PathDirectionCCW)) {
		return nil
	}

	if dir == enums.PathDirectionCCW {
		y = -y
	}

	// Determine which quadrant the angle falls in
	// 0 == [0..90), 1 == [90..180), 2 == [180..270), 3 == [270..360)
	quadrant := computeQuadrant(x, y)

	// Create one conic per full quadrant
	conics := make([]Conic, 0, MaxConicsForArc)
	for i := 0; i < quadrant; i++ {
		conics = append(conics, Conic{
			Pts: [3]Point{
				quadrantPts[i*2],
				quadrantPts[i*2+1],
				quadrantPts[(i*2+2)%8],
			},
			W: ScalarRoot2Over2,
		})
	}

	// Create the final sub-quadrant conic (if needed)
	finalP := Point{X: x, Y: y}
	lastQ := quadrantPts[quadrant*2]
	dot := lastQ.X*finalP.X + lastQ.Y*finalP.Y

	if !math.IsNaN(float64(dot)) && dot < 1 {
		// Compute bisector (off-curve point)
		offCurve := Point{X: lastQ.X + x, Y: lastQ.Y + y}

		// weight = cos(θ/2), length = 1/cos(θ/2)
		cosThetaOver2 := Scalar(math.Sqrt(float64((1 + dot) / 2)))
		if cosThetaOver2 > ScalarNearlyZero {
			length := 1 / cosThetaOver2
			offLen := Scalar(math.Sqrt(float64(offCurve.X*offCurve.X + offCurve.Y*offCurve.Y)))
			if offLen > ScalarNearlyZero {
				offCurve.X = offCurve.X / offLen * length
				offCurve.Y = offCurve.Y / offLen * length

				// Only add if control point differs from last quadrant point
				if !pointsNearlyEqual(lastQ, offCurve) {
					conics = append(conics, Conic{
						Pts: [3]Point{lastQ, offCurve, finalP},
						W:   cosThetaOver2,
					})
				}
			}
		}
	}

	// Apply rotation matrix to align with actual start direction
	sinCos := Matrix{
		Values: [9]Scalar{
			uStart.X, -uStart.Y, 0,
			uStart.Y, uStart.X, 0,
			0, 0, 1,
		},
	}

	if dir == enums.PathDirectionCCW {
		// Scale Y by -1
		sinCos.Values[3] = -sinCos.Values[3]
		sinCos.Values[4] = -sinCos.Values[4]
	}

	if userMatrix != nil {
		sinCos = sinCos.Concat(*userMatrix)
	}

	// Transform all conic points
	for i := range conics {
		for j := 0; j < 3; j++ {
			conics[i].Pts[j] = sinCos.MapPoint(conics[i].Pts[j])
		}
	}

	return conics
}

// computeQuadrant determines which quadrant [0-3] the angle (x,y) falls in
func computeQuadrant(x, y Scalar) int {
	if y == 0 {
		return 2 // 180°
	}
	if x == 0 {
		if y > 0 {
			return 1
		} // 90°
		return 3 // 270°
	}
	quadrant := 0
	if y < 0 {
		quadrant += 2
	}
	if (x < 0) != (y < 0) {
		quadrant += 1
	}
	return quadrant
}

// pointsNearlyEqual checks if two points are nearly equal
func pointsNearlyEqual(a, b Point) bool {
	return NearlyEqual(a.X, b.X) && NearlyEqual(a.Y, b.Y)
}

// Matrix is a simple 3x3 transformation matrix for conic transformation
type Matrix struct {
	Values [9]Scalar // row-major: [0-2] first row, [3-5] second row, [6-8] third row
}

// MapPoint applies the matrix transformation to a point
func (m *Matrix) MapPoint(p Point) Point {
	return Point{
		X: m.Values[0]*p.X + m.Values[1]*p.Y + m.Values[2],
		Y: m.Values[3]*p.X + m.Values[4]*p.Y + m.Values[5],
	}
}

// Concat concatenates this matrix with another (this * other)
func (m Matrix) Concat(other Matrix) Matrix {
	var result Matrix
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sum := Scalar(0)
			for k := 0; k < 3; k++ {
				sum += m.Values[i*3+k] * other.Values[k*3+j]
			}
			result.Values[i*3+j] = sum
		}
	}
	return result
}
