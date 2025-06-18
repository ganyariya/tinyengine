package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/ganyariya/tinyengine/internal/renderer"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	runtime.LockOSThread()

	// GLFW初期化
	if err := glfw.Init(); err != nil {
		log.Fatalf("GLFWの初期化に失敗しました: %v", err)
	}
	defer glfw.Terminate()

	// OpenGLレンダラー作成
	openglRenderer, err := renderer.NewOpenGLRendererWithWindow(800, 600, "Simple Test")
	if err != nil {
		log.Fatalf("OpenGLレンダラーの作成に失敗しました: %v", err)
	}
	defer openglRenderer.(*renderer.OpenGLRenderer).Destroy()

	fmt.Println("Simple test - 中央に大きな赤い矩形を表示")

	// GLFWウィンドウの取得
	glRenderer := openglRenderer.(*renderer.OpenGLRenderer)
	window := glRenderer.GetWindow()

	for !window.ShouldClose() {
		// ESCキーで終了
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		// 画面クリア（黒背景）
		glRenderer.Clear()

		// 中央に大きな赤い矩形を描画（画面中央付近）
		// 座標: (200, 200) から (600, 400) の矩形
		glRenderer.DrawRectangleColor(200, 200, 400, 200, 1.0, 0.0, 0.0, 1.0)

		// 描画内容を画面に表示
		glRenderer.Present()
	}
}