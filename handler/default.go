package handler

import (
	"github.com/yz89122/go-package-replacer/config"
	"github.com/yz89122/go-package-replacer/gotool/flag"
)

type DefaultHandler struct {
	*baseHandler
	config *config.Config
	flags  *flag.Flags
}

func NewDefaultHandler(config *config.Config, toolCommand string, toolArgs []string) *DefaultHandler {
	var flags = flag.NewFlags(toolCommand, toolArgs)
	return &DefaultHandler{
		baseHandler: &baseHandler{flags: flags},
		config:      config,
		flags:       flags,
	}
}

func (h *DefaultHandler) Handle() error {
	return h.runWithStdio()
}
