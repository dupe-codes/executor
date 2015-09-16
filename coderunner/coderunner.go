package coderunner

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type CodeRun struct {
	Name     string `bson:"name" json:"name"`
	Code     string `bson:"code" json:"code"`
	Language string `bson:"language" json:"language"`
}

type RunResult struct {
	Output string `bson:"output" json:"output"`
	Error  error  `bson:"error" json "error"`
}

// Runs the desired code run, returning an error if one
// is encountered during code execution
func (cr *CodeRun) Run() (string, error) {
	ext, err := getExtension(cr.Language)
	if err != nil {
		return "", err
	}
	tmpFile := fmt.Sprintf("%s_run.%s", cr.Name, ext)
	dumpCodeContents(tmpFile, cr.Code)
	defer cleanTmpFiles(tmpFile)

	// Actually might not make much sense to run this as a goroutine
	// FIXME: Review this later
	resultChan := make(chan *RunResult)
	go func(lang string, codeFile string, resultChan chan *RunResult) {
		var cmd string
		var args []string
		switch lang {
		case "python":
			cmd = "python"
			args = []string{codeFile}
		}

		output, err := exec.Command(cmd, args...).Output()
		if err != nil {
			result := &RunResult{"", err}
			resultChan <- result
			return
		}

		result := &RunResult{string(output), nil}
		resultChan <- result
	}(cr.Language, tmpFile, resultChan)

	result := <-resultChan
	if result.Error != nil {
		return "", result.Error
	}

	return result.Output, nil
}

func cleanTmpFiles(tmpFile string) {
	// Finally, clean up generated tmp files
	err := os.Remove(tmpFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occurred removing the temp code file")
	}
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

// Dump code contents into tmp file for running
func dumpCodeContents(tmpFile string, code string) error {
	file, err := os.Create(tmpFile)
	if err != nil {
		return errors.New("Encountered an error writing tmp code file")
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(code)
	writer.Flush()

	return nil
}
