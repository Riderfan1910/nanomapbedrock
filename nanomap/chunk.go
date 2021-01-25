package nanomap

type Chunk struct {
	x, y, z	int32
	data	[][]byte
}

func SetChunk(data [][]byte, x, y, z int32) *Chunk {
	return &Chunk {
		x: x,
		y: y,
		z: z,
		data: data,
	}
}
