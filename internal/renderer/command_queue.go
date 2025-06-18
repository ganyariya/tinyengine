package renderer

import (
	"github.com/ganyariya/tinyengine/pkg/tinyengine"
)

// CommandType は描画コマンドの種類を表す
type CommandType int

const (
	// ClearCommand は画面クリアコマンド
	ClearCommand CommandType = iota
	// RectangleCommand は矩形描画コマンド
	RectangleCommand
)

// RenderCommand は描画コマンドを表す
type RenderCommand struct {
	Type   CommandType
	Params map[string]interface{}
}

// CommandQueue は描画コマンドキューを管理する
type CommandQueue struct {
	commands []RenderCommand
}

// NewCommandQueue は新しいCommandQueueを作成する
func NewCommandQueue() *CommandQueue {
	return &CommandQueue{
		commands: make([]RenderCommand, 0),
	}
}

// AddClearCommand は画面クリアコマンドを追加する
func (q *CommandQueue) AddClearCommand() {
	command := RenderCommand{
		Type:   ClearCommand,
		Params: make(map[string]interface{}),
	}
	q.commands = append(q.commands, command)
}

// AddRectangleCommand は矩形描画コマンドを追加する
func (q *CommandQueue) AddRectangleCommand(x, y, width, height float32) {
	command := RenderCommand{
		Type: RectangleCommand,
		Params: map[string]interface{}{
			"x":      x,
			"y":      y,
			"width":  width,
			"height": height,
		},
	}
	q.commands = append(q.commands, command)
}

// Execute はキューに蓄積されたコマンドを実行する
func (q *CommandQueue) Execute(renderer tinyengine.Renderer) {
	for _, command := range q.commands {
		switch command.Type {
		case ClearCommand:
			renderer.Clear()
		case RectangleCommand:
			x := command.Params["x"].(float32)
			y := command.Params["y"].(float32)
			width := command.Params["width"].(float32)
			height := command.Params["height"].(float32)
			renderer.DrawRectangle(x, y, width, height)
		}
	}
}

// Clear はキューをクリアする
func (q *CommandQueue) Clear() {
	q.commands = q.commands[:0]
}

// Size はキューに蓄積されているコマンド数を返す
func (q *CommandQueue) Size() int {
	return len(q.commands)
}

// GetCommands はキューに蓄積されているコマンドを取得する（テスト用）
func (q *CommandQueue) GetCommands() []RenderCommand {
	return q.commands
}
