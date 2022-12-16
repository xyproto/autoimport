package importmatcher

import (
	"fmt"
	"testing"
)

func TestFindJava(t *testing.T) {
	javaPath, err := FindJava()
	fmt.Println(javaPath, err)
}
