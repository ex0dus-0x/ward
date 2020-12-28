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
    FilePath string
    Protector *elf.File
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


// Create a new Injector interface to provision a runtime application
func NewInjector(binpath string, protector string) (*Injector, error) {

    // read bytes from target binary path
    targetBytes, err := ioutil.ReadFile(binpath)
    if err != nil {
        return nil, err
    }

    // parse protector as ELF binary
    binary, err := elf.Open(protector)
    if err != nil {
        return nil, err
    }

    return &Injector {
        protector,
        binary,
        targetBytes,
    }, nil
}


// Helper used to inject the original host into the new protector one through the commonly
// weaponized PT_NOTE to PT_LOAD infection vector.
func (inj *Injector) InjectBinary() {

    // TODO: define
    var injectSize uint64
    var fsize uint64

    for _, p := range inj.Protector.Progs {
        if p.Type == elf.PT_NOTE {

            // change to PT_LOAD segment
            p.Type = elf.PT_LOAD

            // allow read + exec
            p.Flags = elf.PF_R | elf.PF_X

            // define virtual memory offset for injected source
            p.Vaddr = 0xc000000 + uint64(fsize)

            // adjust size to account for injected code
            p.Filesz += injectSize
            p.Memsz += injectSize

            // set offset to end of original binary
            p.Off = uint64(fsize)
        }
    }
    inj.Protector.InsertionEOF = inj.Target
    return inj.Protector.WriteFile(inj.FilePath)
}

