package importmatcher

import (
	"fmt"
	"testing"
)

func TestFindKotlin(t *testing.T) {
	kotlinPath, err := FindKotlin()
	fmt.Println(kotlinPath, err)
}
