package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/xyproto/autoimport"
)

const versionString = "autoimport 1.0.0"

// Args defines the possible command line arguments
type Args struct {
	StartOfClassName  string `arg:"positional"`
	SourceFile        string `arg:"-f,--file"`
	ShortestMatchOnly bool   `arg:"-s,--shortest"`
	JavaOnly          bool   `arg:"-j,--java"`
	Exact             bool   `arg:"-e,--exact"`
	Verbose           bool   `arg:"-V,--verbose"`
	NoGlob            bool   `arg:"-n,--noglob"`
}

// Version will output the current program name and version
func (Args) Version() string {
	return versionString
}

func main() {
	var args Args
	arg.MustParse(&args)

	var ima *autoimport.ImportMatcher
	var err error

	if args.SourceFile != "" {
		if strings.HasSuffix(args.SourceFile, ".kt") {
			ima, err = autoimport.New(false)
		} else {
			ima, err = autoimport.New(true)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		imports, err := ima.FileImports(args.SourceFile, args.Verbose)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if args.NoGlob {
			for _, deGlobbedImport := range autoimport.DeGlob(imports) {
				fmt.Println(deGlobbedImport)
			}
		} else {
			fmt.Println(imports)
		}
		return
	}

	ima, err = autoimport.New(args.JavaOnly)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var foundClasses, foundImports []string
	if args.ShortestMatchOnly {
		// Output a single class + import, if found
		var foundClass, foundImport string
		foundClass, foundImport = ima.StarPath(args.StartOfClassName)
		if foundClass == "" {
			fmt.Fprintf(os.Stderr, "could not find the %s class\n", args.StartOfClassName)
			os.Exit(1)
		}
		fmt.Printf("import %s; // %s\n", foundImport, foundClass)
		return
	}

	// Output several classes and imports, if found

	if args.Exact {
		foundClasses, foundImports = ima.StarPathAllExact(args.StartOfClassName)
		if len(foundClasses) == 0 {
			fmt.Fprintf(os.Stderr, "could not find the %s class\n", args.StartOfClassName)
			os.Exit(1)
		}
	} else {
		foundClasses, foundImports = ima.StarPathAll(args.StartOfClassName)
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
