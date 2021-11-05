package pass

import (
	"fmt"
	"gfuzz/pkg/inst"
	"go/ast"
	"go/token"
	"path/filepath"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

// SelEfcm, select enforcement pass, instrument the 'select' keyword,
// turn it into a select with multiple cases, each case represent one
// original select's case and a timeout case.
type SelEfcmPass struct{}

func (p *SelEfcmPass) Name() string {
	return "selefcm"
}

func (p *SelEfcmPass) Before(iCtx *inst.InstContext) {
	iCtx.SetMetadata(MetadataKeyRequiredOrtImport, false)
}

func (p *SelEfcmPass) After(iCtx *inst.InstContext) {
	needOracleRtImportItf, _ := iCtx.GetMetadata(MetadataKeyRequiredOrtImport)
	needOracleRtImport := needOracleRtImportItf.(bool)
	if needOracleRtImport {
		inst.AddImport(iCtx.FS, iCtx.AstFile, oraclertImportName, oraclertImportPath)
	}
}

func (p *SelEfcmPass) Deps() []string {
	return nil
}

func (p *SelEfcmPass) GetPostApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	return nil
}

func (p *SelEfcmPass) GetPreApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	var (
		numOfSelects uint32
	)

	return func(c *astutil.Cursor) bool {
		defer func() {
			if r := recover(); r != nil { // This is allowed. If we insert node into nodes not in slice, we will meet a panic
				// For example, we may identified a receive in select and wanted to insert a function call before it, then this function will panic

				//fmt.Printf("Recover in pre(): c.Name(): %s\n", c.Name())
				//position := currentFSet.Position(c.Node().Pos())
				//position.Filename, _ = filepath.Abs(position.Filename)
				//fmt.Printf("\tLocation: %s\n", position.Filename + ":" + strconv.Itoa(position.Line))
			}
		}()

		switch concrete := c.Node().(type) {

		case *ast.SelectStmt:
			numOfSelects += 1
			positionOriSelect := iCtx.FS.Position(concrete.Select)
			positionOriSelect.Filename, _ = filepath.Abs(positionOriSelect.Filename)

			// store the original select
			oriSelect := SelectStruct{
				StmtSelect:    concrete,
				VecCommClause: nil,
				VecOp:         nil,
				VecBody:       nil,
			}
			for _, stmtCommClause := range concrete.Body.List {
				commClause, _ := stmtCommClause.(*ast.CommClause)
				oriSelect.VecCommClause = append(oriSelect.VecCommClause, commClause)
				oriSelect.VecOp = append(oriSelect.VecOp, commClause.Comm)
				vecContent := []ast.Stmt{}
				vecContent = append(vecContent, commClause.Body...)
				oriSelect.VecBody = append(oriSelect.VecBody, vecContent)
			}

			// create a switch
			newSwitch := &ast.SwitchStmt{
				Switch: 0,
				Init:   nil,
				Tag: NewArgCall(oraclertImportName, "GetSelEfcmSwitchCaseIdx", []ast.Expr{
					&ast.BasicLit{ // first parameter: filename:linenumber
						ValuePos: 0,
						Kind:     token.STRING,
						Value:    fmt.Sprintf("\"%s\"", positionOriSelect.Filename),
					},
					&ast.BasicLit{ // second parameter: linenumber of original select
						ValuePos: 0,
						Kind:     token.STRING,
						Value:    fmt.Sprintf("\"%s\"", strconv.Itoa(positionOriSelect.Line)),
					},
					&ast.BasicLit{
						ValuePos: 0,
						Kind:     token.INT,
						Value:    strconv.Itoa(len(oriSelect.VecCommClause)),
					}}),
				Body: &ast.BlockStmt{
					Lbrace: 0,
					List:   nil,
					Rbrace: 0,
				},
			}
			iCtx.SetMetadata(MetadataKeyRequiredOrtImport, true)

			vecCaseClause := []ast.Stmt{}
			// The number of switch case is (the number of non-default select cases + 1)
			for i, stmtOp := range oriSelect.VecOp {

				// if the case's expression is nil, it means this case is a default. We don't intrument this case
				// but we instrument other cases
				if stmtOp == nil {
					continue
				}

				newCaseClause := &ast.CaseClause{
					Case:  0,
					List:  nil,
					Colon: 0,
					Body:  nil,
				}
				newBasicLit := &ast.BasicLit{
					ValuePos: 0,
					Kind:     token.INT,
					Value:    strconv.Itoa(i),
				}
				newCaseClause.List = []ast.Expr{newBasicLit}

				// the case's content is one select statement
				newSelect := &ast.SelectStmt{
					Select: 0,
					Body:   &ast.BlockStmt{},
				}
				firstSelectCase := &ast.CommClause{
					Case:  0,
					Comm:  oriSelect.VecOp[i],
					Colon: 0,
					Body:  copyStmtBody(oriSelect.VecBody[i]),
				}
				secondSelectCase := &ast.CommClause{
					Case: 0,
					Comm: &ast.ExprStmt{X: &ast.UnaryExpr{
						OpPos: 0,
						Op:    token.ARROW,
						X:     NewArgCall(oraclertImportName, "SelEfcmTimeout", nil),
					}},
					Colon: 0,
					Body: []ast.Stmt{
						// The first line is a call to gooracle.StoreLastMySwitchChoice(-1)
						// The second line is a copy of original select
						&ast.ExprStmt{X: NewArgCall(oraclertImportName, "StoreLastMySwitchChoice", []ast.Expr{&ast.UnaryExpr{
							OpPos: 0,
							Op:    token.SUB,
							X: &ast.BasicLit{
								ValuePos: 0,
								Kind:     token.INT,
								Value:    "1",
							},
						}})},
						copySelect(oriSelect.StmtSelect)},
				}
				newSelect.Body.List = append(newSelect.Body.List, firstSelectCase, secondSelectCase)
				iCtx.SetMetadata(MetadataKeyRequiredOrtImport, true)

				newCaseClause.Body = []ast.Stmt{newSelect}

				// add the created case to vector
				vecCaseClause = append(vecCaseClause, newCaseClause)
			}

			// add one default case to switch
			newCaseClauseDefault := &ast.CaseClause{
				Case:  0,
				List:  nil,
				Colon: 0,
				Body: []ast.Stmt{
					// The first line is a call to gooracle.StoreLastMySwitchChoice(-1)
					// The second line is a copy of original select
					&ast.ExprStmt{X: NewArgCall(oraclertImportName, "StoreLastMySwitchChoice", []ast.Expr{&ast.UnaryExpr{
						OpPos: 0,
						Op:    token.SUB,
						X: &ast.BasicLit{
							ValuePos: 0,
							Kind:     token.INT,
							Value:    "1",
						},
					}})},
					copySelect(oriSelect.StmtSelect)},
			}
			iCtx.SetMetadata(MetadataKeyRequiredOrtImport, true)
			vecCaseClause = append(vecCaseClause, newCaseClauseDefault)

			newSwitch.Body.List = vecCaseClause

			// Insert the new switch before the select
			c.InsertBefore(newSwitch)

			// Delete the original select
			c.Delete()

		default:
		}
		return true
	}

}

var sliceStrNoInstr = []string{
	"src/runtime",
	"src/gooracle",
	"src/sync",
	"src/reflect",
	"src/syscall",
	"src/bufio",
	"src/fmt",
	"src/os",
	"src/strconv",
	"src/strings",
	"src/time",
	"src/bytes",
	"src/hash",
}

type RecvAndFirstStmt struct {
	recvName  string
	firstStmt ast.Stmt
	recvObj   *ast.Object
}

type SelectStruct struct {
	StmtSelect    *ast.SelectStmt   // StmtSelect.Body.List is a vec of CommClause
	VecCommClause []*ast.CommClause // a CommClause is a case and its content in select
	VecOp         []ast.Stmt        // The operations of cases. Nil is default
	VecBody       [][]ast.Stmt      // The content of cases
}

type SwitchStruct struct {
	StmtSwitch    *ast.SwitchStmt // StmtSwitch.Body.List is a vector of CaseClause
	Tag           ast.Expr
	VecCaseClause []*ast.CaseClause // a CaseClause is a case and its content in switch.
	VecVecExpr    [][]ast.Expr      // The expressions of each case.
	VecBody       [][]ast.Stmt      // The content of cases
}

// Deprecated:
func copyOp(stmtOp ast.Stmt) ast.Stmt {
	var result ast.Stmt
	// the stmtOp is either *ast.SendStmt or *ast.ExprStmt
	switch concrete := stmtOp.(type) {
	// TODO: could be "x := <-ch"
	case *ast.SendStmt:
		oriChanIdent, _ := concrete.Chan.(*ast.Ident)
		newSend := &ast.SendStmt{
			Chan: &ast.Ident{
				NamePos: 0,
				Name:    oriChanIdent.Name,
				Obj:     oriChanIdent.Obj,
			},
			Arrow: 0,
			Value: concrete.Value,
		}
		result = newSend
	case *ast.ExprStmt:
		oriUnaryExpr, _ := concrete.X.(*ast.UnaryExpr)
		newRecv := &ast.ExprStmt{X: &ast.UnaryExpr{
			OpPos: 0,
			Op:    token.ARROW,
			X:     oriUnaryExpr.X,
		}}
		result = newRecv
	}

	return result
}

func copyStmtBody(stmtBody []ast.Stmt) []ast.Stmt {
	result := []ast.Stmt{}
	for _, stmt := range stmtBody {
		result = append(result, stmt)
	}
	return result
}

func copySelect(oriSelect *ast.SelectStmt) *ast.SelectStmt {
	result := &ast.SelectStmt{
		Select: 0,
		Body:   oriSelect.Body,
	}
	return result
}

// imports reports whether f has an import with the specified name and path.
func imports(f *ast.File, name, path string) bool {
	for _, s := range f.Imports {
		importedName := importName(s)
		importedPath := importPath(s)
		if importedName == name && importedPath == path {
			return true
		}
	}
	return false
}

// importName returns the name of s,
// or "" if the import is not named.
func importName(s *ast.ImportSpec) string {
	if s.Name == nil {
		return ""
	}
	return s.Name.Name
}

// importPath returns the unquoted import path of s,
// or "" if the path is not properly quoted.
func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err != nil {
		return ""
	}
	return t
}

func NewArgCall(strPkg, strCallee string, vecExprArg []ast.Expr) *ast.CallExpr {
	newIdentPkg := &ast.Ident{
		NamePos: token.NoPos,
		Name:    strPkg,
		Obj:     nil,
	}
	newIdentCallee := &ast.Ident{
		NamePos: token.NoPos,
		Name:    strCallee,
		Obj:     nil,
	}
	newCallSelector := &ast.SelectorExpr{
		X:   newIdentPkg,
		Sel: newIdentCallee,
	}
	newCall := &ast.CallExpr{
		Fun:      newCallSelector,
		Lparen:   token.NoPos,
		Args:     vecExprArg,
		Ellipsis: token.NoPos,
		Rparen:   token.NoPos,
	}

	return newCall
}

func NewArgCallExpr(strPkg, strCallee string, vecExprArg []ast.Expr) *ast.ExprStmt {
	newCall := NewArgCall(strPkg, strCallee, vecExprArg)
	newExpr := &ast.ExprStmt{X: newCall}
	return newExpr
}

func handleCallExpr(ce *ast.CallExpr) (ast.Node, bool) {
	name := getCallExprLiteral(ce)
	switch name {
	case "errors.Wrap":
		return rewriteWrap(ce), true
	case "errors.Wrapf":
		return rewriteWrap(ce), true
	case "errors.Errorf":
		return newErrorfExpr(ce.Args), true
	default:
		return ce, true
	}
}

func handleImportDecl(gd *ast.GenDecl) (ast.Node, bool) {
	// Ignore GenDecl's that aren't imports.
	if gd.Tok != token.IMPORT {
		return gd, true
	}
	// Push "errors" to the front of specs so formatting will sort it with
	// core libraries and discard pkg/errors.
	newSpecs := []ast.Spec{
		&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"errors"`}},
	}
	for _, s := range gd.Specs {
		im, ok := s.(*ast.ImportSpec)
		if !ok {
			continue
		}
		if im.Path.Value == `"github.com/pkg/errors"` {
			continue
		}
		newSpecs = append(newSpecs, s)
	}
	gd.Specs = newSpecs
	return gd, true
}

func rewriteWrap(ce *ast.CallExpr) *ast.CallExpr {
	// Rotate err to the end of a new args list
	newArgs := make([]ast.Expr, len(ce.Args)-1)
	copy(newArgs, ce.Args[1:])
	newArgs = append(newArgs, ce.Args[0])

	// If the format string is a fmt.Sprintf call, we can unwrap it.
	c, name := getCallExpr(newArgs[0])
	if c != nil && name == "fmt.Sprintf" {
		newArgs = append(c.Args, newArgs[1:]...)
	}

	// If the format string is a literal, we can rewrite it:
	//     "......" -> "......: %w"
	// Otherwise, we replace it with a binary op to add the wrap code:
	//     SomeNonLiteral -> SomeNonLiteral + ": %w"
	fmtStr, ok := newArgs[0].(*ast.BasicLit)
	if ok {
		// Strip trailing `"` and append wrap code and new trailing `"`
		fmtStr.Value = fmtStr.Value[:len(fmtStr.Value)-1] + `: %w"`
	} else {
		binOp := &ast.BinaryExpr{
			X:  newArgs[0],
			Op: token.ADD,
			Y:  &ast.BasicLit{Kind: token.STRING, Value: `": %w"`},
		}
		newArgs[0] = binOp
	}

	return newErrorfExpr(newArgs)
}

func getCallExpr(n ast.Node) (*ast.CallExpr, string) {
	c, ok := n.(*ast.CallExpr)
	if !ok {
		return nil, ""
	}
	name := getCallExprLiteral(c)
	if name == "" {
		return nil, ""
	}
	return c, name
}

func getCallExprLiteral(c *ast.CallExpr) string {
	s, ok := c.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}

	i, ok := s.X.(*ast.Ident)
	if !ok {
		return ""
	}

	return i.Name + "." + s.Sel.Name
}

func newErrorfExpr(args []ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "fmt"},
			Sel: &ast.Ident{Name: "Errorf"},
		},
		Args: args,
	}
}
