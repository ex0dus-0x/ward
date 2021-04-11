package main

import (
    "os"
    "path"
    "bytes"
    "errors"
    "runtime"
    "os/exec"
    "io/ioutil"
    "path/filepath"

    // support for mutating and writing ELFs
    "github.com/Binject/debug/elf"
)

const (
    Compiler string = "clang"
    StubSectionName string = ".ward.protect"
)

// Helper that compiles a new protection runtime application with `clang` for use with
// the defensive injector. Returns the bytes of the final blob compiled.
func Provision(name string, overwrite bool) (*string, error) {

    // find directory to stub path in package
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        return nil, errors.New("Cannot find package path to provision stub program.")
    }

    // get current path before changing to stub path
    cwd, err := os.Getwd()
    if err != nil {
        return nil, err
    }

    // get path to stub workspace
    stubPath := filepath.Join(path.Dir(filename), "stub")
    if err := os.Chdir(stubPath); err != nil {
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
    cmd := exec.Command(Compiler, "-Wall", "-O2", "-D_FORTIFY_SOURCE=2", "-o",
        out, "main.c", "runtime.c", "-lelf")

    // execute compilation routine to generate a new binary
    if err := cmd.Run(); err != nil {
        return nil, err
    }

    // go back to original work directory
    if err := os.Chdir(cwd); err != nil {
        return nil, err
    }

    // return path to newly compiled stub executable
    return &out, nil
}

// Defines an Injector object that consumes a path to a compiled stub and
// target binary and creates a protected binary.
type Injector struct {
    Filepath string         // path to the stub host
    Filesize int64          // size of the stub host
    StubProgram *elf.File     // parsed ELF of the stub host
    Target []byte           // parsed bytes of the target binary to protect
}

// Create a new Injector interface to provision a runtime application
func NewInjector(binpath string, stub string) (*Injector, error) {

    // read bytes from target binary path we want to protect
    targetBytes, err := ioutil.ReadFile(binpath)
    if err != nil {
        return nil, err
    }

    // open to parse filesize
    f, err := os.Stat(stub)
    if err != nil {
        return nil, err
    }
    fsize := f.Size()

    // reopen and parse stub as ELF binary
    binary, err := elf.Open(stub)
    if err != nil {
        return nil, err
    }

    return &Injector {
        stub,
        fsize,
        binary,
        targetBytes,
    }, nil
}


// Helper used to modify the state of the shstrtab section in the ELF binary with the
// new section name string.
func (inj *Injector) OverwriteSection() error {

    var newShstrtab []byte
    var overwriteOffset int64

    // store index to section string table
    for _, sec := range inj.StubProgram.Sections {
        if sec.SectionHeader.Name == ".shstrtab" {
            shstrtab, err := sec.Data()
            if err != nil {
                return err
            }

            // change section name to one with same length
            newShstrtab = bytes.Replace(shstrtab, []byte(".note.ABI-tag"), []byte(StubSectionName), 1)
            overwriteOffset = int64(sec.SectionHeader.Offset)
            break
        }
    }

    // commit back to stub ELF by reopening for writing
    tempfile, err := os.OpenFile(inj.Filepath, os.O_RDWR, 0644)
    if err != nil {
        return err
    }

    // seek to offset and write
    if _, err := tempfile.WriteAt(newShstrtab, overwriteOffset); err != nil {
        return err
    }
    tempfile.Close()
    return nil
}


// Method used to inject the original host into the new stub one through the 
// commonly weaponized PT_NOTE to PT_LOAD infection vector.
func (inj *Injector) InjectBinary() error {

    // align code address to be congruent to file offset
    offset := (len(inj.Target) % 4096) - (0xc000000 % 4096)

    // overwrite the section name in shstrtab
    inj.OverwriteSection()

    // find code section to rename and rewrite for appended code
    for _, sec := range inj.StubProgram.Sections {
        if sec.SectionHeader.Name == StubSectionName {
            sec.SectionHeader.Type = elf.SHT_PROGBITS
            sec.SectionHeader.Flags = elf.SHF_ALLOC | elf.SHF_EXECINSTR
            sec.SectionHeader.Addr = 0xc000000 + uint64(inj.Filesize)
            sec.SectionHeader.Offset = uint64(offset)
            sec.SectionHeader.Size = uint64(inj.Filesize)
            sec.SectionHeader.Link = 0
            sec.SectionHeader.Info = 0
            sec.SectionHeader.Addralign = uint64(16)
            sec.SectionHeader.Entsize = 0
        }
    }

    // find a rewritable program header that has PT_NOTE segment, point to new section
    // with the injected target bytes we want to parse out in our stub.
    for _, seg := range inj.StubProgram.Progs {
        if seg.Type == elf.PT_NOTE {
            seg.Type = elf.PT_LOAD
            seg.Vaddr = 0xc000000 + uint64(inj.Filesize)
            seg.Flags = elf.PF_R | elf.PF_X
            seg.Filesz += uint64(len(inj.Target))
            seg.Memsz += uint64(len(inj.Target))
            seg.Off = uint64(inj.Filesize)
            break
        }
    }

    // append target binary to the end of the stub host
    inj.StubProgram.InsertionEOF = inj.Target

    // get bytes from final stub state
    elfBytes, err := inj.StubProgram.Bytes()
    if err != nil {
        return nil
    }

    // close stub after mutating and parsing bytes
    inj.StubProgram.Close()

    // overwrite original stub with changes in ELF format
    f, err := os.OpenFile(inj.Filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    if err != nil {
        return err
    }

    // write bytes and close
    f.Write(elfBytes)
    f.Close()
    return nil
}
