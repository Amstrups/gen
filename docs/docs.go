package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type FileArg struct {
	Module   string
	EvalArgs []EvalArg
}

type EvalArg struct {
	Commands     []string
	Preprocess   []string
	FunctionName string
}

var (
	nlReg    = `\n*\s*`
	startReg = `^/\*`
	endReg   = `\*/$`
	seqReg   = `{\s*([a-z]*,)*([a-z]*)\s*}`

	funcReg       = `FUNCTION:\s*` + seqReg
	helpReg       = `HELP:\s*{\[[A-z]|\s]*}`
	preReg        = `PRE:\s*` + seqReg
	preOptReg     = "(" + preReg + "){0,1}"
	genCommentReg = startReg + nlReg +
		funcReg + nlReg +
		helpReg + nlReg +
		preOptReg + nlReg +
		endReg
)

func getNames(decls []ast.Decl) {
	defRegex, err := regexp.Compile(genCommentReg)
	if err != nil {
		panic(err)
	}

	for _, d := range decls {
		switch dt := d.(type) {
		case *ast.FuncDecl:
			fmt.Println(dt.Name)
			for _, c := range dt.Doc.List {
				fmt.Println(defRegex.MatchString(c.Text))
				fmt.Println(c.Text)
				finds := defRegex.FindStringSubmatch(c.Text)
				fmt.Println(finds)
				for _, f := range finds {
					fmt.Print("Find: ")
					fmt.Println(f)
				}
			}
		default:
			fmt.Println("unknown type")
		}
	}
}

func getTree(fname string) *ast.File {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, fname, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	//spew.Dump(f)
	return f
}

func Make() {
	f := getTree("../arit/modules/numbers/numbers.go")

	getNames(f.Decls)
	//spew.Dump(f.Comments)
}

func Make2() {
	numberArgs := EvalArg{
		Commands:     []string{"num", "number"},
		Preprocess:   []string{},
		FunctionName: "number",
	}

	nextArgs := EvalArg{
		Commands:     []string{"next", "n"},
		Preprocess:   []string{"trim"},
		FunctionName: "next",
	}

	fileArgs := FileArg{
		Module:   "number",
		EvalArgs: []EvalArg{numberArgs, nextArgs},
	}

	evalFuncMap := template.FuncMap{
		"len":  func(list []string) int { return len(list) },
		"join": func(list []string) string { return strings.Join(list, "\", \"") },
	}

	tmplFile := "eval.tmpl"
	tmpl, err := template.New(tmplFile).Funcs(evalFuncMap).ParseFiles(tmplFile)
	if err != nil {
		panic(tmpl)
	}

	err = tmpl.Execute(os.Stdout, fileArgs)
	if err != nil {
		panic(err)
	}

}
