//go:build js && wasm
// +build js,wasm

package main

import (
	"syscall/js"
)

func main() {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "glcanvas")
	gl := canvas.Call("getContext", "webgl")

	vertSrc := `
    attribute vec3 aPosition;
    void main(void) {
        gl_Position = vec4(aPosition, 1.0);
    }`
	fragSrc := `
    void main(void) {
        gl_FragColor = vec4(1, 0, 0, 1);
    }`

	vert := compileShader(gl, gl.Get("VERTEX_SHADER"), vertSrc)
	frag := compileShader(gl, gl.Get("FRAGMENT_SHADER"), fragSrc)

	prog := gl.Call("createProgram")
	gl.Call("attachShader", prog, vert)
	gl.Call("attachShader", prog, frag)
	gl.Call("linkProgram", prog)
	gl.Call("useProgram", prog)

	vertices := []float32{
		0, 1, 0,
		-1, -1, 0,
		1, -1, 0,
	}

	buf := gl.Call("createBuffer")
	gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), buf)

	jsArray := js.Global().Get("Float32Array").New(len(vertices))
	for i, v := range vertices {
		jsArray.SetIndex(i, v)
	}
	gl.Call("bufferData", gl.Get("ARRAY_BUFFER"), jsArray, gl.Get("STATIC_DRAW"))

	aPos := gl.Call("getAttribLocation", prog, "aPosition")
	gl.Call("enableVertexAttribArray", aPos)
	gl.Call("vertexAttribPointer", aPos, 3, gl.Get("FLOAT"), false, 0, 0)

	gl.Call("clearColor", 0, 0, 0, 1)
	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))
	gl.Call("drawArrays", gl.Get("TRIANGLES"), 0, 3)

	select {}
}

func compileShader(gl js.Value, shaderType js.Value, src string) js.Value {
	shader := gl.Call("createShader", shaderType)
	gl.Call("shaderSource", shader, src)
	gl.Call("compileShader", shader)
	return shader
}
