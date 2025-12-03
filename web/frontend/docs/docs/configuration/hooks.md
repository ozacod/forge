# Git Hooks Configuration

Configure git hooks in `cpx.yaml` to automatically run code quality checks.

## Configuration

Add hooks configuration to your `cpx.yaml`:

```yaml
hooks:
  precommit:
    - fmt      # Format code before commit
    - lint     # Run linter before commit
  prepush:
    - test     # Run tests before push
```

## Supported Hook Checks

- `fmt` - Format code with clang-format
- `lint` - Run clang-tidy static analysis
- `test` - Run tests (blocking for pre-push)
- `flawfinder` - Run Flawfinder security analysis
- `cppcheck` - Run Cppcheck static analysis
- `check` - Run code check

## Installation

Hooks are automatically installed when creating a project from a template with hooks configured. You can also install them manually:

```bash
cpx hooks install
```

## Behavior

- **Hooks configured in cpx.yaml** → Creates actual hook files (e.g., `pre-commit`)
- **Hooks NOT configured** → Creates `.sample` files (e.g., `pre-commit.sample`)
- **No cpx.yaml** → Uses defaults (fmt, lint for pre-commit; test for pre-push)

