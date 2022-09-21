package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := NewFileParser("./demo_test.go").Parse()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	f.Print()
}

func TestGen(t *testing.T) {
	fset := token.NewFileSet()
	absPath, err := filepath.Abs("./demo_test.go")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	f, err := parser.ParseFile(fset, absPath, nil, parser.AllErrors)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	ast.Print(fset, f)
}

func rangeTypeDecls(f *ast.File) {

	if f == nil {
		return
	}

	for _, decl := range f.Decls {

		gd, ok := decl.(*ast.GenDecl)

		if !ok {
			continue
		}

		if gd.Tok != token.TYPE {
			continue
		}

		rangeTypeSpecs(gd)

	}

}

func rangeTypeSpecs(decl *ast.GenDecl) {

	if decl == nil {
		return
	}

	spec := decl.Specs[0]

	tspec := spec.(*ast.TypeSpec)
	structName := tspec.Name.Name

	st := tspec.Type.(*ast.StructType)
	if st.Fields == nil {
		return
	}

	for _, f := range st.Fields.List {
		se := f.Type.(*ast.SelectorExpr)
		fname := se.X.(*ast.Ident).Name + "." + se.Sel.Name
		// has var name
		if len(f.Names) > 0 {
			fname = f.Names[0].Name + " " + fname
		}
		fmt.Printf("%s - [%s]\n", structName, fname)
	}

}

func rangeFuncDecls(f *ast.File) {

	if f == nil {
		return
	}

	for _, decl := range f.Decls {

		fd, ok := decl.(*ast.FuncDecl)

		if !ok {
			continue
		}

		rangeClassFunc(fd)
		rangeGlobalFunc(fd)

	}

}

func rangeClassFunc(fd *ast.FuncDecl) {

	if fd.Recv == nil {
		return
	}

	structName := fd.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name

	fmt.Println("class func -", structName+"."+fd.Name.Name)

}

func rangeGlobalFunc(fd *ast.FuncDecl) {

	if fd.Recv != nil {
		return
	}

	fmt.Println("global func -", fd.Name.Name)

}
