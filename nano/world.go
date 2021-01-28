package nano

import (
	"path/filepath"
	"encoding/binary"
	"github.com/beito123/goleveldb/leveldb"
)

type World struct {
	Path			string
	Database	*leveldb.DB
	chunks 		Chunks
	edges 		[]int
}

func Load(path string) (*World, error) {
	world := &World{
		Path: path,
		Database: nil,
		chunks: Chunks{},
		edges: make([]int, 4),
	}
	var err error

	world.Database, err = leveldb.OpenFile(filepath.Join(path, "db"), nil)
	if err != nil {
		_ = world.Database.Close()
		return world, err
	}

	return world, nil
}

func (world *World) LoadChunks() (Chunks, error) {
	iter := world.Database.NewIterator(nil, nil)

	chunks := make(Chunks)

	subChunksTotal := 0
	var x, z, xmin, xmax, zmin, zmax int
	
	for iter.Next() {
		key := iter.Key()
		tmp := make([]byte, len(key))

		if len(key) > 8 && key[8] == 47 {
			// Update the min & max coordinate of the world.
			x, z, xmin, xmax, zmin, zmax = CalcEdges(key, xmin, xmax, zmin, zmax)

			subChunkData := iter.Value()

			_chunks := chunks[SetXZPos(x, z)]
			chunk := NewChunk(subChunkData, x, z)
			chunks[SetXZPos(x, z)] = append(_chunks, chunk)
			// superChunks.xmin, superChunks.xmax, superChunks.zmin, superChunks.zmax = xmin, xmax, zmin, zmax

			subChunksTotal++
		}
	
		copy(tmp, key)
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return chunks, err
	}

	edges := []int{xmin, xmax, zmin, zmax}

	world.chunks = chunks
	world.edges = edges

	defer world.Database.Close()

	return chunks, nil
}

func (world *World) Chunks() (Chunks, error) {
	if len(world.chunks) != 0 {
		return world.chunks, nil
	}

	chunks, err := world.LoadChunks()
	if err != nil {
		return nil, err
	}

	return chunks, nil
}

func (world *World) Edges() []int {
	return world.edges
}

func (world *World) Xmin() int {
	return world.edges[0]
}

func (world *World) Xmax() int {
	return world.edges[1]
}

func (world *World) Zmin() int {
	return world.edges[2]
}

func (world *World) Zmax() int {
	return world.edges[3]
}

// Get min & max coordinate from leveldb key.
func CalcEdges(key []byte, _xmin, _xmax, _zmin, _zmax int) (x, z, xmin, xmax, zmin, zmax int) {
	x = int(int32(binary.LittleEndian.Uint32(key[0:4])))
	z = int(int32(binary.LittleEndian.Uint32(key[4:8])))

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
