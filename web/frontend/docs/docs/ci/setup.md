# CI Setup

1) Define targets in `cpx.ci` (include musl targets if you need Alpine-compatible binaries).
2) Ensure Docker is available in your CI runner.
3) (Optional) Run `cpx ci init --github-actions` or `--gitlab` to scaffold workflows.
4) Cache vcpkg and build artifacts where possible for faster builds.
5) Use `cpx ci --rebuild` only after changing Dockerfiles.

### GitHub Actions snippet
```yaml
- uses: actions/checkout@v4
- name: Build with cpx
  run: |
    cpx ci --target linux-amd64
```

### Artifacts
Collect built binaries from your configured `cpx.ci` output directory.
