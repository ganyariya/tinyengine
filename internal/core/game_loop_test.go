package core

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestGameLoop_DeltaTime(t *testing.T) {
	loop := NewGameLoop()
	
	// 初回はデルタタイムが0に近い値
	deltaTime := loop.GetDeltaTime()
	assert.GreaterOrEqual(t, deltaTime, 0.0)
	
	// 少し待ってから再度取得
	time.Sleep(10 * time.Millisecond)
	deltaTime = loop.GetDeltaTime()
	assert.Greater(t, deltaTime, 0.0)
	assert.Less(t, deltaTime, 1.0) // 1秒以下であることを確認
}

func TestGameLoop_FrameRate(t *testing.T) {
	loop := NewGameLoop()
	loop.SetTargetFPS(60)
	
	assert.Equal(t, 60, loop.GetTargetFPS())
	
	// フレーム時間の計算確認
	expectedFrameTime := 1.0 / 60.0
	frameTime := loop.GetTargetFrameTime()
	assert.InDelta(t, expectedFrameTime, frameTime, 0.001)
}