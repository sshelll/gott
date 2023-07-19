package main

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"fmt"
	"log"
	"strings"

	"github.com/sshelll/fuckflag"
	"github.com/sshelll/gott/core"
	"github.com/sshelll/sinfra/ast"
)

func main() {
	initialize()

	var (
		testName string
		done     bool
	)

	switch true {
	// version
	case *versionFlag:
		versionMode()
		return

	// pos
	case *posFlag != "":
		testName, done = posMode()

	// run file
	case *runFileFlag != "":
		testName, done = runFileMode()

	// interactive
	default:
		testName, done = interactiveMode()
	}

	if done {
		return
	}

	if *pFlag {
		print(testName)
		return
	}

	if len(strings.TrimSpace(testName)) > 0 {
		core.ExecGoTest(testName, fuckflag.Extends()...)
	}
}

func initialize() {
	fuckflag.Usage = usage
	fuckflag.Parse()
	if *runFileFlag != "" && *posFlag != "" {
		log.Fatalln("[gott] -pos and -runFile flag cannot co-exist.")
	}
}

func versionMode() {
	println(fmt.Sprintf("gott version %s\nto get more detail please visit '%s'.", version, website))
}

func interactiveMode() (testName string, done bool) {
	f, ok := core.ChooseTestFile()
	if !ok {
		log.Println("[gott] no files were chosen, exit...")
		return "", true
	}

	fInfo, err := ast.Parse(f)
	if err != nil {
		log.Fatalln("[gott] ast parse failed:", err.Error())
	}

	goTests, testifyTests := core.ExtractTestFuncs(fInfo), core.ExtractTestifySuiteTestMethods(fInfo)
	testList := append(goTests, testifyTests...)

	if len(testList) == 0 {
		log.Println("[gott] no tests were found, exit...")
		return "", true
	}

	testName, testAll, ok := core.ChooseTest(testList)
	if !ok {
		log.Println("[gott] no tests were chosen, exit...")
		return "", true
	}

	if testAll {
		testName = core.BuildTestAllExpr(goTests)
	}

	return
}

func posMode() (testName string, done bool) {
	uri := *posFlag

	f, pos, err := core.ParseURI(uri)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fInfo, err := ast.Parse(f)
	if err != nil {
		log.Fatalln("[gott] ast parse failed:", err.Error())
	}

	testName, ok := core.FindClosestTestFunc(fInfo, pos)
	if !ok {
		log.Println("[gott] no tests were found, exit...")
		return "", true
	}

	testName = fmt.Sprintf("^%s$", testName)

	return
}

func runFileMode() (testName string, done bool) {
	f := *runFileFlag
	fInfo, err := ast.Parse(f)
	if err != nil {
		log.Fatalln("[gott] ast parse failed:", err.Error())
	}

	goTests := core.ExtractTestFuncs(fInfo)

	if len(goTests) == 0 {
		log.Println("[gott] no tests were found, exit...")
		return "", true
	}

	testName = core.BuildTestAllExpr(goTests)
	return
}
