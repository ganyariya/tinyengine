package core

import (
	"time"
)

// GameLoop はゲームループの管理を行う
type GameLoop struct {
	lastTime    time.Time
	targetFPS   int
	frameTime   float64
}

// NewGameLoop は新しいゲームループインスタンスを作成する
func NewGameLoop() *GameLoop {
	return &GameLoop{
		lastTime:  time.Now(),
		targetFPS: 60, // デフォルト60FPS
		frameTime: 1.0 / 60.0,
	}
}

// GetDeltaTime は前フレームからの経過時間を返す（秒）
func (gl *GameLoop) GetDeltaTime() float64 {
	now := time.Now()
	deltaTime := now.Sub(gl.lastTime).Seconds()
	gl.lastTime = now
	return deltaTime
}

// SetTargetFPS は目標フレームレートを設定する
func (gl *GameLoop) SetTargetFPS(fps int) {
	gl.targetFPS = fps
	gl.frameTime = 1.0 / float64(fps)
}

// GetTargetFPS は目標フレームレートを返す
func (gl *GameLoop) GetTargetFPS() int {
	return gl.targetFPS
}

// GetTargetFrameTime は目標フレーム時間を返す（秒）
func (gl *GameLoop) GetTargetFrameTime() float64 {
	return gl.frameTime
}

// SleepForFrameRate はフレームレート制限のためのスリープを行う
func (gl *GameLoop) SleepForFrameRate() {
	sleepDuration := time.Duration(gl.frameTime * float64(time.Second))
	time.Sleep(sleepDuration)
}