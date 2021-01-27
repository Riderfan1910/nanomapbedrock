package nanomap

import (
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

func (world *World) GenerateSuperChunks() (*SuperChunks, error) {
	iter := world.DB.NewIterator(nil, nil)

	superChunks := &SuperChunks{
		data: map[XZPos][]*SuperChunk{},
	}

	chunksTotal := 0
	var x, z, xmin, xmax, zmin, zmax int32
	
	for iter.Next() {
		key := iter.Key()
		tmp := make([]byte, len(key))

		if len(key) > 8 && key[8] == 47 {
			// Update the min & max coordinate of the world.
			x, z, xmin, xmax, zmin, zmax = GetEdges(key, xmin, xmax, zmin, zmax)

			chunkData := iter.Value()

			_superChunks := superChunks.data[SetXZPos(x, z)]
			superChunk := SetSuperChunk(chunkData, x, z)
			superChunks.data[SetXZPos(x, z)] = append(_superChunks, superChunk)
			superChunks.xmin, superChunks.xmax, superChunks.zmin, superChunks.zmax = xmin, xmax, zmin, zmax

			chunksTotal++
		}
	
		copy(tmp, key)
	}
	iter.Release()

	world = &World{
		Path: world.Path,
		DB: world.DB,
	}

	err := iter.Error()
	if err != nil {
		return superChunks, err
	}

	defer world.DB.Close()

	return superChunks, nil
}

// Get min & max coordinate from leveldb key.
func GetEdges(key []byte, _xmin, _xmax, _zmin, _zmax int32) (x, z, xmin, xmax, zmin, zmax int32) {
	x = int32(binary.LittleEndian.Uint32(key[0:4]))
	z = int32(binary.LittleEndian.Uint32(key[4:8]))

	if x <= _xmin {
		xmin = x
	} else {
		xmin = _xmin
	}

	if x >= _xmax {
		xmax = x
	} else {
		xmax = _xmax
	}

	if z <= _zmin {
		zmin = z
	} else {
		zmin = _zmin
	}

	if z >= _zmax {
		zmax = z
	} else {
		zmax = _zmax
	}

	return
}
