package sculpt

// Some hard coded meshes to try the new sculpt system.

import (
	"github.com/go-gl/gl"
	"glw"
)

func quadInstNorm(programs glw.Programs, vertexData []glw.VertexXyzNorUv) glw.Renderer {
	srefs := glw.ShaderRefs{
		glw.VSH_NOR_UV_INSTANCED,
		glw.FSH_NOR_UV,
	}

	vertices := glw.NewVerticesXyzNorUv(gl.STATIC_DRAW)
	vertices.SetData(vertexData)

	elementData := []gl.GLubyte{0, 1, 2, 3}
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

// FloorInstNorm creates a simple floor mesh.
// If there is some text written on the texture of the floor, then the floor
// at rest (non rotated) is readable by a character looking in the +x direction.
func FloorInstNorm(programs glw.Programs) glw.Renderer {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	vertexData := []glw.VertexXyzNorUv{
		// position xyz, normal xyz, uv
		glw.VertexXyzNorUv{m, m, 0, 0, 0, 1, 1, 0},
		glw.VertexXyzNorUv{p, m, 0, 0, 0, 1, 1, 1},
		glw.VertexXyzNorUv{m, p, 0, 0, 0, 1, 0, 0},
		glw.VertexXyzNorUv{p, p, 0, 0, 0, 1, 0, 1},
	}
	return quadInstNorm(programs, vertexData)
}

// CeilingInstNorm creates a simple ceiling mesh.
// If there is some text written on the texture of the ceiling, then the ceiling
// at rest (non rotated) is readable by a character looking in the +x direction.
func CeilingInstNorm(programs glw.Programs) glw.Renderer {
	const p = .5 // Plus sign.
	const m = -p // Minus sign.

	vertexData := []glw.VertexXyzNorUv{
		// position xyz, normal xyz, uv
		glw.VertexXyzNorUv{m, p, 1, 0, 0, -1, 0, 1},
		glw.VertexXyzNorUv{p, p, 1, 0, 0, -1, 0, 0},
		glw.VertexXyzNorUv{m, m, 1, 0, 0, -1, 1, 1},
		glw.VertexXyzNorUv{p, m, 1, 0, 0, -1, 1, 0},
	}
	return quadInstNorm(programs, vertexData)
}

// Creates a wall mesh.
// At rest (non rotated), the wall faces a player looking in the +x direction.
// In other words, the normal to the wall at rest is -x.
func WallInstNorm(programs glw.Programs) glw.Renderer {
	// Horizontal coordinates.
	const p = .5 // Plus sign.
	const m = -p // Minus sign.
	// Vertical coordinates.
	const P = 1 // Plus sign.
	const M = 0 // Minus sign.

	vertexData := []glw.VertexXyzNorUv{
		// position xyz, normal xyz, uv
		glw.VertexXyzNorUv{p, p, M, 1, 0, 0, 0, 0},
		glw.VertexXyzNorUv{p, m, M, 1, 0, 0, 1, 0},
		glw.VertexXyzNorUv{p, p, P, 1, 0, 0, 0, 1},
		glw.VertexXyzNorUv{p, m, P, 1, 0, 0, 1, 1},
	}
	return quadInstNorm(programs, vertexData)
}
