package main

import(
	"fmt"
	"os"
	"github.com/opwire/opwire-testa/cli"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Testing execution has cracked, error: %s\n", err)
		}
	}()

	manifest := &Manifest{}

	cmd, err := cli.NewCommander(manifest)
	if err != nil {
		fmt.Printf("Cannot create Commander, error: %s\n", err.Error())
		os.Exit(2)
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Cannot process command, error: %s\n", err.Error())
		os.Exit(1)
	}
}
