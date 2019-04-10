package main

import(
	"fmt"
	"os"
	"github.com/opwire/opwire-qakit/cli"
)

func main() {
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
