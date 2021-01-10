### v0.0.5
\+ Implemented possibility to installs extra packages<br>
\+ Added -s option for install command<br>
\+ Added print of handled size when using a commands<br>
\+ Some optimizations

### v0.0.3
\+ Implemented beauty uninstall and flag -w (see documentation)<br>
\- Fixed critical errors with a scan dirs, when you can write one letter of package name and uninstall them

### v0.0.2b
\+ Changed chmod from unsafe (777) to safe (755)

### v0.0.2
\+ Support for environment markers. Implemented two base markers: `python_version` and `sys_platform`<br>
\- Fixed mode perm, when you installed package and could not use it (especially on linux)<br>
\- Fixed error when downloaded archive with package was not removed because it was already used (on windows when installing new package)
