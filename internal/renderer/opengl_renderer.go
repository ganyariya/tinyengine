package renderer

import (
	"fmt"
	"runtime"

	"github.com/ganyariya/tinyengine/pkg/tinyengine"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// OpenGL設定の定数
const (
	OpenGLMajorVersion    = 4
	OpenGLMinorVersion    = 1
	VertexPositionAttrib  = 0
	VertexPositionSize    = 3
	FloatSizeBytes        = 4
	DefaultBufferPoolSize = 100
)

// デフォルトカラー設定
var (
	DefaultClearColor = [4]float32{0.0, 0.0, 0.0, 1.0} // 黒背景
)

// デフォルトシェーダーソースコード
const (
	BasicVertexShaderSource = `#version 410 core
layout (location = 0) in vec3 aPos;

uniform mat4 u_transform;

void main()
{
    gl_Position = u_transform * vec4(aPos, 1.0);
}`

	BasicFragmentShaderSource = `#version 410 core
out vec4 FragColor;

uniform vec4 u_color;

void main()
{
    FragColor = u_color;
}`
)


// OpenGLRenderer はOpenGLを使用した描画を提供する
type OpenGLRenderer struct {
	width         int
	height        int
	window        *glfw.Window
	shaderManager *ShaderManager
	bufferPool    *BufferPool
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
	glfw.WindowHint(glfw.ContextVersionMajor, OpenGLMajorVersion)
	glfw.WindowHint(glfw.ContextVersionMinor, OpenGLMinorVersion)
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
	
	// シェーダーマネージャーでシェーダーを読み込み
	if err := shaderManager.LoadShader("basic", BasicVertexShaderSource, BasicFragmentShaderSource); err != nil {
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
		bufferPool:    NewBufferPool(DefaultBufferPoolSize),
	}

	return renderer, nil
}

// Clear は画面をクリアする
func (r *OpenGLRenderer) Clear() {
	gl.ClearColor(DefaultClearColor[0], DefaultClearColor[1], DefaultClearColor[2], DefaultClearColor[3])
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
	// より効率的な描画のためにDrawPrimitiveを使用
	rect := NewRectangle(x, y, width, height, NewColor(1.0, 1.0, 1.0, 1.0))
	r.DrawPrimitive(rect)
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
func (r *OpenGLRenderer) DrawRectangleColor(x, y, width, height float32, red, green, blue, alpha float32) {
	color := NewColor(red, green, blue, alpha)
	rect := NewRectangle(x, y, width, height, color)
	r.DrawPrimitive(rect)
}

// DrawCircle は円を描画する
func (r *OpenGLRenderer) DrawCircle(x, y, radius float32, red, green, blue, alpha float32) {
	color := NewColor(red, green, blue, alpha)
	circle := NewCircle(x, y, radius, color)
	r.DrawPrimitive(circle)
}

// DrawLine は線を描画する
func (r *OpenGLRenderer) DrawLine(x1, y1, x2, y2 float32, red, green, blue, alpha float32) {
	color := NewColor(red, green, blue, alpha)
	line := NewLine(x1, y1, x2, y2, color)
	r.DrawPrimitive(line)
}

// drawVertices は頂点データを描画する共通メソッド
func (r *OpenGLRenderer) drawVertices(vertices []float32, indices []uint32, color Color, primitiveType PrimitiveType) {
	if r.shaderManager == nil {
		return // シェーダーマネージャーが初期化されていない場合は何もしない
	}
	

	// 現在のシェーダーを取得
	currentShaderName := r.shaderManager.GetCurrentShader()
	if currentShaderName == "" {
		return
	}
	
	shader := r.shaderManager.GetShader(currentShaderName)
	if shader == nil {
		return
	}

	// VBO, VAO, EBO取得（プールから再利用 or 新規作成）
	vao := r.bufferPool.GetVAO()
	vbo := r.bufferPool.GetVBO()
	ebo := r.bufferPool.GetEBO()
	
	// defer文でリソースの確実な返却を保証
	defer func() {
		gl.BindVertexArray(0)
		r.bufferPool.ReturnVAO(vao)
		r.bufferPool.ReturnVBO(vbo)
		r.bufferPool.ReturnEBO(ebo)
	}()

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
	// ピクセル座標 (0,0) = 左上 → NDC (-1,1)
	// ピクセル座標 (width,height) = 右下 → NDC (1,-1)
	
	// 現在のフレームバッファサイズを取得（ウィンドウサイズ変更に対応）
	var fbWidth, fbHeight int32
	if r.window != nil {
		w, h := r.window.GetFramebufferSize()
		fbWidth, fbHeight = int32(w), int32(h)
		// ビューポートも現在のサイズに合わせて更新
		gl.Viewport(0, 0, fbWidth, fbHeight)
	} else {
		fbWidth, fbHeight = int32(r.width), int32(r.height)
	}
	
	width := float32(fbWidth)
	height := float32(fbHeight)
	
	// 正射投影行列 (Orthographic Projection)
	// 左上原点のピクセル座標系 → OpenGL NDC座標系
	// ピクセル座標 Y=0 (上) → NDC Y=1 (上)
	// ピクセル座標 Y=height (下) → NDC Y=-1 (下)
	transformMatrix := [16]float32{
		2.0/width,   0,            0, 0,  // X: [0,width] → [-1,1]
		0,           -2.0/height,  0, 0,  // Y: [0,height] → [1,-1] (反転)
		0,           0,            1, 0,  // Z: そのまま
		-1,          1,            0, 1,  // 平行移動: (0,0)→(-1,1)
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
	
	// クリーンアップはdefer文で処理
}

// GetWindow はGLFWウィンドウを取得する
func (r *OpenGLRenderer) GetWindow() *glfw.Window {
	return r.window
}

// Destroy はOpenGLリソースを解放する
func (r *OpenGLRenderer) Destroy() {
	if r.bufferPool != nil {
		r.bufferPool.Destroy()
	}
	if r.shaderManager != nil {
		r.shaderManager.DeleteAllShaders()
	}
	if r.window != nil {
		r.window.Destroy()
		glfw.Terminate()
	}
}
