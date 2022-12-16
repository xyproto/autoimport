package importmatcher

import (
	"errors"
	"fmt"
)

const kotlinPath = "/usr/share/kotlin/lib"

// FindKotlin finds the most likely location of a Kotlin installation
// (with subfolders with .jar files) on the system.
func FindKotlin() (string, error) {
	// 1. Find out if "kotlinc" is in the $PATH
	if javaExecutablePath := which("kotlinc"); javaExecutablePath != "" {
		// TODO: Find the definition of KOTLIN_HOME within the kotlinc script
		fmt.Println("TO IMPLEMENT: examine kotlinc")
	}
	// 2. Consider typical path, for Arch Linux
	if isDir(kotlinPath) {
		return kotlinPath, nil
	}
	return "", errors.New("could not find an installation of Kotlin")
}
