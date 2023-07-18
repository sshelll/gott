package core

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseURI parse uri to filepath and pos, uri format is like '/xxx/xxx/xxx_test.go:59'
func ParseURI(uri string) (filepath string, pos int, err error) {
	splited := strings.Split(uri, ":")
	if len(splited) != 2 {
		err = fmt.Errorf("[gott] uri format error, please check if it's like '/xxx/xxx/xxx_test.go:59', exit...")
		return
	}

	filepath, posStr := splited[0], splited[1]

	// check if filepath is a go test file
	if !strings.HasSuffix(filepath, "_test.go") {
		err = fmt.Errorf("[gott] filepath [%s] is not a test file, exit...", filepath)
		return
	}

	// check if pos is number
	pos, err = strconv.Atoi(posStr)
	if err != nil {
		err = fmt.Errorf("[gott] pos [%s] is not a number, exit...", posStr)
		return
	}

	return
}
