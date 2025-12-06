# Features

cpx brings a Cargo-like workflow to C++ with opinionated defaults, fast scaffolding, and first-class tooling.

## Project creation
- Interactive TUI (`cpx new`) for project type, test framework, formatting style, hooks, C++ standard, and vcpkg.
- Generates CMake presets, vcpkg manifest, `.clang-format`, tests, git hooks, and optional CI targets.

## Dependency management
- vcpkg-first: everything goes through `vcpkg.json`.
- Direct passthrough of vcpkg commands (`cpx install`, `cpx upgrade`, `cpx search`, etc.).

## Build & test
- CMake + presets with sensible build types.
- Sanitizers: ASan, TSan, MSan, UBSan.
- GoogleTest and Catch2 templates out of the box.

## Code quality
- clang-format, clang-tidy, Cppcheck, Flawfinder built in.
- Optional git hooks (pre-commit/pre-push) running fmt, lint, tests, and security checks.

## CI & cross-compilation
- Docker targets for linux-amd64/arm64, windows-amd64, macOS placeholders, plus Alpine musl images.
- `cpx ci` builds all configured targets; `cpx ci init` can scaffold workflows.

## Releases & self-update
- Release binaries embed the tag version (`cpx --version` matches the downloaded release).
- `cpx upgrade` fetches and replaces the CLI from the latest GitHub release.

