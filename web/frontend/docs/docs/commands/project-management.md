# Project Management Commands

## cpx new
Interactive scaffolding. No flags needed.
```bash
cpx new
```
You pick project type, test framework, hooks, formatting style, C++ standard, and vcpkg usage. Outputs CMake presets, vcpkg manifest, git hooks, tests, and optional CI targets.

## cpx release
Bump semantic version in `CMakeLists.txt` (and `include/<name>/version.hpp` if present).
```bash
cpx release patch   # or: minor / major
```

## cpx doc
Generate docs (where supported) for the current project.

## cpx upgrade
Self-upgrade cpx to the latest GitHub release.

