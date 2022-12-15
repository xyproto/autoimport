package main

import (
	"fmt"
	"os"

	"github.com/xyproto/importmatcher"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "syntax: w00t [START OF CLASSNAME]")
		os.Exit(1)

	}

	arg := os.Args[1]

	impl, err := importmatcher.New(true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)

	}

	foundClass, foundImport := impl.StarPath(arg)
	if foundClass == "" {
		fmt.Fprintf(os.Stderr, "found no class starting with %s\n", arg)
		os.Exit(1)
	}

	fmt.Printf("import %s; // %s\n", foundImport, foundClass)
}
