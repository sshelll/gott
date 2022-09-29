package util

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/SCU-SJL/menuscreen"
)

func ChooseTestFile() (file string, ok bool) {
	testFiles := lsTestFiles()
	if len(testFiles) == 0 {
		return
	}
	screen := buildScreen()
	defer screen.Fini()
	_, file, ok = screen.SetTitle("GO TEST FILES").
		SetLines(testFiles...).
		Start().
		ChosenLine()
	if ok {
		file = "./" + file
	}
	return
}

func ChooseTest(testList []string) (tname string, ok bool) {
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

func buildScreen() *menuscreen.MenuScreen {
	screen, err := menuscreen.NewMenuScreen()
	if err != nil || screen == nil {
		log.Fatalf("init screen controller failed: %v\n", err)
	}
	return screen
}
