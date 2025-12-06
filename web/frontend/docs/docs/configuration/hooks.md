# Git Hooks

Choose hook checks during `cpx new` (fmt, lint, test, flawfinder, cppcheck, check). Install them with:
```bash
cpx hooks install
```

Default if you skip selection: fmt + lint on pre-commit, tests on pre-push.

Use hooks to keep formatting, linting, and tests consistent across contributors.
