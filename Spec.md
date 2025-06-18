# TinyEngine プロジェクト仕様書

## プロジェクト概要

**TinyEngine** は Go 言語で作成する教育的な小さなゲームエンジンです。一般的なゲームエンジンの構成に沿いながら、理解しやすく実装しやすい、クリーンで美しい設計を目指します。

### 目標
- 教育目的に最適化された小さなゲームエンジン
- クリーンアーキテクチャとSOLID原則に基づいた設計
- TDD（テスト駆動開発）による高品質なコード
- 段階的な実装による学習効果の最大化

## 技術スタック

- **言語**: Go 1.19+
- **描画**: OpenGL 4 Core Profile + GLFW
- **数学ライブラリ**: mathgl
- **テストフレームワーク**: Go標準testing + testify
- **オーディオ**: Beep（後期段階で追加）

## アーキテクチャ設計

### 一般的なゲームエンジン構造

```
tinyengine/
├── cmd/
│   └── tinyengine/
│       └── main.go                 # エントリーポイント
├── internal/
│   ├── core/                       # エンジンコア
│   │   ├── engine.go
│   │   ├── application.go
│   │   └── game_loop.go
│   ├── scene/                      # シーンシステム
│   │   ├── scene.go
│   │   ├── scene_manager.go
│   │   └── actor.go
│   ├── renderer/                   # 描画システム
│   │   ├── renderer.go
│   │   ├── opengl_renderer.go
│   │   ├── shader.go
│   │   ├── texture.go
│   │   └── primitive.go
│   ├── input/                      # 入力システム
│   │   ├── input_manager.go
│   │   ├── keyboard.go
│   │   └── mouse.go
│   ├── audio/                      # オーディオシステム
│   │   ├── audio_manager.go
│   │   ├── sound.go
│   │   └── music.go
│   ├── ui/                         # UIシステム
│   │   ├── ui_manager.go
│   │   ├── button.go
│   │   ├── label.go
│   │   └── panel.go
│   ├── camera/                     # カメラシステム
│   │   ├── camera.go
│   │   └── viewport.go
│   ├── collision/                  # 衝突判定システム
│   │   ├── collision_manager.go
│   │   ├── aabb.go
│   │   └── circle.go
│   ├── math/                       # 数学ライブラリ
│   │   ├── vector.go
│   │   ├── matrix.go
│   │   ├── transform.go
│   │   └── color.go
│   └── platform/                   # プラットフォーム層
│       ├── window.go
│       └── timer.go
├── pkg/                           # 公開パッケージ
│   └── tinyengine/
│       └── interfaces.go          # 公開インターフェース
├── test/                          # テストデータ・ヘルパー
│   ├── fixtures/
│   └── helpers/
├── assets/                        # ゲームアセット
│   ├── shaders/
│   ├── textures/
│   ├── sounds/
│   └── fonts/
├── examples/                      # サンプルゲーム
│   ├── basic/
│   ├── platformer/
│   └── pong/
├── books/                         # 教育ガイドブック
│   ├── 01-getting-started.md
│   ├── 02-engine-architecture.md
│   ├── 03-rendering-system.md
│   ├── 04-input-system.md
│   ├── 05-audio-system.md
│   ├── 06-ui-system.md
│   ├── 07-camera-system.md
│   ├── 08-collision-system.md
│   └── 09-complete-game.md
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## コア機能仕様

### 1. 基本ループインターフェース

すべてのゲームオブジェクトは以下のインターフェースに従います：

```go
type GameObject interface {
    Initialize() error
    Update(deltaTime float64)
    Render(renderer Renderer)
    Destroy()
}
```

### 2. エンジンコア機能

#### 2.1 エンジン (Engine)
- ゲームループの管理
- システム初期化・終了処理
- デルタタイム計算
- フレームレート管理

#### 2.2 シーンマネージャー (SceneManager)
- シーンの切り替え
- シーン間のデータ受け渡し
- シーンスタックの管理

#### 2.3 アクターシステム (Actor System)
- 基本ゲームオブジェクト
- コンポーネントベースアーキテクチャ
- 階層構造サポート

### 3. 描画システム

#### 3.1 レンダラー (Renderer)
- OpenGL抽象化レイヤー
- バッチレンダリング
- 描画コマンドキュー

#### 3.2 図形描画
- プリミティブ描画（矩形、円、線）
- カスタムシェイプ
- ワイヤーフレーム・塗りつぶし

#### 3.3 画像描画
- テクスチャ読み込み・管理
- スプライト描画
- アニメーション

#### 3.4 テキスト描画
- フォント読み込み・管理
- テキスト描画・配置
- マルチライン対応

### 4. 入力システム

#### 4.1 キーボード入力
- キー状態管理
- キーイベント処理
- キーマッピング

#### 4.2 マウス入力
- マウス座標取得
- クリック・ドラッグ検出
- ホイール入力

### 5. オーディオシステム

#### 5.1 サウンド再生
- WAV/OGG対応
- 音量・ピッチ制御
- 3Dサウンド（基本実装）

#### 5.2 音楽再生
- BGM管理
- ループ再生
- フェードイン・アウト

### 6. UIシステム

#### 6.1 基本UI要素
- ボタン
- ラベル
- パネル

#### 6.2 レイアウト
- 絶対座標
- 相対座標
- アンカー

### 7. カメラシステム

#### 7.1 2Dカメラ
- 座標変換
- ズーム機能
- カメラ追従

#### 7.2 ビューポート
- 画面分割
- マルチカメラ

### 8. 衝突判定システム

#### 8.1 基本衝突判定
- AABB（軸平行境界ボックス）
- 円形衝突
- 点と図形

#### 8.2 衝突応答
- 衝突イベント
- 物理応答（基本）

## SOLID原則の適用

### Single Responsibility Principle (SRP)
- 各クラス・モジュールは単一の責任を持つ
- 例: `Renderer`は描画のみ、`InputManager`は入力のみ

### Open/Closed Principle (OCP)
- インターフェースによる拡張性
- 新しいレンダリングバックエンドの追加が容易

### Liskov Substitution Principle (LSP)
- インターフェースの適切な実装
- `GameObject`の実装はすべて互換性がある

### Interface Segregation Principle (ISP)
- 小さく特化したインターフェース
- クライアントは不要なメソッドに依存しない

### Dependency Inversion Principle (DIP)
- 抽象に依存し、実装に依存しない
- 依存性注入による疎結合

## テスト戦略

### テスト駆動開発（TDD）フロー

1. **Red**: 失敗するテストを書く
2. **Green**: テストを通すための最小限の実装
3. **Refactor**: コードをクリーンに改善

### テスト分類

#### 1. 単体テスト (Unit Tests)
```go
func TestVector2_Add(t *testing.T) {
    v1 := Vector2{1, 2}
    v2 := Vector2{3, 4}
    result := v1.Add(v2)
    assert.Equal(t, Vector2{4, 6}, result)
}
```

#### 2. 統合テスト (Integration Tests)
```go
func TestRenderer_DrawPrimitive(t *testing.T) {
    renderer := NewMockRenderer()
    primitive := NewRectangle(10, 10)
    
    renderer.DrawPrimitive(primitive)
    
    assert.True(t, renderer.WasDrawCalled())
}
```

#### 3. エンドツーエンドテスト (E2E Tests)
```go
func TestEngine_BasicGameLoop(t *testing.T) {
    engine := NewEngine(800, 600, "Test")
    game := NewTestGame()
    
    engine.SetGame(game)
    // 数フレーム実行してテスト
}
```

### モックとテストダブル

- インターフェースベースのモック
- テスト用のフェイク実装
- スタブによる外部依存の隔離

## 実装フェーズ

**🎯 重要方針: 各フェーズで必ずビジュアルサンプルを作成し、動作を目視確認すること**

### フェーズ 1: 基盤システム (週1-2)
- プロジェクト構造セットアップ
- 基本インターフェース定義
- エンジンコア実装
- ウィンドウ管理
- **ビジュアル目標**: 空の黒いウィンドウの表示

### フェーズ 2: 描画システム (週3-4)
- OpenGL レンダラー
- 基本図形描画
- シェーダー管理
- 座標変換
- **ビジュアル目標**: カラフルな図形の描画、変換エフェクト

### フェーズ 3: 入力システム (週5)
- キーボード・マウス入力
- 入力マッピング
- イベントシステム
- **ビジュアル目標**: マウス・キーボードで操作可能なインタラクティブデモ

### フェーズ 4: アセット管理 (週6)
- 画像読み込み
- テクスチャ管理
- フォントシステム
- **ビジュアル目標**: 画像・テキストを含むリッチなコンテンツ表示

### フェーズ 5: オーディオシステム (週7)
- サウンド再生
- 音楽管理
- 音量制御
- **ビジュアル目標**: 音声付きインタラクティブデモ

### フェーズ 6: UIシステム (週8)
- 基本UI要素
- レイアウトシステム
- イベント処理
- **ビジュアル目標**: クリック可能なボタン・メニューシステム

### フェーズ 7: カメラシステム (週9)
- 2Dカメラ実装
- ビューポート管理
- カメラ制御
- **ビジュアル目標**: スクロール・ズーム可能な世界表示

### フェーズ 8: 衝突判定 (週10)
- 基本衝突判定
- 衝突応答
- 物理演算（基本）
- **ビジュアル目標**: 物理法則に従って動くオブジェクト

### フェーズ 9: 統合・最適化 (週11-12)
- パフォーマンス最適化
- バグ修正
- ドキュメント整備
- **ビジュアル目標**: 完成したゲーム（Pong, Platformer等）

## 品質保証

### コード品質指標
- テストカバレッジ: 80%以上
- サイクロマティック複雑度: 関数あたり10以下
- Linter: golangci-lint使用

### ビジュアル品質指標
- **必須**: 各フェーズでビジュアルサンプルの実装・動作確認
- **必須**: 人間による目視確認（自動テストでは検証困難）
- **推奨**: スクリーンショット・動画での動作記録
- **推奨**: 異なるOS環境での動作確認

### CI/CD パイプライン
```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.19'
      - run: make test
      - run: make lint
```

### Makefile
```makefile
.PHONY: test lint build clean

test:
	go test ./... -v -race -coverprofile=coverage.out

lint:
	golangci-lint run

build:
	go build -o bin/tinyengine cmd/tinyengine/main.go

clean:
	rm -rf bin/ coverage.out
```

## ドキュメント戦略

### コードドキュメント
- Godoc形式のコメント
- 使用例を含む
- パッケージレベルの説明

### 教育ガイドブック (/books)

`/books`配下には、TinyEngineを段階的に実装するための詳細な教育ガイドブックを配置します。各章は**理論説明 + 実装手順 + 動作原理**の3部構成で、このガイドブック通りに実装すればゲームエンジンが完成するように設計します。

#### 各章の構成と内容

**1. 01-getting-started.md**
- Go言語環境セットアップ
- OpenGL/GLFW のインストールと設定
- プロジェクト構造の作成
- 基本的なウィンドウ表示
- **実装**: 空のウィンドウを表示するコード
- **動作原理**: GLFWとOpenGLの初期化プロセス

**2. 02-engine-architecture.md**
- ゲームエンジンの基本概念
- ゲームループとは何か？
- GameObject インターフェースの設計思想
- エンジンコアの実装
- **実装**: Engine構造体とゲームループ
- **動作原理**: フレーム処理、デルタタイム、更新と描画の分離

**3. 03-rendering-system.md**
- 描画システムの基礎理論
- OpenGLの座標系と行列変換
- シェーダーの基本概念
- **実装**: Renderer, Shader, 基本図形描画
- **動作原理**: GPU描画パイプライン、バーテックスバッファ、シェーダープログラム

**4. 04-math-system.md**
- ゲーム数学の基礎
- ベクトル演算、行列変換
- **実装**: Vector2, Matrix, Transform構造体
- **動作原理**: 座標変換、回転、スケール計算

**5. 05-scene-system.md**
- シーンとアクターの概念
- オブジェクト管理システム
- **実装**: Scene, SceneManager, Actor
- **動作原理**: オブジェクトライフサイクル、階層構造

**6. 06-input-system.md**
- 入力処理の基本
- イベントドリブンvs状態ベース
- **実装**: InputManager, Keyboard, Mouse
- **動作原理**: GLFWイベントシステム、入力状態管理

**7. 07-texture-sprite-system.md**
- テクスチャとスプライトの概念
- 画像の読み込みと描画
- **実装**: Texture, Sprite, ImageLoader
- **動作原理**: テクスチャメモリ、UV座標、アルファブレンディング

**8. 08-text-rendering.md**
- フォントレンダリングの仕組み
- ビットマップフォントvsTrueTypeフォント
- **実装**: Font, TextRenderer
- **動作原理**: グリフアトラス、文字配置計算

**9. 09-audio-system.md**
- デジタルオーディオの基礎
- サウンド再生とミキシング
- **実装**: AudioManager, Sound, Music
- **動作原理**: PCMデータ、サンプリングレート、オーディオバッファ

**10. 10-ui-system.md**
- ゲームUIの設計原則
- イベント処理とレイアウト
- **実装**: UIManager, Button, Label, Panel
- **動作原理**: UI座標系、イベントバブリング、レイアウト計算

**11. 11-camera-system.md**
- カメラの概念と座標変換
- ビューポートとプロジェクション
- **実装**: Camera2D, Viewport
- **動作原理**: ビュー行列、ワールド座標とスクリーン座標

**12. 12-collision-system.md**
- 衝突判定アルゴリズム
- 空間分割の基礎
- **実装**: CollisionManager, AABB, Circle
- **動作原理**: 境界ボックス、分離軸定理、最適化手法

**13. 13-complete-game.md**
- 全システムを統合したゲーム作成
- パフォーマンス最適化
- **実装**: 完全なPongゲーム
- **動作原理**: ゲーム状態管理、最適化テクニック

#### 各章の詳細構成

各MDファイルは以下の構造で記述：

```markdown
# [章タイトル]

## 理論編：[システム名]とは何か？
- 基本概念の説明
- なぜそのシステムが必要なのか？
- 一般的なゲームエンジンでの実装例

## 設計編：どう設計するか？
- インターフェース設計
- 依存関係の整理
- SOLID原則の適用

## 実装編：段階的実装手順
### ステップ1: 基本構造
- コード例とテスト
### ステップ2: 機能拡張
- コード例とテスト
### ステップ3: 統合
- コード例とテスト

## 動作原理編：どう動いているか？
- 内部処理フローの詳細説明
- メモリ配置とパフォーマンス
- トラブルシューティング

## 発展編：さらに学ぶために
- より高度な実装方法
- 他のゲームエンジンとの比較
- 参考資料

## 章末課題
- 実装確認用の小さなタスク
- 理解度チェック問題
```

#### ガイドブックの使用方法

1. **順次実装**: 各章を順番に進めることで段階的にエンジンが完成
2. **テストファースト**: 各章でテストコードから始める
3. **動作確認**: 章末でサンプルを実行して理解を深める
4. **応用課題**: 章末課題で知識を定着

#### 完成時の学習成果

このガイドブックを完了すると以下が身につきます：

- ✅ ゲームエンジンの全体アーキテクチャ理解
- ✅ OpenGL描画の仕組みと実装
- ✅ リアルタイムシステムの設計手法
- ✅ Go言語でのオブジェクト指向設計
- ✅ テスト駆動開発の実践
- ✅ 独自ゲームエンジンの拡張能力

### サンプルゲーム
- **Basic**: 最小限のゲーム
- **Platformer**: 2Dプラットフォーマー
- **Pong**: 古典的なPongゲーム

## 拡張可能性

### プラグインアーキテクチャ
- インターフェースベースの拡張
- 外部レンダリングバックエンド対応
- カスタムコンポーネント追加

### 将来の拡張予定
- 3D描画サポート
- ネットワーク機能
- 物理エンジン統合
- スクリプト言語組み込み

## 成功指標

### 技術指標
- ✅ 全機能が仕様通り動作
- ✅ テストカバレッジ80%以上
- ✅ クリーンなアーキテクチャ

### 教育指標
- ✅ ガイドブックに従ってエンジンを再実装可能
- ✅ サンプルゲームが動作
- ✅ 基本的なゲーム開発概念の理解
- ✅ Go言語でのソフトウェア設計スキル向上

この仕様書に基づいて、段階的にTinyEngineを実装していきます。各フェーズでTDDを実践し、クリーンなアーキテクチャを維持しながら、教育的価値の高いゲームエンジンを作成します。