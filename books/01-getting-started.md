# フェーズ1: プロジェクト開始とGo環境セットアップ

## 学習目標
- Go言語でのプロジェクト構造を理解する
- モジュール管理とパッケージ設計を学ぶ
- OpenGL/GLFWの基本概念を把握する
- TDDの基本的な流れを実践する

## 理論: ゲームエンジンの基本構造

### ゲームエンジンとは何か？
ゲームエンジンは、ゲーム開発に必要な基本的な機能を提供するフレームワークです：

1. **描画システム**: 2D/3Dグラフィックスの表示
2. **入力システム**: キーボード、マウス、ゲームパッドの処理
3. **オーディオシステム**: 音楽、効果音の再生
4. **物理システム**: 衝突判定、物理演算
5. **リソース管理**: テクスチャ、サウンド、モデルの読み込み

### なぜGoでゲームエンジンを作るのか？
- **シンプルな言語仕様**: 学習コストが低い
- **強力な標準ライブラリ**: 豊富な機能がビルトイン
- **優れたテストサポート**: TDDがやりやすい
- **並行処理**: ゴルーチンによる効率的な並行処理
- **クロスプラットフォーム**: 複数OSで動作

## 実装手順

### ステップ1: プロジェクト初期化

```bash
# プロジェクトディレクトリ作成
mkdir tinyengine
cd tinyengine

# Goモジュール初期化
go mod init github.com/yourname/tinyengine

# 基本ディレクトリ構造作成
mkdir -p {cmd/tinyengine,internal/{core,renderer,input,audio,scene,ui,platform,math,collision,camera},pkg/tinyengine,examples/{basic,platformer,pong},assets/{shaders,textures,sounds,fonts},test/{fixtures,helpers},books}
```

### ステップ2: 依存関係の追加

```bash
# OpenGL関連ライブラリ
go get github.com/go-gl/gl/v4.1-core/gl
go get github.com/go-gl/glfw/v3.3/glfw

# 数学ライブラリ
go get github.com/go-gl/mathgl/mgl32

# テストライブラリ
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock

# オーディオライブラリ（後で使用）
go get github.com/faiface/beep
```

### ステップ3: 基本インターフェース設計

最初に、エンジンの核となるインターフェースを定義します：

```go
// pkg/tinyengine/interfaces.go
package tinyengine

// GameObject はゲーム内のすべてのオブジェクトが実装するインターフェース
type GameObject interface {
    Initialize() error
    Update(deltaTime float64)
    Render(renderer Renderer)
    Destroy()
}

// Renderer は描画を担当するインターフェース
type Renderer interface {
    Clear()
    Present()
    DrawRectangle(x, y, width, height float32)
}

// Application はアプリケーション全体を管理するインターフェース
type Application interface {
    Initialize() error
    Run() error
    Shutdown()
}
```

### ステップ4: TDDでのインターフェーステスト

```go
// pkg/tinyengine/interfaces_test.go
package tinyengine

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockGameObject はテスト用のGameObjectモック
type MockGameObject struct {
    mock.Mock
}

func (m *MockGameObject) Initialize() error {
    args := m.Called()
    return args.Error(0)
}

func (m *MockGameObject) Update(deltaTime float64) {
    m.Called(deltaTime)
}

func (m *MockGameObject) Render(renderer Renderer) {
    m.Called(renderer)
}

func (m *MockGameObject) Destroy() {
    m.Called()
}

func TestGameObject_Interface(t *testing.T) {
    // Arrange
    mockObj := new(MockGameObject)
    mockObj.On("Initialize").Return(nil)
    mockObj.On("Update", 0.016).Return()
    mockObj.On("Destroy").Return()
    
    // Act & Assert
    err := mockObj.Initialize()
    assert.NoError(t, err)
    
    mockObj.Update(0.016)
    mockObj.Destroy()
    
    mockObj.AssertExpectations(t)
}
```

### ステップ5: 基本的なウィンドウシステム

```go
// internal/platform/window.go
package platform

import (
    "fmt"
    "runtime"
    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
    window *glfw.Window
    width  int
    height int
    title  string
}

func NewWindow(width, height int, title string) (*Window, error) {
    runtime.LockOSThread()
    
    if err := glfw.Init(); err != nil {
        return nil, fmt.Errorf("failed to initialize GLFW: %v", err)
    }
    
    // OpenGL バージョン設定
    glfw.WindowHint(glfw.ContextVersionMajor, 4)
    glfw.WindowHint(glfw.ContextVersionMinor, 1)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
    
    window, err := glfw.CreateWindow(width, height, title, nil, nil)
    if err != nil {
        glfw.Terminate()
        return nil, fmt.Errorf("failed to create window: %v", err)
    }
    
    window.MakeContextCurrent()
    
    if err := gl.Init(); err != nil {
        return nil, fmt.Errorf("failed to initialize OpenGL: %v", err)
    }
    
    return &Window{
        window: window,
        width:  width,
        height: height,
        title:  title,
    }, nil
}

func (w *Window) ShouldClose() bool {
    return w.window.ShouldClose()
}

func (w *Window) SwapBuffers() {
    w.window.SwapBuffers()
}

func (w *Window) PollEvents() {
    glfw.PollEvents()
}

func (w *Window) Destroy() {
    w.window.Destroy()
    glfw.Terminate()
}
```

### ステップ6: 最初のサンプルアプリケーション

```go
// examples/basic/main.go
package main

import (
    "log"
    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/yourname/tinyengine/internal/platform"
)

func main() {
    // ウィンドウ作成
    window, err := platform.NewWindow(800, 600, "TinyEngine - 基本ウィンドウ")
    if err != nil {
        log.Fatalf("Failed to create window: %v", err)
    }
    defer window.Destroy()
    
    // メインループ
    for !window.ShouldClose() {
        // 画面をクリア（濃い青色）
        gl.ClearColor(0.2, 0.3, 0.8, 1.0)
        gl.Clear(gl.COLOR_BUFFER_BIT)
        
        // バッファ交換とイベント処理
        window.SwapBuffers()
        window.PollEvents()
    }
    
    log.Println("ウィンドウが正常に閉じられました")
}
```

## ビジュアル確認

このフェーズを完了すると、以下が実現できます：

### 期待される結果
- 濃い青色の800x600のウィンドウが表示される
- ウィンドウが正常に閉じることができる
- Goのモジュール構造が正しく設定されている

### 確認項目
- [ ] プロジェクト構造が正しく作成されている
- [ ] 依存関係が正しくインストールされている
- [ ] 基本ウィンドウが表示される
- [ ] ウィンドウが正常に閉じられる
- [ ] テストが実行される（`go test ./...`）

## トラブルシューティング

### よくある問題と解決法

1. **OpenGL関連エラー**
   ```bash
   # macOSの場合
   export CGO_CFLAGS_ALLOW="-Xpreprocessor"
   
   # Linuxの場合
   sudo apt-get install libgl1-mesa-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev
   ```

2. **モジュール解決エラー**
   ```bash
   go mod tidy
   go mod download
   ```

3. **テスト実行エラー**
   ```bash
   go clean -testcache
   go test -v ./...
   ```

## 次のステップ

フェーズ1を完了したら、次はフェーズ2.1（ゲームループとエンジンコア）に進みます。ここで学んだ基本的なプロジェクト構造とTDDのアプローチが、より複雑な機能の実装において重要な基盤となります。

## 理解度チェック

1. Goのインターフェースとは何か説明できますか？
2. TDDの「Red-Green-Refactor」サイクルを説明できますか？
3. OpenGLとGLFWの役割の違いを理解していますか？
4. モックオブジェクトがなぜテストで重要なのか説明できますか？