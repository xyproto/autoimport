package importmatcher

import (
	"os"
	"os/exec"
	"regexp"
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
