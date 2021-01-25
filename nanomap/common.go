package nanomap

type XZPos struct {
	x, z int32
}

const (
	ChunkSizeXZ = 16
	ChunkSizeY  = 256
)

func SetXZPos(x, z int32) XZPos {
	return XZPos{
		x: x,
		z: z,
	}
}