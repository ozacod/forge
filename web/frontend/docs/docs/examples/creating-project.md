# Creating a Project

Use the interactive TUI:
```bash
cpx new
```
Select project type (exe/lib), test framework (GoogleTest, Catch2, doctest), benchmark framework (Google Benchmark, Catch2, nanobench), hooks, formatting style, and C++ standard.

The generator produces CMake presets, vcpkg manifest, tests, benchmarks, git hooks, and optional CI targets.

Build, run, and test:
```bash
cd <project>
cpx build && cpx run
cpx test
cpx bench    # run benchmarks
```
