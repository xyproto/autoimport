package autoimport

import (
	"testing"
)

// Mock implementation of ImportMatcher for testing
type MockImportMatcher struct{}

func (m *MockImportMatcher) StarPath(word string) (string, string) {
	switch word {
	case "ArrayList":
		return "ArrayList", "java.util.*"
	case "Map":
		return "Map", "java.util.*"
	default:
		return "", ""
	}
}

func TestFixImports(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple Java imports",
			input: `
package com.example;

import java.util.Map;
import java.util.ArrayList;

public class Main {
    ArrayList<String> list = new ArrayList<>();
    Map<String, String> map;
}
`,
			expected: `
package com.example;

import java.util.*; // ArrayList, Map

public class Main {
    ArrayList<String> list = new ArrayList<>();
    Map<String, String> map;
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ima, err := New(true)
			if err != nil {
				t.Fatalf("Error initializing the import matcher: %v", err)
			}
			got, err := ima.FixImports([]byte(tt.input), false)
			if err != nil {
				t.Fatalf("Error processing FixImports: %v", err)
			}
			if string(got) != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, got)
				//os.WriteFile("expected.java", []byte(tt.expected), 0644)
				//os.WriteFile("got.java", []byte(got), 0644)
			}
		})
	}
}
