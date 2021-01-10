After installation, when you start the Fext, it tries to parse config file which is located in:
`%CONFIG_PATH%/fext`. If it is unsuccessful, Fext will create a new one and then will be found all python directories.
If there was found only one directory, that directory selected by default.
Otherwise, if more than one, you will need to select it from list of found directories.
If directory not found in the path where Fext were looking for, you will need to enter the path to python manually.

Default search paths:

OS|Path
---|---
Windows|`C:\Users\%USERNAME%\AppData\Local\Programms\Python`
Linux|`/usr/lib`
Darwin (Mac OS)|`/usr/local/lib`

<br>The directory is needed so that Fext can interact with packages, since it does not work with a python directly.

Finally, you can simply and conveniently use the package manager.

Install package: `fext i requests`<br>
Uninstall package: `fext u requests`<br>
List of packages: `fext freeze`

More info you can find on [commands](../cmd/cli) page
