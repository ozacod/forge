# Cross-Compilation

cpx uses Docker targets defined in `cpx.ci` to build for multiple platforms.

## Sample `cpx.ci`
```yaml
targets:
  - name: linux-amd64
    dockerfile: Dockerfile.linux-amd64
    image: cpx-linux-amd64
    triplet: x64-linux
    platform: linux/amd64

  - name: linux-arm64
    dockerfile: Dockerfile.linux-arm64
    image: cpx-linux-arm64
    triplet: arm64-linux
    platform: linux/arm64

  - name: linux-amd64-musl
    dockerfile: Dockerfile.linux-amd64-musl
    image: cpx-linux-amd64-musl
    triplet: x64-linux
    platform: linux/amd64

  - name: linux-arm64-musl
    dockerfile: Dockerfile.linux-arm64-musl
    image: cpx-linux-arm64-musl
    triplet: arm64-linux
    platform: linux/arm64

  - name: windows-amd64
    dockerfile: Dockerfile.windows-amd64
    image: cpx-windows-amd64
    triplet: x64-windows
    platform: linux/amd64
```
macOS Dockerfiles remain placeholders until an osxcross toolchain is provided.

## Build
```bash
cpx ci                 # all targets
cpx ci --target linux-amd64
```

## Rebuild images
```bash
cpx ci --rebuild
```
Use sparingly—only when Dockerfiles change.
