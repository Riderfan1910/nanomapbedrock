package nanomap

type SuperChunks struct {
	data	map[XZPos][]*SuperChunk
}

type SuperChunk struct {
	x, z int32
	data []byte
}

func SetSuperChunk(data []byte, x, z int32) *SuperChunk {
	return &SuperChunk {
		x: x,
		z: z,
		data: data,
	}
}
