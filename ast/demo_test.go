package ast

import (
	"context"
	"testing"

	testifySuite "github.com/stretchr/testify/suite"
)

type DummyInt int

type FooTestSuite struct {
	testifySuite.Suite
	Ctx context.Context
}

func TestFoo(t *testing.T) {
	testifySuite.Run(t, &FooTestSuite{})
}

func (s *FooTestSuite) BeforeTest(suiteName, testName string) {
	s.T().Logf("FOO BEFORE TEST - [%s-%s]", suiteName, testName)
}

func (*FooTestSuite) TestCase() {

}

func (FooTestSuite) TestCase2() {

}

type BarTestSuite struct {
	testifySuite.Suite
	Ctx    context.Context
	S1, S2 string
}

func TestBar(t *testing.T) {
	tt := new(BarTestSuite)
	testifySuite.Run(t, tt)
}

func (s *BarTestSuite) BeforeTest(suiteName, testName string) {
	s.T().Logf("BAR BEFORE TEST - [%s-%s]", suiteName, testName)
}

func (*BarTestSuite) TestCase1() {

}

func (*BarTestSuite) TestCase2() {

}

func (*BarTestSuite) OtherFunc() {

}

func (*BarTestSuite) privateFunc() {

}
