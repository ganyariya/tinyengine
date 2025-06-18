package main

import (
	"log"
	"time"
	
	"github.com/ganyariya/tinyengine/internal/renderer"
)

func main() {
	log.Println("フェーズ2.1 ビジュアルサンプル: 基本レンダラーシステム")
	log.Println("黒い背景ウィンドウを3秒間表示します...")
	
	// OpenGLRendererでウィンドウ作成
	openglRenderer, err := renderer.NewOpenGLRendererWithWindow(800, 600, "TinyEngine Phase 2.1 - Basic Renderer")
	if err != nil {
		log.Fatalf("OpenGLRenderer作成に失敗しました: %v", err)
	}
	defer func() {
		if destroyer, ok := openglRenderer.(interface{ Destroy() }); ok {
			destroyer.Destroy()
		}
	}()
	
	log.Println("✅ ウィンドウが正常に作成されました")
	log.Println("📱 黒い背景のウィンドウが表示されていることを確認してください")
	
	// 3秒間表示
	startTime := time.Now()
	for time.Since(startTime) < 3*time.Second {
		// 画面をクリア（黒い背景）
		openglRenderer.Clear()
		
		// 画面に表示
		openglRenderer.Present()
	}
	
	log.Println("✅ フェーズ2.1のビジュアルサンプル完了")
	log.Println("")
	log.Println("確認項目:")
	log.Println("- [  ] 800x600のウィンドウが表示された")
	log.Println("- [  ] ウィンドウの背景が黒い色で表示された")
	log.Println("- [  ] ウィンドウタイトルが正しく表示された")
	log.Println("- [  ] 3秒後にウィンドウが正常に閉じた")
}