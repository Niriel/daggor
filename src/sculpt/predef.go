package sculpt

// Some hard coded meshes to try the new sculpt system.

import (
	"github.com/go-gl/gl"
	"glw"
)

func FloorInstNorm(programs glw.Programs) glw.Renderer {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.

	srefs := glw.ShaderRefs{
		glw.VSH_NOR_UV_INSTANCED,
		glw.FSH_NOR_UV,
	}

	vertexData := []glw.VertexXyzNorUv{
		// position xyz, normal xyz, uv
		glw.VertexXyzNorUv{m, m, 0, m, m, 0, 0, 0},
		glw.VertexXyzNorUv{m, p, 0, m, p, 0, 0, 1},
		glw.VertexXyzNorUv{p, m, 0, p, m, 0, 1, 0},
		glw.VertexXyzNorUv{p, p, 0, p, p, 0, 1, 1},
		glw.VertexXyzNorUv{m * .75, m * .75, 0, 0, 0, 1, .5, .5},
		glw.VertexXyzNorUv{m * .75, m * .75, 0, 0, 0, -1, .5, .5},
	}
	vertices := glw.NewVerticesXyzNorUv(gl.STATIC_DRAW)
	vertices.SetData(vertexData)

	elementData := []gl.GLubyte{
		0, 2, 4,
		2, 3, 4,
		3, 1, 4,
		1, 0, 4,
		2, 0, 5,
		3, 2, 5,
		1, 3, 5,
		0, 1, 5,
	}
	elements := glw.NewElementsUbyte(gl.STATIC_DRAW)
	elements.SetData(elementData)

	instances := glw.NewModelMatInstances(gl.STREAM_DRAW)

	uniforms := new(glw.UniformsLocInstanced)

	drawer := glw.MakeDrawElementInstanced(
		gl.TRIANGLES,
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
