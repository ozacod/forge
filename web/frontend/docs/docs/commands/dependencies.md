# Dependency Commands

vcpkg-first: all dependency state lives in `vcpkg.json`.

## cpx add port `<pkg>`
Add a vcpkg port to the manifest.
```bash
cpx add port spdlog
cpx add port fmt
```

## cpx remove `<pkg>`
Remove a dependency.
```bash
cpx remove spdlog
```

## cpx list / search / info
Inspect available or installed ports.
```bash
cpx list
cpx search json
cpx info fmt
```

## vcpkg passthrough
All vcpkg commands work through cpx.
```bash
cpx install <package>
cpx upgrade
cpx show <package>
```

## Update manifest
Use `cpx update` to surface what’s in `vcpkg.json`; use `vcpkg upgrade` to bump ports.
