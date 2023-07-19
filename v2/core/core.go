package core

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/sshelll/sinfra/ast"
	sutil "github.com/sshelll/sinfra/util"
)

const (
	testingPkgName = `testing`
	testifyPkgName = `github.com/stretchr/testify/suite`
)

func ExecGoTest(workDir string, testName string, testArgs ...string) error {
	// change work dir
	if len(strings.TrimSpace(workDir)) == 0 {
		workDir = "."
	}
	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("change work dir to '%s' failed: %v", workDir, err)
	}

	// exec go test
	goTestCmd := fmt.Sprintf("go test %s -test.run %s",
		strings.Join(testArgs, " "), testName)
	if err := sutil.ExecBashCmd(goTestCmd); err != nil {
		return fmt.Errorf("command: '%s', err: '%v'", goTestCmd, err)
	}
	return nil
}

// FindClosestTestFunc find the closest test function from a file, includes go test funcs and testify test methods.
func FindClosestTestFunc(f *ast.File, targetPos int) (fn string, ok bool) {
	goTestFunc, goTestDistance, goTestFound := ExtractClosestTestFunc(f, targetPos)
	testifyTestFunc, testifyTestDistance, testifyTestFound := ExtractClosestTestifySuiteTestMethod(f, targetPos)

	if !goTestFound && !testifyTestFound {
		return "", false
	}

	if !goTestFound {
		return testifyTestFunc, true
	}

	if !testifyTestFound {
		return goTestFunc, true
	}

	if goTestDistance < testifyTestDistance {
		return goTestFunc, true
	}

	return testifyTestFunc, true
}

// ExtractTestFuncs extract all test functions from a file.
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

// ExtractClosestTestFunc extract the closest test function from a file.
func ExtractClosestTestFunc(f *ast.File, targetPos int) (testName string, distance int, found bool) {
	goTestFuncs := make([]*ast.Func, 0, len(f.FuncList))
	testingPkg := findTestingPkgName(f.ImportList)

	for _, fn := range f.FuncList {
		if ast.IsGoTestFunc(fn, &testingPkg) {
			goTestFuncs = append(goTestFuncs, fn)
		}
	}

	distance = math.MaxInt
	for _, fn := range goTestFuncs {
		pos := fn.AstDecl.Pos()
		end := fn.AstDecl.End()
		if int(pos) <= targetPos && targetPos <= int(end) {
			return fn.Name, 0, true
		}
		dis := math.Min(math.Abs(float64(int(pos)-targetPos)), math.Abs(float64(int(end)-targetPos)))
		if testName == "" || int(dis) < distance {
			found = true
			testName = fn.Name
			distance = int(dis)
		}
	}

	return
}

// ExtractTestifySuiteTestMethods extract all testify suite test methods from a file.
func ExtractTestifySuiteTestMethods(f *ast.File) []string {
	suiteEntryMap := extractTestifySuiteEntryMap(f)

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

// ExtractClosestTestifySuiteTestMethod extract the closest testify suite test method from a file.
func ExtractClosestTestifySuiteTestMethod(f *ast.File, targetPos int) (testName string, distance int, found bool) {
	suiteEntryMap := extractTestifySuiteEntryMap(f)

	distance = math.MaxInt
	for _, s := range f.StructList {
		entryName, ok := suiteEntryMap[s.Name]
		if !ok {
			continue
		}
		for _, m := range s.MethodList {
			if !strings.HasPrefix(m.Name, "Test") {
				continue
			}
			pos := m.AstDecl.Pos()
			end := m.AstDecl.End()
			if int(pos) <= targetPos && targetPos <= int(end) {
				return fmt.Sprintf("%s/%s", entryName, m.Name), 0, true
			}
			dis := math.Min(math.Abs(float64(int(pos)-targetPos)), math.Abs(float64(int(end)-targetPos)))
			if testName == "" || int(dis) < distance {
				found = true
				testName = fmt.Sprintf("%s/%s", entryName, m.Name)
				distance = int(dis)
			}
		}
	}

	return
}

// extractTestifySuiteEntryMap extracts testify suite struct name to the test entry func name map.
func extractTestifySuiteEntryMap(f *ast.File) map[string]string {
	testingPkg := findTestingPkgName(f.ImportList)
	testifyPkg := findTestifyPkgName(f.ImportList)

	suiteEntryMap := make(map[string]string)
	for _, fn := range f.FuncList {
		suiteName, ok := ast.IsTestifySuiteEntryFunc(fn, &testingPkg, &testifyPkg)
		if ok {
			suiteEntryMap[suiteName] = fn.Name
		}
	}
	return suiteEntryMap
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
