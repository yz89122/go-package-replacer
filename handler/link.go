package handler

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/yz89122/go-package-replacer/config"
	"github.com/yz89122/go-package-replacer/gotool/flag"
	"github.com/yz89122/go-package-replacer/importcfg"
)

type LinkHandler struct {
	*baseHandler
	config *config.Config
	flags  *flag.LinkFlags
}

var (
	packageLinkImportCfgPath = filepath.Join("b001", "importcfg.link")
)

func NewLinkHandler(config *config.Config, toolCommand string, toolArgs []string) *LinkHandler {
	var flags = flag.NewLinkFlags(toolCommand, toolArgs)
	return &LinkHandler{
		baseHandler: &baseHandler{flags: flags.Flags},
		config:      config,
		flags:       flags,
	}
}

func (h *LinkHandler) Handle() error {
	if h.config == nil || len(h.config.Replacements) == 0 || h.flags.IsCheckVersion() {
		return h.runWithStdio()
	}

	var importCfgPath = h.flags.ImportCfg()

	var newPackageTmpDirHolders []string
	{
		var mainPackageBuildDir = filepath.Dir(importCfgPath)
		var tmpBuildDir = filepath.Dir(mainPackageBuildDir)
		var tmpBuildDirFS = os.DirFS(tmpBuildDir)
		var err error
		newPackageTmpDirHolders, err = fs.Glob(tmpBuildDirFS, filepath.Join("*", newPackageWorkDirHolder))
		if err != nil {
			return fmt.Errorf("failed to glob: %w", err)
		}
		for index, holder := range newPackageTmpDirHolders {
			newPackageTmpDirHolders[index] = filepath.Join(tmpBuildDir, holder) // to abs path
		}
	}

	var newPackageTmpDirs []string
	{
		for _, holder := range newPackageTmpDirHolders {
			var tmpDir string
			if data, err := os.ReadFile(holder); err != nil {
				return fmt.Errorf("failed to read [%s]: %w", holder, err)
			} else {
				tmpDir = string(data)
			}
			newPackageTmpDirs = append(newPackageTmpDirs, tmpDir)
		}
	}

	// update importcfg
	{
		var importCfg = importcfg.NewImportCfg()

		// importcfg of replaced packages
		for _, tmpDir := range newPackageTmpDirs {
			if data, err := os.ReadFile(filepath.Join(tmpDir, packageLinkImportCfgPath)); err != nil {
				return fmt.Errorf("failed to read importcfg (%s): %w", tmpDir, err)
			} else if err := importCfg.ParseString(string(data)); err != nil {
				return fmt.Errorf("failed to parese importcfg: %w", err)
			}
		}

		// origin importcfg
		if data, err := os.ReadFile(importCfgPath); err != nil {
			return fmt.Errorf("failed to read importcfg (%s): %w", importCfgPath, err)
		} else if err := importCfg.ParseString(string(data)); err != nil {
			return fmt.Errorf("failed to parese importcfg (%s): %w", importCfgPath, err)
		}

		// write updated importcfg
		if err := os.WriteFile(importCfgPath, []byte(importCfg.String()), 0644); err != nil {
			return fmt.Errorf("failed to update importcfg: %w", err)
		}
	}

	if err := h.runWithStdio(); err != nil {
		return err
	}

	for _, tmpDir := range newPackageTmpDirs {
		if err := os.RemoveAll(tmpDir); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "failed to remove temporary directory", tmpDir, "err:", err)
		}
	}

	return nil
}
