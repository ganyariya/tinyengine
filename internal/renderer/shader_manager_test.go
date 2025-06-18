package renderer

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
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
	
	// Act
	err := manager.LoadShader(shaderName, testVertexShaderSource, testFragmentShaderSource)
	
	// Assert
	if err != nil {
		// OpenGL環境がない場合はエラーが発生することを許容
		assert.Contains(t, err.Error(), "OpenGL")
	} else {
		// OpenGL環境がある場合は正常にシェーダーが読み込まれる
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