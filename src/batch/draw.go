package batch

import (
	"fmt"
	"github.com/go-gl/gl"
	"glw"
)

type DrawElementsBufferedBatch struct {
	primitive gl.GLenum
	elements  glw.ElementIndexFormat
}

func MakeDrawElementsBufferedBatch(
	primitive gl.GLenum,
	elements glw.ElementIndexFormat,
) DrawElementsBufferedBatch {
	return DrawElementsBufferedBatch{
		primitive: primitive,
		elements:  elements,
	}
}

func (batch DrawElementsBufferedBatch) Enter() {}
func (batch DrawElementsBufferedBatch) Exit()  {}
func (batch DrawElementsBufferedBatch) Run() {
	gl.DrawElements(
		batch.primitive,
		batch.elements.Len(),
		batch.elements.GlType(),
		nil, // Indices are in a buffer, never give them directly.
	)
	if err := glw.CheckGlError(); err != nil {
		err.Description = fmt.Sprintf("DrawElementsBuffered %v", batch)
		panic(err)
	}
}
