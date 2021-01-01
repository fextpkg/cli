The Fext currently does not support all commands for now, but it will certainly support in the future. But, base functional of course Fext provides.

#### Available commands for now:
* [install](#install) (Install a packages)
* [uninstall](#uninstall) (Uninstall a packages)
* [freeze](#freeze) (List of installed packages)
* [debug](#debug) (Debug information)

## Install
Syntax: `fext i [options] <package(s)>`

Installs selected packages, dependencies and dependencies of dependencies. By default it search latest package version, but you can specify comparison operator. You can write them however as you want.
<br>Examples:
```bash
fext i aiohttp<=3
fext i "aiohttp <=3"
fext i "aiohttp (<=3 >=1)"
```
<br>

Available options:

- missing here

Planned options:

* `-s`, `--single` - Install single package without dependencies
* `-t`, `--thread` - Enable multi-threading download
* `-g`, `--global` - Install package globally (for avoid virtualenv)
* `-S`, `--safe` - Safe install. This means if the package requires a different version of the installed package as a dependency, a fatal error will be thrown. (E.g package required `yarl <= 1`, but package `%packageName%` installed in system required `yarl >= 2`, there will be a fatal error)

## Uninstall
Syntax: `fext u [options] <package(s)>`

Uninstall single packages without dependencies.<br>

Available options:

* `-w`, `--with-dependencies` - Uninstall dependencies also

## Freeze
Syntax: `fext freeze`

Show list of installed packages

## Debug
Syntax: `fext debug`

Show debug info.
