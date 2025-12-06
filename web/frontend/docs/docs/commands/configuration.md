# Configuration Commands

## Global vcpkg root
```bash
cpx config set-vcpkg-root /path/to/vcpkg
cpx config get-vcpkg-root
```

## Hooks
```bash
cpx hooks install   # install configured hooks
```

## Doc
```bash
cpx doc             # generate docs where supported
```

## Versioning
```bash
cpx release patch   # or minor / major
```
Updates `CMakeLists.txt` (and `include/<name>/version.hpp` if present).

## Self-update
```bash
cpx upgrade
```
Downloads the latest GitHub release and replaces the current binary (release version is embedded via ldflags).
