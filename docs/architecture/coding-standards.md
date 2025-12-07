# Coding Standards
## Go-Skia-Support Project

**Last Updated:** 2025-01-27

---

## Language & Style

### Go Conventions
- Follow standard Go formatting (`gofmt`)
- Use `golint` and `go vet` for code quality
- Follow Go naming conventions:
  - Exported: `PascalCase`
  - Unexported: `camelCase`
  - Constants: `PascalCase` or `UPPER_CASE` for exported

### Code Organization
- **Interfaces:** `skia/interfaces/` - Public API contracts
- **Implementations:** `skia/impl/` - Concrete implementations
- **Models:** `skia/models/` - Data structures
- **Enums:** `skia/enums/` - Enumerations
- **Helpers:** `skia/helpers/` - Utility functions
- **Base:** `skia/base/` - Core types and constants

### File Size Guidelines
- **Keep source files small:** Files will grow quickly with reference documentation (C++ source comments, test references, etc.)
- **Target size:** Aim for < 500 lines per file when possible
- **Split large files:** When files exceed ~1000 lines, consider splitting:
  - Separate helper functions into `*_helpers.go` files
  - Extract complex algorithms into dedicated files
  - Move test-related code to test files
- **Documentation overhead:** Account for extensive documentation (C++ references, algorithm explanations, edge case notes)
- **Refactor proactively:** Don't wait until files become unmanageable

---

## Skia C++ Parity Requirements

### Critical Rules
1. **API Parity:** Go interfaces must match Skia C++ API semantics
2. **Calculation Accuracy:** Mathematical operations must match C++ implementation exactly
3. **Edge Cases:** Handle all edge cases from C++ tests (see `docs/cpp-test-reference.md`)
4. **Naming:** Match C++ method names, adapted to Go conventions:
   - C++ `getXxx()` → Go `GetXxx()`
   - C++ `setXxx()` → Go `SetXxx()`
   - C++ `isXxx()` → Go `IsXxx()`

### Floating-Point Precision
- Use `float32` (Go `Scalar`) to match C++ `float`
- Document acceptable tolerance for comparisons
- Test with Skia's test vectors

---

## Testing Standards

### Test Organization
- Test files: `*_test.go` alongside implementation
- Test package: `package impl_test` or `package impl` (as appropriate)
- Test naming: `TestFunctionName` or `TestMethodName_Scenario`

### Test Requirements
1. **Port Skia Tests:** All tests from `docs/cpp-test-reference.md` must be ported
2. **Edge Cases:** Test all edge cases from C++ tests
3. **Precision Tests:** Verify floating-point accuracy
4. **Property Tests:** Test mathematical properties (associativity, identity, etc.)
5. **Dependent Component Coverage:** Tests must extend to all dependent components and helpers:
   - Test helper functions (`skia/helpers/`) independently
   - Test model operations (`skia/models/`) used by implementations
   - Test enum behavior (`skia/enums/`) where applicable
   - Test base constants and utilities (`skia/base/`)
   - Verify integration between components through implementation tests

### Test Coverage Scope
- **Direct Tests:** Test public API methods directly
- **Helper Tests:** Test helper functions (`matrix_helpers.go`, `paint_helpers.go`, `path_helper.go`) independently
- **Model Tests:** Test model operations (Point, Rect, RRect, Color4f) used by implementations
- **Integration Tests:** Verify components work together correctly
- **Edge Case Tests:** Test edge cases in helpers and models, not just main implementations

### Test Helpers
- Use tolerance-based comparison for floating-point
- Create test utilities matching Skia's test helpers
- Reference: `docs/cpp-test-reference.md` for helper functions

---

## Documentation Standards

### Code Comments
- **Public APIs:** Full godoc comments explaining behavior
- **Complex Algorithms:** Inline comments explaining logic
- **C++ References:** Comment with C++ source file/line when porting

### Example Format
```go
// MapPoint transforms a point using the matrix.
// For affine matrices: x' = x*scaleX + y*skewX + transX, y' = x*skewY + y*scaleY + transY
// For perspective matrices: applies perspective division
//
// Ported from: skia-source/src/core/SkMatrix.cpp:mapPoints()
func (m Matrix) MapPoint(pt Point) Point {
    // Implementation
}
```

---

## Error Handling

### Panic vs Error
- **Panic:** Only for programming errors (nil pointer, index out of bounds)
- **Error:** Use return values for recoverable errors
- **Bool Returns:** Use `(value, bool)` pattern for optional operations (e.g., `Invert()`)

### Edge Cases
- **NaN/Inf:** Handle gracefully, return appropriate values
- **Zero Division:** Check before division, return appropriate defaults
- **Empty Inputs:** Handle empty slices/paths gracefully

---

## Performance Considerations

### Optimization
- Match C++ optimization paths (identity checks, scale+translate shortcuts)
- Profile critical paths
- Document performance characteristics

### Memory
- Avoid unnecessary allocations
- Reuse buffers where possible
- Document memory characteristics

---

## Code Review Checklist

- [ ] API matches C++ Skia API
- [ ] Tests ported from C++ tests
- [ ] Tests cover dependent components and helpers
- [ ] Edge cases handled
- [ ] Floating-point precision verified
- [ ] Documentation complete
- [ ] No linting errors
- [ ] Code compiles successfully
- [ ] Test coverage adequate (including helpers and models)
- [ ] File sizes reasonable (< 1500 lines, consider splitting if larger)

---

## References

- **Skia C++ Source:** `../skia-source`
- **Test Reference:** `docs/cpp-test-reference.md`
- **API Parity:** `docs/api-parity-verification.md`
- **Verification Plan:** `docs/functional-parity-verification-plan.md`
- **File Size Guidelines:** `docs/architecture/file-size-guidelines.md`

