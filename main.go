package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/SCU-SJL/gott/ast"
	"github.com/SCU-SJL/menuscreen"
)

func main() {

	f, ok := chooseFile()
	if !ok {
		println("no files were chosen, exit...")
		return
	}

	fInfo, err := ast.NewFileParser(f).Parse()
	if err != nil {
		log.Fatalln("ast parse failed:", err.Error())
	}

	testList := make([]string, 0, 16)
	testList = append(testList, extractTestFuncs(fInfo)...)
	testList = append(testList, extractSuiteTestMethods(fInfo)...)

	if len(testList) == 0 {
		println("no tests were found, exit...")
		return
	}

	testName, ok := chooseTest(testList)
	if !ok {
		println("no tests were chosen, exit...")
		return
	}

	args := bytes.Buffer{}
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			args.WriteString(os.Args[i])
			args.WriteString(" ")
		}
	}

	gotestCmd := fmt.Sprintf("go test %s -test.run %s", args.String(), testName)
	cmd := exec.Command("bash", "-c", gotestCmd)
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("get exec output failed: %v\n", err)
	}
	println(string(out))

}

func chooseFile() (fname string, ok bool) {
	testFiles := lsTestFiles()
	if len(testFiles) == 0 {
		return
	}
	screen := buildScreen()
	defer screen.Fini()
	_, fname, ok = screen.SetTitle("GO TEST FILES").
		SetLines(testFiles...).
		Start().
		ChosenLine()
	if ok {
		fname = "./" + fname
	}
	return
}

func chooseTest(testList []string) (tname string, ok bool) {
	screen := buildScreen()
	defer screen.Fini()
	_, v, ok := screen.SetTitle("GO TEST LIST").
		SetLines(testList...).
		Start().
		ChosenLine()
	if ok {
		v = "^" + v + "$"
	}
	return v, ok
}

func buildScreen() *menuscreen.MenuScreen {
	screen, err := menuscreen.NewMenuScreen()
	if err != nil || screen == nil {
		log.Fatalf("init screen controller failed: %v\n", err)
	}
	return screen
}

func lsTestFiles() []string {
	files := make([]string, 0, 16)
	fileInfos, err := ioutil.ReadDir("./")
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

func extractTestFuncs(f *ast.File) []string {
	fnList := make([]string, 0, len(f.FuncList))
	for _, fn := range f.FuncList {
		if fn.IsTest {
			fnList = append(fnList, fn.Name)
		}
	}
	return fnList
}

func extractSuiteTestMethods(f *ast.File) []string {

	suiteEntryMap := make(map[string]string)
	for _, fn := range f.FuncList {
		if fn.IsSuiteEntry {
			suiteEntryMap[fn.SuiteName] = fn.Name
		}
	}

	methodList := make([]string, 0, 16)
	for _, s := range f.StructList {
		entryName, ok := suiteEntryMap[s.Name]
		if !ok {
			continue
		}
		for _, m := range s.MethodList {
			if m.IsTest() {
				methodList = append(methodList, fmt.Sprintf("%s/%s", entryName, m.Name))
			}
		}
	}

	return methodList

}
