Installed the go bindings for glfw3
===================================
go get github.com/go-gl/glfw3

Installing glfw3
================
Download page: http://www.glfw.org/download.html
I chose to install from source.
Because `go-gl/glfw3` documentation (https://github.com/go-gl/glfw3) says I need
to install glfw as dynamically linked libraries, which requires me to compile
glfw3 myself.

So I download this:
    http://sourceforge.net/projects/glfw/files/glfw/3.0.2/glfw-3.0.2.zip/download
into
    ~/sources/glfw3

I unzip.
    > unzip glfw-3.0.2.zip 
    > cd glfw-3.0.2

There is a `README.md` in there.  It vaguely explains how to compile.
It asks me to satisfy dependencies.
I'm using gnu with X11, I should apparently install `xorg-dev` and
`libglu1-mesa-dev`.
    > sudo apt-get install xorg-dev
Apparently I already had it.
    > sudo apt-get install libglu1-mesa-dev
I had that one too.
I note that they say that, despite speaking about MESA, it won't tie my glfw
to MESA's implementation of OpenGL; good, since I want hardware acceleration.

Let's start the compilation:
    > cmake .
I got that wrong, I had to specify the option.  Let's try again.
The readme file isn't clear on how to do that, but `go-gl/glfw3/README.md` asks
me to use `-DBUILD_SHARED_LIBS=on`.
    > cmake -DBUILD_SHARED_LIBS=on .
The output was shorter this time, probably it did not have to redo everything.

I don't know how to actually install.  So I try
    > make install
It does stuff, and complains in the end that permissions are denied.  Good.
    > sudo !!
No complaining this time.  Apparently it installed
    /usr/local/lib/libglfw.so.3.0
    /usr/local/lib/libglfw.so.3
    /usr/local/lib/libglfw.so
and more stuff.

I want to see if I can use glfw3 from go now.
I use a tiny bit of code:
```
// game project main.go
package main

    import (
    "fmt"
	glfw "github.com/go-gl/glfw3"
) 

func main() {
	fmt.Println("Hello World!")
	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()
}
```
I build and run from LiteIde with `Ctrl-R` and get this error message:

    /home/delforge/go/src/github.com/niriel/daggor/game/game: error while
    loading shared libraries: libglfw.so.3: cannot open shared object file: No
    such file or directory
    Error: process exited with code 127.

Maybe I need to restart LiteIde?
Nope, that doesn't work.
I delete my executable and run `go build` on the command line.  It creates the
executable again, which I try to run but it fails in the same way.  I may miss
a `LD_LIBRARY_PATH` or something.
So I can build, but not run.
    > export LD_LIBRARY_PATH="/usr/local/lib/"
And I run the game again.  This time it runs.
Or I could run `ldconfig`.
    > sudo ldconfig
Looks like it works now.  I do not have to add my library path all the time.
Yay, we have installed `glfw3`.

