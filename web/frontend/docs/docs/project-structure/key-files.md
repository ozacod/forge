# Key Files

- `CMakeLists.txt` — project definition and version (updated by `cpx release`).
- `CMakePresets.json` — generated presets for IDEs and CLI builds.
- `vcpkg.json` — dependency manifest; edited by `cpx add/remove`.
- `cpx.ci` — Docker/CI targets for `cpx ci`.
- `include/<name>/version.hpp` — generated version macros; kept in sync with `CMakeLists.txt`.
- `.clang-format` — code style (generated from your TUI choice).
- `tests/` — test scaffolding (GoogleTest or Catch2).
