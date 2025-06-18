package platform

import (
	"time"
)

// Timer は時間管理を行う
type Timer struct {
	startTime time.Time
}

// NewTimer は新しいタイマーインスタンスを作成する
func NewTimer() *Timer {
	return &Timer{
		startTime: time.Now(),
	}
}

// GetTime は開始からの経過時間を秒で返す
func (t *Timer) GetTime() float64 {
	return time.Since(t.startTime).Seconds()
}

// Reset はタイマーをリセットする
func (t *Timer) Reset() {
	t.startTime = time.Now()
}