package main

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sshelll/fuckflag"
	"github.com/sshelll/gott/v2/core"
	"github.com/sshelll/gott/v2/core/ui"
	"github.com/sshelll/gott/v2/core/util"
	"github.com/sshelll/sinfra/ast"
)

var (
	uiMode     bool
	goTestFile string
	goTestArgs []string
)

func main() {
	initialize()

	if *versionFlag {
		fmt.Fprintf(os.Stderr, "gott@%s\nvisit '%s' for more detail.\n", version, website)
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
	if err := core.ExecGoTest(goTestRegExpr, goTestArgs...); err != nil {
		log.Fatalf("[gott] exec go test failed: %v\n", err)
	}
}

func initialize() {
	fuckflag.Usage = usage
	fuckflag.Parse()
	goTestArgs = fuckflag.Extends()

	failConds := []bool{
		*posFlag != "" && *subFlag,                   // -sub with -pos is meaningless
		*posFlag != "" && *fileFlag == "",            // -pos must be used with -file
		*subFlag && (!*printFlag || *fileFlag == ""), // -sub must be used with -print and -file
	}

	for _, failed := range failConds {
		if failed {
			log.Fatalln("[gott] invalid flag combination, please use gott -h to get help.")
		}
	}

	uiMode = !*versionFlag && *posFlag == "" && *fileFlag == "" && !*subFlag
	if uiMode {
		return
	}

	file := *fileFlag
	if file == "" {
		return
	}
	if !strings.HasSuffix(file, "_test.go") {
		log.Fatalf("[gott] '%s' is not a go test file\n", file)
	}
	if f, err := os.Stat(*fileFlag); f == nil || err != nil {
		log.Fatalf("[gott] invalid file path '%s'\n", *fileFlag)
	}
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
	if *posFlag == "" && !*subFlag {
		testFuncs := core.ExtractTestFuncs(f)
		return util.BuildGoTestRegExpr(testFuncs...)
	}

	// pos mode
	if *posFlag != "" {
		tf, found := core.FindClosestTestFunc(f, util.StrToInt(*posFlag))
		if !found {
			return ""
		}
		return util.BuildGoTestRegExpr(tf)
	}

	// sub mode
	if *subFlag {
		testFuncs := core.ExtractTestFuncs(f)
		testifyMethods := core.ExtractTestifySuiteTestMethods(f)
		return util.BuildGoTestRegExpr(append(testFuncs, testifyMethods...)...)
	}

	return ""
}
