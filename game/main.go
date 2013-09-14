// game project main.go
package main

import (
	"fmt"
	glfw "github.com/go-gl/glfw3"
)

func main() {
	fmt.Println("Hello World!")

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()
}
