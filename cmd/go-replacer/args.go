package main

import (
	"fmt"
	"os"
)

func getArgsAndCommand() (args []string, toolCommand string, toolArgs []string) {
	var inputArgs = os.Args[1:]
	if len(inputArgs) == 0 {
		fmt.Fprintln(os.Stderr, "should be used with go build -toolexec")
		os.Exit(1)
	}

	var index, length = 0, len(inputArgs)

	for ; index < length; index++ {
		var arg = inputArgs[index]
		if arg == "--" {
			break
		}
		args = append(args, arg)
	}

	if index >= length {
		// without "--"
		toolCommand = args[0]
		if len(args) > 0 {
			toolArgs = args[1:]
		}
		args = nil
		return
	}

	// with "--"
	toolCommand = inputArgs[index+1]
	if index < length {
		toolArgs = inputArgs[index+2:]
	}
	return
}
