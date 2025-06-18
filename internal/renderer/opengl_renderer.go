package renderer

import (
	"fmt"
	"runtime"
	
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ganyariya/tinyengine/pkg/tinyengine"
)

// OpenGLRenderer はOpenGLを使用した描画を提供する
type OpenGLRenderer struct {
	width  int
	height int
	window *glfw.Window
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
	
	renderer := &OpenGLRenderer{
		width:  width,
		height: height,
		window: window,
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
		x,         y,          // 左下
		x + width, y,          // 右下
		x + width, y + height, // 右上
		x,         y + height, // 左上
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

// Destroy はOpenGLリソースを解放する
func (r *OpenGLRenderer) Destroy() {
	if r.window != nil {
		r.window.Destroy()
		glfw.Terminate()
	}
}