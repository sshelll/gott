package main

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sshelll/fuckflag"
	"github.com/sshelll/gott/v2/core"
	"github.com/sshelll/gott/v2/core/ui"
	"github.com/sshelll/gott/v2/core/util"
	"github.com/sshelll/sinfra/ast"
)

var (
	uiMode     bool
	goTestArgs []string
)

func main() {
	workDir := initialize()

	if versionIsSet && *versionFlag {
		fmt.Fprintf(os.Stdout, "gott@%s\nvisit '%s' for more detail.\n", version, website)
		return
	}

	goTestRegExpr := calGoTestRegRxpr()
	if goTestRegExpr == "" {
		log.Fatalln("[gott] no tests were found, exit...")
	}

	// print only
	if *printFlag {
		print(goTestRegExpr)
		return
	}

	// exec
	if err := core.ExecGoTest(workDir, goTestRegExpr, goTestArgs...); err != nil {
		log.Fatalf("[gott] exec go test failed: %v\n", err)
	}
}

func initialize() (workDir string) {
	parseFlags()

	type cond struct {
		failed bool
		errMsg string
	}

	failConds := []cond{
		{posIsSet && subIsSet, "-sub with -pos is meaningless"},
		{posIsSet && !fileIsSet, "-pos must be used with -file"},
		{subIsSet && (!printIsSet || !fileIsSet), "-sub must be used with -print and -file"},
		{posIsSet && util.StrToInt(*posFlag) < 0, "-pos must be a positive integer"},
		{fileIsSet && !strings.HasSuffix(*fileFlag, "_test.go"), "-file must be a go test file"},
	}

	for _, cond := range failConds {
		if cond.failed {
			log.Fatalf("[gott] invalid flag, %s\n", cond.errMsg)
		}
	}

	if fileIsSet {
		f, err := os.Stat(*fileFlag)
		if f == nil || err != nil {
			log.Fatalf("[gott] invalid file path '%s'\n", *fileFlag)
		}
	}

	workDir = "."
	uiMode = !versionIsSet && !posIsSet && !fileIsSet && !subIsSet
	if uiMode {
		return
	}

	if !fileIsSet {
		return
	}
	return filepath.Dir(*fileFlag)
}

func parseFlags() {
	// init fuckflag
	fuckflag.Usage = usage
	fuckflag.CommandLine.SetOutput(os.Stdout)

	// parse flags
	fuckflag.Parse()
	goTestArgs = fuckflag.Extends()

	// init set status
	posIsSet = fuckflag.IsSet(posFlagName)
	subIsSet = fuckflag.IsSet(subFlagName)
	fileIsSet = fuckflag.IsSet(fileFlagName)
	printIsSet = fuckflag.IsSet(printFlagName)
	versionIsSet = fuckflag.IsSet(versionFlagName)
}

func calGoTestRegRxpr() string {
	if uiMode {
		return ui.NewUI().Run()
	}

	f, err := ast.Parse(*fileFlag)
	if err != nil {
		log.Fatalf("[gott] parse file as ast failed, err = %v\n", err)
	}

	// file mode
	if !posIsSet && !subIsSet {
		testFuncs := core.ExtractTestFuncs(f)
		return util.BuildGoTestRegExpr(testFuncs...)
	}

	// pos mode
	if posIsSet {
		tf, found := core.FindClosestTestFunc(f, util.StrToInt(*posFlag))
		if !found {
			return ""
		}
		return util.BuildGoTestRegExpr(tf)
	}

	// sub mode
	if subIsSet && *subFlag {
		testFuncs := core.ExtractTestFuncs(f)
		testifyMethods := core.ExtractTestifySuiteTestMethods(f)
		return util.BuildGoTestRegExpr(append(testFuncs, testifyMethods...)...)
	}

	return ""
}
