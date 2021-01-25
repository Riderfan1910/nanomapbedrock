package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli/v2"
	"nanomap/nanomap"
)

func main() {
	app := cli.NewApp()
	app.Name = "nanomap"
	app.Version = "0.1.0"
	app.Copyright = "(c) 2021 Hal Sakuragi"
	app.Usage = "A simple map viewer for Minecraft: Bedrock Edition."
	
	app.Flags = []cli.Flag {
		&cli.StringFlag {
			Name: "world",
			Aliases: []string{"w"},
			Usage: "Path of the Minecraft world folder.",
			Required: true,
		},
	}
	
	app.Action = func (context *cli.Context) error {
		_world, err := nanomap.OpenWorld(context.String("world"))
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		defer _world.Close()

		world, chunks, err := _world.SetupChunk()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		fmt.Println(world, chunks)

    return nil
	}

	app.Run(os.Args)
}
