package core

import (
	"time"
	"github.com/ganyariya/tinyengine/pkg/tinyengine"
)

// Engine はゲームエンジンのコア機能を提供する
type Engine struct {
	title       string
	width       int
	height      int
	running     bool
	application tinyengine.GameObject
	lastTime    time.Time
}

// NewEngine は新しいエンジンインスタンスを作成する
func NewEngine(title string, width, height int) *Engine {
	return &Engine{
		title:  title,
		width:  width,
		height: height,
	}
}

// SetApplication はエンジンで実行するアプリケーションを設定する
func (e *Engine) SetApplication(app tinyengine.GameObject) {
	e.application = app
}

// Run はゲームループを開始する
func (e *Engine) Run() error {
	if e.application == nil {
		return ErrApplicationNotSet
	}

	// アプリケーションの初期化
	if err := e.application.Initialize(); err != nil {
		return NewEngineError("core", "application initialization", err)
	}

	e.running = true
	e.lastTime = time.Now()

	// ゲームループ
	for e.running {
		// デルタタイムの計算
		now := time.Now()
		deltaTime := now.Sub(e.lastTime).Seconds()
		e.lastTime = now

		// 更新処理
		e.application.Update(deltaTime)

		// 描画処理（レンダラーは後で実装）
		e.application.Render(nil)

		// フレームレート制限（60FPS）
		time.Sleep(DefaultFrameTimeMs)
	}

	// 終了処理
	e.application.Destroy()
	return nil
}

// Stop はゲームループを停止する
func (e *Engine) Stop() {
	e.running = false
}

// IsRunning はエンジンが動作中かを返す
func (e *Engine) IsRunning() bool {
	return e.running
}