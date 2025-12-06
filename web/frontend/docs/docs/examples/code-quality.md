# Code Quality Example

Run formatting, linting, and security checks:
```bash
cpx fmt
cpx lint
cpx flawfinder
cpx cppcheck --enable style,performance
```
Install hooks so fmt/lint/tests run before pushes:
```bash
cpx hooks install
```
