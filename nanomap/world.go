package nanomap

import (
	"path/filepath"
	"encoding/binary"
	"encoding/hex"
	"github.com/midnightfreddie/goleveldb/leveldb"
)

type Chunks struct {
	data	map[XZPos]*Chunk
}

type World struct {
	Path		string
	DB 			*leveldb.DB
	chunksTotal	int
}

func OpenWorld(path string) (*World, error) {
	result := &World{
		Path: path,
		DB: nil,
		chunksTotal: 0,
	}
	var err error

	result.DB, err = leveldb.OpenFile(filepath.Join(path, "db"), nil)
	if err != nil {
		_ = result.DB.Close()
		return result, err
	}

	return result, nil
}

func (world *World) Close() error {
	err := world.DB.Close()
	return err
}

func (_world *World) SetupChunk() (*World, *Chunks, error) {
	iter := _world.DB.NewIterator(nil, nil)
	chunksTotal := 0
	var x, y, z, xmin, xmax, zmin, zmax int32
	_chunks := map[string][][]byte{}
	chunks := &Chunks{
		data: map[XZPos]*Chunk{},
	}
	
	for iter.Next() {
		key := iter.Key()
		tmp := make([]byte, len(key))
		if len(key) > 8 && key[8] == 47 {
			subChunkPrefix := hex.EncodeToString(key[0:8])
			_chunks[subChunkPrefix] = append(_chunks[subChunkPrefix], tmp)

			// Update the min & max coordinates of the world.
			x, z, xmin, xmax, zmin, zmax = GetCoords(key, xmin, xmax, zmin, zmax)

			chunk := SetChunk(_chunks[subChunkPrefix], x, y, z)
			chunks.data[SetXZPos(x, z)] = chunk

			chunksTotal++
		}
		copy(tmp, key)
	}
	iter.Release()

	world := &World{
		Path: _world.Path,
		DB: _world.DB,
		chunksTotal: chunksTotal,
	}


	err := iter.Error()
	if err != nil {
		return world, chunks, err
	}

	return world, chunks, nil
}

// Get min & max coordinates from leveldb key.
func GetCoords(key []byte, xmin, xmax, zmin, zmax int32) (int32, int32, int32, int32, int32, int32) {
	_x := int32(binary.LittleEndian.Uint32(key[0:4]))
	_z := int32(binary.LittleEndian.Uint32(key[4:8]))

	switch {
		case _x <= xmin:
			xmin = _x
		case _x >= xmax:
			xmax = _x
		case _z <= zmin:
			zmin = _z
		case _z >= zmax:
			zmax = _z
	}

	return _x, _z, xmin, xmax, zmin, zmax
}
