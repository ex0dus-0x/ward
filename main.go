package main

import (
    "os"
    "fmt"
    "log"
    "errors"

    "github.com/urfave/cli/v2"
)

const (
    Compiler string = "clang"
    Description string = "Security-hardened binary notary"
)



func main() {
    app := &cli.App {
        Name: "ward",
        Usage: Description,
        Commands: []*cli.Command{
            {
                Name: "protect",
                Usage: "Given a target binary, inject a self-protection runtime.",
                Action: func(c *cli.Context) error {
                    binary := c.Args().First()
                    if binary == "" {
                        return errors.New("No binary specified for protection runtime injection.")
                    }

                    // check if valid path and valid ELF binary

                    // if everything is set, provision new

                    fmt.Println("[*] Injecting self-protection into ", binary)
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
