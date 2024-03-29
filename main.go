package main

import (
	"debug/elf"
	"errors"
	"flag"
	"log"
	"os"
)

func RunWard() error {
	overwrite := flag.Bool("overwrite", false, "Overwrite original executable path (NOT RECOMMENDED)")
	compress := flag.Bool("compress", true, "Compress executable when packing with zlib (on by default)")
	protect := flag.Bool("protect", true, "Incorporate code injection prevention (on by default)")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		return errors.New("Specify a single ELF binary for packing.")
	}
	binary := args[0]
	log.Println("Starting ward to pack", binary)

	_, err := os.Stat(binary)
	if os.IsNotExist(err) {
		return errors.New("ELF file not found at path.")
	}

	log.Println("Checking if valid ELF binary")
	if _, err := elf.Open(binary); err != nil {
		return errors.New("Cannot open and parse target as ELF binary.")
	}

	log.Println("Provisioning stub program for packing")
	protector, err := Provision(binary, *overwrite, *compress, *protect)
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
}

func main() {
	if err := RunWard(); err != nil {
		log.Fatal(err)
	}
}
