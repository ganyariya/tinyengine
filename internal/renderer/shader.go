package renderer

import (
	"fmt"
	"os"
	"strings"
	"unsafe"
	
	"github.com/go-gl/gl/v4.1-core/gl"
)

// Shader はOpenGLシェーダープログラムを管理する
type Shader struct {
	programID      uint32
	vertexShaderID uint32
	fragmentShaderID uint32
}

// NewShader は新しいShaderを作成する
func NewShader() *Shader {
	return &Shader{
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
	// OpenGL初期化チェック
	if !isOpenGLInitialized() {
		return fmt.Errorf("OpenGL is not initialized")
	}
	
	// シェーダー作成
	*shaderID = gl.CreateShader(shaderType)
	if *shaderID == 0 {
		return fmt.Errorf("failed to create shader")
	}
	
	// ソースコード設定
	cSource, free := gl.Strs(source + "\x00")
	defer free()
	gl.ShaderSource(*shaderID, 1, cSource, nil)
	
	// コンパイル
	gl.CompileShader(*shaderID)
	
	// コンパイル結果確認
	var success int32
	gl.GetShaderiv(*shaderID, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(*shaderID, gl.INFO_LOG_LENGTH, &logLength)
		
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(*shaderID, logLength, nil, gl.Str(log))
		
		gl.DeleteShader(*shaderID)
		*shaderID = 0
		
		return fmt.Errorf("shader compilation failed: %s", log)
	}
	
	return nil
}

// LinkProgram はシェーダープログラムをリンクする
func (s *Shader) LinkProgram() error {
	// OpenGL初期化チェック
	if !isOpenGLInitialized() {
		return fmt.Errorf("OpenGL is not initialized")
	}
	
	// 頂点・フラグメントシェーダーがロードされているかチェック
	if s.vertexShaderID == 0 {
		return fmt.Errorf("vertex shader not loaded")
	}
	if s.fragmentShaderID == 0 {
		return fmt.Errorf("fragment shader not loaded")
	}
	
	// プログラム作成
	s.programID = gl.CreateProgram()
	if s.programID == 0 {
		return fmt.Errorf("failed to create shader program")
	}
	
	// シェーダーをアタッチ
	gl.AttachShader(s.programID, s.vertexShaderID)
	gl.AttachShader(s.programID, s.fragmentShaderID)
	
	// リンク
	gl.LinkProgram(s.programID)
	
	// リンク結果確認
	var success int32
	gl.GetProgramiv(s.programID, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.programID, gl.INFO_LOG_LENGTH, &logLength)
		
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.programID, logLength, nil, gl.Str(log))
		
		return fmt.Errorf("shader program linking failed: %s", log)
	}
	
	// シェーダーをデタッチ・削除（プログラムにリンク済み）
	gl.DetachShader(s.programID, s.vertexShaderID)
	gl.DetachShader(s.programID, s.fragmentShaderID)
	gl.DeleteShader(s.vertexShaderID)
	gl.DeleteShader(s.fragmentShaderID)
	
	s.vertexShaderID = 0
	s.fragmentShaderID = 0
	
	return nil
}

// Use はシェーダープログラムを使用する
func (s *Shader) Use() {
	if s.programID != 0 && isOpenGLInitialized() {
		gl.UseProgram(s.programID)
	}
}

// Delete はシェーダープログラムを削除する
func (s *Shader) Delete() {
	if s.programID != 0 && isOpenGLInitialized() {
		gl.DeleteProgram(s.programID)
		s.programID = 0
	}
	
	// 個別シェーダーも削除
	if s.vertexShaderID != 0 && isOpenGLInitialized() {
		gl.DeleteShader(s.vertexShaderID)
		s.vertexShaderID = 0
	}
	if s.fragmentShaderID != 0 && isOpenGLInitialized() {
		gl.DeleteShader(s.fragmentShaderID)
		s.fragmentShaderID = 0
	}
}

// GetProgramID はシェーダープログラムIDを取得する
func (s *Shader) GetProgramID() uint32 {
	return s.programID
}

// GetUniformLocation はユニフォーム変数の位置を取得する
func (s *Shader) GetUniformLocation(name string) int32 {
	if s.programID == 0 || !isOpenGLInitialized() {
		return -1
	}
	
	cName := gl.Str(name + "\x00")
	return gl.GetUniformLocation(s.programID, cName)
}

// SetUniformMat4 は4x4行列のユニフォーム変数を設定する
func (s *Shader) SetUniformMat4(location int32, matrix [16]float32) {
	if location >= 0 && isOpenGLInitialized() {
		gl.UniformMatrix4fv(location, 1, false, (*float32)(unsafe.Pointer(&matrix[0])))
	}
}

// SetUniformVec3 は3次元ベクトルのユニフォーム変数を設定する
func (s *Shader) SetUniformVec3(location int32, vector [3]float32) {
	if location >= 0 && isOpenGLInitialized() {
		gl.Uniform3fv(location, 1, (*float32)(unsafe.Pointer(&vector[0])))
	}
}

// SetUniformFloat は浮動小数点数のユニフォーム変数を設定する
func (s *Shader) SetUniformFloat(location int32, value float32) {
	if location >= 0 && isOpenGLInitialized() {
		gl.Uniform1f(location, value)
	}
}

// SetUniformInt は整数のユニフォーム変数を設定する
func (s *Shader) SetUniformInt(location int32, value int32) {
	if location >= 0 && isOpenGLInitialized() {
		gl.Uniform1i(location, value)
	}
}

// isOpenGLInitialized はOpenGLが初期化されているかを簡易チェックする
func isOpenGLInitialized() bool {
	// CI環境やテスト環境ではOpenGLが利用できない場合が多い
	if os.Getenv("CI") != "" {
		return false
	}
	
	// テストコンテキストかどうかをチェック
	if os.Getenv("GO_TEST") != "" {
		return false
	}
	
	// 実際のアプリケーション実行時はtrueを返す
	// （gl.Init()が事前に呼び出されていることを前提）
	return true
}