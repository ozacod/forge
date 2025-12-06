# Code Quality Commands

## cpx fmt
Format with clang-format.
```bash
cpx fmt
cpx fmt --check   # report-only
```

## cpx lint
Run clang-tidy.
```bash
cpx lint
cpx lint --fix    # apply safe fixes
```

## cpx flawfinder
Security scanner for C/C++.
```bash
cpx flawfinder
cpx flawfinder --html
cpx flawfinder --csv
cpx flawfinder --dataflow
```

## cpx cppcheck
Static analysis via Cppcheck.
```bash
cpx cppcheck
cpx cppcheck --xml
cpx cppcheck --enable style,performance
```

## Hooks
Install git hooks to enforce fmt/lint/tests before commits or pushes.
```bash
cpx hooks install
```
Select the checks during `cpx new`; defaults are fmt+lint on pre-commit and tests on pre-push.
