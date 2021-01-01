After installation, when you start the Fext, it tries to parsing config file which is located in: `%CONFIG_PATH%/fext`. If it unsucessfull, Fext will create a new and will be find all python directores. If there was found one directory, that directory selected by default. Otherwise, if more then one, you will be offered a choice. If directory not found in the path where Fext were looking for, you will need to enter path to python manually.

Default search paths:

OS|Path
--|--
Windows|`C:\Users\%USERNAME%\AppData\Local\Programms\Python`
Linux|`/usr/lib`
Darwin (Mac OS)|`/usr/local/lib`

<br>Directory is needed so that Fext can interact with packages, since it does not work with python directly.

And after that all, you can simplicity and conveniently use the package manager.

Install package: `fext i requests`<br>
Uninstall package: `fext u requests`<br>
List of packages: `fext freeze`

More info you can found on [commands](../cmd/cli) page
