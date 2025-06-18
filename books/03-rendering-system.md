# フェーズ2.2: シェーダーシステム実装

## 学習目標
- OpenGLシェーダーの基本概念を理解する
- 依存性注入によるテスト可能な設計を学ぶ
- シェーダープログラムの管理方法を習得する
- 実際のGPUプログラミングの基礎を体験する

## 理論: シェーダーとGPUプログラミング

### シェーダーとは？
シェーダーは、GPU上で実行される小さなプログラムです：

1. **頂点シェーダー**: 頂点の位置や属性を処理
2. **フラグメントシェーダー**: ピクセルの色を決定
3. **ジオメトリシェーダー**: プリミティブの変形（今回は使用しない）

### GLSL（OpenGL Shading Language）
シェーダーはGLSLという言語で記述されます：

```glsl
// 頂点シェーダーの例
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 vertexColor;

void main() {
    gl_Position = vec4(aPos, 1.0);
    vertexColor = aColor;
}

// フラグメントシェーダーの例
#version 410 core
in vec3 vertexColor;
out vec4 FragColor;

uniform float alpha;

void main() {
    FragColor = vec4(vertexColor, alpha);
}
```

### 依存性注入の重要性
OpenGLはシステム依存のAPIのため、テストが困難です。依存性注入により：
- **テスト可能**: モックオブジェクトでテスト実行
- **保守性**: OpenGL実装を抽象化
- **移植性**: 異なるグラフィックAPIへの対応が容易

## 実装手順

### ステップ1: OpenGLバックエンド抽象化

```go
// internal/renderer/opengl_backend.go
package renderer

// OpenGLBackend はOpenGL APIの抽象化インターフェース
type OpenGLBackend interface {
    // シェーダー関連
    CreateShader(shaderType uint32) uint32
    ShaderSource(shader uint32, source string)
    CompileShader(shader uint32)
    GetShaderiv(shader uint32, pname uint32) int32
    GetShaderInfoLog(shader uint32) string
    DeleteShader(shader uint32)
    
    // プログラム関連
    CreateProgram() uint32
    AttachShader(program, shader uint32)
    DetachShader(program, shader uint32)
    LinkProgram(program uint32)
    GetProgramiv(program uint32, pname uint32) int32
    GetProgramInfoLog(program uint32) string
    UseProgram(program uint32)
    DeleteProgram(program uint32)
    
    // ユニフォーム変数関連
    GetUniformLocation(program uint32, name string) int32
    Uniform1f(location int32, value float32)
    Uniform1i(location int32, value int32)
    Uniform3fv(location int32, vector [3]float32)
    UniformMatrix4fv(location int32, matrix [16]float32)
}
```

### ステップ2: 実際のOpenGL実装

```go
// internal/renderer/real_opengl_backend.go
package renderer

import (
    "unsafe"
    "github.com/go-gl/gl/v4.1-core/gl"
)

// RealOpenGLBackend は実際のOpenGL APIを使用する実装
type RealOpenGLBackend struct{}

// NewRealOpenGLBackend は新しいRealOpenGLBackendを作成する
func NewRealOpenGLBackend() *RealOpenGLBackend {
    return &RealOpenGLBackend{}
}

func (r *RealOpenGLBackend) CreateShader(shaderType uint32) uint32 {
    return gl.CreateShader(shaderType)
}

func (r *RealOpenGLBackend) ShaderSource(shader uint32, source string) {
    cSource, free := gl.Strs(source + "\x00")
    defer free()
    gl.ShaderSource(shader, 1, cSource, nil)
}

func (r *RealOpenGLBackend) CompileShader(shader uint32) {
    gl.CompileShader(shader)
}

func (r *RealOpenGLBackend) GetShaderiv(shader uint32, pname uint32) int32 {
    var result int32
    gl.GetShaderiv(shader, pname, &result)
    return result
}

func (r *RealOpenGLBackend) GetShaderInfoLog(shader uint32) string {
    var logLength int32
    gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
    
    if logLength == 0 {
        return ""
    }
    
    log := make([]byte, logLength)
    gl.GetShaderInfoLog(shader, logLength, nil, &log[0])
    return string(log)
}

// 他のメソッドも同様に実装...
```

### ステップ3: テスト用モック実装

```go
// internal/renderer/mock_opengl_backend.go
package renderer

import (
    "strings"
    "github.com/stretchr/testify/mock"
)

// MockOpenGLBackend はテスト用のOpenGLバックエンドモック
type MockOpenGLBackend struct {
    mock.Mock
    
    // モック用の内部状態
    shaders       map[uint32]*MockShader
    programs      map[uint32]*MockProgram
    nextShaderID  uint32
    nextProgramID uint32
}

// MockShader はモック用のシェーダー情報
type MockShader struct {
    ID           uint32
    Type         uint32
    Source       string
    Compiled     bool
    CompileError string
}

// MockProgram はモック用のプログラム情報
type MockProgram struct {
    ID        uint32
    Shaders   []uint32
    Linked    bool
    LinkError string
    Uniforms  map[string]int32
    InUse     bool
}

// NewMockOpenGLBackend は新しいMockOpenGLBackendを作成する
func NewMockOpenGLBackend() *MockOpenGLBackend {
    return &MockOpenGLBackend{
        shaders:       make(map[uint32]*MockShader),
        programs:      make(map[uint32]*MockProgram),
        nextShaderID:  1,
        nextProgramID: 1,
    }
}

func (m *MockOpenGLBackend) CreateShader(shaderType uint32) uint32 {
    args := m.Called(shaderType)
    
    // モックの戻り値を取得
    id := args.Get(0).(uint32)
    
    // 内部状態にシェーダーを作成
    m.shaders[id] = &MockShader{
        ID:       id,
        Type:     shaderType,
        Compiled: false,
    }
    
    return id
}

func (m *MockOpenGLBackend) ShaderSource(shader uint32, source string) {
    m.Called(shader, source)
    
    if s, exists := m.shaders[shader]; exists {
        s.Source = source
    }
}

func (m *MockOpenGLBackend) CompileShader(shader uint32) {
    m.Called(shader)
    
    if s, exists := m.shaders[shader]; exists {
        // デフォルトではコンパイル成功
        s.Compiled = true
        
        // ソースコードに"ERROR"が含まれている場合はエラーにする
        if strings.Contains(s.Source, "ERROR") {
            s.Compiled = false
            s.CompileError = "Mock compile error"
        }
    }
}

// ヘルパーメソッド：テスト用
func (m *MockOpenGLBackend) GetShader(id uint32) *MockShader {
    return m.shaders[id]
}

func (m *MockOpenGLBackend) GetProgram(id uint32) *MockProgram {
    return m.programs[id]
}
```

### ステップ4: Shader構造体実装

```go
// internal/renderer/shader.go
package renderer

import (
    "fmt"
    "github.com/go-gl/gl/v4.1-core/gl"
)

// Shader はOpenGLシェーダープログラムを管理する
type Shader struct {
    backend          OpenGLBackend
    programID        uint32
    vertexShaderID   uint32
    fragmentShaderID uint32
}

// NewShader は新しいShaderを作成する
func NewShader(backend OpenGLBackend) *Shader {
    return &Shader{
        backend:          backend,
        programID:        0,
        vertexShaderID:   0,
        fragmentShaderID: 0,
    }
}

// LoadVertexShader は頂点シェーダーを読み込む
func (s *Shader) LoadVertexShader(source string) error {
    return s.loadShader(source, gl.VERTEX_SHADER, &s.vertexShaderID)
}

// LoadFragmentShader はフラグメントシェーダーを読み込む
func (s *Shader) LoadFragmentShader(source string) error {
    return s.loadShader(source, gl.FRAGMENT_SHADER, &s.fragmentShaderID)
}

// loadShader は指定された種類のシェーダーを読み込む
func (s *Shader) loadShader(source string, shaderType uint32, shaderID *uint32) error {
    // シェーダー作成
    *shaderID = s.backend.CreateShader(shaderType)
    if *shaderID == 0 {
        return fmt.Errorf("failed to create shader")
    }
    
    // ソースコード設定
    s.backend.ShaderSource(*shaderID, source)
    
    // コンパイル
    s.backend.CompileShader(*shaderID)
    
    // コンパイル結果確認
    success := s.backend.GetShaderiv(*shaderID, gl.COMPILE_STATUS)
    if success == gl.FALSE {
        log := s.backend.GetShaderInfoLog(*shaderID)
        s.backend.DeleteShader(*shaderID)
        *shaderID = 0
        return fmt.Errorf("shader compilation failed: %s", log)
    }
    
    return nil
}

// LinkProgram はシェーダープログラムをリンクする
func (s *Shader) LinkProgram() error {
    // 頂点・フラグメントシェーダーがロードされているかチェック
    if s.vertexShaderID == 0 {
        return fmt.Errorf("vertex shader not loaded")
    }
    if s.fragmentShaderID == 0 {
        return fmt.Errorf("fragment shader not loaded")
    }
    
    // プログラム作成
    s.programID = s.backend.CreateProgram()
    if s.programID == 0 {
        return fmt.Errorf("failed to create shader program")
    }
    
    // シェーダーをアタッチ
    s.backend.AttachShader(s.programID, s.vertexShaderID)
    s.backend.AttachShader(s.programID, s.fragmentShaderID)
    
    // リンク
    s.backend.LinkProgram(s.programID)
    
    // リンク結果確認
    success := s.backend.GetProgramiv(s.programID, gl.LINK_STATUS)
    if success == gl.FALSE {
        log := s.backend.GetProgramInfoLog(s.programID)
        return fmt.Errorf("shader program linking failed: %s", log)
    }
    
    // シェーダーをデタッチ・削除（プログラムにリンク済み）
    s.backend.DetachShader(s.programID, s.vertexShaderID)
    s.backend.DetachShader(s.programID, s.fragmentShaderID)
    s.backend.DeleteShader(s.vertexShaderID)
    s.backend.DeleteShader(s.fragmentShaderID)
    
    s.vertexShaderID = 0
    s.fragmentShaderID = 0
    
    return nil
}

// Use はシェーダープログラムを使用する
func (s *Shader) Use() {
    if s.programID != 0 {
        s.backend.UseProgram(s.programID)
    }
}

// GetUniformLocation はユニフォーム変数の位置を取得する
func (s *Shader) GetUniformLocation(name string) int32 {
    if s.programID == 0 {
        return -1
    }
    
    return s.backend.GetUniformLocation(s.programID, name)
}

// SetUniformFloat は浮動小数点数のユニフォーム変数を設定する
func (s *Shader) SetUniformFloat(location int32, value float32) {
    if location >= 0 {
        s.backend.Uniform1f(location, value)
    }
}
```

### ステップ5: 意味のあるテスト実装

```go
// internal/renderer/shader_test.go
package renderer

import (
    "testing"
    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/stretchr/testify/assert"
)

// テスト用の基本的なシェーダーソースコード
const (
    validVertexShaderSource = `
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 vertexColor;

void main() {
    gl_Position = vec4(aPos, 1.0);
    vertexColor = aColor;
}
`

    validFragmentShaderSource = `
#version 410 core
in vec3 vertexColor;
out vec4 FragColor;

uniform float alpha;

void main() {
    FragColor = vec4(vertexColor, alpha);
}
`

    invalidShaderSource = `
#version 410 core
ERROR This is invalid syntax
`
)

func TestNewShader_WithValidBackend(t *testing.T) {
    // Arrange
    mockBackend := NewMockOpenGLBackend()
    
    // Act
    shader := NewShader(mockBackend)
    
    // Assert
    assert.NotNil(t, shader)
    assert.Equal(t, uint32(0), shader.GetProgramID())
}

func TestShader_LoadVertexShader_Success(t *testing.T) {
    // Arrange
    mockBackend := NewMockOpenGLBackend()
    shader := NewShader(mockBackend)
    
    // モックの設定：正常ケース
    mockBackend.On("CreateShader", uint32(gl.VERTEX_SHADER)).Return(uint32(1))
    mockBackend.On("ShaderSource", uint32(1), validVertexShaderSource).Return()
    mockBackend.On("CompileShader", uint32(1)).Return()
    mockBackend.On("GetShaderiv", uint32(1), uint32(gl.COMPILE_STATUS)).Return(int32(1))
    
    // Act
    err := shader.LoadVertexShader(validVertexShaderSource)
    
    // Assert
    assert.NoError(t, err)
    mockBackend.AssertExpectations(t)
    
    // 内部状態の確認
    mockShader := mockBackend.GetShader(1)
    assert.NotNil(t, mockShader)
    assert.Equal(t, validVertexShaderSource, mockShader.Source)
    assert.True(t, mockShader.Compiled)
}

func TestShader_LoadVertexShader_CompilationError(t *testing.T) {
    // Arrange
    mockBackend := NewMockOpenGLBackend()
    shader := NewShader(mockBackend)
    
    // モックの設定：コンパイルエラーケース
    mockBackend.On("CreateShader", uint32(gl.VERTEX_SHADER)).Return(uint32(1))
    mockBackend.On("ShaderSource", uint32(1), invalidShaderSource).Return()
    mockBackend.On("CompileShader", uint32(1)).Return()
    mockBackend.On("GetShaderiv", uint32(1), uint32(gl.COMPILE_STATUS)).Return(int32(0))
    mockBackend.On("GetShaderInfoLog", uint32(1)).Return("Mock compile error")
    mockBackend.On("DeleteShader", uint32(1)).Return()
    
    // Act
    err := shader.LoadVertexShader(invalidShaderSource)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "shader compilation failed")
    assert.Contains(t, err.Error(), "Mock compile error")
    mockBackend.AssertExpectations(t)
}

func TestShader_FullWorkflow_Integration(t *testing.T) {
    // Arrange
    mockBackend := NewMockOpenGLBackend()
    shader := NewShader(mockBackend)
    
    // 全フローのモック設定
    // 頂点シェーダー
    mockBackend.On("CreateShader", uint32(gl.VERTEX_SHADER)).Return(uint32(1))
    mockBackend.On("ShaderSource", uint32(1), validVertexShaderSource).Return()
    mockBackend.On("CompileShader", uint32(1)).Return()
    mockBackend.On("GetShaderiv", uint32(1), uint32(gl.COMPILE_STATUS)).Return(int32(1))
    
    // フラグメントシェーダー
    mockBackend.On("CreateShader", uint32(gl.FRAGMENT_SHADER)).Return(uint32(2))
    mockBackend.On("ShaderSource", uint32(2), validFragmentShaderSource).Return()
    mockBackend.On("CompileShader", uint32(2)).Return()
    mockBackend.On("GetShaderiv", uint32(2), uint32(gl.COMPILE_STATUS)).Return(int32(1))
    
    // プログラムリンク
    mockBackend.On("CreateProgram").Return(uint32(3))
    mockBackend.On("AttachShader", uint32(3), uint32(1)).Return()
    mockBackend.On("AttachShader", uint32(3), uint32(2)).Return()
    mockBackend.On("LinkProgram", uint32(3)).Return()
    mockBackend.On("GetProgramiv", uint32(3), uint32(gl.LINK_STATUS)).Return(int32(1))
    mockBackend.On("DetachShader", uint32(3), uint32(1)).Return()
    mockBackend.On("DetachShader", uint32(3), uint32(2)).Return()
    mockBackend.On("DeleteShader", uint32(1)).Return()
    mockBackend.On("DeleteShader", uint32(2)).Return()
    
    // 使用とユニフォーム設定
    mockBackend.On("UseProgram", uint32(3)).Return()
    mockBackend.On("GetUniformLocation", uint32(3), "alpha").Return(int32(0))
    mockBackend.On("Uniform1f", int32(0), float32(0.5)).Return()
    
    // Act: 完全なワークフロー
    err := shader.LoadVertexShader(validVertexShaderSource)
    assert.NoError(t, err)
    
    err = shader.LoadFragmentShader(validFragmentShaderSource)
    assert.NoError(t, err)
    
    err = shader.LinkProgram()
    assert.NoError(t, err)
    
    shader.Use()
    
    location := shader.GetUniformLocation("alpha")
    assert.Equal(t, int32(0), location)
    
    shader.SetUniformFloat(location, 0.5)
    
    // Assert
    mockBackend.AssertExpectations(t)
    assert.Equal(t, uint32(3), shader.GetProgramID())
}
```

### ステップ6: ビジュアルサンプル実装

```go
// examples/phase2-2/main.go
package main

import (
    "log"
    "math"
    "runtime"
    "time"
    
    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
    "github.com/yourname/tinyengine/internal/renderer"
)

// 三角形の頂点データ（位置 + 色）
var triangleVertices = []float32{
    // 位置 (x, y, z)      色 (r, g, b)
    0.0, 0.5, 0.0,        1.0, 0.0, 0.0, // 上の頂点（赤）
    -0.5, -0.5, 0.0,      0.0, 1.0, 0.0, // 左下の頂点（緑）
    0.5, -0.5, 0.0,       0.0, 0.0, 1.0, // 右下の頂点（青）
}

func main() {
    log.Println("フェーズ2.2 ビジュアルサンプル: カラフルな三角形表示")
    log.Println("シェーダーシステムを使用したカラフルな三角形を表示します...")
    
    // OpenGL/GLFW初期化
    if err := initOpenGL(); err != nil {
        log.Fatalf("OpenGL初期化に失敗しました: %v", err)
    }
    defer glfw.Terminate()
    
    // ウィンドウ作成
    window, err := createWindow()
    if err != nil {
        log.Fatalf("ウィンドウ作成に失敗しました: %v", err)
    }
    defer window.Destroy()
    
    // OpenGL設定
    gl.Viewport(0, 0, 800, 600)
    
    // シェーダー作成
    shader, err := createColoredTriangleShader()
    if err != nil {
        log.Fatalf("シェーダー作成に失敗しました: %v", err)
    }
    defer shader.Delete()
    
    // VAO, VBO作成
    vao, vbo := createTriangleGeometry()
    defer deleteGeometry(vao, vbo)
    
    log.Println("✅ OpenGL初期化とシェーダー作成が完了しました")
    log.Println("📱 カラフルな三角形が表示されることを確認してください")
    log.Println("🎨 上が赤、左下が緑、右下が青のグラデーション三角形")
    
    // メインループ（5秒間表示）
    startTime := time.Now()
    for !window.ShouldClose() && time.Since(startTime) < 5*time.Second {
        // 画面クリア
        gl.ClearColor(0.1, 0.1, 0.1, 1.0)
        gl.Clear(gl.COLOR_BUFFER_BIT)
        
        // シェーダー使用
        shader.Use()
        
        // 時間ベースのアルファ値設定
        elapsed := float32(time.Since(startTime).Seconds())
        alpha := 0.5 + 0.5*float32(math.Sin(float64(elapsed*2.0)))
        alphaLocation := shader.GetUniformLocation("alpha")
        shader.SetUniformFloat(alphaLocation, alpha)
        
        // 三角形描画
        gl.BindVertexArray(vao)
        gl.DrawArrays(gl.TRIANGLES, 0, 3)
        gl.BindVertexArray(0)
        
        // バッファスワップとイベント処理
        window.SwapBuffers()
        glfw.PollEvents()
    }
    
    log.Println("✅ フェーズ2.2のビジュアルサンプル完了")
    log.Println("")
    log.Println("確認項目:")
    log.Println("- [  ] カラフルな三角形が表示された")
    log.Println("- [  ] 上の頂点が赤色で表示された")
    log.Println("- [  ] 左下の頂点が緑色で表示された")
    log.Println("- [  ] 右下の頂点が青色で表示された")
    log.Println("- [  ] 三角形にグラデーション効果があった")
    log.Println("- [  ] 時間と共に明度が変化していた")
}

func createColoredTriangleShader() (*renderer.Shader, error) {
    vertexShaderSource := `#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 vertexColor;

void main() {
    gl_Position = vec4(aPos, 1.0);
    vertexColor = aColor;
}
`

    fragmentShaderSource := `#version 410 core
in vec3 vertexColor;
out vec4 FragColor;

uniform float alpha;

void main() {
    FragColor = vec4(vertexColor, alpha);
}
`

    shader := renderer.NewShader(renderer.NewRealOpenGLBackend())
    
    if err := shader.LoadVertexShader(vertexShaderSource); err != nil {
        return nil, err
    }
    
    if err := shader.LoadFragmentShader(fragmentShaderSource); err != nil {
        return nil, err
    }
    
    if err := shader.LinkProgram(); err != nil {
        return nil, err
    }
    
    return shader, nil
}
```

## ビジュアル確認

このフェーズを完了すると、以下が実現できます：

### 期待される結果
- カラフルな三角形（上:赤、左下:緑、右下:青）が表示される
- グラデーション効果が美しく表現される
- 時間と共に明度が変化する（アルファ値のアニメーション）
- 5秒間表示された後、自動的に終了する

### 確認項目
- [ ] シェーダーが正常にコンパイル・リンクされる
- [ ] 頂点属性（位置・色）が正しく設定される
- [ ] ユニフォーム変数が正しく動作する
- [ ] 全テストが成功する（OpenGL環境なしでも実行可能）
- [ ] 依存性注入による設計が理解できる

## 重要な概念の理解

### なぜ依存性注入が重要なのか？
1. **テスト容易性**: OpenGL環境がなくてもテスト実行可能
2. **保守性**: OpenGL APIの変更に対応しやすい
3. **可読性**: 実際のロジックとOpenGL呼び出しを分離
4. **拡張性**: 他のグラフィックAPI（Vulkan、DirectX）への対応が容易

### シェーダープログラムのライフサイクル
1. **作成**: CreateShader() でシェーダーオブジェクト作成
2. **ソース設定**: ShaderSource() でGLSLコード設定
3. **コンパイル**: CompileShader() でGPUコードにコンパイル
4. **プログラム作成**: CreateProgram() でシェーダープログラム作成
5. **アタッチ**: AttachShader() でシェーダーをプログラムに結合
6. **リンク**: LinkProgram() で実行可能なプログラムを生成
7. **使用**: UseProgram() でGPUに送信
8. **削除**: DeleteShader(), DeleteProgram() でリソース解放

## 次のステップ

フェーズ2.2を完了したら、次はフェーズ2.3（基本図形描画）に進みます。ここで実装したシェーダーシステムを使用して、より複雑な図形や効果を描画していきます。

## 理解度チェック

1. 頂点シェーダーとフラグメントシェーダーの役割の違いを説明できますか？
2. ユニフォーム変数とは何で、なぜ重要なのか理解していますか？
3. 依存性注入がなぜテストにおいて重要なのか説明できますか？
4. シェーダープログラムのコンパイル→リンク→使用の流れを理解していますか？