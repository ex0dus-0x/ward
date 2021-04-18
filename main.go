package main

import (
	"debug/elf"
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ward",
		Usage: "Dumb ELF packer",
		Commands: []*cli.Command{
			{
				Name:  "pack",
				Usage: "Pack a target binary, and inject a self-protection runtime.",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "overwrite",
						Usage:   "If set, overwrite the original target binary (NOT RECOMMENDED).",
						Aliases: []string{"o"},
					},
				},
				Action: func(c *cli.Context) error {
					log.Println("Starting up ward")

					binary := c.Args().First()
					if binary == "" {
						return errors.New("No binary specified for packing.")
					}

					_, err := os.Stat(binary)
					if os.IsNotExist(err) {
						return errors.New("ELF file not found at path.")
					}

					log.Println("Checking if valid ELF binary")
					if _, err := elf.Open(binary); err != nil {
						return errors.New("Cannot open and parse target as ELF binary.")
					}

					overwrite := c.Bool("overwrite")

					log.Println("Provisioning stub program for packing")
					protector, err := Provision(binary, overwrite)
					if err != nil {
						return err
					}

					log.Println("Packing original executable into stub", binary)
					injector, err := NewInjector(binary, *protector)
					if err != nil {
						return err
					}

					injector.InjectBinary()
					log.Println("Done! Find the packed application at", *protector)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
