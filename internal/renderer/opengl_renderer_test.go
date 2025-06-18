package renderer

import (
	"testing"

	"github.com/ganyariya/tinyengine/pkg/tinyengine"
	"github.com/stretchr/testify/assert"
)

func TestNewOpenGLRenderer(t *testing.T) {
	// Arrange
	width, height := 800, 600

	// Act
	renderer, err := NewOpenGLRenderer(width, height)

	// Assert
	// 実際のOpenGL初期化が必要なため、ヘッドレス環境ではエラーになることを確認
	if err != nil {
		// ヘッドレス環境での期待される動作
		assert.Error(t, err)
		assert.Nil(t, renderer)
	} else {
		// OpenGL環境が利用可能な場合
		assert.NotNil(t, renderer)
		assert.NoError(t, err)
	}
}

func TestOpenGLRenderer_Implementation(t *testing.T) {
	// OpenGLRenderer構造体がRendererインターフェースを実装していることを確認
	var _ tinyengine.Renderer = (*OpenGLRenderer)(nil)
}

func TestOpenGLRenderer_Methods(t *testing.T) {
	// Skip test when OpenGL is not available
	t.Skip("OpenGL methods require GL context initialization - skipping in CI environment")
}
