package ast

import (
	"fmt"
	"strings"
)

type File struct {
	Path       string
	ImportList []*Import
	StructList []*Struct
	FuncList   []*Func
}

type Import struct {
	Alias string
	Pkg   string
}

type Struct struct {
	Name       string
	FieldList  []*Field
	MethodList []*Method
}

type Field struct {
	NameList []string
	TypeName string
}

type Method struct {
	Name      string
	TypeName  string
	IsPtrRecv bool
}

type Func struct {
	Name         string
	IsTest       bool
	IsSuiteEntry bool
	SuiteName    string
}

func (f *File) Print() {
	fmt.Println("-----------------------------")
	fmt.Printf("********* File Path *********\n%s\n", f.Path)
	fmt.Println("********* Import List *********")
	for _, imp := range f.ImportList {
		imp.Print()
	}
	fmt.Println("********* Func List *********")
	for _, fn := range f.FuncList {
		fn.Print()
	}
	fmt.Println("******** Struct List ********")
	for _, s := range f.StructList {
		s.Print()
		fmt.Println()
	}
	fmt.Println("-----------------------------")
}

func (i *Import) Print() {
	if i == nil {
		return
	}
	fmt.Printf("Alias: %s, Pkg: %s\n", i.Alias, i.Pkg)
}

func (s *Struct) Print() {
	if s == nil {
		return
	}
	fmt.Println("Type Name: " + s.Name)
	fmt.Println("Fields:")
	for _, f := range s.FieldList {
		f.Print()
	}
	for _, m := range s.MethodList {
		m.Print()
	}
}

func (f *Field) Print() {
	if f == nil {
		return
	}
	fmt.Printf("Field: %v %s\n", f.NameList, f.TypeName)
}

func (m *Method) Print() {
	if m == nil {
		return
	}
	if m.IsPtrRecv {
		fmt.Printf("Method: (recv *%s) %s(...)\n", m.TypeName, m.Name)
	} else {
		fmt.Printf("Method: (recv %s) %s(...)\n", m.TypeName, m.Name)
	}
}

func (m *Method) IsTest() bool {
	return m != nil && strings.HasPrefix(m.Name, "Test")
}

func (fn *Func) Print() {
	if fn == nil {
		return
	}
	if fn.IsTest {
		fmt.Printf("Func: %s(...) - Test\n", fn.Name)
	} else {
		fmt.Printf("Func: %s(...) - Normal\n", fn.Name)
	}
}
