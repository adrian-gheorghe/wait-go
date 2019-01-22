package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var appVersionFile = "./VERSION"

func readVersion(appVersionFile string) string {
	dat, _ := ioutil.ReadFile(appVersionFile)
	return string(dat)
}

// PathDetector detects if binaries are in path
type PathDetector interface {
	inPath(command string) bool
}
type localPathDetector struct {
}

func (pathDetector localPathDetector) inPath(command string) bool {
	_, err := exec.LookPath(command)
	if err != nil {
		return false
	}
	return true
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "String"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func waitFor(waitsFlags arrayFlags, commandFlags arrayFlags, timeoutFlag int, intervalFlag int, shell string) {
	for _, wait := range waitsFlags {
		processWait(wait, timeoutFlag, intervalFlag, shell)
	}

	for _, command := range commandFlags {
		processCommandExec(command, timeoutFlag, intervalFlag, shell)
	}

}

func processWait(wait string, timeoutFlag int, intervalFlag int, shell string) {
	pattern, _ := regexp.Compile("(.*):(.*)")
	matches := pattern.FindAllStringSubmatch(wait, -1)
	if len(matches) > 0 {
		dbHost := matches[0][1]
		dbPort := matches[0][2]

		for {
			conn, err := net.DialTimeout("tcp", net.JoinHostPort(dbHost, dbPort), time.Duration(intervalFlag)*time.Second)
			if err != nil {
				log.Println(err)
				log.Printf("Sleeping %d seconds waiting for host\n", intervalFlag)
				time.Sleep(time.Duration(intervalFlag) * time.Second)
			}
			if conn != nil {
				conn.Close()
				break
			}
		}
	} else {
		for {
			out, err := exec.Command(shell, "-c", wait).Output()
			if err != nil {
				log.Printf("Sleeping %d seconds waiting for command - %s - to return\n", intervalFlag, wait)
				time.Sleep(time.Duration(intervalFlag) * time.Second)
			} else {
				log.Println(string(out))
				break
			}

		}
	}
}

func processCommandExec(command string, timeoutFlag int, intervalFlag int, shell string) {
	out, err := exec.Command(shell, "-c", command).Output()
	if err != nil {
		log.Printf("Sleeping %d seconds waiting for command - %s - to return\n", intervalFlag, command)
		time.Sleep(time.Duration(intervalFlag) * time.Second)
	} else {
		log.Println(string(out))
	}
}

func chooseShell(pathDetector PathDetector) string {
	if pathDetector.inPath("bash") {
		return "bash"
	} else if pathDetector.inPath("sh") {
		return "sh"
	} else {
		panic("Neither bash or sh present on system")
	}
}

func mainExecution(waitsFlags arrayFlags, commandFlags arrayFlags, timeoutFlag int, intervalFlag int, version bool, pathDetector localPathDetector) int {
	if version {
		log.Println(readVersion(appVersionFile))
		return 0
	}

	if len(waitsFlags) == 0 || len(commandFlags) == 0 {
		log.Println("You must specify at least a wait and a command. Please see --help for more information.")
		return 1
	}
	shell := chooseShell(pathDetector)
	waitFor(waitsFlags, commandFlags, timeoutFlag, intervalFlag, shell)
	return 2
}

func main() {
	var pathDetector = localPathDetector{}
	// Set custom logger
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	var waitsFlags arrayFlags
	var commandFlags arrayFlags

	flag.Var(&waitsFlags, "wait", "You can specify the HOST and TCP PORT using the format HOST:PORT, or you can specify a command that should return an output. Multiple wait flags can be added.")
	flag.Var(&commandFlags, "command", "Command that should be run when all waits are accessible. Multiple commands can be added.")
	timeoutFlag := flag.Int("timeout", 600, "Timeout untill script is killed.")
	intervalFlag := flag.Int("interval", 15, "Interval between calls")
	version := flag.Bool("version", false, "Prints current version")
	flag.Parse()
	returnValue := mainExecution(waitsFlags, commandFlags, *timeoutFlag, *intervalFlag, *version, pathDetector)
	os.Exit(returnValue)
}
