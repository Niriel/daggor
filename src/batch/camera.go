package batch

import (
	"glm"
	"glw"
)

type CameraBatch struct {
	BaseBatch
	context *glw.GlContext
	proj    glm.Matrix4
}

func MakeCameraBatch(context *glw.GlContext, proj glm.Matrix4) CameraBatch {
	return CameraBatch{
		context: context,
		proj:    proj,
	}
}

func (batch CameraBatch) Enter() {
	batch.context.SetCameraProj(batch.proj)
}

func (batch CameraBatch) Exit() {}
