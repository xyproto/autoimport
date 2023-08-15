package autoimport

import (
	"fmt"
	"testing"
)

func TestFindKotlin(_ *testing.T) {
	kotlinPath, err := FindKotlin()
	if err != nil {
		fmt.Printf("Could not find Kotlin: %s\n", err)
	}
	fmt.Printf("Found Kotlin at %s\n", kotlinPath)
}
