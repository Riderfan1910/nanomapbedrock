package nanomap

type SuperChunks struct {
	data	map[XZPos][]*SuperChunk
	xmin, xmax, zmin, zmax int32
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

func (superChunk *SuperChunk) GetBlock(x, y, z int) error {
	return nil
}
