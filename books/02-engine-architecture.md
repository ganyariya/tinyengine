# フェーズ2.1: ゲームループとエンジンコア実装

## 学習目標
- ゲームループの仕組みを理解する
- デルタタイムとフレームレート管理を学ぶ
- エンジンアーキテクチャの設計パターンを習得する
- アプリケーションライフサイクル管理を実装する

## 理論: ゲームループの重要性

### ゲームループとは？
ゲームループは、ゲームが実行されている間に継続的に実行される処理サイクルです：

```
while (game is running) {
    1. 入力処理 (Handle Input)
    2. ゲーム状態更新 (Update Game State)
    3. 描画 (Render)
    4. タイミング調整 (Frame Rate Control)
}
```

### デルタタイムの重要性
デルタタイムは前回のフレームからの経過時間で、以下の利点があります：
- **フレームレート非依存**: 異なる性能のマシンで一貫した動作
- **スムーズな動き**: 可変フレームレートに対応
- **物理演算の精度**: 時間ベースの計算が可能

### アプリケーションライフサイクル
1. **初期化**: リソースの読み込み、システムの準備
2. **実行**: メインゲームループの実行
3. **終了**: リソースの解放、状態の保存

## 実装手順

### ステップ1: タイマーシステム実装

```go
// internal/platform/timer.go
package platform

import (
    "time"
)

// Timer はゲーム用の高精度タイマーを提供する
type Timer struct {
    startTime    time.Time
    lastTime     time.Time
    deltaTime    float64
    totalTime    float64
    frameCount   int
    fps          float64
    fpsUpdateTime float64
}

// NewTimer は新しいTimerを作成する
func NewTimer() *Timer {
    now := time.Now()
    return &Timer{
        startTime:    now,
        lastTime:     now,
        deltaTime:    0.0,
        totalTime:    0.0,
        frameCount:   0,
        fps:          0.0,
        fpsUpdateTime: 0.0,
    }
}

// Update はタイマーを更新し、デルタタイムを計算する
func (t *Timer) Update() {
    now := time.Now()
    t.deltaTime = now.Sub(t.lastTime).Seconds()
    t.totalTime = now.Sub(t.startTime).Seconds()
    t.lastTime = now
    
    // FPS計算（1秒間隔で更新）
    t.frameCount++
    t.fpsUpdateTime += t.deltaTime
    if t.fpsUpdateTime >= 1.0 {
        t.fps = float64(t.frameCount) / t.fpsUpdateTime
        t.frameCount = 0
        t.fpsUpdateTime = 0.0
    }
}

// GetDeltaTime は前回フレームからの経過時間を取得する
func (t *Timer) GetDeltaTime() float64 {
    return t.deltaTime
}

// GetTotalTime はタイマー開始からの総経過時間を取得する
func (t *Timer) GetTotalTime() float64 {
    return t.totalTime
}

// GetFPS は現在のフレームレートを取得する
func (t *Timer) GetFPS() float64 {
    return t.fps
}
```

### ステップ2: ゲームループ実装

```go
// internal/core/game_loop.go
package core

import (
    "fmt"
    "log"
    "github.com/yourname/tinyengine/internal/platform"
    "github.com/yourname/tinyengine/pkg/tinyengine"
)

// GameLoop はメインゲームループを管理する
type GameLoop struct {
    window   *platform.Window
    timer    *platform.Timer
    renderer tinyengine.Renderer
    running  bool
    
    // FPS制限
    targetFPS      float64
    frameTime      float64
    maxDeltaTime   float64
}

// NewGameLoop は新しいGameLoopを作成する
func NewGameLoop(window *platform.Window, renderer tinyengine.Renderer) *GameLoop {
    return &GameLoop{
        window:      window,
        timer:       platform.NewTimer(),
        renderer:    renderer,
        running:     false,
        targetFPS:   60.0,
        frameTime:   1.0 / 60.0,
        maxDeltaTime: 1.0 / 15.0, // 15FPS以下でもゲームロジックを制限
    }
}

// SetTargetFPS は目標フレームレートを設定する
func (gl *GameLoop) SetTargetFPS(fps float64) {
    gl.targetFPS = fps
    gl.frameTime = 1.0 / fps
}

// Start はゲームループを開始する
func (gl *GameLoop) Start() {
    gl.running = true
    log.Printf("ゲームループ開始 (目標FPS: %.1f)", gl.targetFPS)
}

// Stop はゲームループを停止する
func (gl *GameLoop) Stop() {
    gl.running = false
    log.Println("ゲームループ停止")
}

// IsRunning はゲームループが実行中かを確認する
func (gl *GameLoop) IsRunning() bool {
    return gl.running && !gl.window.ShouldClose()
}

// Update は1フレーム分のゲームループを実行する
func (gl *GameLoop) Update() error {
    if !gl.IsRunning() {
        return fmt.Errorf("game loop is not running")
    }
    
    // タイマー更新
    gl.timer.Update()
    deltaTime := gl.timer.GetDeltaTime()
    
    // デルタタイムを制限（スパイク防止）
    if deltaTime > gl.maxDeltaTime {
        deltaTime = gl.maxDeltaTime
    }
    
    // 入力処理（後で実装）
    gl.handleInput()
    
    // ゲーム状態更新（後で実装）
    gl.updateGameState(deltaTime)
    
    // 描画
    gl.render()
    
    // フレームレート制御（後で実装）
    gl.controlFrameRate()
    
    return nil
}

// handleInput は入力を処理する（フェーズ3で詳細実装）
func (gl *GameLoop) handleInput() {
    gl.window.PollEvents()
}

// updateGameState はゲーム状態を更新する
func (gl *GameLoop) updateGameState(deltaTime float64) {
    // このフェーズでは基本的な更新のみ
    // 実際のゲームオブジェクトの更新は後のフェーズで実装
}

// render は描画を実行する
func (gl *GameLoop) render() {
    gl.renderer.Clear()
    
    // 基本的な描画（テスト用）
    gl.renderer.DrawRectangle(100, 100, 50, 50)
    
    gl.renderer.Present()
    gl.window.SwapBuffers()
}

// controlFrameRate はフレームレートを制御する
func (gl *GameLoop) controlFrameRate() {
    // 現在は基本実装のみ
    // より高度なフレームレート制御は後で実装
}

// GetStats はゲームループの統計情報を取得する
func (gl *GameLoop) GetStats() (float64, float64, float64) {
    return gl.timer.GetFPS(), gl.timer.GetDeltaTime(), gl.timer.GetTotalTime()
}
```

### ステップ3: エンジンコア実装

```go
// internal/core/engine.go
package core

import (
    "fmt"
    "log"
    "github.com/yourname/tinyengine/internal/platform"
    "github.com/yourname/tinyengine/internal/renderer"
    "github.com/yourname/tinyengine/pkg/tinyengine"
)

// Engine はゲームエンジンのメインコントローラー
type Engine struct {
    window   *platform.Window
    renderer tinyengine.Renderer
    gameLoop *GameLoop
    
    // エンジン設定
    config *Config
}

// Config はエンジンの設定を保持する
type Config struct {
    WindowWidth  int
    WindowHeight int
    WindowTitle  string
    TargetFPS    float64
    VSync        bool
}

// DefaultConfig はデフォルトのエンジン設定を返す
func DefaultConfig() *Config {
    return &Config{
        WindowWidth:  800,
        WindowHeight: 600,
        WindowTitle:  "TinyEngine",
        TargetFPS:    60.0,
        VSync:        true,
    }
}

// NewEngine は新しいEngineを作成する
func NewEngine(config *Config) (*Engine, error) {
    if config == nil {
        config = DefaultConfig()
    }
    
    engine := &Engine{
        config: config,
    }
    
    return engine, nil
}

// Initialize はエンジンを初期化する
func (e *Engine) Initialize() error {
    log.Println("エンジン初期化開始...")
    
    // ウィンドウ作成
    window, err := platform.NewWindow(
        e.config.WindowWidth,
        e.config.WindowHeight,
        e.config.WindowTitle,
    )
    if err != nil {
        return fmt.Errorf("failed to create window: %v", err)
    }
    e.window = window
    
    // レンダラー作成
    renderer, err := renderer.NewOpenGLRenderer(
        e.config.WindowWidth,
        e.config.WindowHeight,
    )
    if err != nil {
        return fmt.Errorf("failed to create renderer: %v", err)
    }
    e.renderer = renderer
    
    // ゲームループ作成
    e.gameLoop = NewGameLoop(e.window, e.renderer)
    e.gameLoop.SetTargetFPS(e.config.TargetFPS)
    
    log.Println("エンジン初期化完了")
    return nil
}

// Run はエンジンを実行する
func (e *Engine) Run() error {
    if e.gameLoop == nil {
        return fmt.Errorf("engine not initialized")
    }
    
    log.Println("エンジン実行開始")
    e.gameLoop.Start()
    
    // メインループ
    for e.gameLoop.IsRunning() {
        if err := e.gameLoop.Update(); err != nil {
            return fmt.Errorf("game loop error: %v", err)
        }
        
        // 統計情報の表示（デバッグ用）
        if fps, deltaTime, totalTime := e.gameLoop.GetStats(); totalTime > 0 {
            if int(totalTime)%5 == 0 && int(totalTime*10)%10 == 0 {
                log.Printf("FPS: %.1f, DeltaTime: %.4f, TotalTime: %.1f", 
                    fps, deltaTime, totalTime)
            }
        }
    }
    
    log.Println("エンジン実行終了")
    return nil
}

// Shutdown はエンジンを終了する
func (e *Engine) Shutdown() {
    log.Println("エンジン終了処理開始...")
    
    if e.gameLoop != nil {
        e.gameLoop.Stop()
    }
    
    if e.renderer != nil {
        // レンダラーの終了処理
    }
    
    if e.window != nil {
        e.window.Destroy()
    }
    
    log.Println("エンジン終了処理完了")
}
```

### ステップ4: アプリケーション実装

```go
// internal/core/application.go
package core

import (
    "log"
    "github.com/yourname/tinyengine/pkg/tinyengine"
)

// Application はアプリケーション全体を管理する
type Application struct {
    engine   *Engine
    config   *Config
    initialized bool
}

// NewApplication は新しいApplicationを作成する
func NewApplication(config *Config) *Application {
    return &Application{
        config: config,
        initialized: false,
    }
}

// Initialize はアプリケーションを初期化する
func (app *Application) Initialize() error {
    if app.initialized {
        return nil
    }
    
    log.Println("アプリケーション初期化開始...")
    
    // エンジン作成と初期化
    engine, err := NewEngine(app.config)
    if err != nil {
        return err
    }
    
    if err := engine.Initialize(); err != nil {
        return err
    }
    
    app.engine = engine
    app.initialized = true
    
    log.Println("アプリケーション初期化完了")
    return nil
}

// Run はアプリケーションを実行する
func (app *Application) Run() error {
    if !app.initialized {
        if err := app.Initialize(); err != nil {
            return err
        }
    }
    
    return app.engine.Run()
}

// Shutdown はアプリケーションを終了する
func (app *Application) Shutdown() {
    if app.engine != nil {
        app.engine.Shutdown()
    }
    app.initialized = false
}
```

### ステップ5: ビジュアルサンプル実装

```go
// examples/phase2-1/main.go
package main

import (
    "log"
    "github.com/yourname/tinyengine/internal/core"
)

func main() {
    log.Println("フェーズ2.1 ビジュアルサンプル: 基本ゲームループ")
    log.Println("動く四角形でゲームループとデルタタイムを確認します...")
    
    // アプリケーション設定
    config := &core.Config{
        WindowWidth:  800,
        WindowHeight: 600,
        WindowTitle:  "TinyEngine Phase 2.1 - Game Loop Demo",
        TargetFPS:    60.0,
        VSync:        true,
    }
    
    // アプリケーション作成
    app := core.NewApplication(config)
    
    // 実行
    if err := app.Run(); err != nil {
        log.Fatalf("アプリケーション実行エラー: %v", err)
    }
    
    // 終了処理
    app.Shutdown()
    
    log.Println("✅ フェーズ2.1のビジュアルサンプル完了")
    log.Println("")
    log.Println("確認項目:")
    log.Println("- [  ] ウィンドウが表示された")
    log.Println("- [  ] 四角形が表示された")
    log.Println("- [  ] FPS情報がコンソールに表示された")
    log.Println("- [  ] ウィンドウを閉じることができた")
}
```

## TDDテスト実装

### ステップ6: 包括的なテスト作成

```go
// internal/core/game_loop_test.go
package core

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockWindow はテスト用のWindowモック
type MockWindow struct {
    mock.Mock
    shouldClose bool
}

func (m *MockWindow) ShouldClose() bool {
    return m.shouldClose
}

func (m *MockWindow) SwapBuffers() {
    m.Called()
}

func (m *MockWindow) PollEvents() {
    m.Called()
}

func TestGameLoop_Creation(t *testing.T) {
    // Arrange
    mockWindow := new(MockWindow)
    mockRenderer := new(MockRenderer)
    
    // Act
    gameLoop := NewGameLoop(mockWindow, mockRenderer)
    
    // Assert
    assert.NotNil(t, gameLoop)
    assert.False(t, gameLoop.IsRunning())
}

func TestGameLoop_StartStop(t *testing.T) {
    // Arrange
    mockWindow := new(MockWindow)
    mockRenderer := new(MockRenderer)
    gameLoop := NewGameLoop(mockWindow, mockRenderer)
    
    // Act & Assert
    assert.False(t, gameLoop.IsRunning())
    
    gameLoop.Start()
    assert.True(t, gameLoop.running)
    
    gameLoop.Stop()
    assert.False(t, gameLoop.running)
}

func TestGameLoop_Update(t *testing.T) {
    // Arrange
    mockWindow := new(MockWindow)
    mockRenderer := new(MockRenderer)
    
    mockWindow.On("PollEvents").Return()
    mockWindow.On("SwapBuffers").Return()
    mockRenderer.On("Clear").Return()
    mockRenderer.On("DrawRectangle", float32(100), float32(100), float32(50), float32(50)).Return()
    mockRenderer.On("Present").Return()
    
    gameLoop := NewGameLoop(mockWindow, mockRenderer)
    gameLoop.Start()
    
    // Act
    err := gameLoop.Update()
    
    // Assert
    assert.NoError(t, err)
    mockWindow.AssertExpectations(t)
    mockRenderer.AssertExpectations(t)
}
```

## ビジュアル確認

このフェーズを完了すると、以下が実現できます：

### 期待される結果
- ウィンドウに四角形が表示される
- コンソールにFPS情報が5秒間隔で表示される
- スムーズなゲームループが動作する
- ウィンドウを閉じると正常に終了する

### 確認項目
- [ ] ゲームループが正常に動作する
- [ ] デルタタイムが正しく計算される
- [ ] FPS表示が機能する
- [ ] アプリケーションライフサイクルが正常に動作する
- [ ] 全テストが成功する

## パフォーマンス考慮事項

### フレームレート最適化
- **VSync**: 画面リフレッシュレートに同期
- **デルタタイム制限**: スパイクによる異常動作防止
- **ガベージコレクション**: アロケーションの最小化

### メモリ管理
- **オブジェクトプール**: 頻繁に作成されるオブジェクトの再利用
- **リソース管理**: 適切な初期化と解放

## 次のステップ

フェーズ2.1を完了したら、次はフェーズ2.2（シェーダーシステム）に進みます。ここで実装したゲームループが、より複雑な描画システムの基盤となります。

## 理解度チェック

1. デルタタイムがなぜフレームレート非依存のゲームに重要なのか説明できますか？
2. ゲームループの各ステップ（入力処理、更新、描画）の役割を理解していますか？
3. アプリケーションライフサイクルの重要性を説明できますか？
4. FPS制限がなぜ必要なのか理解していますか？