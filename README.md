# Fext
[![Latest](https://img.shields.io/github/v/release/fextpkg/cli)](https://github.com/fextpkg/cli/releases?latest) [![Gover](https://img.shields.io/github/go-mod/go-version/fextpkg/cli?filename=fext%2Fgo.mod)](https://golang.org/dl/) [![License](https://img.shields.io/github/license/fextpkg/cli)](https://github.com/fextpkg/cli/blob/main/LICENSE)

Fext is a modern, small, fast, Go powered package manager for Python.

## Features
- **Speed.** Fast packages downloading and installation process.
- **Dependencies free.** Fext is designed with minimum use of dependencies.
- **Shortcuts.** Every command has it's own shortcut for fastest user experience.
- **Backward compatibility.** Every function designed to be compatible with PIP.

## Installation
You need to download install script and run it: 
```bash
curl https://cdn.lunte.dev/get-fext.py -o get-fext.py
python get-fext.py
```

## Basic usage
To install a package use:
```
fext i(nstall) <package(s)> # fext i aiohttp
```
To uninstall:
```
fext u(ninstall) <package(s)> # fext u aiohttp
```
More commands can be found in [documentation](https://fext.lunte.dev/commands.html).

## Bugs
If you encounter some bugs or problems, please be free to use [Issues page](https://github.com/fextpkg/cli/issues).
