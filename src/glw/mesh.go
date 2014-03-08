package glw

import (
	"fmt"
	"github.com/go-gl/gl"
	"glm"
)

type Vao struct {
	programs  Programs
	srefs     ShaderRefs
	Vertices  Buffer
	Elements  Buffer
	Instances Buffer
	Uniforms  Uniforms

	// Filled in by the SetUpVao method.
	vao gl.VertexArray
}

func (vao *Vao) bind() {
	vao.vao.Bind()
	if err := CheckGlError(); err != nil {
		err.Description = "vao.vao.Bind()"
		panic(err)
	}
}

func (vao *Vao) unbind() {
	gl.VertexArray(0).Bind()
	if err := CheckGlError(); err != nil {
		err.Description = "gl.VertexArray(0).Bind()"
		panic(err)
	}
}

func (vao *Vao) updateBuffers() {
	vao.Vertices.Update()
	vao.Elements.Update()
	if vao.Instances != nil {
		vao.Instances.Update()
	}
}

func (vao *Vao) setUp() {
	if vao == nil {
		panic("setting up a nil vao")
	}
	if vao.vao != 0 {
		panic("setUp already called.")
	}
	vao.vao = gl.GenVertexArray() // Cannot fail.
	vao.bind()
	program, err := vao.programs.Serve(vao.srefs)
	if err != nil {
		panic(err)
	}

	vao.Vertices.SetUpVao(program)
	vao.Elements.SetUpVao(program)
	if vao.Instances != nil {
		vao.Instances.SetUpVao(program)
	}
	vao.Uniforms.SetUpVao(program)

	vao.unbind()
}

func (vao *Vao) DeleteVao() {
	vao.vao.Delete()
	vao.vao = 0
}

//-----------------------------------------------------------------------------
type MeshDrawer interface {
	DrawMesh([]glm.Matrix4)
}

type InstancedMesh struct {
	Vao
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

func (mesh *InstancedMesh) DrawMesh(locations []glm.Matrix4) {
	if mesh == nil {
		panic("drawing nil mesh")
	}
	if len(locations) == 0 {
		return
	}
	if mesh.vao == 0 {
		mesh.setUp()
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
	Vao
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

func (mesh *UninstancedMesh) DrawMesh(locations []glm.Matrix4) {
	if mesh == nil {
		panic("drawing nil mesh")
	}
	if len(locations) == 0 {
		return
	}
	if mesh.vao == 0 {
		mesh.setUp()
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
