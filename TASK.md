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

#### 1.5 教育ガイドブック作成 📚
- [ ] **基盤システム実装ガイド** (`books/01-getting-started.md`)
  - Go環境セットアップ手順
  - OpenGL/GLFW環境構築方法
  - 基本ウィンドウ表示の実装解説
- [ ] **エンジンアーキテクチャガイド** (`books/02-engine-architecture.md`)
  - ゲームエンジンの基本概念
  - ゲームループとは何か？
  - デルタタイム管理の重要性
  - TinyEngineのアーキテクチャ設計思想

### フェーズ2: 描画システム (3-4週目)

#### 2.1 基本レンダラー
- [x] Renderer インターフェース実装 (`internal/renderer/renderer.go`)
- [x] OpenGLRenderer 実装 (`internal/renderer/opengl_renderer.go`)
- [x] 基本的な描画コマンドキュー
- [x] **ビジュアルサンプル**: 黒い背景ウィンドウの表示確認 (`examples/phase2-1/`)

#### 2.2 シェーダーシステム
- [x] Shader 構造体実装 (`internal/renderer/shader.go`)
- [x] 基本的な頂点・フラグメントシェーダー
- [x] シェーダープログラム管理
- [x] **ビジュアルサンプル**: カラフルな三角形の表示 (`examples/phase2-2/`)

#### 2.3 基本図形描画
- [x] Primitive 基底クラス (`internal/renderer/primitive.go`)
- [x] 矩形描画機能
- [x] 円形描画機能
- [x] 線描画機能
- [x] **ビジュアルサンプル**: カラフルな図形ギャラリー表示 (`examples/phase2-3/`)

#### 2.4 座標変換システム
- [ ] 行列変換実装 (`internal/math/matrix.go`)
- [ ] 座標系変換機能
- [ ] ビューポート管理
- [ ] **ビジュアルサンプル**: 回転・拡大・移動する矩形デモ (`examples/phase2-4/`)
- [ ] **リファクタする**: マジック変数を消す。クリーンなコードにする。その場限りの実装をやめる。

#### 2.5 教育ガイドブック作成 📚
- [ ] **描画システム基礎ガイド** (`books/03-rendering-system.md`)
  - OpenGLとは何か？現代グラフィックスAPIの仕組み
  - GPU vs CPU：なぜGPUが描画に適しているのか？
  - 描画パイプラインの基本概念
  - 座標系と変換行列の数学的理解
  - TinyEngineの描画アーキテクチャ実装解説
- [ ] **シェーダーシステムガイド** (`books/04-shader-system.md`)
  - シェーダーとは何か？頂点シェーダーとフラグメントシェーダー
  - GLSL（OpenGL Shading Language）基本文法
  - シェーダープログラムの作成と管理
  - Uniform変数とAttribute変数の違い
  - カラフルな三角形から学ぶシェーダー実装

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
- [ ] **ビジュアルサンプル**: ベクトル演算の視覚化デモ (`examples/phase3/`)
- [ ] **リファクタする**: マジック変数を消す。クリーンなコードにする。その場限りの実装をやめる。

#### 3.3 教育ガイドブック作成 📚
- [ ] **ゲーム数学基礎ガイド** (`books/05-math-system.md`)
  - ゲーム開発に必要な数学の基礎
  - ベクトルとは何か？2Dベクトルの概念と応用
  - 行列変換の数学的意味と実装
  - 座標変換（平行移動、回転、スケール）の仕組み
  - TinyEngineの数学ライブラリ実装解説

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
- [ ] **ビジュアルサンプル**: インタラクティブな図形操作デモ (`examples/phase4/`)
- [ ] **リファクタする**: マジック変数を消す。クリーンなコードにする。その場限りの実装をやめる。

#### 4.4 教育ガイドブック作成 📚
- [ ] **入力システムガイド** (`books/06-input-system.md`)
  - ゲームの入力処理基本概念
  - イベントドリブン vs 状態ベース入力処理
  - キーボード入力の実装と状態管理
  - マウス入力と座標変換
  - TinyEngineの入力アーキテクチャ実装解説

### フェーズ5: シーンシステム (6週目)

#### 5.1 基本シーンシステム
- [ ] Scene インターフェース (`internal/scene/scene.go`)
- [ ] SceneManager 実装 (`internal/scene/scene_manager.go`)
- [ ] シーン切り替え機能

#### 5.2 アクターシステム
- [ ] Actor 基底クラス (`internal/scene/actor.go`)
- [ ] 階層構造サポート
- [ ] コンポーネントシステム基盤
- [ ] **ビジュアルサンプル**: 複数シーンの切り替えデモ (`examples/phase5/`)
- [ ] **リファクタする**: マジック変数を消す。クリーンなコードにする。その場限りの実装をやめる。

#### 5.3 教育ガイドブック作成 📚
- [ ] **シーンシステムガイド** (`books/07-scene-system.md`)
  - シーンとアクターの概念
  - ゲームオブジェクト階層管理の重要性
  - コンポーネントベースアーキテクチャ入門
  - シーン切り替えとライフサイクル管理
  - TinyEngineのシーンシステム実装解説

### フェーズ6: アセット管理 (7週目)

#### 6.1 テクスチャシステム
- [ ] Texture 実装 (`internal/renderer/texture.go`)
- [ ] 画像読み込み機能
- [ ] テクスチャ管理システム

#### 6.2 フォントシステム
- [ ] Font 実装 (フォント管理システム)
- [ ] テキスト描画機能
- [ ] 文字配置計算
- [ ] **ビジュアルサンプル**: 画像・テキスト表示ギャラリー (`examples/phase6/`)
- [ ] **リファクタする**: マジック変数を消す。クリーンなコードにする。その場限りの実装をやめる。

#### 6.3 教育ガイドブック作成 📚
- [ ] **テクスチャ・スプライトシステムガイド** (`books/08-texture-sprite-system.md`)
  - テクスチャとスプライトの概念
  - 画像ファイル読み込みの仕組み
  - UV座標とテクスチャマッピング
  - スプライトバッチングによる描画最適化
  - TinyEngine のテクスチャシステム実装解説
- [ ] **テキストレンダリングガイド** (`books/09-text-rendering.md`)
  - フォントレンダリングの基本仕組み
  - ビットマップフォント vs ベクターフォント
  - 文字配置とレイアウト計算
  - 多言語対応（日本語）の考慮点
  - TinyEngineのフォントシステム実装解説

### フェーズ7: オーディオシステム (8週目)

#### 7.1 オーディオ管理
- [ ] AudioManager 実装 (`internal/audio/audio_manager.go`)
- [ ] オーディオ初期化・終了処理

#### 7.2 サウンド再生
- [ ] Sound 実装 (`internal/audio/sound.go`)
- [ ] Music 実装 (`internal/audio/music.go`)
- [ ] 音量・ピッチ制御
- [ ] **ビジュアルサンプル**: サウンド付きインタラクティブデモ (`examples/phase7/`)
- [ ] **リファクタする**: マジック変数を消す。クリーンなコードにする。その場限りの実装をやめる。

#### 7.3 教育ガイドブック作成 📚
- [ ] **オーディオシステムガイド** (`books/10-audio-system.md`)
  - デジタルオーディオの基礎知識
  - サンプリングレートとビット深度
  - サウンドエフェクト vs 背景音楽の管理
  - オーディオミキシングとボリューム制御
  - TinyEngineのオーディオアーキテクチャ実装解説

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
- [ ] **ビジュアルサンプル**: インタラクティブUIデモ (`examples/phase8/`)
- [ ] **リファクタする**: マジック変数を消す。クリーンなコードにする。その場限りの実装をやめる。

#### 8.4 教育ガイドブック作成 📚
- [ ] **UIシステムガイド** (`books/11-ui-system.md`)
  - ゲームUI設計の基本原則
  - イベント処理とUIレスポンス
  - レイアウトシステムの概念と実装
  - アンカーとスケーリング対応
  - TinyEngineのUIアーキテクチャ実装解説

### フェーズ9: カメラシステム (10週目)

#### 9.1 2Dカメラ
- [ ] Camera2D 実装 (`internal/camera/camera.go`)
- [ ] カメラ座標変換
- [ ] ズーム機能

#### 9.2 ビューポート
- [ ] Viewport 実装 (`internal/camera/viewport.go`)
- [ ] 画面分割機能
- [ ] マルチカメラサポート
- [ ] **ビジュアルサンプル**: カメラ操作・ズームデモ (`examples/phase9/`)
- [ ] **リファクタする**: マジック変数を消す。クリーンなコードにする。その場限りの実装をやめる。

#### 9.3 教育ガイドブック作成 📚
- [ ] **カメラシステムガイド** (`books/12-camera-system.md`)
  - カメラの概念と座標変換
  - ビューポートとプロジェクション
  - 2Dカメラのズームと移動実装
  - 複数カメラとスクリーン分割
  - TinyEngineのカメラアーキテクチャ実装解説

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
- [ ] **ビジュアルサンプル**: 物理衝突デモ (`examples/phase10/`)
- [ ] **リファクタする**: マジック変数を消す。クリーンなコードにする。その場限りの実装をやめる。

#### 10.4 教育ガイドブック作成 📚
- [ ] **衝突判定システムガイド** (`books/13-collision-system.md`)
  - 衝突判定アルゴリズムの基礎
  - AABB（軸平行境界ボックス）と円形衝突判定
  - 空間分割によるパフォーマンス最適化
  - 衝突応答と物理シミュレーション基礎
  - TinyEngineの衝突判定アーキテクチャ実装解説

### フェーズ11: 統合・最適化 (12週目)

#### 11.1 サンプルゲーム作成
- [ ] Basic サンプル (`examples/basic/`)
- [ ] Pong ゲーム (`examples/pong/`)
- [ ] Platformer サンプル (`examples/platformer/`)
- [ ] **ビジュアルサンプル**: 完成ゲームデモ（操作可能）

#### 11.2 パフォーマンス最適化
- [ ] 描画バッチング最適化
- [ ] メモリ使用量最適化
- [ ] プロファイリングとベンチマーク
- [ ] **リファクタする**: マジック変数を消す。クリーンなコードにする。その場限りの実装をやめる。

#### 11.3 ドキュメント整備
- [ ] API ドキュメント完成
- [ ] README 更新

#### 11.4 教育ガイドブック作成 📚
- [ ] **完全ゲーム開発ガイド** (`books/14-complete-game.md`)
  - 全システム統合によるゲーム作成
  - パフォーマンス最適化の実践
  - デバッグとプロファイリング手法
  - ゲーム公開と配布の考慮点
  - TinyEngineで学んだ技術の応用方法
- [ ] **総合学習ガイド** (`books/15-learning-path.md`)
  - TinyEngine学習カリキュラム全体像
  - 各章の学習目標と評価方法
  - 発展学習への道筋（3D、ネットワーク、物理エンジン等）
  - ゲーム開発キャリアパス
  - おすすめ学習リソースと書籍

## 進捗状況

- **開始日**: 2025-06-18
- **現在のフェーズ**: フェーズ2（描画システム）
- **完了タスク数**: フェーズ1完了 (16/16), フェーズ2.1完了 (4/4), フェーズ2.2完了 (4/4), フェーズ2.3完了 (5/5)
- **次回作業**: 座標変換システム（行列変換実装、座標系変換機能、ビューポート管理）

## 注意事項

1. **各タスク完了後、必ずこのファイルを更新してください**
2. **TDD を厳守し、テストファーストで開発してください**
3. **各フェーズ完了時には全体的なテストとビルドを実行してください**
4. **コミットメッセージは日本語で分かりやすく記述してください**
5. **問題が発生した場合、このファイルに記録してください**
6. **📚 教育ガイドブック作成は各フェーズの重要な仕上げタスクです**
   - 実装だけでなく、その仕組みや理論的背景を説明
   - OpenGL、GPU、ゲームエンジンアーキテクチャの基本概念を含む
   - 段階的学習で初心者でも理解できる構成
   - 実装コードと理論解説を組み合わせた実践的内容

## 🎯 重要：ビジュアル確認の義務

**各フェーズのビジュアルサンプルは必須実装項目です。**

### 学習者モチベーション維持のため：
- 各フェーズで「動く結果」が見えることで達成感を提供
- 理論だけでなく実際の動作確認で理解を深める
- 段階的な成長を視覚的に実感できる

### 実装品質確認のため：
- 人間の目で実際にウィンドウに正しく表示されているかを確認
- 自動テストでは検証しきれない視覚的なバグの発見
- OpenGL描画の正常動作確認（環境依存問題の早期発見）

### 各フェーズのビジュアルサンプル実行手順：
1. `examples/phase2-1/` 等のディレクトリを作成
2. そのフェーズで実装した機能を使用するサンプルアプリケーションを実装
3. `go run` で実行して動作確認
4. README.md にスクリーンショットと実行方法を記述
5. 必ず人間が実際にウィンドウを確認すること

**⚠️ ビジュアルサンプルなしでフェーズ完了とみなしてはいけません**

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