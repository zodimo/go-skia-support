# Source Tree Documentation
## Go-Skia-Support Project Structure

**Last Updated:** 2025-01-27

---

## Overview

This document describes the organization of the Go-Skia-Support codebase, mapping it to the Skia C++ source structure for easy reference.

---

## Package Structure

### `skia/base/`
**Purpose:** Core foundational types and constants

**Files:**
- `base.go` - Scalar type definition
- `constants.go` - Matrix indices, segment masks, mathematical constants

**C++ Equivalent:**
- `include/core/SkScalar.h`
- `include/core/SkMatrix.h` (constants section)
- `include/private/base/SkFloatingPoint.h`

**Key Types:**
- `Scalar` - Floating-point type (`float32`)
- Matrix indices: `KMScaleX`, `KMSkewX`, etc.
- Constants: `SkScalarNearlyZero`, `ScalarRoot2Over2`

---

### `skia/interfaces/`
**Purpose:** Public API interfaces matching Skia C++ API

**Files:**
- `matix.go` - Matrix interface (note: typo in filename preserved)
- `paint.go` - Paint interface
- `path.go` - Path interface
- `alias.go` - Type aliases

**C++ Equivalent:**
- `include/core/SkMatrix.h` (public API)
- `include/core/SkPaint.h` (public API)
- `include/core/SkPath.h` (public API)

**Key Interfaces:**
- `SkMatrix` - 3x3 transformation matrix
- `SkPaint` - Paint properties for drawing
- `SkPath` - 2D path geometry
- Filter interfaces: `Shader`, `ColorFilter`, `ImageFilter`, `MaskFilter`, `PathEffect`, `Blender`

---

### `skia/impl/`
**Purpose:** Concrete implementations of interfaces

**Files:**
- `matrix.go` - Matrix implementation (~725 lines)
- `matrix_helpers.go` - Matrix helper functions
- `paint.go` - Paint implementation
- `paint_helpers.go` - Paint helper functions
- `paint_models.go` - Paint internal models
- `path.go` - Path implementation (~927 lines)
- `path_helper.go` - Path helper functions
- `path_models.go` - Path internal models
- `alias.go` - Implementation aliases and constants

**C++ Equivalent:**
- `src/core/SkMatrix.cpp` (~1800 lines)
- `src/core/SkMatrixInvert.cpp`
- `src/core/SkPaint.cpp`
- `src/core/SkPath.cpp` (multiple files)
- `src/core/SkPathBuilder.cpp`
- `src/core/SkPathData.cpp`
- `src/core/SkPathRef.cpp`

**Key Implementations:**
- `Matrix` - 3x3 matrix with affine and perspective support
- `Paint` - Complete paint implementation
- `pathImpl` - Verb-based path implementation

---

### `skia/models/`
**Purpose:** Data structures representing geometric and color primitives

**Files:**
- `point.go` - 2D point (X, Y)
- `rect.go` - Rectangle (Left, Top, Right, Bottom)
- `rrect.go` - Rounded rectangle
- `color.go` - Color4f (RGBA float)
- `alias.go` - Model type aliases

**C++ Equivalent:**
- `include/core/SkPoint.h`
- `include/core/SkRect.h`
- `include/core/SkRRect.h`
- `include/core/SkColor.h`

**Key Types:**
- `Point` - 2D point
- `Rect` - Axis-aligned rectangle
- `RRect` - Rounded rectangle with corner radii
- `Color4f` - RGBA color (unpremultiplied)

---

### `skia/enums/`
**Purpose:** Enumerations for graphics operations

**Files:**
- `enums.go` - Core enumerations
- `blendmode.go` - Blend mode enumeration

**C++ Equivalent:**
- `include/core/SkBlendMode.h`
- `include/core/SkPathTypes.h`
- `include/core/SkPaint.h` (enums section)

**Key Enums:**
- `BlendMode` - Porter-Duff and advanced blend modes
- `PaintStyle` - Fill, Stroke, StrokeAndFill
- `PaintCap` - Stroke cap styles
- `PaintJoin` - Stroke join styles
- `PathFillType` - Fill rules
- `PathVerb` - Path commands
- `PathDirection` - Contour direction
- `PathConvexity` - Convexity classification
- `MatrixType` - Matrix classification flags

---

### `skia/helpers/`
**Purpose:** Utility functions for mathematical operations

**Files:**
- `helpers.go` - Mathematical helpers
- `alias.go` - Helper type aliases

**C++ Equivalent:**
- `include/private/base/SkMath.h`
- `src/core/SkGeometry.h` (some functions)

**Key Functions:**
- `CrossProduct` - 2D cross product
- `DotProduct` - 2D dot product
- `Sign` - Sign function
- `ScalarPin` - Clamp scalar values

---

## File Mapping: C++ → Go

### Matrix
| C++ File | Go File | Notes |
|----------|---------|-------|
| `include/core/SkMatrix.h` | `interfaces/matix.go` | Public API |
| `src/core/SkMatrix.cpp` | `impl/matrix.go` | Main implementation |
| `src/core/SkMatrixInvert.cpp` | `impl/matrix.go` | Inversion logic integrated |
| `src/core/SkMatrixUtils.h` | `impl/matrix_helpers.go` | Helper functions |

### Paint
| C++ File | Go File | Notes |
|----------|---------|-------|
| `include/core/SkPaint.h` | `interfaces/paint.go` | Public API |
| `src/core/SkPaint.cpp` | `impl/paint.go` | Main implementation |
| `src/core/SkPaintPriv.cpp` | `impl/paint_helpers.go` | Helper functions |
| `src/core/SkPaintDefaults.h` | `impl/paint.go` | Default values |

### Path
| C++ File | Go File | Notes |
|----------|---------|-------|
| `include/core/SkPath.h` | `interfaces/path.go` | Public API |
| `src/core/SkPath.cpp` | `impl/path.go` | Main implementation |
| `src/core/SkPathBuilder.cpp` | `impl/path.go` | Builder pattern integrated |
| `src/core/SkPathData.cpp` | `impl/path_models.go` | Data structures |
| `src/core/SkPathRef.cpp` | `impl/path.go` | Reference counting not needed in Go |
| `src/core/SkPath_editing.cpp` | `impl/path.go` | Editing methods |
| `src/core/SkPath_pathdata.cpp` | `impl/path.go` | Path data access |
| `src/core/SkPathRaw.cpp` | `impl/path_helper.go` | Raw path operations |
| `src/core/SkPathRawShapes.cpp` | `impl/path_helper.go` | Shape generation |

---

## Code Organization Principles

### Separation of Concerns
1. **Interfaces** - Public API contracts only
2. **Implementations** - Concrete logic, can use private helpers
3. **Models** - Pure data structures
4. **Enums** - Type-safe enumerations
5. **Helpers** - Reusable utility functions

### Go-Specific Adaptations
- **No inheritance** - Use composition and interfaces
- **No operators** - Use explicit methods (`Equals()` instead of `==`)
- **No overloading** - Use different method names (`MoveTo()` and `MoveToPoint()`)
- **Error handling** - Use return values (`(value, bool)` tuples)
- **Memory management** - Automatic (no manual allocation)

---

## Test Organization

### Test Files (To Be Created)
- `impl/matrix_test.go` - Matrix tests
- `impl/paint_test.go` - Paint tests
- `impl/path_test.go` - Path tests

### Test Structure
- Port tests from `tests/MatrixTest.cpp`
- Port tests from `tests/PaintTest.cpp`
- Port tests from `tests/PathTest.cpp`
- Reference: `docs/cpp-test-reference.md`

---

## Documentation Files

### Project Documentation
- `README.md` - Project overview
- `docs/functional-parity-verification-plan.md` - Verification strategy
- `docs/api-parity-verification.md` - API comparison
- `docs/cpp-test-reference.md` - Test reference guide
- `docs/api-implementation-summary.md` - Implementation status

### Architecture Documentation
- `docs/architecture/coding-standards.md` - Coding standards
- `docs/architecture/tech-stack.md` - Technology stack
- `docs/architecture/source-tree.md` - This file

---

## References

- **Skia C++ Source:** `/home/jaco/SecondBrain/1-Projects/GoProjects/Development/skia-source`
- **Test Reference:** `docs/cpp-test-reference.md`
- **API Parity:** `docs/api-parity-verification.md`

