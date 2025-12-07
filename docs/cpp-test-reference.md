# C++ Test Reference Guide
## Skia Test Cases to Match in Go Implementation

**Purpose:** This document lists all C++ tests from Skia source that need to be ported/matched in the Go implementation. Use this as a reference when verifying functional parity and when in doubt about expected behavior.

**Source Version:** Skia chrome/m144  
**Last Updated:** 2025-01-27

---

## Test File Locations

### Primary Test Files
- **Matrix Tests:** `skia-source/tests/MatrixTest.cpp`
- **Paint Tests:** `skia-source/tests/PaintTest.cpp`
- **Path Tests:** `skia-source/tests/PathTest.cpp`

### Related Test Files
- `skia-source/tests/EmptyPathTest.cpp` - Empty path edge cases
- `skia-source/tests/DrawPathTest.cpp` - Path drawing tests
- `skia-source/tests/MatrixProcsTest.cpp` - Matrix procedure tests
- `skia-source/tests/PathOpsExtendedTest.cpp` - Path operations (advanced)

---

## Matrix Tests (`MatrixTest.cpp`)

### Core Test Functions

#### `test_set9(reporter)`
**Location:** Line 122  
**Purpose:** Tests setting matrix from 9-element buffer  
**What it tests:**
- `set9()` with identity matrix
- `set9()` with scale matrix
- `set9()` after post-translate
- `get9()` returns correct values
- `rc(row, col)` accessor matches `get9()` values

**Go Equivalent:** Test `Set9()` and `Get9()` methods

---

#### `test_matrix_recttorect(reporter)`
**Location:** Line 144  
**Purpose:** Tests `Rect2Rect()` static factory method  
**What it tests:**
- Identity matrix when src == dst
- Translate-only matrix when dst is offset
- Scale+translate matrix
- Scale-only matrix
- Failure case with empty src rect
- `setRectToRect()` failure resets to identity

**Go Equivalent:** Test `MapRectToRect()` method

---

#### `test_flatten(reporter, const SkMatrix& m)`
**Location:** Line 190  
**Purpose:** Tests matrix serialization/flattening  
**What it tests:**
- `WriteToMemory()` size calculation
- `ReadFromMemory()` reconstruction
- Round-trip serialization preserves matrix
- Multiple serializations produce identical output

**Go Equivalent:** Not applicable (no serialization in current API)

---

#### `test_matrix_min_max_scale(reporter)`
**Location:** Line 210  
**Purpose:** Tests scale factor calculations  
**What it tests:**
- Identity matrix: min=1, max=1
- Scale matrix: min/max match scale factors
- Rotated+scaled matrix
- Pure rotation: min≈1, max≈1
- Translate-only: min=1, max=1
- Perspective matrices: returns -1 (invalid)
- Edge cases: very large values, negative nearly-zeros
- `getMinMaxScales()` output array
- Relationship between `getMinScale()`, `getMaxScale()`, and `getMinMaxScales()`
- Vector scaling bounds verification

**Go Equivalent:** Not applicable (methods not in current API - advanced feature)

---

#### `test_matrix_preserve_shape(reporter)`
**Location:** Line 345  
**Purpose:** Tests shape preservation queries  
**What it tests:**
- `isSimilarity()` - preserves circles (uniform scale + rotation)
- `preservesRightAngles()` - preserves rectangles
- Identity: both true
- Translation: both true
- Uniform scale: both true
- Non-uniform scale: similarity=false, rightAngles=true
- Skew: both false
- Scale at pivot point
- Skew at pivot point

**Go Equivalent:** Test `PreservesRightAngles()` method (similarity not implemented)

---

#### `test_matrix_decomposition(reporter)`
**Location:** Line 518  
**Purpose:** Tests matrix decomposition  
**What it tests:**
- Decomposition of various matrix types
- Reconstruction from decomposed parts

**Go Equivalent:** Not applicable (advanced feature)

---

#### `test_matrix_homogeneous(reporter)`
**Location:** Line 682  
**Purpose:** Tests homogeneous coordinate transformations  
**What it tests:**
- Perspective transformations
- Homogeneous point mapping
- W-plane clipping

**Go Equivalent:** Not applicable (advanced feature)

---

#### `test_decompScale(reporter)`
**Location:** Line 841  
**Purpose:** Tests scale decomposition  
**What it tests:**
- Extracting scale factors from matrices
- Scale decomposition accuracy

**Go Equivalent:** Not applicable (advanced feature)

---

### Main Test Cases (DEF_TEST)

#### `DEF_TEST(Matrix, reporter)`
**Location:** Line 861  
**Purpose:** Comprehensive matrix operations test  
**What it tests:**
- Matrix inversion:
  - Identity after M * M^-1
  - Identity after M^-1 * M
  - Scale matrix inversion
  - Scale+translate inversion
  - Scale+rotate inversion
  - Zero scale (non-invertible)
  - Non-finite inversion results
  - NaN handling in scale+translate paths
- `rectStaysRect()` with various 2x2 matrices
- `asAffine()` conversion
- Zero sign differences (-0 vs +0)
- NaN equality (NaN != NaN)
- Calls all helper test functions

**Go Equivalent:** Primary test suite - port all assertions

---

#### `DEF_TEST(Matrix_Concat, r)`
**Location:** Line 1028  
**Purpose:** Tests matrix concatenation  
**What it tests:**
- `setConcat(a, b)` matches `Concat(a, b)` static method
- Translate * Scale concatenation

**Go Equivalent:** Test `SetConcat()` method

---

#### `DEF_TEST(Matrix_maprects, r)`
**Location:** Line 1042  
**Purpose:** Tests all mapRect variants  
**What it tests:**
- `mapPoints()` on rect corners matches `mapRect()`
- `mapRectScaleTranslate()` optimization
- `mapRect()` return value
- Non-finite rect handling after mapping
- Large scale factors (1e20)

**Go Equivalent:** Test `MapRect()` and `MapPoints()` consistency

---

#### `DEF_TEST(Matrix_mapRect_skbug12335, r)`
**Location:** Line 1082  
**Purpose:** Bug fix regression test  
**What it tests:**
- Perspective matrix mapping with very small w values
- Rect should not become empty when w is small but positive
- Specific matrix values that caused bug

**Go Equivalent:** Edge case test for perspective mapping

---

#### `DEF_TEST(Matrix_Ctor, r)`
**Location:** Line 1096  
**Purpose:** Tests default constructor  
**What it tests:**
- Default constructor equals identity matrix

**Go Equivalent:** Test `NewMatrixIdentity()`

---

#### `DEF_TEST(Matrix_LookAt, r)`
**Location:** Line 1100  
**Purpose:** Tests LookAt matrix (3D)  
**What it tests:**
- Degenerate LookAt inputs don't crash
- Returns identity for degenerate case

**Go Equivalent:** Not applicable (3D matrices not in scope)

---

#### `DEF_TEST(Matrix_SetRotateSnap, r)`
**Location:** Line 1106  
**Purpose:** Tests rotation snapping  
**What it tests:**
- Rotations by multiples of 90° snap correctly
- Sin/cos values snap to exact 0, 1, -1

**Go Equivalent:** Test `SetRotate()` with 90° multiples

---

#### `DEF_TEST(Matrix_rectStaysRect_zeroScale, r)`
**Location:** Line 1125  
**Purpose:** Tests zero scale edge case  
**What it tests:**
- Zero scale matrices don't crash `rectStaysRect()`
- Zero scale returns false for `rectStaysRect()`

**Go Equivalent:** Edge case test for `RectStaysRect()`

---

## Paint Tests (`PaintTest.cpp`)

### Main Test Cases

#### `DEF_TEST(Paint_copy, reporter)`
**Location:** Line 44  
**Purpose:** Tests paint copying  
**What it tests:**
- Copy constructor creates equal paint
- Assignment operator creates equal paint
- `reset()` returns to initial state
- Style, stroke width, mask filter copying

**Go Equivalent:** Test paint equality and reset

---

#### `DEF_TEST(Paint_regression_cubic, reporter)`
**Location:** Line 71  
**Purpose:** Regression test for cubic path handling  
**What it tests:**
- Stroke bounds don't explode with degenerate cubic
- Stroke bounds stay within expected maximum
- Miter limit calculation

**Go Equivalent:** Test `ComputeFastStrokeBounds()` with edge cases

---

#### `DEF_TEST(Paint_flattening, reporter)`
**Location:** Line 103  
**Purpose:** Tests paint serialization  
**What it tests:**
- All cap types serialize correctly
- All join types serialize correctly
- All style types serialize correctly
- Round-trip serialization preserves paint

**Go Equivalent:** Not applicable (no serialization in current API)

---

#### `DEF_TEST(Paint_regression_measureText, reporter)`
**Location:** Line 149  
**Purpose:** Regression test for text measurement  
**What it tests:**
- Empty string measurement resets rect
- Rect initialized correctly for zero-length strings

**Go Equivalent:** Not applicable (text measurement not in scope)

---

#### `DEF_TEST(Paint_MoreFlattening, r)`
**Location:** Line 164  
**Purpose:** Additional serialization tests  
**What it tests:**
- Color serialization
- Blend mode serialization
- Round-trip accuracy

**Go Equivalent:** Not applicable (no serialization)

---

#### `DEF_TEST(Paint_nothingToDraw, r)`
**Location:** Line 184  
**Purpose:** Tests `nothingToDraw()` method  
**What it tests:**
- Default paint: `nothingToDraw()` = false
- Zero alpha: `nothingToDraw()` = true
- Dst blend mode: `nothingToDraw()` = true
- Color filter that preserves alpha: `nothingToDraw()` = true
- Color filter that modifies alpha: `nothingToDraw()` = false

**Go Equivalent:** Test `NothingToDraw()` method

---

#### `DEF_TEST(Font_getpos, r)`
**Location:** Line 208  
**Purpose:** Font positioning tests  
**What it tests:**
- Font glyph positioning
- Subpixel positioning
- Hinting modes

**Go Equivalent:** Not applicable (font API not in scope)

---

#### `DEF_TEST(Paint_dither, reporter)`
**Location:** Line 246  
**Purpose:** Tests dithering  
**What it tests:**
- Dither flag setting
- Dither decision logic

**Go Equivalent:** Test `SetDither()` and `IsDither()` methods

---

## Path Tests (`PathTest.cpp`)

### Helper Test Functions (100+ functions)

#### Core Functionality Tests

##### `test_addrect(reporter)`
**Location:** Line 688  
**Purpose:** Tests `addRect()` method  
**What it tests:**
- Rect addition with various directions
- Start index handling
- Bounds calculation
- Verb sequence correctness

**Go Equivalent:** Test `AddRect()` method

---

##### `test_addrect_isfinite(reporter)`
**Location:** Line 721  
**Purpose:** Tests finite rect addition  
**What it tests:**
- Non-finite rect handling
- Empty rect handling
- Infinite rect handling

**Go Equivalent:** Edge case test for `AddRect()`

---

##### `test_bounds(reporter)`
**Location:** Line 1234  
**Purpose:** Tests bounds calculation  
**What it tests:**
- Empty path bounds
- Single point bounds
- Line bounds
- Curve bounds
- Multiple contour bounds
- Bounds invalidation after edits
- `updateBoundsCache()` behavior

**Go Equivalent:** Test `Bounds()` and `UpdateBoundsCache()` methods

---

##### `test_close(reporter)`
**Location:** Line 1318  
**Purpose:** Tests `close()` method  
**What it tests:**
- Closing open contours
- Multiple close calls
- Close after moveTo
- Close with no points
- Bounds after close

**Go Equivalent:** Test `Close()` method

---

##### `test_convexity(reporter)`
**Location:** Line 1670  
**Purpose:** Tests convexity detection  
**What it tests:**
- Simple convex shapes
- Concave shapes
- Self-intersecting paths
- Degenerate paths
- Convexity caching
- Convexity invalidation

**Go Equivalent:** Test `Convexity()` and `IsConvex()` methods

---

##### `test_convexity2(reporter)`
**Location:** Line 1406  
**Purpose:** Additional convexity tests  
**What it tests:**
- Complex convex shapes
- Edge cases for convexity

**Go Equivalent:** Extended convexity tests

---

##### `test_convexity_doubleback(reporter)`
**Location:** Line 1573  
**Purpose:** Tests double-back paths  
**What it tests:**
- Paths that double back on themselves
- Convexity with overlapping segments

**Go Equivalent:** Edge case convexity test

---

##### `test_direction(reporter)`
**Location:** Line 1138  
**Purpose:** Tests path direction  
**What it tests:**
- Clockwise vs counter-clockwise
- Direction with `addRect()`
- Direction with `addOval()`
- Direction with `addCircle()`
- Direction with `addRRect()`

**Go Equivalent:** Test direction parameter in shape addition methods

---

##### `test_isLine(reporter)`
**Location:** Line 1857  
**Purpose:** Tests `isLine()` method  
**What it tests:**
- Single line detection
- Move + Line verb sequence
- Non-line paths return false
- Empty paths return false

**Go Equivalent:** Test `IsLine()` method

---

##### `test_isRect(reporter)`
**Location:** Line 2172  
**Purpose:** Tests `isRect()` method  
**What it tests:**
- Rect detection
- Rect with different fill types
- Non-rect paths
- Rect bounds extraction
- Trailing moveTo handling

**Go Equivalent:** Not applicable (method not in current API)

---

##### `test_transform(reporter)`
**Location:** Line 2825  
**Purpose:** Tests path transformation  
**What it tests:**
- Transform with identity matrix
- Transform with scale matrix
- Transform with rotation matrix
- Transform with perspective matrix
- Transform preserves path structure
- Bounds after transform

**Go Equivalent:** Test `Transform()` method

---

##### `test_addPath(reporter)`
**Location:** Line 4024  
**Purpose:** Tests `addPath()` method  
**What it tests:**
- Adding path with offset
- Adding path with matrix
- Adding empty path
- Adding path to empty path
- Verb and point copying

**Go Equivalent:** Test `AddPath()`, `AddPathNoOffset()`, `AddPathMatrix()` methods

---

##### `test_addPathMode(reporter, bool explicitMoveTo, bool extend)`
**Location:** Line 4044  
**Purpose:** Tests add path modes  
**What it tests:**
- Append mode behavior
- Extend mode behavior
- Explicit moveTo handling
- Line injection in extend mode

**Go Equivalent:** Test `AddPathMode` enum values

---

##### `test_circle(reporter)`
**Location:** Line 3612  
**Purpose:** Tests circle addition  
**What it tests:**
- Circle with different radii
- Circle with different directions
- Circle bounds
- Circle point count
- Circle verb sequence

**Go Equivalent:** Test `AddCircle()` method

---

##### `test_oval(reporter)`
**Location:** Line 3643  
**Purpose:** Tests oval addition  
**What it tests:**
- Oval with different rects
- Oval with different directions
- Oval bounds
- Oval point count

**Go Equivalent:** Test `AddOval()` method

---

##### `test_rrect(reporter)`
**Location:** Line 3733  
**Purpose:** Tests rounded rect addition  
**What it tests:**
- RRect with various corner radii
- RRect with different directions
- Degenerate RRect (rect)
- Degenerate RRect (oval)
- RRect bounds

**Go Equivalent:** Test `AddRRect()` method

---

##### `test_arc(reporter)`
**Location:** Line 3797  
**Purpose:** Tests arc addition  
**What it tests:**
- Arc with various angles
- Arc sweep direction
- Arc bounds

**Go Equivalent:** Not applicable (arc methods not in current API)

---

##### `test_arcTo(reporter)`
**Location:** Line 3959  
**Purpose:** Tests `arcTo()` method  
**What it tests:**
- ArcTo with various radii
- ArcTo sweep flags
- ArcTo point sequences

**Go Equivalent:** Not applicable (arcTo not in current API)

---

##### `test_flattening(reporter)`
**Location:** Line 2778  
**Purpose:** Tests path serialization  
**What it tests:**
- Serialization of all verb types
- Point array serialization
- Conic weight serialization
- Fill type serialization
- Round-trip accuracy

**Go Equivalent:** Not applicable (no serialization)

---

##### `test_iter(reporter)`
**Location:** Line 3051  
**Purpose:** Tests path iteration  
**What it tests:**
- Iterating through all verbs
- Point extraction during iteration
- Conic weight extraction
- Iteration with empty path
- Iteration edge cases

**Go Equivalent:** Not applicable (iteration API not in current scope)

---

##### `test_segment_masks(reporter)`
**Location:** Line 2993  
**Purpose:** Tests segment mask queries  
**What it tests:**
- Line segment mask
- Quad segment mask
- Conic segment mask
- Cubic segment mask
- Combined masks

**Go Equivalent:** Not applicable (segment masks not in current API)

---

##### `test_zero_length_paths(reporter)`
**Location:** Line 2922  
**Purpose:** Tests zero-length path handling  
**What it tests:**
- Zero-length lines
- Zero-length curves
- Bounds of zero-length paths
- Convexity of zero-length paths

**Go Equivalent:** Edge case tests

---

##### `test_isfinite(reporter)`
**Location:** Line 1014  
**Purpose:** Tests finite path detection  
**What it tests:**
- Path with finite points: `isFinite()` = true
- Path with NaN points: `isFinite()` = false
- Path with Inf points: `isFinite()` = false
- Mixed finite/non-finite points

**Go Equivalent:** Test `IsFinite()` method

---

##### `test_isfinite_after_transform(reporter)`
**Location:** Line 842  
**Purpose:** Tests finiteness after transform  
**What it tests:**
- Transform that produces NaN
- Transform that produces Inf
- Finite path + finite matrix = finite result
- Non-finite detection after transform

**Go Equivalent:** Test `IsFinite()` after `Transform()`

---

### Bug Regression Tests

#### `test_skbug_3469(reporter)`
**Location:** Line 79  
**Purpose:** Bug fix regression  
**What it tests:** Specific bug scenario

#### `test_skbug_3239(reporter)`
**Location:** Line 88  
**Purpose:** Bug fix regression  
**What it tests:** Specific bug scenario

#### `test_path_crbug364224()`
**Location:** Line 172  
**Purpose:** Chrome bug regression  
**What it tests:** Specific crash scenario

#### `test_fuzz_crbug_638223()`
**Location:** Line 193  
**Purpose:** Fuzzer-found bug  
**What it tests:** Specific crash scenario

#### `test_fuzz_crbug_643933()`
**Location:** Line 203  
**Purpose:** Fuzzer-found bug  
**What it tests:** Specific crash scenario

#### `test_fuzz_crbug_647922()`
**Location:** Line 220  
**Purpose:** Fuzzer-found bug  
**What it tests:** Specific crash scenario

#### `test_fuzz_crbug_662780()`
**Location:** Line 230  
**Purpose:** Fuzzer-found bug  
**What it tests:** Specific crash scenario

#### `test_fuzz_crbug_668907()`
**Location:** Line 274  
**Purpose:** Fuzzer-found bug  
**What it tests:** Specific crash scenario

#### `test_fuzz_crbug_627414(reporter)`
**Location:** Line 813  
**Purpose:** Fuzzer-found bug  
**What it tests:** Specific crash scenario

#### `test_fuzz_crbug_662952(reporter)`
**Location:** Line 4571  
**Purpose:** Fuzzer-found bug  
**What it tests:** Specific crash scenario

#### `test_fuzz_crbug_662730(reporter)`
**Location:** Line 4612  
**Purpose:** Fuzzer-found bug  
**What it tests:** Specific crash scenario

**Go Equivalent:** Port all fuzzer-found bug tests as edge case tests

---

### Main Test Cases (DEF_TEST)

#### `DEF_TEST(PathInterp, reporter)`
**Location:** Line 4890  
**Purpose:** Tests path interpolation  
**What it tests:**
- `isInterpolatable()` method
- `makeInterpolate()` method
- Interpolation with various weights

**Go Equivalent:** Not applicable (interpolation not in current API)

---

#### `DEF_TEST(Path_multipleMoveTos, reporter)`
**Location:** Line 4895  
**Purpose:** Tests multiple moveTo calls  
**What it tests:**
- Multiple consecutive moveTo calls
- Last point tracking
- Bounds with multiple moves

**Go Equivalent:** Test `MoveTo()` behavior

---

#### `DEF_TEST(PathBigCubic, reporter)`
**Location:** Line 4915  
**Purpose:** Tests large cubic curves  
**What it tests:**
- Very large coordinate values
- Degenerate cubic handling
- No assertion/crash with extreme values

**Go Equivalent:** Edge case test for `CubicTo()`

---

#### `DEF_TEST(PathContains, reporter)`
**Location:** Line 4931  
**Purpose:** Tests point containment  
**What it tests:**
- `contains()` method with various points
- Fill type affects containment
- Edge cases

**Go Equivalent:** Not applicable (contains not in current API)

---

#### `DEF_TEST(Paths, reporter)`
**Location:** Line 4935  
**Purpose:** Comprehensive path test suite  
**What it tests:**
- Calls all helper test functions
- Self-assignment
- Swap operation
- Empty path behavior
- Bounds calculation
- Segment masks
- Point/verb extraction

**Go Equivalent:** Primary test suite - port all assertions

---

#### `DEF_TEST(conservatively_contains_rect, reporter)`
**Location:** Line 5116  
**Purpose:** Tests conservative rect containment  
**What it tests:**
- `conservativelyContainsRect()` method
- Conservative vs exact containment

**Go Equivalent:** Not applicable (method not in current API)

---

#### `DEF_TEST(skbug_6450, r)`
**Location:** Line 5132  
**Purpose:** Bug fix regression  
**What it tests:** Specific bug scenario

---

#### `DEF_TEST(PathRefSerialization, reporter)`
**Location:** Line 5155  
**Purpose:** Tests path reference serialization  
**What it tests:** Internal serialization

**Go Equivalent:** Not applicable (internal API)

---

#### `DEF_TEST(NonFinitePathIteration, reporter)`
**Location:** Line 5189  
**Purpose:** Tests iteration with non-finite paths  
**What it tests:**
- Iteration doesn't crash with NaN/Inf
- Graceful handling of non-finite values

**Go Equivalent:** Edge case test

---

#### `DEF_TEST(AndroidArc, reporter)`
**Location:** Line 5197  
**Purpose:** Android-specific arc test  
**What it tests:** Arc behavior on Android

**Go Equivalent:** Not applicable (platform-specific)

---

#### `DEF_TEST(HugeGeometry, reporter)`
**Location:** Line 5222  
**Purpose:** Tests very large geometry  
**What it tests:**
- Very large coordinates
- Overflow handling
- Bounds with huge values

**Go Equivalent:** Edge case test

---

#### `DEF_TEST(ClipPath_nonfinite, reporter)`
**Location:** Line 5254  
**Purpose:** Tests clipping with non-finite paths  
**What it tests:**
- Clipping doesn't crash with NaN/Inf
- Graceful degradation

**Go Equivalent:** Edge case test

---

#### `DEF_TEST(Path_isRect, reporter)`
**Location:** Line 5283  
**Purpose:** Tests `isRect()` method  
**What it tests:**
- Various rect detection scenarios
- Edge cases

**Go Equivalent:** Not applicable (method not in current API)

---

#### `DEF_TEST(Path_self_add, reporter)`
**Location:** Line 5438  
**Purpose:** Tests adding path to itself  
**What it tests:**
- Self-addition doesn't cause infinite loop
- Correct behavior when adding self

**Go Equivalent:** Edge case test for `AddPath()`

---

#### `DEF_TEST(triangle_onehalf, reporter)`
**Location:** Line 5474  
**Purpose:** Tests triangle path  
**What it tests:** Specific triangle scenario

---

#### `DEF_TEST(triangle_big, reporter)`
**Location:** Line 5485  
**Purpose:** Tests large triangle  
**What it tests:** Large coordinate handling

---

#### `DEF_TEST(Path_setLastPt, r)`
**Location:** Line 5502  
**Purpose:** Tests setting last point  
**What it tests:**
- `setLastPoint()` method
- Last point tracking

**Go Equivalent:** Not applicable (method not in current API)

---

#### `DEF_TEST(Path_increserve_handle_neg_crbug_883666, r)`
**Location:** Line 5518  
**Purpose:** Bug fix regression  
**What it tests:** Negative reserve handling

---

#### `DEF_TEST(Path_survive_transform, r)`
**Location:** Line 5596  
**Purpose:** Tests path survival after transform  
**What it tests:**
- Path remains valid after transform
- No corruption

**Go Equivalent:** Test `Transform()` robustness

---

#### `DEF_TEST(path_last_move_to_index, r)`
**Location:** Line 5617  
**Purpose:** Tests last moveTo index tracking  
**What it tests:**
- Last moveTo index correctness
- Index updates correctly

**Go Equivalent:** Not applicable (internal detail)

---

#### `DEF_TEST(pathedger, r)`
**Location:** Line 5806  
**Purpose:** Tests path edge generation  
**What it tests:** Edge generation for rendering

**Go Equivalent:** Not applicable (rendering detail)

---

#### `DEF_TEST(path_addpath_crbug_1153516, r)`
**Location:** Line 5826  
**Purpose:** Bug fix regression  
**What it tests:** Specific addPath bug

---

#### `DEF_TEST(path_convexity_scale_way_down, r)`
**Location:** Line 5852  
**Purpose:** Tests convexity with very small scale  
**What it tests:**
- Convexity calculation with tiny values
- Numerical stability

**Go Equivalent:** Edge case test for `Convexity()`

---

#### `DEF_TEST(path_moveto_addrect, r)`
**Location:** Line 5865  
**Purpose:** Tests moveTo followed by addRect  
**What it tests:**
- MoveTo + AddRect interaction
- Bounds calculation

**Go Equivalent:** Test `MoveTo()` + `AddRect()` sequence

---

#### `DEF_TEST(path_moveto_twopass_convexity, r)`
**Location:** Line 5900  
**Purpose:** Tests convexity with two-pass calculation  
**What it tests:** Convexity caching

**Go Equivalent:** Test convexity caching

---

#### `DEF_TEST(path_walk_simple_edges_1154864, r)`
**Location:** Line 5920  
**Purpose:** Bug fix regression  
**What it tests:** Edge walking bug

---

#### `DEF_TEST(path_walk_edges_concave_large_dx, r)`
**Location:** Line 5939  
**Purpose:** Tests edge walking with large dx  
**What it tests:** Numerical stability

---

#### `DEF_TEST(path_filltype_utils, r)`
**Location:** Line 5958  
**Purpose:** Tests fill type utility functions  
**What it tests:**
- Fill type conversions
- Inverse fill type handling

**Go Equivalent:** Test fill type methods

---

#### `DEF_TEST(path_computeTightBounds, reporter)`
**Location:** Line 6079  
**Purpose:** Tests tight bounds calculation  
**What it tests:**
- `computeTightBounds()` accuracy
- Tight vs regular bounds difference

**Go Equivalent:** Test `ComputeTightBounds()` method

---

#### `DEF_TEST(path_trivial_isrect, reporter)`
**Location:** Line 6107  
**Purpose:** Tests trivial rect detection  
**What it tests:** Simple rect cases

---

#### `DEF_TEST(path_infinite_transform, reporter)`
**Location:** Line 6152  
**Purpose:** Tests transform with infinite matrix  
**What it tests:**
- Handling of infinite matrix values
- Graceful degradation

**Go Equivalent:** Edge case test for `Transform()`

---

#### `DEF_TEST(path_factory_inverted_bounds, reporter)`
**Location:** Line 6197  
**Purpose:** Tests factory methods with inverted bounds  
**What it tests:**
- Rect with left>right, top>bottom
- Factory method handling

**Go Equivalent:** Edge case test for static factories

---

## Test Coverage Summary

### Matrix Tests
- **Total Test Functions:** ~16
- **Core Functionality:** ✅ High coverage
- **Edge Cases:** ✅ Good coverage
- **Advanced Features:** ⚠️ Some not applicable

### Paint Tests
- **Total Test Functions:** ~8
- **Core Functionality:** ✅ Good coverage
- **Edge Cases:** ✅ Some coverage
- **Advanced Features:** ⚠️ Some not applicable

### Path Tests
- **Total Test Functions:** ~100+
- **Core Functionality:** ✅ Excellent coverage
- **Edge Cases:** ✅ Excellent coverage (many fuzzer tests)
- **Advanced Features:** ⚠️ Some not applicable

---

## Priority for Porting

### High Priority (Core Functionality)
1. ✅ Matrix inversion tests (`DEF_TEST(Matrix, reporter)`)
2. ✅ Matrix concatenation tests
3. ✅ Matrix mapRect tests
4. ✅ Paint equality and reset tests
5. ✅ Paint `nothingToDraw()` tests
6. ✅ Path bounds tests
7. ✅ Path convexity tests
8. ✅ Path transform tests
9. ✅ Path addPath tests
10. ✅ Path shape addition tests (rect, oval, circle, rrect)

### Medium Priority (Edge Cases)
1. ⚠️ Zero scale matrix tests
2. ⚠️ Non-finite value tests
3. ⚠️ Empty path tests
4. ⚠️ Zero-length path tests
5. ⚠️ Fuzzer-found bug tests

### Low Priority (Advanced/Not Applicable)
1. ❌ Serialization tests (no serialization API)
2. ❌ Font tests (not in scope)
3. ❌ 3D matrix tests (not in scope)
4. ❌ Path iteration tests (API not implemented)
5. ❌ Path contains tests (API not implemented)

---

## Usage Notes

1. **When in doubt:** Check the C++ test file for expected behavior
2. **Edge cases:** Many fuzzer tests reveal important edge cases
3. **Bug regressions:** Port bug fix tests to prevent regressions
4. **Test helpers:** Many helper functions can be ported as test utilities
5. **Test data:** Extract test matrices, paths, and paints from tests

---

## Test Helper Functions to Port

### Matrix Helpers
- `nearly_equal_scalar()` - Floating point comparison with tolerance
- `nearly_equal()` - Matrix comparison with tolerance
- `is_identity()` - Identity matrix check
- `assert9()` - Matrix value assertion helper

### Path Helpers
- `test_empty()` - Empty path verification
- Various bounds checking helpers

---

**Document Status:** Reference Guide - Keep Updated as Tests are Ported  
**Next Update:** After initial test porting begins

