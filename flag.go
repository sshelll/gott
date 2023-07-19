package main

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"fmt"
	"os"

	"github.com/sshelll/fuckflag"
)

const (
	version = "v1.4.2"
	website = "https://github.com/sshelll/gott"
)

var (
	versionFlag = fuckflag.Bool("version", false, "print version of gott")

	posFlag = fuckflag.String("pos", "", "the uri with absolute filepath to exec the closest test\n\n"+
		"[EXAMPLE 1]:\n\t'gott --pos=/Users/sshelll/go/src/xx_test.go:104'\n"+
		"\tthis will exec the closest test to the 104(byte pos) in the file xx_test.go\n\n"+
		"[EXAMPLE 2]:\n\t'gott --pos=/Users/sshelll/go/src/xx_test.go:235 -v -race'\n"+
		"\tthis will exec the closest test to the 235(byte pos) in the file xx_test.go with -v -race flags\n",
	)

	runFileFlag = fuckflag.String("runFile", "", "the uri with absolute filepath to exec all test in the file\n\n"+
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

	listAllFlag = fuckflag.Bool("listAll", false, "print all tests of a file, this flag shoule be used with '-runFile'\n\n"+
		"NOTE: the difference between '-listAll' and '-p' is that '-p' will not print 'testify test methods'\n"+
		"For Example, if you got a test file like this:\n"+
		`
// package test
//
// import (
// 	"testing"
// 	"github.com/stretchr/testify/suite"
// )
// 
// func TestXxx(t *testing.T) {
// 	println("hello world")
// }
// 
// type FooTestSuite struct {
// 	suite.Suite
// }
// 
// func TestFoo(t *testing.T) {
// 	testifySuite.Run(t, &FooTestSuite{})
// }
// 
// func (*FooTestSuite) TestCase() {
// 
// }

`+
		"The -p flag will only print 'TestXXX' and 'TestFoo' for you,\n"+
		"while -listAll will print both of them and 'TestFoo/TestCase'.\n",
	)
)

func usage() {
	fmt.Fprintf(os.Stderr, "Gott is a alternative to 'go test' command, it can help you to choose a specific test to run\nwith UI.\n\n")
	fmt.Fprintf(os.Stderr, "Also, it has some useful features, such as:\n \t-find the closest test func by byte pos\n\t-find all tests in a go test file\n\n")
	fmt.Fprintf(os.Stderr, "If you want to use it as a command line tool, just try 'gott [go test args]' in your terminal,\n for example:\n")
	fmt.Fprintf(os.Stderr, "\tgott -v ==> go test -v\n")
	fmt.Fprintf(os.Stderr, "\tgott -race ==> go test -race\n")
	fmt.Fprintf(os.Stderr, "\tgott -v -race ==> go test -v -race\n\n")
	fmt.Fprintf(os.Stderr, "If you want to use it as a binary tool to help get a test name, just try \n'gott -p [other flags]' in your terminal, for example:\n")
	fmt.Fprintf(os.Stderr, "\tgott -p --pos=/Users/sshelll/go/src/xx_test.go:235 ==> print the closest\n\ttest to the 235(byte pos) in the file xx_test.go\n")
	fmt.Fprintf(os.Stderr, "\tgott -p --runFile=/Users/sshelll/go/src/xx_test.go ==> print all tests\n\tin the file xx_test.go\n\n")
	fmt.Fprintf(os.Stderr, "TIPS: the result of 'gott -p' is a regexp test name, you can use it with \n'go test' or 'dlv' command, for example:\n")
	fmt.Fprintf(os.Stderr, "\tgo test -v -test.run $(gott -p --pos=/Users/sshelll/go/src/xx_test.go:235)\n")
	fmt.Fprintf(os.Stderr, "\tdlv test --build-flags=-test.run $(gott -p --pos=/Users/sshelll/go/src/xx_test.go:235)\n\n")
	fmt.Fprintf(os.Stderr, "WARN: The work dir of gott should be a go package dir, otherwise it will not work.\n")
	fmt.Fprintf(os.Stderr, "Also, although gott is a alternative command to 'go test', 'go' is still required.\n\n")
	fmt.Fprintf(os.Stderr, "For more detail please visit '%s'. Have fun with gott!\n\n", website)
	fmt.Fprintf(os.Stderr, "----------------------------------------------------------------------------------------------\n\n")
	fmt.Fprintf(os.Stderr, "Usage of Gott:\n")
	fuckflag.PrintDefaults()
}
