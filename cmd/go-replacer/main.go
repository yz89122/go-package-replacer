package main

import (
	"fmt"
	"os"

	"github.com/yz89122/go-package-replacer/handler"
)

func main() {
	var args, toolCommand, toolArgs = getArgsAndCommand()
	// var toolName = filepath.Base(toolCommand)

	// fmt.Fprintln(os.Stderr, "args:", args)
	// fmt.Fprintln(os.Stderr, "tool name:", toolName)
	// fmt.Fprintln(os.Stderr, "tool command:", toolCommand)
	// fmt.Fprintln(os.Stderr, "tool args:", toolArgs)

	var h handler.Handler
	{
		var config, err = getConfigFromArgs(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to get config:", err)
			os.Exit(1)
		}

		h = handler.NewHandler(
			config,
			toolCommand,
			toolArgs,
		)
	}

	if err := h.Handle(); err != nil {
		fmt.Fprintln(os.Stderr, "err:", err)
		os.Exit(1)
	}
}
