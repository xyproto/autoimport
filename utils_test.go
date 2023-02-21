package autoimport

import (
	"testing"
)

func TestFindJavaInEtcEnvironment(t *testing.T) {
	FindInEtcEnvironment("JAVA_HOME")
}

func TestFindKotlinInEtcEnvironment(t *testing.T) {
	FindInEtcEnvironment("KOTLIN_HOME")
}
