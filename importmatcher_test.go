package autoimport

import (
	"testing"
)

func TestString(t *testing.T) {
	impM, err := New(false)
	if err != nil {
		t.Fatalf("Could not initialize ImportMatcher: %s\n", err)
	}
	if impM.String() == "" {
		t.Fail()
	}
}

func TestSocketInputStream(t *testing.T) {
	impM, err := New(false)
	if err != nil {
		t.Fatalf("Could not initialize ImportMatcher: %s\n", err)
	}
	foundClass, foundImport := impM.StarPath("FileSto")
	if foundClass != "FileStore" {
		t.Fatalf("FileSto did not match with FileStore but with %s\n", foundClass)
	}
	if foundImport != "java.nio.file.*" {
		t.Fatalf("Expected java.nio.file.*, got %s\n", foundImport)
	}
}

func TestFileInputStream(t *testing.T) {
	impM, err := New(false)
	if err != nil {
		t.Fatalf("Could not initialize ImportMatcher: %s\n", err)
	}
	foundClass, foundImport := impM.StarPath("FileInputS")
	if foundClass != "FileInputStream" {
		t.Fatalf("FileInputS did not match with FileInputStream but with %s\n", foundClass)
	}
	if foundImport != "java.io.*" {
		t.Fatalf("Expected java.io.*, got %s\n", foundImport)
	}
}
