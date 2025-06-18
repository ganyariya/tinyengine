package renderer

import (
	"testing"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// テスト用の基本的なシェーダーソースコード
const (
	validVertexShaderSource = `
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 vertexColor;

void main() {
    gl_Position = vec4(aPos, 1.0);
    vertexColor = aColor;
}
`

	validFragmentShaderSource = `
#version 410 core
in vec3 vertexColor;
out vec4 FragColor;

uniform float alpha;

void main() {
    FragColor = vec4(vertexColor, alpha);
}
`

	invalidShaderSource = `
#version 410 core
ERROR This is invalid syntax
`
)

func TestNewShader_WithValidBackend(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()

	// Act
	shader := NewShader(mockBackend)

	// Assert
	assert.NotNil(t, shader)
	assert.Equal(t, uint32(0), shader.GetProgramID())
}

func TestShader_LoadVertexShader_Success(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)

	// モックの設定：正常ケース
	mockBackend.On("CreateShader", uint32(gl.VERTEX_SHADER)).Return(uint32(1))
	mockBackend.On("ShaderSource", uint32(1), validVertexShaderSource).Return()
	mockBackend.On("CompileShader", uint32(1)).Return()
	mockBackend.On("GetShaderiv", uint32(1), uint32(gl.COMPILE_STATUS)).Return(int32(1))

	// Act
	err := shader.LoadVertexShader(validVertexShaderSource)

	// Assert
	assert.NoError(t, err)
	mockBackend.AssertExpectations(t)

	// 内部状態の確認
	mockShader := mockBackend.GetShader(1)
	assert.NotNil(t, mockShader)
	assert.Equal(t, validVertexShaderSource, mockShader.Source)
	assert.True(t, mockShader.Compiled)
}

func TestShader_LoadVertexShader_CompilationError(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)

	// モックの設定：コンパイルエラーケース
	mockBackend.On("CreateShader", uint32(gl.VERTEX_SHADER)).Return(uint32(1))
	mockBackend.On("ShaderSource", uint32(1), invalidShaderSource).Return()
	mockBackend.On("CompileShader", uint32(1)).Return()
	mockBackend.On("GetShaderiv", uint32(1), uint32(gl.COMPILE_STATUS)).Return(int32(0))
	mockBackend.On("GetShaderInfoLog", uint32(1)).Return("Mock compile error")
	mockBackend.On("DeleteShader", uint32(1)).Return()

	// Act
	err := shader.LoadVertexShader(invalidShaderSource)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "shader compilation failed")
	assert.Contains(t, err.Error(), "Mock compile error")
	mockBackend.AssertExpectations(t)
}

func TestShader_LoadFragmentShader_Success(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)

	// モックの設定
	mockBackend.On("CreateShader", uint32(gl.FRAGMENT_SHADER)).Return(uint32(2))
	mockBackend.On("ShaderSource", uint32(2), validFragmentShaderSource).Return()
	mockBackend.On("CompileShader", uint32(2)).Return()
	mockBackend.On("GetShaderiv", uint32(2), uint32(gl.COMPILE_STATUS)).Return(int32(1))

	// Act
	err := shader.LoadFragmentShader(validFragmentShaderSource)

	// Assert
	assert.NoError(t, err)
	mockBackend.AssertExpectations(t)
}

func TestShader_LinkProgram_Success(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)

	// 頂点シェーダーのロード（内部状態を設定）
	shader.vertexShaderID = 1
	shader.fragmentShaderID = 2

	// モックの設定
	mockBackend.On("CreateProgram").Return(uint32(3))
	mockBackend.On("AttachShader", uint32(3), uint32(1)).Return()
	mockBackend.On("AttachShader", uint32(3), uint32(2)).Return()
	mockBackend.On("LinkProgram", uint32(3)).Return()
	mockBackend.On("GetProgramiv", uint32(3), uint32(gl.LINK_STATUS)).Return(int32(1))
	mockBackend.On("DetachShader", uint32(3), uint32(1)).Return()
	mockBackend.On("DetachShader", uint32(3), uint32(2)).Return()
	mockBackend.On("DeleteShader", uint32(1)).Return()
	mockBackend.On("DeleteShader", uint32(2)).Return()

	// Act
	err := shader.LinkProgram()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint32(3), shader.GetProgramID())
	assert.Equal(t, uint32(0), shader.vertexShaderID)   // クリーンアップされている
	assert.Equal(t, uint32(0), shader.fragmentShaderID) // クリーンアップされている
	mockBackend.AssertExpectations(t)
}

func TestShader_LinkProgram_MissingVertexShader(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)

	// フラグメントシェーダーのみロード
	shader.fragmentShaderID = 2

	// Act
	err := shader.LinkProgram()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vertex shader not loaded")
}

func TestShader_LinkProgram_LinkError(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)

	shader.vertexShaderID = 1
	shader.fragmentShaderID = 2

	// モックの設定：リンクエラーケース
	mockBackend.On("CreateProgram").Return(uint32(3))
	mockBackend.On("AttachShader", uint32(3), uint32(1)).Return()
	mockBackend.On("AttachShader", uint32(3), uint32(2)).Return()
	mockBackend.On("LinkProgram", uint32(3)).Return()
	mockBackend.On("GetProgramiv", uint32(3), uint32(gl.LINK_STATUS)).Return(int32(0))
	mockBackend.On("GetProgramInfoLog", uint32(3)).Return("Mock link error")

	// Act
	err := shader.LinkProgram()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "shader program linking failed")
	assert.Contains(t, err.Error(), "Mock link error")
	mockBackend.AssertExpectations(t)
}

func TestShader_Use(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)
	shader.programID = 123

	mockBackend.On("UseProgram", uint32(123)).Return()

	// Act
	shader.Use()

	// Assert
	mockBackend.AssertExpectations(t)
}

func TestShader_Delete(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)
	shader.programID = 123
	shader.vertexShaderID = 1
	shader.fragmentShaderID = 2

	mockBackend.On("DeleteProgram", uint32(123)).Return()
	mockBackend.On("DeleteShader", uint32(1)).Return()
	mockBackend.On("DeleteShader", uint32(2)).Return()

	// Act
	shader.Delete()

	// Assert
	assert.Equal(t, uint32(0), shader.programID)
	assert.Equal(t, uint32(0), shader.vertexShaderID)
	assert.Equal(t, uint32(0), shader.fragmentShaderID)
	mockBackend.AssertExpectations(t)
}

func TestShader_GetUniformLocation(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)
	shader.programID = 123

	mockBackend.On("GetUniformLocation", uint32(123), "testUniform").Return(int32(5))

	// Act
	location := shader.GetUniformLocation("testUniform")

	// Assert
	assert.Equal(t, int32(5), location)
	mockBackend.AssertExpectations(t)
}

func TestShader_GetUniformLocation_NoProgramID(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)
	// programID = 0 (初期値)

	// Act
	location := shader.GetUniformLocation("testUniform")

	// Assert
	assert.Equal(t, int32(-1), location)
	// バックエンドは呼び出されない
}

func TestShader_SetUniformFloat(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)

	mockBackend.On("Uniform1f", int32(5), float32(1.5)).Return()

	// Act
	shader.SetUniformFloat(5, 1.5)

	// Assert
	mockBackend.AssertExpectations(t)
}

func TestShader_SetUniformFloat_InvalidLocation(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)

	// Act
	shader.SetUniformFloat(-1, 1.5)

	// Assert
	// バックエンドは呼び出されない（invalid location）
	mockBackend.AssertNotCalled(t, "Uniform1f", mock.Anything, mock.Anything)
}

func TestShader_FullWorkflow_Integration(t *testing.T) {
	// Arrange
	mockBackend := NewMockOpenGLBackend()
	shader := NewShader(mockBackend)

	// 全フローのモック設定
	// 頂点シェーダー
	mockBackend.On("CreateShader", uint32(gl.VERTEX_SHADER)).Return(uint32(1))
	mockBackend.On("ShaderSource", uint32(1), validVertexShaderSource).Return()
	mockBackend.On("CompileShader", uint32(1)).Return()
	mockBackend.On("GetShaderiv", uint32(1), uint32(gl.COMPILE_STATUS)).Return(int32(1))

	// フラグメントシェーダー
	mockBackend.On("CreateShader", uint32(gl.FRAGMENT_SHADER)).Return(uint32(2))
	mockBackend.On("ShaderSource", uint32(2), validFragmentShaderSource).Return()
	mockBackend.On("CompileShader", uint32(2)).Return()
	mockBackend.On("GetShaderiv", uint32(2), uint32(gl.COMPILE_STATUS)).Return(int32(1))

	// プログラムリンク
	mockBackend.On("CreateProgram").Return(uint32(3))
	mockBackend.On("AttachShader", uint32(3), uint32(1)).Return()
	mockBackend.On("AttachShader", uint32(3), uint32(2)).Return()
	mockBackend.On("LinkProgram", uint32(3)).Return()
	mockBackend.On("GetProgramiv", uint32(3), uint32(gl.LINK_STATUS)).Return(int32(1))
	mockBackend.On("DetachShader", uint32(3), uint32(1)).Return()
	mockBackend.On("DetachShader", uint32(3), uint32(2)).Return()
	mockBackend.On("DeleteShader", uint32(1)).Return()
	mockBackend.On("DeleteShader", uint32(2)).Return()

	// 使用とユニフォーム設定
	mockBackend.On("UseProgram", uint32(3)).Return()
	mockBackend.On("GetUniformLocation", uint32(3), "alpha").Return(int32(0))
	mockBackend.On("Uniform1f", int32(0), float32(0.5)).Return()

	// クリーンアップ
	mockBackend.On("DeleteProgram", uint32(3)).Return()

	// Act: 完全なワークフロー
	err := shader.LoadVertexShader(validVertexShaderSource)
	assert.NoError(t, err)

	err = shader.LoadFragmentShader(validFragmentShaderSource)
	assert.NoError(t, err)

	err = shader.LinkProgram()
	assert.NoError(t, err)

	shader.Use()

	location := shader.GetUniformLocation("alpha")
	assert.Equal(t, int32(0), location)

	shader.SetUniformFloat(location, 0.5)

	shader.Delete()

	// Assert
	mockBackend.AssertExpectations(t)
	assert.Equal(t, uint32(0), shader.GetProgramID())
}
