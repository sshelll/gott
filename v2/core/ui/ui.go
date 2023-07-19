package ui

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/sshelll/gott/v2/core"
	"github.com/sshelll/gott/v2/core/util"
	"github.com/sshelll/menuscreen"
	"github.com/sshelll/sinfra/ast"
)

const (
	TestAllOption = "→TEST ALL!"
)

type UI struct {
	workDir string
}

func NewUI() *UI {
	return &UI{
		workDir: "./",
	}
}

func (ui *UI) Run() (testName string) {
	f, ok := ui.ChooseTestFile()
	if !ok {
		log.Fatalln("[gott] no files were chosen, exit...")
	}

	fInfo, err := ast.Parse(f)
	if err != nil {
		log.Fatalln("[gott] ast parse failed:", err.Error())
	}

	goTests, testifyTests := core.ExtractTestFuncs(fInfo), core.ExtractTestifySuiteTestMethods(fInfo)
	testList := append(goTests, testifyTests...)

	if len(testList) == 0 {
		log.Fatalln("[gott] no tests were found in this file, exit...")
	}

	testName, testAll, ok := ui.ChooseTest(testList)
	if !ok {
		log.Fatalln("[gott] no tests were chosen, exit...")
	}

	if testAll {
		testName = util.BuildGoTestRegExpr(goTests...)
	}

	return
}

func (ui *UI) ChooseTestFile() (file string, ok bool) {
	testFiles := ui.lsTestFiles()
	if len(testFiles) == 0 {
		return
	}
	screen := ui.buildScreen()
	defer screen.Fini()
	_, file, ok = screen.SetTitle("GO TEST FILES").
		SetLines(testFiles...).
		Start().
		ChosenLine()
	if ok {
		file = ui.workDir + file
	}
	return
}

func (ui *UI) ChooseTest(testList []string) (tname string, testAll, ok bool) {
	screen := ui.buildScreen()
	defer screen.Fini()
	_, tname, ok = screen.SetTitle("GO TEST LIST").
		SetLines(testList...).
		AppendLines(TestAllOption).
		Start().
		ChosenLine()
	if tname == TestAllOption {
		return tname, true, true
	}
	if ok {
		tname = "^" + tname + "$"
	}
	return
}

func (ui *UI) lsTestFiles() []string {
	files := make([]string, 0, 16)
	fileInfos, err := ioutil.ReadDir(ui.workDir)
	if err != nil {
		log.Fatalf("read current dir failed: %v", err)
	}
	for _, f := range fileInfos {
		if !f.IsDir() && strings.HasSuffix(f.Name(), "_test.go") {
			files = append(files, f.Name())
		}
	}
	return files
}

func (ui *UI) buildScreen() *menuscreen.MenuScreen {
	screen, err := menuscreen.NewMenuScreen()
	if err != nil || screen == nil {
		log.Fatalf("init screen controller failed: %v\n", err)
	}
	return screen
}
