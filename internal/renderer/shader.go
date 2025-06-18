package renderer

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Shader はOpenGLシェーダープログラムを管理する
type Shader struct {
	backend          OpenGLBackend
	programID        uint32
	vertexShaderID   uint32
	fragmentShaderID uint32
}

// NewShader は新しいShaderを作成する
func NewShader(backend OpenGLBackend) *Shader {
	return &Shader{
		backend:          backend,
		programID:        0,
		vertexShaderID:   0,
		fragmentShaderID: 0,
	}
}

// LoadVertexShader は頂点シェーダーを読み込む
func (s *Shader) LoadVertexShader(source string) error {
	return s.loadShader(source, gl.VERTEX_SHADER, &s.vertexShaderID)
}

// LoadFragmentShader はフラグメントシェーダーを読み込む
func (s *Shader) LoadFragmentShader(source string) error {
	return s.loadShader(source, gl.FRAGMENT_SHADER, &s.fragmentShaderID)
}

// loadShader は指定された種類のシェーダーを読み込む
func (s *Shader) loadShader(source string, shaderType uint32, shaderID *uint32) error {
	// シェーダー作成
	*shaderID = s.backend.CreateShader(shaderType)
	if *shaderID == 0 {
		return fmt.Errorf("failed to create shader")
	}

	// ソースコード設定
	s.backend.ShaderSource(*shaderID, source)

	// コンパイル
	s.backend.CompileShader(*shaderID)

	// コンパイル結果確認
	success := s.backend.GetShaderiv(*shaderID, gl.COMPILE_STATUS)
	if success == gl.FALSE {
		log := s.backend.GetShaderInfoLog(*shaderID)
		s.backend.DeleteShader(*shaderID)
		*shaderID = 0
		return fmt.Errorf("shader compilation failed: %s", log)
	}

	return nil
}

// LinkProgram はシェーダープログラムをリンクする
func (s *Shader) LinkProgram() error {
	// 頂点・フラグメントシェーダーがロードされているかチェック
	if s.vertexShaderID == 0 {
		return fmt.Errorf("vertex shader not loaded")
	}
	if s.fragmentShaderID == 0 {
		return fmt.Errorf("fragment shader not loaded")
	}

	// プログラム作成
	s.programID = s.backend.CreateProgram()
	if s.programID == 0 {
		return fmt.Errorf("failed to create shader program")
	}

	// シェーダーをアタッチ
	s.backend.AttachShader(s.programID, s.vertexShaderID)
	s.backend.AttachShader(s.programID, s.fragmentShaderID)

	// リンク
	s.backend.LinkProgram(s.programID)

	// リンク結果確認
	success := s.backend.GetProgramiv(s.programID, gl.LINK_STATUS)
	if success == gl.FALSE {
		log := s.backend.GetProgramInfoLog(s.programID)
		return fmt.Errorf("shader program linking failed: %s", log)
	}

	// シェーダーをデタッチ・削除（プログラムにリンク済み）
	s.backend.DetachShader(s.programID, s.vertexShaderID)
	s.backend.DetachShader(s.programID, s.fragmentShaderID)
	s.backend.DeleteShader(s.vertexShaderID)
	s.backend.DeleteShader(s.fragmentShaderID)

	s.vertexShaderID = 0
	s.fragmentShaderID = 0

	return nil
}

// Use はシェーダープログラムを使用する
func (s *Shader) Use() {
	if s.programID != 0 {
		s.backend.UseProgram(s.programID)
	}
}

// Delete はシェーダープログラムを削除する
func (s *Shader) Delete() {
	if s.programID != 0 {
		s.backend.DeleteProgram(s.programID)
		s.programID = 0
	}

	// 個別シェーダーも削除
	if s.vertexShaderID != 0 {
		s.backend.DeleteShader(s.vertexShaderID)
		s.vertexShaderID = 0
	}
	if s.fragmentShaderID != 0 {
		s.backend.DeleteShader(s.fragmentShaderID)
		s.fragmentShaderID = 0
	}
}

// GetProgramID はシェーダープログラムIDを取得する
func (s *Shader) GetProgramID() uint32 {
	return s.programID
}

// GetUniformLocation はユニフォーム変数の位置を取得する
func (s *Shader) GetUniformLocation(name string) int32 {
	if s.programID == 0 {
		return -1
	}

	return s.backend.GetUniformLocation(s.programID, name)
}

// SetUniformMat4 は4x4行列のユニフォーム変数を設定する
func (s *Shader) SetUniformMat4(location int32, matrix [16]float32) {
	if location >= 0 {
		s.backend.UniformMatrix4fv(location, matrix)
	}
}

// SetUniformVec3 は3次元ベクトルのユニフォーム変数を設定する
func (s *Shader) SetUniformVec3(location int32, vector [3]float32) {
	if location >= 0 {
		s.backend.Uniform3fv(location, vector)
	}
}

// SetUniformFloat は浮動小数点数のユニフォーム変数を設定する
func (s *Shader) SetUniformFloat(location int32, value float32) {
	if location >= 0 {
		s.backend.Uniform1f(location, value)
	}
}

// SetUniformInt は整数のユニフォーム変数を設定する
func (s *Shader) SetUniformInt(location int32, value int32) {
	if location >= 0 {
		s.backend.Uniform1i(location, value)
	}
}
