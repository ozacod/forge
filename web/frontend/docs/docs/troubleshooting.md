# Troubleshooting

## `cpx --version` shows an old version
- Ensure you installed from the latest release.
- Run `cpx upgrade` to replace the binary; releases embed the tag version.

## Dockerfiles missing after install
- Re-run the installer: `curl -fsSL https://raw.githubusercontent.com/ozacod/cpx/master/install.sh | sh`.
- Dockerfiles live in `~/.config/cpx/dockerfiles`.

## vcpkg not found
- Set the root: `cpx config set-vcpkg-root /path/to/vcpkg`.
- If you don’t have vcpkg, rerun the installer or install vcpkg manually.

## CI cross-compilation fails
- Make sure Docker is available in CI.
- Use `cpx ci --rebuild` after changing Dockerfiles.
- Check that `cpx.ci` targets reference the right Dockerfile names (including musl if needed).

## Hooks aren’t running
- Install them: `cpx hooks install`.
- Verify the repo has `.git` and hooks are executable.
