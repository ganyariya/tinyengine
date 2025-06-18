package renderer

// OpenGLBackend はOpenGL呼び出しを抽象化するインターフェース
type OpenGLBackend interface {
	// シェーダー関連
	CreateShader(shaderType uint32) uint32
	ShaderSource(shader uint32, source string)
	CompileShader(shader uint32)
	GetShaderiv(shader uint32, pname uint32) int32
	GetShaderInfoLog(shader uint32) string
	DeleteShader(shader uint32)

	// プログラム関連
	CreateProgram() uint32
	AttachShader(program, shader uint32)
	DetachShader(program, shader uint32)
	LinkProgram(program uint32)
	GetProgramiv(program uint32, pname uint32) int32
	GetProgramInfoLog(program uint32) string
	UseProgram(program uint32)
	DeleteProgram(program uint32)

	// ユニフォーム関連
	GetUniformLocation(program uint32, name string) int32
	UniformMatrix4fv(location int32, matrix [16]float32)
	Uniform3fv(location int32, vector [3]float32)
	Uniform1f(location int32, value float32)
	Uniform1i(location int32, value int32)
}
