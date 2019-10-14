package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// ERRORCODE Default error code
const ERRORCODE = 111

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(ERRORCODE)
}

// GetFileEnv reads and parses variable value from file
// (only the first line)
// From "man envdir":
// "Spaces and tabs at the end of t are removed. Nulls in t
//       are changed to newlines in the environment variable."
func GetFileEnv(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	value := scanner.Text()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	// trim the trailing newlines, tabs and spaces
	value = strings.TrimRight(value, " \t\n")
	if value == "" {
		return "", nil
	}
	// envdir replaces null bytes with newlines
	return strings.Replace(value, "\x00", "\n", -1), nil
}

// ReadEnvDir reads environment variables from directory files
func ReadEnvDir(d string) ([]string, error) {
	f, err := os.Open(d)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fInfos, err1 := f.Readdir(0)
	if err1 != nil {
		return nil, err1
	}
	envs := make([]string, len(fInfos))
	for i, fInfo := range fInfos {
		path := filepath.Join(d, fInfo.Name())
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		value, err := GetFileEnv(f)
		if err != nil {
			return nil, err
		}
		envs[i] = fmt.Sprintf("%s=%s", fInfo.Name(), value)
	}
	return envs, nil
}

// BuildEnv enrich the original environment with updated values
// for existing variables, and remove variables with empty values
func BuildEnv(oldEnv []string, newEnv []string) []string {
	m := make(map[string]string)
	for _, env := range oldEnv {
		nameVal := strings.SplitN(env, "=", 2)
		m[nameVal[0]] = nameVal[1]
	}
	for _, env := range newEnv {
		nameVal := strings.SplitN(env, "=", 2)
		if nameVal[1] != "" {
			m[nameVal[0]] = nameVal[1]
		} else {
			delete(m, nameVal[0])
		}
	}
	resultEnv := make([]string, len(m))
	for name, val := range m {
		resultEnv = append(resultEnv, fmt.Sprintf("%s=%s", name, val))
	}
	return resultEnv
}

// RunCmd runs a child process, passes std* handles transparently,
// in the case of a failure exits from the main process
// with the child's exitcode
func RunCmd(envs []string, name string, args []string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = BuildEnv(os.Environ(), envs)
	if err := cmd.Start(); err != nil {
		fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			if status, ok := ee.Sys().(*syscall.WaitStatus); ok {
				exitcode := status.ExitStatus()
				os.Exit(exitcode)
			}
		}
		fatal(err)
	}
}

func main() {
	if len(os.Args) < 3 {
		fatal("Usage: goenvdir d command [args]")
	}
	envs, err := ReadEnvDir(os.Args[1])
	if err != nil {
		fatal("failed to parse dir: ", err)
	}
	RunCmd(envs, os.Args[2], os.Args[3:])
}
