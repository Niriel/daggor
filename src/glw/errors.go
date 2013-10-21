package glw

import (
	"bytes"
	"fmt"
	"github.com/go-gl/gl"
	"runtime"
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

func retreiveStackTrace() *bytes.Buffer {
	const SIZE = 2048
	var stack []byte
	size := SIZE
	for i := 0; i <= 6; i++ {
		stack = make([]byte, size)
		written := runtime.Stack(stack, true)
		if written < size {
			stack = stack[:written]
			break
		}
		size *= 2
	}
	return bytes.NewBuffer(stack)
}

type GlError struct {
	Error_code  gl.GLenum
	Stack       *bytes.Buffer
	Description string
}

func (self GlError) Error() string {
	return self.String()
}
func (self GlError) String() string {
	error_name, ok := ERROR_NAMES[self.Error_code]
	if !ok {
		error_name = fmt.Sprintf("unknown error code %v", self.Error_code)
	}
	var stack_string string
	if self.Stack == nil {
		stack_string = "No stack."
	} else {
		stack_string = self.Stack.String()
	}
	return fmt.Sprintf(
		"OpenGL [%v] %v\nStacktrace:\n%v\n",
		error_name,
		self.Description,
		stack_string,
	)
}

func CheckGlError() *GlError {
	ec := gl.GetError()
	if ec == gl.NO_ERROR {
		return nil
	}
	err := GlError{
		Error_code: ec,
		Stack:      retreiveStackTrace(),
	}
	return &err
}
