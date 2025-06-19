package core

import "time"

// Frame rate constants
const (
	DefaultTargetFPS        = 60
	DefaultFrameTimeSeconds = 1.0 / DefaultTargetFPS
	DefaultFrameTimeMs      = time.Millisecond * 16 // ~60FPS
)