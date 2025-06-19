package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatrix3x3_NewIdentity(t *testing.T) {
	matrix := NewIdentityMatrix3x3()
	
	expected := Matrix3x3{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	
	assert.Equal(t, expected, matrix)
}

func TestMatrix3x3_NewTranslation(t *testing.T) {
	matrix := NewTranslationMatrix3x3(5.0, 3.0)
	
	expected := Matrix3x3{
		{1, 0, 5},
		{0, 1, 3},
		{0, 0, 1},
	}
	
	assert.Equal(t, expected, matrix)
}

func TestMatrix3x3_NewScale(t *testing.T) {
	matrix := NewScaleMatrix3x3(2.0, 1.5)
	
	expected := Matrix3x3{
		{2.0, 0, 0},
		{0, 1.5, 0},
		{0, 0, 1},
	}
	
	assert.Equal(t, expected, matrix)
}

func TestMatrix3x3_NewRotation(t *testing.T) {
	angle := math.Pi / 4 // 45度
	matrix := NewRotationMatrix3x3(angle)
	
	cos45 := math.Cos(angle)
	sin45 := math.Sin(angle)
	
	expected := Matrix3x3{
		{cos45, -sin45, 0},
		{sin45, cos45, 0},
		{0, 0, 1},
	}
	
	assert.InDelta(t, expected[0][0], matrix[0][0], Epsilon)
	assert.InDelta(t, expected[0][1], matrix[0][1], Epsilon)
	assert.InDelta(t, expected[1][0], matrix[1][0], Epsilon)
	assert.InDelta(t, expected[1][1], matrix[1][1], Epsilon)
}

func TestMatrix3x3_Multiply(t *testing.T) {
	m1 := Matrix3x3{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	
	m2 := Matrix3x3{
		{9, 8, 7},
		{6, 5, 4},
		{3, 2, 1},
	}
	
	result := m1.Multiply(m2)
	
	expected := Matrix3x3{
		{30, 24, 18},
		{84, 69, 54},
		{138, 114, 90},
	}
	
	assert.Equal(t, expected, result)
}

func TestMatrix3x3_MultiplyVector(t *testing.T) {
	matrix := Matrix3x3{
		{2, 0, 5},
		{0, 3, 2},
		{0, 0, 1},
	}
	
	vector := Vector3{1, 2, 1}
	result := matrix.MultiplyVector(vector)
	
	expected := Vector3{7, 8, 1}
	assert.Equal(t, expected, result)
}

func TestMatrix3x3_TransformPoint(t *testing.T) {
	// Translation matrix
	translation := NewTranslationMatrix3x3(10, 5)
	point := Vector2{3, 4}
	
	result := translation.TransformPoint(point)
	expected := Vector2{13, 9}
	
	assert.Equal(t, expected, result)
}

func TestMatrix3x3_TransformVector(t *testing.T) {
	// Scale matrix
	scale := NewScaleMatrix3x3(2, 3)
	vector := Vector2{4, 2}
	
	result := scale.TransformVector(vector)
	expected := Vector2{8, 6}
	
	assert.Equal(t, expected, result)
}

func TestMatrix3x3_Determinant(t *testing.T) {
	matrix := Matrix3x3{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	
	det := matrix.Determinant()
	assert.Equal(t, 0.0, det)
	
	matrix2 := Matrix3x3{
		{2, 0, 0},
		{0, 3, 0},
		{0, 0, 1},
	}
	
	det2 := matrix2.Determinant()
	assert.Equal(t, 6.0, det2)
}

func TestMatrix3x3_Inverse(t *testing.T) {
	// Invertible matrix
	matrix := Matrix3x3{
		{2, 0, 1},
		{1, 1, 0},
		{0, 1, 1},
	}
	
	inverse, err := matrix.Inverse()
	assert.NoError(t, err)
	
	// Check if matrix * inverse = identity
	identity := matrix.Multiply(inverse)
	expectedIdentity := NewIdentityMatrix3x3()
	
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			assert.InDelta(t, expectedIdentity[i][j], identity[i][j], Epsilon)
		}
	}
}

func TestMatrix3x3_Inverse_Singular(t *testing.T) {
	// Singular matrix (determinant = 0)
	matrix := Matrix3x3{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	
	_, err := matrix.Inverse()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "singular matrix")
}

func TestTransformationChain(t *testing.T) {
	// Translation -> Rotation -> Scale
	translation := NewTranslationMatrix3x3(2, 3)
	rotation := NewRotationMatrix3x3(math.Pi / 2) // 90度
	scale := NewScaleMatrix3x3(2, 2)
	
	// Combined transformation
	combined := scale.Multiply(rotation).Multiply(translation)
	
	// Apply to point
	point := Vector2{1, 0}
	result := combined.TransformPoint(point)
	
	// Expected: translate(1,0) -> (3,3), rotate 90° -> (-3,3), scale 2x -> (-6,6)
	assert.InDelta(t, -6.0, result.X, Epsilon)
	assert.InDelta(t, 6.0, result.Y, Epsilon)
}