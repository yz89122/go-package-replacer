package config

type Config struct {
	Replacements []*Replacement
}

type Replacement struct {
	// default see: github.com/yz89122/go-package-replacer/utils.ImportPathToAlias
	PackageAlias     string `json:",omitempty"`
	OriginImportPath string
	NewImportPath    string
	// The dir that contains the go.mod file of the new package.
	// The path should be an absolute path.
	NewPackageRoot string

	// TODO: add replacing strategies for different situations
}
