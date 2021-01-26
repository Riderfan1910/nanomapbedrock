package nanomap

import (
	"fmt"
	"path/filepath"
	"encoding/binary"
	"github.com/beito123/goleveldb/leveldb"
)

type World struct {
	Path		string
	DB 			*leveldb.DB
}

func OpenWorld(path string) (*World, error) {
	world := &World{
		Path: path,
		DB: nil,
	}
	var err error

	world.DB, err = leveldb.OpenFile(filepath.Join(path, "db"), nil)
	if err != nil {
		_ = world.DB.Close()
		return world, err
	}

	return world, nil
}

func ReadWorld(path string) (*World, map[XZPos][][]byte, error) {
	world, err := OpenWorld(path)
	if err != nil {
		return nil, nil, err
	}
	
	iter := world.DB.NewIterator(nil, nil)

	_chunks := map[XZPos][][]byte{}
	chunks := &Chunks{
		data: map[XZPos]*Chunk{},
	}

	chunksTotal := 0
	var x, y, z, xmin, xmax, zmin, zmax int32
	
	for iter.Next() {
		key := iter.Key()
		tmp := make([]byte, len(key))

		if len(key) > 8 && key[8] == 47 {
			// Update the min & max coordinates of the world.
			x, z, xmin, xmax, zmin, zmax = GetEdges(key, xmin, xmax, zmin, zmax)
			
			chunkData := iter.Value()

			_chunks[SetXZPos(x, z)] = append(_chunks[SetXZPos(x, z)], chunkData)
			chunk := SetChunk(_chunks[SetXZPos(x, z)], x, y, z)
			chunks.data[SetXZPos(x, z)] = chunk
			
			// ReadChunk(chunkData, 0)
			fmt.Println(key)

			chunksTotal++
		}
	
		copy(tmp, key)
	}
	iter.Release()

	world = &World{
		Path: world.Path,
		DB: world.DB,
	}

	err = iter.Error()
	if err != nil {
		return world, _chunks, err
	}

	defer world.DB.Close()

	return world, _chunks, nil
}

// Get min & max coordinates from leveldb key.
func GetEdges(key []byte, _xmin, _xmax, _zmin, _zmax int32) (x, z, xmin, xmax, zmin, zmax int32) {
	x = int32(binary.LittleEndian.Uint32(key[0:4]))
	z = int32(binary.LittleEndian.Uint32(key[4:8]))

	switch {
		case x <= _xmin:
			xmin = x
		case x >= _xmax:
			xmax = x
		case z <= _zmin:
			zmin = z
		case z >= _zmax:
			zmax = z
	}

	return
}
