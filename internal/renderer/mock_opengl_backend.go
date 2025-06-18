package renderer

import (
	"strings"

	"github.com/stretchr/testify/mock"
)

// MockOpenGLBackend はテスト用のOpenGLバックエンドモック
type MockOpenGLBackend struct {
	mock.Mock

	// モック用の内部状態
	shaders       map[uint32]*MockShader
	programs      map[uint32]*MockProgram
	nextShaderID  uint32
	nextProgramID uint32
}

// MockShader はモック用のシェーダー情報
type MockShader struct {
	ID           uint32
	Type         uint32
	Source       string
	Compiled     bool
	CompileError string
}

// MockProgram はモック用のプログラム情報
type MockProgram struct {
	ID        uint32
	Shaders   []uint32
	Linked    bool
	LinkError string
	Uniforms  map[string]int32
	InUse     bool
}

// NewMockOpenGLBackend は新しいMockOpenGLBackendを作成する
func NewMockOpenGLBackend() *MockOpenGLBackend {
	return &MockOpenGLBackend{
		shaders:       make(map[uint32]*MockShader),
		programs:      make(map[uint32]*MockProgram),
		nextShaderID:  1,
		nextProgramID: 1,
	}
}

// CreateShader は新しいシェーダーオブジェクトを作成する
func (m *MockOpenGLBackend) CreateShader(shaderType uint32) uint32 {
	args := m.Called(shaderType)

	// モックの戻り値を取得
	id := args.Get(0).(uint32)

	// 内部状態にシェーダーを作成
	m.shaders[id] = &MockShader{
		ID:       id,
		Type:     shaderType,
		Compiled: false,
	}

	return id
}

// ShaderSource はシェーダーオブジェクトにソースコードを設定する
func (m *MockOpenGLBackend) ShaderSource(shader uint32, source string) {
	m.Called(shader, source)

	if s, exists := m.shaders[shader]; exists {
		s.Source = source
	}
}

// CompileShader はシェーダーをコンパイルする
func (m *MockOpenGLBackend) CompileShader(shader uint32) {
	m.Called(shader)

	if s, exists := m.shaders[shader]; exists {
		// デフォルトではコンパイル成功
		s.Compiled = true

		// ソースコードに"ERROR"が含まれている場合はエラーにする
		if strings.Contains(s.Source, "ERROR") {
			s.Compiled = false
			s.CompileError = "Mock compile error"
		}
	}
}

// GetShaderiv はシェーダーパラメータを取得する
func (m *MockOpenGLBackend) GetShaderiv(shader uint32, pname uint32) int32 {
	args := m.Called(shader, pname)

	// デフォルトの動作
	if args.Get(0) == nil {
		if s, exists := m.shaders[shader]; exists {
			switch pname {
			case 0x8B81: // GL_COMPILE_STATUS
				if s.Compiled {
					return 1
				}
				return 0
			case 0x8B84: // GL_INFO_LOG_LENGTH
				return int32(len(s.CompileError))
			}
		}
		return 0
	}

	return args.Get(0).(int32)
}

// GetShaderInfoLog はシェーダーのコンパイル情報ログを取得する
func (m *MockOpenGLBackend) GetShaderInfoLog(shader uint32) string {
	args := m.Called(shader)

	// デフォルトの動作
	if args.Get(0) == nil {
		if s, exists := m.shaders[shader]; exists {
			return s.CompileError
		}
		return ""
	}

	return args.Get(0).(string)
}

// DeleteShader はシェーダーオブジェクトを削除する
func (m *MockOpenGLBackend) DeleteShader(shader uint32) {
	m.Called(shader)
	delete(m.shaders, shader)
}

// CreateProgram は新しいプログラムオブジェクトを作成する
func (m *MockOpenGLBackend) CreateProgram() uint32 {
	args := m.Called()

	// モックの戻り値を取得
	id := args.Get(0).(uint32)

	// 内部状態にプログラムを作成
	m.programs[id] = &MockProgram{
		ID:       id,
		Shaders:  make([]uint32, 0),
		Linked:   false,
		Uniforms: make(map[string]int32),
	}

	return id
}

// AttachShader はシェーダーをプログラムにアタッチする
func (m *MockOpenGLBackend) AttachShader(program, shader uint32) {
	m.Called(program, shader)

	if p, exists := m.programs[program]; exists {
		p.Shaders = append(p.Shaders, shader)
	}
}

// DetachShader はシェーダーをプログラムからデタッチする
func (m *MockOpenGLBackend) DetachShader(program, shader uint32) {
	m.Called(program, shader)

	if p, exists := m.programs[program]; exists {
		for i, s := range p.Shaders {
			if s == shader {
				p.Shaders = append(p.Shaders[:i], p.Shaders[i+1:]...)
				break
			}
		}
	}
}

// LinkProgram はプログラムをリンクする
func (m *MockOpenGLBackend) LinkProgram(program uint32) {
	m.Called(program)

	if p, exists := m.programs[program]; exists {
		// デフォルトではリンク成功
		p.Linked = true

		// アタッチされたシェーダーが2つ未満の場合はエラー
		if len(p.Shaders) < 2 {
			p.Linked = false
			p.LinkError = "Mock link error: insufficient shaders"
		}

		// アタッチされたシェーダーにコンパイルエラーがある場合はエラー
		for _, shaderID := range p.Shaders {
			if s, exists := m.shaders[shaderID]; exists && !s.Compiled {
				p.Linked = false
				p.LinkError = "Mock link error: shader compile failed"
				break
			}
		}

		// リンク成功時はユニフォーム変数を登録
		if p.Linked {
			p.Uniforms["testUniform"] = 0
			p.Uniforms["model"] = 1
			p.Uniforms["view"] = 2
			p.Uniforms["projection"] = 3
		}
	}
}

// GetProgramiv はプログラムパラメータを取得する
func (m *MockOpenGLBackend) GetProgramiv(program uint32, pname uint32) int32 {
	args := m.Called(program, pname)

	// デフォルトの動作
	if args.Get(0) == nil {
		if p, exists := m.programs[program]; exists {
			switch pname {
			case 0x8B82: // GL_LINK_STATUS
				if p.Linked {
					return 1
				}
				return 0
			case 0x8B84: // GL_INFO_LOG_LENGTH
				return int32(len(p.LinkError))
			}
		}
		return 0
	}

	return args.Get(0).(int32)
}

// GetProgramInfoLog はプログラムのリンク情報ログを取得する
func (m *MockOpenGLBackend) GetProgramInfoLog(program uint32) string {
	args := m.Called(program)

	// デフォルトの動作
	if args.Get(0) == nil {
		if p, exists := m.programs[program]; exists {
			return p.LinkError
		}
		return ""
	}

	return args.Get(0).(string)
}

// UseProgram はプログラムを使用する
func (m *MockOpenGLBackend) UseProgram(program uint32) {
	m.Called(program)

	// 全プログラムの使用状態をリセット
	for _, p := range m.programs {
		p.InUse = false
	}

	// 指定されたプログラムを使用中にする
	if p, exists := m.programs[program]; exists {
		p.InUse = true
	}
}

// DeleteProgram はプログラムオブジェクトを削除する
func (m *MockOpenGLBackend) DeleteProgram(program uint32) {
	m.Called(program)
	delete(m.programs, program)
}

// GetUniformLocation はユニフォーム変数の位置を取得する
func (m *MockOpenGLBackend) GetUniformLocation(program uint32, name string) int32 {
	args := m.Called(program, name)

	// デフォルトの動作
	if args.Get(0) == nil {
		if p, exists := m.programs[program]; exists {
			if location, exists := p.Uniforms[name]; exists {
				return location
			}
		}
		return -1
	}

	return args.Get(0).(int32)
}

// UniformMatrix4fv は4x4行列のユニフォーム変数を設定する
func (m *MockOpenGLBackend) UniformMatrix4fv(location int32, matrix [16]float32) {
	m.Called(location, matrix)
}

// Uniform3fv は3次元ベクトルのユニフォーム変数を設定する
func (m *MockOpenGLBackend) Uniform3fv(location int32, vector [3]float32) {
	m.Called(location, vector)
}

// Uniform1f は浮動小数点数のユニフォーム変数を設定する
func (m *MockOpenGLBackend) Uniform1f(location int32, value float32) {
	m.Called(location, value)
}

// Uniform1i は整数のユニフォーム変数を設定する
func (m *MockOpenGLBackend) Uniform1i(location int32, value int32) {
	m.Called(location, value)
}

// ヘルパーメソッド：テスト用
func (m *MockOpenGLBackend) GetShader(id uint32) *MockShader {
	return m.shaders[id]
}

func (m *MockOpenGLBackend) GetProgram(id uint32) *MockProgram {
	return m.programs[id]
}

func (m *MockOpenGLBackend) SetShaderCompileError(id uint32, err string) {
	if s, exists := m.shaders[id]; exists {
		s.Compiled = false
		s.CompileError = err
	}
}
