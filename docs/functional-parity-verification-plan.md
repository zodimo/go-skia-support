# Functional Parity Verification Plan
## Go-Skia-Support Port Verification Strategy

**Project:** go-skia-support  
**Source Version:** Skia chrome/m144  
**Critical Requirement:** Calculation accuracy is paramount  
**Date:** 2025-01-27

---

## Executive Summary

This document outlines a comprehensive strategy to verify functional parity between the Go port of Skia's core types (`SkMatrix`, `SkPaint`, `SkPath`) and the original C++ implementation from chrome/m144. Given the calculation-heavy nature of this library and the critical importance of accuracy, this plan establishes systematic verification methodologies, test strategies, and comparison frameworks.

---

## Problem Statement

### Current State
- **Ported Components:**
  - `SkMatrix` - 3x3 transformation matrix with affine and perspective support
  - `SkPaint` - Paint object specifying drawing properties (color, stroke, blend modes, filters)
  - `SkPath` - 2D path with verb-based construction and manipulation
  - Supporting types: `Point`, `Rect`, `RRect`, `Color4f`
  - Enums: `BlendMode`, `PaintStyle`, `PaintCap`, `PaintJoin`, `PathFillType`, etc.

- **Source Reference:**
  - C++ sources located in `/home/jaco/SecondBrain/1-Projects/GoProjects/Development/skia-source`
  - Primary implementations:
    - `src/core/SkMatrix.cpp` (~1800 lines)
    - `src/core/SkPaint.cpp`
    - `src/core/SkPath.cpp` (multiple files)
  - Headers: `include/core/SkMatrix.h`, `include/core/SkPaint.h`, `include/core/SkPath.h`

### Challenge
Ensuring bit-for-bit or near-bit-for-bit accuracy in mathematical operations across language boundaries (C++ → Go) where:
- Floating-point precision differences may occur
- Different compiler optimizations may affect intermediate calculations
- Type system differences (C++ `float` vs Go `float32`) need verification
- Edge cases and corner cases must be thoroughly tested

---

## Verification Strategy Overview

### Three-Tier Verification Approach

1. **Unit-Level Verification** - Direct function-by-function comparison
2. **Integration Verification** - Complex operation sequences
3. **Regression Testing** - Continuous validation against Skia test suite

---

## Detailed Verification Methodology

### Phase 0: API/Interface Parity Verification (CRITICAL FOR DEVELOPER EXPERIENCE)

**Objective:** Ensure Go interfaces match Skia C++ API so developers familiar with Skia feel at home.

**Actions:**
- [ ] Complete method-by-method comparison of C++ headers vs Go interfaces
- [ ] Verify all public methods are present (see `docs/api-parity-verification.md`)
- [ ] Verify method signatures match (parameters, return types)
- [ ] Document intentional differences and rationale
- [ ] Verify static factory methods exist
- [ ] Verify operator equivalents exist (Go doesn't support operators)
- [ ] Create migration guide for C++ developers

**Deliverable:** API parity verification report (see `docs/api-parity-verification.md`)

**Status:** ✅ Initial analysis complete - see `docs/api-parity-verification.md` for detailed findings

**Key Findings:**
- SkPaint API: ✅ Complete
- SkMatrix API: ⚠️ Missing ~15 methods (getters, setters, static factories)
- SkPath API: ⚠️ Missing static factories and some advanced methods

### Phase 1: Source Code Analysis & Mapping

#### 1.1 Function Inventory
**Objective:** Create comprehensive mapping of C++ functions to Go implementations

**Actions:**
- [ ] Extract all public methods from C++ headers (`SkMatrix.h`, `SkPaint.h`, `SkPath.h`)
- [ ] Map each C++ method to corresponding Go method
- [ ] Identify any missing functions in Go port
- [ ] Document intentional omissions (if any)

**Deliverable:** Function mapping matrix (C++ → Go)

#### 1.2 Algorithm Comparison
**Objective:** Verify algorithmic equivalence

**Actions:**
- [ ] Compare core algorithms line-by-line:
  - Matrix multiplication (`SetConcat`)
  - Matrix inversion (`Invert`)
  - Point transformation (`MapPoint`, `MapPoints`)
  - Rect transformation (`MapRect`)
  - Path bounds calculation
  - Paint bounds computation
- [ ] Document any algorithmic differences
- [ ] Verify optimization paths match (e.g., identity matrix shortcuts, scale+translate optimizations)

**Deliverable:** Algorithm comparison report

#### 1.3 Constant & Enum Verification
**Objective:** Ensure all constants match exactly

**Actions:**
- [ ] Compare matrix indices (`kMScaleX`, `kMSkewX`, etc.)
- [ ] Verify enum values match C++ definitions
- [ ] Check default values (e.g., `PaintDefaultsMiterLimit = 4.0`)
- [ ] Verify mathematical constants (PI, epsilon values)

**Deliverable:** Constants verification checklist

---

### Phase 2: Test-Driven Verification

#### 2.1 Skia Test Suite Extraction
**Objective:** Leverage Skia's existing test suite

**Actions:**
- [ ] Locate Skia test files:
  - `tests/MatrixTest.cpp`
  - `tests/PaintTest.cpp`
  - `tests/PathTest.cpp`
- [ ] Extract test cases and adapt to Go test framework
- [ ] Create Go test equivalents using `testing` package
- [ ] Run tests and compare results

**Deliverable:** Ported test suite with results comparison

#### 2.2 Edge Case Testing
**Objective:** Verify behavior at boundaries and special cases

**Critical Test Cases:**

**Matrix:**
- [ ] Identity matrix operations
- [ ] Zero scale factors
- [ ] Negative scale factors
- [ ] Very large/small values (overflow/underflow)
- [ ] NaN and Inf handling
- [ ] Perspective division by zero
- [ ] Determinant near zero (singular matrices)
- [ ] Matrix inversion of singular matrices
- [ ] Concatenation with identity matrices
- [ ] Scale+translate optimization paths

**Paint:**
- [ ] Default paint values
- [ ] Zero-width strokes (hairline)
- [ ] Negative stroke width handling
- [ ] Blend mode edge cases
- [ ] Alpha channel edge cases (0, 1, >1)
- [ ] Filter combinations

**Path:**
- [ ] Empty paths
- [ ] Single-point paths
- [ ] Degenerate curves (zero-length)
- [ ] Self-intersecting paths
- [ ] Path bounds with empty/zero-area paths
- [ ] Convexity detection edge cases

**Deliverable:** Edge case test suite with pass/fail results

#### 2.3 Precision Testing
**Objective:** Verify floating-point accuracy

**Actions:**
- [ ] Create test vectors with known expected results
- [ ] Compare Go output vs C++ output for identical inputs
- [ ] Measure deviation (ULP - Units in Last Place)
- [ ] Document acceptable tolerance thresholds
- [ ] Test cumulative error in chained operations

**Test Vectors:**
- Common transformations (90° rotations, 2x scales)
- Arbitrary transformations
- Chained operations (10+ concatenations)
- Inverse operations (M * M^-1 should equal identity)

**Deliverable:** Precision analysis report with tolerance specifications

---

### Phase 3: Comparative Execution Testing

#### 3.1 C++ Bridge Testing
**Objective:** Direct comparison using CGO or shared library

**Approach Options:**

**Option A: CGO Wrapper**
- Create CGO bindings to original Skia C++ code
- Run identical test cases in both Go and C++
- Compare outputs programmatically

**Option B: Shared Library**
- Compile Skia as shared library
- Load from Go using cgo
- Execute side-by-side comparisons

**Option C: Test Harness**
- Create C++ test harness that outputs results to file
- Create Go test harness with same inputs
- Compare output files

**Deliverable:** Comparative test framework

#### 3.2 Property-Based Testing
**Objective:** Verify mathematical properties hold

**Properties to Verify:**

**Matrix Properties:**
- [ ] Associativity: `(A * B) * C == A * (B * C)`
- [ ] Identity: `M * I == I * M == M`
- [ ] Inverse: `M * M^-1 == I` (for invertible matrices)
- [ ] Transpose properties
- [ ] Determinant properties

**Path Properties:**
- [ ] Bounds consistency: `bounds.Contains(path.points)`
- [ ] Transform consistency: `path.Transform(M).bounds == path.bounds.Transform(M)`
- [ ] AddPath consistency

**Deliverable:** Property-based test suite

---

### Phase 4: Performance & Correctness Benchmarking

#### 4.1 Performance Comparison
**Objective:** Ensure Go implementation doesn't introduce performance regressions

**Actions:**
- [ ] Benchmark critical operations:
  - Matrix multiplication
  - Point transformation (single and batch)
  - Path bounds calculation
  - Paint bounds computation
- [ ] Compare against C++ benchmarks (if available)
- [ ] Document performance characteristics

**Deliverable:** Performance benchmark report

#### 4.2 Correctness Under Load
**Objective:** Verify correctness during intensive operations

**Actions:**
- [ ] Run long sequences of operations (1000+ transformations)
- [ ] Verify no accumulation errors beyond acceptable tolerance
- [ ] Test memory stability (no leaks, no corruption)

**Deliverable:** Stress test results

---

## Critical Areas Requiring Special Attention

### 1. Floating-Point Precision
**Risk:** Go `float32` vs C++ `float` may have subtle differences

**Mitigation:**
- Use identical IEEE 754 compliance
- Document acceptable ULP differences
- Test with various input ranges
- Verify rounding modes match

### 2. Matrix Type Classification
**Risk:** The C++ `SkMatrix` uses cached type masks for optimization

**Verification Points:**
- [ ] Type mask computation matches C++ logic
- [ ] Type mask invalidation on mutations matches
- [ ] Optimization paths triggered correctly

### 3. Path Bounds Calculation
**Risk:** Bounds computation is complex and accuracy-critical

**Verification Points:**
- [ ] Exact bounds match C++ implementation
- [ ] Tight bounds calculation matches
- [ ] Bounds invalidation/caching logic matches

### 4. Paint Bounds Computation
**Risk:** Stroke width and effects affect bounds calculation

**Verification Points:**
- [ ] Stroke width correctly expands bounds
- [ ] Filter effects correctly expand bounds
- [ ] Blend mode doesn't affect bounds (correct?)

---

## Verification Tools & Infrastructure

### Recommended Tools

1. **Test Framework:**
   - Go `testing` package
   - Property-based testing: `gopter` or `rapid`
   - Benchmarking: `testing.B`

2. **Comparison Tools:**
   - Custom comparison functions with configurable tolerance
   - Test result diffing utilities
   - Visualization tools for matrix/geometry operations

3. **C++ Integration:**
   - CGO for direct C++ calls
   - Or separate C++ test harness with file I/O

4. **Documentation:**
   - Test coverage reports
   - Verification status tracking
   - Deviation documentation

---

## Success Criteria

### Must Have (Blocking Issues)
- ✅ All Skia test cases pass (or documented acceptable deviations)
- ✅ Edge cases handled identically to C++
- ✅ Floating-point precision within acceptable tolerance (< 1 ULP for most operations)
- ✅ No crashes or panics on valid inputs
- ✅ Correct handling of invalid inputs (NaN, Inf, singular matrices)

### Should Have (Quality)
- ✅ Performance within 2x of C++ implementation
- ✅ 100% code coverage of ported functions
- ✅ Comprehensive documentation of any deviations
- ✅ Property-based tests verify mathematical correctness

### Nice to Have (Future)
- ✅ Automated regression testing against Skia updates
- ✅ Performance profiling and optimization
- ✅ Extended test coverage beyond Skia's test suite

---

## Implementation Plan

### Week 1: Foundation
- [ ] Complete source code analysis (Phase 1)
- [ ] Set up test infrastructure
- [ ] Create function mapping document

### Week 2: Core Testing
- [ ] Port Skia test suite (Phase 2.1)
- [ ] Implement edge case tests (Phase 2.2)
- [ ] Begin precision testing (Phase 2.3)

### Week 3: Comparative Testing
- [ ] Set up C++ bridge/testing harness (Phase 3)
- [ ] Run comparative tests
- [ ] Analyze deviations

### Week 4: Validation & Documentation
- [ ] Complete property-based testing
- [ ] Performance benchmarking
- [ ] Document all findings
- [ ] Create verification report

---

## Risk Mitigation

### Risk: Floating-Point Differences
**Mitigation:** Establish clear tolerance thresholds, document acceptable deviations

### Risk: Missing Test Coverage
**Mitigation:** Extract and port all Skia tests, supplement with additional edge cases

### Risk: Performance Degradation
**Mitigation:** Benchmark early, optimize hot paths if needed

### Risk: Undiscovered Bugs
**Mitigation:** Property-based testing, fuzzing, extensive edge case coverage

---

## Next Steps

1. **Immediate Actions:**
   - [ ] Review this plan and adjust priorities
   - [ ] Set up test infrastructure
   - [ ] Begin Phase 1: Source code analysis
   - [ ] Extract Skia test cases

2. **Questions to Resolve:**
   - What is acceptable floating-point tolerance? (ULP threshold)
   - Should we prioritize performance or exact precision?
   - Are there specific use cases that are most critical?
   - Do we have access to compile/run C++ Skia for comparison?

3. **Resource Needs:**
   - Test infrastructure setup
   - C++ compilation environment (if doing comparative testing)
   - Time allocation for thorough verification

---

## Appendix: Reference Files

### C++ Source Files
- `skia-source/src/core/SkMatrix.cpp`
- `skia-source/include/core/SkMatrix.h`
- `skia-source/src/core/SkMatrixInvert.cpp`
- `skia-source/src/core/SkPaint.cpp`
- `skia-source/include/core/SkPaint.h`
- `skia-source/src/core/SkPath.cpp`
- `skia-source/include/core/SkPath.h`

### Go Implementation Files
- `go-skia-support/skia/impl/matrix.go`
- `go-skia-support/skia/impl/paint.go`
- `go-skia-support/skia/impl/path.go`

### Test Files (to locate/port)
- `skia-source/tests/MatrixTest.cpp`
- `skia-source/tests/PaintTest.cpp`
- `skia-source/tests/PathTest.cpp`

---

**Document Status:** Draft - Ready for Review and Refinement  
**Owner:** Mary (Business Analyst)  
**Next Review:** After initial verification work begins

