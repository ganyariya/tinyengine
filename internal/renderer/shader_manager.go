package renderer

import (
	"fmt"
	"sort"
)

// ShaderManager は複数のシェーダープログラムを管理する
type ShaderManager struct {
	shaders       map[string]*Shader
	currentShader string
}

// NewShaderManager は新しいShaderManagerを作成する
func NewShaderManager() *ShaderManager {
	return &ShaderManager{
		shaders:       make(map[string]*Shader),
		currentShader: "",
	}
}

// LoadShader はシェーダーソースコードからシェーダーを読み込む
func (sm *ShaderManager) LoadShader(name, vertexSource, fragmentSource string) error {
	// 同名のシェーダーが既に存在する場合は先に削除
	if sm.HasShader(name) {
		sm.DeleteShader(name)
	}

	// 新しいシェーダー作成（実際のOpenGLバックエンドを使用）
	shader := NewShader(NewRealOpenGLBackend())

	// 頂点シェーダー読み込み
	if err := shader.LoadVertexShader(vertexSource); err != nil {
		return fmt.Errorf("failed to load vertex shader '%s': %v", name, err)
	}

	// フラグメントシェーダー読み込み
	if err := shader.LoadFragmentShader(fragmentSource); err != nil {
		return fmt.Errorf("failed to load fragment shader '%s': %v", name, err)
	}

	// プログラムリンク
	if err := shader.LinkProgram(); err != nil {
		return fmt.Errorf("failed to link shader program '%s': %v", name, err)
	}

	// マネージャーに登録
	sm.shaders[name] = shader

	return nil
}

// LoadShaderFromFiles はファイルからシェーダーを読み込む
func (sm *ShaderManager) LoadShaderFromFiles(name, vertexPath, fragmentPath string) error {
	shader, err := CreateShaderFromFiles(vertexPath, fragmentPath)
	if err != nil {
		return fmt.Errorf("failed to create shader '%s' from files: %v", name, err)
	}

	// 同名のシェーダーが既に存在する場合は先に削除
	if sm.HasShader(name) {
		sm.DeleteShader(name)
	}

	sm.shaders[name] = shader
	return nil
}

// LoadBuiltinShader は組み込みシェーダーを読み込む
func (sm *ShaderManager) LoadBuiltinShader(name string) error {
	shader, err := CreateBuiltinShader(name)
	if err != nil {
		return fmt.Errorf("failed to load builtin shader '%s': %v", name, err)
	}

	// 同名のシェーダーが既に存在する場合は先に削除
	if sm.HasShader(name) {
		sm.DeleteShader(name)
	}

	sm.shaders[name] = shader
	return nil
}

// GetShader は指定された名前のシェーダーを取得する
func (sm *ShaderManager) GetShader(name string) *Shader {
	if shader, exists := sm.shaders[name]; exists {
		return shader
	}
	return nil
}

// HasShader は指定された名前のシェーダーが存在するかを確認する
func (sm *ShaderManager) HasShader(name string) bool {
	_, exists := sm.shaders[name]
	return exists
}

// UseShader は指定された名前のシェーダーを使用する
func (sm *ShaderManager) UseShader(name string) bool {
	if shader, exists := sm.shaders[name]; exists {
		shader.Use()
		sm.currentShader = name
		return true
	}
	return false
}

// GetCurrentShader は現在使用中のシェーダー名を取得する
func (sm *ShaderManager) GetCurrentShader() string {
	return sm.currentShader
}

// DeleteShader は指定された名前のシェーダーを削除する
func (sm *ShaderManager) DeleteShader(name string) bool {
	if shader, exists := sm.shaders[name]; exists {
		shader.Delete()
		delete(sm.shaders, name)

		// 現在使用中のシェーダーが削除された場合はクリア
		if sm.currentShader == name {
			sm.currentShader = ""
		}

		return true
	}
	return false
}

// DeleteAllShaders はすべてのシェーダーを削除する
func (sm *ShaderManager) DeleteAllShaders() {
	for name, shader := range sm.shaders {
		shader.Delete()
		delete(sm.shaders, name)
	}
	sm.currentShader = ""
}

// GetShaderCount は登録されているシェーダー数を取得する
func (sm *ShaderManager) GetShaderCount() int {
	return len(sm.shaders)
}

// GetShaderNames は登録されているシェーダー名のリストを取得する
func (sm *ShaderManager) GetShaderNames() []string {
	names := make([]string, 0, len(sm.shaders))
	for name := range sm.shaders {
		names = append(names, name)
	}

	// アルファベット順にソート
	sort.Strings(names)

	return names
}

// SetUniformMat4 は現在のシェーダーに4x4行列ユニフォームを設定する
func (sm *ShaderManager) SetUniformMat4(name string, matrix [16]float32) bool {
	if sm.currentShader == "" {
		return false
	}

	shader := sm.GetShader(sm.currentShader)
	if shader == nil {
		return false
	}

	location := shader.GetUniformLocation(name)
	if location < 0 {
		return false
	}

	shader.SetUniformMat4(location, matrix)
	return true
}

// SetUniformVec3 は現在のシェーダーに3次元ベクトルユニフォームを設定する
func (sm *ShaderManager) SetUniformVec3(name string, vector [3]float32) bool {
	if sm.currentShader == "" {
		return false
	}

	shader := sm.GetShader(sm.currentShader)
	if shader == nil {
		return false
	}

	location := shader.GetUniformLocation(name)
	if location < 0 {
		return false
	}

	shader.SetUniformVec3(location, vector)
	return true
}

// SetUniformFloat は現在のシェーダーに浮動小数点数ユニフォームを設定する
func (sm *ShaderManager) SetUniformFloat(name string, value float32) bool {
	if sm.currentShader == "" {
		return false
	}

	shader := sm.GetShader(sm.currentShader)
	if shader == nil {
		return false
	}

	location := shader.GetUniformLocation(name)
	if location < 0 {
		return false
	}

	shader.SetUniformFloat(location, value)
	return true
}
