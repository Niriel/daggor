package glw

import (
	"fmt"
	"github.com/go-gl/gl"
)

var ERROR_NAMES = map[gl.GLenum]string{
	gl.NO_ERROR:                      "no error",
	gl.INVALID_ENUM:                  "invalid enum",
	gl.INVALID_VALUE:                 "invalid value",
	gl.INVALID_OPERATION:             "invalid operation",
	gl.STACK_OVERFLOW:                "stack overflow",
	gl.STACK_UNDERFLOW:               "stack underflow",
	gl.OUT_OF_MEMORY:                 "out of memory",
	gl.INVALID_FRAMEBUFFER_OPERATION: "invalid frame buffer operation",
	gl.TABLE_TOO_LARGE:               "table too large",
}

type GlError struct {
	error_code  gl.GLenum
	Description string
}

func (self GlError) Error() string {
	error_name, ok := ERROR_NAMES[self.error_code]
	if !ok {
		error_name = fmt.Sprintf("unknown error code %v", self.error_code)
	}
	return fmt.Sprintf("OpenGL [%v] %v", error_name, self.Description)
}

func CheckGlError() *GlError {
	ec := gl.GetError()
	if ec == gl.NO_ERROR {
		return nil
	}
	err := GlError{error_code: ec}
	return &err
}
