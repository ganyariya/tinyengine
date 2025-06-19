package math

// Camera2D represents a 2D camera for viewport management
type Camera2D struct {
	Position Vector2
	Zoom     float64
	Rotation float64 // in radians
}

// NewCamera2D creates a new 2D camera with default values
func NewCamera2D() Camera2D {
	return Camera2D{
		Position: Vector2{X: 0, Y: 0},
		Zoom:     1.0,
		Rotation: 0.0,
	}
}

// NewCamera2DWithValues creates a new 2D camera with specified values
func NewCamera2DWithValues(position Vector2, zoom, rotation float64) Camera2D {
	return Camera2D{
		Position: position,
		Zoom:     zoom,
		Rotation: rotation,
	}
}

// GetViewMatrix returns the view matrix for this camera
func (c Camera2D) GetViewMatrix() Matrix3x3 {
	// Create inverse transformation matrix for view
	// Camera transformation is the inverse of object transformation
	
	// Scale by inverse zoom
	scale := NewScaleMatrix3x3(c.Zoom, c.Zoom)
	
	// Rotate by negative rotation
	rotation := NewRotationMatrix3x3(-c.Rotation)
	
	// Translate by negative position
	translation := NewTranslationMatrix3x3(-c.Position.X, -c.Position.Y)
	
	// Combine transformations: Translate -> Rotate -> Scale
	return scale.Multiply(rotation).Multiply(translation)
}

// GetProjectionMatrix returns the projection matrix for screen space conversion
func (c Camera2D) GetProjectionMatrix(screenWidth, screenHeight float64) Matrix3x3 {
	// Convert from world space to screen space
	// Screen space: (0,0) at top-left, (width,height) at bottom-right
	// World space: (0,0) at center, (-1,-1) to (1,1) normalized
	
	halfWidth := screenWidth / 2.0
	halfHeight := screenHeight / 2.0
	
	// Scale and translate to screen coordinates
	return Matrix3x3{
		{halfWidth, 0, halfWidth},
		{0, -halfHeight, halfHeight}, // Flip Y-axis (screen Y goes down)
		{0, 0, 1},
	}
}

// GetViewProjectionMatrix returns the combined view-projection matrix
func (c Camera2D) GetViewProjectionMatrix(screenWidth, screenHeight float64) Matrix3x3 {
	view := c.GetViewMatrix()
	projection := c.GetProjectionMatrix(screenWidth, screenHeight)
	return projection.Multiply(view)
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c Camera2D) ScreenToWorld(screenPos Vector2, screenWidth, screenHeight float64) Vector2 {
	// Convert screen coordinates to normalized device coordinates (-1 to 1)
	halfWidth := screenWidth / 2.0
	halfHeight := screenHeight / 2.0
	
	normalized := Vector2{
		X: (screenPos.X - halfWidth) / halfWidth,
		Y: (halfHeight - screenPos.Y) / halfHeight, // Flip Y-axis
	}
	
	// Apply inverse camera transformation
	viewMatrix := c.GetViewMatrix()
	inverseView, err := viewMatrix.Inverse()
	if err != nil {
		return Vector2{X: 0, Y: 0} // Fallback to origin if matrix is singular
	}
	
	return inverseView.TransformPoint(normalized)
}

// WorldToScreen converts world coordinates to screen coordinates
func (c Camera2D) WorldToScreen(worldPos Vector2, screenWidth, screenHeight float64) Vector2 {
	// Apply camera transformation
	viewMatrix := c.GetViewMatrix()
	viewPos := viewMatrix.TransformPoint(worldPos)
	
	// Convert to screen coordinates
	halfWidth := screenWidth / 2.0
	halfHeight := screenHeight / 2.0
	
	return Vector2{
		X: viewPos.X*halfWidth + halfWidth,
		Y: halfHeight - viewPos.Y*halfHeight, // Flip Y-axis
	}
}

// SetPosition sets the camera position
func (c *Camera2D) SetPosition(position Vector2) {
	c.Position = position
}

// SetZoom sets the camera zoom level
func (c *Camera2D) SetZoom(zoom float64) {
	if zoom > 0 {
		c.Zoom = zoom
	}
}

// SetRotation sets the camera rotation in radians
func (c *Camera2D) SetRotation(rotation float64) {
	c.Rotation = rotation
}

// SetRotationDegrees sets the camera rotation in degrees
func (c *Camera2D) SetRotationDegrees(degrees float64) {
	c.Rotation = DegreesToRad(degrees)
}

// Move moves the camera by the given offset
func (c *Camera2D) Move(offset Vector2) {
	c.Position = c.Position.Add(offset)
}

// ZoomBy multiplies the current zoom by the given factor
func (c *Camera2D) ZoomBy(factor float64) {
	if factor > 0 {
		c.Zoom *= factor
	}
}

// Rotate rotates the camera by the given angle in radians
func (c *Camera2D) Rotate(angle float64) {
	c.Rotation += angle
}

// RotateDegrees rotates the camera by the given angle in degrees
func (c *Camera2D) RotateDegrees(degrees float64) {
	c.Rotation += DegreesToRad(degrees)
}

// GetBounds returns the world space bounds visible by this camera
func (c Camera2D) GetBounds(screenWidth, screenHeight float64) (Vector2, Vector2) {
	// Get the four corners of the screen in world space
	topLeft := c.ScreenToWorld(Vector2{X: 0, Y: 0}, screenWidth, screenHeight)
	topRight := c.ScreenToWorld(Vector2{X: screenWidth, Y: 0}, screenWidth, screenHeight)
	bottomLeft := c.ScreenToWorld(Vector2{X: 0, Y: screenHeight}, screenWidth, screenHeight)
	bottomRight := c.ScreenToWorld(Vector2{X: screenWidth, Y: screenHeight}, screenWidth, screenHeight)
	
	// Find min and max bounds
	minX := topLeft.X
	maxX := topLeft.X
	minY := topLeft.Y
	maxY := topLeft.Y
	
	corners := []Vector2{topRight, bottomLeft, bottomRight}
	for _, corner := range corners {
		if corner.X < minX {
			minX = corner.X
		}
		if corner.X > maxX {
			maxX = corner.X
		}
		if corner.Y < minY {
			minY = corner.Y
		}
		if corner.Y > maxY {
			maxY = corner.Y
		}
	}
	
	return Vector2{X: minX, Y: minY}, Vector2{X: maxX, Y: maxY}
}

// LookAt makes the camera look at a specific world position
func (c *Camera2D) LookAt(target Vector2) {
	c.Position = target
}

// FollowTarget smoothly follows a target position
func (c *Camera2D) FollowTarget(target Vector2, followSpeed float64, deltaTime float64) {
	if followSpeed <= 0 || deltaTime <= 0 {
		return
	}
	
	// Interpolate towards target
	direction := target.Sub(c.Position)
	distance := direction.Length()
	
	if distance > ZeroThreshold { // ジッターを避けるため非常に近い場合は無視
		maxDistance := followSpeed * deltaTime
		if distance > maxDistance {
			direction = direction.Normalize().Scale(maxDistance)
		}
		c.Position = c.Position.Add(direction)
	}
}