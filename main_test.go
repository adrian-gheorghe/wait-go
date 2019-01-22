package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

type testPathDetector struct {
	ShellReturn string
}

func (pathDetector testPathDetector) inPath(command string) bool {
	if command == pathDetector.ShellReturn {
		return true
	}
	return false
}

func readVersion() string {
	dat, _ := ioutil.ReadFile("./VERSION")
	return string(dat)
}

func TestChooseShell(t *testing.T) {

	cases := []struct {
		Name     string
		Shell    string
		Expected string
	}{
		{"Test if system has bash in path", "bash", "bash"},
		{"Test if system has sh in path", "sh", "sh"},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := tc.Expected
			testPathDetector := testPathDetector{}
			testPathDetector.ShellReturn = tc.Shell
			if actual != chooseShell(testPathDetector) {
				t.Fatal("failure")
			}
		})
	}
}

func TestMainExecution(t *testing.T) {
	var pathDetector = localPathDetector{}
	var testWaitsFlags arrayFlags
	var testCommandFlags arrayFlags
	var testTimeoutFlag = 10
	var testIntervalFlag = 5
	var testVersion = true

	var buf bytes.Buffer
	log.SetOutput(&buf)
	mainExecution(testWaitsFlags, testCommandFlags, testTimeoutFlag, testIntervalFlag, testVersion, pathDetector)
	log.SetOutput(os.Stderr)
	out := buf.String()

	t.Fatal("aaa" + out)
}
