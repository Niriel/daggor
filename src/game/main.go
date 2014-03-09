// game project main.go
package main

import (
	"batch"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"glm"
	"glw"
	"ia"
	"runtime"
	"sculpt"
	"time"
	"world"
)

func init() {
	// OpenGL and GLFW want to run on the main thread.
	// Or at least, want to run always from the same thread.
	runtime.LockOSThread()
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

const (
	cubeID = iota
	pyramidID
	floorID
	wallID
	columnID
	ceilingID
	monsterID
)

func viewMatrix(pos world.Position) glm.Matrix4 {
	R := glm.RotZ(float64(-90 * pos.F.Value()))
	T := glm.Vector3{float64(-pos.X), float64(-pos.Y), -.5}.Translation()
	return R.Mult(T)
}

type glState struct {
	Window           *glfw.Window
	glfwKeyEventList *glfwKeyEventList
	Shapes           [7]glw.Renderer
	context          *glw.GlContext
}

type programState struct {
	Gl    glState     // Highly mutable, impure.
	World world.World // Immutable, pure.
}

func main() {
	var programState programState
	var err error
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("GLFW initialization failed.")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.SrgbCapable, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	programState.Gl.Window, err = glfw.CreateWindow(640, 480, "Daggor", nil, nil)
	if err != nil {
		panic(err)
	}
	defer programState.Gl.Window.Destroy()

	programState.Gl.glfwKeyEventList = makeGlfwKeyEventList()
	programState.Gl.Window.SetKeyCallback(programState.Gl.glfwKeyEventList.Callback)

	programState.Gl.Window.MakeContextCurrent()
	if ec := gl.Init(); ec != 0 {
		panic(fmt.Sprintf("OpenGL initialization failed with code %v.", ec))
	}
	// For some reason, here, the OpenGL error flag for me contains "Invalid enum".
	// This is weird since I have not done anything yet.  I imagine that something
	// goes wrong in gl.Init.  Reading the error flag clears it, so I do it.
	// Here's the reason:
	//     https://github.com/go-gl/glfw3/issues/50
	// Maybe I should not even ask for a core profile anyway.
	// What are the advantages are asking for a core profile?
	if err := glw.CheckGlError(); err != nil {
		err.Description = "OpenGL has this error right after init for some reason."
		//fmt.Println(err)
	}
	major := programState.Gl.Window.GetAttribute(glfw.ContextVersionMajor)
	minor := programState.Gl.Window.GetAttribute(glfw.ContextVersionMinor)
	fmt.Printf("OpenGL version %v.%v.\n", major, minor)
	if (major < 3) || (major == 3 && minor < 3) {
		panic("OpenGL version 3.3 required, your video card/driver does not seem to support it.")
	}

	programState.Gl.context = glw.NewGlContext()

	//programState.Gl.Shapes[cubeID] = glw.Cube(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[pyramidID] = glw.Pyramid(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[floorID] = glw.Floor(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[wallID] = glw.Wall(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[columnID] = glw.Column(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[ceilingID] = glw.Ceiling(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.Shapes[monsterID] = glw.Monster(programState.Gl.context.Programs, UniformBinding)
	//programState.Gl.DynaPyramid = glw.DynaPyramid(programState.Gl.context.Programs)
	programState.Gl.Shapes[floorID] = sculpt.FloorInst(programState.Gl.context.Programs)

	// I do not like the default reference frame of OpenGl.
	// By default, we look in the direction -z, and y points up.
	// I want z to point up, and I want to look in the direction +x
	// by default.  That way, I move on an xy plane where z is the
	// altitude, instead of having the altitude stuffed between
	// the two things I use the most.  And my reason for pointing
	// toward +x is that I use the convention for trigonometry:
	// an angle of 0 points to the right (east) of the trigonometric
	// circle.  Bonus point: this matches Blender's reference frame.
	myFrame := glm.ZUP.Mult(glm.RotZ(90))
	projectionMatrix := glm.PerspectiveProj(110, 640./480., .1, 100).Mult(myFrame)
	programState.Gl.context.SetCameraProj(projectionMatrix)

	gl.Enable(gl.FRAMEBUFFER_SRGB)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)

	//blah := make([]int32, 1)
	//for _, thing := range []gl.GLenum{gl.UNIFORM_BUFFER_OFFSET_ALIGNMENT,
	//	gl.UNIFORM_BLOCK_DATA_SIZE,
	//} {
	//	gl.GetIntegerv(thing, blah)
	//	fmt.Println(blah)
	//}

	programState.World = world.MakeWorld()
	mainLoop(programState)
}

func mainLoop(programState programState) programState {
	const tickPeriod = 1000000000 / 60
	ticker := time.NewTicker(tickPeriod * time.Nanosecond)
	keepTicking := true
	for keepTicking {
		select {
		case _, ok := <-ticker.C:
			{
				if ok {
					programState, keepTicking = onTick(programState, tickPeriod)
					if !keepTicking {
						fmt.Println("No more ticks.")
						ticker.Stop()
					}
				} else {
					fmt.Println("Ticker closed, weird.")
					keepTicking = false
				}
			}
		}
	}
	return programState
}

func onTick(programState programState, dt uint64) (programState, bool) {
	glfw.PollEvents()
	keepTicking := !programState.Gl.Window.ShouldClose()
	if keepTicking {
		// Read raw inputs.
		keys := programState.Gl.glfwKeyEventList.Freeze()
		// Analyze the inputs, see what they mean.
		commands := commands(keys)
		// One of these commands may correspond to an action of the player's actor.
		// We take it out so that we can process it in the IA phase.
		// The remaining commands are kept for further processing.
		playerAction, commands := commandsToAction(commands, programState.World.Player_id)
		// Evolve the program one step.
		programState.World.Time += dt // No side effect, we own a copy.
		// $$$ THERE COULD BE SIDE EFFECTS HERE ACTUALLY:  IF I GAVE A POINTER
		// TO THE WORLD OR PROGRAM STATE TO SOMETHING.  NEED TO CORRECT THAT.
		programState = executeCommands(programState, commands)
		//
		programState.World = runAI(programState.World, playerAction)
		// render on screen.
		render(programState)
		programState.Gl.Window.SwapBuffers()
	}
	return programState, keepTicking
}

func runAI(w world.World, playerAction ia.Action) world.World {
	var action ia.Action
	// It's like on a board game.  Every one plays when it is their turn.
	// This function is called every frame.

	// Temporary: Any creature that is not scheduled yet is added to the
	// scheduler.
	schedule := w.Level.ActorSchedule
	for actorID := range w.Level.Actors.Content() {
		index := schedule.PosActorID(actorID)
		if index == -1 {
			fmt.Println("Force scheduling", actorID)
			schedule = schedule.Add(actorID, w.Time)
		}
	}
	w = w.SetActorSchedule(schedule)

	for {
		actorTime, ok := w.Level.ActorSchedule.Next(w.Time)
		if !ok {
			// Actions can modify the list of actors, so I cannot loop over
			// all the actors.  This is why I break the loop this way.
			break // No more actors to process.
		}
		newSchedule, ok := w.Level.ActorSchedule.Remove(actorTime)
		if !ok {
			panic("Could not find actor to remove from scheduler")
		}
		w = w.SetActorSchedule(newSchedule)
		if actorTime.Actor_id == w.Player_id {
			action = playerAction
		} else {
			action = ia.DecideAction(actorTime.Actor_id)
		}
		if action != nil {
			var err error
			newSchedule = newSchedule.Add(actorTime.Actor_id, actorTime.Time+100000000)
			w = w.SetActorSchedule(newSchedule)
			w, err = action.Execute(w)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// Nil actions should only happen for the player.  The player is the
			// only actor who can decide not to act.  All other actors decide
			// an action, even if it is just a waiting action.
			if actorTime.Actor_id == w.Player_id {
				// Reschedule the player for next turn.
				newSchedule = newSchedule.Add(actorTime.Actor_id, w.Time+1)
				w = w.SetActorSchedule(newSchedule)
			} else {
				panic("Only the player is allowed to idle.")
			}
		}
	}
	return w
}

func render(programState programState) {
	clearBatch := batch.MakeClearBatch(
		gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT,
		[4]gl.GLclampf{0.0, 0.0, 0.4, 0.0},
		1,
	)

	actorID := programState.World.Player_id
	position, ok := programState.World.Level.ActorPosition(actorID)
	if !ok {
		panic("Could not find player's character position.")
	}

	camBatch := batch.MakeCameraBatch(
		programState.Gl.context,
		programState.Gl.context.CameraProj(),
	)

	clearBatch.Batches = append(clearBatch.Batches, camBatch)

	clearBatch.Enter()
	clearBatch.Run()
	clearBatch.Exit()

	view := viewMatrix(position)
	renderBuildings(
		programState.World.Level.Floors,
		0, 0,
		nil,
		view,
		programState.Gl,
	)
}

func renderBuildings(
	buildings world.Buildings,
	offsetX, offsetY float64,
	defaultR *glm.Matrix4, // Can be nil.
	view glm.Matrix4,
	glState glState,
) {
	locations := make(map[world.ModelId][]glm.Matrix4)
	for coords, building := range buildings {
		m := glm.Vector3{
			float64(coords.X) + offsetX,
			float64(coords.Y) + offsetY,
			0,
		}.Translation()
		facer, ok := building.(world.Facer)
		if ok {
			// We obey the facing of the buildings that have one.
			r := glm.RotZ(float64(90 * facer.Facing().Value()))
			m = m.Mult(r)
		} else {
			// Buildings without facing receive the provided default facing.
			// It is given as a precalculated rotation matrix `defaultR`.
			m = m.Mult(*defaultR)
		}
		m = view.Mult(m) // Shaders work in view space.
		modelID := building.Model()
		meshlocs := locations[modelID]
		meshlocs = append(meshlocs, m)
		locations[modelID] = meshlocs
	}
	//
	for modelID, locs := range locations {
		mesh := glState.Shapes[modelID]
		mesh.Render(locs)
	}
}
