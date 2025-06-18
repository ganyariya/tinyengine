package core

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/ganyariya/tinyengine/pkg/tinyengine"
)

// テスト用のアプリケーション実装
type testApplication struct {
	initialized bool
	updated     bool
	rendered    bool
	destroyed   bool
	updateCount int
}

func (app *testApplication) Initialize() error {
	app.initialized = true
	return nil
}

func (app *testApplication) Update(deltaTime float64) {
	app.updated = true
	app.updateCount++
}

func (app *testApplication) Render(renderer tinyengine.Renderer) {
	app.rendered = true
}

func (app *testApplication) Destroy() {
	app.destroyed = true
}

func TestEngine_Initialize(t *testing.T) {
	engine := NewEngine("テストエンジン", 800, 600)
	
	assert.Equal(t, "テストエンジン", engine.title)
	assert.Equal(t, 800, engine.width)
	assert.Equal(t, 600, engine.height)
	assert.False(t, engine.running)
}

func TestEngine_SetApplication(t *testing.T) {
	engine := NewEngine("テスト", 800, 600)
	app := &testApplication{}
	
	engine.SetApplication(app)
	assert.Equal(t, app, engine.application)
}

func TestEngine_GameLoop(t *testing.T) {
	engine := NewEngine("テスト", 800, 600)
	app := &testApplication{}
	engine.SetApplication(app)
	
	// ゲームループを短時間実行
	go func() {
		time.Sleep(50 * time.Millisecond)
		engine.Stop()
	}()
	
	err := engine.Run()
	assert.NoError(t, err)
	
	// アプリケーションのメソッドが呼び出されたことを確認
	assert.True(t, app.initialized)
	assert.True(t, app.updated)
	assert.True(t, app.rendered)
	assert.True(t, app.destroyed)
	assert.Greater(t, app.updateCount, 0)
}