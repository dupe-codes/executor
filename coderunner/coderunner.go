package coderunner

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type CodeRun struct {
	Name     string `bson:"name" json:"name"`
	Code     string `bson:"code" json:"code"`
	Language string `bson:"language" json:"language"`
}

// Runs the desired code run, returning an error if one
// is encountered during code execution
func (cr *CodeRun) Run() error {
	ext, err := getExtension(cr.Language)
	if err != nil {
		return err
	}
	tmpFile := fmt.Sprintf("%s_run.%s", cr.Name, ext)

	// Dump code contents into tmp file for running
	file, err := os.Create(tmpFile)
	if err != nil {
		return errors.New("Encountered an error writing tmp code file")
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(cr.Code)
	writer.Flush()

	return nil
}

// Returns the appropriate extension type for the given
// programming language
func getExtension(lang string) (string, error) {
	switch lang {
	case "python":
		return "py", nil
	}
	return "", errors.New("No matching extension for given language found")
}
