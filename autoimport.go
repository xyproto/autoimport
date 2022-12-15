package importmatcher

import (
	"regexp"
)

// extractWords can extract words that starts with an uppercase letter from the given source code
func extractWords(sourceCode string) []string {
	re := regexp.MustCompile(`\b[A-Z][a-z]*([A-Z][a-z]*)*\b`)
	return re.FindAllString(sourceCode, -1)
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

// FindImports can find words that looks like classes, and then look up
// appropriate import package paths.
func (impM *ImportMatcher) FindImports(sourceCode string) []string {
	var foundImports []string
	for _, word := range extractWords(sourceCode) {
		foundPath := impM.ImportPathExact(word)
		if foundPath != "" {
			if !hasS(foundImports, foundPath) {
				foundImports = append(foundImports, foundPath)
			}
		}
	}
	return foundImports
}
