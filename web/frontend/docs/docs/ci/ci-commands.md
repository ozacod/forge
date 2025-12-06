# CI Commands

## cpx ci
Build all targets defined in `cpx.ci` using Docker.
```bash
cpx ci                 # all targets
cpx ci --target linux-amd64   # specific target
cpx ci --rebuild       # force rebuild Docker images
```

## cpx ci init
Scaffold CI workflows.
- `cpx ci init --github-actions`
- `cpx ci init --gitlab`

## Tips
- Keep Dockerfiles pinned to specific toolchain versions.
- Use musl images (`Dockerfile.linux-*-musl`) if you need Alpine-compatible artifacts.
- Run `cpx ci --rebuild` only when Dockerfiles change.
