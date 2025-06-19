package math

import (
	stdmath "math"
)

// Transform represents a 2D transformation with position, rotation, and scale
type Transform struct {
	Position Vector2
	Rotation float64 // in radians
	Scale    Vector2
}

// NewTransform creates a new transform with default values
func NewTransform() Transform {
	return Transform{
		Position: Vector2{X: 0, Y: 0},
		Rotation: 0,
		Scale:    Vector2{X: 1, Y: 1},
	}
}

// NewTransformWithValues creates a new transform with specified values
func NewTransformWithValues(position Vector2, rotation float64, scale Vector2) Transform {
	return Transform{
		Position: position,
		Rotation: rotation,
		Scale:    scale,
	}
}

// ToMatrix converts the transform to a transformation matrix
func (t Transform) ToMatrix() Matrix3x3 {
	// Create individual transformation matrices
	translation := NewTranslationMatrix3x3(t.Position.X, t.Position.Y)
	rotation := NewRotationMatrix3x3(t.Rotation)
	scale := NewScaleMatrix3x3(t.Scale.X, t.Scale.Y)
	
	// Combine transformations: Scale -> Rotate -> Translate (SRT order)
	return translation.Multiply(rotation).Multiply(scale)
}

// ToInverseMatrix converts the transform to an inverse transformation matrix
func (t Transform) ToInverseMatrix() (Matrix3x3, error) {
	matrix := t.ToMatrix()
	return matrix.Inverse()
}

// TransformPoint transforms a point using this transform
func (t Transform) TransformPoint(point Vector2) Vector2 {
	matrix := t.ToMatrix()
	return matrix.TransformPoint(point)
}

// TransformVector transforms a direction vector using this transform
func (t Transform) TransformVector(vector Vector2) Vector2 {
	matrix := t.ToMatrix()
	return matrix.TransformVector(vector)
}

// InverseTransformPoint transforms a point using the inverse of this transform
func (t Transform) InverseTransformPoint(point Vector2) (Vector2, error) {
	inverse, err := t.ToInverseMatrix()
	if err != nil {
		return Vector2{}, err
	}
	return inverse.TransformPoint(point), nil
}

// InverseTransformVector transforms a direction vector using the inverse of this transform
func (t Transform) InverseTransformVector(vector Vector2) (Vector2, error) {
	inverse, err := t.ToInverseMatrix()
	if err != nil {
		return Vector2{}, err
	}
	return inverse.TransformVector(vector), nil
}

// SetPosition sets the position
func (t *Transform) SetPosition(position Vector2) {
	t.Position = position
}

// SetRotation sets the rotation in radians
func (t *Transform) SetRotation(rotation float64) {
	t.Rotation = rotation
}

// SetRotationDegrees sets the rotation in degrees
func (t *Transform) SetRotationDegrees(degrees float64) {
	t.Rotation = DegreesToRad(degrees)
}

// SetScale sets the scale
func (t *Transform) SetScale(scale Vector2) {
	t.Scale = scale
}

// SetUniformScale sets uniform scale (same for X and Y)
func (t *Transform) SetUniformScale(scale float64) {
	t.Scale = Vector2{X: scale, Y: scale}
}

// Translate moves the transform by the given offset
func (t *Transform) Translate(offset Vector2) {
	t.Position = t.Position.Add(offset)
}

// Rotate rotates the transform by the given angle in radians
func (t *Transform) Rotate(angle float64) {
	t.Rotation += angle
}

// RotateDegrees rotates the transform by the given angle in degrees
func (t *Transform) RotateDegrees(degrees float64) {
	t.Rotation += DegreesToRad(degrees)
}

// ScaleBy multiplies the current scale by the given scale
func (t *Transform) ScaleBy(scale Vector2) {
	t.Scale.X *= scale.X
	t.Scale.Y *= scale.Y
}

// ScaleByUniform multiplies the current scale by the given uniform scale
func (t *Transform) ScaleByUniform(scale float64) {
	t.Scale.X *= scale
	t.Scale.Y *= scale
}

// GetRotationDegrees returns the rotation in degrees
func (t Transform) GetRotationDegrees() float64 {
	return RadToDegrees(t.Rotation)
}

// Forward returns the forward direction vector (after rotation)
func (t Transform) Forward() Vector2 {
	return Vector2{
		X: stdmath.Cos(t.Rotation),
		Y: stdmath.Sin(t.Rotation),
	}
}

// Right returns the right direction vector (after rotation)
func (t Transform) Right() Vector2 {
	return Vector2{
		X: -stdmath.Sin(t.Rotation),
		Y: stdmath.Cos(t.Rotation),
	}
}

// Up returns the up direction vector (after rotation)
func (t Transform) Up() Vector2 {
	return Vector2{
		X: stdmath.Sin(t.Rotation),
		Y: -stdmath.Cos(t.Rotation),
	}
}

// Combine combines this transform with another transform
func (t Transform) Combine(other Transform) Transform {
	// Transform the other transform's position by this transform
	transformedPosition := t.TransformPoint(other.Position)
	
	return Transform{
		Position: transformedPosition,
		Rotation: t.Rotation + other.Rotation,
		Scale:    Vector2{X: t.Scale.X * other.Scale.X, Y: t.Scale.Y * other.Scale.Y},
	}
}

// Equals checks if two transforms are equal (within tolerance)
func (t Transform) Equals(other Transform) bool {
	return t.Position.Distance(other.Position) < Epsilon &&
		stdmath.Abs(t.Rotation-other.Rotation) < Epsilon &&
		stdmath.Abs(t.Scale.X-other.Scale.X) < Epsilon &&
		stdmath.Abs(t.Scale.Y-other.Scale.Y) < Epsilon
}

// Reset resets the transform to default values
func (t *Transform) Reset() {
	t.Position = Vector2{X: 0, Y: 0}
	t.Rotation = 0
	t.Scale = Vector2{X: 1, Y: 1}
}