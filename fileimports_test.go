package autoimport

import (
	"os"
	"testing"
)

// Sample imaementation of ImportMatcher for the purpose of testing.
type SampleImportMatcher struct{}

func (s *SampleImportMatcher) StarPath(word string) (string, string) {
	switch word {
	case "ArrayList":
		return "ArrayList", "java.util.*"
	case "HashMap":
		return "HashMap", "java.util.*"
	default:
		return "", ""
	}
}

func TestFileImports(t *testing.T) {
	// Mock java content
	javaContent := `
package com.example;

import java.util.Map;
import java.util.List;

public class Sample {
    List<String> list = new ArrayList<>();
    Map<String, String> map = new HashMap<>();
}
`
	// Write to a temporary file
	tmpfile, err := os.CreateTemp("", "sample.*.java")
	if err != nil {
		t.Fatalf("Could not create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // cleanup

	_, err = tmpfile.Write([]byte(javaContent))
	if err != nil {
		t.Fatalf("Could not write to temporary file: %v", err)
	}
	tmpfile.Close()

	// Create instance of our test ImportMatcher
	ima, err := New(true)
	if err != nil {
		t.Fatalf("Error initializing the import matcher: %v", err)
	}

	const verbose = false

	// Get imports
	imports, err := ima.FileImports(tmpfile.Name(), verbose)
	if err != nil {
		t.Fatalf("Error generating imports: %v", err)
	}

	// Expected imports based on the provided imaementation of ImportMatcher
	expected := `import java.util.*; // ArrayList, HashMap, List, Map`

	if imports != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, imports)
	}
}
