package nanomap

import (
	// "fmt"
	"math"
	// "encoding/binary"
	// "github.com/hrqsn/nbt"
)

type Chunks struct {
	data	map[XZPos]*Chunk
}

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

func ReadChunk(data []byte, y int) error {
	_offset := 0
	subChunkVersion := data[_offset]
	_offset++
	// subChankYOffset := 16 * y
	storages := 1

	switch subChunkVersion {
		case 8:
			storages = int(data[_offset])
			_offset++

			for storage := 0; storage < storages; storage++ {
				paletteAndFlag := int8(data[_offset])
				_offset++
				isRuntime := (paletteAndFlag & 1) != 0
				bitsPerBlock := paletteAndFlag >> 1
        blocksPerWord := math.Floor(32 / float64(bitsPerBlock))
				wordCount := int(math.Ceil(4096 / blocksPerWord))
				
				// indexBlocks := _offset
				_offset += wordCount * 4

				if !isRuntime {
					// localPalette.size = int(binary.LittleEndian.Uint16(data[_offset:]))
					_offset += 4

					// for paletteID := 0; paletteID < localPalette.size; paletteID++ {
					// 	stream, err := nbt.FromBytes(data[_offset:], nbt.LittleEndian)
					// 	if err != nil {
					// 		panic(err)
					// 	}

					// 	// tag, err := stream.ReadTag()
					// 	// if err != nil {
					// 	// 	panic(err)
					// 	// }

					// 	// str, err := tag.ToString()
					// 	// if err != nil {
					// 	// 	panic(err)
					// 	// }

					// 	_offset += stream.GetPayload()
					// }
				}
			}
			break
		default:
			break
	}

	return nil
}
