package math

import (
	stdmath "math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamera2D_NewCamera2D(t *testing.T) {
	camera := NewCamera2D()
	
	assert.Equal(t, Vector2{X: 0, Y: 0}, camera.Position)
	assert.Equal(t, 1.0, camera.Zoom)
	assert.Equal(t, 0.0, camera.Rotation)
}

func TestCamera2D_NewCamera2DWithValues(t *testing.T) {
	position := Vector2{X: 10, Y: 5}
	zoom := 2.0
	rotation := stdmath.Pi / 4
	
	camera := NewCamera2DWithValues(position, zoom, rotation)
	
	assert.Equal(t, position, camera.Position)
	assert.Equal(t, zoom, camera.Zoom)
	assert.Equal(t, rotation, camera.Rotation)
}

func TestCamera2D_GetViewMatrix(t *testing.T) {
	camera := NewCamera2DWithValues(
		Vector2{X: 5, Y: 3},
		2.0,
		0, // no rotation for simplicity
	)
	
	viewMatrix := camera.GetViewMatrix()
	
	// Test transforming a world point
	worldPoint := Vector2{X: 10, Y: 8}
	viewPoint := viewMatrix.TransformPoint(worldPoint)
	
	// Expected: translate(-5,-3) -> (5,5), scale(2,2) -> (10,10)
	assert.InDelta(t, 10.0, viewPoint.X, Epsilon)
	assert.InDelta(t, 10.0, viewPoint.Y, Epsilon)
}

func TestCamera2D_GetProjectionMatrix(t *testing.T) {
	camera := NewCamera2D()
	screenWidth := 800.0
	screenHeight := 600.0
	
	projMatrix := camera.GetProjectionMatrix(screenWidth, screenHeight)
	
	// Test transforming normalized coordinates to screen coordinates
	normalizedPoint := Vector2{X: 0, Y: 0} // Center
	screenPoint := projMatrix.TransformPoint(normalizedPoint)
	
	// Expected: center of screen
	assert.Equal(t, 400.0, screenPoint.X) // screenWidth / 2
	assert.Equal(t, 300.0, screenPoint.Y) // screenHeight / 2
}

func TestCamera2D_ScreenToWorld(t *testing.T) {
	camera := NewCamera2D()
	screenWidth := 800.0
	screenHeight := 600.0
	
	// Test center of screen
	screenCenter := Vector2{X: 400, Y: 300}
	worldPoint := camera.ScreenToWorld(screenCenter, screenWidth, screenHeight)
	
	// Should be at world origin
	assert.InDelta(t, 0.0, worldPoint.X, Epsilon)
	assert.InDelta(t, 0.0, worldPoint.Y, Epsilon)
}

func TestCamera2D_WorldToScreen(t *testing.T) {
	camera := NewCamera2D()
	screenWidth := 800.0
	screenHeight := 600.0
	
	// Test world origin
	worldOrigin := Vector2{X: 0, Y: 0}
	screenPoint := camera.WorldToScreen(worldOrigin, screenWidth, screenHeight)
	
	// Should be at screen center
	assert.InDelta(t, 400.0, screenPoint.X, Epsilon)
	assert.InDelta(t, 300.0, screenPoint.Y, Epsilon)
}

func TestCamera2D_ScreenToWorld_WorldToScreen_RoundTrip(t *testing.T) {
	camera := NewCamera2DWithValues(
		Vector2{X: 10, Y: 5},
		1.5,
		stdmath.Pi/6, // 30 degrees
	)
	screenWidth := 800.0
	screenHeight := 600.0
	
	originalWorldPoint := Vector2{X: 20, Y: 15}
	
	// World -> Screen -> World
	screenPoint := camera.WorldToScreen(originalWorldPoint, screenWidth, screenHeight)
	backToWorldPoint := camera.ScreenToWorld(screenPoint, screenWidth, screenHeight)
	
	assert.InDelta(t, originalWorldPoint.X, backToWorldPoint.X, 1e-8)
	assert.InDelta(t, originalWorldPoint.Y, backToWorldPoint.Y, 1e-8)
}

func TestCamera2D_SetPosition(t *testing.T) {
	camera := NewCamera2D()
	newPosition := Vector2{X: 15, Y: 25}
	
	camera.SetPosition(newPosition)
	
	assert.Equal(t, newPosition, camera.Position)
}

func TestCamera2D_SetZoom(t *testing.T) {
	camera := NewCamera2D()
	
	camera.SetZoom(2.5)
	assert.Equal(t, 2.5, camera.Zoom)
	
	// Test invalid zoom (should not change)
	camera.SetZoom(-1.0)
	assert.Equal(t, 2.5, camera.Zoom)
	
	camera.SetZoom(0.0)
	assert.Equal(t, 2.5, camera.Zoom)
}

func TestCamera2D_SetRotation(t *testing.T) {
	camera := NewCamera2D()
	rotation := stdmath.Pi / 3
	
	camera.SetRotation(rotation)
	
	assert.Equal(t, rotation, camera.Rotation)
}

func TestCamera2D_SetRotationDegrees(t *testing.T) {
	camera := NewCamera2D()
	
	camera.SetRotationDegrees(90)
	
	assert.InDelta(t, stdmath.Pi/2, camera.Rotation, Epsilon)
}

func TestCamera2D_Move(t *testing.T) {
	camera := NewCamera2DWithValues(Vector2{X: 5, Y: 3}, 1.0, 0.0)
	offset := Vector2{X: 2, Y: 4}
	
	camera.Move(offset)
	
	expected := Vector2{X: 7, Y: 7}
	assert.Equal(t, expected, camera.Position)
}

func TestCamera2D_ZoomBy(t *testing.T) {
	camera := NewCamera2DWithValues(Vector2{X: 0, Y: 0}, 2.0, 0.0)
	
	camera.ZoomBy(1.5)
	assert.Equal(t, 3.0, camera.Zoom)
	
	// Test invalid zoom factor (should not change)
	camera.ZoomBy(-1.0)
	assert.Equal(t, 3.0, camera.Zoom)
	
	camera.ZoomBy(0.0)
	assert.Equal(t, 3.0, camera.Zoom)
}

func TestCamera2D_Rotate(t *testing.T) {
	camera := NewCamera2DWithValues(Vector2{X: 0, Y: 0}, 1.0, stdmath.Pi/4)
	
	camera.Rotate(stdmath.Pi / 4)
	
	assert.InDelta(t, stdmath.Pi/2, camera.Rotation, Epsilon)
}

func TestCamera2D_RotateDegrees(t *testing.T) {
	camera := NewCamera2D()
	
	camera.RotateDegrees(45)
	camera.RotateDegrees(45)
	
	assert.InDelta(t, stdmath.Pi/2, camera.Rotation, Epsilon)
}

func TestCamera2D_GetBounds(t *testing.T) {
	camera := NewCamera2D()
	screenWidth := 800.0
	screenHeight := 600.0
	
	minBounds, maxBounds := camera.GetBounds(screenWidth, screenHeight)
	
	// With default camera (no zoom, no rotation, centered at origin)
	// bounds should be symmetrical around origin
	assert.True(t, minBounds.X < 0)
	assert.True(t, minBounds.Y < 0)
	assert.True(t, maxBounds.X > 0)
	assert.True(t, maxBounds.Y > 0)
}

func TestCamera2D_LookAt(t *testing.T) {
	camera := NewCamera2D()
	target := Vector2{X: 100, Y: 50}
	
	camera.LookAt(target)
	
	assert.Equal(t, target, camera.Position)
}

func TestCamera2D_FollowTarget(t *testing.T) {
	camera := NewCamera2D()
	target := Vector2{X: 10, Y: 5}
	followSpeed := 5.0
	deltaTime := 1.0
	
	camera.FollowTarget(target, followSpeed, deltaTime)
	
	// Should move towards target
	distance := camera.Position.Distance(Vector2{X: 0, Y: 0})
	assert.True(t, distance > 0)
	assert.True(t, distance <= followSpeed*deltaTime)
}

func TestCamera2D_FollowTarget_AlreadyAtTarget(t *testing.T) {
	target := Vector2{X: 5, Y: 3}
	camera := NewCamera2DWithValues(target, 1.0, 0.0)
	followSpeed := 5.0
	deltaTime := 1.0
	
	camera.FollowTarget(target, followSpeed, deltaTime)
	
	// Should remain at target
	assert.Equal(t, target, camera.Position)
}

func TestCamera2D_FollowTarget_InvalidParameters(t *testing.T) {
	camera := NewCamera2D()
	target := Vector2{X: 10, Y: 5}
	
	originalPosition := camera.Position
	
	// Test invalid follow speed
	camera.FollowTarget(target, -1.0, 1.0)
	assert.Equal(t, originalPosition, camera.Position)
	
	// Test invalid delta time
	camera.FollowTarget(target, 5.0, -1.0)
	assert.Equal(t, originalPosition, camera.Position)
}