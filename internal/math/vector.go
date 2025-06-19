package math

import (
	"math"
)

// Vector2 represents a 2D vector
type Vector2 struct {
	X, Y float64
}

// Vector3 represents a 3D vector (used for homogeneous coordinates in 2D)
type Vector3 struct {
	X, Y, Z float64
}

// NewVector2 creates a new 2D vector
func NewVector2(x, y float64) Vector2 {
	return Vector2{X: x, Y: y}
}

// NewVector3 creates a new 3D vector
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

// Add adds two vectors
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub subtracts another vector from this vector
func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Scale multiplies the vector by a scalar
func (v Vector2) Scale(scalar float64) Vector2 {
	return Vector2{X: v.X * scalar, Y: v.Y * scalar}
}

// Dot calculates the dot product
func (v Vector2) Dot(other Vector2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Length calculates the length of the vector
func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// LengthSquared calculates the squared length (faster than Length)
func (v Vector2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Normalize returns a normalized vector (length = 1)
func (v Vector2) Normalize() Vector2 {
	length := v.Length()
	if IsZero(length) {
		return Vector2{X: 0, Y: 0}
	}
	return Vector2{X: v.X / length, Y: v.Y / length}
}

// Distance calculates the distance between two vectors
func (v Vector2) Distance(other Vector2) float64 {
	return v.Sub(other).Length()
}

// ToVector3 converts Vector2 to Vector3 with Z=1 (for homogeneous coordinates)
func (v Vector2) ToVector3() Vector3 {
	return Vector3{X: v.X, Y: v.Y, Z: 1.0}
}

// ToVector2 converts Vector3 to Vector2 by dividing by Z (perspective divide)
func (v Vector3) ToVector2() Vector2 {
	if v.Z == 0 {
		return Vector2{X: v.X, Y: v.Y}
	}
	return Vector2{X: v.X / v.Z, Y: v.Y / v.Z}
}

// Add adds two 3D vectors
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

// Sub subtracts another 3D vector from this vector
func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

// Scale multiplies the 3D vector by a scalar
func (v Vector3) Scale(scalar float64) Vector3 {
	return Vector3{X: v.X * scalar, Y: v.Y * scalar, Z: v.Z * scalar}
}

// Dot calculates the dot product of two 3D vectors
func (v Vector3) Dot(other Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross calculates the cross product of two 3D vectors
func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Length calculates the length of the 3D vector
func (v Vector3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns a normalized 3D vector
func (v Vector3) Normalize() Vector3 {
	length := v.Length()
	if IsZero(length) {
		return Vector3{X: 0, Y: 0, Z: 0}
	}
	return Vector3{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}