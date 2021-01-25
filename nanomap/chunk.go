package nanomap

import (
	"fmt"
)

type Chunks struct {
	data	map[XZPos]*Chunk
}

type Chunk struct {
	x, y, z	int32
	data	[][]byte
}

func SetPreChunk(data [][]byte, x, y, z int32) *Chunk {
	return &Chunk {
		x: x,
		y: y,
		z: z,
		data: data,
	}
}

func ReadChunk(data []byte, chunk *Chunk, y int32) error {
	_offset := 0
	// subChunkVersion := data[_offset]
	// subChankYOffset := 16 * y
	fmt.Println(data[_offset])
	return nil
}
