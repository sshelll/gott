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
	"github.com/sshelll/gott/core"
	"github.com/sshelll/sinfra/ast"
)

const (
	version = "v1.5.0"
	website = "https://github.com/sshelll/gott"
)

var (
	versionFlag = fuckflag.Bool("version", false, "print version of gott")

	posFlag     = fuckflag.String("pos", "", "the uri with absolute filepath to exec the closest test\n\n"+
		"[EXAMPLE 1]:\n\t'gott --pos=/Users/sshelll/go/src/xx_test.go:104'\n"+
		"\tthis will exec the closest test to the 104(byte pos) in the file xx_test.go\n\n"+
		"[EXAMPLE 2]:\n\t'gott --pos=/Users/sshelll/go/src/xx_test.go:235 -v -race'\n"+
		"\tthis will exec the closest test to the 235(byte pos) in the file xx_test.go with -v -race flags\n",
	)
	
	runFileFlag = fuckflag.String("runFile", "", "the uri with absolute filepath to exec all test in the file\n"+
		"[EXAMPLE 1]:\n\t'gott --runFile=/Users/sshelll/go/src/xx_test.go'\n"+
		"\tthis will exec all tests in 'xx_test.go'\n"+
		"[EXAMPLE 2]:\n\t'gott --runFile=/Users/sshelll/go/src/xx_test.go -v -race'\n"+
		"\tthis will exec all tests in 'xx_test.go' with -v -race flags\n",
	)

	pFlag = fuckflag.Bool("p", false, "print the test name instead of exec it\n\n"+
		"[EXAMPLE 1]:\n\t'gott -p --pos=/Users/sshelll/go/src/xx_test.go:104'\n"+
		"\tthis will print the test name instead of exec it\n\n"+
		"[EXAMPLE 2]:\n\t'gott -p --runFile=/Users/sshelll/go/src/xx_test.go'\n"+
		"\tthis will print name of all tests in the file instead of exec them\n",
	)
)

func main() {
	fuckflag.Usage = usage
	fuckflag.Parse()
	prechck()

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

func prechck() {
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

func usage() {
	fmt.Fprintf(os.Stderr, "Gott is a alternative to 'go test' command, it can help you to choose a specific test to run with UI.\n")
	fmt.Fprintf(os.Stderr, "Also, it has some useful features, such as:\n \t-find the closest test func by byte pos\n\t-find all tests in a go test file\n\n")
	fmt.Fprintf(os.Stderr, "If you want to use it as a command line tool, just try 'gott [go test args]' in your terminal, for example:\n")
	fmt.Fprintf(os.Stderr, "\tgott -v ==> go test -v\n")
	fmt.Fprintf(os.Stderr, "\tgott -race ==> go test -race\n")
	fmt.Fprintf(os.Stderr, "\tgott -v -race ==> go test -v -race\n\n")
	fmt.Fprintf(os.Stderr, "If you want to use it as a binary tool to help get a test name, just try 'gott -p [other flags]' in your terminal, for example:\n")
	fmt.Fprintf(os.Stderr, "\tgott -p --pos=/Users/sshelll/go/src/xx_test.go:235 ==> print the closest test to the 235(byte pos) in the file xx_test.go\n")
	fmt.Fprintf(os.Stderr, "\tgott -p --runFile=/Users/sshelll/go/src/xx_test.go ==> print all tests in the file xx_test.go\n\n")
	fmt.Fprintf(os.Stderr, "TIPS: the result of 'gott -p' is a regexp test name, you can use it with 'go test' or 'dlv' command, for example:\n")
	fmt.Fprintf(os.Stderr, "\tgo test -v -test.run $(gott -p --pos=/Users/sshelll/go/src/xx_test.go:235)\n")
	fmt.Fprintf(os.Stderr, "\tdlv test --build-flags=-test.run $(gott -p --pos=/Users/sshelll/go/src/xx_test.go:235)\n\n")
	fmt.Fprintf(os.Stderr, "WARN: The work dir of gott should be a go package dir, otherwise it will not work.\n")
	fmt.Fprintf(os.Stderr, "Also, although gott is a alternative command to 'go test', 'go' is still required.\n\n")
	fmt.Fprintf(os.Stderr, "For more detail please visit '%s'. Have fun with gott!\n\n", website)
	fmt.Fprintf(os.Stderr, "Usage of Gott:\n")
	fuckflag.PrintDefaults()
}
