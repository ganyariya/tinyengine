# TASK.md - TinyEngine 実装タスク管理

このファイルは TinyEngine プロジェクトの実装タスクとその進捗を管理します。

## 重要：TDD開発フロー

各タスクは以下の手順で実行してください：

1. **Red**: 失敗するテストを書く
2. **Green**: テストを通すための最小限の実装
3. **Refactor**: コードをクリーンに改善
4. **テスト実行**: `go test ./...` でテストを実行
5. **ビルド確認**: `go build cmd/tinyengine/main.go` でビルドを確認
6. **リント確認**: `golangci-lint run` でコード品質を確認
7. **コミット**: 問題なければコミット
8. **TASK.md更新**: タスク完了後にこのファイルを更新

## タスク状況

### フェーズ1: 基盤システム (1-2週目)

#### 1.1 プロジェクト初期化
- [x] go.mod ファイル作成
- [x] 基本的なディレクトリ構造作成
- [x] Makefile 作成
- [x] GitHub Actions CI設定
- [x] 基本的な main.go 作成

#### 1.2 コアインターフェース定義
- [x] GameObject インターフェース定義 (`pkg/tinyengine/interfaces.go`)
- [x] Renderer インターフェース定義
- [x] InputManager インターフェース定義
- [x] AudioManager インターフェース定義

#### 1.3 エンジンコア実装
- [x] Engine 構造体実装 (`internal/core/engine.go`)
- [x] Application 構造体実装 (`internal/core/application.go`)
- [x] GameLoop 実装 (`internal/core/game_loop.go`)
- [x] デルタタイム計算機能

#### 1.4 プラットフォーム層
- [x] Window 管理実装 (`internal/platform/window.go`)
- [x] Timer 実装 (`internal/platform/timer.go`)
- [x] GLFW 初期化処理

### フェーズ2: 描画システム (3-4週目)

#### 2.1 基本レンダラー
- [ ] Renderer インターフェース実装 (`internal/renderer/renderer.go`)
- [ ] OpenGLRenderer 実装 (`internal/renderer/opengl_renderer.go`)
- [ ] 基本的な描画コマンドキュー

#### 2.2 シェーダーシステム
- [ ] Shader 構造体実装 (`internal/renderer/shader.go`)
- [ ] 基本的な頂点・フラグメントシェーダー
- [ ] シェーダープログラム管理

#### 2.3 基本図形描画
- [ ] Primitive 基底クラス (`internal/renderer/primitive.go`)
- [ ] 矩形描画機能
- [ ] 円形描画機能
- [ ] 線描画機能

#### 2.4 座標変換システム
- [ ] 行列変換実装 (`internal/math/matrix.go`)
- [ ] 座標系変換機能
- [ ] ビューポート管理

### フェーズ3: 数学ライブラリ (5週目前半)

#### 3.1 基本数学構造体
- [ ] Vector2 実装 (`internal/math/vector.go`)
- [ ] Matrix 実装 (`internal/math/matrix.go`)
- [ ] Transform 実装 (`internal/math/transform.go`)
- [ ] Color 実装 (`internal/math/color.go`)

#### 3.2 数学演算
- [ ] ベクトル演算（加算、減算、内積、外積）
- [ ] 行列演算（乗算、逆行列、転置）
- [ ] 変換演算（平行移動、回転、スケール）

### フェーズ4: 入力システム (5週目後半)

#### 4.1 入力管理
- [ ] InputManager 実装 (`internal/input/input_manager.go`)
- [ ] 入力状態管理システム

#### 4.2 キーボード入力
- [ ] Keyboard 実装 (`internal/input/keyboard.go`)
- [ ] キー状態追跡
- [ ] キーイベント処理

#### 4.3 マウス入力
- [ ] Mouse 実装 (`internal/input/mouse.go`)
- [ ] マウス座標取得
- [ ] クリック・ドラッグ検出

### フェーズ5: シーンシステム (6週目)

#### 5.1 基本シーンシステム
- [ ] Scene インターフェース (`internal/scene/scene.go`)
- [ ] SceneManager 実装 (`internal/scene/scene_manager.go`)
- [ ] シーン切り替え機能

#### 5.2 アクターシステム
- [ ] Actor 基底クラス (`internal/scene/actor.go`)
- [ ] 階層構造サポート
- [ ] コンポーネントシステム基盤

### フェーズ6: アセット管理 (7週目)

#### 6.1 テクスチャシステム
- [ ] Texture 実装 (`internal/renderer/texture.go`)
- [ ] 画像読み込み機能
- [ ] テクスチャ管理システム

#### 6.2 フォントシステム
- [ ] Font 実装 (フォント管理システム)
- [ ] テキスト描画機能
- [ ] 文字配置計算

### フェーズ7: オーディオシステム (8週目)

#### 7.1 オーディオ管理
- [ ] AudioManager 実装 (`internal/audio/audio_manager.go`)
- [ ] オーディオ初期化・終了処理

#### 7.2 サウンド再生
- [ ] Sound 実装 (`internal/audio/sound.go`)
- [ ] Music 実装 (`internal/audio/music.go`)
- [ ] 音量・ピッチ制御

### フェーズ8: UIシステム (9週目)

#### 8.1 UI管理
- [ ] UIManager 実装 (`internal/ui/ui_manager.go`)
- [ ] UI座標系管理

#### 8.2 基本UI要素
- [ ] Button 実装 (`internal/ui/button.go`)
- [ ] Label 実装 (`internal/ui/label.go`)
- [ ] Panel 実装 (`internal/ui/panel.go`)

#### 8.3 レイアウトシステム
- [ ] 絶対座標レイアウト
- [ ] 相対座標レイアウト
- [ ] アンカーシステム

### フェーズ9: カメラシステム (10週目)

#### 9.1 2Dカメラ
- [ ] Camera2D 実装 (`internal/camera/camera.go`)
- [ ] カメラ座標変換
- [ ] ズーム機能

#### 9.2 ビューポート
- [ ] Viewport 実装 (`internal/camera/viewport.go`)
- [ ] 画面分割機能
- [ ] マルチカメラサポート

### フェーズ10: 衝突判定システム (11週目)

#### 10.1 衝突判定管理
- [ ] CollisionManager 実装 (`internal/collision/collision_manager.go`)
- [ ] 衝突判定最適化

#### 10.2 基本衝突判定
- [ ] AABB 実装 (`internal/collision/aabb.go`)
- [ ] Circle 実装 (`internal/collision/circle.go`)
- [ ] 点と図形の衝突判定

#### 10.3 衝突応答
- [ ] 衝突イベントシステム
- [ ] 基本的な物理応答

### フェーズ11: 統合・最適化 (12週目)

#### 11.1 サンプルゲーム作成
- [ ] Basic サンプル (`examples/basic/`)
- [ ] Pong ゲーム (`examples/pong/`)
- [ ] Platformer サンプル (`examples/platformer/`)

#### 11.2 パフォーマンス最適化
- [ ] 描画バッチング最適化
- [ ] メモリ使用量最適化
- [ ] プロファイリングとベンチマーク

#### 11.3 ドキュメント整備
- [ ] API ドキュメント完成
- [ ] 教育ガイドブック作成 (`books/` ディレクトリ)
- [ ] README 更新

## 進捗状況

- **開始日**: 2025-06-18
- **現在のフェーズ**: フェーズ3（数学ライブラリ）
- **完了タスク数**: フェーズ1完了 (16/16)
- **次回作業**: 数学ライブラリ実装（Vector2、Matrix、Transform、Color）

## 注意事項

1. **各タスク完了後、必ずこのファイルを更新してください**
2. **TDD を厳守し、テストファーストで開発してください**
3. **各フェーズ完了時には全体的なテストとビルドを実行してください**
4. **コミットメッセージは日本語で分かりやすく記述してください**
5. **問題が発生した場合、このファイルに記録してください**

## コミット例

```
feat: Vector2構造体とベクトル演算機能を実装

- Vector2構造体の基本実装
- 加算、減算、スカラー倍の演算機能
- 単体テストによる動作確認
- TDD手法により品質確保

closes #1
```

最終更新: [日付を記入]
更新者: Claude Code