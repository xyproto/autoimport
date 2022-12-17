package importmatcher

import (
	"bytes"
	"errors"
	"os"
	"strings"
)

const kotlinPath = "/usr/share/kotlin/lib"

// FindKotlin finds the most likely location of a Kotlin installation
// (with subfolders with .jar files) on the system.
func FindKotlin() (string, error) {
	// 1. Find out if "kotlinc" is in the $PATH
	if kotlinExecutablePath := which("kotlinc"); kotlinExecutablePath != "" {
		// TODO: Find the definition of KOTLIN_HOME within the kotlinc script
		data, err := os.ReadFile(kotlinExecutablePath)
		if err != nil {
			return "", err
		}
		lines := bytes.Split(data, []byte{'\n'})
		for _, line := range lines {
			if bytes.Contains(line, []byte("KOTLIN_HOME")) && bytes.Count(line, []byte("=")) == 1 {
				fields := bytes.SplitN(line, []byte("="), 2)
				kotlinPath := strings.TrimSpace(string(fields[1]))
				return kotlinPath, nil
			}
		}

	}
	// 2. Consider typical path, for Arch Linux
	if isDir(kotlinPath) {
		return kotlinPath, nil
	}
	return "", errors.New("could not find an installation of Kotlin")
}
