package importmatcher

import (
	"errors"
	"path/filepath"

	"github.com/xyproto/env"
)

const archJavaPath = "/usr/lib/jvm/default"
const debianJavaPath = "/usr/lib/jvm/default-java"

// FindJava finds the most likely location of a Java installation
// (with subfolders with .jar files) on the system.
func FindJava() (string, error) {
	// 1. Respect $JAVA_HOME, if it's set
	if javaHomePath := env.Str("JAVA_HOME"); javaHomePath != "" && isDir(javaHomePath) {
		if isDir(javaHomePath) {
			return javaHomePath, nil
		}
	}
	// 2. Find out if "java" is in the $PATH
	if javaExecutablePath := which("java"); javaExecutablePath != "" {
		// Follow the symlink as far as it takes us
		for isSymlink(javaExecutablePath) {
			javaExecutablePath = followSymlink(javaExecutablePath)
		}
		// Return the grandparent directory of the java executable (since it's typically in the "bin" directory")
		if grandParentDirectory := filepath.Dir(filepath.Dir(javaExecutablePath)); isDir(grandParentDirectory) {
			return grandParentDirectory, nil
		}
	}
	// 3. Consider typical paths, for Arch Linux and Debian/Ubuntu
	if isDir(archJavaPath) {
		return archJavaPath, nil
	} else if isDir(debianJavaPath) {
		return debianJavaPath, nil
	}
	return "", errors.New("could not find an installation of Java")
}
