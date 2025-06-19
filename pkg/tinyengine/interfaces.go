package tinyengine

// GameObject は全てのゲームオブジェクトが実装する基本インターフェース
type GameObject interface {
	// Initialize はオブジェクトの初期化を行う
	Initialize() error
	
	// Update はフレーム毎の更新処理を行う
	// deltaTime は前フレームからの経過時間（秒）
	Update(deltaTime float64)
	
	// Render は描画処理を行う
	Render(renderer Renderer)
	
	// Destroy はオブジェクトの破棄処理を行う
	Destroy()
}

// Renderer は描画機能を提供するインターフェース
type Renderer interface {
	// Clear は画面をクリアする
	Clear()
	
	// Present は描画内容を画面に表示する
	Present()
	
	// DrawRectangle は矩形を描画する
	DrawRectangle(x, y, width, height float32)
	
	// DrawPrimitive はプリミティブを描画する
	DrawPrimitive(primitive interface{})
	
	// DrawRectangleColor は色付き矩形を描画する
	DrawRectangleColor(x, y, width, height float32, red, green, blue, alpha float32)
	
	// DrawCircle は円を描画する
	DrawCircle(x, y, radius float32, red, green, blue, alpha float32)
	
	// DrawLine は線を描画する
	DrawLine(x1, y1, x2, y2 float32, red, green, blue, alpha float32)
}

// InputManager は入力管理機能を提供するインターフェース
type InputManager interface {
	// Update は入力状態を更新する
	Update()
	
	// IsKeyPressed は指定されたキーが押されているかを確認する
	IsKeyPressed(key int) bool
	
	// GetMousePosition はマウス座標を取得する
	GetMousePosition() (float64, float64)
	
	// IsMouseButtonPressed はマウスボタンが押されているかを確認する
	IsMouseButtonPressed(button int) bool
}

// AudioManager はオーディオ機能を提供するインターフェース
type AudioManager interface {
	// Initialize はオーディオシステムを初期化する
	Initialize() error
	
	// PlaySound はサウンドを再生する
	PlaySound(filename string) error
	
	// PlayMusic は音楽を再生する
	PlayMusic(filename string) error
	
	// StopMusic は音楽を停止する
	StopMusic()
	
	// SetVolume は音量を設定する（0.0〜1.0）
	SetVolume(volume float32)
	
	// Destroy はオーディオシステムを破棄する
	Destroy()
}