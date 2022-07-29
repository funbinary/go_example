package bshell

import (
	"bytes"
	"github.com/funbinary/go_example/pkg/bfile"
	"runtime"
	"strings"
)

// ShellExec executes given command `cmd` synchronously and returns the command result.
func ShellExec(cmd string, environment ...[]string) (result string, err error) {
	var (
		buf = bytes.NewBuffer(nil)
		p   = NewProcess(
			getShell(),
			append([]string{getShellOption()}, parseCommand(cmd)...),
			environment...,
		)
	)
	p.Stdout = buf
	p.Stderr = buf
	err = p.Run()
	result = buf.String()
	return
}

// SearchBinary searches the binary `file` in current working folder and PATH environment.
func SearchBinary(file string) string {
	// Check if it is absolute path of exists at current working directory.
	if bfile.Exists(file) {
		return file
	}
	return SearchBinaryPath(file)
}
func SearchBinaryPath(file string) string {
	array := ([]string)(nil)
	switch runtime.GOOS {
	case "windows":
		envPath := GetEnv("PATH")
		if strings.Contains(envPath, ";") {
			array = strings.Split(envPath, ";")
		} else if strings.Contains(envPath, ":") {
			array = strings.Split(envPath, ":")
		}
		if bfile.Ext(file) != ".exe" {
			file += ".exe"
		}

	default:
		array = strings.Split(GetEnv("PATH"), ":")
	}
	if len(array) > 0 {
		path := ""
		for _, v := range array {
			path = v + bfile.Separator + file
			if bfile.Exists(path) && bfile.IsFile(path) {
				return path
			}
		}
	}
	return ""
}

func getShellOption() string {
	switch runtime.GOOS {
	case "windows":
		return "/c"
	default:
		return "-c"
	}
}

func parseCommand(cmd string) (args []string) {
	if runtime.GOOS != "windows" {
		return []string{cmd}
	}
	// Just for "cmd.exe" in windows.
	var argStr string
	var firstChar, prevChar, lastChar1, lastChar2 byte
	array := strings.Split(cmd, " ")
	for _, v := range array {
		if len(argStr) > 0 {
			argStr += " "
		}
		firstChar = v[0]
		lastChar1 = v[len(v)-1]
		lastChar2 = 0
		if len(v) > 1 {
			lastChar2 = v[len(v)-2]
		}
		if prevChar == 0 && (firstChar == '"' || firstChar == '\'') {
			// It should remove the first quote char.
			argStr += v[1:]
			prevChar = firstChar
		} else if prevChar != 0 && lastChar2 != '\\' && lastChar1 == prevChar {
			// It should remove the last quote char.
			argStr += v[:len(v)-1]
			args = append(args, argStr)
			argStr = ""
			prevChar = 0
		} else if len(argStr) > 0 {
			argStr += v
		} else {
			args = append(args, v)
		}
	}
	return
}

func getShell() string {
	switch runtime.GOOS {
	case "windows":
		return SearchBinary("cmd.exe")
	default:
		// Check the default binary storage path.
		if bfile.Exists("/bin/bash") {
			return "/bin/bash"
		}
		if bfile.Exists("/bin/sh") {
			return "/bin/sh"
		}
		// Else search the env PATH.
		path := SearchBinary("bash")
		if path == "" {
			path = SearchBinary("sh")
		}
		return path
	}
}
