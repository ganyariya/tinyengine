package renderer

import (
	"fmt"
	"runtime"

	"github.com/ganyariya/tinyengine/pkg/tinyengine"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)


// OpenGLRenderer はOpenGLを使用した描画を提供する
type OpenGLRenderer struct {
	width         int
	height        int
	window        *glfw.Window
	shaderManager *ShaderManager
}

// NewOpenGLRenderer は新しいOpenGLRendererを作成する
func NewOpenGLRenderer(width, height int) (tinyengine.Renderer, error) {
	// GLFWの初期化はプラットフォーム層で行われているため、ここでは行わない
	// ウィンドウ作成とOpenGL初期化のみ行う

	renderer := &OpenGLRenderer{
		width:  width,
		height: height,
	}

	// ヘッドレス環境のテスト対応
	if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		// CI環境などではOpenGLが利用できない可能性があるためエラーを返す
		return nil, fmt.Errorf("OpenGL not available in headless environment")
	}

	return renderer, nil
}

// NewOpenGLRendererWithWindow はウィンドウ付きのOpenGLRendererを作成する
func NewOpenGLRendererWithWindow(width, height int, title string) (tinyengine.Renderer, error) {
	runtime.LockOSThread()

	// GLFW初期化確認
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize GLFW: %v", err)
	}

	// OpenGLヒント設定
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// ウィンドウ作成
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("failed to create window: %v", err)
	}

	window.MakeContextCurrent()

	// OpenGL初期化
	if err := gl.Init(); err != nil {
		window.Destroy()
		glfw.Terminate()
		return nil, fmt.Errorf("failed to initialize OpenGL: %v", err)
	}

	// ビューポート設定
	gl.Viewport(0, 0, int32(width), int32(height))

	// シェーダーマネージャー作成
	shaderManager := NewShaderManager()
	
	// 基本的な頂点シェーダー
	vertexShaderSource := `#version 410 core
layout (location = 0) in vec3 aPos;

uniform mat4 u_transform;

void main()
{
    gl_Position = u_transform * vec4(aPos, 1.0);
}`
	
	// 基本的なフラグメントシェーダー
	fragmentShaderSource := `#version 410 core
out vec4 FragColor;

uniform vec4 u_color;

void main()
{
    FragColor = u_color;
}`
	
	// シェーダーマネージャーでシェーダーを読み込み
	if err := shaderManager.LoadShader("basic", vertexShaderSource, fragmentShaderSource); err != nil {
		window.Destroy()
		glfw.Terminate()
		return nil, fmt.Errorf("failed to load basic shader: %v", err)
	}
	
	shaderManager.UseShader("basic")

	renderer := &OpenGLRenderer{
		width:         width,
		height:        height,
		window:        window,
		shaderManager: shaderManager,
	}

	return renderer, nil
}

// Clear は画面をクリアする
func (r *OpenGLRenderer) Clear() {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

// Present は描画内容を画面に表示する
func (r *OpenGLRenderer) Present() {
	if r.window != nil {
		r.window.SwapBuffers()
		glfw.PollEvents()
	}
}

// DrawRectangle は矩形を描画する
func (r *OpenGLRenderer) DrawRectangle(x, y, width, height float32) {
	// 基本的な矩形描画（頂点データを使用）
	vertices := []float32{
		// 位置
		x, y, // 左下
		x + width, y, // 右下
		x + width, y + height, // 右上
		x, y + height, // 左上
	}

	indices := []uint32{
		0, 1, 2, // 最初の三角形
		2, 3, 0, // 二番目の三角形
	}

	// VBO, VAO, EBO作成と描画
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// 描画
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

	// クリーンアップ
	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
}

// DrawPrimitive はプリミティブを描画する
func (r *OpenGLRenderer) DrawPrimitive(primitive interface{}) {
	if p, ok := primitive.(Primitive); ok {
		vertices := p.GetVertices()
		indices := p.GetIndices()
		color := p.GetColor()
		
		r.drawVertices(vertices, indices, color, p.GetType())
	}
}

// DrawRectangleColor は色付き矩形を描画する
func (r *OpenGLRenderer) DrawRectangleColor(x, y, width, height float32, red, g, b, a float32) {
	color := NewColor(red, g, b, a)
	rect := NewRectangle(x, y, width, height, color)
	r.DrawPrimitive(rect)
}

// DrawCircle は円を描画する
func (r *OpenGLRenderer) DrawCircle(x, y, radius float32, red, g, b, a float32) {
	color := NewColor(red, g, b, a)
	circle := NewCircle(x, y, radius, color)
	r.DrawPrimitive(circle)
}

// DrawLine は線を描画する
func (r *OpenGLRenderer) DrawLine(x1, y1, x2, y2 float32, red, g, b, a float32) {
	color := NewColor(red, g, b, a)
	line := NewLine(x1, y1, x2, y2, color)
	r.DrawPrimitive(line)
}

// drawVertices は頂点データを描画する共通メソッド
func (r *OpenGLRenderer) drawVertices(vertices []float32, indices []uint32, color Color, primitiveType PrimitiveType) {
	if r.shaderManager == nil {
		return // シェーダーマネージャーが初期化されていない場合は何もしない
	}
	
	// デバッグ情報（本来はログレベルで制御）
	// fmt.Printf("Drawing primitive type %d with color R:%.2f G:%.2f B:%.2f A:%.2f\n", 
	//	primitiveType, color.R, color.G, color.B, color.A)

	// 現在のシェーダーを取得
	currentShaderName := r.shaderManager.GetCurrentShader()
	if currentShaderName == "" {
		return
	}
	
	shader := r.shaderManager.GetShader(currentShaderName)
	if shader == nil {
		return
	}

	// VBO, VAO, EBO作成
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	// 頂点データをVBOに設定
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// インデックスデータをEBOに設定
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// 頂点属性の設定（位置のみ: x, y, z）
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// シェーダーを使用
	shader.Use()

	// 正規化デバイス座標系への変換を実行
	// 左上原点のピクセル座標系をOpenGLのNDC座標系に変換
	width := float32(r.width)
	height := float32(r.height)
	
	// 変換行列: ピクセル座標(0,0)→(-1,1), (width,height)→(1,-1)
	// 列優先行列 (column-major)
	transformMatrix := [16]float32{
		2.0/width,  0,            0, 0,
		0,          -2.0/height,  0, 0,
		0,          0,            1, 0,
		-1,         1,            0, 1,
	}
	
	// Uniform変数を設定
	transformLoc := shader.GetUniformLocation("u_transform")
	if transformLoc != -1 {
		gl.UniformMatrix4fv(transformLoc, 1, false, &transformMatrix[0])
	}
	
	colorLoc := shader.GetUniformLocation("u_color")
	if colorLoc != -1 {
		gl.Uniform4f(colorLoc, color.R, color.G, color.B, color.A)
	}
	
	// 描画タイプに応じて描画
	var drawMode uint32
	switch primitiveType {
	case PrimitiveTypeLine:
		drawMode = gl.LINES
	case PrimitiveTypeTriangle:
		drawMode = gl.TRIANGLES
	case PrimitiveTypeRectangle, PrimitiveTypeCircle:
		drawMode = gl.TRIANGLES
	default:
		drawMode = gl.TRIANGLES
	}

	// 描画実行
	gl.DrawElements(drawMode, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

	// クリーンアップ
	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
}

// GetWindow はGLFWウィンドウを取得する
func (r *OpenGLRenderer) GetWindow() *glfw.Window {
	return r.window
}

// Destroy はOpenGLリソースを解放する
func (r *OpenGLRenderer) Destroy() {
	if r.shaderManager != nil {
		r.shaderManager.DeleteAllShaders()
	}
	if r.window != nil {
		r.window.Destroy()
		glfw.Terminate()
	}
}
