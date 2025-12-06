# Global Configuration

cpx stores global config at:
- Linux/macOS: `~/.config/cpx/config.yaml`
- Windows: `%APPDATA%/cpx/config.yaml`

Example:
```yaml
vcpkg_root: "/path/to/vcpkg"
```

Use commands to manage it:
```bash
cpx config set-vcpkg-root /path/to/vcpkg
cpx config get-vcpkg-root
```

The installer creates the config file and will update it when vcpkg is installed or detected.
