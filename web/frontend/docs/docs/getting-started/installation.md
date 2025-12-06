# Installation

Pick the quick script or download a release asset manually.

## Quick install (recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/ozacod/cpx/master/install.sh | sh
```
The script:
- Detects OS/arch and downloads the latest release
- Installs/bootstraps vcpkg if needed
- Writes Dockerfiles (including Alpine musl) to `~/.config/cpx/dockerfiles`
- Ensures `cpx` is on your PATH

## Manual install
1) Download the right asset from [Releases](https://github.com/ozacod/cpx/releases/latest).  
2) Install it:
```bash
chmod +x cpx-<os>-<arch>
mv cpx-<os>-<arch> /usr/local/bin/cpx
```
3) Point cpx to vcpkg:
```bash
cpx config set-vcpkg-root /path/to/vcpkg
```

## Verify
```bash
cpx --version
```
Release binaries embed their tag version, so this should match the downloaded release. Use `cpx upgrade` later to replace the binary from the latest release. The installer can be rerun safely to refresh Dockerfiles.

