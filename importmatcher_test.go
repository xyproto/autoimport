package importmatcher

import (
	"testing"
)

func TestString(t *testing.T) {
	impl, err := New(true)
	if err != nil {
		t.Errorf("Could not initialize ImportMatcher: %s\n", err)
	}
	if impl.String() == "" {
		t.Fail()
	}
}

func TestSocketInputStream(t *testing.T) {
	impl, err := New(true)
	if err != nil {
		t.Errorf("Could not initialize ImportMatcher: %s\n", err)
	}
	foundClass, foundImport := impl.StarPath("SocketInputS")
	if foundClass != "SocketInputStream" {
		t.Errorf("SocketInputS did not match with SocketInputStream but with %s\n", foundClass)
	}
	if foundImport != "java.net.*" {
		t.Errorf("Expected java.net.*, got %s\n", foundImport)
	}
}

func TestFileInputStream(t *testing.T) {
	impl, err := New(true)
	if err != nil {
		t.Errorf("Could not initialize ImportMatcher: %s\n", err)
	}
	foundClass, foundImport := impl.StarPath("FileInputS")
	if foundClass != "FileInputStream" {
		t.Errorf("FileInputS did not match with FileInputStream but with %s\n", foundClass)
	}
	if foundImport != "java.io.*" {
		t.Errorf("Expected java.io.*, got %s\n", foundImport)
	}
}
