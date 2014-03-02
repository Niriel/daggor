package sculpt

import (
	"fmt"
	"github.com/go-gl/gl"
	"glw"
)

type Mesh struct {
	// Filled in by the CreateMesh constructor.
	programs *glw.Programs
	srefs    glw.ShaderRefs
	Vertices Vertices
	Elements Elements
	Uniforms Uniforms
	drawer   Drawer
	// Filled in by the SetUpVao method.
	vao gl.VertexArray
}

// MakeMesh creates a mesh from the provided components.
// It does not make any OpenGL calls.
func MakeMesh(p *glw.Programs, s glw.ShaderRefs, v Vertices, e Elements, u Uniforms, d Drawer) Mesh {
	return Mesh{
		programs: p,
		srefs:    s,
		Vertices: v,
		Elements: e,
		Uniforms: u,
		drawer:   d,
	}
}

func (mesh Mesh) bind() {
	mesh.vao.Bind()
	if err := glw.CheckGlError(); err != nil {
		err.Description = "mesh.vao.Bind()"
		panic(err)
	}
}

func (mesh Mesh) unbind() {
	gl.VertexArray(0).Bind()
	if err := glw.CheckGlError(); err != nil {
		err.Description = "gl.VertexArray(0).Bind()"
		panic(err)
	}
}

func (mesh Mesh) updateBuffers() {
	mesh.Vertices.UpdateBuffer()
	if mesh.Vertices.BufferName() == 0 {
		panic("Vertices buffer name is still 0.")
	}
	mesh.Elements.UpdateBuffer()
	if mesh.Vertices.BufferName() == 0 {
		panic("Elements buffer name is still 0.")
	}
}

func (mesh *Mesh) SetUpVao() {
	if mesh.vao != 0 {
		panic("SetUpVao already called.")
	}
	mesh.vao = gl.GenVertexArray() // Cannot fail.
	mesh.bind()
	program, err := mesh.programs.Serve(mesh.srefs)
	if err != nil {
		panic(err)
	}

	// The buffers need to be created in order to be linked to the VAO.
	// The call to updateBuffers will create the buffers.
	mesh.updateBuffers()
	mesh.Vertices.SetUpVao(program)
	mesh.Uniforms.SetUpVao(program)
	mesh.Elements.SetUpVao()

	mesh.unbind()
}

func (mesh *Mesh) DeleteVao() {
	mesh.vao.Delete()
	mesh.vao = 0
}

func (mesh *Mesh) Draw() {
	// Assume the program is used.
	// Assume the textures are bound.
	// We also assume that each mesh has its own ebo.
	program, err := mesh.programs.Serve(mesh.srefs)
	if err != nil {
		panic(err)
	}
	program.Use()
	mesh.updateBuffers()
	mesh.bind()
	mesh.Vertices.bind()
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
	mesh.drawer.Draw()
	mesh.Vertices.unbind()
	mesh.unbind()
	gl.ProgramUnuse()
}
