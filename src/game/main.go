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
	"sort"
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

func viewMatrix(pos world.Position) (glm.Matrix4, glm.Matrix4) {
	const eyeZ = .5
	Rd := glm.RotZ(float64(-90 * pos.F.Value()))
	Td := glm.Vector3{float64(-pos.X), float64(-pos.Y), -eyeZ}.Translation()
	Ri := glm.RotZ(float64(90 * pos.F.Value()))
	Ti := glm.Vector3{float64(pos.X), float64(pos.Y), eyeZ}.Translation()
	Vd := Rd.Mult(Td) // Direct.
	Vi := Ti.Mult(Ri) // Inverse.
	return Vd, Vi
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

	programState.Gl.Shapes[floorID] = sculpt.FloorInstNorm(programState.Gl.context.Programs)
	programState.Gl.Shapes[ceilingID] = sculpt.CeilingInstNorm(programState.Gl.context.Programs)
	programState.Gl.Shapes[wallID] = sculpt.WallInstNorm(programState.Gl.context.Programs)

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
	gl.Enable(gl.TEXTURE_CUBE_MAP_SEAMLESS)

	glw.LoadSkybox()

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

//=============================================================================

// Positions holds model-to-eye matrices.
// Positions satisfies sort.Interface, sorting from closest to farthest.
type Positions []glm.Matrix4

// Len is required by sort.Interface.
func (positions Positions) Len() int {
	return len(positions)
}

// Less is required by sort.Interface.
func (positions Positions) Less(i, j int) bool {
	ix, iy := positions[i][12], positions[i][13]
	jx, jy := positions[j][12], positions[j][13]
	id := ix*ix + iy*iy
	jd := jx*jx + jy*jy
	return id < jd
}

// Swap is required by sort.Interface.
func (positions Positions) Swap(i, j int) {
	positions[i], positions[j] = positions[j], positions[i]
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

	worldToEye, eyeToWorld := viewMatrix(position)
	programState.Gl.context.SetCameraViewI(eyeToWorld)

	//
	verticalPositions := make(map[world.ModelId]Positions)
	horizontalPositions := make(map[world.ModelId]Positions)

	gatherBuildingsPositions(
		horizontalPositions,
		programState.World.Level.Floors,
		0, 0,
		nil,
		worldToEye,
	)
	gatherBuildingsPositions(
		horizontalPositions,
		programState.World.Level.Ceilings,
		0, 0,
		nil,
		worldToEye,
	)
	for i := 0; i < 4; i++ {
		rot := glm.RotZ(180 + 90*float64(i))
		gatherBuildingsPositions(
			verticalPositions,
			programState.World.Level.Walls[i],
			0, 0,
			&rot,
			worldToEye,
		)
	}
	// Finally render all the things.
	// Reduce fill rate by drawing the closest objects first and making use of
	// the depth test to cull fragments before expensive lightings computations.
	for rendererID, pos := range verticalPositions {
		sort.Sort(pos)
		programState.Gl.Shapes[rendererID].Render(pos)
	}
	for rendererID, pos := range horizontalPositions {
		sort.Sort(pos)
		programState.Gl.Shapes[rendererID].Render(pos)
	}
}

func gatherBuildingsPositions(
	rendererPositions map[world.ModelId]Positions,
	buildings world.Buildings,
	offsetX, offsetY float64,
	defaultR *glm.Matrix4, // Can be nil.
	worldToEye glm.Matrix4,
) {
	for coords, building := range buildings {
		position := glm.Vector3{
			float64(coords.X) + offsetX,
			float64(coords.Y) + offsetY,
			0,
		}.Translation()
		facer, ok := building.(world.Facer)
		if ok {
			// We obey the facing of the buildings that have one.
			r := glm.RotZ(float64(90 * facer.Facing().Value()))
			position = position.Mult(r)
		} else {
			// Buildings without facing receive the provided default facing.
			// It is given as a precalculated rotation matrix `defaultR`.
			position = position.Mult(*defaultR)
		}
		position = worldToEye.Mult(position) // Shaders work in view space.
		rendererID := building.Model()
		positions := rendererPositions[rendererID]
		positions = append(positions, position)
		rendererPositions[rendererID] = positions
	}
}
