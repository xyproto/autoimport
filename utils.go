package autoimport

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// which tries to find the given executable name in the $PATH
// Returns an empty string if not found.
func which(executable string) string {
	p, err := exec.LookPath(executable)
	if err != nil {
		return ""
	}
	return p
}

// hasS checks if the given string slice contains the given string
func hasS(sl []string, e string) bool {
	for _, s := range sl {
		if s == e {
			return true
		}
	}
	return false
}

// extractWords can extract words that starts with an uppercase letter from the given source code
func extractWords(sourceCode string) []string {
	re := regexp.MustCompile(`\b[A-Z][a-z]*([A-Z][a-z]*)*\b`)
	return re.FindAllString(sourceCode, -1)
}

// isDir checks if the given path is a directory
func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func isSymlink(path string) bool {
	_, err := os.Readlink(path)
	if err != nil {
		return false
	}
	return true
}

func followSymlink(path string) string {
	s, err := os.Readlink(path)
	if err != nil {
		return path
	}
	return s
}

func FindInEtcEnvironment(envVarName string) (string, error) {
	// Find the definition of ie. JAVA_HOME within /etc/environment
	data, err := os.ReadFile("/etc/environment")
	if err != nil {
		return "", err
	}
	lines := bytes.Split(data, []byte{'\n'})
	for _, line := range lines {
		if bytes.Contains(line, []byte(envVarName)) && bytes.Count(line, []byte("=")) == 1 {
			fields := bytes.SplitN(line, []byte("="), 2)
			javaPath := strings.TrimSpace(string(fields[1]))
			if !isDir(javaPath) {
				continue
			}
			return javaPath, nil
		}
	}
	return "", fmt.Errorf("could not find the value of %s in /etc/environment", envVarName)
}
