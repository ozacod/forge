---
slug: /
---

# Quick Start

Get going in minutes.

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/ozacod/cpx/master/install.sh | sh
```

## Create your first project

```bash
# Launch the interactive creator
cpx new
```

Pick project type, test framework, hooks, formatting style, C++ standard, and vcpkg.

After the TUI finishes, enter the project folder and build.

## Build and run

```bash
cd <project_name>
cpx build
cpx run
cpx test
cpx fmt
```

## Upgrade later
```bash
cpx upgrade
```
Replaces your binary with the latest tagged release (version is embedded at build time).

