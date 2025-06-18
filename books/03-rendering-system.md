# 第3章: 描画システム

## 理論編：描画システムとは何か？

### 3.1 描画システムの基本概念

描画システム（Rendering System）は、ゲームエンジンの中核となるコンポーネントの一つです。ゲーム内のオブジェクトやUIを画面に表示する責任を持ちます。

#### なぜ描画システムが必要なのか？

1. **抽象化**: 複雑なグラフィックスAPIを使いやすい形で提供
2. **最適化**: バッチング、カリング、Z-オーダリングなどの最適化を一箇所で管理
3. **ポータビリティ**: 異なるグラフィックスAPI（OpenGL、DirectX、Vulkan）に対応可能な統一インターフェース
4. **デバッグ**: 描画デバッグ、パフォーマンス計測の統一実装

### 3.2 一般的なゲームエンジンでの描画システム

主要なゲームエンジンでの描画システムの役割：

- **Unity**: Graphics.DrawMesh(), Camera.Render()
- **Unreal Engine**: UE Rendering Pipeline
- **Godot**: RenderingServer, CanvasItem

## 設計編：どう設計するか？

### 3.3 インターフェース設計

TinyEngineの描画システムは以下の設計原則に従います：

```go
// 基本的な描画インターフェース
type Renderer interface {
    Clear()                                         // 画面をクリア
    Present()                                       // 画面に表示
    DrawRectangle(x, y, width, height float32)      // 矩形描画
}
```

### 3.4 依存関係の整理

```
描画システムの依存関係:
GameObject -> Renderer Interface <- OpenGLRenderer
                |
                v
        CommandQueue System
                |
                v
         OpenGL/Graphics API
```

### 3.5 SOLID原則の適用

- **SRP**: Rendererは描画のみを担当
- **OCP**: 新しいレンダラー（DirectXRenderer等）を追加可能
- **LSP**: すべてのRenderer実装は互換性がある
- **ISP**: 必要最小限のメソッドのみ公開
- **DIP**: 実装ではなく抽象（interface）に依存

## 実装編：段階的実装手順

### ステップ1: 基本構造

#### 1.1 基本Rendererインターフェースの実装

```go
// pkg/tinyengine/interfaces.go
type Renderer interface {
    Clear()
    Present()
    DrawRectangle(x, y, width, height float32)
}
```

#### 1.2 BaseRendererの実装

```go
// internal/renderer/renderer.go
type BaseRenderer struct {
    width  int
    height int
}

func NewBaseRenderer(width, height int) tinyengine.Renderer {
    return &BaseRenderer{
        width:  width,
        height: height,
    }
}

func (r *BaseRenderer) Clear() {
    // 基本実装: 何もしない
}

func (r *BaseRenderer) Present() {
    // 基本実装: 何もしない
}

func (r *BaseRenderer) DrawRectangle(x, y, width, height float32) {
    // 基本実装: 何もしない
}
```

### ステップ2: OpenGLRenderer実装

#### 2.1 OpenGLライブラリの追加

```bash
go get github.com/go-gl/gl/v4.1-core/gl
go get github.com/go-gl/glfw/v3.3/glfw
```

#### 2.2 OpenGLRendererの実装

```go
// internal/renderer/opengl_renderer.go
type OpenGLRenderer struct {
    width  int
    height int
    window *glfw.Window
}

func NewOpenGLRenderer(width, height int) (tinyengine.Renderer, error) {
    renderer := &OpenGLRenderer{
        width:  width,
        height: height,
    }
    
    // ヘッドレス環境対応
    if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
        return nil, fmt.Errorf("OpenGL not available in headless environment")
    }
    
    return renderer, nil
}

func (r *OpenGLRenderer) Clear() {
    gl.ClearColor(0.0, 0.0, 0.0, 1.0)
    gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (r *OpenGLRenderer) Present() {
    if r.window != nil {
        r.window.SwapBuffers()
        glfw.PollEvents()
    }
}

func (r *OpenGLRenderer) DrawRectangle(x, y, width, height float32) {
    // 頂点データ作成
    vertices := []float32{
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
    
    gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
    gl.EnableVertexAttribArray(0)
    
    // 描画
    gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
    
    // クリーンアップ
    gl.BindVertexArray(0)
    gl.DeleteVertexArrays(1, &vao)
    gl.DeleteBuffers(1, &vbo)
    gl.DeleteBuffers(1, &ebo)
}
```

### ステップ3: コマンドキューシステム

#### 3.1 描画コマンドの定義

```go
// internal/renderer/command_queue.go
type CommandType int

const (
    ClearCommand CommandType = iota
    RectangleCommand
)

type RenderCommand struct {
    Type   CommandType
    Params map[string]interface{}
}
```

#### 3.2 CommandQueueの実装

```go
type CommandQueue struct {
    commands []RenderCommand
}

func NewCommandQueue() *CommandQueue {
    return &CommandQueue{
        commands: make([]RenderCommand, 0),
    }
}

func (q *CommandQueue) AddClearCommand() {
    command := RenderCommand{
        Type:   ClearCommand,
        Params: make(map[string]interface{}),
    }
    q.commands = append(q.commands, command)
}

func (q *CommandQueue) AddRectangleCommand(x, y, width, height float32) {
    command := RenderCommand{
        Type: RectangleCommand,
        Params: map[string]interface{}{
            "x":      x,
            "y":      y,
            "width":  width,
            "height": height,
        },
    }
    q.commands = append(q.commands, command)
}

func (q *CommandQueue) Execute(renderer tinyengine.Renderer) {
    for _, command := range q.commands {
        switch command.Type {
        case ClearCommand:
            renderer.Clear()
        case RectangleCommand:
            x := command.Params["x"].(float32)
            y := command.Params["y"].(float32)
            width := command.Params["width"].(float32)
            height := command.Params["height"].(float32)
            renderer.DrawRectangle(x, y, width, height)
        }
    }
}
```

## 動作原理編：どう動いているか？

### 4.1 OpenGL描画パイプライン

TinyEngineのOpenGL描画処理は以下の流れで実行されます：

1. **頂点データ作成**: 矩形の4つの頂点座標を生成
2. **バッファオブジェクト作成**: VBO（頂点バッファ）, EBO（インデックスバッファ）, VAO（頂点配列オブジェクト）
3. **データ転送**: CPUからGPUへ頂点データを転送
4. **描画実行**: GPU上でピクセルを描画
5. **リソース解放**: 一時的なOpenGLリソースをクリーンアップ

### 4.2 コマンドキューの利点

1. **遅延実行**: 描画コマンドを蓄積して最適なタイミングで実行
2. **バッチング**: 似たような描画コマンドをまとめて効率化
3. **デバッグ**: コマンド履歴を確認可能
4. **並列化**: コマンド生成と実行を並列化可能

### 4.3 メモリ配置とパフォーマンス

```
描画データのメモリフロー:
CPU Memory -> VBO (GPU Memory) -> Vertex Shader -> Fragment Shader -> Framebuffer
```

- **VBO**: 頂点データをGPUメモリに保存
- **EBO**: インデックスデータでメモリ使用量を削減
- **VAO**: 頂点属性の設定を保存して再利用

### 4.4 座標系変換

TinyEngineでは以下の座標系を使用：

- **ワールド座標**: ゲーム世界の絶対座標
- **スクリーン座標**: 実際の画面ピクセル座標
- **正規化デバイス座標(NDC)**: OpenGLの-1.0〜1.0の範囲

## テスト戦略

### 5.1 単体テスト

```go
func TestBaseRenderer_DrawRectangle(t *testing.T) {
    // Arrange
    renderer := NewBaseRenderer(800, 600)
    
    // Act & Assert
    assert.NotPanics(t, func() {
        renderer.DrawRectangle(10, 20, 100, 50)
    })
}
```

### 5.2 モックテスト

```go
type MockRenderer struct {
    mock.Mock
}

func (m *MockRenderer) DrawRectangle(x, y, width, height float32) {
    m.Called(x, y, width, height)
}

func TestCommandQueue_Execute(t *testing.T) {
    // Arrange
    queue := NewCommandQueue()
    mockRenderer := new(MockRenderer)
    mockRenderer.On("DrawRectangle", float32(10), float32(20), float32(100), float32(50)).Return()
    
    // Act
    queue.AddRectangleCommand(10, 20, 100, 50)
    queue.Execute(mockRenderer)
    
    // Assert
    mockRenderer.AssertExpectations(t)
}
```

### 5.3 統合テスト

```go
func TestOpenGLRenderer_Integration(t *testing.T) {
    // OpenGL環境が必要なため、CI環境ではスキップ
    if testing.Short() {
        t.Skip("Integration test requires OpenGL context")
    }
    
    renderer, err := NewOpenGLRendererWithWindow(800, 600, "Test")
    require.NoError(t, err)
    defer renderer.Destroy()
    
    // 基本的な描画テスト
    renderer.Clear()
    renderer.DrawRectangle(10, 20, 100, 50)
    renderer.Present()
}
```

## 発展編：さらに学ぶために

### 6.1 より高度な実装方法

1. **バッチレンダリング**: 複数の矩形を一度に描画
2. **テクスチャサポート**: 画像を使った描画
3. **シェーダーシステム**: カスタムシェーダーサポート
4. **3D描画**: 3Dオブジェクトの描画

### 6.2 他のゲームエンジンとの比較

| エンジン | 描画API | 特徴 |
|---------|---------|------|
| Unity | Graphics API | 高レベル抽象化、エディタ統合 |
| Unreal | RHI | 高パフォーマンス、リアルタイムレンダリング |
| Godot | RenderingServer | オープンソース、軽量 |
| TinyEngine | Renderer Interface | 教育目的、シンプル設計 |

### 6.3 参考資料

- [Learn OpenGL](https://learnopengl.com/)
- [OpenGL Tutorial](http://www.opengl-tutorial.org/)
- [Real-Time Rendering](https://www.realtimerendering.com/)
- [Game Engine Architecture](https://www.gameenginebook.com/)

## 章末課題

### 課題1: 実装確認

以下のコードを実行して、描画システムが正常に動作することを確認してください：

```go
func main() {
    // BaseRendererのテスト
    renderer := NewBaseRenderer(800, 600)
    renderer.Clear()
    renderer.DrawRectangle(100, 100, 200, 150)
    renderer.Present()
    
    // CommandQueueのテスト
    queue := NewCommandQueue()
    queue.AddClearCommand()
    queue.AddRectangleCommand(50, 50, 100, 100)
    queue.Execute(renderer)
}
```

### 課題2: 理解度チェック

1. Rendererインターフェースの設計原則は何ですか？
2. CommandQueueを使う利点を3つ挙げてください
3. OpenGLの描画パイプラインの主要なステップを説明してください

### 課題3: 拡張実装

以下の機能を追加実装してみてください：

1. 円形描画機能の추가
2. 色指定機能の추加
3. 描画統計情報（FPS、描画コマンド数）の収集

## まとめ

この章では、TinyEngineの描画システムを実装しました。学習したポイント：

- ✅ 描画システムの基本概念と必要性
- ✅ インターフェース設計とSOLID原則の適用
- ✅ OpenGLを使った実際の描画実装
- ✅ コマンドキューによる描画最適化
- ✅ テスト駆動開発による品質保証

次の章では、ゲーム数学システム（ベクトル、行列、変換）を実装していきます。