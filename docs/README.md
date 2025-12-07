# Documentation Index
## Go-Skia-Support Project

This directory contains all project documentation organized by BMAD framework structure.

---

## Project Documentation

### Verification & Planning
- **[Functional Parity Verification Plan](./functional-parity-verification-plan.md)** - Comprehensive strategy for verifying Go port matches C++ Skia
- **[API Parity Verification](./api-parity-verification.md)** - Detailed API comparison (C++ vs Go)
- **[C++ Test Reference Guide](./cpp-test-reference.md)** - Complete reference of all C++ tests to port
- **[API Implementation Summary](./api-implementation-summary.md)** - Status of implemented APIs

---

## Architecture Documentation

See `docs/architecture/` directory:
- **[Coding Standards](./architecture/coding-standards.md)** - Go coding standards and Skia parity requirements
- **[Technology Stack](./architecture/tech-stack.md)** - Technology choices and build configuration
- **[Source Tree](./architecture/source-tree.md)** - Code organization and C++ → Go mapping

---

## BMAD Framework Directories

### `docs/architecture/`
Architecture documentation and technical specifications.

### `docs/stories/`
User stories and development tasks (to be populated during development).

### `docs/qa/`
- `assessments/` - QA assessments
- `gates/` - QA gate documents

### `docs/prd/`
Sharded PRD documents (if PRD is created).

### `docs/source-ref/`
Reference materials from Skia C++ source.

---

## Quick Links

- **Skia C++ Source:** `/home/jaco/SecondBrain/1-Projects/GoProjects/Development/skia-source`
- **BMAD Core:** `.bmad-core/`
- **Test Reference:** `docs/cpp-test-reference.md`
- **API Comparison:** `docs/api-parity-verification.md`

---

## Project Status

**Current Phase:** Functional Parity Verification  
**Focus:** API implementation and test porting  
**Next Steps:** Begin porting high-priority tests from C++ test suite

