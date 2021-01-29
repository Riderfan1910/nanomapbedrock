package nano

import (
	"fmt"
	"github.com/beito123/binary"
	"github.com/beito123/nbt"
	"math"
	"strings"
)

type Chunks map[XZPos]*Chunk

type Chunk struct {
	x, z	int
	subChunks []*SubChunk
}

type SubChunk struct {
	y byte
	data []*BlockStorage
}

type BlockStorage struct {
	Palettes []*RawBlockState
	Blocks   []uint16
}

type RawBlockState struct {
	name  string
	value int
}

func NewChunk(x, z int) *Chunk {
	return &Chunk {
		x: x,
		z: z,
		subChunks: make([]*SubChunk, 16),
	}
}

const BlockStorageSize = 16 * 16 * 16
func NewBlockStorage() *BlockStorage {
	return &BlockStorage{
		Blocks: make([]uint16, BlockStorageSize),
	}
}

func NewRawBlockState(name string, value int) *RawBlockState {
	return &RawBlockState{
		name:  strings.ToLower(name),
		value: value,
	}
}

func (BlockStorage) At(x, y, z int) int {
	return x<<8 | z<<4 | y
}

func (chunk *Chunk) GetSubChunk(index int) (*SubChunk, bool) {
	if index >= len(chunk.subChunks) {
		return nil, false
	}

	return chunk.subChunks[index], chunk.subChunks[index] != nil
}

func (chunk *Chunk) AtSubChunk(y int) (*SubChunk, bool) {
	return chunk.GetSubChunk(y / 16)
}

func (chunk *Chunk) GetBlock(x, y, z int) (*RawBlockState, error) {
	subChunk, ok := chunk.AtSubChunk(y)
	if !ok {
		return NewRawBlockState("minecraft:air", 0), nil
	}

	return subChunk.GetBlock(x, y&15, z, 0)
}

func (chunk *Chunk) SubChunks() []*SubChunk {
	return chunk.subChunks
}

func ReadBlockStorage(stream *binary.Stream) (*BlockStorage, error) {
	storage := NewBlockStorage()

	flags, err := stream.Byte()
	if err != nil {
		return nil, err
	}

	bitsPerBlock := int(flags >> 1)
	isRuntime := (flags & 0x01) != 0

	if bitsPerBlock > 16 {
		return nil, fmt.Errorf("unsupported bits per block, wants 1-16 bits")
	}

	mask := uint16((1 << uint(bitsPerBlock)) - 1)

	wordBits := 8 * 4 // 1byte * 4
	blocksPerWord := wordBits / bitsPerBlock

	wordCount := int(math.Ceil(float64(BlockStorageSize) / float64(blocksPerWord)))

	count := 0
	for i := 0; i < wordCount; i++ {
		word, err := stream.LInt()
		if err != nil {
			return nil, err
		}

		for j := 0; j < blocksPerWord && count < BlockStorageSize; j++ {
			id := uint16(word>>uint(j*bitsPerBlock)) & mask

			storage.Blocks[count] = id

			count++
		}
	}

	paletteSize, err := stream.LInt()
	if err != nil {
		return nil, err
	}

	if isRuntime {
		return nil, fmt.Errorf("unsupported runtime id")
	}

	nbtStream := nbt.NewStreamBytes(nbt.LittleEndian, stream.Bytes())

	for i := 0; i < int(paletteSize); i++ {
		tag, err := nbtStream.ReadTag()
		if err != nil {
			return nil, err
		}

		com, ok := tag.(*nbt.Compound)
		if !ok {
			return nil, fmt.Errorf("unexpected tag %s (%d)", tag.Name(), tag.ID())
		}

		name, err := com.GetString("name")
		if err != nil {
			return nil, err
		}

		val := 0

		state := NewRawBlockState(name, int(val))

		storage.Palettes = append(storage.Palettes, state)
	}

	stream.Skip(nbtStream.Stream.Off())

	return storage, nil
}

func GetChunkKey(x, z int, tag byte, sid int) []byte {
	base := []byte{
		byte(x),
		byte(x >> 8),
		byte(x >> 16),
		byte(x >> 24),
		byte(z),
		byte(z >> 8),
		byte(z >> 16),
		byte(z >> 24),
	}

	return []byte{
		base[0],
		base[1],
		base[2],
		base[3],
		base[4],
		base[5],
		base[6],
		base[7],
		tag,
	}
}
