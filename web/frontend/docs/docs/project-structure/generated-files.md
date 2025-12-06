# Generated Files

When you run `cpx new`, you get:

- `CMakeLists.txt` / `CMakePresets.json`
- `vcpkg.json`
- `.clang-format`
- `include/<name>/<name>.hpp`
- `include/<name>/version.hpp`
- `src/` and `tests/` scaffolding (GoogleTest or Catch2)
- `cpx.ci` (optional Docker/CI targets)
- Git hooks configuration (installed later via `cpx hooks install`)

These files are intended to be checked into version control and edited as your project evolves.
