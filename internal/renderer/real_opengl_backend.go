package renderer

import (
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// RealOpenGLBackend は実際のOpenGL関数を呼び出す実装
type RealOpenGLBackend struct{}

// NewRealOpenGLBackend は新しいRealOpenGLBackendを作成する
func NewRealOpenGLBackend() OpenGLBackend {
	return &RealOpenGLBackend{}
}

// CreateShader は新しいシェーダーオブジェクトを作成する
func (b *RealOpenGLBackend) CreateShader(shaderType uint32) uint32 {
	return gl.CreateShader(shaderType)
}

// ShaderSource はシェーダーオブジェクトにソースコードを設定する
func (b *RealOpenGLBackend) ShaderSource(shader uint32, source string) {
	cSource, free := gl.Strs(source + "\x00")
	defer free()
	gl.ShaderSource(shader, 1, cSource, nil)
}

// CompileShader はシェーダーをコンパイルする
func (b *RealOpenGLBackend) CompileShader(shader uint32) {
	gl.CompileShader(shader)
}

// GetShaderiv はシェーダーパラメータを取得する
func (b *RealOpenGLBackend) GetShaderiv(shader uint32, pname uint32) int32 {
	var value int32
	gl.GetShaderiv(shader, pname, &value)
	return value
}

// GetShaderInfoLog はシェーダーのコンパイル情報ログを取得する
func (b *RealOpenGLBackend) GetShaderInfoLog(shader uint32) string {
	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

	if logLength == 0 {
		return ""
	}

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
	return log
}

// DeleteShader はシェーダーオブジェクトを削除する
func (b *RealOpenGLBackend) DeleteShader(shader uint32) {
	gl.DeleteShader(shader)
}

// CreateProgram は新しいプログラムオブジェクトを作成する
func (b *RealOpenGLBackend) CreateProgram() uint32 {
	return gl.CreateProgram()
}

// AttachShader はシェーダーをプログラムにアタッチする
func (b *RealOpenGLBackend) AttachShader(program, shader uint32) {
	gl.AttachShader(program, shader)
}

// DetachShader はシェーダーをプログラムからデタッチする
func (b *RealOpenGLBackend) DetachShader(program, shader uint32) {
	gl.DetachShader(program, shader)
}

// LinkProgram はプログラムをリンクする
func (b *RealOpenGLBackend) LinkProgram(program uint32) {
	gl.LinkProgram(program)
}

// GetProgramiv はプログラムパラメータを取得する
func (b *RealOpenGLBackend) GetProgramiv(program uint32, pname uint32) int32 {
	var value int32
	gl.GetProgramiv(program, pname, &value)
	return value
}

// GetProgramInfoLog はプログラムのリンク情報ログを取得する
func (b *RealOpenGLBackend) GetProgramInfoLog(program uint32) string {
	var logLength int32
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

	if logLength == 0 {
		return ""
	}

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
	return log
}

// UseProgram はプログラムを使用する
func (b *RealOpenGLBackend) UseProgram(program uint32) {
	gl.UseProgram(program)
}

// DeleteProgram はプログラムオブジェクトを削除する
func (b *RealOpenGLBackend) DeleteProgram(program uint32) {
	gl.DeleteProgram(program)
}

// GetUniformLocation はユニフォーム変数の位置を取得する
func (b *RealOpenGLBackend) GetUniformLocation(program uint32, name string) int32 {
	cName := gl.Str(name + "\x00")
	return gl.GetUniformLocation(program, cName)
}

// UniformMatrix4fv は4x4行列のユニフォーム変数を設定する
func (b *RealOpenGLBackend) UniformMatrix4fv(location int32, matrix [16]float32) {
	gl.UniformMatrix4fv(location, 1, false, (*float32)(unsafe.Pointer(&matrix[0])))
}

// Uniform3fv は3次元ベクトルのユニフォーム変数を設定する
func (b *RealOpenGLBackend) Uniform3fv(location int32, vector [3]float32) {
	gl.Uniform3fv(location, 1, (*float32)(unsafe.Pointer(&vector[0])))
}

// Uniform1f は浮動小数点数のユニフォーム変数を設定する
func (b *RealOpenGLBackend) Uniform1f(location int32, value float32) {
	gl.Uniform1f(location, value)
}

// Uniform1i は整数のユニフォーム変数を設定する
func (b *RealOpenGLBackend) Uniform1i(location int32, value int32) {
	gl.Uniform1i(location, value)
}
