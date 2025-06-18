package renderer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewColor(t *testing.T) {
	color := NewColor(1.0, 0.5, 0.2, 0.8)
	assert.Equal(t, float32(1.0), color.R)
	assert.Equal(t, float32(0.5), color.G)
	assert.Equal(t, float32(0.2), color.B)
	assert.Equal(t, float32(0.8), color.A)
}

func TestNewColorRGB(t *testing.T) {
	color := NewColorRGB(1.0, 0.5, 0.2)
	assert.Equal(t, float32(1.0), color.R)
	assert.Equal(t, float32(0.5), color.G)
	assert.Equal(t, float32(0.2), color.B)
	assert.Equal(t, float32(1.0), color.A) // アルファ値は1.0
}

func TestNewRectangle(t *testing.T) {
	color := NewColorRGB(1.0, 0.0, 0.0)
	rect := NewRectangle(10, 20, 100, 50, color)
	
	assert.Equal(t, float32(10), rect.X)
	assert.Equal(t, float32(20), rect.Y)
	assert.Equal(t, float32(100), rect.Width)
	assert.Equal(t, float32(50), rect.Height)
	assert.Equal(t, color, rect.Color)
}

func TestRectangleGetVertices(t *testing.T) {
	color := NewColorRGB(1.0, 0.0, 0.0)
	rect := NewRectangle(0, 0, 10, 20, color)
	vertices := rect.GetVertices()
	
	expected := []float32{
		// 左下
		0, 20, 0,
		// 右下
		10, 20, 0,
		// 右上
		10, 0, 0,
		// 左上
		0, 0, 0,
	}
	
	assert.Equal(t, expected, vertices)
}

func TestRectangleGetIndices(t *testing.T) {
	color := NewColorRGB(1.0, 0.0, 0.0)
	rect := NewRectangle(0, 0, 10, 20, color)
	indices := rect.GetIndices()
	
	expected := []uint32{
		0, 1, 2, // 第1三角形
		2, 3, 0, // 第2三角形
	}
	
	assert.Equal(t, expected, indices)
}

func TestRectangleInterface(t *testing.T) {
	color := NewColorRGB(1.0, 0.0, 0.0)
	rect := NewRectangle(0, 0, 10, 20, color)
	
	assert.Equal(t, color, rect.GetColor())
	assert.Equal(t, PrimitiveTypeRectangle, rect.GetType())
}

func TestNewCircle(t *testing.T) {
	color := NewColorRGB(0.0, 1.0, 0.0)
	circle := NewCircle(50, 50, 25, color)
	
	assert.Equal(t, float32(50), circle.X)
	assert.Equal(t, float32(50), circle.Y)
	assert.Equal(t, float32(25), circle.Radius)
	assert.Equal(t, color, circle.Color)
	assert.Equal(t, 32, circle.Segments) // デフォルト値
}

func TestNewCircleWithSegments(t *testing.T) {
	color := NewColorRGB(0.0, 1.0, 0.0)
	circle := NewCircleWithSegments(50, 50, 25, color, 16)
	
	assert.Equal(t, 16, circle.Segments)
}

func TestCircleGetVertices(t *testing.T) {
	color := NewColorRGB(0.0, 1.0, 0.0)
	circle := NewCircleWithSegments(0, 0, 10, color, 4) // 正方形に近い形
	vertices := circle.GetVertices()
	
	// 頂点数の確認: 中心点(1) + 外周点(4) + 最初の外周点(1) = 6点
	expectedVertexCount := (4 + 2) * 3 // 各点は3要素(x,y,z)
	assert.Equal(t, expectedVertexCount, len(vertices))
	
	// 中心点の確認
	assert.Equal(t, float32(0), vertices[0]) // x
	assert.Equal(t, float32(0), vertices[1]) // y
	assert.Equal(t, float32(0), vertices[2]) // z
	
	// 最初の外周点の確認（角度0度 = (半径, 0)）
	assert.InDelta(t, 10.0, vertices[3], 0.001) // x
	assert.InDelta(t, 0.0, vertices[4], 0.001)  // y
	assert.Equal(t, float32(0), vertices[5])    // z
}

func TestCircleGetIndices(t *testing.T) {
	color := NewColorRGB(0.0, 1.0, 0.0)
	circle := NewCircleWithSegments(0, 0, 10, color, 4)
	indices := circle.GetIndices()
	
	// インデックス数の確認: セグメント数 * 3
	expectedIndexCount := 4 * 3
	assert.Equal(t, expectedIndexCount, len(indices))
	
	// 最初の三角形の確認
	assert.Equal(t, uint32(0), indices[0]) // 中心点
	assert.Equal(t, uint32(1), indices[1]) // 最初の外周点
	assert.Equal(t, uint32(2), indices[2]) // 次の外周点
}

func TestCircleInterface(t *testing.T) {
	color := NewColorRGB(0.0, 1.0, 0.0)
	circle := NewCircle(0, 0, 10, color)
	
	assert.Equal(t, color, circle.GetColor())
	assert.Equal(t, PrimitiveTypeCircle, circle.GetType())
}

func TestNewLine(t *testing.T) {
	color := NewColorRGB(0.0, 0.0, 1.0)
	line := NewLine(0, 0, 10, 20, color)
	
	assert.Equal(t, float32(0), line.X1)
	assert.Equal(t, float32(0), line.Y1)
	assert.Equal(t, float32(10), line.X2)
	assert.Equal(t, float32(20), line.Y2)
	assert.Equal(t, color, line.Color)
	assert.Equal(t, float32(1.0), line.Width) // デフォルト幅
}

func TestLineGetVertices(t *testing.T) {
	color := NewColorRGB(0.0, 0.0, 1.0)
	line := NewLine(5, 10, 15, 30, color)
	vertices := line.GetVertices()
	
	expected := []float32{
		5, 10, 0,  // 開始点
		15, 30, 0, // 終了点
	}
	
	assert.Equal(t, expected, vertices)
}

func TestLineGetIndices(t *testing.T) {
	color := NewColorRGB(0.0, 0.0, 1.0)
	line := NewLine(0, 0, 10, 20, color)
	indices := line.GetIndices()
	
	expected := []uint32{0, 1}
	assert.Equal(t, expected, indices)
}

func TestLineInterface(t *testing.T) {
	color := NewColorRGB(0.0, 0.0, 1.0)
	line := NewLine(0, 0, 10, 20, color)
	
	assert.Equal(t, color, line.GetColor())
	assert.Equal(t, PrimitiveTypeLine, line.GetType())
}

func TestCircleVerticesCorrectness(t *testing.T) {
	// より厳密な円の頂点計算テスト
	color := NewColorRGB(1.0, 0.0, 0.0)
	circle := NewCircleWithSegments(100, 100, 50, color, 8)
	vertices := circle.GetVertices()
	
	// 中心点は固定
	assert.Equal(t, float32(100), vertices[0])
	assert.Equal(t, float32(100), vertices[1])
	
	// 各外周点の距離が半径と一致することを確認
	for i := 1; i <= 8; i++ {
		x := vertices[i*3]
		y := vertices[i*3+1]
		
		// 中心からの距離を計算
		dx := x - 100
		dy := y - 100
		distance := math.Sqrt(float64(dx*dx + dy*dy))
		
		// 半径50との誤差を確認
		assert.InDelta(t, 50.0, distance, 0.001, "外周点%dの距離が正しくありません", i)
	}
}