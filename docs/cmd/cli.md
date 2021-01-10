The Fext does not support all commands that PIP provides for now, but will be in the future.

#### Available commands for now:
* [install](#install) (Install a packages)
* [uninstall](#uninstall) (Uninstall a packages)
* [freeze](#freeze) (List of installed packages)
* [debug](#debug) (Debug information)

## Install
Syntax: `fext i [options] <package(s)>`

Installs selected packages, dependencies and dependencies of dependencies.
By default, it searches the latest package version, but you can specify comparison operator.
You can write them however as you want.
<br>Examples:
```bash
fext i aiohttp<=3
fext i "aiohttp <=3"
fext i "aiohttp (<=3 >=1)"
```
<br>
Also, you can install extra packages as:

```bash
fext i requests[socks,security]
```

Available options:

* `-s`, `--single` - Installs single package without dependencies

Planned options:

* `-t`, `--thread` - Enables multi-threading download
* `-g`, `--global` - Installs package globally (for avoid virtualenv)
* `-S`, `--safe` - Safe installation. This means if the package requires a different version of the installed package as a dependency, a fatal error will be thrown. (E.g package required `yarl <= 1`, but package `%packageName%` installed in system required `yarl >= 2`, there will be a fatal error)

## Uninstall
Syntax: `fext u [options] <package(s)>`

Uninstalls single packages without dependencies.<br>

Available options:

* `-w`, `--with-dependencies` - Uninstalls dependencies also

## Freeze
Syntax: `fext freeze`

Shows list of installed packages

## Debug
Syntax: `fext debug`

Shows debug info.
