package autoimport

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// isDir checks if the given path is a directory (could also be a symlink)
func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// isDir checks if the given path is a symlink
func isSymlink(path string) bool {
	_, err := os.Readlink(path)
	return err == nil
}

// exists checks if the given path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// followSymlink follows the given path
func followSymlink(path string) string {
	s, err := os.Readlink(path)
	if err != nil {
		return path
	}
	if !exists(s) && !strings.HasPrefix(s, "/") { // relative symlink
		s = filepath.Join(path, "..", s)
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

// keys will return the keys in a map[string]bool map as a string slice
func keys(m map[string]bool) []string {
	var keyStrings []string
	for k := range m {
		keyStrings = append(keyStrings, k)
	}
	return keyStrings
}

// unique will return all unique strings from a given string slice
func unique(xs []string) []string {
	// initialize the capacity of the map with the length of the given string slice
	uniqueStrings := make(map[string]bool, len(xs))
	for _, x := range xs {
		if _, ok := uniqueStrings[x]; !ok {
			uniqueStrings[x] = true
		}
	}
	return keys(uniqueStrings)
}
