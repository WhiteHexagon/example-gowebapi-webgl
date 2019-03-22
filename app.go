package main

import (
	"fmt"
	"syscall/js"

	"github.com/gowebapi/webapi"
	"github.com/gowebapi/webapi/graphics/webgl"
	"github.com/gowebapi/webapi/html/canvas"
)

//see https://github.com/golang/go/wiki/WebAssembly
//see https://github.com/bobcob7/wasm-basic-triangle
func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Go/WASM loaded")

	addCanvas()
	<-c
}

func addCanvas() {
	doc := webapi.GetWindow().Document()
	app := doc.GetElementById("app")
	body := doc.GetElementById("body")
	width := body.ClientWidth()
	height := body.ClientHeight()

	canvasE := webapi.GetWindow().Document().CreateElement("canvas", &webapi.Union{js.ValueOf("dom.Node")}) //TODO this seems wrong since we have to use js. for node
	canvasE.SetId("canvas42")
	app.AppendChild(&canvasE.Node)
	canvasHTML := canvas.HTMLCanvasElementFromJS(canvasE)
	canvasHTML.SetWidth(uint(width))
	canvasHTML.SetHeight(uint(height))
	//canvasHTML.RequestFullscreen(&dom.FullscreenOptions{})	//TODO find a way to do fullscreen request

	contextU := canvasHTML.GetContext("webgl", nil)
	gl := webgl.WebGLRenderingContextFromJS(contextU)

	vBuffer, iBuffer, icount := createBuffers(gl)

	//// Shaders ////
	prog := setupShaders(gl)

	//// Associating shaders to buffer objects ////

	// Bind vertex buffer object
	gl.BindBuffer(webgl.ARRAYBUFFER_WebGLRenderingContext, vBuffer)

	// Bind index buffer object
	gl.BindBuffer(webgl.ELEMENTARRAYBUFFER_WebGLRenderingContext, iBuffer)

	// Get the attribute location
	coord := gl.GetAttribLocation(prog, "coordinates")

	// Point an attribute to the currently bound VBO
	gl.VertexAttribPointer(uint(coord), 3, webgl.FLOAT_WebGLRenderingContext, false, 0, 0)

	// Enable the attribute
	gl.EnableVertexAttribArray(uint(coord))

	//// Drawing the triangle ////

	// Clear the canvas
	gl.ClearColor(0.5, 0.5, 0.5, 0.9)
	gl.Clear(webgl.COLORBUFFERBIT_WebGLRenderingContext)

	// Enable the depth test
	gl.Enable(webgl.DEPTHTEST_WebGLRenderingContext)

	// Set the view port
	gl.Viewport(0, 0, width, height)

	// Draw the triangle
	gl.DrawElements(webgl.TRIANGLES_WebGLRenderingContext, icount, webgl.UNSIGNEDSHORT_WebGLRenderingContext, 0)

	fmt.Println("done")
}

func createBuffers(gl *webgl.WebGLRenderingContext) (*webgl.WebGLBuffer, *webgl.WebGLBuffer, int) {
	//// VERTEX BUFFER ////
	var verticesNative = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
	var vertices = js.TypedArrayOf(verticesNative)
	// Create buffer
	vBuffer := gl.CreateBuffer()
	// Bind to buffer
	gl.BindBuffer(webgl.ARRAYBUFFER_WebGLRenderingContext, vBuffer)
	// Pass data to buffer
	gl.BufferData2(webgl.ARRAYBUFFER_WebGLRenderingContext, webgl.UnionFromJS(vertices.Value), webgl.STATICDRAW_WebGLRenderingContext)
	// Unbind buffer
	gl.BindBuffer(webgl.ARRAYBUFFER_WebGLRenderingContext, &webgl.WebGLBuffer{})

	// INDEX BUFFER ////
	var indicesNative = []uint32{
		2, 1, 0,
	}
	var indices = js.TypedArrayOf(indicesNative)

	// Create buffer
	iBuffer := gl.CreateBuffer()

	// Bind to buffer
	gl.BindBuffer(webgl.ELEMENTARRAYBUFFER_WebGLRenderingContext, iBuffer)

	// Pass data to buffer
	gl.BufferData2(webgl.ELEMENTARRAYBUFFER_WebGLRenderingContext, webgl.UnionFromJS(indices.Value), webgl.STATICDRAW_WebGLRenderingContext)

	// Unbind buffer
	gl.BindBuffer(webgl.ELEMENTARRAYBUFFER_WebGLRenderingContext, &webgl.WebGLBuffer{})
	return vBuffer, iBuffer, len(indicesNative)
}

func setupShaders(gl *webgl.WebGLRenderingContext) *webgl.WebGLProgram {
	// Vertex shader source code
	vertCode := `
	attribute vec3 coordinates;

	void main(void) {
		gl_Position = vec4(coordinates, 1.0);
	}`

	// Create a vertex shader object
	vShader := gl.CreateShader(webgl.VERTEXSHADER_WebGLRenderingContext)

	// Attach vertex shader source code
	gl.ShaderSource(vShader, vertCode)

	// Compile the vertex shader
	gl.CompileShader(vShader)

	//fragment shader source code
	fragCode := `
	void main(void) {
		gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
	}`

	// Create fragment shader object
	fShader := gl.CreateShader(webgl.FRAGMENTSHADER_WebGLRenderingContext)

	// Attach fragment shader source code
	gl.ShaderSource(fShader, fragCode)

	// Compile the fragmentt shader
	gl.CompileShader(fShader)

	// Create a shader program object to store
	// the combined shader program
	prog := gl.CreateProgram()

	// Attach a vertex shader
	gl.AttachShader(prog, vShader)

	// Attach a fragment shader
	gl.AttachShader(prog, fShader)

	// Link both the programs
	gl.LinkProgram(prog)

	// Use the combined shader program object
	gl.UseProgram(prog)

	return prog
}
