package sculpt

// Some hard coded meshes to try the new sculpt system.

import (
	"github.com/go-gl/gl"
	"glw"
)

// Floor creates the mesh for a floor.
// It makes no call to OpenGL whatsoever.
// This can even be called before the context is created.
func Floor(programs glw.Programs) glw.MeshDrawer {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.

	srefs := glw.ShaderRefs{glw.VSH_COL3, glw.FSH_VCOL}

	vertexData := []glw.VertexXyzRgb{
		// x y z r v b
		glw.VertexXyzRgb{m, m, 0, .1, .1, .5},
		glw.VertexXyzRgb{m, p, 0, .1, .1, .5},
		glw.VertexXyzRgb{p, m, 0, 0, 1, 0},
		glw.VertexXyzRgb{p, p, 0, 1, 0, 0},
	}
	vertices := glw.NewVerticesXyzRgb(gl.STATIC_DRAW)
	vertices.SetData(vertexData)

	elementData := []gl.GLubyte{0, 2, 1, 3}
	elements := glw.NewElementsUbyte(gl.STATIC_DRAW)
	elements.SetData(elementData)

	uniforms := new(glw.UniformsLoc)

	drawer := glw.MakeDrawElement(
		gl.TRIANGLE_STRIP,
		len(elementData),
		gl.UNSIGNED_BYTE,
	)

	return glw.NewUninstancedMesh(
		programs,
		srefs,
		vertices,
		elements,
		uniforms,
		&drawer,
	)
}

// Floor creates the mesh for a floor.
// It makes no call to OpenGL whatsoever.
// This can even be called before the context is created.
func FloorInst(programs glw.Programs) glw.MeshDrawer {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.

	srefs := glw.ShaderRefs{glw.VSH_COL3_INSTANCED, glw.FSH_VCOL}

	vertexData := []glw.VertexXyzRgb{
		// x y z r v b
		glw.VertexXyzRgb{m, m, 0, .1, .1, .5},
		glw.VertexXyzRgb{m, p, 0, .1, .1, .5},
		glw.VertexXyzRgb{p, m, 0, 0, 1, 0},
		glw.VertexXyzRgb{p, p, 0, 1, 0, 0},
	}
	vertices := glw.NewVerticesXyzRgb(gl.STATIC_DRAW)
	vertices.SetData(vertexData)

	elementData := []gl.GLubyte{0, 2, 1, 3}
	elements := glw.NewElementsUbyte(gl.STATIC_DRAW)
	elements.SetData(elementData)

	instances := glw.NewModelMatInstances(gl.STREAM_DRAW)

	uniforms := new(glw.UniformsLocInstanced)

	drawer := glw.MakeDrawElementInstanced(
		gl.TRIANGLE_STRIP,
		len(elementData),
		gl.UNSIGNED_BYTE,
	)
	return glw.NewInstancedMesh(
		programs,
		srefs,
		vertices,
		elements,
		instances,
		uniforms,
		&drawer,
	)
}
