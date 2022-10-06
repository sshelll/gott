package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/SCU-SJL/gott/util"
	"github.com/SCU-SJL/sinfra/ast"
)

func main() {

	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		println("use -p to print test func name only, or else gott would exec 'go test'")
		return
	}

	f, ok := util.ChooseTestFile()
	if !ok {
		println("[gott] no files were chosen, exit...")
		return
	}

	fInfo, err := ast.Parse(f)
	if err != nil {
		log.Fatalln("[gott] ast parse failed:", err.Error())
	}

	testList := append(util.ExtractTestFuncs(fInfo), util.ExtractTestifySuiteEntryFuncs(fInfo)...)

	if len(testList) == 0 {
		println("[gott] no tests were found, exit...")
		return
	}

	testName, ok := util.ChooseTest(testList)
	if !ok {
		println("[gott] no tests were chosen, exit...")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "-p" {
		fmt.Print(testName)
		return
	}

	execGoTest(testName)

}

func execGoTest(testName string) {

	args := bytes.Buffer{}
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			args.WriteString(os.Args[i])
			args.WriteString(" ")
		}
	}

	goTestCmd := fmt.Sprintf("go test %s -test.run %s", args.String(), testName)
	execCmd(goTestCmd, true)

}

func execCmd(sh string, useStdIO bool) {

	cmd := exec.Command("bash", "-c", sh)
	if useStdIO {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		log.Fatalf("[gott] exec cmd '%s' failed, err = %v\n", cmd.String(), err)
	}

}
