// Hook Echo is a simply utility used for testing the Webhook package.

package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"io"
)

func checkPrefix(prefixMap map[string]struct{}, prefix string, arg string) bool {
	if _, found := prefixMap[prefix]; found {
		fmt.Printf("prefix specified more then once: %s", arg)
		os.Exit(-1)
	}

	if strings.HasPrefix(arg, prefix) {
		prefixMap[prefix] = struct{}{}
		return true
	}

	return false
}

func main() {
	var outputStream io.Writer
	outputStream = os.Stdout
	seenPrefixes := make(map[string]struct{})
	exit_code := 0

	for _, arg := range os.Args[1:] {
		if checkPrefix(seenPrefixes, "stream=", arg) {
			switch arg {
			case "stream=stdout":
				outputStream = os.Stdout
			case "stream=stderr":
				outputStream = os.Stderr
			case "stream=both":
				outputStream = io.MultiWriter(os.Stdout, os.Stderr)
			default:
				fmt.Printf("unrecognized stream specification: %s", arg)
				os.Exit(-1)
			}
		} else if checkPrefix(seenPrefixes, "exit=", arg) {
			exit_code_str := arg[5:]
			var err error
			exit_code_conv, err := strconv.Atoi(exit_code_str)
			exit_code = exit_code_conv
			if err != nil {
				fmt.Printf("Exit code %s not an int!", exit_code_str)
				os.Exit(-1)
			}
		}
	}

	if len(os.Args) > 1 {
		fmt.Fprintf(outputStream, "arg: %s\n", strings.Join(os.Args[1:], " "))
	}

	var env []string
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "HOOK_") {
			env = append(env, v)
		}
	}

	if len(env) > 0 {
		fmt.Fprintf(outputStream, "env: %s\n", strings.Join(env, " "))
	}

	os.Exit(exit_code)
}
