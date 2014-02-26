package batch

import (
	"glm"
)

type CameraBatch struct {
	BaseBatch
	view glm.Matrix4
	proj glm.Matrix4
}

func MakeCameraBatch(view, proj glm.Matrix4) CameraBatch {
	return CameraBatch{view: view, proj: proj}
}

func (batch CameraBatch) Enter() {
	GlobalGlState.SetCameraView(batch.view)
	GlobalGlState.SetCameraProj(batch.proj)
}

func (batch CameraBatch) Exit() {}
