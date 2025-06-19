package renderer

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// BufferPool manages reusable OpenGL buffers
type BufferPool struct {
	vaoPool chan uint32
	vboPool chan uint32
	eboPool chan uint32
	maxSize int
}

// NewBufferPool creates a new buffer pool
func NewBufferPool(maxSize int) *BufferPool {
	return &BufferPool{
		vaoPool: make(chan uint32, maxSize),
		vboPool: make(chan uint32, maxSize),
		eboPool: make(chan uint32, maxSize),
		maxSize: maxSize,
	}
}

// GetVAO gets a VAO from the pool or creates a new one
func (bp *BufferPool) GetVAO() uint32 {
	select {
	case vao := <-bp.vaoPool:
		return vao
	default:
		var vao uint32
		gl.GenVertexArrays(1, &vao)
		return vao
	}
}

// GetVBO gets a VBO from the pool or creates a new one
func (bp *BufferPool) GetVBO() uint32 {
	select {
	case vbo := <-bp.vboPool:
		return vbo
	default:
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		return vbo
	}
}

// GetEBO gets an EBO from the pool or creates a new one
func (bp *BufferPool) GetEBO() uint32 {
	select {
	case ebo := <-bp.eboPool:
		return ebo
	default:
		var ebo uint32
		gl.GenBuffers(1, &ebo)
		return ebo
	}
}

// ReturnVAO returns a VAO to the pool
func (bp *BufferPool) ReturnVAO(vao uint32) {
	select {
	case bp.vaoPool <- vao:
	default:
		// Pool is full, delete the buffer
		gl.DeleteVertexArrays(1, &vao)
	}
}

// ReturnVBO returns a VBO to the pool
func (bp *BufferPool) ReturnVBO(vbo uint32) {
	select {
	case bp.vboPool <- vbo:
	default:
		// Pool is full, delete the buffer
		gl.DeleteBuffers(1, &vbo)
	}
}

// ReturnEBO returns an EBO to the pool
func (bp *BufferPool) ReturnEBO(ebo uint32) {
	select {
	case bp.eboPool <- ebo:
	default:
		// Pool is full, delete the buffer
		gl.DeleteBuffers(1, &ebo)
	}
}

// Destroy cleans up all buffers in the pool
func (bp *BufferPool) Destroy() {
	// Clean VAO pool
	for {
		select {
		case vao := <-bp.vaoPool:
			gl.DeleteVertexArrays(1, &vao)
		default:
			goto cleanVBO
		}
	}

cleanVBO:
	// Clean VBO pool
	for {
		select {
		case vbo := <-bp.vboPool:
			gl.DeleteBuffers(1, &vbo)
		default:
			goto cleanEBO
		}
	}

cleanEBO:
	// Clean EBO pool
	for {
		select {
		case ebo := <-bp.eboPool:
			gl.DeleteBuffers(1, &ebo)
		default:
			return
		}
	}
}