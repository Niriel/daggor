Installing the go bindings for glfw3.
=====================================
    
    go get github.com/go-gl/glfw3

The last commit is 408297defc6a24b5dd60fca13504b37d4e35b6bf.



Installing glfw3.
=================

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



Installing the go bindings for OpenGL.
======================================

    go get github.com/go-gl/gl

The last commit is 4b3131d48842f804af76cd82f64f7520677cbece.


Interactivity.
==============

Keyboard input.
---------------
Safety before speed.
Several keys can be pressed or released during a frame.
I do not want to be interrupted all the time during the computations of a frame,
though.
Therefore I want that, for a given frame, there be a fixed list of events to
process.  They will be the events received since the last frame.
I double buffer keyboard events to avoid the list to grow while I am processing
it.

State.
------
The game has a state now.  Currently, the state is just a player position stored
as an `int`.  The state is immutable.  Well, it technically IS mutable, but I
refuse to mutate it.
Keyboard events are translated into commands, and commands affect the state.
A new state is then created.
This is how, by pressing W and S, we can change the player position (displayed
in the console).

Move in a 3D world.
-------------------
I want to link the player position from the state discussed above to an actual
modification of what is displayed.
So I created a 3D cube, and extended its shader to take a MVP matrix as a
uniform.
Then I brought my the `glm` library that I had written in the past and adapted
it to `go-gl`.
By pressing the W and S keys, the red cube moves because the player moves.

The next step is to be able to display several models using different shaders.
I added a green pyramid to the code, next to the red cube.
Now I can display a bunch of these with a loop.
Note that the location of the pyramid and cubes are not owned by these shapes.
These shapes act as functions: draw a cube here, draw a pyramid there.

I implemented the movements along x and y.  Use WASD and QE.  W goes forward,
S goes backward, A strafes left, D strafes right, Q rotates 90 degrees direct
(makes you face what used to be on your left) and E rotates 90 degrees
retrograde (turn right).

To do that, I had to augment the state of the world with two variables, for a
total of three: X and Y for the position of the player, and F for his facing.
I keep them all integers.  A player position can produce a corresponding view
matrix.  At first I got the multiplication of the rotation and the translation
wrong.  We want to apply the translation first, then rotate, not the other way
around.

I also had to decide a reference frame for the facing.  In accordance with
trigonometric conventions, a facing (angle) of 0 makes you look toward +x.
That, plus the fact that I like my z axis pointing up, forces me to depart from
the natural OpenGL reference frame.  It is not a big change though, just a
rotation matrix to apply before doing the projection.  Actually, I even embed
this rotation inside the projection matrix, then I can forget about it.
Note that having z pointing up matches Blender's default frame, which is nice.

The world in memory.
====================
I created a structure that would contain the state of the world.  For now, it
just contains the landscape (name chosen to avoid clashing with the keyword
"map").  That structure, as always, is immutable.  Each frame, that landscape
is analyzed and draw calls are performed.  Visually, it looks like before.  But
now, I should be able to change the landscape while playing.

To prove it, I added three commands linked to key presses:

* `C` creates a cube,
* `P` creates a pyramid,
* `Delete` removes the cube or pyramid, if any.

The edition is limited to the 16x16 grid of the landscape, even though the
player can venture beyond.
The tile affected is not the one on which the player stands, but the one in
front.

Note that with my current 80 degrees FOV and the size of the tiles, I cannot see
the walls directly on my left and right.  Although it is correct considering the
geometry, it is poor for the gameplay as you never know if you stand in a
corridor or an open area.  I must come up with a solution, probably wider tiles.

Save and load game.
===================
My immutable structures should make it easy.  Let's get it over with as soon as
possible.  As always with serializing stuff, the pointers are a mess.  Some state
is also contained inside OpenGL, and this one should NOT be saved to disk.

So, the same data may be expressed in two ways: the way that is stored in the
file, which contains no pointer, and maybe the version in memory that has
pointers.  Well, using pointers is a type of optimization, I can do without.
So I should work at removing all pointers from the game state.

This is still going to be a problem with slices, as these contain pointer
internally.  Therefore I should not dump slices to disk directly.  Maps probably
have the same problem.

Something without pointers...  Isn't that a relational database?

Let us remove pointers from the landscape.

