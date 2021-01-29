package main

import (
	"fmt"
	"os"
	"github.com/pkg/errors"
	"nanomap/nano"
	"image"
	"image/png"
	"image/color"
	"github.com/cheggaaa/pb/v3"
)

const (
	worldPath = "./data/2"
	outputPath = "./maps"
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

	fmt.Println("generating maps...")
	err = GenerateMaps(world)
	if err != nil {
		return err
	}
	fmt.Println("complete!")

	return nil
}

func GenerateMap(world *nano.World, cx, cz, cminx, cminz int, img *image.RGBA) (*image.RGBA, error) {

	chunk, err := world.LoadChunk(cminx, cminz)
	if err != nil {
		return nil, err
	}

	gcx, gcz := 16*cx, 16*cz

	for y := 0; y < 256; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				block, err := chunk.GetBlock(x, y, z)
				if err != nil {
					panic(err)
				}
				
				if block.Name() == "minecraft:air" {
					continue
				}

				bx, bz := gcx+x, gcz+z

				if block.Name() != "minecraft:stone" {
					img.Set(bx, bz, color.RGBA{125, 125, 125, 255})
				} else {
					img.Set(bx, bz, color.RGBA{78, 118, 42, 255})
				}
			}
		}
	}

	return img, nil
}

func GenerateMaps(world *nano.World) error {
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))

	tmpl := `{{percent .}} {{ bar . "[" "=" ">" "_" "]"}} {{counters . }} {{speed . | rndcolor }}`
	bar := pb.ProgressBarTemplate(tmpl).Start(256)

	var err error

	// 描画するチャンク(16x16)の最小値
	cminx, cminz := 0, 0

	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			img, err = GenerateMap(world, x, z, cminx+x, cminz+z, img)
			if err != nil {
				return err
			}

			bar.Increment()
		}
	}

	bar.Finish()

	defer world.Database.Close()

	f, err := os.Create(outputPath + "/map.png")
	if err != nil {
		return err
	}

	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return err
	}

	return nil
}
