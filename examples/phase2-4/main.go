package main

import (
	"fmt"
	stdmath "math"
	"runtime"
	"time"

	mathlib "github.com/ganyariya/tinyengine/internal/math"
	"github.com/ganyariya/tinyengine/internal/renderer"
	"github.com/ganyariya/tinyengine/pkg/tinyengine"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	// ウィンドウ設定
	WindowWidth  = mathlib.DefaultWindowWidth
	WindowHeight = mathlib.DefaultWindowHeight
	WindowTitle  = "Phase 2-4: Transform Demo - Rotating, Scaling, Moving Rectangles"
	
	// アニメーション設定
	DefaultRotationSpeed = mathlib.DefaultRotationSpeed // 1.0 ラジアン/秒
	DefaultScaleSpeed    = mathlib.DefaultScaleSpeed    // 0.5 スケール変化速度
	DefaultMoveSpeed     = mathlib.DefaultMoveSpeed     // 50.0 ピクセル/秒
	ScaleOscillation     = mathlib.ScaleOscillation     // 0.3 スケール振動幅
	MinAnimationScale    = mathlib.MinAnimationScale    // 0.1 最小アニメーションスケール
	CircularRadius       = mathlib.DefaultRadius        // 100.0 円運動の半径
	CircularSpeedDivisor = mathlib.CircularSpeedDivisor // 100.0 円運動速度の除数
	
	// 矩形サイズ設定
	RedRectWidth   = 60.0
	RedRectHeight  = 40.0
	GreenRectWidth = 80.0
	GreenRectHeight = 30.0
	BlueRectSize   = 50.0 // 正方形
	
	// 各矩形の固有設定
	GreenRotationSpeed = 0.5
	GreenScaleSpeed    = 1.0
	GreenMoveSpeed     = 30.0
	GreenBaseScale     = 1.2
	
	BlueRotationSpeed  = -0.3 // 逆回転
	BlueScaleSpeed     = 0.2
	BlueMoveSpeed      = -40.0 // 逆移動
	BlueBaseScale      = 0.8
	
	// FPS表示設定
	FPSDisplayInterval = 1.0 // 1秒間隔
	FallbackFrameLimit = 300 // フォールバック時のフレーム数制限（約5秒 @ 60fps）
)

func init() {
	// OpenGLコンテキストはメインスレッドで実行する必要がある
	runtime.LockOSThread()
}

// TransformableRectangle 変形可能な矩形を表現する構造体
type TransformableRectangle struct {
	transform  mathlib.Transform // 座標変換情報（位置、回転、スケール）
	size       mathlib.Vector2   // 矩形のサイズ
	color      renderer.Color    // 描画色
	
	// アニメーション特性
	rotationSpeed float64 // 回転速度（ラジアン/秒）
	scaleSpeed    float64 // スケール変化速度（スケール単位/秒）
	moveSpeed     float64 // 移動速度（ピクセル/秒）
	
	// アニメーション状態
	time          float64 // 経過時間
	baseScale     float64 // 基準スケール値
}

// NewTransformableRectangle 新しい変形可能な矩形を作成
func NewTransformableRectangle(position mathlib.Vector2, size mathlib.Vector2, color renderer.Color) *TransformableRectangle {
	return &TransformableRectangle{
		transform: mathlib.NewTransformWithValues(position, 0, mathlib.Vector2{X: 1, Y: 1}),
		size:      size,
		color:     color,
		rotationSpeed: DefaultRotationSpeed,
		scaleSpeed:    DefaultScaleSpeed,
		moveSpeed:     DefaultMoveSpeed,
		baseScale:     1.0,
	}
}

// Update 矩形のアニメーションを更新
func (tr *TransformableRectangle) Update(deltaTime float64) {
	tr.time += deltaTime
	
	// 回転アニメーション
	tr.transform.Rotate(tr.rotationSpeed * deltaTime)
	
	// スケールアニメーション（振動）
	scaleOffset := stdmath.Sin(tr.time * tr.scaleSpeed) * ScaleOscillation
	newScale := tr.baseScale + scaleOffset
	if newScale > MinAnimationScale { // 負の値や小さすぎるスケールを防ぐ
		tr.transform.SetUniformScale(newScale)
	}
	
	// 位置アニメーション（円運動）
	centerX := float64(WindowWidth) * 0.5
	centerY := float64(WindowHeight) * 0.5
	radius := CircularRadius
	
	x := centerX + radius * stdmath.Cos(tr.time * tr.moveSpeed / CircularSpeedDivisor)
	y := centerY + radius * stdmath.Sin(tr.time * tr.moveSpeed / CircularSpeedDivisor)
	
	tr.transform.SetPosition(mathlib.Vector2{X: x, Y: y})
}

// Render 指定されたレンダラーを使用して矩形を描画
func (tr *TransformableRectangle) Render(r tinyengine.Renderer) {
	// 数学ライブラリから変換行列を取得
	transformMatrix := tr.transform.ToMatrix()
	
	// OpenGL用に3x3行列を4x4行列に変換
	transform4x4 := convert3x3To4x4(transformMatrix)
	
	// 原点に基本的な矩形を作成してから変形
	halfWidth := float32(tr.size.X * 0.5)
	halfHeight := float32(tr.size.Y * 0.5)
	
	vertices := []float32{
		-halfWidth, -halfHeight, 0.0, // Bottom left
		 halfWidth, -halfHeight, 0.0, // Bottom right
		 halfWidth,  halfHeight, 0.0, // Top right
		-halfWidth,  halfHeight, 0.0, // Top left
	}
	
	indices := []uint32{
		0, 1, 2, // First triangle
		2, 3, 0, // Second triangle
	}
	
	// 変換を適用して描画
	transformedVertices := applyTransformToVertices(vertices, transform4x4)
	
	// 描画用のプリミティブを作成
	rect := &TransformedRectangle{
		vertices: transformedVertices,
		indices:  indices,
		color:    tr.color,
	}
	
	r.DrawPrimitive(rect)
}

// TransformedRectangle Primitiveインターフェースを実装する変換済み矩形
type TransformedRectangle struct {
	vertices []float32
	indices  []uint32
	color    renderer.Color
}

func (tr *TransformedRectangle) GetVertices() []float32 {
	return tr.vertices
}

func (tr *TransformedRectangle) GetIndices() []uint32 {
	return tr.indices
}

func (tr *TransformedRectangle) GetColor() renderer.Color {
	return tr.color
}

func (tr *TransformedRectangle) GetType() renderer.PrimitiveType {
	return renderer.PrimitiveTypeRectangle
}

// convert3x3To4x4 OpenGL用に3x3行列を4x4行列に変換
func convert3x3To4x4(m3 mathlib.Matrix3x3) [16]float32 {
	return [16]float32{
		float32(m3[0][0]), float32(m3[1][0]), 0, float32(m3[2][0]),
		float32(m3[0][1]), float32(m3[1][1]), 0, float32(m3[2][1]),
		0,                 0,                 1, 0,
		float32(m3[0][2]), float32(m3[1][2]), 0, float32(m3[2][2]),
	}
}

// applyTransformToVertices 4x4変換行列を頂点に適用
func applyTransformToVertices(vertices []float32, transform [16]float32) []float32 {
	transformed := make([]float32, len(vertices))
	
	// 頂点を3個ずつ（x, y, z）のグループで処理
	for i := 0; i < len(vertices); i += 3 {
		x, y, z := vertices[i], vertices[i+1], vertices[i+2]
		
		// 4x4変換行列を適用
		transformed[i] = transform[0]*x + transform[4]*y + transform[8]*z + transform[12]   // new x
		transformed[i+1] = transform[1]*x + transform[5]*y + transform[9]*z + transform[13] // new y
		transformed[i+2] = transform[2]*x + transform[6]*y + transform[10]*z + transform[14] // new z
	}
	
	return transformed
}

// createRedRectangle 赤い矩形を作成
func createRedRectangle() *TransformableRectangle {
	return NewTransformableRectangle(
		mathlib.Vector2{X: float64(WindowWidth) * 0.5, Y: float64(WindowHeight) * 0.5},
		mathlib.Vector2{X: RedRectWidth, Y: RedRectHeight},
		renderer.NewColor(1.0, 0.0, 0.0, 1.0),
	)
}

// createGreenRectangle 緑の矩形を作成
func createGreenRectangle() *TransformableRectangle {
	return &TransformableRectangle{
		transform: mathlib.NewTransformWithValues(
			mathlib.Vector2{X: float64(WindowWidth) * 0.5, Y: float64(WindowHeight) * 0.5},
			0,
			mathlib.Vector2{X: 1, Y: 1},
		),
		size:          mathlib.Vector2{X: GreenRectWidth, Y: GreenRectHeight},
		color:         renderer.NewColor(0.0, 1.0, 0.0, 1.0),
		rotationSpeed: GreenRotationSpeed,
		scaleSpeed:    GreenScaleSpeed,
		moveSpeed:     GreenMoveSpeed,
		baseScale:     GreenBaseScale,
	}
}

// createBlueRectangle 青い矩形を作成
func createBlueRectangle() *TransformableRectangle {
	return &TransformableRectangle{
		transform: mathlib.NewTransformWithValues(
			mathlib.Vector2{X: float64(WindowWidth) * 0.5, Y: float64(WindowHeight) * 0.5},
			0,
			mathlib.Vector2{X: 1, Y: 1},
		),
		size:          mathlib.Vector2{X: BlueRectSize, Y: BlueRectSize},
		color:         renderer.NewColor(0.0, 0.0, 1.0, 1.0),
		rotationSpeed: BlueRotationSpeed, // 逆回転
		scaleSpeed:    BlueScaleSpeed,
		moveSpeed:     BlueMoveSpeed,     // 逆移動
		baseScale:     BlueBaseScale,
	}
}

// initializeRenderer レンダラーとウィンドウを初期化
func initializeRenderer() (tinyengine.Renderer, *glfw.Window, error) {
	r, err := renderer.NewOpenGLRendererWithWindow(WindowWidth, WindowHeight, WindowTitle)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create renderer: %w", err)
	}
	
	// テスト：基本的なインターフェースが動作することを確認
	r.Clear()
	
	// 入力処理のためのGLFWウィンドウアクセス
	var window *glfw.Window
	if openglRenderer, ok := r.(*renderer.OpenGLRenderer); ok {
		window = openglRenderer.GetWindow()
	}
	
	if window == nil {
		fmt.Println("Warning: Could not access GLFW window, input handling disabled")
	}
	
	return r, window, nil
}

// FPSCounter FPS計測のためのヘルパー構造体
type FPSCounter struct {
	frameCount  int
	lastTime    time.Time
	lastFPSTime time.Time
}

// NewFPSCounter 新しいFPSカウンターを作成
func NewFPSCounter() *FPSCounter {
	now := time.Now()
	return &FPSCounter{
		frameCount:  0,
		lastTime:    now,
		lastFPSTime: now,
	}
}

// Update デルタタイムを計算し、FPSを表示（必要に応じて）
func (fps *FPSCounter) Update() float64 {
	currentTime := time.Now()
	deltaTime := currentTime.Sub(fps.lastTime).Seconds()
	fps.lastTime = currentTime
	
	fps.frameCount++
	if time.Since(fps.lastFPSTime).Seconds() >= FPSDisplayInterval {
		currentFPS := float64(fps.frameCount) / time.Since(fps.lastFPSTime).Seconds()
		fmt.Printf("FPS: %.1f\n", currentFPS)
		fps.frameCount = 0
		fps.lastFPSTime = time.Now()
	}
	
	return deltaTime
}

// GetFrameCount フレーム数を取得（フォールバック用）
func (fps *FPSCounter) GetFrameCount() int {
	return fps.frameCount
}

// handleInput 入力処理
func handleInput(window *glfw.Window, frameCount int) bool {
	if window != nil {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
		return !window.ShouldClose()
	} else {
		// フォールバック：ウィンドウがない場合は一定フレーム後に終了
		return frameCount <= FallbackFrameLimit
	}
}

// updateRectangles 全ての矩形を更新
func updateRectangles(rectangles []*TransformableRectangle, deltaTime float64) {
	for _, rect := range rectangles {
		rect.Update(deltaTime)
	}
}

// renderRectangles 全ての矩形を描画
func renderRectangles(r tinyengine.Renderer, rectangles []*TransformableRectangle) {
	r.Clear()
	for _, rect := range rectangles {
		rect.Render(r)
	}
	r.Present()
}

// runTransformDemo トランスフォームデモのメインループを実行
func runTransformDemo(r tinyengine.Renderer, window *glfw.Window, rectangles []*TransformableRectangle) {
	fmt.Println("Transform Demo Controls:")
	fmt.Println("- ESC: Exit")
	fmt.Println("- Watch the rectangles rotate, scale, and move in circular patterns!")
	
	fpsCounter := NewFPSCounter()
	
	// メインレンダーループ
	for {
		deltaTime := fpsCounter.Update()
		
		// 入力処理
		if !handleInput(window, fpsCounter.GetFrameCount()) {
			break
		}
		
		// 全ての矩形を更新
		updateRectangles(rectangles, deltaTime)
		
		// 描画
		renderRectangles(r, rectangles)
	}
}

func main() {
	fmt.Println("Starting Phase 2-4 Transform Demo...")
	
	// レンダラーとウィンドウの初期化
	r, window, err := initializeRenderer()
	if err != nil {
		fmt.Printf("Initialization failed: %v\n", err)
		return
	}
	defer func() {
		if openglRenderer, ok := r.(*renderer.OpenGLRenderer); ok {
			openglRenderer.Destroy()
		}
	}()
	
	// 様々な特性を持つ変形可能な矩形を作成
	rectangles := []*TransformableRectangle{
		createRedRectangle(),   // 高速回転、中程度のスケール振動
		createGreenRectangle(), // 中程度の回転、高速スケール振動
		createBlueRectangle(),  // 低速回転、低速スケール振動
	}
	
	// デモの実行
	runTransformDemo(r, window, rectangles)
	
	fmt.Println("Transform Demo finished.")
}