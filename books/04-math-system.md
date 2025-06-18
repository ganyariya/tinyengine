# フェーズ3: ゲーム数学システム実装

## 学習目標
- ゲーム開発に必要な数学の基礎を理解する
- ベクトル演算と行列変換を実装する
- 2D座標系とトランスフォーム操作を習得する
- 効率的な数学ライブラリの設計を学ぶ

## 理論: ゲーム数学の基礎

### なぜゲーム数学が重要なのか？
ゲーム開発では以下の場面で数学が不可欠です：

1. **オブジェクトの位置・回転・スケール**: Transform行列
2. **衝突判定**: ベクトルの内積・外積
3. **カメラ制御**: ビュー行列・プロジェクション行列
4. **アニメーション**: 補間計算
5. **物理演算**: 力学・運動方程式

### 2Dゲームで使用する主な概念

#### ベクトル (Vector)
```
Vector2: (x, y)
- 位置: オブジェクトの座標
- 速度: 移動方向と速さ
- 方向: 正規化されたベクトル
```

#### 行列 (Matrix)
```
Matrix3x3: 2D変換用
[sx*cos(θ)  -sy*sin(θ)  tx]
[sx*sin(θ)   sy*cos(θ)  ty]
[    0           0       1 ]

sx, sy: スケール
θ: 回転角
tx, ty: 平行移動
```

#### 変換の順序
```
最終変換 = 平行移動 × 回転 × スケール × 頂点座標
(順序が重要！)
```

## 実装手順

### ステップ1: Vector2実装

```go
// internal/math/vector2.go
package math

import (
    "math"
)

// Vector2 は2次元ベクトルを表す
type Vector2 struct {
    X, Y float32
}

// NewVector2 は新しいVector2を作成する
func NewVector2(x, y float32) Vector2 {
    return Vector2{X: x, Y: y}
}

// Zero はゼロベクトルを返す
func Zero() Vector2 {
    return Vector2{X: 0, Y: 0}
}

// One は単位ベクトル(1,1)を返す
func One() Vector2 {
    return Vector2{X: 1, Y: 1}
}

// Up は上方向ベクトル(0,1)を返す
func Up() Vector2 {
    return Vector2{X: 0, Y: 1}
}

// Right は右方向ベクトル(1,0)を返す
func Right() Vector2 {
    return Vector2{X: 1, Y: 0}
}

// Add はベクトルの加算を行う
func (v Vector2) Add(other Vector2) Vector2 {
    return Vector2{
        X: v.X + other.X,
        Y: v.Y + other.Y,
    }
}

// Sub はベクトルの減算を行う
func (v Vector2) Sub(other Vector2) Vector2 {
    return Vector2{
        X: v.X - other.X,
        Y: v.Y - other.Y,
    }
}

// Mul はスカラー倍を行う
func (v Vector2) Mul(scalar float32) Vector2 {
    return Vector2{
        X: v.X * scalar,
        Y: v.Y * scalar,
    }
}

// Div はスカラー除算を行う
func (v Vector2) Div(scalar float32) Vector2 {
    if scalar == 0 {
        return v // ゼロ除算回避
    }
    return Vector2{
        X: v.X / scalar,
        Y: v.Y / scalar,
    }
}

// Dot は内積を計算する
func (v Vector2) Dot(other Vector2) float32 {
    return v.X*other.X + v.Y*other.Y
}

// Cross は外積のZ成分を計算する（2Dでは擬似外積）
func (v Vector2) Cross(other Vector2) float32 {
    return v.X*other.Y - v.Y*other.X
}

// Length はベクトルの長さを返す
func (v Vector2) Length() float32 {
    return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// LengthSquared はベクトルの長さの二乗を返す（高速）
func (v Vector2) LengthSquared() float32 {
    return v.X*v.X + v.Y*v.Y
}

// Normalize はベクトルを正規化する
func (v Vector2) Normalize() Vector2 {
    length := v.Length()
    if length == 0 {
        return Zero() // ゼロベクトルは正規化できない
    }
    return v.Div(length)
}

// Distance は2点間の距離を計算する
func (v Vector2) Distance(other Vector2) float32 {
    return v.Sub(other).Length()
}

// DistanceSquared は2点間の距離の二乗を計算する（高速）
func (v Vector2) DistanceSquared(other Vector2) float32 {
    return v.Sub(other).LengthSquared()
}

// Lerp は線形補間を行う
func (v Vector2) Lerp(other Vector2, t float32) Vector2 {
    // t=0でv, t=1でotherを返す
    return Vector2{
        X: v.X + (other.X-v.X)*t,
        Y: v.Y + (other.Y-v.Y)*t,
    }
}

// Angle は2つのベクトル間の角度を返す（ラジアン）
func (v Vector2) Angle(other Vector2) float32 {
    dot := v.Dot(other)
    lengths := v.Length() * other.Length()
    if lengths == 0 {
        return 0
    }
    return float32(math.Acos(float64(dot / lengths)))
}
```

### ステップ2: Matrix3実装

```go
// internal/math/matrix3.go
package math

import (
    "math"
)

// Matrix3 は3x3行列を表す（2D変換用）
type Matrix3 struct {
    // 行優先で格納
    // [m00 m01 m02]
    // [m10 m11 m12]
    // [m20 m21 m22]
    M [9]float32
}

// NewMatrix3 は新しいMatrix3を作成する
func NewMatrix3() Matrix3 {
    return Matrix3{}
}

// Identity は単位行列を返す
func Identity() Matrix3 {
    return Matrix3{
        M: [9]float32{
            1, 0, 0,
            0, 1, 0,
            0, 0, 1,
        },
    }
}

// Translation は平行移動行列を作成する
func Translation(x, y float32) Matrix3 {
    return Matrix3{
        M: [9]float32{
            1, 0, x,
            0, 1, y,
            0, 0, 1,
        },
    }
}

// Rotation は回転行列を作成する（ラジアン）
func Rotation(angle float32) Matrix3 {
    cos := float32(math.Cos(float64(angle)))
    sin := float32(math.Sin(float64(angle)))
    
    return Matrix3{
        M: [9]float32{
            cos, -sin, 0,
            sin,  cos, 0,
            0,    0,   1,
        },
    }
}

// Scale はスケール行列を作成する
func Scale(sx, sy float32) Matrix3 {
    return Matrix3{
        M: [9]float32{
            sx, 0,  0,
            0,  sy, 0,
            0,  0,  1,
        },
    }
}

// Multiply は行列の乗算を行う
func (m Matrix3) Multiply(other Matrix3) Matrix3 {
    result := Matrix3{}
    
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            result.M[i*3+j] = 0
            for k := 0; k < 3; k++ {
                result.M[i*3+j] += m.M[i*3+k] * other.M[k*3+j]
            }
        }
    }
    
    return result
}

// Transform はベクトルを変換する
func (m Matrix3) Transform(v Vector2) Vector2 {
    x := m.M[0]*v.X + m.M[1]*v.Y + m.M[2]
    y := m.M[3]*v.X + m.M[4]*v.Y + m.M[5]
    return Vector2{X: x, Y: y}
}

// TransformDirection は方向ベクトルを変換する（平行移動を無視）
func (m Matrix3) TransformDirection(v Vector2) Vector2 {
    x := m.M[0]*v.X + m.M[1]*v.Y
    y := m.M[3]*v.X + m.M[4]*v.Y
    return Vector2{X: x, Y: y}
}
```

### ステップ3: Transform構造体実装

```go
// internal/math/transform.go
package math

// Transform は2Dオブジェクトの変換情報を管理する
type Transform struct {
    Position Vector2
    Rotation float32 // ラジアン
    Scale    Vector2
    
    // キャッシュ
    matrix      Matrix3
    matrixDirty bool
}

// NewTransform は新しいTransformを作成する
func NewTransform() *Transform {
    return &Transform{
        Position:    Zero(),
        Rotation:    0,
        Scale:       One(),
        matrixDirty: true,
    }
}

// SetPosition は位置を設定する
func (t *Transform) SetPosition(pos Vector2) {
    t.Position = pos
    t.matrixDirty = true
}

// SetRotation は回転を設定する（度数法）
func (t *Transform) SetRotation(degrees float32) {
    t.Rotation = degrees * (math.Pi / 180.0)
    t.matrixDirty = true
}

// SetRotationRadians は回転を設定する（ラジアン）
func (t *Transform) SetRotationRadians(radians float32) {
    t.Rotation = radians
    t.matrixDirty = true
}

// SetScale はスケールを設定する
func (t *Transform) SetScale(scale Vector2) {
    t.Scale = scale
    t.matrixDirty = true
}

// SetUniformScale は均等スケールを設定する
func (t *Transform) SetUniformScale(scale float32) {
    t.Scale = Vector2{X: scale, Y: scale}
    t.matrixDirty = true
}

// GetMatrix は変換行列を取得する
func (t *Transform) GetMatrix() Matrix3 {
    if t.matrixDirty {
        t.updateMatrix()
    }
    return t.matrix
}

// updateMatrix は変換行列を更新する
func (t *Transform) updateMatrix() {
    // S * R * T の順序で計算
    scaleMatrix := Scale(t.Scale.X, t.Scale.Y)
    rotationMatrix := Rotation(t.Rotation)
    translationMatrix := Translation(t.Position.X, t.Position.Y)
    
    // 行列の合成: T * R * S
    t.matrix = translationMatrix.Multiply(rotationMatrix).Multiply(scaleMatrix)
    t.matrixDirty = false
}

// TransformPoint は点を変換する
func (t *Transform) TransformPoint(point Vector2) Vector2 {
    return t.GetMatrix().Transform(point)
}

// TransformDirection は方向ベクトルを変換する
func (t *Transform) TransformDirection(direction Vector2) Vector2 {
    return t.GetMatrix().TransformDirection(direction)
}

// GetForward は前方向ベクトルを取得する
func (t *Transform) GetForward() Vector2 {
    return t.TransformDirection(Up())
}

// GetRight は右方向ベクトルを取得する
func (t *Transform) GetRight() Vector2 {
    return t.TransformDirection(Right())
}
```

### ステップ4: 数学ユーティリティ実装

```go
// internal/math/utils.go
package math

import (
    "math"
)

const (
    Pi      = math.Pi
    Pi2     = math.Pi * 2
    PiHalf  = math.Pi / 2
    Deg2Rad = math.Pi / 180
    Rad2Deg = 180 / math.Pi
    Epsilon = 1e-6
)

// Clamp は値を指定範囲に制限する
func Clamp(value, min, max float32) float32 {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

// Lerp は線形補間を行う
func Lerp(a, b, t float32) float32 {
    return a + (b-a)*t
}

// LerpAngle は角度の線形補間を行う
func LerpAngle(a, b, t float32) float32 {
    delta := WrapAngle(b - a)
    return a + delta*t
}

// WrapAngle は角度を-π〜πの範囲に正規化する
func WrapAngle(angle float32) float32 {
    for angle > Pi {
        angle -= Pi2
    }
    for angle < -Pi {
        angle += Pi2
    }
    return angle
}

// Sin はsin関数（float32版）
func Sin(angle float32) float32 {
    return float32(math.Sin(float64(angle)))
}

// Cos はcos関数（float32版）
func Cos(angle float32) float32 {
    return float32(math.Cos(float64(angle)))
}

// Sqrt は平方根（float32版）
func Sqrt(value float32) float32 {
    return float32(math.Sqrt(float64(value)))
}

// Abs は絶対値（float32版）
func Abs(value float32) float32 {
    return float32(math.Abs(float64(value)))
}

// Min は最小値を返す
func Min(a, b float32) float32 {
    if a < b {
        return a
    }
    return b
}

// Max は最大値を返す
func Max(a, b float32) float32 {
    if a > b {
        return a
    }
    return b
}

// Approximately は近似比較を行う
func Approximately(a, b float32) bool {
    return Abs(a-b) < Epsilon
}
```

### ステップ5: 包括的なテスト実装

```go
// internal/math/vector2_test.go
package math

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestVector2_Creation(t *testing.T) {
    // Act
    v := NewVector2(3, 4)
    
    // Assert
    assert.Equal(t, float32(3), v.X)
    assert.Equal(t, float32(4), v.Y)
}

func TestVector2_Constants(t *testing.T) {
    // Assert
    assert.Equal(t, Vector2{0, 0}, Zero())
    assert.Equal(t, Vector2{1, 1}, One())
    assert.Equal(t, Vector2{0, 1}, Up())
    assert.Equal(t, Vector2{1, 0}, Right())
}

func TestVector2_Add(t *testing.T) {
    // Arrange
    v1 := NewVector2(1, 2)
    v2 := NewVector2(3, 4)
    
    // Act
    result := v1.Add(v2)
    
    // Assert
    assert.Equal(t, Vector2{4, 6}, result)
}

func TestVector2_Length(t *testing.T) {
    // Arrange
    v := NewVector2(3, 4)
    
    // Act
    length := v.Length()
    
    // Assert
    assert.InDelta(t, 5.0, length, 0.001)
}

func TestVector2_Normalize(t *testing.T) {
    // Arrange
    v := NewVector2(3, 4)
    
    // Act
    normalized := v.Normalize()
    
    // Assert
    assert.InDelta(t, 1.0, normalized.Length(), 0.001)
    assert.InDelta(t, 0.6, normalized.X, 0.001)
    assert.InDelta(t, 0.8, normalized.Y, 0.001)
}

func TestVector2_Dot(t *testing.T) {
    // Arrange
    v1 := NewVector2(1, 0)
    v2 := NewVector2(0, 1)
    
    // Act
    dot := v1.Dot(v2)
    
    // Assert
    assert.Equal(t, float32(0), dot) // 垂直なベクトルの内積は0
}

func TestVector2_Lerp(t *testing.T) {
    // Arrange
    start := NewVector2(0, 0)
    end := NewVector2(10, 10)
    
    // Act
    middle := start.Lerp(end, 0.5)
    
    // Assert
    assert.Equal(t, Vector2{5, 5}, middle)
}
```

### ステップ6: ビジュアルサンプル実装

```go
// examples/phase3/main.go
package main

import (
    "log"
    "math"
    "runtime"
    "time"
    
    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/glfw/v3.3/glfw"
    gamemath "github.com/yourname/tinyengine/internal/math"
)

func main() {
    log.Println("フェーズ3 ビジュアルサンプル: 数学システムデモ")
    log.Println("Transform, Vector2, Matrix3を使用した動的オブジェクト...")
    
    // OpenGL初期化
    if err := initOpenGL(); err != nil {
        log.Fatalf("OpenGL初期化失敗: %v", err)
    }
    defer glfw.Terminate()
    
    // ウィンドウ作成
    window, err := createWindow()
    if err != nil {
        log.Fatalf("ウィンドウ作成失敗: %v", err)
    }
    defer window.Destroy()
    
    // OpenGL設定
    gl.Viewport(0, 0, 800, 600)
    
    // 数学システムのデモ
    runMathDemo(window)
    
    log.Println("✅ フェーズ3のビジュアルサンプル完了")
    log.Println("")
    log.Println("確認項目:")
    log.Println("- [  ] 四角形が円軌道で回転していた")
    log.Println("- [  ] スケールが時間と共に変化していた")
    log.Println("- [  ] 複数のオブジェクトが同期して動いていた")
    log.Println("- [  ] ベクトル演算と行列変換が正常に動作していた")
}

func runMathDemo(window *glfw.Window) {
    // Transform作成
    transform1 := gamemath.NewTransform()
    transform2 := gamemath.NewTransform()
    transform3 := gamemath.NewTransform()
    
    startTime := time.Now()
    
    for !window.ShouldClose() && time.Since(startTime) < 10*time.Second {
        elapsed := float32(time.Since(startTime).Seconds())
        
        // 画面クリア
        gl.ClearColor(0.1, 0.1, 0.1, 1.0)
        gl.Clear(gl.COLOR_BUFFER_BIT)
        
        // 円軌道運動のデモ
        angle1 := elapsed * 1.0 // 1 rad/s
        radius := float32(0.3)
        pos1 := gamemath.NewVector2(
            radius*gamemath.Cos(angle1),
            radius*gamemath.Sin(angle1),
        )
        transform1.SetPosition(pos1)
        transform1.SetRotation(elapsed * 45) // 45度/秒で回転
        
        // スケールアニメーション
        scale := 0.5 + 0.3*gamemath.Sin(elapsed*3.0)
        transform1.SetUniformScale(scale)
        
        // 楕円軌道運動のデモ
        angle2 := elapsed * 2.0 // 2 rad/s
        pos2 := gamemath.NewVector2(
            0.4*gamemath.Cos(angle2),
            0.2*gamemath.Sin(angle2),
        )
        transform2.SetPosition(pos2)
        transform2.SetRotation(-elapsed * 90) // -90度/秒で回転
        
        // リサージュ曲線のデモ
        pos3 := gamemath.NewVector2(
            0.5*gamemath.Sin(elapsed*1.5),
            0.4*gamemath.Cos(elapsed*2.3),
        )
        transform3.SetPosition(pos3)
        
        // 各オブジェクトを描画
        drawSquare(transform1, gamemath.NewVector2(1.0, 0.0)) // 赤
        drawSquare(transform2, gamemath.NewVector2(0.0, 1.0)) // 緑
        drawSquare(transform3, gamemath.NewVector2(0.0, 0.8)) // 青
        
        // パフォーマンス情報表示
        if int(elapsed*10)%30 == 0 {
            fps := 1.0 / (time.Since(startTime).Seconds() - float64(elapsed))
            log.Printf("Time: %.1fs, FPS: %.1f", elapsed, fps)
            
            // 数学計算のデモ
            v1 := gamemath.NewVector2(1, 0)
            v2 := gamemath.NewVector2(0, 1)
            dot := v1.Dot(v2)
            cross := v1.Cross(v2)
            log.Printf("Vector demo: v1·v2=%.2f, v1×v2=%.2f", dot, cross)
        }
        
        window.SwapBuffers()
        glfw.PollEvents()
    }
}

func drawSquare(transform *gamemath.Transform, color gamemath.Vector2) {
    // 四角形の基本頂点（ローカル座標）
    vertices := []gamemath.Vector2{
        gamemath.NewVector2(-0.05, -0.05),
        gamemath.NewVector2( 0.05, -0.05),
        gamemath.NewVector2( 0.05,  0.05),
        gamemath.NewVector2(-0.05,  0.05),
    }
    
    // Transform適用
    transformedVertices := make([]float32, 0, 8)
    for _, vertex := range vertices {
        worldPos := transform.TransformPoint(vertex)
        transformedVertices = append(transformedVertices, worldPos.X, worldPos.Y)
    }
    
    // OpenGLで描画
    var vao, vbo uint32
    gl.GenVertexArrays(1, &vao)
    gl.GenBuffers(1, &vbo)
    
    gl.BindVertexArray(vao)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len(transformedVertices)*4, gl.Ptr(transformedVertices), gl.STATIC_DRAW)
    
    gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
    gl.EnableVertexAttribArray(0)
    
    // 色設定（簡易実装）
    gl.Color3f(color.X, color.Y, 0.5)
    gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
    
    // クリーンアップ
    gl.DeleteVertexArrays(1, &vao)
    gl.DeleteBuffers(1, &vbo)
}
```

## ビジュアル確認

このフェーズを完了すると、以下が実現できます：

### 期待される結果
- 複数の四角形が異なる軌道で動く
- 円軌道、楕円軌道、リサージュ曲線の動きが見える
- オブジェクトが回転とスケール変化を行う
- ベクトル計算の結果がコンソールに表示される

### 確認項目
- [ ] Vector2の基本操作（加算、減算、内積等）が動作する
- [ ] Matrix3の行列演算が正しく動作する
- [ ] Transformによるオブジェクト変換が機能する
- [ ] 数学ユーティリティ関数が正常に動作する
- [ ] 全テストが成功する

## パフォーマンス考慮事項

### 最適化のポイント
1. **行列キャッシュ**: 変更がない場合は再計算を避ける
2. **SIMD活用**: ベクトル演算の並列化
3. **メモリ配置**: 構造体のメモリレイアウト最適化
4. **関数インライン**: 頻繁に呼ばれる小さな関数

### 精度とパフォーマンスのトレードオフ
- **float32 vs float64**: ゲームでは通常float32で十分
- **三角関数**: ルックアップテーブルによる高速化
- **平方根**: 逆平方根による高速化

## 次のステップ

フェーズ3を完了したら、次はフェーズ4（入力システム）に進みます。ここで実装した数学システムを使用して、マウスやキーボード入力に基づくオブジェクト制御を実装していきます。

## 理解度チェック

1. ベクトルの内積と外積の意味と用途を説明できますか？
2. 2D変換行列の構成要素（平行移動、回転、スケール）を理解していますか？
3. 行列の乗算順序がなぜ重要なのか説明できますか？
4. ゲームでの座標系変換の流れを理解していますか？