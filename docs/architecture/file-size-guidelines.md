# File Size Guidelines
## Managing Source File Growth

**Last Updated:** 2025-01-27

---

## Why File Size Matters

Source files in this project will grow quickly due to:
- **Extensive documentation:** C++ source references, algorithm explanations, edge case notes
- **Inline comments:** Detailed explanations of ported algorithms
- **Test references:** Comments linking to C++ test cases
- **Edge case handling:** Additional code for robustness

**Goal:** Keep files maintainable and navigable.

---

## Size Guidelines

### Target Sizes
- **Optimal:** < 500 lines per file
- **Acceptable:** < 1000 lines per file
- **Action Required:** > 1000 lines - consider splitting

### Current File Sizes (Reference)
- `impl/matrix.go` - ~725 lines ✅
- `impl/path.go` - ~927 lines ✅
- `impl/paint.go` - Check size

---

## When to Split Files

### Indicators for Splitting
1. **File exceeds 1000 lines**
2. **Multiple distinct responsibilities** (e.g., helpers mixed with main implementation)
3. **Difficult to navigate** - hard to find specific functions
4. **Test file growing large** - consider splitting test files too

### Splitting Strategies

#### 1. Extract Helper Functions
**Pattern:** Move helper functions to `*_helpers.go` files

**Example:**
- `matrix.go` - Main implementation
- `matrix_helpers.go` - Helper functions (`rowcol3`, `muladdmul`, etc.)
- `matrix_test.go` - Tests

#### 2. Extract Complex Algorithms
**Pattern:** Move complex algorithms to dedicated files

**Example:**
- `path.go` - Main path implementation
- `path_bounds.go` - Bounds calculation logic
- `path_convexity.go` - Convexity detection logic

#### 3. Extract Model-Specific Code
**Pattern:** Group model-related operations

**Example:**
- `paint.go` - Main paint implementation
- `paint_bounds.go` - Bounds computation
- `paint_models.go` - Internal data structures (already done)

#### 4. Split Test Files
**Pattern:** Organize tests by feature area

**Example:**
- `matrix_test.go` - Core matrix tests
- `matrix_inversion_test.go` - Inversion-specific tests
- `matrix_transformation_test.go` - Transformation tests

---

## Refactoring Guidelines

### Proactive Refactoring
- **Don't wait** until files become unmanageable
- **Refactor early** when approaching 1200-1300 lines
- **Plan splits** before implementing new features

### Refactoring Process
1. **Identify logical boundaries** - What can be separated?
2. **Extract to new file** - Move related code together
3. **Update imports** - Ensure all dependencies resolved
4. **Update tests** - Move related tests if needed
5. **Verify compilation** - Ensure everything still works
6. **Update documentation** - Reflect new structure

---

## Documentation Overhead

### Expected Documentation Per File
- **Function comments:** ~5-10 lines per public function
- **Algorithm explanations:** ~10-20 lines for complex algorithms
- **C++ references:** ~2-3 lines per ported function
- **Edge case notes:** ~5-10 lines per edge case handler

### Example Documentation Overhead
For a file with:
- 20 public functions × 8 lines = 160 lines
- 5 complex algorithms × 15 lines = 75 lines
- 10 edge cases × 7 lines = 70 lines
- **Total overhead:** ~305 lines

**Implication:** A 1000-line file might have ~700 lines of code + ~300 lines of documentation.

---

## File Organization Best Practices

### Keep Related Code Together
- **Main implementation** in primary file
- **Helpers** in `*_helpers.go` files
- **Models** in `*_models.go` files (if internal)
- **Tests** in `*_test.go` files

### Naming Conventions
- `{component}.go` - Main implementation
- `{component}_helpers.go` - Helper functions
- `{component}_models.go` - Internal data structures
- `{component}_test.go` - Tests
- `{component}_{feature}_test.go` - Feature-specific tests

---

## Monitoring File Sizes

### Regular Checks
- **Before adding features:** Check current file size
- **After major additions:** Verify size is still manageable
- **During code review:** Flag files approaching limits

### Tools
```bash
# Count lines in Go files
find skia/impl -name "*.go" -exec wc -l {} + | sort -n

# Check specific file
wc -l skia/impl/matrix.go
```

---

## Examples

### Good: Well-Sized Files
```
impl/
├── matrix.go              (725 lines) ✅
├── matrix_helpers.go      (150 lines) ✅
├── matrix_test.go         (400 lines) ✅
└── paint.go               (500 lines) ✅
```

### Needs Refactoring: Large File
```
impl/
└── path.go                (1800 lines) ⚠️ Split needed
```

### After Refactoring
```
impl/
├── path.go                (800 lines) ✅
├── path_bounds.go         (300 lines) ✅
├── path_convexity.go      (250 lines) ✅
├── path_helper.go         (450 lines) ✅
└── path_test.go           (600 lines) ✅
```

---

## Checklist

When working on a file:
- [ ] Check current line count
- [ ] Estimate new code + documentation overhead
- [ ] If approaching 1200 lines, plan refactoring
- [ ] If exceeding 1500 lines, refactor immediately
- [ ] Extract helpers to separate files
- [ ] Split tests if test file grows large
- [ ] Update documentation to reflect structure

---

## References

- **Coding Standards:** `docs/architecture/coding-standards.md`
- **Source Tree:** `docs/architecture/source-tree.md`
- **Verification Plan:** `docs/functional-parity-verification-plan.md`

