package nanomap

import (
	"fmt"
)

type Level struct {
	data map[XZPos]int
}

func Main(path string) error {
	fmt.Println("#*-- nanomap --*#")
	
	world, err := OpenWorld(path)
	if err != nil {
		return err
	}

	superChunks, err := world.GenerateSuperShunks()
	if err != nil {
		return err
	}

	fmt.Println(superChunks)

	return nil
}
