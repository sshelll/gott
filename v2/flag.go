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
	version = "v2.0.0"
	website = "https://github.com/sshelll/gott"
)

var (
	versionFlag = fuckflag.Bool("version", false, "print version of gott\n")

	posFlag = fuckflag.String("pos", "", "byte pos of a test file, this flag should be used with '-file'\n\n"+
		"[EXAMPLE 1]:\n\t'gott -pos=104 -file=/Users/sshelll/go/src/xx_test.go'\n"+
		"\tthis will exec the closest test to the 104(byte pos) in the file xx_test.go\n\n"+
		"[EXAMPLE 2]:\n\t'gott -pos=235 -file=/Users/sshelll/go/src/xx_test.go -v -race'\n"+
		"\tthis will exec the closest test to the 235(byte pos) in the file xx_test.go\nwith -v -race flags\n",
	)

	fileFlag = fuckflag.String("file", "", "filepath of a go test file\n\n"+
		"NOTE: if the path is not start with '/', then the final file path would be './file'\n\n"+
		"[EXAMPLE 1]:\n\t'gott -file=/Users/sshelll/go/src/xx_test.go'\n"+
		"\tthis will exec all tests in the file\n"+
		"[EXAMPLE 2]:\n\t'gott -file=xx_test.go -v -race'\n"+
		"\tthis will exec all tests in './xx_test.go' with -v -race flags\n",
	)

	printFlag = fuckflag.Bool("print", false, "print the result instead of exec\n\n"+
		"[EXAMPLE 1]:\n\t'gott -print -pos=104 -file=xx_test.go'\n"+
		"\tthis will print the closest test name instead of exec it\n\n"+
		"[EXAMPLE 2]:\n\t'gott -print -file=xx_test.go'\n"+
		"\tthis will print name of all tests in the file instead of exec them\n",
	)

	subFlag = fuckflag.Bool("sub", false, "get all tests of a file(including sub-tests), this flag must be used\n"+
		"with '-print' and '-file'\n\n"+
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
		"'gott -print -file=xx_test.go' will only print 'TestXXX' and 'TestFoo' for you,\n"+
		"'gott -print -file=xx_test.go -sub' will print both of them and 'TestFoo/TestCase'.\n",
	)
)

func usage() {
	fmt.Fprintf(os.Stderr, "Gott is a alternative to 'go test' command, it can help you to choose a specific test to run\nwith UI.\n\n")
	fmt.Fprintf(os.Stderr, "Also, it has some useful features, such as:\n \t-find the closest test func by byte pos\n\t-find all tests in a go test file\n\n")
	fmt.Fprintf(os.Stderr, "If you want to use it as a interactive test runner with UI, just try 'gott [go test args]'!\nFor example:\n")
	fmt.Fprintf(os.Stderr, "\tgott -v ==> go test -v\n")
	fmt.Fprintf(os.Stderr, "\tgott -race ==> go test -race\n")
	fmt.Fprintf(os.Stderr, "\tgott -v -race ==> go test -v -race\n\n")
	fmt.Fprintf(os.Stderr, "If you want to use it to find / exec tests without UI, I mean if you wish to treat it\nlike 'github.com/josharian/impl', please check the Usage below.\n\n")
	fmt.Fprintf(os.Stderr, "For more detail please visit '%s'. Have fun with gott!\n\n", website)
	fmt.Fprintf(os.Stderr, "----------------------------------------------------------------------------------------------\n\n")
	fmt.Fprintf(os.Stderr, "Usage of Gott:\n")
	fuckflag.PrintDefaults()
}
