package core

import (
	"log"
	"github.com/ganyariya/tinyengine/pkg/tinyengine"
)

// Application は基本的なアプリケーション実装を提供する
type Application struct {
	// 将来的に追加する予定のフィールド
	// - シーンマネージャー
	// - 入力マネージャー
	// - オーディオマネージャー
}

// NewApplication は新しいアプリケーションインスタンスを作成する
func NewApplication() *Application {
	return &Application{}
}

// Initialize はアプリケーションを初期化する
func (app *Application) Initialize() error {
	log.Println("アプリケーションを初期化しています...")
	// TODO: システムの初期化
	// - レンダラーの初期化
	// - 入力システムの初期化
	// - オーディオシステムの初期化
	// - シーンの読み込み
	return nil
}

// Update はフレーム毎の更新処理を行う
func (app *Application) Update(deltaTime float64) {
	// TODO: システムの更新
	// - 入力の更新
	// - シーンの更新
	// - 物理演算
	// - 衝突判定
}

// Render は描画処理を行う
func (app *Application) Render(renderer tinyengine.Renderer) {
	// TODO: 描画処理
	// - 画面クリア
	// - シーンの描画
	// - UIの描画
	// - 画面表示
}

// Destroy はアプリケーションの終了処理を行う
func (app *Application) Destroy() {
	log.Println("アプリケーションを終了しています...")
	// TODO: システムの終了処理
	// - オーディオシステムの終了
	// - レンダラーの終了
	// - リソースの解放
}