package nano

type Block struct {
	ID uint16
	Data byte
}

type RGB uint32

func (x RGB) RGBA() (r, g, b, a uint32) {
	a = 0xffff
	r = uint32((x >> 16) << 8)
	g = uint32(((x >> 8) & 0xff) << 8)
	b = uint32((x & 0xff) << 8)
	return
}

func GetGlobalCoords(cx, cz, rbx, rbz int) (gx, gz int) {
	gx = cx*ChunkSizeXZ + rbx
	gz = cz*ChunkSizeXZ + rbz
	return
}

func (block *RawBlockState) Name() string {
	return block.name
}
