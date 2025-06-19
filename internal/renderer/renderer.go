package renderer

import (
	"github.com/ganyariya/tinyengine/pkg/tinyengine"
)

// BaseRenderer は基本的な描画機能を提供する構造体
type BaseRenderer struct {
	width  int
	height int
}

// NewBaseRenderer は新しいBaseRendererを作成する
func NewBaseRenderer(width, height int) tinyengine.Renderer {
	return &BaseRenderer{
		width:  width,
		height: height,
	}
}

// Clear は画面をクリアする
func (r *BaseRenderer) Clear() {
	// 基本実装: 何もしない（OpenGLRendererでオーバーライド）
}

// Present は描画内容を画面に表示する
func (r *BaseRenderer) Present() {
	// 基本実装: 何もしない（OpenGLRendererでオーバーライド）
}

// DrawRectangle は矩形を描画する
func (r *BaseRenderer) DrawRectangle(x, y, width, height float32) {
	// 基本実装: 何もしない（OpenGLRendererでオーバーライド）
}

// DrawPrimitive はプリミティブを描画する
func (r *BaseRenderer) DrawPrimitive(primitive interface{}) {
	// 基本実装: 何もしない（OpenGLRendererでオーバーライド）
}

// DrawRectangleColor は色付き矩形を描画する
func (r *BaseRenderer) DrawRectangleColor(x, y, width, height float32, red, green, blue, alpha float32) {
	// 基本実装: 何もしない（OpenGLRendererでオーバーライド）
}

// DrawCircle は円を描画する
func (r *BaseRenderer) DrawCircle(x, y, radius float32, red, green, blue, alpha float32) {
	// 基本実装: 何もしない（OpenGLRendererでオーバーライド）
}

// DrawLine は線を描画する
func (r *BaseRenderer) DrawLine(x1, y1, x2, y2 float32, red, green, blue, alpha float32) {
	// 基本実装: 何もしない（OpenGLRendererでオーバーライド）
}

// GetSize は描画領域のサイズを取得する
func (r *BaseRenderer) GetSize() (int, int) {
	return r.width, r.height
}
