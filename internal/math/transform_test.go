package math

import (
	stdmath "math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransform_NewTransform(t *testing.T) {
	transform := NewTransform()
	
	assert.Equal(t, Vector2{X: 0, Y: 0}, transform.Position)
	assert.Equal(t, 0.0, transform.Rotation)
	assert.Equal(t, Vector2{X: 1, Y: 1}, transform.Scale)
}

func TestTransform_NewTransformWithValues(t *testing.T) {
	position := Vector2{X: 5, Y: 3}
	rotation := stdmath.Pi / 4
	scale := Vector2{X: 2, Y: 1.5}
	
	transform := NewTransformWithValues(position, rotation, scale)
	
	assert.Equal(t, position, transform.Position)
	assert.Equal(t, rotation, transform.Rotation)
	assert.Equal(t, scale, transform.Scale)
}

func TestTransform_ToMatrix(t *testing.T) {
	transform := NewTransformWithValues(
		Vector2{X: 10, Y: 5},
		stdmath.Pi/2, // 90 degrees
		Vector2{X: 2, Y: 2},
	)
	
	matrix := transform.ToMatrix()
	
	// Test a point transformation
	point := Vector2{X: 1, Y: 0}
	result := matrix.TransformPoint(point)
	
	// Expected: scale(1,0) -> (2,0), rotate 90° -> (0,2), translate -> (10,7)
	assert.InDelta(t, 10.0, result.X, Epsilon)
	assert.InDelta(t, 7.0, result.Y, Epsilon)
}

func TestTransform_TransformPoint(t *testing.T) {
	transform := NewTransformWithValues(
		Vector2{X: 5, Y: 3},
		0, // no rotation
		Vector2{X: 2, Y: 2},
	)
	
	point := Vector2{X: 1, Y: 1}
	result := transform.TransformPoint(point)
	
	// Expected: scale(1,1) -> (2,2), translate -> (7,5)
	expected := Vector2{X: 7, Y: 5}
	assert.Equal(t, expected, result)
}

func TestTransform_TransformVector(t *testing.T) {
	transform := NewTransformWithValues(
		Vector2{X: 100, Y: 100}, // translation should not affect vectors
		stdmath.Pi/2,            // 90 degrees
		Vector2{X: 2, Y: 2},
	)
	
	vector := Vector2{X: 1, Y: 0}
	result := transform.TransformVector(vector)
	
	// Expected: scale(1,0) -> (2,0), rotate 90° -> (0,2), no translation for vectors
	assert.InDelta(t, 0.0, result.X, Epsilon)
	assert.InDelta(t, 2.0, result.Y, Epsilon)
}

func TestTransform_InverseTransformPoint(t *testing.T) {
	transform := NewTransformWithValues(
		Vector2{X: 5, Y: 3},
		stdmath.Pi/4, // 45 degrees
		Vector2{X: 2, Y: 2},
	)
	
	// Transform a point and then inverse transform it
	originalPoint := Vector2{X: 1, Y: 1}
	transformedPoint := transform.TransformPoint(originalPoint)
	inversePoint, err := transform.InverseTransformPoint(transformedPoint)
	
	assert.NoError(t, err)
	assert.InDelta(t, originalPoint.X, inversePoint.X, Epsilon)
	assert.InDelta(t, originalPoint.Y, inversePoint.Y, Epsilon)
}

func TestTransform_SetPosition(t *testing.T) {
	transform := NewTransform()
	newPosition := Vector2{X: 10, Y: 20}
	
	transform.SetPosition(newPosition)
	
	assert.Equal(t, newPosition, transform.Position)
}

func TestTransform_SetRotation(t *testing.T) {
	transform := NewTransform()
	rotation := stdmath.Pi / 3
	
	transform.SetRotation(rotation)
	
	assert.Equal(t, rotation, transform.Rotation)
}

func TestTransform_SetRotationDegrees(t *testing.T) {
	transform := NewTransform()
	degrees := 90.0
	
	transform.SetRotationDegrees(degrees)
	
	assert.InDelta(t, stdmath.Pi/2, transform.Rotation, Epsilon)
}

func TestTransform_SetScale(t *testing.T) {
	transform := NewTransform()
	newScale := Vector2{X: 2, Y: 3}
	
	transform.SetScale(newScale)
	
	assert.Equal(t, newScale, transform.Scale)
}

func TestTransform_SetUniformScale(t *testing.T) {
	transform := NewTransform()
	scale := 2.5
	
	transform.SetUniformScale(scale)
	
	expected := Vector2{X: 2.5, Y: 2.5}
	assert.Equal(t, expected, transform.Scale)
}

func TestTransform_Translate(t *testing.T) {
	transform := NewTransformWithValues(Vector2{X: 5, Y: 3}, 0, Vector2{X: 1, Y: 1})
	offset := Vector2{X: 2, Y: 4}
	
	transform.Translate(offset)
	
	expected := Vector2{X: 7, Y: 7}
	assert.Equal(t, expected, transform.Position)
}

func TestTransform_Rotate(t *testing.T) {
	transform := NewTransformWithValues(Vector2{X: 0, Y: 0}, stdmath.Pi/4, Vector2{X: 1, Y: 1})
	additionalRotation := stdmath.Pi / 4
	
	transform.Rotate(additionalRotation)
	
	assert.InDelta(t, stdmath.Pi/2, transform.Rotation, Epsilon)
}

func TestTransform_RotateDegrees(t *testing.T) {
	transform := NewTransform()
	
	transform.RotateDegrees(90)
	
	assert.InDelta(t, stdmath.Pi/2, transform.Rotation, Epsilon)
}

func TestTransform_ScaleBy(t *testing.T) {
	transform := NewTransformWithValues(Vector2{X: 0, Y: 0}, 0, Vector2{X: 2, Y: 3})
	scaleBy := Vector2{X: 1.5, Y: 2}
	
	transform.ScaleBy(scaleBy)
	
	expected := Vector2{X: 3, Y: 6}
	assert.Equal(t, expected, transform.Scale)
}

func TestTransform_ScaleByUniform(t *testing.T) {
	transform := NewTransformWithValues(Vector2{X: 0, Y: 0}, 0, Vector2{X: 2, Y: 3})
	
	transform.ScaleByUniform(2)
	
	expected := Vector2{X: 4, Y: 6}
	assert.Equal(t, expected, transform.Scale)
}

func TestTransform_GetRotationDegrees(t *testing.T) {
	transform := NewTransformWithValues(Vector2{X: 0, Y: 0}, stdmath.Pi/2, Vector2{X: 1, Y: 1})
	
	degrees := transform.GetRotationDegrees()
	
	assert.InDelta(t, 90.0, degrees, Epsilon)
}

func TestTransform_Forward(t *testing.T) {
	transform := NewTransformWithValues(Vector2{X: 0, Y: 0}, 0, Vector2{X: 1, Y: 1})
	
	forward := transform.Forward()
	expected := Vector2{X: 1, Y: 0}
	
	assert.InDelta(t, expected.X, forward.X, Epsilon)
	assert.InDelta(t, expected.Y, forward.Y, Epsilon)
}

func TestTransform_Right(t *testing.T) {
	transform := NewTransformWithValues(Vector2{X: 0, Y: 0}, 0, Vector2{X: 1, Y: 1})
	
	right := transform.Right()
	expected := Vector2{X: 0, Y: 1}
	
	assert.InDelta(t, expected.X, right.X, Epsilon)
	assert.InDelta(t, expected.Y, right.Y, Epsilon)
}

func TestTransform_Combine(t *testing.T) {
	parent := NewTransformWithValues(
		Vector2{X: 5, Y: 3},
		stdmath.Pi/4,
		Vector2{X: 2, Y: 2},
	)
	
	child := NewTransformWithValues(
		Vector2{X: 1, Y: 0},
		stdmath.Pi/4,
		Vector2{X: 0.5, Y: 0.5},
	)
	
	combined := parent.Combine(child)
	
	// Check that rotations are added
	assert.InDelta(t, stdmath.Pi/2, combined.Rotation, Epsilon)
	
	// Check that scales are multiplied
	assert.InDelta(t, 1.0, combined.Scale.X, Epsilon)
	assert.InDelta(t, 1.0, combined.Scale.Y, Epsilon)
}

func TestTransform_Equals(t *testing.T) {
	t1 := NewTransformWithValues(
		Vector2{X: 1, Y: 2},
		stdmath.Pi/4,
		Vector2{X: 2, Y: 3},
	)
	
	t2 := NewTransformWithValues(
		Vector2{X: 1, Y: 2},
		stdmath.Pi/4,
		Vector2{X: 2, Y: 3},
	)
	
	t3 := NewTransformWithValues(
		Vector2{X: 1.1, Y: 2},
		stdmath.Pi/4,
		Vector2{X: 2, Y: 3},
	)
	
	assert.True(t, t1.Equals(t2))
	assert.False(t, t1.Equals(t3))
}

func TestTransform_Reset(t *testing.T) {
	transform := NewTransformWithValues(
		Vector2{X: 10, Y: 20},
		stdmath.Pi,
		Vector2{X: 5, Y: 5},
	)
	
	transform.Reset()
	
	expected := NewTransform()
	assert.True(t, transform.Equals(expected))
}