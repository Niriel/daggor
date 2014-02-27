package batch

import (
	"glm"
)

type CameraBatch struct {
	BaseBatch
	context *GlContext
	view    glm.Matrix4
	proj    glm.Matrix4
}

func MakeCameraBatch(context *GlContext, view, proj glm.Matrix4) CameraBatch {
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
