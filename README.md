# Go Package Replacer

## Config

See [`./config/config.go`](./config/config.go)

### Example

```json
// config.json
{
	"replacements": [
		{
			"PackageAlias": "pkg",
			"OriginImportPath": "github.com/origin/pkg",
			"NewImportPath": "github.com/replaced/the-new-package",
			"NewPackageRoot": "/the/path/to/the/source/code/located"
		}
	]
}
```

## Run

```bash
# build this tool
go install 'github.com/yz89122/go-package-replacer/cmd/go-replacer@latest'
# uses this tool when building
go build -a -toolexec "$(go env GOPATH)/bin/go-replacer /abs/path/to/the/config.json --" ./main.go
```

## Known issues

 - The new package cannot cause import cycle. This will cause the generated package archive incompatible.
   - **NOTE:** The error might occurred in any form.
   - If anyone known how to fix this, please open a PR.
