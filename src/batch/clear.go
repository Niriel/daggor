package batch

import (
	"github.com/go-gl/gl"
	"glw"
)

type ClearBatch struct {
	BaseBatch
	clearMask  gl.GLbitfield
	clearColor [4]gl.GLclampf
	clearDepth gl.GLclampd
}

func MakeClearBatch(
	clearMask gl.GLbitfield,
	clearColor [4]gl.GLclampf,
	clearDepth gl.GLclampd,
) ClearBatch {
	return ClearBatch{
		clearMask:  clearMask,
		clearColor: clearColor,
		clearDepth: clearDepth,
	}
}

func (batch ClearBatch) Enter() {
	if batch.clearMask&gl.COLOR_BUFFER_BIT != 0 {
		gl.ClearColor(
			batch.clearColor[0],
			batch.clearColor[1],
			batch.clearColor[2],
			batch.clearColor[3],
		)
	}
	if batch.clearMask&gl.DEPTH_BUFFER_BIT != 0 {
		gl.ClearDepth(batch.clearDepth)
	}
	gl.Clear(batch.clearMask)
	if err := glw.CheckGlError(); err != nil {
		err.Description = "gl.Clear"
		panic(err)
	}
}
