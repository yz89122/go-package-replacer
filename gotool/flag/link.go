package flag

type LinkFlags struct {
	*Flags
}

func NewLinkFlags(cmd string, flags []string) *LinkFlags {
	return &LinkFlags{Flags: NewFlags(cmd, flags)}
}

func (f *LinkFlags) ImportCfg() string {
	return f.getFlagValue("importcfg")
}
