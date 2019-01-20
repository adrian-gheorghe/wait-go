package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"time"
)

const appVersion = "0.1.0"

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
				fmt.Println(err)
				fmt.Printf("Sleeping %d seconds waiting for host\n", intervalFlag)
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
				fmt.Printf("Sleeping %d seconds waiting for command - %s - to return\n", intervalFlag, wait)
				time.Sleep(time.Duration(intervalFlag) * time.Second)
			} else {
				fmt.Println(string(out))
				break
			}

		}
	}
}

func processCommandExec(command string, timeoutFlag int, intervalFlag int, shell string) {
	out, err := exec.Command(shell, "-c", command).Output()
	if err != nil {
		fmt.Printf("Sleeping %d seconds waiting for command - %s - to return\n", intervalFlag, command)
		time.Sleep(time.Duration(intervalFlag) * time.Second)
	} else {
		fmt.Println(string(out))
	}
}

func inPath(command string) bool {
	_, err := exec.LookPath(command)
	if err != nil {
		return false
	}
	return true
}

func chooseShell() string {
	if inPath("bash") {
		return "bash"
	} else if inPath("sh") {
		return "sh"
	} else {
		panic("Neither bash or sh present on system")
	}
}

func main() {
	var waitsFlags arrayFlags
	var commandFlags arrayFlags

	flag.Var(&waitsFlags, "wait", "You can specify the HOST and TCP PORT using the format HOST:PORT, or you can specify a command that should return an output. Multiple wait flags can be added.")
	flag.Var(&commandFlags, "command", "Command that should be run when all waits are accessible. Multiple commands can be added.")
	timeoutFlag := flag.Int("timeout", 600, "Timeout untill script is killed.")
	intervalFlag := flag.Int("interval", 15, "Interval between calls")
	version := flag.Bool("version", false, "Prints current version")
	flag.Parse()

	if *version {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	if len(waitsFlags) == 0 || len(commandFlags) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	shell := chooseShell()
	waitFor(waitsFlags, commandFlags, *timeoutFlag, *intervalFlag, shell)
}
