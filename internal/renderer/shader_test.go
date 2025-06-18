package renderer

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

// テスト用の基本的なシェーダーソースコード
const (
	testVertexShaderSource = `
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 vertexColor;

void main() {
    gl_Position = vec4(aPos, 1.0);
    vertexColor = aColor;
}
`

	testFragmentShaderSource = `
#version 410 core
in vec3 vertexColor;
out vec4 FragColor;

void main() {
    FragColor = vec4(vertexColor, 1.0);
}
`
)

func TestNewShader(t *testing.T) {
	// Act
	shader := NewShader()
	
	// Assert
	assert.NotNil(t, shader)
	assert.Equal(t, uint32(0), shader.GetProgramID())
}

func TestShader_LoadVertexShader(t *testing.T) {
	// Arrange
	shader := NewShader()
	
	// Act
	err := shader.LoadVertexShader(testVertexShaderSource)
	
	// Assert
	// ヘッドレス環境ではOpenGLが利用できないため、エラーまたは成功の両方を許容
	if err != nil {
		assert.Contains(t, err.Error(), "OpenGL")
	} else {
		// OpenGL環境が利用可能な場合
		assert.NoError(t, err)
	}
}

func TestShader_LoadFragmentShader(t *testing.T) {
	// Arrange
	shader := NewShader()
	
	// Act
	err := shader.LoadFragmentShader(testFragmentShaderSource)
	
	// Assert
	// ヘッドレス環境ではOpenGLが利用できないため、エラーまたは成功の両方を許容
	if err != nil {
		assert.Contains(t, err.Error(), "OpenGL")
	} else {
		// OpenGL環境が利用可能な場合
		assert.NoError(t, err)
	}
}

func TestShader_LinkProgram(t *testing.T) {
	// Arrange
	shader := NewShader()
	
	// Act
	err := shader.LinkProgram()
	
	// Assert
	// シェーダーが読み込まれていない状態でのリンクはエラーになるべき
	assert.Error(t, err)
}

func TestShader_Use(t *testing.T) {
	// Arrange
	shader := NewShader()
	
	// Act & Assert
	// OpenGL環境なしでも基本的な使用メソッドがパニックしないことを確認
	assert.NotPanics(t, func() {
		shader.Use()
	})
}

func TestShader_Delete(t *testing.T) {
	// Arrange
	shader := NewShader()
	
	// Act & Assert
	// Delete操作がパニックしないことを確認
	assert.NotPanics(t, func() {
		shader.Delete()
	})
}

func TestShader_GetUniformLocation(t *testing.T) {
	// Arrange
	shader := NewShader()
	uniformName := "testUniform"
	
	// Act
	location := shader.GetUniformLocation(uniformName)
	
	// Assert
	// プログラムが作成されていない状態では-1が返されるべき
	assert.Equal(t, int32(-1), location)
}

func TestShader_SetUniformMat4(t *testing.T) {
	// Arrange
	shader := NewShader()
	matrix := [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	
	// Act & Assert
	// 無効なlocationでも例外が発生しないことを確認
	assert.NotPanics(t, func() {
		shader.SetUniformMat4(-1, matrix)
	})
}

func TestShader_SetUniformVec3(t *testing.T) {
	// Arrange
	shader := NewShader()
	vector := [3]float32{1.0, 0.5, 0.2}
	
	// Act & Assert
	// 無効なlocationでも例外が発生しないことを確認
	assert.NotPanics(t, func() {
		shader.SetUniformVec3(-1, vector)
	})
}

func TestShader_FullWorkflow(t *testing.T) {
	// Skip in CI environment without OpenGL
	t.Skip("Full shader workflow requires OpenGL context - skipping in CI environment")
}