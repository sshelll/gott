package core

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func BuildTestAllExpr(testList []string) string {
	buf := strings.Builder{}
	cnt := len(testList)

	for i, testName := range testList {
		buf.WriteString("^")
		buf.WriteString(testName)
		buf.WriteString("$")
		if i < cnt-1 {
			buf.WriteString("\\|")
		}
	}

	return buf.String()
}

func ExecGoTest(testName string, flags ...string) {
	args := bytes.Buffer{}
	flagLen := len(flags)
	for i := 0; i < flagLen; i++ {
		args.WriteString(flags[i])
		args.WriteString(" ")
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
