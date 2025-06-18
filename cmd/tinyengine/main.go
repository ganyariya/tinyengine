package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ganyariya/tinyengine/internal/core"
	"github.com/ganyariya/tinyengine/internal/platform"
)

func main() {
	fmt.Println("TinyEngine - 教育的な小さなゲームエンジン")

	// エンジンの初期化
	engine := core.NewEngine("TinyEngine Demo", 800, 600)

	// 基本アプリケーションの作成
	app := core.NewApplication()
	engine.SetApplication(app)

	// ウィンドウの初期化テスト（実際の表示はまだ実装しない）
	windowConfig := platform.WindowConfig{
		Title:  "TinyEngine",
		Width:  800,
		Height: 600,
	}
	window := platform.NewWindow(windowConfig)

	// ヘッドレス環境での実行を考慮
	if err := window.Initialize(); err != nil {
		log.Printf("ウィンドウ初期化をスキップ（ヘッドレス環境）: %v", err)
	} else {
		defer window.Destroy()
		log.Println("ウィンドウ初期化成功")
	}

	// エンジンの実行（短時間で終了）
	log.Println("エンジンを開始...")

	time.Sleep(2 * time.Second) // 2秒間のスリープ

	// 簡単なデモ実行
	timer := platform.NewTimer()
	log.Printf("実行開始時刻: %.3f秒", timer.GetTime())

	log.Println("TinyEngine実行完了")
}
