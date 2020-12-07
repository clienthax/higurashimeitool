package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "HigurashiMeiTool",
		Usage: "A tool for reading Higurashi Mei's resource files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "list",
				Usage:    "The list.bin file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "data",
				Usage:    "The data.bin file",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "hash",
				Usage: "The hash.bin file",
			},
			&cli.StringFlag{
				Name:     "out",
				Usage:    "Output directory",
				Required: true,
			},
		},
		Action: process,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
