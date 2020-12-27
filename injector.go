package main

import (
    "os"
    "path"
    "errors"
    "runtime"
    "os/exec"
    "path/filepath"
    "io/ioutil"
    "debug/elf"
)


type Injector struct {
    Protector elf.File
    Target []byte
}

// Helper that compiles a new protection runtime application with `clang` for use with
// the defensive injector. Returns the bytes of the final blob compiled.
func Provision(name string) (*string, error) {

    // find directory to protector path in package
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        return nil, errors.New("Cannot find package path with protector.")
    }

    // get path to protector workspace
    protectorPath := filepath.Join(path.Dir(filename), "protector")
    if err := os.Chdir(protectorPath); err != nil {
        return nil, err
    }

    // create compilation command
    cmd := exec.Command(Compiler, "-Wall", "-O2", "-o",
        name + "_out", "protect.c", "runtime.c", "-lelf")

    // execute compilation routine to generate a new binary
    if err := cmd.Run(); err != nil {
        return nil, err
    }

    // get path to compiled executable
    protector := filepath.Join(path.Dir(filename), name + "_out")
    return &protector, nil
}


// Helper to validate a blob of data as an ELF binary with a vacant PT_NOTE header
func BinaryCheck(data []byte) error {
    return nil
}


// Helper used to inject the original host into the new protector one through the commonly
// weaponized PT_NOTE to PT_LOAD infection vector.
func InjectBinary(protector elf.File, host []byte) {
    /*
    var injectSize uint64
    var shellcode []byte

    // get original entry point to host
    originalEntryPoint := protector.FileHeader.Entry
    */
}


// Create a new Injector interface to provision a runtime application
func NewInjector(binpath string, protector string) (*Injector, error) {

    // read bytes from target binary path
    targetBytes, err := ioutil.ReadFile(binpath)
    if err != nil {
        return nil, err
    }

    // parse protector as ELF binary

    // validate protector as an injectable ELF binary
    if err = BinaryCheck(targetBytes); err != nil {
        return nil, err
    }

    return &Injector {

    }, nil
}
