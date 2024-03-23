package main

import (
	"os"

	"github.com/debabky/pem-inclusion-prover-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
