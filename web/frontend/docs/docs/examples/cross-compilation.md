# Cross-Compilation Example

Configure targets in `cpx.ci` and run `cpx ci`:
```yaml
targets:
  - name: linux-amd64
    dockerfile: Dockerfile.linux-amd64
  - name: linux-arm64
    dockerfile: Dockerfile.linux-arm64
  - name: linux-amd64-musl
    dockerfile: Dockerfile.linux-amd64-musl
  - name: linux-arm64-musl
    dockerfile: Dockerfile.linux-arm64-musl
  - name: windows-amd64
    dockerfile: Dockerfile.windows-amd64
```

Build:
```bash
cpx ci
cpx ci --target linux-amd64
```
