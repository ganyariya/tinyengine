package platform

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestTimer_Basic(t *testing.T) {
	timer := NewTimer()
	
	// 初期値の確認
	assert.Greater(t, timer.GetTime(), 0.0)
	
	// 少し待ってから時間の進行を確認
	startTime := timer.GetTime()
	time.Sleep(10 * time.Millisecond)
	endTime := timer.GetTime()
	
	assert.Greater(t, endTime, startTime)
}

func TestTimer_Reset(t *testing.T) {
	timer := NewTimer()
	
	// 時間を進める
	time.Sleep(10 * time.Millisecond)
	
	// リセット前の時間を記録
	timeBeforeReset := timer.GetTime()
	assert.Greater(t, timeBeforeReset, 0.0)
	
	// リセット
	timer.Reset()
	timeAfterReset := timer.GetTime()
	
	// リセット後は時間が小さくなっている
	assert.Less(t, timeAfterReset, timeBeforeReset)
}