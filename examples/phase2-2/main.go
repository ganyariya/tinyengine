package main

import (
	"log"
	"math"
	"runtime"
	"time"
	
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ganyariya/tinyengine/internal/renderer"
)

// 三角形の頂点データ（位置 + 色）
var triangleVertices = []float32{
	// 位置 (x, y, z)      色 (r, g, b)
	0.0, 0.5, 0.0,        1.0, 0.0, 0.0, // 上の頂点（赤）
	-0.5, -0.5, 0.0,      0.0, 1.0, 0.0, // 左下の頂点（緑）
	0.5, -0.5, 0.0,       0.0, 0.0, 1.0, // 右下の頂点（青）
}

func main() {
	log.Println("フェーズ2.2 ビジュアルサンプル: カラフルな三角形表示")
	log.Println("シェーダーシステムを使用したカラフルな三角形を表示します...")
	
	// OpenGL/GLFW初期化
	if err := initOpenGL(); err != nil {
		log.Fatalf("OpenGL初期化に失敗しました: %v", err)
	}
	defer glfw.Terminate()
	
	// ウィンドウ作成
	window, err := createWindow()
	if err != nil {
		log.Fatalf("ウィンドウ作成に失敗しました: %v", err)
	}
	defer window.Destroy()
	
	// OpenGL設定
	gl.Viewport(0, 0, 800, 600)
	
	// シェーダー作成
	shader, err := createColoredTriangleShader()
	if err != nil {
		log.Fatalf("シェーダー作成に失敗しました: %v", err)
	}
	defer shader.Delete()
	
	// VAO, VBO作成
	vao, vbo := createTriangleGeometry()
	defer deleteGeometry(vao, vbo)
	
	log.Println("✅ OpenGL初期化とシェーダー作成が完了しました")
	log.Println("📱 カラフルな三角形が表示されることを確認してください")
	log.Println("🎨 上が赤、左下が緑、右下が青のグラデーション三角形")
	
	// メインループ（5秒間表示）
	startTime := time.Now()
	for !window.ShouldClose() && time.Since(startTime) < 5*time.Second {
		// 画面クリア
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		
		// シェーダー使用
		shader.Use()
		
		// 時間ベースのアルファ値設定
		elapsed := float32(time.Since(startTime).Seconds())
		alpha := 0.5 + 0.5*float32(math.Sin(float64(elapsed*2.0)))
		alphaLocation := shader.GetUniformLocation("alpha")
		shader.SetUniformFloat(alphaLocation, alpha)
		
		// 時間値設定
		timeLocation := shader.GetUniformLocation("time")
		shader.SetUniformFloat(timeLocation, elapsed)
		
		// 三角形描画
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.BindVertexArray(0)
		
		// バッファスワップとイベント処理
		window.SwapBuffers()
		glfw.PollEvents()
	}
	
	log.Println("✅ フェーズ2.2のビジュアルサンプル完了")
	log.Println("")
	log.Println("確認項目:")
	log.Println("- [  ] カラフルな三角形が表示された")
	log.Println("- [  ] 上の頂点が赤色で表示された")
	log.Println("- [  ] 左下の頂点が緑色で表示された")
	log.Println("- [  ] 右下の頂点が青色で表示された")
	log.Println("- [  ] 三角形にグラデーション効果があった")
	log.Println("- [  ] 時間と共に明度が変化していた")
}

func initOpenGL() error {
	runtime.LockOSThread()
	
	if err := glfw.Init(); err != nil {
		return err
	}
	
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	
	return nil
}

func createWindow() (*glfw.Window, error) {
	window, err := glfw.CreateWindow(800, 600, "TinyEngine Phase 2.2 - Colored Triangle", nil, nil)
	if err != nil {
		return nil, err
	}
	
	window.MakeContextCurrent()
	
	if err := gl.Init(); err != nil {
		return nil, err
	}
	
	return window, nil
}

func createColoredTriangleShader() (*renderer.Shader, error) {
	vertexShaderSource := `#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 vertexColor;

void main() {
    gl_Position = vec4(aPos, 1.0);
    vertexColor = aColor;
}
`

	fragmentShaderSource := `#version 410 core
in vec3 vertexColor;
out vec4 FragColor;

uniform float alpha;
uniform float time;

void main() {
    // 基本的な色の出力
    vec3 color = vertexColor;
    
    // 時間ベースの明度変化
    float brightness = 0.5 + 0.5 * sin(time * 2.0);
    color *= brightness;
    
    FragColor = vec4(color, alpha);
}
`

	shader := renderer.NewShader()
	
	if err := shader.LoadVertexShader(vertexShaderSource); err != nil {
		return nil, err
	}
	
	if err := shader.LoadFragmentShader(fragmentShaderSource); err != nil {
		return nil, err
	}
	
	if err := shader.LinkProgram(); err != nil {
		return nil, err
	}
	
	return shader, nil
}

func createTriangleGeometry() (uint32, uint32) {
	var vao, vbo uint32
	
	// VAO作成
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	
	// VBO作成
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangleVertices)*4, gl.Ptr(triangleVertices), gl.STATIC_DRAW)
	
	// 位置属性（location = 0）
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	
	// 色属性（location = 1）
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
	
	gl.BindVertexArray(0)
	
	return vao, vbo
}

func deleteGeometry(vao, vbo uint32) {
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
}