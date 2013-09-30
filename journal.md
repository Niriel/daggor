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

I replaced pointers by indices into arrays.  It was easy with the very simple
world we have now.  Things may get tricker once we start add/remove creatures
and items, as some arrays will start having holes.  It's like removing rows
from databases.  We'll see in good time.

So now, we can save (press F4) and load (press F5).  The `<nil>` in the console
means there was no error.  Things are likely to explode if you load an old saved
world, and I will not protect myself against that for now.

Also, note that the way things work now, the entire world is held in memory.
This will be bad once we get a big world.

A real world.
=============
Cubes and pyramid suck.  I want to see floors, ceilings and walls.  Not even
counting creatures, chests and other decorations.  The world will be a
collection of such things, attached to positions.  I could create a Tile
structure that contains a ceiling, a floor and four walls, but I don't think it
would actually be useful for anything.  Actually, I think that having a list of
ceilings, a list of walls, etc., will improve performance.  Not that I should
think about that now, but still: I'd rather copy a list of ceiling in case of
modification (immutable, remember?) than copy a list of tiles.

I start with keeping cubes and pyramids but I make them technically floor
elements.  This is because I do not wish to change looks and logic at the same
time.

It is now time to change the appearance of things.  I can keep the cube and
pyramid around to use them as place holders later, but right now I need floors.
A floor should be easy, just a square, two triangles.  One detail though: I want
it to have a facing, I should be able to rotate it.  Since I do not want to
play with textures for now, I will have to put colors in the vertices.  That
requires two new shaders.

Adding more shaders calls for a shader management tools.  I want to come up with
a function that abstracts the whole shader and program process, creating/
compiling them on the fly.

I created a new 3D model.  Just a flat square for the floor.  Its vertices have
different colors so you can see its orientation.  You can place it with `F` and
rotate it with `[` and `]`.  You cannot rotate cubes and pyramids, not just
because it would not show, but because they do not have a facing information.
Indeed, I now have two types of buildings, one with and one without facing.
They both share a common interface so they can be processed together.  Type
assertions allows me to check whether I can rotate.

Lots of refactoring.

Time to get walls.  I hardcoded a wall, just a vertical square.  It is still
used as a floor tile which makes no sense, but I will start activating the walls
and ceilings.

Placing walls was surprisingly easy.  I have a loop over the four possible
facings in order to render the four lists of walls.  The walls that I have
designed are, on purpose, not touching the border of the tiles.  This simulates
their thickness.  Of course, this results in the ability to see between walls
when it should not be possible.  That is OK:  the artist will come up with
different meshes for different walls, adapted to the various configurations of
walls.

I am also thinking of having columns.  Columns live on the corners of tiles.
They are not walls, they are they own type of building.  They never block
movement, they are decorative.  One possible use of columns is to hide the
seams between walls.  The columns live in a coordinate space that is shifted
relatively to the center of the tiles.  They are shifted toward -x -y.
Therefore, the column that has coordinates (4, 6) in its list corresponds to the
south-west corner of the tile (4, 6), so will be centered on (3.5, 5.5).
Columns are orientable buildings.  By the way, if the artist wishes to place a
column in a middle of a tile for some reason, he can, it's just a prop on a tile
like any other, but it is not what we technically call column inside the engine.

Kay, ceilings done.  Loads of copy and paste.

Collision, passability of tiles.
--------------------------------

I want to include collisions.  This is a requirement for pathfinding, which is
a requirement for IA, which is a requirement for creatures.

There are three sources of collisions in the landscape: floor, walls, objects
on the floor.  Collision on the floor may be a bit weird, but it merely means
that you are not allowed to walk over a pit, or over lava, this kind of things.
Maybe in the future I will differentiate between flying/swimming/walking
creatures, but let us forget it for now.  Some floors are unpassable.

Regarding the walls, it is obvious.  Some walls block, some do not.  Remember
that each tile is separated by two walls, this allows to make one-way walls.

As for the object placed on the floor, these are not implemented yet.  But you
can imagine massive columns, statues, big furniture.

Creatures should also block passage.

So, when issuing a command such as "Go forward", we must first check if is
possible to actually do so.  Going from a tile to a neighboring tile requires
that the floor of the neighboring tile is passable, and the wall of the current
tile leading to the neighboring tile is either absent or passable.

How do we deal with tiles that have no floor?  It is up to us.  I chose that if
there is no floor, then the tile is not passable.

Finally, I need a flag somewhere which enables to override all the passability
tests.  This is required for tests, debug and level design.  This must be
attached to the player-controlled creature only, otherwise all the creatures
will suddenly be able to walk through everything as soon as noclip mode is
activated.

Now, I can go two ways here.  I can set a `passable` flag on each tile, or I can
set a `passable` flag on each TYPE of tile.  For example, I can say that "lava
is unpassable", and then every time I put some lava it inherits the
impassability.  Or I totally decouple the lava look from its effects, and let
the level designer make sure that its level is consistent.  This is after all
the way I chose when I decided that both faces of a wall are totally distinct
objects.  I want to continue with the most free solution, and make rules an
optional and higher level abstraction.

Things that move.
-----------------
I am going to add creatures soon.  I want them to show animations.
I want some things to be able to move around and to change shape/color/texture/
light-intensity, etc.

Technically, there are two types of things that can change: uniforms and vertex
data.  Uniforms deal with things like position, rotation, illumination of the
whole mesh.  Vertex data work at the vertex level.  For example, a torch can
be implemented with a fixed vertex data but a changing fire texture via a
uniform.  A creature animation can be implemented as a complete reupload of all
its vertices.  Actually it is probably possible to do skeletal animation in the
vertex shader, but I can also just stream vertices.  In any case, it is good to
have.

I am not sure how things will end up looking, so I am trying a very quick and
dirty prototype, just to check whether I am understanding the basics.  The world
is created with one and only one dynamic object in it.  Both its position and
its vertices are modified each frame.  The dynamic object knows when it was born
and what time it is now, so it knows how to display itself.

Artificial Intelligence.
========================

Actions.
--------

Intelligence, whether artificial or natural (that of the player) results in a
creature performing an action.

Examples of actions: drink a potion, step to the left, read a scroll, pull a
lever, cast a spell, wait.

Should things like "put on a helmet" be an action?  I guess so.  We can either
consider that the character suddenly has a helmet on their head (because divine
intervention of the player) or that the character actually puts on a helmet.  I
really do not plan for the IA to manage clothes, but I don't see any reason to
prevent it either.

Turning is also an action.

Actions take time.  So, an action need to carry some state over several frames.
Each creature has a current action.  When we loop over all the creatures to run
their IA, we give that current action to the IA function.  If the IA function
sees that that action is over, then it comes up with a new action, otherwise it
returns that action.

Some actions seem not to need an object.  For example, wait, or move forward.
Some seem to need an object, like drink a potion, attack someone.
What if the object disappears?  What if a creature is attacking another, but
that other creature is killed before the hit lands?  Should I stop the action
immediately or let it reach its end?  In the case of a step, or a weapon blow,
it feels weird to interrupt suddenly because of mechanical momentum.  In the
case of reading a scroll or drinking a potion, it's less clear.  When eating a
meal or writing something, then it is clear that we need to be interrupted.

So actions should be able to be interrupted, unless they would leave the game in
an insane state.  For example, a creature needs to be standing on one tile, so
its movement must not be interrupted before it is finished.  But it is perfectly
valid to eat half your bread.  So, in a sense, it is linked to the granularity
of the action.  A bread can have a granularity of 10, indicating ten bites,
while walking forward has a granularity of 1.  Seen like that, each bite in a
bread has a granularity of 1 and cannot be interrupted.  So, in the end, no
action can be interrupted.  The actions that can are actually meta actions that
put several small actions in sequence.

If a creature has no action going, then it comes up with one the next time it
thinks.  If it has an action going, then it progresses it.  What when an action
finished during a frame?  Should a new action be generated, or should we return
the finished action?  I believe that we should return the finished action so
that we can see that 100% sometimes.  What if a finished action is sent to a
creature ?  Then it is equivalent to no action at all, the creature comes up
with a new one.  In short: really finish what you're doing before starting
something else.

Now, we know how long it takes to do something.  We have yet to know what we are
doing.  Some actions are simple: "Move forward".  Some are complicated: "Cast
fireball at goblin".  I don't think it can get more complicated.  So, four
cases?

 * verb (wait)
 * verb object (drink potion, wield sword, push button)
 * verb target (go there)
 * verb object target (use sword on goblin, use gem on sword)

Objects can be targets, but targets cannot always be object?
What about transitivity?  Is Use sword on goblin the same as use goblin on
sword?  Probably.  But Cast Goblin at Fireball?  Nope.

Maybe there is always a target.  Drink potion: use potion on self.  Push button:
use hand on button.  Wield sword: equip sword to self/hand, something like that.

Maybe there is always an object.  Go there: move self to there.

Damn, this is a mess.  It's a mess because that system is more adapted to a real
time game than a turned based game.  Turn base cut the momentum, so the
perception of speed by the player is bad.  That makes it very difficult to know
whether or not an attack will land on time, all that stuff.  I want real turns.
One turn is one move or one attack or one spell, etc.

I do not want every creature to move at the same speed, because that creates
very boring situations where a creature can follow you forever.  Let's imagine
three speeds for now: slow, normal and fast.  Fast would act every turn.  Normal
would act every other turn.  Slow would act every three turns.  So there is
some period: fast=1, normal=2, slow=3 for example.  That means we have an action
clock in addition to the tick clock.  The tick clock is used for animations, it
is driven by the game loop.  The action clock is independent from the tick clock
although it can be driven by it in order to give a real-time feeling.

In turned-based mode, the action clock stops when it is the player turn to play.
In real-time mode, it does not stop.  If the player does not take an action when
it is able to do so, then no problem, the action clock keep ticking and the
player character does nothing.  When the player decides to start an action, then
that action starts at the current action tick.

Action ticks should be quick enough, something of the order of ten times per
second.  This is to provide granularity for two different things:

 * Action periods of the order of 10 instead of 1 allow for more fine tuning of
   speed.
 * Keeps the game reactive in real-time mode.


