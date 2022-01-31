package handler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yz89122/go-package-replacer/config"
	"github.com/yz89122/go-package-replacer/gotool/flag"
)

type Handler interface {
	Handle() error
}

func NewHandler(c *config.Config, toolCommand string, toolArgs []string) Handler {
	var toolName = filepath.Base(toolCommand)

	switch toolName {
	case "compile":
		return NewCompileHandler(c, toolCommand, toolArgs)
	case "link":
		return NewLinkHandler(c, toolCommand, toolArgs)
	}

	return NewDefaultHandler(c, toolCommand, toolArgs)
}

type baseHandler struct {
	flags *flag.Flags
}

func (h *baseHandler) runWithStdio() error {
	var cmd = exec.Command(h.flags.Cmd(), h.flags.Flags()...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run command failed: %w", err)
	}
	return nil
}
