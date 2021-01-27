package nanomap

import (
	"os"
	"fmt"
	"image"
	"image/png"
)

type Level struct {
	data map[XZPos]int
}

func Main(worldPath string, outputPath string) error {
	fmt.Println("#*-- nanomap --*#")
	
	world, err := OpenWorld(worldPath)
	if err != nil {
		return err
	}

	superChunks, err := world.GenerateSuperChunks()
	if err != nil {
		return err
	}

	w := int((superChunks.xmax - superChunks.xmin) * ChunkSizeXZ)
	h := int((superChunks.zmax - superChunks.zmin) * ChunkSizeXZ)
	
	if err := GenerateMaps(outputPath, superChunks, w, h); err != nil {
		return err
	}

	return nil
}

func GenerateMaps(outputPath string, superChunks *SuperChunks, w, h int) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for pos, superChunksData := range superChunks.data {
		for range superChunksData {
			xmin, zmin := int(superChunks.xmin), int(superChunks.zmin)
			for x := 0; x < ChunkSizeXZ; x++ {
				scanZ:
				for z := 0; z < ChunkSizeXZ; z++ {
					gx, gz := GetGlobalCoords(int(pos.x), int(pos.z), x, z)
					for y := 0; y < ChunkSizeY; y++ {
						img.Set(gx-(xmin*ChunkSizeXZ), gz-(zmin*ChunkSizeXZ), rgb(0x000000))
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
