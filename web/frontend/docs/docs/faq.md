# FAQ

**How do I install?**
Use the one-liner: `curl -fsSL https://raw.githubusercontent.com/ozacod/cpx/master/install.sh | sh`.

**How do I update cpx?**
Run `cpx upgrade`. It downloads the latest release and replaces the binary; `--version` will match the release tag.

**Where are Dockerfiles stored?**
`~/.config/cpx/dockerfiles` (installer writes them, including musl variants).

**Can I use vcpkg commands directly?**
Yes. `cpx install/upgrade/show` forward to vcpkg while keeping `vcpkg.json` as the source of truth.

**How do I bump my project version?**
`cpx release <major|minor|patch>` updates `CMakeLists.txt` and `include/<name>/version.hpp` if present.

**Does cpx support sanitizers?**
Yes: `--asan`, `--tsan`, `--msan`, `--ubsan` on `cpx build` or `cpx check`.
