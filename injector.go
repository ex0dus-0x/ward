package main

import (
    "os"
    "path"
    "errors"
    "runtime"
    "os/exec"
    "io/ioutil"
    "debug/elf"
    "path/filepath"
)


type Injector struct {
    FilePath string
    Protector *elf.File
    Target []byte
}

// Helper that compiles a new protection runtime application with `clang` for use with
// the defensive injector. Returns the bytes of the final blob compiled.
func Provision(name string, overwrite bool) (*string, error) {

    // find directory to protector path in package
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        return nil, errors.New("Cannot find package path with protector.")
    }

    // get current path before changing to protector path
    cwd, err := os.Getwd()
    if err != nil {
        return nil, err
    }

    // get path to protector workspace
    protectorPath := filepath.Join(path.Dir(filename), "protector")
    if err := os.Chdir(protectorPath); err != nil {
        return nil, err
    }

    // if overwrite is set, rewrite the original path (might be dangerous)
    var out string
    if overwrite != true {
        out = filepath.Join(cwd, name + "_out")
    } else {
        out, err = filepath.Abs(name)
        if err != nil {
            return nil, err
        }
    }

    // create compilation command
    cmd := exec.Command(Compiler, "-Wall", "-O2", "-o",
        out, "protector.c", "runtime.c", "-lelf")

    // execute compilation routine to generate a new binary
    if err := cmd.Run(); err != nil {
        return nil, err
    }

    // go back to original work directory
    if err := os.Chdir(cwd); err != nil {
        return nil, err
    }

    // return path to newly compiled protector executable
    return &out, nil
}


// Create a new Injector interface to provision a runtime application
func NewInjector(binpath string, protector string) (*Injector, error) {

    // read bytes from target binary path we want to protect
    targetBytes, err := ioutil.ReadFile(binpath)
    if err != nil {
        return nil, err
    }

    // before parsing as ELF, write payload blob to end of file, save and open again,
    // since the ELF parser in Golang doesn't enable writing at EOF anymore
    f, err := os.OpenFile(protector, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        return nil, err
    }

    // write target bytes and close
    if _, err = f.Write(targetBytes); err != nil {
        return nil, err
    }
    f.Close()

    // reopen and parse protector as ELF binary
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
    injectSize := len(inj.Target)

    // TODO
    var fsize int

    // modify the protector's PT_NOTE segment
    for _, p := range inj.Protector.Progs {
        if p.Type == elf.PT_NOTE {

            // change to PT_LOAD segment
            p.Type = elf.PT_LOAD

            // allow read + exec
            p.Flags = elf.PF_R | elf.PF_X

            // define virtual memory offset for injected source
            p.Vaddr = 0xc000000 + uint64(injectSize)

            // adjust size to account for injected code
            p.Filesz += uint64(injectSize)
            p.Memsz += uint64(injectSize)

            // set offset to end of original binary
            p.Off = uint64(fsize)
        }
    }
}

