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

	drawCube(gl)

	select {}
}

func drawCube(gl js.Value) {
	vertSrc := `
    attribute vec3 aPosition;
    uniform mat4 uMVP;
    void main(void) {
        gl_Position = uMVP * vec4(aPosition, 1.0);
    }`
	fragSrc := `
    void main(void) {
        gl_FragColor = vec4(0.2, 0.7, 1.0, 1.0);
    }`

	vert := compileShader(gl, gl.Get("VERTEX_SHADER"), vertSrc)
	frag := compileShader(gl, gl.Get("FRAGMENT_SHADER"), fragSrc)

	prog := gl.Call("createProgram")
	gl.Call("attachShader", prog, vert)
	gl.Call("attachShader", prog, frag)
	gl.Call("linkProgram", prog)
	gl.Call("useProgram", prog)

	vertices := []float32{
		// Front face
		-1, -1, 1,
		1, -1, 1,
		1, 1, 1,
		-1, 1, 1,
		// Back face
		-1, -1, -1,
		1, -1, -1,
		1, 1, -1,
		-1, 1, -1,
	}
	indices := []uint16{
		// Front
		0, 1, 2, 2, 3, 0,
		// Right
		1, 5, 6, 6, 2, 1,
		// Back
		5, 4, 7, 7, 6, 5,
		// Left
		4, 0, 3, 3, 7, 4,
		// Top
		3, 2, 6, 6, 7, 3,
		// Bottom
		4, 5, 1, 1, 0, 4,
	}

	vbuf := gl.Call("createBuffer")
	gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), vbuf)
	jsVerts := js.Global().Get("Float32Array").New(len(vertices))
	for i, v := range vertices {
		jsVerts.SetIndex(i, v)
	}
	gl.Call("bufferData", gl.Get("ARRAY_BUFFER"), jsVerts, gl.Get("STATIC_DRAW"))

	ibuf := gl.Call("createBuffer")
	gl.Call("bindBuffer", gl.Get("ELEMENT_ARRAY_BUFFER"), ibuf)
	jsInds := js.Global().Get("Uint16Array").New(len(indices))
	for i, v := range indices {
		jsInds.SetIndex(i, v)
	}
	gl.Call("bufferData", gl.Get("ELEMENT_ARRAY_BUFFER"), jsInds, gl.Get("STATIC_DRAW"))

	aPos := gl.Call("getAttribLocation", prog, "aPosition")
	gl.Call("enableVertexAttribArray", aPos)
	gl.Call("vertexAttribPointer", aPos, 3, gl.Get("FLOAT"), false, 0, 0)

	// Simple MVP matrix (static view)
	mvp := []float32{
		0.7, 0, 0, 0,
		0, 0.7, 0, 0,
		0, 0, 0.7, 0,
		0, 0, -5, 1,
	}
	uMVP := gl.Call("getUniformLocation", prog, "uMVP")
	jsMVP := js.Global().Get("Float32Array").New(len(mvp))
	for i, v := range mvp {
		jsMVP.SetIndex(i, v)
	}
	gl.Call("uniformMatrix4fv", uMVP, false, jsMVP)

	gl.Call("clearColor", 0, 0, 0, 1)
	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))
	gl.Call("drawElements", gl.Get("TRIANGLES"), len(indices), gl.Get("UNSIGNED_SHORT"), 0)
}

func compileShader(gl js.Value, shaderType js.Value, src string) js.Value {
	shader := gl.Call("createShader", shaderType)
	gl.Call("shaderSource", shader, src)
	gl.Call("compileShader", shader)
	return shader
}
