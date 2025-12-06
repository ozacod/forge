# Creating a Project

Use the interactive TUI:
```bash
cpx new
```
Select project type (exe/lib), test framework, hooks, formatting style, C++ standard, and vcpkg. The generator produces CMake presets, vcpkg manifest, tests, git hooks, and optional CI targets.

Then build and run:
```bash
cd <project>
cpx build
cpx run
```
