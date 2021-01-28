package main

import (
	"fmt"
	"os"
	"github.com/pkg/errors"
	"nanomap/nano"
	"image"
	"image/png"
)

const (
	worldPath = "/Users/hal/Desktop/Assets/Projects/nanomap/world_3"
)

func main() {
	err := test()
	if err != nil {
		fmt.Printf("Error: %s", errors.WithStack(err))
	}
}

func test() error {
	fmt.Println("#*-- nanomap --*#")
	
	world, err := nano.Load(worldPath)
	if err != nil {
		return err
	}

	chunks, err := world.Chunks()
	if err != nil {
		return err
	}

	w := int((world.Xmax() - world.Xmin()) * nano.ChunkSizeXZ)
	h := int((world.Zmax() - world.Zmin()) * nano.ChunkSizeXZ)
	
	if err := GenerateMaps("./map.png", world, chunks, w, h); err != nil {
		return err
	}

	return nil
}

func GenerateMaps(outputPath string, world *nano.World, chunks nano.Chunks, w, h int) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for pos, superChunksData := range chunks {
		for range superChunksData {
			xmin, zmin := world.Xmin(), world.Zmin()
			for x := 0; x < nano.ChunkSizeXZ; x++ {
				scanZ:
				for z := 0; z < nano.ChunkSizeXZ; z++ {
					gx, gz := nano.GetGlobalCoords(pos.X, pos.Z, x, z)
					for y := 0; y < nano.ChunkSizeY; y++ {
						img.Set(gx-(xmin*nano.ChunkSizeXZ), gz-(zmin*nano.ChunkSizeXZ), nano.RGB(0x000000))
						continue scanZ
					}
				}
			}
		}
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return err
	}

	return nil
}
