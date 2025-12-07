# API Parity Verification
## Ensuring Go Interfaces Match Skia C++ API

**Objective:** Verify that Go interfaces match Skia C++ API so developers familiar with Skia feel at home when using this library.

**Source Version:** Skia chrome/m144  
**Date:** 2025-01-27

---

## Verification Methodology

### 1. Method-by-Method Comparison
- Compare each public method in C++ headers with Go interfaces
- Verify method names match (or follow Go conventions appropriately)
- Verify method signatures match (parameters, return types)
- Document intentional differences and rationale

### 2. Static Factory Methods
- Verify all static factory methods are present
- Check constructor equivalents exist

### 3. Operator Overloads
- Document how C++ operators are handled in Go
- Ensure equivalent functionality exists

### 4. Method Grouping & Organization
- Verify logical grouping matches C++ API
- Check that related methods are together

---

## SkMatrix API Comparison

### C++ Public API (from `include/core/SkMatrix.h`)

#### Static Factory Methods
| C++ Method | Go Equivalent | Status | Notes |
|-----------|---------------|--------|-------|
| `SkMatrix()` (default constructor) | `NewMatrixIdentity()` | ✅ | Go uses explicit factory |
| `SkMatrix::Scale(sx, sy)` | `NewMatrixScale(sx, sy)` | ✅ | |
| `SkMatrix::Translate(dx, dy)` | `NewMatrixTranslate(dx, dy)` | ✅ | |
| `SkMatrix::RotateDeg(deg)` | `NewMatrixRotate(deg)` | ✅ | |
| `SkMatrix::RotateDeg(deg, pt)` | `NewMatrixRotate(deg)` + `SetRotate(deg, px, py)` | ⚠️ | Missing static factory with pivot |
| `SkMatrix::RotateRad(rad)` | ❌ | ❌ | Missing - could add `NewMatrixRotateRad` |
| `SkMatrix::Skew(kx, ky)` | `NewMatrixSkew(kx, ky)` | ✅ | |
| `SkMatrix::MakeAll(...)` | ❌ | ❌ | Missing - should add `NewMatrixAll` |
| `SkMatrix::ScaleTranslate(sx, sy, tx, ty)` | ❌ | ❌ | Missing static factory |

#### Instance Methods - Getters
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `getScaleX()` | `GetScaleX()` | ✅ | |
| `getScaleY()` | `GetScaleY()` | ✅ | |
| `getSkewX()` | `getSkewX()` | ✅ | |
| `getSkewY()` | `getSkewY()` | ✅ | |
| `getTranslateX()` | `GetTranslateX()` | ✅ | |
| `getTranslateY()` | `GetTranslateY()` | ✅ | |
| `getPerspX()` | `GetPerspX()` | ✅ | |
| `getPerspY()` | `GetPerspY()` | ✅ | |
| `getType()` | `GetType()` | ✅ | |
| `operator[](int)` | ❌ | ❌ | Missing - Go doesn't support operators |
| `get(int)` | ❌ | ❌ | Missing - should add `Get(index int)` |
| `rc(int r, int c)` | ❌ | ❌ | Missing - should add `GetRC(row, col int)` |
| `get9(SkScalar buffer[9])` | ❌ | ❌ | Missing - should add `Get9() [9]Scalar` |

#### Instance Methods - Setters
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `setScale(sx, sy)` | `SetScale(sx, sy)` | ✅ | |
| `setTranslate(dx, dy)` | `SetTranslate(dx, dy)` | ✅ | |
| `setRotate(deg)` | `SetRotate(deg, 0, 0)` | ⚠️ | Requires explicit pivot (0,0) |
| `setRotate(deg, px, py)` | `SetRotate(degrees, px, py)` | ✅ | |
| `setSkew(kx, ky)` | `SetSkew(kx, ky)` | ✅ | |
| `setAll(...)` | ❌ | ❌ | Missing - should add `SetAll(...)` |
| `set9(const SkScalar buffer[9])` | ❌ | ❌ | Missing - should add `Set9([9]Scalar)` |
| `set(int index, SkScalar value)` | ❌ | ❌ | Missing - should add `Set(index int, value Scalar)` |
| `setScaleX(v)` | ❌ | ❌ | Missing - should add `SetScaleX(v)` |
| `setScaleY(v)` | ❌ | ❌ | Missing - should add `SetScaleY(v)` |
| `setSkewX(v)` | ❌ | ❌ | Missing - should add `SetSkewX(v)` |
| `setSkewY(v)` | ❌ | ❌ | Missing - should add `SetSkewY(v)` |
| `setTranslateX(v)` | ❌ | ❌ | Missing - should add `SetTranslateX(v)` |
| `setTranslateY(v)` | ❌ | ❌ | Missing - should add `SetTranslateY(v)` |
| `setPerspX(v)` | ❌ | ❌ | Missing - should add `SetPerspX(v)` |
| `setPerspY(v)` | ❌ | ❌ | Missing - should add `SetPerspY(v)` |
| `reset()` | `Reset()` | ✅ | |
| `setIdentity()` | `SetIdentity()` | ✅ | |

#### Instance Methods - Transformations
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `preTranslate(dx, dy)` | `PreTranslate(dx, dy)` | ✅ | |
| `preScale(sx, sy)` | `PreScale(sx, sy)` | ✅ | |
| `preRotate(deg)` | `PreRotate(deg, 0, 0)` | ⚠️ | Requires explicit pivot |
| `preRotate(deg, px, py)` | `PreRotate(degrees, px, py)` | ✅ | |
| `preSkew(kx, ky)` | `PreSkew(kx, ky)` | ✅ | |
| `preConcat(other)` | `PreConcat(other)` | ✅ | |
| `postTranslate(dx, dy)` | `PostTranslate(dx, dy)` | ✅ | |
| `postScale(sx, sy)` | `PostScale(sx, sy)` | ✅ | |
| `postRotate(deg)` | `PostRotate(deg, 0, 0)` | ⚠️ | Requires explicit pivot |
| `postRotate(deg, px, py)` | `PostRotate(degrees, px, py)` | ✅ | |
| `postSkew(kx, ky)` | `PostSkew(kx, ky)` | ✅ | |
| `postConcat(other)` | `PostConcat(other)` | ✅ | |
| `setConcat(a, b)` | `SetConcat(a, b)` | ✅ | |

#### Instance Methods - Mapping
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `mapPoint(pt)` | `MapPoint(pt)` | ✅ | |
| `mapPoints(dst, src, count)` | `MapPoints(dst, src) int` | ✅ | Returns count instead of taking it |
| `mapRect(rect)` | `MapRect(rect)` | ✅ | |
| `mapRectToRect(src, dst)` | `MapRectToRect(src, dst) bool` | ✅ | |
| `mapHomogeneousPoints(dst, src, count)` | ❌ | ❌ | Missing - advanced feature |
| `mapXY(x, y)` | ❌ | ❌ | Missing - convenience method |

#### Instance Methods - Queries
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `isIdentity()` | `IsIdentity()` | ✅ | |
| `isScaleTranslate()` | `IsScaleTranslate()` | ✅ | |
| `isTranslate()` | ❌ | ❌ | Missing - should add |
| `hasPerspective()` | `HasPerspective()` | ✅ | |
| `rectStaysRect()` | `RectStaysRect()` | ✅ | |
| `preservesAxisAlignment()` | ❌ | ❌ | Missing - alias for `rectStaysRect()` |
| `preservesRightAngles()` | `PreservesRightAngles()` | ✅ | |
| `isSimilarity(tol)` | ❌ | ❌ | Missing - advanced query |

#### Instance Methods - Advanced
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `invert()` | `Invert() (SkMatrix, bool)` | ✅ | Returns tuple instead of bool param |
| `getMinScale()` | ❌ | ❌ | Missing - advanced feature |
| `getMaxScale()` | ❌ | ❌ | Missing - advanced feature |
| `getMinMaxScale(min, max)` | ❌ | ❌ | Missing - advanced feature |

#### Operators (C++ only - no direct equivalent)
| C++ Operator | Go Equivalent | Status | Notes |
|--------------|---------------|--------|-------|
| `operator==` | `Equals()` method? | ❌ | Not in interface - should add |
| `operator!=` | `!Equals()` | ❌ | Not in interface |
| `operator*` (matrix multiply) | `SetConcat()` or helper | ⚠️ | No operator syntax in Go |

---

## SkPaint API Comparison

### C++ Public API (from `include/core/SkPaint.h`)

#### Constructors
| C++ Method | Go Equivalent | Status | Notes |
|-----------|---------------|--------|-------|
| `SkPaint()` | `NewPaint()` | ✅ | |
| `SkPaint(const SkColor4f& color)` | `NewPaintWithColor(color)` | ✅ | |

#### Getters
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `getColor()` | `GetColorInt()` | ✅ | Returns `uint32` |
| `getColor4f()` | `GetColor()` | ✅ | Returns `Color4f` |
| `getAlpha()` | `GetAlpha()` | ✅ | Returns `uint8` |
| `getAlphaf()` | `GetAlphaf()` | ✅ | Returns `Scalar` |
| `getStyle()` | `GetStyle()` | ✅ | |
| `getStrokeWidth()` | `GetStrokeWidth()` | ✅ | |
| `getStrokeCap()` | `GetStrokeCap()` | ✅ | |
| `getStrokeJoin()` | `GetStrokeJoin()` | ✅ | |
| `getStrokeMiter()` | `GetStrokeMiter()` | ✅ | |
| `getBlendMode()` | `AsBlendMode()` | ⚠️ | Different name - returns `(BlendMode, bool)` |
| `getBlender()` | `GetBlender()` | ✅ | |
| `getShader()` | `GetShader()` | ✅ | |
| `getColorFilter()` | `GetColorFilter()` | ✅ | |
| `getImageFilter()` | `GetImageFilter()` | ✅ | |
| `getMaskFilter()` | `GetMaskFilter()` | ✅ | |
| `getPathEffect()` | `GetPathEffect()` | ✅ | |
| `isAntiAlias()` | `IsAntiAlias()` | ✅ | |
| `isDither()` | `IsDither()` | ✅ | |

#### Setters
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `setColor(SkColor)` | `SetColorInt(color)` | ✅ | |
| `setColor(const SkColor4f&)` | `SetColor(color)` | ✅ | |
| `setARGB(a, r, g, b)` | `SetARGB(a, r, g, b)` | ✅ | |
| `setAlpha(U8CPU)` | `SetAlpha(a)` | ✅ | |
| `setAlphaf(float)` | `SetAlphaf(a)` | ✅ | |
| `setStyle(Style)` | `SetStyle(style)` | ✅ | |
| `setStroke(bool)` | `SetStroke(isStroke)` | ✅ | |
| `setStrokeWidth(SkScalar)` | `SetStrokeWidth(width)` | ✅ | |
| `setStrokeCap(Cap)` | `SetStrokeCap(cap)` | ✅ | |
| `setStrokeJoin(Join)` | `SetStrokeJoin(join)` | ✅ | |
| `setStrokeMiter(SkScalar)` | `SetStrokeMiter(limit)` | ✅ | |
| `setBlendMode(BlendMode)` | `SetBlendMode(mode)` | ✅ | |
| `setBlender(SkBlender*)` | `SetBlender(blender)` | ✅ | |
| `setShader(SkShader*)` | `SetShader(shader)` | ✅ | |
| `setColorFilter(SkColorFilter*)` | `SetColorFilter(filter)` | ✅ | |
| `setImageFilter(SkImageFilter*)` | `SetImageFilter(filter)` | ✅ | |
| `setMaskFilter(SkMaskFilter*)` | `SetMaskFilter(filter)` | ✅ | |
| `setPathEffect(SkPathEffect*)` | `SetPathEffect(effect)` | ✅ | |
| `setAntiAlias(bool)` | `SetAntiAlias(aa)` | ✅ | |
| `setDither(bool)` | `SetDither(dither)` | ✅ | |
| `reset()` | `Reset()` | ✅ | |

#### Advanced Methods
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `computeFastBounds(orig, storage)` | `ComputeFastBounds(orig, storage)` | ✅ | |
| `computeFastStrokeBounds(orig, storage)` | `ComputeFastStrokeBounds(orig, storage)` | ✅ | |
| `getInflationRadius(style, matrixScale)` | `GetInflationRadius(style, matrixScale...)` | ✅ | Uses variadic for optional param |
| `nothingToDraw()` | `NothingToDraw()` | ✅ | |
| `canComputeFastBounds()` | `CanComputeFastBounds()` | ✅ | |
| `isSrcOver()` | `IsSrcOver()` | ✅ | |
| `equals(other)` | `Equals(other)` | ✅ | |

---

## SkPath API Comparison

### C++ Public API (from `include/core/SkPath.h`)

#### Static Factory Methods
| C++ Method | Go Equivalent | Status | Notes |
|-----------|---------------|--------|-------|
| `SkPath(SkPathFillType)` | `NewSkPath(fillType)` | ✅ | |
| `SkPath::Rect(...)` | ❌ | ❌ | Missing static factory |
| `SkPath::Oval(...)` | ❌ | ❌ | Missing static factory |
| `SkPath::Circle(...)` | ❌ | ❌ | Missing static factory |
| `SkPath::RRect(...)` | ❌ | ❌ | Missing static factory |
| `SkPath::Polygon(...)` | ❌ | ❌ | Missing static factory |
| `SkPath::Line(a, b)` | ❌ | ❌ | Missing static factory |
| `SkPath::Raw(...)` | ❌ | ❌ | Missing static factory |

#### Getters
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `getFillType()` | `FillType()` | ✅ | |
| `isInverseFillType()` | `IsInverseFillType()` | ✅ | |
| `getBounds()` | `Bounds()` | ✅ | |
| `computeTightBounds()` | `ComputeTightBounds()` | ✅ | |
| `isConvex()` | `IsConvex()` | ✅ | |
| `getConvexity()` | `Convexity()` | ✅ | |
| `isEmpty()` | `IsEmpty()` | ✅ | |
| `isFinite()` | `IsFinite()` | ✅ | |
| `isLine()` | `IsLine()` | ✅ | |
| `countPoints()` | `CountPoints()` | ✅ | |
| `getPoint(index)` | `Point(index)` | ✅ | |
| `getPoints(points, max)` | `GetPoints(points) int` | ✅ | Returns count |
| `countVerbs()` | `CountVerbs()` | ✅ | |
| `getVerbs(verbs, max)` | `GetVerbs(verbs) int` | ✅ | Returns count |
| `getConicWeights()` | `ConicWeights() []Scalar` | ✅ | Returns slice |
| `getLastPoint()` | `GetLastPoint() (Point, bool)` | ✅ | Returns tuple |

#### Setters
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `setFillType(fillType)` | `SetFillType(fillType)` | ✅ | |
| `toggleInverseFillType()` | `ToggleInverseFillType()` | ✅ | |
| `reset()` | `Reset()` | ✅ | |

#### Path Construction
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `moveTo(x, y)` | `MoveTo(x, y)` | ✅ | |
| `moveTo(pt)` | `MoveToPoint(p)` | ✅ | |
| `lineTo(x, y)` | `LineTo(x, y)` | ✅ | |
| `lineTo(pt)` | `LineToPoint(p)` | ✅ | |
| `quadTo(cx, cy, x, y)` | `QuadTo(cx, cy, x, y)` | ✅ | |
| `quadTo(c, p)` | `QuadToPoint(c, p)` | ✅ | |
| `conicTo(cx, cy, x, y, w)` | `ConicTo(cx, cy, x, y, w)` | ✅ | |
| `conicTo(c, p, w)` | `ConicToPoint(c, p, w)` | ✅ | |
| `cubicTo(cx1, cy1, cx2, cy2, x, y)` | `CubicTo(cx1, cy1, cx2, cy2, x, y)` | ✅ | |
| `cubicTo(c1, c2, p)` | `CubicToPoint(c1, c2, p)` | ✅ | |
| `close()` | `Close()` | ✅ | |
| `addRect(...)` | `AddRect(...)` | ✅ | |
| `addOval(...)` | `AddOval(...)` | ✅ | |
| `addCircle(...)` | `AddCircle(...)` | ✅ | |
| `addRRect(...)` | `AddRRect(...)` | ✅ | |
| `addPath(path, dx, dy, addMode)` | `AddPath(path, dx, dy, addMode)` | ✅ | |
| `addPath(path, addMode)` | `AddPathNoOffset(path, addMode)` | ⚠️ | Different name |
| `addPath(path, matrix, addMode)` | `AddPathMatrix(path, matrix, addMode)` | ⚠️ | Different name |

#### Path Manipulation
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `transform(matrix)` | `Transform(matrix)` | ✅ | |
| `offset(dx, dy)` | `Offset(dx, dy)` | ✅ | |

#### Advanced Methods (Missing)
| C++ Method | Go Method | Status | Notes |
|-----------|-----------|--------|-------|
| `makeFillType(newFillType)` | ❌ | ❌ | Missing - returns new path |
| `makeToggleInverseFillType()` | ❌ | ❌ | Missing - returns new path |
| `snapshot()` | ❌ | ❌ | Missing - returns copy |
| `isInterpolatable(other)` | ❌ | ❌ | Missing - advanced feature |
| `makeInterpolate(other, weight)` | ❌ | ❌ | Missing - advanced feature |
| `interpolate(other, weight, out)` | ❌ | ❌ | Missing - advanced feature |
| `updateBoundsCache()` | `UpdateBoundsCache()` | ✅ | |

---

## Summary of Missing APIs

### SkMatrix Missing Methods
1. **Static Factories:**
   - `NewMatrixRotateRad(rad Scalar)` - radians version
   - `NewMatrixRotateWithPivot(deg Scalar, px, py Scalar)` - static with pivot
   - `NewMatrixAll(...)` - all 9 values
   - `NewMatrixScaleTranslate(sx, sy, tx, ty Scalar)`

2. **Getters:**
   - `Get(index int) Scalar` - indexed access
   - `GetRC(row, col int) Scalar` - row/column access
   - `Get9() [9]Scalar` - get all values
   - `IsTranslate() bool` - translate-only check

3. **Setters:**
   - `Set(index int, value Scalar)` - indexed set
   - `SetAll(...)` - set all 9 values
   - `Set9([9]Scalar)` - set from array
   - Individual setters: `SetScaleX`, `SetScaleY`, `SetSkewX`, `SetSkewY`, `SetTranslateX`, `SetTranslateY`, `SetPerspX`, `SetPerspY`

4. **Convenience:**
   - `MapXY(x, y Scalar) (Scalar, Scalar)` - map single x,y pair

5. **Equality:**
   - `Equals(other SkMatrix) bool` - comparison method

### SkPaint Missing Methods
- **None identified** - API appears complete ✅

### SkPath Missing Methods
1. **Static Factories:**
   - `NewPathRect(...)` - static rectangle factory
   - `NewPathOval(...)` - static oval factory
   - `NewPathCircle(...)` - static circle factory
   - `NewPathRRect(...)` - static rounded rect factory
   - `NewPathPolygon(...)` - static polygon factory
   - `NewPathLine(a, b Point)` - static line factory

2. **Advanced:**
   - `MakeFillType(newFillType) SkPath` - returns new path
   - `MakeToggleInverseFillType() SkPath` - returns new path
   - `Snapshot() SkPath` - returns copy
   - `IsInterpolatable(other SkPath) bool`
   - `MakeInterpolate(other SkPath, weight Scalar) SkPath`

---

## Recommendations

### High Priority (Core API Parity)
1. **Add missing Matrix getters/setters:**
   - `Get(index int) Scalar`
   - `Get9() [9]Scalar`
   - `Set9([9]Scalar)`
   - `SetAll(...)`
   - Individual setters for each matrix element

2. **Add Matrix convenience methods:**
   - `IsTranslate() bool`
   - `MapXY(x, y Scalar) (Scalar, Scalar)`

3. **Add Path static factories:**
   - `NewPathRect(...)`
   - `NewPathOval(...)`
   - `NewPathCircle(...)`
   - `NewPathRRect(...)`

### Medium Priority (Developer Experience)
1. **Add Matrix equality:**
   - `Equals(other SkMatrix) bool`

2. **Add Path convenience methods:**
   - `Snapshot() SkPath`
   - `MakeFillType(newFillType) SkPath`

### Low Priority (Advanced Features)
1. Path interpolation methods
2. Matrix advanced queries (`isSimilarity`, etc.)
3. Matrix homogeneous point mapping

---

## Go Conventions Applied

### Naming Conventions
- **Getters:** `GetXxx()` instead of `getXxx()` (Go convention)
- **Setters:** `SetXxx()` instead of `setXxx()` (Go convention)
- **Queries:** `IsXxx()` instead of `isXxx()` (Go convention)
- **Static factories:** `NewXxx()` instead of constructors (Go convention)

### Return Values
- **Multiple returns:** Go uses tuples `(value, bool)` instead of output parameters
- **Error handling:** `Invert()` returns `(SkMatrix, bool)` instead of `bool invert(SkMatrix* out)`

### Method Overloading
- **Go doesn't support overloading:** Use different method names
  - `MoveTo(x, y)` and `MoveToPoint(p)`
  - `LineTo(x, y)` and `LineToPoint(p)`
  - `QuadTo(...)` and `QuadToPoint(...)`

### Operators
- **No operator overloading:** Use explicit methods
  - `SetConcat()` instead of `operator*`
  - `Equals()` instead of `operator==`

---

## Next Steps

1. **Create API Gap Analysis:** Detailed comparison document
2. **Implement Missing Methods:** Prioritize high-priority items
3. **Update Documentation:** Ensure examples match Skia C++ patterns
4. **Create Migration Guide:** Help C++ developers transition to Go API
5. **Add API Tests:** Verify new methods match C++ behavior

---

**Document Status:** Draft - Ready for Review  
**Next Review:** After implementing missing high-priority APIs

