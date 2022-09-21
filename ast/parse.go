package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type fileParser struct {
	path    string
	absPath string
	fset    *token.FileSet
	astf    *ast.File
}

func NewFileParser(path string) *fileParser {
	return &fileParser{
		path: path,
	}
}

func (fp *fileParser) Parse() (*File, error) {

	if err := fp.parseAST(); err != nil {
		return nil, err
	}

	file := &File{
		Path: fp.absPath,
	}

	tDecls, fDecls := fp.extractStructAndFunc()

	var (
		funcList        []*Func
		structMethodMap map[string][]*Method = make(map[string][]*Method)
	)

	for _, fDecl := range fDecls {
		m, f := fp.parseFuncDecl(fDecl)
		if m != nil {
			structMethodMap[m.TypeName] = append(structMethodMap[m.TypeName], m)
		}
		if f != nil {
			funcList = append(funcList, f)
		}
	}
	file.FuncList = funcList

	for _, tDecl := range tDecls {
		structInfo := fp.parseStructDecl(tDecl)
		if structInfo != nil {
			structInfo.MethodList = structMethodMap[structInfo.Name]
			file.StructList = append(file.StructList, structInfo)
		}
	}

	return file, nil

}

func (fp *fileParser) parseAST() error {
	fp.fset = token.NewFileSet()
	absPath, err := filepath.Abs(fp.path)
	if err != nil {
		return err
	}
	fp.absPath = absPath
	fp.astf, err = parser.ParseFile(fp.fset, absPath, nil, parser.AllErrors)
	return err
}

func (fp *fileParser) extractStructAndFunc() (structs, funcs []ast.Decl) {

	if fp.astf == nil {
		return
	}

	for i := range fp.astf.Decls {

		decl := fp.astf.Decls[i]

		gDecl, ok := decl.(*ast.GenDecl)
		if ok && gDecl.Tok == token.TYPE {
			structs = append(structs, decl)
		}

		if _, ok = decl.(*ast.FuncDecl); ok {
			funcs = append(funcs, decl)
		}

	}

	return

}

func (fp *fileParser) parseStructDecl(decl ast.Decl) *Struct {

	gDecl := decl.(*ast.GenDecl)

	spec := gDecl.Specs[0].(*ast.TypeSpec)

	st, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	structInfo := &Struct{
		Name: spec.Name.Name,
	}

	for _, field := range st.Fields.List {

		fieldInfo := &Field{}

		// extract field names
		fNames := make([]string, 0, len(field.Names))
		for _, name := range field.Names {
			fNames = append(fNames, name.Name)
		}
		fieldInfo.NameList = fNames

		// extract field type
		if idt, ok := field.Type.(*ast.Ident); ok {
			fieldInfo.TypeName = idt.Name
		} else if expr, ok := field.Type.(*ast.SelectorExpr); ok {
			pkg := expr.X.(*ast.Ident).Name
			clz := expr.Sel.Name
			fieldInfo.TypeName = pkg + "." + clz
		} else {
			continue
		}

		structInfo.FieldList = append(structInfo.FieldList, fieldInfo)

	}

	return structInfo

}

func (fp *fileParser) parseFuncDecl(decl ast.Decl) (method *Method, fn *Func) {

	fDecl := decl.(*ast.FuncDecl)

	// is func
	if fDecl.Recv == nil {
		fn = &Func{
			Name:   fDecl.Name.Name,
			IsTest: fp.isTestFunc(fDecl),
		}
		if fn.IsTest {
			suiteName, isSuiteEntry := fp.isSuiteEntry(fDecl)
			fn.IsSuiteEntry = isSuiteEntry
			fn.SuiteName = suiteName
		}
		return
	}

	// is method
	var isPtrRecv bool
	var typeName, methodName string

	t := fDecl.Recv.List[0].Type
	if starExpr, ok := t.(*ast.StarExpr); ok {
		isPtrRecv = true
		typeName = starExpr.X.(*ast.Ident).Name
	} else {
		isPtrRecv = false
		typeName = t.(*ast.Ident).Name
	}

	methodName = fDecl.Name.Name

	method = &Method{
		Name:      methodName,
		TypeName:  typeName,
		IsPtrRecv: isPtrRecv,
	}

	return

}

func (fp *fileParser) isTestFunc(fDecl *ast.FuncDecl) bool {

	if !strings.HasPrefix(fDecl.Name.Name, "Test") {
		return false
	}

	params := fDecl.Type.Params

	if params == nil || len(params.List) != 1 {
		return false
	}

	tname := params.List[0].Names
	if len(tname) != 1 {
		return false
	}

	starExpr, ok := params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}

	selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	idt, ok := selectorExpr.X.(*ast.Ident)
	if !ok || idt.Name != "testing" {
		return false
	}

	switch selectorExpr.Sel.Name {
	case "T", "B", "M", "F", "TB", "PB":
		return true
	}

	return false

}

func (fp *fileParser) isSuiteEntry(fDecl *ast.FuncDecl) (suiteName string, isOK bool) {

	if !fp.isTestFunc(fDecl) {
		return "", false
	}

	tName := fDecl.Type.Params.List[0].Names[0].Name

	if fDecl.Body == nil {
		return "", false
	}

	if fDecl.Body == nil {
		return "", false
	}

	for _, stmt := range fDecl.Body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		if exprStmt.X == nil {
			continue
		}
		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		if callExpr.Fun == nil || len(callExpr.Args) != 2 {
			continue
		}
		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		// TODO: add pkg name judgement
		if selectorExpr.Sel == nil || selectorExpr.Sel.Name != "Run" {
			continue
		}
		if len(callExpr.Args) != 2 {
			continue
		}
		idt, ok := callExpr.Args[0].(*ast.Ident)
		if !ok {
			continue
		}
		if idt.Name != tName {
			continue
		}
		// parse new(suite)
		if callNewExpr, ok := callExpr.Args[1].(*ast.CallExpr); ok {
			idt, ok := callNewExpr.Fun.(*ast.Ident)
			if ok && idt.Name == "new" {
				structIdt, ok := callNewExpr.Args[0].(*ast.Ident)
				if ok {
					return structIdt.Name, true
				}
			}
		}
		// parse &suite{}
		if unaryExpr, ok := callExpr.Args[1].(*ast.UnaryExpr); ok {
			if unaryExpr.Op == token.AND {
				compositeLit, ok := unaryExpr.X.(*ast.CompositeLit)
				if ok {
					idt, isIdt := compositeLit.Type.(*ast.Ident)
					if isIdt {
						return idt.Name, true
					}
				}
			}
		}
	}

	return "", false

}
