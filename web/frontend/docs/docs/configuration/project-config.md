# Project Configuration

Most project settings come from the TUI (`cpx new`) and are stored in generated files:
- `CMakeLists.txt` / `CMakePresets.json`
- `vcpkg.json`
- `.clang-format`
- `cpx.ci` (optional)
- `include/<name>/version.hpp` (generated with version macros)

To adjust after creation:
- Update `vcpkg.json` for dependencies.
- Update `CMakeLists.txt` for project name/version; use `cpx release <type>` to bump safely.
- Edit `.clang-format` for style changes.
- Edit `cpx.ci` to tweak CI targets.
