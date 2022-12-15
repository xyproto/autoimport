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
		t.Fail()
	}
	if foundImport != "java.net.*" {
		t.Fail()
	}
}

func TestFileInputStream(t *testing.T) {
	impl, err := New(true)
	if err != nil {
		t.Errorf("Could not initialize ImportMatcher: %s\n", err)
	}
	foundClass, foundImport := impl.StarPath("FileInputS")
	if foundClass != "FileInputStream" {
		t.Fail()
	}
	if foundImport != "java.io.*" {
		t.Fail()
	}
}
