package glw

import (
	"fmt"
	"github.com/go-gl/gl"
	"glm"
)

type Renderer interface {
	Render(locations []glm.Matrix4)
	SetUp()
}

type renderer struct {
	programs  Programs
	srefs     ShaderRefs
	Vertices  Buffer
	Elements  Buffer
	Instances Buffer
	Uniforms  Uniforms

	// Filled in by the SetUp method.
	vao gl.VertexArray
}

func (renderer *renderer) bind() {
	renderer.vao.Bind()
	if err := CheckGlError(); err != nil {
		err.Description = "renderer.vao.Bind()"
		panic(err)
	}
}

func (renderer *renderer) unbind() {
	gl.VertexArray(0).Bind()
	if err := CheckGlError(); err != nil {
		err.Description = "gl.VertexArray(0).Bind()"
		panic(err)
	}
}

func (renderer *renderer) updateBuffers() {
	renderer.Vertices.Update()
	renderer.Elements.Update()
	if renderer.Instances != nil {
		renderer.Instances.Update()
	}
}

func (renderer *renderer) SetUp() {
	if renderer == nil {
		panic("setting up a nil renderer")
	}
	if renderer.vao != 0 {
		return
	}
	renderer.vao = gl.GenVertexArray() // Cannot fail.
	renderer.bind()
	program, err := renderer.programs.Serve(renderer.srefs)
	if err != nil {
		panic(err)
	}

	renderer.Vertices.SetUpVao(program)
	renderer.Elements.SetUpVao(program)
	if renderer.Instances != nil {
		renderer.Instances.SetUpVao(program)
	}
	renderer.Uniforms.SetUpVao(program)

	renderer.unbind()
}

func (renderer *renderer) Delete() {
	renderer.vao.Delete()
	renderer.vao = 0
}

//-----------------------------------------------------------------------------

type InstancedMesh struct {
	renderer
	Drawer InstancedDrawer
}

func NewInstancedMesh(
	p Programs,
	s ShaderRefs,
	v, e, i Buffer,
	u Uniforms, d InstancedDrawer) *InstancedMesh {
	mesh := new(InstancedMesh)
	mesh.programs = p
	mesh.srefs = s
	mesh.Vertices = v
	mesh.Elements = e
	mesh.Instances = i
	mesh.Uniforms = u
	mesh.Drawer = d
	return mesh
}

func (mesh *InstancedMesh) Render(locations []glm.Matrix4) {
	if mesh == nil {
		panic("drawing nil mesh")
	}
	if len(locations) == 0 {
		return
	}
	if mesh.vao == 0 {
		mesh.SetUp()
	}
	// This needs to go away soon.
	program, err := mesh.programs.Serve(mesh.srefs)
	if err != nil {
		panic(err)
	}
	program.Use()

	if instances, ok := mesh.Instances.(locationDataSetter); ok {
		instances.SetLocationData(locations)
	} else {
		panic("mesh instance buffer refuses locations")
	}

	mesh.updateBuffers()
	mesh.bind()
	mesh.Uniforms.SetGl()
	validateProgram(program)
	mesh.Drawer.Draw(len(locations))
	mesh.unbind()
	gl.ProgramUnuse()
}

//-----------------------------------------------------------------------------
type UninstancedMesh struct {
	renderer
	Drawer UninstancedDrawer
}

func NewUninstancedMesh(
	p Programs,
	s ShaderRefs,
	v, e Buffer,
	u Uniforms,
	d UninstancedDrawer,
) *UninstancedMesh {
	mesh := new(UninstancedMesh)
	mesh.programs = p
	mesh.srefs = s
	mesh.Vertices = v
	mesh.Elements = e
	mesh.Instances = nil
	mesh.Uniforms = u
	mesh.Drawer = d
	return mesh
}

func (mesh *UninstancedMesh) Render(locations []glm.Matrix4) {
	if mesh == nil {
		panic("drawing nil mesh")
	}
	if len(locations) == 0 {
		return
	}
	if mesh.vao == 0 {
		mesh.SetUp()
	}
	// This needs to go away soon.
	program, err := mesh.programs.Serve(mesh.srefs)
	if err != nil {
		panic(err)
	}
	program.Use()

	mesh.updateBuffers()
	mesh.bind()
	for _, location := range locations {
		mesh.Uniforms.SetLocation(location)
		mesh.Uniforms.SetGl()
		validateProgram(program)
		mesh.Drawer.Draw()
	}
	mesh.unbind()
	gl.ProgramUnuse()
}

//-----------------------------------------------------------------------------
func validateProgram(program gl.Program) {
	program.Validate()
	if err := CheckGlError(); err != nil {
		err.Description = "program.Validate failed"
		panic(err)
	}
	status := program.Get(gl.VALIDATE_STATUS)
	if err := CheckGlError(); err != nil {
		err.Description = "program.Get(VALIDATE_STATUS) failed"
		panic(err)
	}
	if status == gl.FALSE {
		infolog := program.GetInfoLog()
		gl.GetError() // Clear error flag if infolog derped.
		panic(fmt.Errorf("program validation failed with log: %v", infolog))
	}
}
