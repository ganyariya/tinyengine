package renderer

import (
	"math"
)

// Color は色情報を表す構造体
type Color struct {
	R, G, B, A float32
}

// NewColor は新しいColorを作成する
func NewColor(r, g, b, a float32) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// NewColorRGB はアルファ値1.0でColorを作成する
func NewColorRGB(r, g, b float32) Color {
	return Color{R: r, G: g, B: b, A: 1.0}
}

// Primitive は描画プリミティブの基底インターフェース
type Primitive interface {
	// GetVertices は頂点データを取得する
	GetVertices() []float32
	
	// GetIndices はインデックスデータを取得する（必要な場合）
	GetIndices() []uint32
	
	// GetColor は色情報を取得する
	GetColor() Color
	
	// GetType はプリミティブの種類を取得する
	GetType() PrimitiveType
}

// PrimitiveType はプリミティブの種類を表す
type PrimitiveType int

const (
	PrimitiveTypeTriangle PrimitiveType = iota
	PrimitiveTypeRectangle
	PrimitiveTypeCircle
	PrimitiveTypeLine
)

// Rectangle は矩形プリミティブ
type Rectangle struct {
	X, Y          float32 // 左上角の座標
	Width, Height float32 // 幅と高さ
	Color         Color   // 色
}

// NewRectangle は新しい矩形を作成する
func NewRectangle(x, y, width, height float32, color Color) *Rectangle {
	return &Rectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Color:  color,
	}
}

// GetVertices は矩形の頂点データを取得する
// 各頂点は x, y, z の3要素で構成される
func (r *Rectangle) GetVertices() []float32 {
	return []float32{
		// 左下
		r.X, r.Y + r.Height, 0.0,
		// 右下
		r.X + r.Width, r.Y + r.Height, 0.0,
		// 右上
		r.X + r.Width, r.Y, 0.0,
		// 左上
		r.X, r.Y, 0.0,
	}
}

// GetIndices は矩形のインデックスデータを取得する
func (r *Rectangle) GetIndices() []uint32 {
	return []uint32{
		0, 1, 2, // 第1三角形
		2, 3, 0, // 第2三角形
	}
}

// GetColor は矩形の色を取得する
func (r *Rectangle) GetColor() Color {
	return r.Color
}

// GetType は矩形のプリミティブタイプを取得する
func (r *Rectangle) GetType() PrimitiveType {
	return PrimitiveTypeRectangle
}

// Circle は円プリミティブ
type Circle struct {
	X, Y   float32 // 中心座標
	Radius float32 // 半径
	Color  Color   // 色
	Segments int   // 円を構成する線分数（デフォルト32）
}

// NewCircle は新しい円を作成する
func NewCircle(x, y, radius float32, color Color) *Circle {
	return &Circle{
		X:        x,
		Y:        y,
		Radius:   radius,
		Color:    color,
		Segments: 32, // デフォルト値
	}
}

// NewCircleWithSegments は線分数を指定して新しい円を作成する
func NewCircleWithSegments(x, y, radius float32, color Color, segments int) *Circle {
	return &Circle{
		X:        x,
		Y:        y,
		Radius:   radius,
		Color:    color,
		Segments: segments,
	}
}

// GetVertices は円の頂点データを取得する
func (c *Circle) GetVertices() []float32 {
	vertices := make([]float32, (c.Segments+2)*3) // 中心点 + 外周点 + 最初の外周点
	
	// 中心点
	vertices[0] = c.X
	vertices[1] = c.Y
	vertices[2] = 0.0
	
	// 外周点を計算
	for i := 0; i <= c.Segments; i++ {
		angle := 2.0 * math.Pi * float64(i) / float64(c.Segments)
		x := c.X + c.Radius*float32(math.Cos(angle))
		y := c.Y + c.Radius*float32(math.Sin(angle))
		
		idx := (i + 1) * 3
		vertices[idx] = x
		vertices[idx+1] = y
		vertices[idx+2] = 0.0
	}
	
	return vertices
}

// GetIndices は円のインデックスデータを取得する
func (c *Circle) GetIndices() []uint32 {
	indices := make([]uint32, c.Segments*3)
	
	for i := 0; i < c.Segments; i++ {
		indices[i*3] = 0                    // 中心点
		indices[i*3+1] = uint32(i + 1)      // 現在の外周点
		indices[i*3+2] = uint32(i + 2)      // 次の外周点
	}
	
	return indices
}

// GetColor は円の色を取得する
func (c *Circle) GetColor() Color {
	return c.Color
}

// GetType は円のプリミティブタイプを取得する
func (c *Circle) GetType() PrimitiveType {
	return PrimitiveTypeCircle
}

// Line は線プリミティブ
type Line struct {
	X1, Y1 float32 // 開始点
	X2, Y2 float32 // 終了点
	Color  Color   // 色
	Width  float32 // 線の太さ（将来対応）
}

// NewLine は新しい線を作成する
func NewLine(x1, y1, x2, y2 float32, color Color) *Line {
	return &Line{
		X1:    x1,
		Y1:    y1,
		X2:    x2,
		Y2:    y2,
		Color: color,
		Width: 1.0, // デフォルト幅
	}
}

// GetVertices は線の頂点データを取得する
func (l *Line) GetVertices() []float32 {
	return []float32{
		l.X1, l.Y1, 0.0, // 開始点
		l.X2, l.Y2, 0.0, // 終了点
	}
}

// GetIndices は線のインデックスデータを取得する（線は不要）
func (l *Line) GetIndices() []uint32 {
	return []uint32{0, 1}
}

// GetColor は線の色を取得する
func (l *Line) GetColor() Color {
	return l.Color
}

// GetType は線のプリミティブタイプを取得する
func (l *Line) GetType() PrimitiveType {
	return PrimitiveTypeLine
}