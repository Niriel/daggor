package sculpt

import (
	"fmt"
	"github.com/go-gl/gl"
	"glm"
	"glw"
)

type MeshDrawer interface {
	DrawMesh([]glm.Matrix4)
}

type Mesh struct {
	// Filled in by the MakeMesh constructor.
	programs  *glw.Programs
	srefs     glw.ShaderRefs
	Vertices  Buffer
	Elements  Buffer
	Instances Buffer
	Uniforms  Uniforms

	// Filled in by the SetUpVao method.
	vao gl.VertexArray
}

func (mesh *Mesh) bind() {
	mesh.vao.Bind()
	if err := glw.CheckGlError(); err != nil {
		err.Description = "mesh.vao.Bind()"
		panic(err)
	}
}

func (mesh *Mesh) unbind() {
	gl.VertexArray(0).Bind()
	if err := glw.CheckGlError(); err != nil {
		err.Description = "gl.VertexArray(0).Bind()"
		panic(err)
	}
}

func (mesh *Mesh) updateBuffers() {
	mesh.Vertices.Update()
	mesh.Elements.Update()
	if mesh.Instances != nil {
		mesh.Instances.Update()
	}
}

func (mesh *Mesh) setUpVao() {
	if mesh == nil {
		panic("setting up the vao of a nil mesh")
	}
	if mesh.vao != 0 {
		panic("SetUpVao already called.")
	}
	mesh.vao = gl.GenVertexArray() // Cannot fail.
	mesh.bind()
	program, err := mesh.programs.Serve(mesh.srefs)
	if err != nil {
		panic(err)
	}

	mesh.Vertices.SetUpVao(program)
	mesh.Elements.SetUpVao(program)
	if mesh.Instances != nil {
		mesh.Instances.SetUpVao(program)
	}
	mesh.Uniforms.SetUpVao(program)

	mesh.unbind()
}

func (mesh *Mesh) DeleteVao() {
	mesh.vao.Delete()
	mesh.vao = 0
}

//-----------------------------------------------------------------------------
type InstancedMesh struct {
	Mesh
	Drawer InstancedDrawer
}

func NewInstancedMesh(p *glw.Programs,
	s glw.ShaderRefs,
	v, e, i Buffer,
	u Uniforms, d InstancedDrawer) *InstancedMesh {
	if p == nil {
		panic("trying to create a mesh with programs=nil")
	}
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
		mesh.setUpVao()
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
	Mesh
	Drawer UninstancedDrawer
}

func NewUninstancedMesh(
	p *glw.Programs,
	s glw.ShaderRefs,
	v, e Buffer,
	u Uniforms,
	d UninstancedDrawer,
) *UninstancedMesh {
	if p == nil {
		panic("trying to create a mesh with programs=nil")
	}
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
		mesh.setUpVao()
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
	if err := glw.CheckGlError(); err != nil {
		err.Description = "program.Validate failed"
		panic(err)
	}
	status := program.Get(gl.VALIDATE_STATUS)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "program.Get(VALIDATE_STATUS) failed"
		panic(err)
	}
	if status == gl.FALSE {
		infolog := program.GetInfoLog()
		gl.GetError() // Clear error flag if infolog derped.
		panic(fmt.Errorf("program validation failed with log: %v", infolog))
	}
}
