package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVector2_NewVector2(t *testing.T) {
	v := NewVector2(3.0, 4.0)
	assert.Equal(t, 3.0, v.X)
	assert.Equal(t, 4.0, v.Y)
}

func TestVector2_Add(t *testing.T) {
	v1 := Vector2{X: 1, Y: 2}
	v2 := Vector2{X: 3, Y: 4}
	
	result := v1.Add(v2)
	expected := Vector2{X: 4, Y: 6}
	
	assert.Equal(t, expected, result)
}

func TestVector2_Sub(t *testing.T) {
	v1 := Vector2{X: 5, Y: 7}
	v2 := Vector2{X: 2, Y: 3}
	
	result := v1.Sub(v2)
	expected := Vector2{X: 3, Y: 4}
	
	assert.Equal(t, expected, result)
}

func TestVector2_Scale(t *testing.T) {
	v := Vector2{X: 2, Y: 3}
	
	result := v.Scale(2.5)
	expected := Vector2{X: 5, Y: 7.5}
	
	assert.Equal(t, expected, result)
}

func TestVector2_Dot(t *testing.T) {
	v1 := Vector2{X: 1, Y: 2}
	v2 := Vector2{X: 3, Y: 4}
	
	result := v1.Dot(v2)
	expected := 11.0 // 1*3 + 2*4
	
	assert.Equal(t, expected, result)
}

func TestVector2_Length(t *testing.T) {
	v := Vector2{X: 3, Y: 4}
	
	result := v.Length()
	expected := 5.0 // sqrt(3^2 + 4^2)
	
	assert.Equal(t, expected, result)
}

func TestVector2_LengthSquared(t *testing.T) {
	v := Vector2{X: 3, Y: 4}
	
	result := v.LengthSquared()
	expected := 25.0 // 3^2 + 4^2
	
	assert.Equal(t, expected, result)
}

func TestVector2_Normalize(t *testing.T) {
	v := Vector2{X: 3, Y: 4}
	
	result := v.Normalize()
	expected := Vector2{X: 0.6, Y: 0.8}
	
	assert.InDelta(t, expected.X, result.X, Epsilon)
	assert.InDelta(t, expected.Y, result.Y, Epsilon)
	assert.InDelta(t, 1.0, result.Length(), Epsilon)
}

func TestVector2_Normalize_ZeroVector(t *testing.T) {
	v := Vector2{X: 0, Y: 0}
	
	result := v.Normalize()
	expected := Vector2{X: 0, Y: 0}
	
	assert.Equal(t, expected, result)
}

func TestVector2_Distance(t *testing.T) {
	v1 := Vector2{X: 1, Y: 1}
	v2 := Vector2{X: 4, Y: 5}
	
	result := v1.Distance(v2)
	expected := 5.0 // sqrt((4-1)^2 + (5-1)^2)
	
	assert.Equal(t, expected, result)
}

func TestVector2_ToVector3(t *testing.T) {
	v2 := Vector2{X: 3, Y: 4}
	
	result := v2.ToVector3()
	expected := Vector3{X: 3, Y: 4, Z: 1}
	
	assert.Equal(t, expected, result)
}

func TestVector3_NewVector3(t *testing.T) {
	v := NewVector3(1.0, 2.0, 3.0)
	assert.Equal(t, 1.0, v.X)
	assert.Equal(t, 2.0, v.Y)
	assert.Equal(t, 3.0, v.Z)
}

func TestVector3_ToVector2(t *testing.T) {
	v3 := Vector3{X: 6, Y: 8, Z: 2}
	
	result := v3.ToVector2()
	expected := Vector2{X: 3, Y: 4} // 6/2, 8/2
	
	assert.Equal(t, expected, result)
}

func TestVector3_ToVector2_ZeroZ(t *testing.T) {
	v3 := Vector3{X: 3, Y: 4, Z: 0}
	
	result := v3.ToVector2()
	expected := Vector2{X: 3, Y: 4}
	
	assert.Equal(t, expected, result)
}

func TestVector3_Add(t *testing.T) {
	v1 := Vector3{X: 1, Y: 2, Z: 3}
	v2 := Vector3{X: 4, Y: 5, Z: 6}
	
	result := v1.Add(v2)
	expected := Vector3{X: 5, Y: 7, Z: 9}
	
	assert.Equal(t, expected, result)
}

func TestVector3_Sub(t *testing.T) {
	v1 := Vector3{X: 7, Y: 9, Z: 11}
	v2 := Vector3{X: 2, Y: 3, Z: 4}
	
	result := v1.Sub(v2)
	expected := Vector3{X: 5, Y: 6, Z: 7}
	
	assert.Equal(t, expected, result)
}

func TestVector3_Scale(t *testing.T) {
	v := Vector3{X: 2, Y: 3, Z: 4}
	
	result := v.Scale(2.0)
	expected := Vector3{X: 4, Y: 6, Z: 8}
	
	assert.Equal(t, expected, result)
}

func TestVector3_Dot(t *testing.T) {
	v1 := Vector3{X: 1, Y: 2, Z: 3}
	v2 := Vector3{X: 4, Y: 5, Z: 6}
	
	result := v1.Dot(v2)
	expected := 32.0 // 1*4 + 2*5 + 3*6
	
	assert.Equal(t, expected, result)
}

func TestVector3_Cross(t *testing.T) {
	v1 := Vector3{X: 1, Y: 0, Z: 0}
	v2 := Vector3{X: 0, Y: 1, Z: 0}
	
	result := v1.Cross(v2)
	expected := Vector3{X: 0, Y: 0, Z: 1}
	
	assert.Equal(t, expected, result)
}

func TestVector3_Length(t *testing.T) {
	v := Vector3{X: 2, Y: 3, Z: 6}
	
	result := v.Length()
	expected := 7.0 // sqrt(2^2 + 3^2 + 6^2)
	
	assert.Equal(t, expected, result)
}

func TestVector3_Normalize(t *testing.T) {
	v := Vector3{X: 2, Y: 3, Z: 6}
	
	result := v.Normalize()
	length := result.Length()
	
	assert.InDelta(t, 1.0, length, Epsilon)
}

func TestVector3_Normalize_ZeroVector(t *testing.T) {
	v := Vector3{X: 0, Y: 0, Z: 0}
	
	result := v.Normalize()
	expected := Vector3{X: 0, Y: 0, Z: 0}
	
	assert.Equal(t, expected, result)
}