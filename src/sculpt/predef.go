package sculpt

// Some hard coded meshes to try the new sculpt system.

import (
	"github.com/go-gl/gl"
	"glw"
)

// Floor creates the mesh for a floor.
// It makes no call to OpenGL whatsoever.
// This can even be called before the context is created.
func Floor(programs *glw.Programs) Mesh {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.

	srefs := glw.ShaderRefs{glw.VSH_COL3, glw.FSH_VCOL}

	vertexData := []VertexXyzRgb{
		// x y z r v b
		VertexXyzRgb{m, m, 0, .1, .1, .5},
		VertexXyzRgb{m, p, 0, .1, .1, .5},
		VertexXyzRgb{p, m, 0, 0, 1, 0},
		VertexXyzRgb{p, p, 0, 1, 0, 0},
	}
	vertices := new(VerticesXyzRgb)
	vertices.SetVertexData(vertexData)
	vertices.usage = gl.STATIC_DRAW

	elementData := []gl.GLubyte{0, 2, 1, 3}
	elements := new(ElementsUbyte)
	elements.SetElementData(elementData)

	uniforms := UniformsLoc{}

	drawer := DrawElement{
		mode:    gl.TRIANGLE_STRIP,
		count:   len(elementData),
		typ:     gl.UNSIGNED_BYTE,
		indices: nil,
	}
	return MakeMesh(
		programs,
		srefs,
		vertices,
		elements,
		&uniforms,
		drawer,
	)
}
