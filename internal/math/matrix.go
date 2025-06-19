package math

import (
	"errors"
	"math"
)

// Matrix3x3 represents a 3x3 matrix for 2D transformations
type Matrix3x3 [3][3]float64

// NewIdentityMatrix3x3 creates a new identity matrix
func NewIdentityMatrix3x3() Matrix3x3 {
	return Matrix3x3{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
}

// NewTranslationMatrix3x3 creates a translation matrix
func NewTranslationMatrix3x3(dx, dy float64) Matrix3x3 {
	return Matrix3x3{
		{1, 0, dx},
		{0, 1, dy},
		{0, 0, 1},
	}
}

// NewScaleMatrix3x3 creates a scale matrix
func NewScaleMatrix3x3(sx, sy float64) Matrix3x3 {
	return Matrix3x3{
		{sx, 0, 0},
		{0, sy, 0},
		{0, 0, 1},
	}
}

// NewRotationMatrix3x3 creates a rotation matrix (angle in radians)
func NewRotationMatrix3x3(angle float64) Matrix3x3 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	
	return Matrix3x3{
		{cos, -sin, 0},
		{sin, cos, 0},
		{0, 0, 1},
	}
}

// Multiply multiplies this matrix with another matrix
func (m Matrix3x3) Multiply(other Matrix3x3) Matrix3x3 {
	var result Matrix3x3
	
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				result[i][j] += m[i][k] * other[k][j]
			}
		}
	}
	
	return result
}

// MultiplyVector multiplies the matrix with a 3D vector
func (m Matrix3x3) MultiplyVector(v Vector3) Vector3 {
	return Vector3{
		X: m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z,
		Y: m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z,
		Z: m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z,
	}
}

// TransformPoint transforms a 2D point (treats as homogeneous coordinate with Z=1)
func (m Matrix3x3) TransformPoint(point Vector2) Vector2 {
	v3 := point.ToVector3()
	result := m.MultiplyVector(v3)
	return result.ToVector2()
}

// TransformVector transforms a 2D vector (treats as direction with Z=0)
func (m Matrix3x3) TransformVector(vector Vector2) Vector2 {
	v3 := Vector3{X: vector.X, Y: vector.Y, Z: 0}
	result := m.MultiplyVector(v3)
	return Vector2{X: result.X, Y: result.Y}
}

// Determinant calculates the determinant of the matrix
func (m Matrix3x3) Determinant() float64 {
	return m[0][0]*(m[1][1]*m[2][2]-m[1][2]*m[2][1]) -
		m[0][1]*(m[1][0]*m[2][2]-m[1][2]*m[2][0]) +
		m[0][2]*(m[1][0]*m[2][1]-m[1][1]*m[2][0])
}

// Inverse calculates the inverse matrix
func (m Matrix3x3) Inverse() (Matrix3x3, error) {
	det := m.Determinant()
	if math.Abs(det) < Epsilon {
		return Matrix3x3{}, errors.New("cannot invert singular matrix")
	}
	
	invDet := 1.0 / det
	
	return Matrix3x3{
		{
			(m[1][1]*m[2][2] - m[1][2]*m[2][1]) * invDet,
			(m[0][2]*m[2][1] - m[0][1]*m[2][2]) * invDet,
			(m[0][1]*m[1][2] - m[0][2]*m[1][1]) * invDet,
		},
		{
			(m[1][2]*m[2][0] - m[1][0]*m[2][2]) * invDet,
			(m[0][0]*m[2][2] - m[0][2]*m[2][0]) * invDet,
			(m[0][2]*m[1][0] - m[0][0]*m[1][2]) * invDet,
		},
		{
			(m[1][0]*m[2][1] - m[1][1]*m[2][0]) * invDet,
			(m[0][1]*m[2][0] - m[0][0]*m[2][1]) * invDet,
			(m[0][0]*m[1][1] - m[0][1]*m[1][0]) * invDet,
		},
	}, nil
}

// Transpose returns the transpose of the matrix
func (m Matrix3x3) Transpose() Matrix3x3 {
	return Matrix3x3{
		{m[0][0], m[1][0], m[2][0]},
		{m[0][1], m[1][1], m[2][1]},
		{m[0][2], m[1][2], m[2][2]},
	}
}

// IsIdentity checks if the matrix is an identity matrix
func (m Matrix3x3) IsIdentity() bool {
	identity := NewIdentityMatrix3x3()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(m[i][j]-identity[i][j]) > Epsilon {
				return false
			}
		}
	}
	return true
}

// Equals checks if two matrices are equal (within tolerance)
func (m Matrix3x3) Equals(other Matrix3x3) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(m[i][j]-other[i][j]) > Epsilon {
				return false
			}
		}
	}
	return true
}