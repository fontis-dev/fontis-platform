---
name: rust-hal
description: Develop the Rust HAL crate in core/hal/ following safety, FFI, and error conventions.
---

## rust-hal

### Conventions

- `#![deny(unsafe_code)]` enforced. Unsafe only in `ffi/` module.
- All public items have doc comments.
- All functions return `Result<T, HalError>`.
- FFI module is minimal: thin C-compatible wrappers around safe Rust functions.
- `#[repr(C)]` for all types crossing FFI boundary.
- FFI validates all inputs (pointer validity, buffer sizes, null checks) before calling safe code.
- No shared mutable state across FFI boundary. State owned by Rust side, accessed through handles.

### Testing

- Unit tests in `#[cfg(test)] mod tests {}` blocks within each source file.
- Integration tests in `tests/` using public HAL API only.
- Hardware-dependent tests marked `#[cfg(feature = "hwtest")]` — excluded from `cargo test`.

### Adding a HAL module

1. Create `core/hal/src/<module>/` with `mod.rs` re-exporting public types.
2. Add the module to `lib.rs` as `pub mod <module>;`.
3. Implement safe Rust API, then add FFI wrappers in `core/hal/src/ffi/`.
4. Write unit tests covering all error paths (target: 90%+ coverage).
