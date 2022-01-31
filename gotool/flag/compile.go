package flag

import "strings"

type CompileFlags struct {
	*Flags
}

func NewCompileFlags(cmd string, flags []string) *CompileFlags {
	return &CompileFlags{Flags: NewFlags(cmd, flags)}
}

func (f *CompileFlags) PackageImportPath() string {
	return f.getFlagValue("p")
}

func (f *CompileFlags) GoLanguageVersion() string {
	return f.getFlagValue("lang")
}

func (f *CompileFlags) SetPackageImportPath(value string) {
	f.replaceFlagValue("p", value)
}

func (f *CompileFlags) GoSourceFiles() []string {
	for index, flag := range f.flags {
		if !strings.HasSuffix(flag, ".go") {
			continue
		}

		var files = make([]string, len(f.flags[index:]))
		copy(files, f.flags[index:])
		return files
	}

	return nil
}

func (f *CompileFlags) SetGoSourceFiles(files []string) {
	for index, flag := range f.flags {
		if !strings.HasSuffix(flag, ".go") {
			continue
		}

		f.flags = append(f.flags[:index], files...)
		break
	}
}

func (f *CompileFlags) ReplaceGoSourceFile(origin, newFilePath string) {
	for index, flag := range f.flags {
		if flag == origin {
			f.flags[index] = newFilePath
		}
	}
}
