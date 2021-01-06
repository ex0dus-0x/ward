package main

import (
    "os"
    "fmt"
    "log"
    "errors"
    "debug/elf"

    "github.com/urfave/cli/v2"
)

const (
    Compiler string = "clang"
    Description string = "Experimental security-hardened binary notary for ELFs"
)

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}


func main() {
    app := &cli.App {
        Name: "ward",
        Usage: Description,
        Commands: []*cli.Command{
            {
                Name: "protect",
                Usage: "Given a target binary, inject a self-protection runtime.",
                Flags: []cli.Flag{
                    &cli.BoolFlag{
                        Name: "overwrite",
                        Usage: "If set, overwrite the original target binary. (NOT RECOMMENDED)",
                        Aliases: []string{"o"},
                    },
                },
                Action: func(c *cli.Context) error {
                    binary := c.Args().First()
                    if binary == "" {
                        return errors.New("No binary specified for protection runtime injection.")
                    }

                    // check if valid path
                    if !fileExists(binary) {
                        return errors.New("Target ELF path does not exist for injection.")
                    }

                    // passive open to ensure path is valid ELF
                    if _, err := elf.Open(binary); err != nil {
                        return errors.New("Cannot open and parse target as ELF binary.")
                    }

                    overwrite := c.Bool("overwrite")

                    // start by provisioning a new protector host
                    fmt.Println("[*] Provisioning new protection executable")
                    protector, err := Provision(binary, overwrite)
                    if err != nil {
                        return err
                    }

                    // create new injector with target binary and new protector host
                    fmt.Println("[*] Injecting self-protection into", binary)
                    injector, err := NewInjector(binary, *protector)
                    if err != nil {
                        return err
                    }

                    // run PT_NOTE injection vector to inject target binary into host
                    injector.InjectBinary()
                    fmt.Println("[*] Done! Find the protected application at", *protector)
                    return nil
                },
            },
            {
                Name: "verify",
                Usage: "Ensures a protected binary has proper integrity.",
                Action: func(c *cli.Context) error {
                    binary := c.Args().First()
                    if binary == "" {
                        return errors.New("No binary specified.")
                    }

                    // check if valid path
                    if !fileExists(binary) {
                        return errors.New("Target ELF path does not exist for verification.")
                    }

                    // passive open to ensure path is valid ELF
                    if _, err := elf.Open(binary); err != nil {
                        return errors.New("Cannot open and parse target as ELF binary.")
                    }

                    // calculate digital signature of file

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
