package tinyengine

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// テスト用のGameObject実装
type testGameObject struct {
	initialized bool
	updated     bool
	rendered    bool
	destroyed   bool
}

func (t *testGameObject) Initialize() error {
	t.initialized = true
	return nil
}

func (t *testGameObject) Update(deltaTime float64) {
	t.updated = true
}

func (t *testGameObject) Render(renderer Renderer) {
	t.rendered = true
}

func (t *testGameObject) Destroy() {
	t.destroyed = true
}

func TestGameObjectInterface(t *testing.T) {
	obj := &testGameObject{}
	
	// GameObjectインターフェースとして扱えることを確認
	var gameObj GameObject = obj
	
	// 各メソッドが正常に動作することを確認
	err := gameObj.Initialize()
	assert.NoError(t, err)
	assert.True(t, obj.initialized)
	
	gameObj.Update(0.016) // 60FPS相当のデルタタイム
	assert.True(t, obj.updated)
	
	gameObj.Render(nil) // レンダラーは後で実装
	assert.True(t, obj.rendered)
	
	gameObj.Destroy()
	assert.True(t, obj.destroyed)
}

// テスト用のRenderer実装
type testRenderer struct {
	clearCalled bool
	presentCalled bool
}

func (r *testRenderer) Clear() {
	r.clearCalled = true
}

func (r *testRenderer) Present() {
	r.presentCalled = true
}

func (r *testRenderer) DrawRectangle(x, y, width, height float32) {
	// テスト用の空実装
}

func (r *testRenderer) DrawPrimitive(primitive interface{}) {
	// テスト用の空実装
}

func (r *testRenderer) DrawRectangleColor(x, y, width, height float32, red, g, b, a float32) {
	// テスト用の空実装
}

func (r *testRenderer) DrawCircle(x, y, radius float32, red, g, b, a float32) {
	// テスト用の空実装
}

func (r *testRenderer) DrawLine(x1, y1, x2, y2 float32, red, g, b, a float32) {
	// テスト用の空実装
}

func TestRendererInterface(t *testing.T) {
	renderer := &testRenderer{}
	
	// Rendererインターフェースとして扱えることを確認
	var r Renderer = renderer
	
	r.Clear()
	assert.True(t, renderer.clearCalled)
	
	r.Present()
	assert.True(t, renderer.presentCalled)
	
	// 矩形描画メソッドが呼び出せることを確認
	r.DrawRectangle(10, 20, 100, 50)
	
	// 新しい描画メソッドが呼び出せることを確認
	r.DrawRectangleColor(10, 20, 100, 50, 1.0, 0.0, 0.0, 1.0)
	r.DrawCircle(50, 50, 25, 0.0, 1.0, 0.0, 1.0)
	r.DrawLine(0, 0, 100, 100, 0.0, 0.0, 1.0, 1.0)
	r.DrawPrimitive(nil) // nilでもパニックしないことを確認
}