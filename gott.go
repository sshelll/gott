package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sshelll/gott/util"
	"github.com/sshelll/sinfra/ast"
)

// lua print(require'util'.get_pos_under_cursor())
// lua print(require'gott'.get_test_name('/Users/shaojiale/Codes/github/mine/gott/util/mock_test.go', 423))

func main() {
	var (
		testName string
		done     bool
	)

	switch true {
	// help
	case len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help"):
		helpMode()
		return

	// pos
	case len(os.Args) > 1 && os.Args[1] == "--pos":
		testName, done = posMode()
		if done {
			return
		}
		if len(os.Args) > 3 && os.Args[3] == "-p" {
			print(testName)
			return
		}
		util.ExecGoTest(testName, os.Args[3:]...)

	// interactive
	default:
		testName, done = interactiveMode()
		if done {
			return
		}
		if len(os.Args) > 1 && os.Args[1] == "-p" {
			print(testName)
			return
		}
		util.ExecGoTest(testName, os.Args[1:]...)
	}
}

func helpMode() {
	log.Println("Use --pos to pass an uri with absolute filepath to exec the closest test\n" +
		"\tNOTE: this flag must be the first arg if you try to use it!!!\n" +
		"\tFor example: \n" +
		"\t\t'gott --pos /Users/sshelll/go/src/gott/xxx_test.go:59'\n" +
		"\t\tIn this way, gott would try to exec the closest go test func to the uri with no flags\n" +
		"\t\t'gott --pos /Users/sshelll/go/src/gott/xxx_test.go:59 -p'\n" +
		"\t\tIn this way, gott would print the closest test name of the uri\n" +
		"\t\t'gott --pos /Users/sshelll/go/src/gott/xxx_test.go:59 -v'\n" +
		"\t\tIn this way, gott would try to exec the closest test name of the uri with -v flag\n" +
		"\nUse -p to print the go test name instead of exec it\n" +
		"\tNOTE: This flag must be the first arg or the third arg if you try to use it!!!\n" +
		"\tFor example: \n" +
		"\t\t'gott -p'\n" +
		"\t\tIn this way, gott would print the test name with interactive mode\n" +
		"\t\t'gott --pos xxx_test.go:59 -p'\n" +
		"\t\tIn this way, gott would print the closest test name of the uri\n" +
		"\t\tPlease note that if you want to use --pos and -p together, you should put the --pos in the first, uri in the sec, and -p is the third\n" +
		"\nOtherwise you will exec go test with interactive mode, and other args will be passed to 'go test'\n" +
		"\tFor example: \n" +
		"\t\t'gott' equals 'go test'\n" +
		"\t\t'gott -v equals 'go test -v'\n" +
		"\t\t'gott -v -count=1' equals 'go test -v -count=1'\n" +
		"Thanks")
}

func interactiveMode() (testName string, done bool) {
	f, ok := util.ChooseTestFile()
	if !ok {
		log.Println("[gott] no files were chosen, exit...")
		return "", true
	}

	fInfo, err := ast.Parse(f)
	if err != nil {
		log.Fatalln("[gott] ast parse failed:", err.Error())
	}

	goTests, testifyTests := util.ExtractTestFuncs(fInfo), util.ExtractTestifySuiteTestMethods(fInfo)
	testList := append(goTests, testifyTests...)

	if len(testList) == 0 {
		log.Println("[gott] no tests were found, exit...")
		return "", true
	}

	testName, testAll, ok := util.ChooseTest(testList)
	if !ok {
		log.Println("[gott] no tests were chosen, exit...")
		return "", true
	}

	if testAll {
		testName = util.BuildTestAllExpr(goTests)
	}

	return
}

func posMode() (testName string, done bool) {
	if len(os.Args) < 3 {
		log.Println("[gott] no uri was passed, exit...")
		return "", true
	}

	uri := os.Args[2]
	f, pos, err := util.ParseURI(uri)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fInfo, err := ast.Parse(f)
	if err != nil {
		log.Fatalln("[gott] ast parse failed:", err.Error())
	}

	testName, ok := util.FindClosestTestFunc(fInfo, pos)
	if !ok {
		log.Println("[gott] no tests were found, exit...")
		return "", true
	}

	testName = fmt.Sprintf("^%s$", testName)

	return
}
