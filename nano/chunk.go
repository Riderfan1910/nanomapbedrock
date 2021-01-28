package nano

type Chunks map[XZPos][]*Chunk

type Chunk struct {
	x, z	int
	data	[]byte
}

func NewChunk(data []byte, x, z int) *Chunk {
	return &Chunk {
		data: data,
		x: x,
		z: z,
	}
}
