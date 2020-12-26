package main

import (
    "os"
    "os/exec"
    "io/ioutil"
    "debug/elf"
)


type Injector struct {
    Protector []byte
    Target []byte
}

// Helper that compiles a new protection runtime application with `clang` for use with
// the defensive injector. Returns the bytes of the final blob compiled.
func Provision(name string) error {
    // TODO: find internal path to protector
    if err := os.Chdir("protector"); err != nil {
        return err
    }

    // create compilation command
    cmd := exec.Command(Compiler, "-Wall", "-O2", "-o",
        name, "protect.c", "runtime.c", "-lelf")

    // execute compilation routine to generate a new binary
    if err := cmd.Run(); err != nil {
        return err
    }
    return nil
}


// Helper to validate a blob of data as an ELF binary with a vacant PT_NOTE header
func BinaryCheck(data []byte) error {
    return nil
}


// Create a new Injector interface to provision a runtime application
func InitInjector(binpath string) (*Injector, error) {

    // read bytes from target binary path
    targetBytes, err := ioutil.ReadFile(binpath)
    if err != nil {
        return nil, err
    }

    // validate as an injectable ELF binary
    if err = BinaryCheck(targetBytes); err != nil {
        return nil, err
    }

    // provision new protector app runtime

    return &Injector {

    }, nil
}
