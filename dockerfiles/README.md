# Docker Cross-Compilation Images

This directory contains Dockerfiles for cross-compiling C++ projects to different platforms, similar to Go's cross-compilation experience.

## Available Dockerfiles

- **Dockerfile.linux-amd64** - Linux x86_64 compilation
- **Dockerfile.linux-arm64** - Linux ARM64 compilation (cross-compilation from x86_64)
- **Dockerfile.windows-amd64** - Windows x86_64 compilation (using MinGW-w64)
- **Dockerfile.macos-amd64** - macOS x86_64 compilation (placeholder - requires osxcross setup)
- **Dockerfile.macos-arm64** - macOS ARM64 (Apple Silicon) compilation (placeholder - requires osxcross setup)

## Installation

These Dockerfiles are automatically downloaded to `~/.config/forge/dockerfiles/` (or `%APPDATA%/forge/dockerfiles/` on Windows) during `forge` installation.

## Usage

These Dockerfiles are intended to be used by `forge` commands for cross-compilation. They include:

- Build tools (CMake, Ninja, GCC/Clang)
- vcpkg installation and bootstrapping
- Cross-compilation toolchains (where applicable)
- Proper environment variables for cross-compilation

**Note**: vcpkg is installed in each Docker image at `/opt/vcpkg`. This ensures:
- Self-contained images that work in CI/CD environments
- No dependency on host vcpkg installation
- Reproducible builds across different machines
- Each target can have its own vcpkg packages installed

## Building Images

### Using forge ci command

The easiest way to rebuild Docker images is using the `--rebuild` flag:

```bash
# Rebuild all images and build all targets
forge ci --rebuild

# Rebuild image for specific target only
forge ci --target linux-amd64 --rebuild
```

### Manual Docker build

To build a Docker image manually from the config directory:

```bash
# Get the dockerfiles directory path
dockerfiles_dir=~/.config/forge/dockerfiles

# Build a specific image
docker build -f $dockerfiles_dir/Dockerfile.linux-amd64 -t forge-linux-amd64 $dockerfiles_dir

# Or build all images
for dockerfile in $dockerfiles_dir/Dockerfile.*; do
    image_name="forge-$(basename $dockerfile | sed 's/Dockerfile\.//')"
    docker build -f $dockerfile -t $image_name $dockerfiles_dir
done
```

### When to rebuild

You should rebuild Docker images when:
- Dockerfiles are updated (after running `forge upgrade`)
- You want to ensure you have the latest vcpkg and build tools
- Images are corrupted or outdated
- You've modified the Dockerfiles manually

## Notes

- macOS cross-compilation requires osxcross and macOS SDK, which is complex to set up. These are placeholders for future implementation.
- Windows cross-compilation uses MinGW-w64, which provides good compatibility with most C++ libraries.
- Linux ARM64 cross-compilation uses the `aarch64-linux-gnu` toolchain.

