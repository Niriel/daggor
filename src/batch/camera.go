package batch

import (
	"glm"
	"glw"
)

type CameraBatch struct {
	BaseBatch
	context *glw.GlContext
	view    glm.Matrix4
	proj    glm.Matrix4
}

func MakeCameraBatch(context *glw.GlContext, view, proj glm.Matrix4) CameraBatch {
	return CameraBatch{
		context: context,
		view:    view,
		proj:    proj,
	}
}

func (batch CameraBatch) Enter() {
	batch.context.SetCameraView(batch.view)
	batch.context.SetCameraProj(batch.proj)
}

func (batch CameraBatch) Exit() {}
