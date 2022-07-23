package main

import (
	"bytes"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	ana "github.com/fuji8/gotypeconverter/analysis"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"golang.org/x/tools/go/packages"
	"honnef.co/go/tools/go/ast/astutil"
)

func SuggestedFix(path string, line, character int) (string, *protocol.Range) {
	var srcType, dstType types.Type
	var rng protocol.Range

	lastSlashI := strings.LastIndex(path, "/")
	log.Infof("%v", lastSlashI)
	dir := path[:lastSlashI]

	pkgs, _ := packages.Load(&packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  dir,
	}, "") // mod name or empty
	pkg := pkgs[0]

	fset := pkg.Fset
	// aFile, _ := parser.ParseFile(fset, path+"/tmp.go", nil, 0)
	// get ast.File from package.packgae
	aFile := pkg.Syntax[0]
	// get token.File from fset
	fset.Iterate(func(tFile *token.File) bool {
		if tFile.Name() != path {
			return true
		}
		node, _ := astutil.PathEnclosingInterval(aFile, tFile.LineStart(line)+token.Pos(character), tFile.LineStart(line)+token.Pos(character))
		ast.Inspect(node[1], func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				ff := fset.File(fn.Pos())
				rng.Start.Line = uint32(ff.Line(fn.Pos()))
				rng.End.Line = uint32(ff.Line(fn.End()))
				t := pkg.TypesInfo.TypeOf(fn.Name)
				ts, _ := t.(*types.Signature)
				if ts.Params().Len() == 1 && ts.Results().Len() == 1 {
					srcType = ts.Params().At(0).Type()
					dstType = ts.Results().At(0).Type()
				}
			}
			return true
		})
		return true
	})

	if srcType == nil || dstType == nil {
		return "", nil
	}

	fn := ana.InitFuncMaker(pkg.Types)
	fn.MakeFunc(ana.InitType(dstType, "src"), ana.InitType(srcType, ""))

	buf := &bytes.Buffer{}
	buf.Write(fn.WriteBytes())
	out := buf.String()
	log.Info(out)
	return out, &rng
}
