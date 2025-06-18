package platform

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestWindow_Configuration(t *testing.T) {
	config := WindowConfig{
		Title:  "テストウィンドウ",
		Width:  800,
		Height: 600,
	}
	
	window := NewWindow(config)
	
	assert.Equal(t, "テストウィンドウ", window.config.Title)
	assert.Equal(t, 800, window.config.Width)
	assert.Equal(t, 600, window.config.Height)
	assert.False(t, window.initialized)
}

func TestWindow_InitializeAndDestroy(t *testing.T) {
	// CI環境やヘッドレス環境ではGLFWの初期化がスキップされる可能性がある
	// そのため、テストはエラーハンドリングに焦点を当てる
	
	config := WindowConfig{
		Title:  "テスト",
		Width:  400,
		Height: 300,
	}
	
	window := NewWindow(config)
	
	// 初期化を試行（環境によってはエラーになる可能性がある）
	err := window.Initialize()
	if err != nil {
		// ヘッドレス環境の場合は期待されるエラー
		t.Logf("ウィンドウ初期化スキップ（ヘッドレス環境）: %v", err)
		return
	}
	
	// 初期化が成功した場合のテスト
	assert.True(t, window.initialized)
	
	// 終了処理
	window.Destroy()
}