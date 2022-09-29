package util

import (
	"testing"

	"github.com/SCU-SJL/sinfra/ast"
	"github.com/stretchr/testify/assert"
)

func TestExtractFuncs(t *testing.T) {
	f, err := ast.Parse("./mock_test.go")
	assert.Nil(t, err)
	testList := make([]string, 0, 16)
	testList = append(testList, ExtractTestFuncs(f)...)
	testList = append(testList, ExtractTestifySuiteEntryFuncs(f)...)
	assert.NotZero(t, len(testList))
}
