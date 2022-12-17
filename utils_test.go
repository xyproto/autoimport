package importmatcher

import (
	"fmt"
	"testing"
)

func TestFindJavaInEtcEnvironment(t *testing.T) {
	javaHome, err := FindInEtcEnvironment("JAVA_HOME")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Found $JAVA_HOME in /etc/environment: %s\n", javaHome)
}

func TestFindKotlinInEtcEnvironment(t *testing.T) {
	kotlinHome, err := FindInEtcEnvironment("KOTLIN_HOME")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Found $KOTLIN_HOME in /etc/environment: %s\n", kotlinHome)
}
