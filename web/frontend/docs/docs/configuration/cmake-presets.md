# CMake Presets

cpx generates `CMakePresets.json` so IDEs can pick up configurations automatically.

- Debug and Release presets are created by default.
- Presets honor sanitizer flags when you build with `--asan`, `--tsan`, `--msan`, or `--ubsan`.
- Keep the file under version control; avoid per-user edits.

If you need custom presets, extend the generated file instead of replacing it so cpx updates remain compatible.
