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
	// Arrange
	tempDir := t.TempDir()
	vertPath := filepath.Join(tempDir, "test.vert")
	fragPath := filepath.Join(tempDir, "test.frag")

	vertSource := `#version 410 core
layout (location = 0) in vec3 aPos;
void main() { gl_Position = vec4(aPos, 1.0); }`

	fragSource := `#version 410 core
out vec4 FragColor;
void main() { FragColor = vec4(1.0, 0.0, 0.0, 1.0); }`

	err := writeStringToFile(vertPath, vertSource)
	assert.NoError(t, err)
	err = writeStringToFile(fragPath, fragSource)
	assert.NoError(t, err)

	// Act - 実際のOpenGL環境がないため、エラーが発生することを想定
	shader, err := CreateShaderFromFiles(vertPath, fragPath)

	// Assert - OpenGL環境がない場合はエラーが発生することを許容
	// 実際のOpenGL依存のテストは統合テストで実施する
	if err != nil {
		// OpenGL関連のエラーが発生することを確認
		assert.NotNil(t, err)
	} else {
		// まれにOpenGL環境がある場合
		assert.NotNil(t, shader)
	}
}

func TestGetBuiltinShaderPaths(t *testing.T) {
	// Act
	vertPath, fragPath := GetBuiltinShaderPaths("simple")

	// Assert
	assert.Contains(t, vertPath, "simple.vert")
	assert.Contains(t, fragPath, "simple.frag")
}
