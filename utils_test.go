package autoimport

import (
	"testing"

	"github.com/xyproto/env/v2"
)

func TestFindJavaInEtcEnvironment(t *testing.T) {
	env.EtcEnvironment("JAVA_HOME")
}

func TestFindKotlinInEtcEnvironment(t *testing.T) {
	env.EtcEnvironment("KOTLIN_HOME")
}
