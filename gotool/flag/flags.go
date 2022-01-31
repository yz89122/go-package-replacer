package flag

import (
	"strings"
)

type Flags struct {
	cmd   string
	flags []string
}

func NewFlags(cmd string, flags []string) *Flags {
	var nb = &Flags{
		cmd:   cmd,
		flags: make([]string, len(flags)),
	}
	copy(nb.flags, flags)
	return nb
}

func (f *Flags) Cmd() string {
	return f.cmd
}

func (f *Flags) Flags() []string {
	var flags = make([]string, len(f.flags))
	copy(flags, f.flags)
	return flags
}

func (f *Flags) IsCheckVersion() bool {
	return f.HasFlag("V")
}

func (f *Flags) OutputPath() string {
	return f.getFlagValue("o")
}

func (f *Flags) ImportCfg() string {
	return f.getFlagValue("importcfg")
}

func (f *Flags) HasFlag(targetFlag string) bool {
	targetFlag = "-" + targetFlag

	for _, flag := range f.flags {
		if !strings.HasPrefix(flag, targetFlag) {
			continue
		}

		var inlineFlag = targetFlag + "="        // -flag=
		if strings.HasPrefix(flag, inlineFlag) { // -flag=value
			return true
		}

		if flag != targetFlag {
			continue
		}

		return true
	}

	return false
}

func (f *Flags) getFlagValue(targetFlag string) string {
	targetFlag = "-" + targetFlag // -flag

	for index, flag := range f.flags {
		if !strings.HasPrefix(flag, targetFlag) {
			continue
		}

		var inlineFlag = targetFlag + "="        // -flag=
		if strings.HasPrefix(flag, inlineFlag) { // -flag=value
			return strings.TrimPrefix(flag, inlineFlag)
		}

		if flag != targetFlag {
			continue
		}

		if index+1 < len(f.flags) { // -flag value
			return f.flags[index+1]
		}

		return ""
	}

	return ""
}

func (f *Flags) replaceFlagValue(targetFlag string, value string) {
	targetFlag = "-" + targetFlag // -flag

	for index, flag := range f.flags {
		if !strings.HasPrefix(flag, targetFlag) {
			continue
		}

		var inlineFlag = targetFlag + "="        // -flag=
		if strings.HasPrefix(inlineFlag, flag) { // -flag=value
			f.flags[index] = inlineFlag + value
			return
		}

		if index+1 < len(f.flags) { // -flag value
			f.flags[index+1] = value
			return
		}

		return
	}
}
