// Package importmatcher tries to find which import should be used, given the start of a class name
package importmatcher

import (
	"archive/zip"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"unicode"

	"github.com/stretchr/powerwalk"
	"github.com/xyproto/env"
)

// ImportMatcher is a struct that contains a list of JAR file paths,
// and a lookup map from class names to class paths, which is populated
// when New or NewCustom is called.
type ImportMatcher struct {
	JARPaths []string          // list of paths to examine for .jar files
	classMap map[string]string // map from class name to class path. Shortest class path "wins".
	onlyJava bool              // only Java, or Kotlin too?
}

var numCPU = runtime.NumCPU()

// New creates a new ImportMatcher. If onlyJava is false, /usr/share/kotlin/lib will be added to the .jar file search path.
func New(onlyJava bool) (*ImportMatcher, error) {
	var JARPaths = []string{
		env.Str("JAVA_HOME", "/usr/lib/jvm/default"),
		"/usr/lib/jvm/default-java",
	}
	if javaExecutablePath := which("java"); javaExecutablePath != "" {
		JARPaths = append(JARPaths, filepath.Dir(javaExecutablePath))
	}
	if !onlyJava {
		JARPaths = append(JARPaths, "/usr/share/kotlin/lib")
	}
	return NewCustom(JARPaths, onlyJava)
}

// NewCustom creates a new ImportMatcher, given a slice of paths to search for .jar files.
func NewCustom(JARPaths []string, onlyJava bool) (*ImportMatcher, error) {
	var impM ImportMatcher
	impM.onlyJava = onlyJava

	impM.JARPaths = make([]string, len(JARPaths))
	for i := range JARPaths {
		JARPath := JARPaths[i]
		if !strings.HasSuffix(JARPath, "/") {
			JARPath += "/"
		}
		fi, err := os.Stat(JARPath)
		if err != nil {
			continue
		}
		if fi.IsDir() {
			impM.JARPaths[i] = JARPath
		}
	}

	if len(JARPaths) == 0 {
		return nil, errors.New("no paths to search for JAR files")
	}

	foundClasses := make(chan string, 512)
	impM.classMap = make(map[string]string)

	var (
		m                 sync.RWMutex
		existingClassPath string
		ok                bool
	)

	go func() {
		err := impM.findClasses(foundClasses)
		if err != nil {
			log.Printf("error: %s\n", err)
		}
		close(foundClasses)
	}()

	for classPath := range foundClasses {

		className := classPath
		if strings.Contains(classPath, ".") {
			fields := strings.Split(classPath, ".")
			lastField := fields[len(fields)-1]
			className = lastField
		}

		m.RLock()
		existingClassPath, ok = impM.classMap[className]
		m.RUnlock()

		if ok && existingClassPath != "" && len(existingClassPath) <= len(classPath) {
			continue
		}

		m.Lock()
		impM.classMap[className] = classPath
		m.Unlock()
	}

	return &impM, nil
}

// ClassMap returns the mapping from class names to class paths
func (impM *ImportMatcher) ClassMap() map[string]string {
	return impM.classMap
}

// readJAR returns a list of classes within the given .jar file,
// for instance "SomePath/SomePath/SomeClass"
func (impM *ImportMatcher) readJAR(filePath string, found chan string) error {
	readCloser, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	for _, f := range readCloser.File {

		if strings.HasSuffix(f.Name, ".class") || strings.HasSuffix(f.Name, ".CLASS") {

			className := strings.TrimSuffix(strings.TrimSuffix(f.Name, ".class"), ".CLASS")
			className = strings.ReplaceAll(className, "/", ".")
			className = strings.TrimSuffix(className, "$1")
			className = strings.TrimSuffix(className, "$1")
			if pos := strings.Index(className, "$"); pos >= 0 {
				className = className[:pos]
			}

			if className == "" {
				continue
			}

			// Filter out class names that are only lowercase (and '.')
			allLower := true
			for _, r := range className {
				if !unicode.IsLower(r) && r != '.' {
					allLower = false
				}
			}
			if allLower {
				continue
			}

			found <- className
		}
	}

	return nil
}

// FindClassesInJAR will search the given JAR file for classes,
// and pass them as strings down the "found" chan.
func (impM *ImportMatcher) FindClassesInJAR(JARPath string, found chan string) error {
	return powerwalk.WalkLimit(JARPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		filePath := path
		fileName := info.Name()

		if filepath.Ext(fileName) != ".jar" && filepath.Ext(fileName) != ".JAR" {
			return nil
		}

		return impM.readJAR(filePath, found)
	}, numCPU)
}

func (impM *ImportMatcher) findClasses(found chan string) error {
	for _, JARPath := range impM.JARPaths {
		if err := impM.FindClassesInJAR(JARPath, found); err != nil {
			log.Printf("error: %s\n", err)
		}
	}
	return nil
}

func (impM *ImportMatcher) String() string {
	var sb strings.Builder
	for className, classPath := range impM.classMap {
		sb.WriteString(className + ": " + classPath + "\n")
	}
	return sb.String()
}

// StarPath takes the start of the class name and tries to return the shortest
// found class name, and also the import path like "java.io.*"
// Returns empty strings if there are no matches.
func (impM *ImportMatcher) StarPath(startOfClassName string) (string, string) {
	shortestClassName := ""
	shortestImportPath := ""
	for className, classPath := range impM.classMap {
		if strings.HasPrefix(className, startOfClassName) {
			if shortestClassName == "" || len(className) < len(shortestClassName) {
				shortestClassName = className
				shortestImportPath = strings.Replace(classPath, className, "*", 1)
			} else if len(className) == len(shortestClassName) {
				importPath := strings.Replace(classPath, className, "*", 1)
				if shortestImportPath == "" || len(importPath) < len(shortestImportPath) {
					shortestClassName = className
					shortestImportPath = importPath
				}
			}
		}
	}
	return shortestClassName, shortestImportPath
}

// StarPathExact takes the exact class name and tries to return the shortest
// import path for the matching class, if found, like "java.io.*".
// Returns empty string if there are no matches.
func (impM *ImportMatcher) StarPathExact(exactClassName string) string {
	shortestClassName := ""
	shortestImportPath := ""
	for className, classPath := range impM.classMap {
		if className == exactClassName {
			if shortestClassName == "" {
				shortestClassName = className
				shortestImportPath = strings.Replace(classPath, className, "*", 1)
			} else if len(className) == len(shortestClassName) {
				importPath := strings.Replace(classPath, className, "*", 1)
				if shortestImportPath == "" || len(importPath) < len(shortestImportPath) {
					shortestClassName = className
					shortestImportPath = importPath
				}
			}
		}
	}
	return shortestImportPath
}

// ImportPathExact takes the exact class name and tries to return the shortest
// specific import path for the matching class, if found, like "java.io.File".
// Returns empty string if there are no matches.
func (impM *ImportMatcher) ImportPathExact(exactClassName string) string {
	shortestClassName := ""
	shortestImportPath := ""
	for className, classPath := range impM.classMap {
		if className == exactClassName {
			if shortestClassName == "" {
				shortestClassName = className
				shortestImportPath = classPath
			} else if len(className) == len(shortestClassName) {
				importPath := classPath
				if shortestImportPath == "" || len(importPath) < len(shortestImportPath) {
					shortestClassName = className
					shortestImportPath = importPath
				}
			}
		}
	}
	return shortestImportPath
}

// StarPathAll takes the start of the class name and tries to return all
// found class names, and also the import paths, like "java.io.*".
// Returns empty strings if there are no matches.
func (impM *ImportMatcher) StarPathAll(startOfClassName string) ([]string, []string) {
	allClassNames := make([]string, 0)
	allImportPaths := make([]string, 0)
	for className, classPath := range impM.classMap {
		if strings.HasPrefix(className, startOfClassName) {
			allClassNames = append(allClassNames, className)
			allImportPaths = append(allImportPaths, strings.Replace(classPath, className, "*", 1))
		}
	}
	return allClassNames, allImportPaths
}

// StarPathAllExact takes the exact class name and tries to return all
// matching class names, and also the import paths, like "java.io.*".
// Returns empty strings if there are no matches.
func (impM *ImportMatcher) StarPathAllExact(exactClassName string) ([]string, []string) {
	allClassNames := make([]string, 0)
	allImportPaths := make([]string, 0)
	for className, classPath := range impM.classMap {
		if className == exactClassName {
			allClassNames = append(allClassNames, className)
			allImportPaths = append(allImportPaths, strings.Replace(classPath, className, "*", 1))
		}
	}
	return allClassNames, allImportPaths
}
