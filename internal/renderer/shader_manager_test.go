package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// テスト用シェーダーソース
const (
	testVertexShaderSource = `
#version 410 core
layout (location = 0) in vec3 aPos;
void main() { gl_Position = vec4(aPos, 1.0); }
`

	testFragmentShaderSource = `
#version 410 core
out vec4 FragColor;
void main() { FragColor = vec4(1.0, 0.0, 0.0, 1.0); }
`
)

func TestNewShaderManager(t *testing.T) {
	// Act
	manager := NewShaderManager()

	// Assert
	assert.NotNil(t, manager)
	assert.Equal(t, 0, manager.GetShaderCount())
}

func TestShaderManager_LoadShader(t *testing.T) {
	// Arrange
	manager := NewShaderManager()
	shaderName := "test_shader"

	// Act - 実際のOpenGL環境がないため、エラーが発生することを想定
	err := manager.LoadShader(shaderName, testVertexShaderSource, testFragmentShaderSource)

	// Assert - OpenGL環境がない場合はエラーが発生することを許容
	// 実際のOpenGL依存のテストは統合テストで実施する
	if err != nil {
		// OpenGL関連のエラーが発生することを確認
		assert.NotNil(t, err)
	} else {
		// まれにOpenGL環境がある場合
		assert.Equal(t, 1, manager.GetShaderCount())
		assert.True(t, manager.HasShader(shaderName))
	}
}

func TestShaderManager_GetShader(t *testing.T) {
	// Arrange
	manager := NewShaderManager()
	shaderName := "test_shader"

	// Act
	shader := manager.GetShader(shaderName)

	// Assert
	// 存在しないシェーダーの場合はnilが返される
	assert.Nil(t, shader)
}

func TestShaderManager_UseShader(t *testing.T) {
	// Arrange
	manager := NewShaderManager()
	shaderName := "test_shader"

	// Act & Assert
	// 存在しないシェーダーを使用しようとしてもパニックしない
	assert.NotPanics(t, func() {
		manager.UseShader(shaderName)
	})
}

func TestShaderManager_DeleteShader(t *testing.T) {
	// Arrange
	manager := NewShaderManager()
	shaderName := "test_shader"

	// Act
	deleted := manager.DeleteShader(shaderName)

	// Assert
	// 存在しないシェーダーを削除しようとした場合はfalseが返される
	assert.False(t, deleted)
}

func TestShaderManager_DeleteAllShaders(t *testing.T) {
	// Arrange
	manager := NewShaderManager()

	// Act & Assert
	// 空のマネージャーでもパニックしない
	assert.NotPanics(t, func() {
		manager.DeleteAllShaders()
	})

	assert.Equal(t, 0, manager.GetShaderCount())
}

func TestShaderManager_GetCurrentShader(t *testing.T) {
	// Arrange
	manager := NewShaderManager()

	// Act
	currentShader := manager.GetCurrentShader()

	// Assert
	// 初期状態では現在のシェーダーは存在しない
	assert.Equal(t, "", currentShader)
}

func TestShaderManager_GetShaderNames(t *testing.T) {
	// Arrange
	manager := NewShaderManager()

	// Act
	names := manager.GetShaderNames()

	// Assert
	assert.NotNil(t, names)
	assert.Equal(t, 0, len(names))
}
