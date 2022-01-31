package importcfg

import (
	"bytes"
	"fmt"
	"strings"
)

// ImportCfg stores the content of the `-importcfg` file.
// The format of the file is referencing
// https://github.com/golang/go/blob/go1.17.6/src/cmd/compile/internal/base/flag.go#L392-L423
// and
// https://github.com/golang/go/blob/go1.17.6/src/cmd/link/internal/ld/ld.go#L55-L89
type ImportCfg struct {
	content map[string]map[string]string
}

const (
	VerbImportmap    = "importmap"
	VerbPackagefile  = "packagefile"
	VerbPackageshlib = "packageshlib"
	VerbModinfo      = "modinfo"

	// https://github.com/golang/go/blob/go1.17.6/src/cmd/go/internal/work/exec.go#L761
	ImportCfgBegin = "# import config"
)

func NewImportCfg() *ImportCfg {
	return &ImportCfg{
		content: make(map[string]map[string]string),
	}
}

func (c *ImportCfg) addConfig(verb, before, after string) {
	var verbConfigs = c.content[verb]
	if verbConfigs == nil {
		verbConfigs = make(map[string]string)
		c.content[verb] = verbConfigs
	}

	verbConfigs[before] = after
}

func (c *ImportCfg) RenamePackagefile(origin, newImportPath string) bool {
	var value string
	{
		var ok bool
		value, ok = c.content[VerbPackagefile][origin]
		if !ok {
			return false
		}
	}

	c.content[VerbPackagefile][newImportPath] = value
	delete(c.content[VerbPackagefile], origin)
	return true
}

func (c *ImportCfg) ParseString(data string) error {
	// https://github.com/golang/go/blob/go1.17.6/src/cmd/compile/internal/base/flag.go#L392-L423
	// https://github.com/golang/go/blob/go1.17.6/src/cmd/link/internal/ld/ld.go#L55-L89
	for lineNum, line := range strings.Split(data, "\n") {
		lineNum++ // 1-based
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		var verb, args string
		if i := strings.Index(line, " "); i < 0 {
			verb = line
		} else {
			verb, args = line[:i], strings.TrimSpace(line[i+1:])
		}
		var before, after string
		if i := strings.Index(args, "="); i >= 0 {
			before, after = args[:i], args[i+1:]
		} else {
			before = args
		}
		switch verb {
		default:
			return fmt.Errorf("%d: unknown directive %q", lineNum, verb)
		case VerbImportmap,
			VerbPackagefile,
			VerbPackageshlib,
			VerbModinfo:
			c.addConfig(verb, before, after)
		}
	}
	return nil
}

func (c *ImportCfg) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(ImportCfgBegin)
	buffer.WriteByte('\n')

	for _, verb := range []string{VerbImportmap, VerbPackagefile, VerbPackageshlib, VerbModinfo} {
		for before, after := range c.content[verb] {
			buffer.WriteString(verb)
			buffer.WriteByte(' ')

			switch verb {
			case VerbModinfo:
				// https://github.com/golang/go/blob/a5c0b190809436fd196a348f85eca0416f4de7fe/src/cmd/link/internal/ld/ld.go#L88-L93
				buffer.WriteString(before)
			default:
				buffer.WriteString(before)
				buffer.WriteByte('=')
				buffer.WriteString(after)
			}

			buffer.WriteByte('\n')
		}
	}

	return buffer.String()
}
