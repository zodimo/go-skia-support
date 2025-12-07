# Skia Go Package

A Go port of Skia's core drawing primitives, designed to be backend-agnostic, Bring Your Own Graphics Backend.

## Overview

This package provides a complete implementation of Skia's fundamental graphics types and operations in pure Go. It follows Skia's C++ API closely while adapting to Go idioms and conventions.

## Package Structure

### `base/`
Core foundational types and constants:
- **Scalar**: Floating-point type (`float32`) used throughout the package
- **Constants**: Matrix indices, segment masks, and mathematical constants

### `interfaces/`
Interface definitions for all major graphics primitives:
- **SkPaint**: Paint interface for specifying how geometry is drawn (color, stroke, blend modes, filters)
- **SkPath**: Path interface for creating and manipulating 2D paths
- **SkMatrix**: Matrix interface for 2D transformations
- **Filter Interfaces**: Shader, ColorFilter, ImageFilter, MaskFilter, PathEffect, Blender

### `models/`
Data structures representing geometric and color primitives:
- **Point**: 2D point with X, Y coordinates
- **Rect**: Rectangle with Left, Top, Right, Bottom bounds
- **RRect**: Rounded rectangle with corner radii
- **Color4f**: RGBA color with float components (unpremultiplied)

### `enums/`
Enumerations for various graphics operations:
- **BlendMode**: Porter-Duff and advanced blend modes (SrcOver, Multiply, Screen, etc.)
- **PaintStyle**: Fill, Stroke, or StrokeAndFill
- **PaintCap**: Stroke cap styles (Butt, Round, Square)
- **PaintJoin**: Stroke join styles (Miter, Round, Bevel)
- **PathFillType**: Fill rules (Winding, EvenOdd, Inverse variants)
- **PathVerb**: Path commands (Move, Line, Quad, Conic, Cubic, Close)
- **PathDirection**: Contour direction (CW, CCW)
- **PathConvexity**: Path convexity classification
- **MatrixType**: Matrix classification flags

### `impl/`
Concrete implementations of the interfaces:
- **Paint**: Complete paint implementation with:
  - Color management (RGBA, ARGB, alpha)
  - Stroke properties (width, cap, join, miter limit)
  - Blend modes and custom blenders
  - Filter support (shader, color filter, image filter, mask filter, path effect)
  - Fast bounds computation
- **Path**: Complete path implementation with:
  - Verb-based path construction (MoveTo, LineTo, QuadTo, ConicTo, CubicTo, Close)
  - Path manipulation (Transform, Offset, AddPath)
  - Bounds computation (bounds, tight bounds)
  - Convexity detection
  - Fill type management
- **Matrix**: 3x3 transformation matrix with:
  - Affine transformations (translate, scale, rotate, skew)
  - Perspective transformations
  - Point and rect mapping
  - Matrix concatenation (pre/post)

### `helpers/`
Utility functions for mathematical operations:
- **CrossProduct**: 2D cross product
- **DotProduct**: 2D dot product
- **Sign**: Sign function
- **ScalarPin**: Clamp scalar values

## Key Features

- **Backend Agnostic**: Interfaces allow different backends (OpenGL, Vulkan, Metal, CPU) to implement rendering
- **Complete Paint System**: Full support for colors, strokes, blend modes, and all filter types
- **Robust Path Implementation**: Verb-based path system with proper bounds and convexity tracking
- **Matrix Transformations**: Full 3x3 matrix support including perspective
- **Type Safety**: Strong typing throughout with clear separation of concerns

## Usage Example

```go
import (
    "github.com/zodimo/go-skia-support/skia/interfaces"
    "github.com/zodimo/go-skia-support/skia/impl"
    "github.com/zodimo/go-skia-support/skia/models"
    "github.com/zodimo/go-skia-support/skia/enums"
)

// Create a paint
paint := impl.NewPaint()
paint.SetColor(models.Color4f{R: 1.0, G: 0.0, B: 0.0, A: 1.0})
paint.SetStyle(enums.PaintStyleFill)
paint.SetAntiAlias(true)

// Create a path
path := impl.NewSkPath(enums.PathFillTypeWinding)
path.AddCircle(100, 100, 50, enums.PathDirectionCW)

// Create a matrix
matrix := impl.NewMatrixTranslate(10, 20)
```

## Development Status

The implementation aims to follow  Skia's C++ codebase closely and is designed to be a drop-in replacement for Skia's core types in Go applications.
