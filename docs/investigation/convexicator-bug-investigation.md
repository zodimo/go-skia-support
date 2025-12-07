# Convexicator Concave Detection Bug - Investigation Document

**Date:** 2025-01-27  
**Issue:** Concave quadrilateral test fails - path detected as convex instead of concave  
**Status:** Under Investigation

---

## Table of Contents

1. [Issue Summary](#issue-summary)
2. [C++ Reference Implementation](#c-reference-implementation)
3. [Go Implementation](#go-implementation)
4. [Test Cases](#test-cases)
5. [Helper Functions](#helper-functions)
6. [Path Processing Logic](#path-processing-logic)
7. [Manual Trace Analysis](#manual-trace-analysis)
8. [Potential Root Causes](#potential-root-causes)
9. [Debugging Recommendations](#debugging-recommendations)

---

## Issue Summary

### Problem Statement

The concave quadrilateral test case fails in the Go implementation. The path should be detected as **concave** but is incorrectly detected as **convex**.

### Failing Test Case

**Test:** Concave Quadrilateral  
**Points:** `(0,0) -> (10,10) -> (10,0) -> (0,10) -> close`  
**Expected:** `IsConvex() = false` (concave)  
**Actual:** `IsConvex() = true` (convex)  
**Convexity Type:** `0` (PathConvexityConvexDegenerate or PathConvexityUnknown)

### What Was Fixed

1. ✅ **`addPt()` method:** Changed condition from checking `vec.X == 0 && vec.Y == 0` to `c.lastVec.X == 0 && c.lastVec.Y == 0` to match C++ `fLastVec.equals(0,0)` check
2. ✅ **`setMovePt()` method:** Added explicit resets for `lastVec`, `firstVec`, `reversals`, and `isFinite` to ensure clean state initialization

### What Was Verified

- ✅ `directionChange()` method: Cross product calculation matches C++
- ✅ `addVec()` method: Reversal logic matches C++ (`++fReversals < 3` equivalent)
- ✅ `close()` method: Logic matches C++ (`addPt(fFirstPt) && addVec(fFirstVec)`)
- ✅ Helper functions: `crossProduct()`, `dotProduct()`, and `sign()` implementations verified correct

### Current Status

The implementation structurally matches the C++ code, but concave shapes are not being detected correctly. The bug appears to be a subtle logic error that requires deeper investigation.

---

## C++ Reference Implementation

### Source Location
`skia-source/src/core/SkPathPriv.cpp` lines 419-560

### Complete Convexicator Struct

```cpp
struct Convexicator {

    /** The direction returned is only valid if the path is determined convex */
    SkPathFirstDirection getFirstDirection() const { return fFirstDirection; }

    void setMovePt(const SkPoint& pt) {
        fFirstPt = fLastPt = pt;
        fExpectedDir = kInvalid_DirChange;
    }

    bool addPt(const SkPoint& pt) {
        if (fLastPt == pt) {
            return true;
        }
        // should only be true for first non-zero vector after setMovePt was called. It is possible
        // we doubled backed at the start so need to check if fLastVec is zero or not.
        if (fFirstPt == fLastPt && fExpectedDir == kInvalid_DirChange && fLastVec.equals(0,0)) {
            fLastVec = pt - fLastPt;
            fFirstVec = fLastVec;
        } else if (!this->addVec(pt - fLastPt)) {
            return false;
        }
        fLastPt = pt;
        return true;
    }

    bool close() {
        // If this was an explicit close, there was already a lineTo to fFirstPoint, so this
        // addPt() is a no-op. Otherwise, the addPt implicitly closes the contour. In either case,
        // we have to check the direction change along the first vector in case it is concave.
        return this->addPt(fFirstPt) && this->addVec(fFirstVec);
    }

    bool isFinite() const {
        return fIsFinite;
    }

    int reversals() const {
        return fReversals;
    }

private:
    DirChange directionChange(const SkVector& curVec) {
        SkScalar cross = SkPoint::CrossProduct(fLastVec, curVec);
        if (!SkIsFinite(cross)) {
            return kUnknown_DirChange;
        }
        if (cross == 0) {
            return fLastVec.dot(curVec) < 0 ? kBackwards_DirChange : kStraight_DirChange;
        }
        return 1 == SkScalarSignAsInt(cross) ? kRight_DirChange : kLeft_DirChange;
    }

    bool addVec(const SkVector& curVec) {
        DirChange dir = this->directionChange(curVec);
        switch (dir) {
            case kLeft_DirChange:       // fall through
            case kRight_DirChange:
                if (kInvalid_DirChange == fExpectedDir) {
                    fExpectedDir = dir;
                    fFirstDirection = (kRight_DirChange == dir) ? SkPathFirstDirection::kCW
                                                                : SkPathFirstDirection::kCCW;
                } else if (dir != fExpectedDir) {
                    fFirstDirection = SkPathFirstDirection::kUnknown;
                    return false;
                }
                fLastVec = curVec;
                break;
            case kStraight_DirChange:
                break;
            case kBackwards_DirChange:
                //  allow path to reverse direction twice
                //    Given path.moveTo(0, 0); path.lineTo(1, 1);
                //    - 1st reversal: direction change formed by line (0,0 1,1), line (1,1 0,0)
                //    - 2nd reversal: direction change formed by line (1,1 0,0), line (0,0 1,1)
                fLastVec = curVec;
                return ++fReversals < 3;
            case kUnknown_DirChange:
                return (fIsFinite = false);
            case kInvalid_DirChange:
                SK_ABORT("Use of invalid direction change flag");
                break;
        }
        return true;
    }

    SkPoint              fFirstPt {0, 0};  // The first point of the contour, e.g. moveTo(x,y)
    SkVector             fFirstVec {0, 0}; // The direction leaving fFirstPt to the next vertex

    SkPoint              fLastPt {0, 0};   // The last point passed to addPt()
    SkVector             fLastVec {0, 0};  // The direction that brought the path to fLastPt

    DirChange            fExpectedDir { kInvalid_DirChange };
    SkPathFirstDirection fFirstDirection { SkPathFirstDirection::kUnknown };
    int                  fReversals { 0 };
    bool                 fIsFinite { true };
};
```

### Key C++ Helper Functions

**SkPoint::CrossProduct** (from `include/core/SkPoint.h`):
```cpp
static SkScalar CrossProduct(const SkPoint& a, const SkPoint& b) {
    return a.fX * b.fY - a.fY * b.fX;
}
```

**SkScalarSignAsInt** (from `include/core/SkScalar.h`):
```cpp
static inline int SkScalarSignAsInt(SkScalar x) {
    return x < 0 ? -1 : (x > 0);
}
// Returns: -1 if x < 0, 0 if x == 0, 1 if x > 0
```

### C++ Test Case

**Source:** `skia-source/tests/PathTest.cpp` lines 1692-1710

```cpp
static void test_convexity(skiatest::Reporter* reporter) {
    // ... other test cases ...
    
    static const struct {
        const char*           fPathStr;
        bool                  fExpectedIsConvex;
        SkPathFirstDirection  fExpectedDirection;
    } gRec[] = {
        { "", true, SkPathFirstDirection::kUnknown },
        { "0 0", true, SkPathFirstDirection::kUnknown },
        { "0 0 10 10", true, SkPathFirstDirection::kUnknown },
        { "0 0 10 10 20 20 0 0 10 10", false, SkPathFirstDirection::kUnknown },
        { "0 0 10 10 10 20", true, SkPathFirstDirection::kCW },
        { "0 0 10 10 10 0", true, SkPathFirstDirection::kCCW },
        { "0 0 10 10 10 0 0 10", false, kDontCheckDir },  // <-- CONCAVE QUADRILATERAL
        { "0 0 10 0 0 10 -10 -10", false, SkPathFirstDirection::kCW },
    };

    for (size_t i = 0; i < std::size(gRec); ++i) {
        path = setFromString(gRec[i].fPathStr);
        check_convexity(reporter, path, gRec[i].fExpectedIsConvex);
        check_direction(reporter, path, gRec[i].fExpectedDirection);
    }
}
```

**Test Case:** `"0 0 10 10 10 0 0 10"`  
**Expected:** `fExpectedIsConvex = false` (concave)

---

## Go Implementation

### Source Location
`skia/impl/path_models.go` lines 164-275

### Complete Convexicator Struct

```go
// Convexicator tracks convexity state while iterating through a path
type convexicator struct {
	firstPt        Point
	firstVec       Point // direction leaving firstPt
	lastPt         Point
	lastVec        Point // direction that brought path to lastPt
	expectedDir    enums.DirChange
	firstDirection enums.PathFirstDirection
	reversals      int
	isFinite       bool
}

func newConvexicator() *convexicator {
	return &convexicator{
		expectedDir:    enums.DirChangeInvalid,
		firstDirection: enums.PathFirstDirectionUnknown,
		isFinite:       true,
	}
}

func (c *convexicator) setMovePt(pt Point) {
	c.firstPt = pt
	c.lastPt = pt
	c.expectedDir = enums.DirChangeInvalid
	// Reset vectors to zero to match C++ initialization
	// Ported from: skia-source/src/core/SkPathPriv.cpp:setMovePt() (lines 424-427)
	c.lastVec = Point{X: 0, Y: 0}
	c.firstVec = Point{X: 0, Y: 0}
	c.reversals = 0
	c.isFinite = true
}

func (c *convexicator) addPt(pt Point) bool {
	if c.lastPt == pt {
		return true
	}
	// Should only be true for first non-zero vector after setMovePt was called.
	// It is possible we doubled back at the start so need to check if lastVec is zero or not.
	// Ported from: skia-source/src/core/SkPathPriv.cpp:addPt() (lines 429-443)
	vec := Point{X: pt.X - c.lastPt.X, Y: pt.Y - c.lastPt.Y}
	if c.firstPt == c.lastPt && c.expectedDir == enums.DirChangeInvalid && c.lastVec.X == 0 && c.lastVec.Y == 0 {
		c.lastVec = vec
		c.firstVec = vec
	} else if !c.addVec(vec) {
		return false
	}
	c.lastPt = pt
	return true
}

func (c *convexicator) close() bool {
	// If this was an explicit close, there was already a lineTo to firstPt, so this
	// addPt() is a no-op. Otherwise, the addPt implicitly closes the contour.
	return c.addPt(c.firstPt) && c.addVec(c.firstVec)
}

func (c *convexicator) getFirstDirection() enums.PathFirstDirection {
	return c.firstDirection
}

func (c *convexicator) directionChange(curVec Point) enums.DirChange {
	cross := crossProduct(c.lastVec, curVec)
	if !IsFinite(cross) {
		return enums.DirChangeUnknown
	}
	if cross == 0 {
		dot := dotProduct(c.lastVec, curVec)
		if dot < 0 {
			return enums.DirChangeBackwards
		}
		return enums.DirChangeStraight
	}
	if cross > 0 {
		return enums.DirChangeRight
	}
	return enums.DirChangeLeft
}

func (c *convexicator) addVec(curVec Point) bool {
	dir := c.directionChange(curVec)
	switch dir {
	case enums.DirChangeLeft, enums.DirChangeRight:
		if c.expectedDir == enums.DirChangeInvalid {
			c.expectedDir = dir
			if dir == enums.DirChangeRight {
				c.firstDirection = enums.PathFirstDirectionCW
			} else {
				c.firstDirection = enums.PathFirstDirectionCCW
			}
		} else if dir != c.expectedDir {
			c.firstDirection = enums.PathFirstDirectionUnknown
			return false
		}
		c.lastVec = curVec
	case enums.DirChangeStraight:
		// Continue with same direction
	case enums.DirChangeBackwards:
		// Allow path to reverse direction twice
		c.lastVec = curVec
		c.reversals++
		if c.reversals >= 3 {
			return false
		}
	case enums.DirChangeUnknown:
		c.isFinite = false
		return false
	case enums.DirChangeInvalid:
		// Should not happen
		return false
	}
	return true
}
```

### Go Test Case

**Source:** `skia/impl/path_test.go` lines 262-295

```go
func TestPath_Convexity(t *testing.T) {
	// ... other test cases ...
	
	testCases := []struct {
		name           string
		points         []models.Point
		expectedConvex bool
		description    string
	}{
		// ... other cases ...
		{
			name:           "concave quadrilateral",
			points:         []models.Point{{X: 0, Y: 0}, {X: 10, Y: 10}, {X: 10, Y: 0}, {X: 0, Y: 10}},
			expectedConvex: false,
		},
		// ... other cases ...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := NewSkPath(enums.PathFillTypeDefault)
			if len(tc.points) > 0 {
				path.MoveTo(tc.points[0].X, tc.points[0].Y)
				for i := 1; i < len(tc.points); i++ {
					path.LineTo(tc.points[i].X, tc.points[i].Y)
				}
				path.Close()
			}
			checkConvexity(t, path, tc.expectedConvex)
			// Also verify direct IsConvex() call matches
			if path.IsConvex() != tc.expectedConvex {
				t.Errorf("Direct IsConvex() call: expected %v, got %v", tc.expectedConvex, path.IsConvex())
			}
		})
	}
}
```

**Test Case:** `(0,0) -> (10,10) -> (10,0) -> (0,10) -> close`  
**Expected:** `expectedConvex = false` (concave)  
**Actual:** `IsConvex() = true` (convex) ❌

---

## Helper Functions

### Go Helper Functions

**Source:** `skia/helpers/helpers.go`

```go
package helpers

func CrossProduct(a, b Point) Scalar {
	return a.X*b.Y - a.Y*b.X
}

func DotProduct(a, b Point) Scalar {
	return a.X*b.X + a.Y*b.Y
}

func Sign(x Scalar) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}
```

**Source:** `skia/impl/alias.go`

```go
func crossProduct(a, b Point) Scalar {
	return helpers.CrossProduct(a, b)
}

func dotProduct(a, b Point) Scalar {
	return helpers.DotProduct(a, b)
}

func sign(x Scalar) int {
	return helpers.Sign(x)
}
```

**Verification:** ✅ Helper functions match C++ implementation exactly.

---

## Path Processing Logic

### Go Path Processing

**Source:** `skia/impl/path.go` lines 680-760

```go
func (p *pathImpl) computeConvexity() enums.PathConvexity {
	// ... early returns for empty/simple paths ...
	
	// Quick concave test: check if path changes direction more than three times
	if isConcaveBySign(points) {
		return enums.PathConvexityConcave
	}

	// Iterate through the path and check convexity
	contourCount := 0
	needsClose := false
	state := newConvexicator()

	pointIdx := 0
	conicWeightIdx := 0

	for _, verb := range verbs {
		// Looking for the last moveTo before non-move verbs start
		if contourCount == 0 {
			if verb == enums.PathVerbMove {
				if pointIdx < len(points) {
					state.setMovePt(points[pointIdx])
					pointIdx++
				}
			} else {
				// Starting the actual contour, fall through to add the points
				// Note: This assumes there was a MoveTo (which should always be the case)
				contourCount++
				needsClose = true
			}
		}

		// Accumulating points into the Convexicator until we hit a close or another move
		if contourCount == 1 {
			if verb == enums.PathVerbClose || verb == enums.PathVerbMove {
				if !state.close() {
					return enums.PathConvexityConcave
				}
				needsClose = false
				contourCount++
				if verb == enums.PathVerbMove {
					if pointIdx < len(points) {
						state.setMovePt(points[pointIdx])
						pointIdx++
					}
				}
			} else {
				// Lines add 1 point, cubics add 3, conics and quads add 2
				// These are the points AFTER the start point (which is tracked in state.lastPt)
				count := ptsInVerb(verb)
				if count > 0 && pointIdx+count-1 < len(points) {
					for i := 0; i < count; i++ {
						if !state.addPt(points[pointIdx+i]) {
							return enums.PathConvexityConcave
						}
					}
					pointIdx += count
					if verb == enums.PathVerbConic {
						conicWeightIdx++
					}
				}
			}
		} else {
			// The first contour has closed and anything other than spurious trailing moves means
			// there's multiple contours and the path can't be convex
			if verb != enums.PathVerbMove {
				return enums.PathConvexityConcave
			}
			if pointIdx < len(points) {
				pointIdx++
			}
		}
	}

	// If the path isn't explicitly closed, do so implicitly
	if needsClose && !state.close() {
		return enums.PathConvexityConcave
	}

	firstDir := state.getFirstDirection()
	// ... determine convexity type based on firstDir ...
}
```

---

## Manual Trace Analysis

### Concave Quadrilateral Case

**Path:** `(0,0) -> (10,10) -> (10,0) -> (0,10) -> close`

#### Expected Execution Flow

1. **`setMovePt(0,0)`**
   - `firstPt = (0,0)`
   - `lastPt = (0,0)`
   - `expectedDir = Invalid`
   - `lastVec = (0,0)`
   - `firstVec = (0,0)`

2. **`addPt(10,10)`**
   - `vec = (10,10) - (0,0) = (10,10)`
   - Condition: `firstPt == lastPt` ✅, `expectedDir == Invalid` ✅, `lastVec == (0,0)` ✅
   - **Initialize:** `lastVec = (10,10)`, `firstVec = (10,10)`
   - `lastPt = (10,10)`

3. **`addPt(10,0)`**
   - `vec = (10,0) - (10,10) = (0,-10)`
   - Condition fails (not first point), so call `addVec((0,-10))`
   - `cross = crossProduct((10,10), (0,-10)) = 10*(-10) - 10*0 = -100`
   - `cross < 0` → `DirChangeLeft`
   - `expectedDir == Invalid` → Set `expectedDir = Left`, `firstDirection = CCW`
   - `lastVec = (0,-10)`
   - `lastPt = (10,0)`

4. **`addPt(0,10)`**
   - `vec = (0,10) - (10,0) = (-10,10)`
   - Call `addVec((-10,10))`
   - `cross = crossProduct((0,-10), (-10,10)) = 0*10 - (-10)*(-10) = -100`
   - `cross < 0` → `DirChangeLeft`
   - `expectedDir == Left` ✅, `dir == Left` ✅ → **Match!**
   - `lastVec = (-10,10)`
   - `lastPt = (0,10)`

5. **`close()`**
   - Calls `addPt((0,0))`
     - `vec = (0,0) - (0,10) = (0,-10)`
     - Call `addVec((0,-10))`
     - `cross = crossProduct((-10,10), (0,-10)) = (-10)*(-10) - 10*0 = 100`
     - `cross > 0` → `DirChangeRight`
     - `expectedDir == Left`, `dir == Right` → **MISMATCH!**
     - **Should return `false`** → Path is concave ✅
   - But wait, `addPt((0,0))` updates `lastVec = (0,-10)` before returning
   - Then `close()` calls `addVec(firstVec)` where `firstVec = (10,10)`
   - `cross = crossProduct((0,-10), (10,10)) = 0*10 - (-10)*10 = 100`
   - `cross > 0` → `DirChangeRight`
   - `expectedDir == Left`, `dir == Right` → **MISMATCH!**
   - **Should return `false`** → Path is concave ✅

#### Critical Observation

The logic **should** detect the path as concave because:
- After `addPt((0,0))` in `close()`, `lastVec = (0,-10)` (the closing edge vector)
- When `addVec(firstVec)` is called with `firstVec = (10,10)`, it compares:
  - `crossProduct((0,-10), (10,10)) = 100 > 0` → `DirChangeRight`
  - But `expectedDir = Left` → **Mismatch!** → Returns `false` → Concave ✅

**However, the test fails**, suggesting:
- Either `close()` is not being called correctly
- Or `addVec(firstVec)` is not detecting the direction change
- Or `computeConvexity()` is not properly handling the `false` return from `close()`
- Or there's an issue with how `addPt((0,0))` modifies `lastVec` before `addVec(firstVec)` is called

---

## Potential Root Causes

### Issue #1: Path Processing Logic

**Hypothesis:** `computeConvexity()` may not be calling `close()` correctly for explicitly closed paths.

**Investigation Points:**
- Verify that `PathVerbClose` triggers `state.close()` correctly
- Check if point indexing matches C++ `SkPathIter` behavior
- Verify that `needsClose` flag is set correctly

### Issue #2: Closing Edge Detection

**Hypothesis:** `close()` may not be correctly handling the sequence of `addPt(fFirstPt)` followed by `addVec(fFirstVec)`.

**Investigation Points:**
- Verify that `addPt(fFirstPt)` correctly updates `lastVec` before `addVec(fFirstVec)` is called
- Check if there's an edge case where `addPt(fFirstPt)` is a no-op (when path is explicitly closed)
- Verify that `lastVec` state is correct when `addVec(fFirstVec)` is called

### Issue #3: Direction Change Logic

**Hypothesis:** There may be an edge case in `directionChange()` when comparing the closing edge to the first vector.

**Investigation Points:**
- Verify cross product calculation for the closing edge case
- Check if there's a floating-point precision issue
- Verify that `expectedDir` is correctly maintained throughout

### Issue #4: State Management

**Hypothesis:** The state of `lastVec` may be incorrect when `close()` calls `addVec(fFirstVec)`.

**Investigation Points:**
- Trace the exact state of `lastVec` at each step
- Verify that `addPt(fFirstPt)` correctly updates `lastVec`
- Check if `lastVec` is being modified incorrectly somewhere

---

## Debugging Recommendations

### 1. Add Comprehensive Logging

Add detailed logging to trace execution:

```go
func (c *convexicator) addPt(pt Point) bool {
    log.Printf("addPt: pt=(%v,%v), lastPt=(%v,%v), lastVec=(%v,%v), expectedDir=%v",
        pt.X, pt.Y, c.lastPt.X, c.lastPt.Y, c.lastVec.X, c.lastVec.Y, c.expectedDir)
    // ... rest of implementation
}

func (c *convexicator) addVec(curVec Point) bool {
    cross := crossProduct(c.lastVec, curVec)
    dir := c.directionChange(curVec)
    log.Printf("addVec: curVec=(%v,%v), lastVec=(%v,%v), cross=%v, dir=%v, expectedDir=%v",
        curVec.X, curVec.Y, c.lastVec.X, c.lastVec.Y, cross, dir, c.expectedDir)
    // ... rest of implementation
}

func (c *convexicator) close() bool {
    log.Printf("close: firstPt=(%v,%v), firstVec=(%v,%v), lastPt=(%v,%v), lastVec=(%v,%v)",
        c.firstPt.X, c.firstPt.Y, c.firstVec.X, c.firstVec.Y,
        c.lastPt.X, c.lastPt.Y, c.lastVec.X, c.lastVec.Y)
    result1 := c.addPt(c.firstPt)
    log.Printf("close: after addPt, lastVec=(%v,%v), result1=%v",
        c.lastVec.X, c.lastVec.Y, result1)
    result2 := c.addVec(c.firstVec)
    log.Printf("close: after addVec, result2=%v, final=%v", result2, result1 && result2)
    return result1 && result2
}
```

### 2. Create Unit Test for Convexicator

Create a standalone test that directly tests the convexicator:

```go
func TestConvexicator_ConcaveQuadrilateral(t *testing.T) {
    c := newConvexicator()
    
    c.setMovePt(Point{X: 0, Y: 0})
    assert.True(t, c.addPt(Point{X: 10, Y: 10}))
    assert.True(t, c.addPt(Point{X: 10, Y: 0}))
    assert.True(t, c.addPt(Point{X: 0, Y: 10}))
    assert.False(t, c.close()) // Should return false for concave path
}
```

### 3. Compare Step-by-Step with C++

Run the C++ test case with a debugger and compare:
- State after each `addPt()` call
- State after each `addVec()` call
- State during `close()` execution
- Identify the first point where behavior diverges

### 4. Verify Path Processing

Add logging to `computeConvexity()`:

```go
func (p *pathImpl) computeConvexity() enums.PathConvexity {
    // ... existing code ...
    for _, verb := range verbs {
        log.Printf("computeConvexity: verb=%v, pointIdx=%v, contourCount=%v, needsClose=%v",
            verb, pointIdx, contourCount, needsClose)
        // ... rest of processing
    }
}
```

### 5. Check for Edge Cases

Investigate:
- What happens if the path is explicitly closed vs implicitly closed?
- What if `addPt(fFirstPt)` in `close()` is a no-op (when `lastPt == firstPt`)?
- What is the state of `lastVec` in that case?

---

## Next Steps

1. **Immediate:** Add comprehensive logging to trace execution
2. **Short-term:** Create unit test for convexicator methods independently
3. **Medium-term:** Compare step-by-step with C++ implementation
4. **Long-term:** Fix the identified bug and verify all tests pass

---

## References

- **C++ Source:** `skia-source/src/core/SkPathPriv.cpp` lines 419-560
- **Go Implementation:** `skia/impl/path_models.go` lines 164-275
- **Go Path Processing:** `skia/impl/path.go` lines 680-760
- **Test Cases:** `skia/impl/path_test.go` lines 262-295
- **C++ Tests:** `skia-source/tests/PathTest.cpp` lines 1692-1710
- **Story:** `docs/stories/1.3.fix-convexicator-bug.md`

---

**Document Version:** 1.0  
**Last Updated:** 2025-01-27  
**Status:** Under Investigation

