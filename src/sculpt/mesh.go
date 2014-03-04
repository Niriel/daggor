package sculpt

import (
	"fmt"
	"github.com/go-gl/gl"
	"glw"
)

type Mesh struct {
	// Filled in by the MakeMesh constructor.
	programs  *glw.Programs
	srefs     glw.ShaderRefs
	Vertices  Buffer
	Elements  Buffer
	Instances Buffer
	Uniforms  Uniforms
	Drawer    Drawer
	// Filled in by the SetUpVao method.
	vao gl.VertexArray
}

// MakeMesh creates a mesh from the provided components.
// It does not make any OpenGL calls.
func NewMesh(p *glw.Programs, s glw.ShaderRefs, v, e, i Buffer, u Uniforms, d Drawer) *Mesh {
	if p == nil {
		panic("trying to create a mesh with programs=nil")
	}
	mesh := Mesh{
		programs:  p,
		srefs:     s,
		Vertices:  v,
		Elements:  e,
		Instances: i,
		Uniforms:  u,
		Drawer:    d,
	}
	return &mesh
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

func (mesh *Mesh) SetUpVao() {
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

func (mesh *Mesh) Draw() {
	if mesh == nil {
		panic("drawing nil mesh")
	}
	// This needs to go away soon.
	program, err := mesh.programs.Serve(mesh.srefs)
	if err != nil {
		panic(err)
	}

	program.Use()
	mesh.updateBuffers()
	mesh.bind()
	mesh.Uniforms.SetGl()
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
	mesh.Drawer.Draw()
	mesh.unbind()
	gl.ProgramUnuse()
}
