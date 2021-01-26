package nanomap

import (
	"fmt"
)

func Main(path string) error {
	fmt.Println("#*-- nanomap --*#")
	
	_, _, err := ReadWorld(path)
	if err != nil {
		return err
	}

	return nil
}
