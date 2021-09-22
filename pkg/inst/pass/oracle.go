package pass

import (
	"gfuzz/pkg/inst"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// OraclePass, instrument the oracle entry and defer function call to trigger oracle bug detection
type OraclePass struct{}

func (p *OraclePass) Name() string {
	return "oracle"
}

func (p *OraclePass) Run(iCtx *inst.InstContext) error {
	inst.AddImport(iCtx.FS, iCtx.AstFile, "gooracle", "gooracle")
	iCtx.AstFile = astutil.Apply(iCtx.AstFile, instOracle, nil).(*ast.File)
	return nil
}

func (p *OraclePass) Deps() []string {
	return nil
}

var additionalNode ast.Stmt

func instOracle(c *astutil.Cursor) bool {
	defer func() {
		if r := recover(); r != nil { // This is allowed. If we insert node into nodes not in slice, we will meet a panic
			// For example, we may identified a receive in select and wanted to insert a function call before it, then this function will panic

			//fmt.Printf("Recover in pre(): c.Name(): %s\n", c.Name())
			//position := currentFSet.Position(c.Node().Pos())
			//position.Filename, _ = filepath.Abs(position.Filename)
			//fmt.Printf("\tLocation: %s\n", position.Filename + ":" + strconv.Itoa(position.Line))
		}
	}()
	if additionalNode != nil && c.Node() == additionalNode {
		newAssign := &ast.AssignStmt{
			Lhs:    nil,
			TokPos: 0,
			Tok:    token.DEFINE,
			Rhs:    nil,
		}
		newObject := &ast.Object{
			Kind: ast.Var,
			Name: "oracleEntry",
			Decl: newAssign,
			Data: nil,
			Type: nil,
		}
		newIdent := &ast.Ident{
			Name: "oracleEntry",
			Obj:  newObject,
		}
		newAssign.Lhs = []ast.Expr{
			newIdent,
		}
		newAssign.Rhs = []ast.Expr{
			NewArgCall("gooracle", "BeforeRun", nil),
		}
		c.InsertBefore(newAssign)

		newAfterTestCall := NewArgCall("gooracle", "AfterRun", []ast.Expr{
			newIdent,
		})
		newDefer := &ast.DeferStmt{
			Defer: 0,
			Call:  newAfterTestCall,
		}
		c.InsertBefore(newDefer)
		additionalNode = nil
	}

	if cStmt, ok := c.Node().(ast.Stmt); ok {
		for _, recvAndFirstStmt := range vecRecvAndFirstStmt {
			if recvAndFirstStmt.firstStmt == cStmt {
				newCall := NewArgCallExpr("gooracle", "CurrentGoAddValue", []ast.Expr{
					&ast.Ident{
						Name: recvAndFirstStmt.recvName,
						Obj:  recvAndFirstStmt.recvObj,
					},
					&ast.Ident{
						Name: "nil",
						Obj:  nil,
					},
					&ast.BasicLit{
						ValuePos: 0,
						Kind:     token.INT,
						Value:    "0",
					},
				})
				c.InsertBefore(newCall)
			}
		}
	}

	switch concrete := c.Node().(type) {

	case *ast.FuncDecl:
		if strings.HasPrefix(concrete.Name.Name, "Test") {
			if concrete.Body != nil && len(concrete.Body.List) != 0 {
				firstStmt := concrete.Body.List[0]
				additionalNode = firstStmt
			}

		} else if concrete.Recv != nil && concrete.Body != nil {
			if len(concrete.Recv.List) == 1 && len(concrete.Body.List) > 0 {
				recvField := concrete.Recv.List[0]
				if len(recvField.Names) == 1 {
					ident := recvField.Names[0]
					if ident.Name != "" {
						recvAndFirstStmt := &RecvAndFirstStmt{
							recvName:  ident.Name,
							firstStmt: concrete.Body.List[0],
							recvObj:   ident.Obj,
						}
						vecRecvAndFirstStmt = append(vecRecvAndFirstStmt, recvAndFirstStmt)
					}
				}
			}
		}

	default:
	}

	return true
}