package nano

type XZPos struct {
	X, Z int
}

func SetXZPos(x, z int) XZPos {
	return XZPos{
		X: x,
		Z: z,
	}
}

const (
	ChunkSizeXZ = 16
	ChunkSizeY = 256
)
