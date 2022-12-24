package util

import (
	"testing"

	"github.com/sshelll/sinfra/ast"
	"github.com/stretchr/testify/assert"
)

func TestExtractFuncs(t *testing.T) {
	f, err := ast.Parse("./mock_test.go")
	assert.Nil(t, err)
	testList := make([]string, 0, 16)
	testList = append(testList, ExtractTestFuncs(f)...)
	testList = append(testList, ExtractTestifySuiteTestMethods(f)...)
	assert.NotZero(t, len(testList))
}
