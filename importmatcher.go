// Package importmatcher tries to find which import should be used, given the start of a class name
package importmatcher

import (
	"archive/zip"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"github.com/xyproto/env"
)

// ImportMatcher is a struct that contains a list of JAR file paths,
// and a lookup map from class names to class paths, which is populated
// when New or NewCustom is called.
type ImportMatcher struct {
	JARPaths []string          // list of paths to examine for .jar files
	classMap map[string]string // map from class name to class path
}

func New(kotlinAsWell bool) (*ImportMatcher, error) {
	var JARPaths = []string{
		env.Str("JAVA_HOME", "/usr/lib/jvm/default"),
	}
	if kotlinAsWell {
		JARPaths = append(JARPaths, "/usr/share/kotlin/lib")
	}
	return NewCustom(JARPaths)
}

func NewCustom(JARPaths []string) (*ImportMatcher, error) {
	var impM ImportMatcher
	impM.JARPaths = make([]string, len(JARPaths))
	for i := range JARPaths {
		JARPath := JARPaths[i]
		if !strings.HasSuffix(JARPath, "/") {
			JARPath += "/"
		}
		impM.JARPaths[i] = JARPath
	}
	allClasses, err := impM.findClasses()
	if err != nil {
		return nil, err
	}
	classMap := make(map[string]string)
	for _, classPath := range allClasses {
		className := classPath
		if strings.Contains(classPath, ".") {
			fields := strings.Split(classPath, ".")
			lastField := fields[len(fields)-1]
			className = lastField
		}
		classMap[className] = classPath
	}
	impM.classMap = classMap
	return &impM, nil
}

// ClassMap returns the mapping from class names to class paths
func (impM *ImportMatcher) ClassMap() map[string]string {
	return impM.classMap
}

// readJAR returns a list of classes within the given .jar file,
// for instance "SomePath/SomePath/SomeClass"
func readJAR(filePath string) ([]string, error) {
	readCloser, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	foundClasses := make([]string, 0)

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

			foundClasses = append(foundClasses, className)
		}
	}

	return foundClasses, nil
}

func (impM *ImportMatcher) FindClassesInJAR(JARPath string) ([]string, error) {
	allClasses := make([]string, 0)

	var wg sync.WaitGroup
	var add sync.Mutex

	return allClasses, filepath.Walk(JARPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		filePath := path
		fileName := info.Name()

		if filepath.Ext(fileName) != ".jar" && filepath.Ext(fileName) != ".JAR" {
			return nil
		}

		//packageName := strings.TrimSuffix(strings.TrimSuffix(fileName, ".jar"), ".JAR")

		wg.Add(1)
		go func() {
			defer wg.Done()
			foundClasses, err := readJAR(filePath)
			if err != nil {
				log.Printf("error: %s\n", err)
			}
			add.Lock()
			allClasses = append(allClasses, foundClasses...)
			add.Unlock()
		}()
		wg.Wait()

		return nil
	})
}

func (impM *ImportMatcher) findClasses() ([]string, error) {
	var wg sync.WaitGroup
	var addMut sync.Mutex
	allClasses := make([]string, 0)
	for _, JARPath := range impM.JARPaths {
		wg.Add(1)
		go func() {
			defer wg.Done()
			foundClasses, err := impM.FindClassesInJAR(JARPath)
			if err != nil {
				log.Printf("error: %s\n", err)
			}
			addMut.Lock()
			allClasses = append(allClasses, foundClasses...)
			addMut.Unlock()

		}()
		wg.Wait()
	}
	return allClasses, nil
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

// StarPathAll takes the start of the class name and tries to return all
// found class names, and also the import paths, like "java.io.*".
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