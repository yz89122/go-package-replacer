package utils

import "strings"

func ImportPathToAlias(importPath string) string {
	if len(importPath) == 0 {
		return ""
	}

	var packageName string
	{
		var splitted = strings.Split(importPath, "/")
		packageName = splitted[len(splitted)-1]
	}

	var alias string
	{
		var splitted = strings.Split(packageName, "-")
		alias = splitted[len(splitted)-1]
	}

	return alias
}
