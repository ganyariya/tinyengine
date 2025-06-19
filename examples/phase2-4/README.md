# Phase 2-4: Transform Demo - 回転・拡大・移動する矩形デモ

このデモは TinyEngine の座標変換システム（行列変換、Transform クラス、ビューポート管理）の実装を視覚的に確認するためのサンプルアプリケーションです。

## 実装機能

### 数学ライブラリ
- `Vector2` / `Vector3`: 2D/3D ベクトル演算
- `Matrix3x3`: 3x3 行列演算と変換
- `Transform`: 位置、回転、スケールを統合した変換システム
- `Camera2D`: 2Dカメラとビューポート管理

### 座標変換機能
- 行列変換（平行移動、回転、スケール）
- 複合変換の合成
- 3x3行列から4x4行列への変換（OpenGL対応）
- ピクセル座標系とワールド座標系の変換

## デモの内容

3つの矩形が異なるアニメーションパターンで動作します：

1. **赤い矩形**: 高速回転、中程度のスケール振動、円形移動
2. **緑の矩形**: 中程度の回転、高速スケール振動、異なる速度の円形移動
3. **青い矩形**: 低速逆回転、低速スケール振動、逆方向の円形移動

## 実行方法

```bash
cd examples/phase2-4
go run main.go
```

## 操作

- **ESC**: プログラム終了
- 矩形が自動的に回転、スケール、移動のアニメーションを実行

## 学習ポイント

### 1. 行列変換の理解
```go
// Transform は SRT（Scale->Rotate->Translate）順で合成される
transform := mathlib.NewTransformWithValues(position, rotation, scale)
matrix := transform.ToMatrix()
```

### 2. 座標系変換
```go
// 3x3 行列を 4x4 行列に変換してOpenGLで使用
transform4x4 := convert3x3To4x4(transformMatrix)
```

### 3. アニメーション管理
```go
// フレーム時間ベースのアニメーション
func (tr *TransformableRectangle) Update(deltaTime float64) {
    tr.transform.Rotate(tr.rotationSpeed * deltaTime)
    // ...
}
```

### 4. プリミティブ変換
```go
// 頂点データに変換行列を適用
transformedVertices := applyTransformToVertices(vertices, transform4x4)
```

## 技術詳細

### 変換順序
TinyEngine では **SRT順** (Scale -> Rotate -> Translate) で変換を適用します：

1. **Scale**: 基準スケールで図形を拡大/縮小
2. **Rotate**: 原点中心に回転
3. **Translate**: 最終位置に移動

### 座標系
- **ワールド座標**: ゲーム内の論理座標系
- **ピクセル座標**: 画面上のピクセル位置（左上原点）
- **NDC座標**: OpenGLの正規化デバイス座標（-1〜1の範囲）

### パフォーマンス
- バッファプールによるVBO/VAO再利用
- フレーム時間ベースのスムーズなアニメーション
- リアルタイムFPS表示による性能監視

## 期待される結果

ウィンドウ内で3つの色付き矩形が：
- **滑らかに回転**
- **リズミカルにスケール変更**
- **円形パターンで移動**

する様子が表示されます。各矩形は異なる速度とパターンで動作し、数学ライブラリの正確性と柔軟性を実証します。

## 次のステップ

このデモの成功により、以下の機能が正常に実装されていることが確認できます：

- ✅ 行列変換システム
- ✅ 座標系変換機能
- ✅ ビューポート管理
- ✅ Transform クラスの統合
- ✅ OpenGL との統合

次のフェーズでは、この変換システムを基盤として：
- 入力システム（ユーザーインタラクション）
- シーンシステム（複数オブジェクト管理）
- アセット管理（テクスチャとスプライト）

の実装に進むことができます。