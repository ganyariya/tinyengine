package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommandQueue(t *testing.T) {
	// Act
	queue := NewCommandQueue()

	// Assert
	assert.NotNil(t, queue)
	assert.Equal(t, 0, queue.Size())
}

func TestCommandQueue_AddRectangleCommand(t *testing.T) {
	// Arrange
	queue := NewCommandQueue()
	x, y, width, height := float32(10), float32(20), float32(100), float32(50)

	// Act
	queue.AddRectangleCommand(x, y, width, height)

	// Assert
	assert.Equal(t, 1, queue.Size())
}

func TestCommandQueue_AddClearCommand(t *testing.T) {
	// Arrange
	queue := NewCommandQueue()

	// Act
	queue.AddClearCommand()

	// Assert
	assert.Equal(t, 1, queue.Size())
}

func TestCommandQueue_Execute(t *testing.T) {
	// Arrange
	queue := NewCommandQueue()
	mockRenderer := new(MockRenderer)

	// Set up expectations
	mockRenderer.On("Clear").Return()
	mockRenderer.On("DrawRectangle", float32(10), float32(20), float32(100), float32(50)).Return()

	// Add commands
	queue.AddClearCommand()
	queue.AddRectangleCommand(10, 20, 100, 50)

	// Act
	queue.Execute(mockRenderer)

	// Assert
	mockRenderer.AssertExpectations(t)
}

func TestCommandQueue_Clear(t *testing.T) {
	// Arrange
	queue := NewCommandQueue()
	queue.AddClearCommand()
	queue.AddRectangleCommand(10, 20, 100, 50)

	// Act
	queue.Clear()

	// Assert
	assert.Equal(t, 0, queue.Size())
}
