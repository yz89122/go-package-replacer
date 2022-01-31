package handler

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yz89122/go-package-replacer/config"
	"github.com/yz89122/go-package-replacer/gofile"
	"github.com/yz89122/go-package-replacer/gotool/flag"
	"github.com/yz89122/go-package-replacer/importcfg"
)

type CompileHandler struct {
	*baseHandler
	config *config.Config
	flags  *flag.CompileFlags
}

const (
	tmpGoSourceCodeDir = "code"
	magicTmpMainDir    = "go-package-replacer-magic-main"
	mainFilename       = "main.go"
)

var (
	//go:embed main.go.template
	mainFileTemplate string
	// The first imported package should be b002.
	newPackageArchivePath = filepath.Join("b002", "_pkg_.a")
)

func NewCompileHandler(config *config.Config, toolCommand string, toolArgs []string) *CompileHandler {
	var flags = flag.NewCompileFlags(toolCommand, toolArgs)
	return &CompileHandler{
		baseHandler: &baseHandler{flags: flags.Flags},
		config:      config,
		flags:       flags,
	}
}

func (h *CompileHandler) Handle() error {
	if h.config == nil || len(h.config.Replacements) == 0 || h.flags.IsCheckVersion() {
		return h.runWithStdio()
	}

	if replacementConfig := h.getReplacementConfig(); replacementConfig != nil {
		return h.buildNewPackage(replacementConfig) // return here, not to build the original package
	}

	// This method will update the arguments (h.flags)
	// that is going to be passed to the underlying command (h.runWithStdio()).
	if err := h.updateImportPaths(); err != nil {
		return fmt.Errorf("failed to update import paths of Go source files: %w", err)
	}

	return h.runWithStdio()
}

func (h *CompileHandler) getReplacementConfig() *config.Replacement {
	var importPath = h.flags.PackageImportPath()

	for _, replacement := range h.config.Replacements {
		if replacement.OriginImportPath == importPath {
			return replacement
		}
	}

	return nil
}

func (h *CompileHandler) updateImportPaths() error {
	var importCfgPath = h.flags.ImportCfg()
	var packageBuildDir = filepath.Dir(importCfgPath)

	for _, sourceFilepath := range h.flags.GoSourceFiles() {
		var filename = filepath.Base(sourceFilepath)

		// source code of the Go file
		var sourceCode string
		{
			var bytesData, err = os.ReadFile(sourceFilepath)
			if err != nil {
				return fmt.Errorf("failed to read file [%s]: %w", sourceFilepath, err)
			}
			sourceCode = string(bytesData)
		}

		// create Go file object and update the file content
		var goFile = gofile.NewFile(sourceCode)
		for _, replacement := range h.config.Replacements {
			goFile.ReplacePackage(
				replacement.PackageAlias,
				replacement.OriginImportPath,
				replacement.NewImportPath,
			)
		}

		// only replace argument if the file is updated
		if !goFile.Updated() {
			continue
		}

		// write to temporary directory
		var tmpDir = filepath.Join(packageBuildDir, tmpGoSourceCodeDir)
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return fmt.Errorf("failed to mkdir all [%s]: %w", tmpDir, err)
		}
		var newFilepath = filepath.Join(tmpDir, filename)
		if err := os.WriteFile(newFilepath, []byte(goFile.String()), 0644); err != nil {
			return fmt.Errorf("failed to write file [%s]: %w", newFilepath, err)
		}

		// replace argument
		h.flags.ReplaceGoSourceFile(sourceFilepath, newFilepath)
	}

	// read import cfg
	var importCfg = importcfg.NewImportCfg()
	{
		if data, err := os.ReadFile(importCfgPath); err != nil {
			return fmt.Errorf("failed to read importcfg [%s]: %w", importCfgPath, err)
		} else if err := importCfg.ParseString(string(data)); err != nil {
			return fmt.Errorf("failed to parse importcfg [%s]: %w", importCfgPath, err)
		}
	}

	// update importcfg
	var updated bool
	for _, replacement := range h.config.Replacements {
		if u := importCfg.RenamePackagefile(replacement.OriginImportPath, replacement.NewImportPath); u {
			updated = true
		}
	}

	// write it back
	if updated {
		if err := os.WriteFile(importCfgPath, []byte(importCfg.String()), 0644); err != nil {
			return fmt.Errorf("failed to write updated importcfg [%s]: %w", importCfgPath, err)
		}
	}

	return nil
}

func (h *CompileHandler) buildNewPackage(replacementConfig *config.Replacement) error {
	var packageRoot = replacementConfig.NewPackageRoot

	var mainFilepath, mainFileDir string
	{
		// create temporary main from template
		var buffer bytes.Buffer
		if err := template.Must(template.New("").Parse(mainFileTemplate)).
			Execute(&buffer, map[string]interface{}{
				"importPath": replacementConfig.NewImportPath,
			}); err != nil {
			return fmt.Errorf("failed to render temporary main file: %w", err)
		}

		// create directory for temporary main.go
		mainFileDir = filepath.Join(packageRoot, magicTmpMainDir)
		if err := os.MkdirAll(mainFileDir, 0755); err != nil {
			return fmt.Errorf("failed to create dir [%s]: %w", mainFileDir, err)
		}

		// write temporary main.go to the file system
		mainFilepath = filepath.Join(mainFileDir, mainFilename)
		if err := os.WriteFile(mainFilepath, buffer.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write file [%s]: %w", mainFilepath, err)
		}
	}

	// go build command for build the new package
	var command = exec.Command("go", "build", "-a", "-work", "-x", mainFilepath)
	{
		command.Dir = filepath.Dir(mainFilepath)
	}

	// build the new package and extract temporary directory from stdout
	var newPackageTmpBuildDir string
	{
		var output, err = command.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to build package (%s): %w", replacementConfig.NewImportPath, err)
		}
		var outputStr = string(output)
		var lines = strings.Split(outputStr, "\n")
		outputStr = lines[0]
		newPackageTmpBuildDir = strings.TrimPrefix(strings.TrimSpace(outputStr), "WORK=")
	}

	// cleanup temporary main.go
	if err := os.RemoveAll(mainFileDir); err != nil {
		return fmt.Errorf("failed to cleanup: %w", err)
	}

	// record temporary build directory of the new package
	var packageBuildDir = filepath.Dir(h.flags.ImportCfg())
	if err := os.WriteFile(
		filepath.Join(packageBuildDir, newPackageWorkDirHolder),
		[]byte(newPackageTmpBuildDir),
		0644,
	); err != nil {
		return fmt.Errorf("failed to record temporary build directory for new package: %w", err)
	}

	// copy package archive to output path
	{
		var archivePath = filepath.Join(newPackageTmpBuildDir, newPackageArchivePath)
		var outputPath = h.flags.OutputPath()
		if err := func() error {
			var source *os.File
			{
				var err error
				source, err = os.Open(archivePath)
				if err != nil {
					return fmt.Errorf("failed to open [%s]: %w", archivePath, err)
				}
				defer func() { _ = source.Close() }()
			}

			var destination *os.File
			{
				var err error
				destination, err = os.Create(outputPath)
				if err != nil {
					return fmt.Errorf("failed to create file [%s]: %w", outputPath, err)
				}
				defer func() { _ = destination.Close() }()
			}

			if _, err := io.Copy(destination, source); err != nil {
				return fmt.Errorf("failed to copy: %w", err)
			}
			return nil
		}(); err != nil {
			return fmt.Errorf("failed to copy package archive: %w", err)
		}
	}

	return nil
}
