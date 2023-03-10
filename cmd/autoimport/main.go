package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/xyproto/autoimport"
)

const versionString = "autoimport 1.0.0"

type Args struct {
	StartOfClassName  string `arg:"positional,required"`
	ShortestMatchOnly bool   `arg:"-s,--shortest"`
	JavaOnly          bool   `arg:"-j,--java"`
	Exact             bool   `arg:"-e,--exact"`
}

func (Args) Version() string {
	return versionString
}

func main() {
	var args Args
	arg.MustParse(&args)

	impl, err := autoimport.New(args.JavaOnly)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)

	}

	if args.ShortestMatchOnly {
		// Output a single class + import, if found
		var foundClass, foundImport string
		foundClass, foundImport = impl.StarPath(args.StartOfClassName)
		if foundClass == "" {
			fmt.Fprintf(os.Stderr, "could not find the %s class\n", args.StartOfClassName)
			os.Exit(1)
		}
		fmt.Printf("import %s; // %s\n", foundImport, foundClass)
		return
	}

	// Output several classes and imports, if found
	var foundClasses, foundImports []string
	if args.Exact {
		foundClasses, foundImports = impl.StarPathAllExact(args.StartOfClassName)
		if len(foundClasses) == 0 {
			fmt.Fprintf(os.Stderr, "could not find the %s class\n", args.StartOfClassName)
			os.Exit(1)
		}
	} else {
		foundClasses, foundImports = impl.StarPathAll(args.StartOfClassName)
		if len(foundClasses) == 0 {
			fmt.Fprintf(os.Stderr, "found no class starting with %s\n", args.StartOfClassName)
			os.Exit(1)
		}
	}
	for i := range foundClasses {
		foundClass := foundClasses[i]
		foundImport := foundImports[i]
		fmt.Printf("import %s; // %s\n", foundImport, foundClass)
	}
}
