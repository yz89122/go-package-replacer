package gofile_test

import (
	testing "testing"

	"github.com/yz89122/go-package-replacer/gofile"
)

func TestReplace(t *testing.T) {
	type testcase struct {
		alias         string
		source        string
		expect        string
		origin        string
		newImportPath string
	}
	for index, testcase := range []testcase{
		{
			source:        `package test`,
			expect:        `package test`,
			origin:        "test",
			newImportPath: "replaced/test",
		},
		{
			source: `package test
			import "abc"`,
			expect: `package test
			import abc "replaced/abc"`,
			origin:        "abc",
			newImportPath: "replaced/abc",
		},
		{
			source: `package test
			import "abc"
			func init() {}`,
			expect: `package test
			import "abc"
			func init() {}`,
			origin:        "ddd",
			newImportPath: "replaced/ddd",
		},
		{
			source: `package test
			import aaaa "abc"
			func init() {}`,
			expect: `package test
			import aaaa "replaced/abc"
			func init() {}`,
			origin:        "abc",
			newImportPath: "replaced/abc",
		},
		{
			source: `package test
			import aaaa"abc"
			func init() {}`,
			expect: `package test
			import aaaa"replaced/abc"
			func init() {}`,
			origin:        "abc",
			newImportPath: "replaced/abc",
		},
		{
			source: `package test
			import (
				aaaa"abc"
			)
			func init() {}`,
			expect: `package test
			import (
				aaaa"replaced/abc"
			)
			func init() {}`,
			origin:        "abc",
			newImportPath: "replaced/abc",
		},
		{
			source: `package test
			import (
				"dd"
				aaaa"abc"
			)
			func init() {}`,
			expect: `package test
			import (
				"dd"
				aaaa"replaced/abc"
			)
			func init() {}`,
			origin:        "abc",
			newImportPath: "replaced/abc",
		},
		{
			source: `package test
			import (
				"dd"
				aaaa "abc"
			)
			func init() {}`,
			expect: `package test
			import (
				"dd"
				aaaa "replaced/abc"
			)
			func init() {}`,
			origin:        "abc",
			newImportPath: "replaced/abc",
		},
		{
			source: `package test
			import (
				"dd"
				aaaa                          "abc"
			)
			func init() {}`,
			expect: `package test
			import (
				"dd"
				aaaa                          "replaced/abc"
			)
			func init() {}`,
			origin:        "abc",
			newImportPath: "replaced/abc",
		},
	} {
		var goFile = gofile.NewFile(testcase.source)
		goFile.ReplacePackage(testcase.alias, testcase.origin, testcase.newImportPath)
		if output := goFile.String(); output != testcase.expect {
			t.Error(
				"\nindex:", index,
				"\nsource:\n", testcase.source,
				"\norigin:", testcase.origin,
				"\nnew import path:", testcase.newImportPath,
				"\nexpect:\n", testcase.expect,
				"\ngot:\n", output,
			)
		}
	}
}
