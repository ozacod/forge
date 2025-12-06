# Adding Dependencies

cpx uses vcpkg manifests.
```bash
cpx add port fmt
cpx add port spdlog
```
This updates `vcpkg.json`. To remove:
```bash
cpx remove fmt
```
List/search/info:
```bash
cpx list
cpx search json
cpx info fmt
```
Use `cpx update` to review manifest entries; run `vcpkg upgrade` if you need to bump installed ports.
