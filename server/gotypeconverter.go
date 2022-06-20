package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	ana "github.com/fuji8/gotypeconverter/analysis"
	"github.com/fuji8/gotypeconverter/ui"
	"golang.org/x/tools/go/packages"
	"honnef.co/go/tools/go/ast/astutil"
)

func SuggestedFix() {
	path := "/home/fuji/workspace/lsp/tmp/fuga"
	var srcType, dstType types.Type

	pkgs, _ := packages.Load(&packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  path,
	}, "") // mod name or empty
	pkg := pkgs[0]

	fset := pkg.Fset
	// aFile, _ := parser.ParseFile(fset, path+"/tmp.go", nil, 0)
	// get ast.File from package.packgae
	aFile := pkg.Syntax[0]
	// get token.File from fset
	fset.Iterate(func(tFile *token.File) bool {
		if tFile.Name() != path+"/tmp.go" {
			return true
		}
		node, _ := astutil.PathEnclosingInterval(aFile, tFile.LineStart(20)+5, tFile.LineStart(20)+5)
		fmt.Println(node, pkg)
		ast.Inspect(node[1], func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				t := pkg.TypesInfo.TypeOf(fn.Name)
				ts, _ := t.(*types.Signature)
				if ts.Params().Len() == 1 && ts.Results().Len() == 1 {
					srcType = ts.Params().At(0).Type()
					dstType = ts.Results().At(0).Type()
				}
				fmt.Println()
			}
			return true
		})
		return true
	})

	fn := ana.InitFuncMaker(pkg.Types)
	fn.MakeFunc(ana.InitType(dstType, "src"), ana.InitType(srcType, ""))

	out, _ := ui.NoInfoGeneration(fn)
	fmt.Println(out)
}
