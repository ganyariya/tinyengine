package platform

import (
	"fmt"
	"runtime"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// WindowConfig はウィンドウの設定を保持する
type WindowConfig struct {
	Title  string
	Width  int
	Height int
}

// Window はウィンドウ管理を行う
type Window struct {
	config      WindowConfig
	window      *glfw.Window
	initialized bool
}

// NewWindow は新しいウィンドウインスタンスを作成する
func NewWindow(config WindowConfig) *Window {
	return &Window{
		config: config,
	}
}

// Initialize はウィンドウを初期化する
func (w *Window) Initialize() error {
	// メインスレッドでGLFWを実行する必要がある
	runtime.LockOSThread()
	
	if err := w.initGLFW(); err != nil {
		return fmt.Errorf("GLFW initialization failed: %w", err)
	}
	
	if err := w.createWindow(); err != nil {
		glfw.Terminate()
		return fmt.Errorf("window creation failed: %w", err)
	}
	
	if err := w.initOpenGL(); err != nil {
		w.Destroy()
		return fmt.Errorf("OpenGL initialization failed: %w", err)
	}
	
	w.initialized = true
	return nil
}

// initGLFW initializes GLFW and sets hints
func (w *Window) initGLFW() error {
	if err := glfw.Init(); err != nil {
		return err
	}
	
	// OpenGLバージョン設定
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	
	return nil
}

// createWindow creates the GLFW window
func (w *Window) createWindow() error {
	window, err := glfw.CreateWindow(w.config.Width, w.config.Height, w.config.Title, nil, nil)
	if err != nil {
		return err
	}
	
	w.window = window
	w.window.MakeContextCurrent()
	return nil
}

// initOpenGL initializes OpenGL and sets up VSync
func (w *Window) initOpenGL() error {
	if err := gl.Init(); err != nil {
		return err
	}
	
	// VSync有効化
	glfw.SwapInterval(1)
	return nil
}

// ShouldClose はウィンドウが閉じられるべきかを返す
func (w *Window) ShouldClose() bool {
	if w.window == nil {
		return true
	}
	return w.window.ShouldClose()
}

// SwapBuffers はフロント・バックバッファを交換する
func (w *Window) SwapBuffers() {
	if w.window != nil {
		w.window.SwapBuffers()
	}
}

// PollEvents はイベントをポーリングする
func (w *Window) PollEvents() {
	glfw.PollEvents()
}

// GetSize はウィンドウサイズを返す
func (w *Window) GetSize() (int, int) {
	if w.window != nil {
		return w.window.GetSize()
	}
	return w.config.Width, w.config.Height
}

// Destroy はウィンドウを破棄する
func (w *Window) Destroy() {
	if w.window != nil {
		w.window.Destroy()
		w.window = nil
	}
	
	if w.initialized {
		glfw.Terminate()
		w.initialized = false
	}
}

// IsInitialized はウィンドウが初期化されているかを返す
func (w *Window) IsInitialized() bool {
	return w.initialized
}