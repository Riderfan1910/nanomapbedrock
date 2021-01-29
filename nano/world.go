package nano

import (
	"path/filepath"
	"encoding/binary"
	"github.com/beito123/goleveldb/leveldb"
	"github.com/beito123/goleveldb/leveldb/util"
)

type World struct {
	Path						string
	Database				*leveldb.DB
	chunks 					Chunks
	subChunksTotal	int
	edges 					[]int
}

func Load(path string) (*World, error) {
	world := &World{
		Path: path,
		Database: nil,
		chunks: Chunks{},
		subChunksTotal: 0,
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

func (world *World) LoadChunk(x, z int) (*Chunk, error) {
	chunk := NewChunk(x, z)

	prefix := GetChunkKey(x, z, 47, -1)
	iter := world.Database.NewIterator(util.BytesPrefix(prefix), nil)
	
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		y := (key[len(key)-1]) & 15

		subChunk, err := ReadSubChunk(y, val)
		if err != nil {
			return nil, err
		}

		chunk.subChunks[y] = subChunk
	}

	iter.Release()

	err := iter.Error()
	if err != nil {
		return chunk, err
	}

	return chunk, nil
}

func (world *World) LoadChunks() (Chunks, error) {
	iter := world.Database.NewIterator(nil, nil)

	subChunksTotal := 0
	var x, z, xmin, xmax, zmin, zmax int
	
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		
		if len(key) > 8 && key[8] == 47 {
			x, z, xmin, xmax, zmin, zmax = CalcEdges(key, xmin, xmax, zmin, zmax)
			y := (key[len(key)-1]) & 15

			chunk, err := world.Chunk(x, z)
			if err != nil {
				return nil, err
			}

			subChunk, err := ReadSubChunk(y, val)
			if err != nil {
				return nil, err
			}

			chunk.subChunks[y] = subChunk
			world.chunks[SetXZPos(x, z)] = chunk

			subChunksTotal++
		}
	}

	iter.Release()

	err := iter.Error()
	if err != nil {
		return world.chunks, err
	}

	edges := []int{xmin, xmax, zmin, zmax}

	world.edges = edges
	world.subChunksTotal = subChunksTotal

	defer world.Database.Close()

	return world.chunks, nil
}

func (world *World) Chunk(x, z int) (*Chunk, error) {
	if world.chunks[SetXZPos(x, z)] != nil {
		return world.chunks[SetXZPos(x, z)], nil
	}

	chunk := NewChunk(x, z)

	return chunk, nil
}

func (world *World) Chunks() (Chunks, error) {
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

func (world *World) SubChunksTotal() int {
	return world.subChunksTotal
}

// CalcEdges returns min & max coordinate from leveldb key.
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
