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
		println("exit...")
		return
	}

	fInfo, err := ast.NewFileParser(f).Parse()
	if err != nil {
		log.Fatalln("ast parse failed:", err.Error())
	}

	testList := make([]string, 0, 16)
	testList = append(testList, extractTestFuncs(fInfo)...)
	for _, s := range extractTestSuites(fInfo) {
		testList = append(testList, extractSuiteTestMethods(s)...)
	}

	testName, ok := chooseTest(testList)
	if !ok {
		println("exit...")
		return
	}

	args := bytes.Buffer{}
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			args.WriteString(os.Args[i])
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
	_, v, ok := buildScreen().
		SetTitle("GO TEST FILES").
		SetLines(lsTestFiles()...).
		Start().
		ChosenLine()
	if ok {
		v = "./" + v
	}
	return v, ok
}

func chooseTest(testList []string) (tname string, ok bool) {
	_, v, ok := buildScreen().
		SetTitle("GO TEST LIST").
		SetLines(testList...).
		Start().
		ChosenLine()
	return v, ok
}

func buildScreen() *menuscreen.MenuScreen {
	screen, err := menuscreen.NewMenuScreen()
	if err != nil {
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

func extractTestSuites(f *ast.File) []*ast.Struct {
	if f == nil {
		return nil
	}
	sList := make([]*ast.Struct, 0, len(f.StructList))
	for i, s := range f.StructList {
		if s == nil {
			continue
		}
		for _, field := range s.FieldList {
			if field.TypeName == "suite.Suite" {
				sList = append(sList, f.StructList[i])
			}
		}
	}
	return sList
}

func extractSuiteTestMethods(s *ast.Struct) []string {
	if s == nil {
		return nil
	}
	suiteName := s.Name
	methodList := make([]string, 0, len(s.MethodList))
	for _, m := range s.MethodList {
		methodList = append(methodList, fmt.Sprintf("%s/%s", suiteName, m.Name))
	}
	return methodList
}
