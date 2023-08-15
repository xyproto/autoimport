package autoimport

import (
	"fmt"
	"testing"
)

func TestFindJava(_ *testing.T) {
	javaPath, err := FindJava()
	if err != nil {
		fmt.Printf("Could not find Java: %s\n", err)
	}
	fmt.Printf("Found Java at %s\n", javaPath)
}
