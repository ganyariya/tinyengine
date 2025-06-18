package renderer

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRenderer はテスト用のRendererモック
type MockRenderer struct {
	mock.Mock
}

func (m *MockRenderer) Clear() {
	m.Called()
}

func (m *MockRenderer) Present() {
	m.Called()
}

func (m *MockRenderer) DrawRectangle(x, y, width, height float32) {
	m.Called(x, y, width, height)
}

func TestMockRenderer_Clear(t *testing.T) {
	// Arrange
	mockRenderer := new(MockRenderer)
	mockRenderer.On("Clear").Return()
	
	// Act
	mockRenderer.Clear()
	
	// Assert
	mockRenderer.AssertExpectations(t)
}

func TestMockRenderer_Present(t *testing.T) {
	// Arrange
	mockRenderer := new(MockRenderer)
	mockRenderer.On("Present").Return()
	
	// Act
	mockRenderer.Present()
	
	// Assert
	mockRenderer.AssertExpectations(t)
}

func TestMockRenderer_DrawRectangle(t *testing.T) {
	// Arrange
	mockRenderer := new(MockRenderer)
	x, y, width, height := float32(10), float32(20), float32(100), float32(50)
	mockRenderer.On("DrawRectangle", x, y, width, height).Return()
	
	// Act
	mockRenderer.DrawRectangle(x, y, width, height)
	
	// Assert
	mockRenderer.AssertExpectations(t)
}

func TestNewBaseRenderer(t *testing.T) {
	// Arrange
	width, height := 800, 600
	
	// Act
	renderer := NewBaseRenderer(width, height)
	
	// Assert
	assert.NotNil(t, renderer)
}

func TestBaseRenderer_GetSize(t *testing.T) {
	// Arrange
	width, height := 800, 600
	renderer := NewBaseRenderer(width, height).(*BaseRenderer)
	
	// Act
	w, h := renderer.GetSize()
	
	// Assert
	assert.Equal(t, width, w)
	assert.Equal(t, height, h)
}

func TestBaseRenderer_Clear(t *testing.T) {
	// Arrange
	renderer := NewBaseRenderer(800, 600)
	
	// Act & Assert
	// Clear()は基本実装では何もしないため、パニックしないことを確認
	assert.NotPanics(t, func() {
		renderer.Clear()
	})
}

func TestBaseRenderer_Present(t *testing.T) {
	// Arrange
	renderer := NewBaseRenderer(800, 600)
	
	// Act & Assert
	// Present()は基本実装では何もしないため、パニックしないことを確認
	assert.NotPanics(t, func() {
		renderer.Present()
	})
}

func TestBaseRenderer_DrawRectangle(t *testing.T) {
	// Arrange
	renderer := NewBaseRenderer(800, 600)
	
	// Act & Assert
	// DrawRectangle()は基本実装では何もしないため、パニックしないことを確認
	assert.NotPanics(t, func() {
		renderer.DrawRectangle(10, 20, 100, 50)
	})
}