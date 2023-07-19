package util

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"log"
	"strconv"
)

func StrToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("[gott] convert str to number failed: %v\n", err)
	}
	return n
}
