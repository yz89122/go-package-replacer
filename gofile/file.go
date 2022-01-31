package gofile

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"

	"github.com/yz89122/go-package-replacer/utils"
)

type File struct {
	sourceCode string
	lines      []string
	updated    bool
}

func NewFile(sourceCode string) *File {
	return &File{
		sourceCode: sourceCode,
		lines:      strings.Split(sourceCode, "\n"),
	}
}

func (f *File) Updated() bool {
	return f.updated
}

func (f *File) String() string {
	if f.updated {
		return strings.Join(f.lines, "\n")
	}

	return f.sourceCode
}

func (f *File) lastLineOfImport() (end int) {
	var fset = token.NewFileSet()
	var astFile, err = parser.ParseFile(fset, "", f.sourceCode, parser.ImportsOnly)
	if err != nil {
		return
	}

	for _, decl := range astFile.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok != token.IMPORT {
				continue
			}
			if line := fset.Position(d.TokPos).Line; line > end {
				end = line
			}
			if line := fset.Position(d.Rparen).Line; line > end {
				end = line
			}
		}
	}

	return
}

func (f *File) ReplacePackage(alias, origin, newImportPath string) {
	// uses string replacement instead of AST in order to
	// keep the original styling

	var end = f.lastLineOfImport()

	var matcher = regexp.MustCompilePOSIX(
		fmt.Sprintf(
			`(import[ \t]+)?(([a-zA-Z0-9_]+)[ \t]*)?"%s"`,
			origin,
		),
	)

	for index, line := range f.lines {
		if !matcher.MatchString(line) {
			continue
		} else if index >= end {
			// in case of changing unintended lines
			return
		}

		f.updated = true

		var groups = matcher.FindStringSubmatch(line)
		var actualAlias = groups[3]
		var replacement = fmt.Sprintf(`$1$2"%s"`, newImportPath)
		if len(actualAlias) == 0 {
			if len(alias) > 0 {
				actualAlias = alias
			} else {
				actualAlias = utils.ImportPathToAlias(origin)
			}
			replacement = fmt.Sprintf(`${1}%s "%s"`, actualAlias, newImportPath)
		}

		f.lines[index] = matcher.ReplaceAllString(line, replacement)

		return
	}
}
