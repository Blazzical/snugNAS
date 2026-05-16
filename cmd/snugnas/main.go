package main

import (
	"fmt"
	"os"

	"github.com/Blazzical/snugNAS/internal/cli"
)

func main() {
	if err := cli.Root().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
