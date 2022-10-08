package util

import (
	"fmt"
	"strings"

	"github.com/SCU-SJL/sinfra/ast"
)

const (
	testingPkgName = `testing`
	testifyPkgName = `github.com/stretchr/testify/suite`
)

func ExtractTestFuncs(f *ast.File) []string {
	fnList := make([]string, 0, len(f.FuncList))
	testingPkg := findTestingPkgName(f.ImportList)
	for _, fn := range f.FuncList {
		if ast.IsGoTestFunc(fn, &testingPkg) {
			fnList = append(fnList, fn.Name)
		}
	}
	return fnList
}

func ExtractTestifySuiteTestMethods(f *ast.File) []string {

	testingPkg := findTestingPkgName(f.ImportList)
	testifyPkg := findTestifyPkgName(f.ImportList)

	suiteEntryMap := make(map[string]string)
	for _, fn := range f.FuncList {
		suiteName, ok := ast.IsTestifySuiteEntryFunc(fn, &testingPkg, &testifyPkg)
		if ok {
			suiteEntryMap[suiteName] = fn.Name
		}
	}

	methodList := make([]string, 0, 16)
	for _, s := range f.StructList {
		entryName, ok := suiteEntryMap[s.Name]
		if !ok {
			continue
		}
		for _, m := range s.MethodList {
			if strings.HasPrefix(m.Name, "Test") {
				methodList = append(methodList, fmt.Sprintf("%s/%s", entryName, m.Name))
			}
		}
	}

	return methodList

}

func findTestifyPkgName(importList []*ast.Import) string {
	alias := findPkgAlias(importList, testifyPkgName)
	if alias != nil {
		return *alias
	}
	return "suite"
}

func findTestingPkgName(importList []*ast.Import) string {
	alias := findPkgAlias(importList, testingPkgName)
	if alias != nil {
		return *alias
	}
	return "testing"
}

func findPkgAlias(importList []*ast.Import, pkg string) (alias *string) {
	for _, imp := range importList {
		if imp.Pkg == pkg {
			if imp.Alias == "" {
				return nil
			}
			return &imp.Alias
		}
	}
	return nil
}
