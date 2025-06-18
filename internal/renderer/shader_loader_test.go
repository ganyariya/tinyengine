package renderer

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadShaderFromFile(t *testing.T) {
	// Arrange
	testShaderSource := `#version 410 core
void main() { gl_Position = vec4(0.0); }`

	// Create temporary file for testing
	tempDir := t.TempDir()
	shaderPath := filepath.Join(tempDir, "test.vert")

	// Write test shader to file
	err := writeStringToFile(shaderPath, testShaderSource)
	assert.NoError(t, err)

	// Act
	source, err := LoadShaderFromFile(shaderPath)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, source, "#version 410 core")
	assert.Contains(t, source, "gl_Position")
}

func TestLoadShaderFromFile_FileNotExists(t *testing.T) {
	// Act
	_, err := LoadShaderFromFile("nonexistent_file.vert")

	// Assert
	assert.Error(t, err)
}

func TestCreateShaderFromFiles(t *testing.T) {
	// OpenGL環境が必要なテストのため、CI環境ではスキップ
	t.Skip("OpenGL関数を呼び出すためOpenGLコンテキストが必要 - 統合テストで実施")
}

func TestGetBuiltinShaderPaths(t *testing.T) {
	// Act
	vertPath, fragPath := GetBuiltinShaderPaths("simple")

	// Assert
	assert.Contains(t, vertPath, "simple.vert")
	assert.Contains(t, fragPath, "simple.frag")
}
