package main

import (
	"fmt"

	cr "github.com/njdup/executor/coderunner"
)

func main() {
	test := cr.CodeRun{
		"test_code",
		"print 'Hello this is a test'",
		"python",
	}

	err := test.Run()
	if err != nil {
		fmt.Println("failed!")
	} else {
		fmt.Println("Succeeded!")
	}
}
