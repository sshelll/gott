package util

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"strings"
)

func BuildGoTestRegExpr(testList ...string) string {
	if len(testList) == 0 {
		return ""
	}
	return "^" + strings.Join(testList, "$\\|^") + "$"
}
