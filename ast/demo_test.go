package ast

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DummyInt int

type FooTestSuite struct {
	suite.Suite
	Ctx context.Context
}

func TestFoo(t *testing.T) {
	suite.Run(t, &FooTestSuite{})
}

func (s *FooTestSuite) BeforeTest(suiteName, testName string) {
	s.T().Logf("FOO BEFORE TEST - [%s-%s]", suiteName, testName)
}

func (*FooTestSuite) TestCase() {

}

func (FooTestSuite) TestCase2() {

}

type BarTestSuite struct {
	suite.Suite
	Ctx    context.Context
	S1, S2 string
}

func TestBar(t *testing.T) {
	suite.Run(t, new(BarTestSuite))
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
