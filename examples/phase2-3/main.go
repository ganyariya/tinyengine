package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/ganyariya/tinyengine/internal/renderer"
	"github.com/ganyariya/tinyengine/pkg/tinyengine"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "TinyEngine - Phase 2.3: Primitive Gallery"
)

func main() {
	runtime.LockOSThread()

	// GLFW初期化
	if err := glfw.Init(); err != nil {
		log.Fatalf("GLFWの初期化に失敗しました: %v", err)
	}
	defer glfw.Terminate()

	// OpenGLレンダラー作成
	openglRenderer, err := renderer.NewOpenGLRendererWithWindow(windowWidth, windowHeight, windowTitle)
	if err != nil {
		log.Fatalf("OpenGLレンダラーの作成に失敗しました: %v", err)
	}
	// 安全な型アサーション
	if oglRenderer, ok := openglRenderer.(*renderer.OpenGLRenderer); ok {
		defer oglRenderer.Destroy()
	}

	fmt.Println("=== TinyEngine Phase 2.3: 基本図形描画ギャラリー ===")
	fmt.Println("様々な図形とカラフルな色を表示します")
	fmt.Println("ESCキーで終了します")
	fmt.Println("")

	// 図形ギャラリーのメインループ
	runShapeGallery(openglRenderer)
}

func runShapeGallery(r tinyengine.Renderer) {
	// GLFWウィンドウの取得（安全な型アサーション）
	openglRenderer, ok := r.(*renderer.OpenGLRenderer)
	if !ok {
		log.Fatal("レンダラーがOpenGLRendererではありません")
	}
	window := openglRenderer.GetWindow()

	for !window.ShouldClose() {
		// ESCキーで終了
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		// 画面クリア（黒背景）
		r.Clear()

		// 様々な図形を描画
		drawShapeGallery(r)

		// 描画内容を画面に表示
		r.Present()
	}
}

func drawShapeGallery(r tinyengine.Renderer) {
	// 矩形のギャラリー
	drawRectangles(r)
	
	// 円のギャラリー
	drawCircles(r)
	
	// 線のギャラリー
	drawLines(r)
	
	// Primitiveオブジェクトを直接使用した描画
	drawPrimitives(r)
}

func drawRectangles(r tinyengine.Renderer) {
	// 色とりどりの矩形を描画 - より確実に見える位置に配置
	
	// 赤い矩形（左上）
	r.DrawRectangleColor(50, 50, 120, 80, 1.0, 0.0, 0.0, 1.0)
	
	// 緑の矩形（上中央）
	r.DrawRectangleColor(250, 50, 120, 80, 0.0, 1.0, 0.0, 1.0)
	
	// 青い矩形（右上）
	r.DrawRectangleColor(450, 50, 120, 80, 0.0, 0.0, 1.0, 1.0)
	
	// 黄色い矩形（左下）
	r.DrawRectangleColor(50, 150, 120, 80, 1.0, 1.0, 0.0, 1.0)
	
	// 紫の矩形（右下）
	r.DrawRectangleColor(450, 150, 120, 80, 1.0, 0.0, 1.0, 1.0)
}

func drawCircles(r tinyengine.Renderer) {
	// 色とりどりの円を描画
	
	// シアンの円
	r.DrawCircle(100, 250, 40, 0.0, 1.0, 1.0, 1.0)
	
	// オレンジの円
	r.DrawCircle(250, 250, 50, 1.0, 0.5, 0.0, 1.0)
	
	// ピンクの円
	r.DrawCircle(400, 250, 35, 1.0, 0.0, 0.5, 1.0)
	
	// ライムグリーンの円
	r.DrawCircle(550, 250, 45, 0.5, 1.0, 0.0, 1.0)
	
	// 半透明の水色の円
	r.DrawCircle(700, 250, 38, 0.0, 0.8, 1.0, 0.7)
}

func drawLines(r tinyengine.Renderer) {
	// 色とりどりの線を描画
	
	// 赤い線（水平）
	r.DrawLine(50, 400, 200, 400, 1.0, 0.0, 0.0, 1.0)
	
	// 緑の線（垂直）
	r.DrawLine(250, 380, 250, 440, 0.0, 1.0, 0.0, 1.0)
	
	// 青い線（斜め）
	r.DrawLine(300, 380, 400, 440, 0.0, 0.0, 1.0, 1.0)
	
	// 黄色い線（逆斜め）
	r.DrawLine(450, 440, 550, 380, 1.0, 1.0, 0.0, 1.0)
	
	// 紫の線（ジグザグの一部）
	r.DrawLine(600, 380, 650, 410, 1.0, 0.0, 1.0, 1.0)
	r.DrawLine(650, 410, 700, 380, 1.0, 0.0, 1.0, 1.0)
	r.DrawLine(700, 380, 750, 440, 1.0, 0.0, 1.0, 1.0)
}

func drawPrimitives(r tinyengine.Renderer) {
	// Primitiveオブジェクトを直接作成して描画
	
	// カスタム矩形（グラデーション風配置）
	rect1 := renderer.NewRectangle(100, 480, 60, 40, renderer.NewColorRGB(0.8, 0.2, 0.3))
	r.DrawPrimitive(rect1)
	
	rect2 := renderer.NewRectangle(170, 480, 60, 40, renderer.NewColorRGB(0.7, 0.4, 0.2))
	r.DrawPrimitive(rect2)
	
	rect3 := renderer.NewRectangle(240, 480, 60, 40, renderer.NewColorRGB(0.6, 0.6, 0.1))
	r.DrawPrimitive(rect3)
	
	// カスタム円（異なるセグメント数）
	circle1 := renderer.NewCircleWithSegments(400, 500, 25, renderer.NewColorRGB(0.2, 0.8, 0.6), 6)  // 六角形
	r.DrawPrimitive(circle1)
	
	circle2 := renderer.NewCircleWithSegments(500, 500, 30, renderer.NewColorRGB(0.9, 0.1, 0.7), 8)  // 八角形
	r.DrawPrimitive(circle2)
	
	circle3 := renderer.NewCircleWithSegments(600, 500, 20, renderer.NewColorRGB(0.1, 0.9, 0.9), 16) // なめらかな円
	r.DrawPrimitive(circle3)
	
	// カスタム線（アート風）
	line1 := renderer.NewLine(50, 550, 150, 580, renderer.NewColorRGB(1.0, 0.3, 0.1))
	r.DrawPrimitive(line1)
	
	line2 := renderer.NewLine(150, 580, 250, 550, renderer.NewColorRGB(0.3, 1.0, 0.1))
	r.DrawPrimitive(line2)
	
	line3 := renderer.NewLine(250, 550, 350, 580, renderer.NewColorRGB(0.1, 0.3, 1.0))
	r.DrawPrimitive(line3)
}