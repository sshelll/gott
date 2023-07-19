package util

// Copyright (c) 2023 sshelll, the gott authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildGoTestRegExpr(t *testing.T) {
	assert.Equal(t, "", BuildGoTestRegExpr())
	assert.Equal(t, "^TestFoo$", BuildGoTestRegExpr("TestFoo"))
	assert.Equal(t, "^TestFoo$|^TestBar$", BuildGoTestRegExpr("TestFoo", "TestBar"))
	assert.Equal(t, "^TestFoo$|^TestBar$|^TestBaz$", BuildGoTestRegExpr("TestFoo", "TestBar", "TestBaz"))
}
