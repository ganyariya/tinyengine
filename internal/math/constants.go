package math

import "math"

// 数学定数
const (
	// 角度変換に使用する定数
	DegreesToRadians = math.Pi / 180.0
	RadiansToDegrees = 180.0 / math.Pi
	
	// よく使用される角度（ラジアン）
	HalfPi     = math.Pi / 2.0
	TwoPi      = math.Pi * 2.0
	QuarterPi  = math.Pi / 4.0
	
	// 許容誤差値（浮動小数点計算の比較用）
	Epsilon           = 1e-10  // 一般的な許容誤差
	EpsilonNormal     = 1e-6   // 通常精度の許容誤差
	EpsilonHigh       = 1e-12  // 高精度の許容誤差
	
	// ゼロ近似判定用（ベクトルの正規化などで使用）
	ZeroThreshold     = 1e-8
	
	// スケール値の制限
	MinScale          = 1e-6   // 最小スケール値
	MaxScale          = 1e6    // 最大スケール値
)

// アニメーション関連の定数
const (
	// デフォルトのアニメーション速度
	DefaultRotationSpeed = 1.0    // ラジアン/秒
	DefaultMoveSpeed     = 50.0   // ピクセル/秒
	DefaultScaleSpeed    = 0.5    // スケール変化速度
	
	// アニメーション制限値
	MinAnimationScale    = 0.1    // 最小アニメーションスケール
	MaxAnimationScale    = 5.0    // 最大アニメーションスケール
	ScaleOscillation     = 0.3    // スケール振動幅
	
	// 円運動関連
	DefaultRadius        = 100.0  // デフォルト半径
	CircularSpeedDivisor = 100.0  // 円運動速度の除数
)

// ウィンドウ関連の定数
const (
	// デフォルトウィンドウサイズ
	DefaultWindowWidth  = 800
	DefaultWindowHeight = 600
	
	// アスペクト比
	DefaultAspectRatio = float64(DefaultWindowWidth) / float64(DefaultWindowHeight)
)

// ユーティリティ関数

// IsZero 値がゼロに近いかどうかを判定
func IsZero(value float64) bool {
	return math.Abs(value) < ZeroThreshold
}

// IsEqual 二つの値が等しいかどうかを許容誤差内で判定
func IsEqual(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}

// IsEqualWithTolerance 指定した許容誤差で二つの値が等しいかどうかを判定
func IsEqualWithTolerance(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

// ClampScale スケール値を有効な範囲に制限
func ClampScale(scale float64) float64 {
	if scale < MinScale {
		return MinScale
	}
	if scale > MaxScale {
		return MaxScale
	}
	return scale
}

// ClampAnimationScale アニメーション用スケール値を有効な範囲に制限
func ClampAnimationScale(scale float64) float64 {
	if scale < MinAnimationScale {
		return MinAnimationScale
	}
	if scale > MaxAnimationScale {
		return MaxAnimationScale
	}
	return scale
}

// NormalizeAngle 角度を0〜2πの範囲に正規化
func NormalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += TwoPi
	}
	for angle >= TwoPi {
		angle -= TwoPi
	}
	return angle
}

// DegreesToRad 度をラジアンに変換
func DegreesToRad(degrees float64) float64 {
	return degrees * DegreesToRadians
}

// RadToDegrees ラジアンを度に変換
func RadToDegrees(radians float64) float64 {
	return radians * RadiansToDegrees
}