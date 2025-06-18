package renderer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// LoadShaderFromFile はファイルからシェーダーソースコードを読み込む
func LoadShaderFromFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read shader file %s: %v", filePath, err)
	}
	
	return string(data), nil
}

// CreateShaderFromFiles は頂点・フラグメントシェーダーファイルからShaderを作成する
func CreateShaderFromFiles(vertexPath, fragmentPath string) (*Shader, error) {
	// 頂点シェーダー読み込み
	vertexSource, err := LoadShaderFromFile(vertexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load vertex shader: %v", err)
	}
	
	// フラグメントシェーダー読み込み
	fragmentSource, err := LoadShaderFromFile(fragmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load fragment shader: %v", err)
	}
	
	// Shader作成
	shader := NewShader()
	
	// 頂点シェーダー読み込み
	if err := shader.LoadVertexShader(vertexSource); err != nil {
		return nil, fmt.Errorf("failed to load vertex shader: %v", err)
	}
	
	// フラグメントシェーダー読み込み
	if err := shader.LoadFragmentShader(fragmentSource); err != nil {
		return nil, fmt.Errorf("failed to load fragment shader: %v", err)
	}
	
	// プログラムリンク
	if err := shader.LinkProgram(); err != nil {
		return nil, fmt.Errorf("failed to link shader program: %v", err)
	}
	
	return shader, nil
}

// GetBuiltinShaderPaths は組み込みシェーダーのパスを取得する
func GetBuiltinShaderPaths(shaderName string) (vertexPath, fragmentPath string) {
	assetsDir := "assets/shaders"
	vertexPath = filepath.Join(assetsDir, shaderName+".vert")
	fragmentPath = filepath.Join(assetsDir, shaderName+".frag")
	return
}

// CreateBuiltinShader は組み込みシェーダーからShaderを作成する
func CreateBuiltinShader(shaderName string) (*Shader, error) {
	vertexPath, fragmentPath := GetBuiltinShaderPaths(shaderName)
	return CreateShaderFromFiles(vertexPath, fragmentPath)
}

// writeStringToFile はテスト用のヘルパー関数
func writeStringToFile(filePath, content string) error {
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}