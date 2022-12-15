package importmatcher

import (
	"regexp"
	"sort"
	"strings"
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
// appropriate import package paths. Ignores "java.lang." classes.
func (impM *ImportMatcher) FindImports(sourceCode string) []string {
	var foundImports []string
	for _, word := range extractWords(sourceCode) {
		foundPath := impM.ImportPathExact(word)
		if foundPath != "" && !strings.HasPrefix(foundPath, "java.lang.") {
			if !hasS(foundImports, foundPath) {
				foundImports = append(foundImports, foundPath)
			}
		}
	}
	return foundImports
}

// OrganizedImports generates import statements for packages that belongs to classes
// that are found in the given source code. If onlyJava is true, a semicolon is added
// after each line, and Kotlin jar files are not considered.
func (impM *ImportMatcher) OrganizedImports(sourceCode string, onlyJava bool) string {
	var sb strings.Builder
	imports := impM.FindImports(sourceCode)
	sort.Strings(imports)
	for _, importPackage := range imports {
		sb.WriteString("import " + importPackage)
		if onlyJava {
			sb.WriteString(";\n")
		} else {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
