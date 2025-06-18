package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/ganyariya/tinyengine/pkg/tinyengine"
)

func TestApplication_Lifecycle(t *testing.T) {
	app := NewApplication()
	
	// 初期化テスト
	err := app.Initialize()
	assert.NoError(t, err)
	
	// 更新テスト
	app.Update(0.016)
	// 基本実装では特に状態変化はないが、エラーが出ないことを確認
	
	// 描画テスト
	app.Render(nil) // レンダラーは後で実装
	
	// 破棄テスト
	app.Destroy()
	
	// GameObjectインターフェースを実装していることを確認
	var gameObj tinyengine.GameObject = app
	assert.NotNil(t, gameObj)
}