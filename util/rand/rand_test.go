package rand

import (
	"fmt"
	"testing"
)

func TestGenerateString(t *testing.T) {
	generateString := GenerateString("a")
	fmt.Println(SplitStringToName(generateString))
}
