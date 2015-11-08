Krapula-Engine-2
================

2-day Krapula Engine


===== Dependency compiling instructions =====

=== go-gl/glfw in Windoze ===

1. Install msys2 64bit
    https://msys2.github.io/
2. Download https://github.com/Alexpux/MINGW-packages to your msys2 home directory (or wherever you want)

    $ cd ~/MINGW-packages-master\mingw-w64-glfw
    $ makepkg-mingw
  You will be told that a lot of stuff is missing, install it using pacman or whatever it tells you to do.
    $ pacman -U mingw-w64-x86_64-glfw-*.pkg.tar.xz
3. Rename the libglfw.dll.a library file in C:\msys64\mingw64\lib to libglfwdll.a