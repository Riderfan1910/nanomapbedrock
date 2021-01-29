package nano

import (
	"fmt"
	"github.com/beito123/binary"
)

func NewSubChunk(y byte) *SubChunk {
	return &SubChunk{
		y: y,
	}
}

func (subChunk *SubChunk) GetBlockStorage(index int) (*BlockStorage, bool) {
	if index >= len(subChunk.data) || index < 0 {
		return nil, false
	}

	return subChunk.data[index], true
}

func (subChunk *SubChunk) GetBlock(x, y, z, index int) (*RawBlockState, error) {
	storage, ok := subChunk.GetBlockStorage(index)
	if !ok {
		return nil, fmt.Errorf("invaild storage index")
	}

	return storage.Palettes[storage.Blocks[storage.At(x, y, z)]], nil
} 

func ReadSubChunk(y byte, data []byte) (*SubChunk, error) {
	subChunk := NewSubChunk(y)

	stream := binary.NewStreamBytes(data)

	ver, err := stream.Byte()
	if err != nil {
		return nil, err
	}

	switch ver {
		case 1:
			storage, err := ReadBlockStorage(stream)
			if err != nil {
				return nil, err
			}

			subChunk.data = append(subChunk.data, storage)
		case 8:
			numStorage, err := stream.Byte()
			if err != nil {
				return nil, err
			}
	
			for i := 0; i < int(numStorage); i++ {
				storage, err := ReadBlockStorage(stream)
				if err != nil {
					return nil, err
				}
				
				subChunk.data = append(subChunk.data, storage)
			}
		default:
			return nil, fmt.Errorf("unsupported version: version %d", ver)
	}
	
	return subChunk, nil
}
