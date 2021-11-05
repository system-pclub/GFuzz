package pass

import (
	"gfuzz/pkg/inst"
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

// ChResPass, Channel Record Pass. This pass instrumented at
// following four channel related operations:
// send, recv, make, close
type ChRecPass struct {
}

func (p *ChRecPass) Name() string {
	return "chrec"
}

func (p *ChRecPass) Deps() []string {
	return nil
}

func (p *ChRecPass) GetPostApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	return nil
}

func (p *ChRecPass) GetPreApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	var needOracleRtImport bool
	return func(c *astutil.Cursor) bool {
		defer func() {
			if r := recover(); r != nil { // This is allowed. If we insert node into nodes not in slice, we will meet a panic
				// For example, we may identified a receive in select and wanted to insert a function call before it, then this function will panic
			}
		}()
		switch concrete := c.Node().(type) {

		// channel send operation
		case *ast.SendStmt:
			intID := int(getNewOpID())
			newCall := NewArgCallExpr(oraclertImportName, "StoreOpInfo", []ast.Expr{&ast.BasicLit{
				ValuePos: 0,
				Kind:     token.STRING,
				Value:    "\"Send\"",
			}, &ast.BasicLit{
				ValuePos: 0,
				Kind:     token.INT,
				Value:    strconv.Itoa(intID),
			}})
			c.InsertBefore(newCall) // Insert the call to store this operation's type and ID into goroutine local storage
			addRecord(strconv.Itoa(intID) + ":chsend")
			needOracleRtImport = true

		// channel make operation
		case *ast.AssignStmt:
			if len(concrete.Rhs) == 1 {
				if callExpr, ok := concrete.Rhs[0].(*ast.CallExpr); ok { // calling on the right
					if funcIdent, ok := callExpr.Fun.(*ast.Ident); ok {
						if funcIdent.Name == "make" {
							if len(callExpr.Args) == 1 { // This is a make operation
								if _, ok := callExpr.Args[0].(*ast.ChanType); ok {
									intID := int(getNewOpID())
									newCall := NewArgCallExpr(oraclertImportName, "StoreChMakeInfo", []ast.Expr{
										concrete.Lhs[0],
										&ast.BasicLit{
											ValuePos: 0,
											Kind:     token.INT,
											Value:    strconv.Itoa(intID),
										}})
									c.InsertAfter(newCall)
									addRecord(strconv.Itoa(intID) + ":chmake")
									needOracleRtImport = true
								}
							}
						}
					}
				}
			}

		// channel recv operation
		case *ast.ExprStmt:
			if unaryExpr, ok := concrete.X.(*ast.UnaryExpr); ok {
				if unaryExpr.Op == token.ARROW { // This is a receive operation
					intID := int(getNewOpID())
					newCall := NewArgCallExpr(oraclertImportName, "StoreOpInfo", []ast.Expr{&ast.BasicLit{
						ValuePos: 0,
						Kind:     token.STRING,
						Value:    "\"Recv\"",
					}, &ast.BasicLit{
						ValuePos: 0,
						Kind:     token.INT,
						Value:    strconv.Itoa(intID),
					}})
					c.InsertBefore(newCall)
					addRecord(strconv.Itoa(intID) + ":chrecv")

					needOracleRtImport = true
				}
			} else if callExpr, ok := concrete.X.(*ast.CallExpr); ok { // like `close(ch)` or `mu.Lock()`
				if funcIdent, ok := callExpr.Fun.(*ast.Ident); ok { // like `close(ch)`
					// channel close operation
					if funcIdent.Name == "close" {
						intID := int(getNewOpID())
						newCall := NewArgCallExpr(oraclertImportName, "StoreOpInfo", []ast.Expr{&ast.BasicLit{
							ValuePos: 0,
							Kind:     token.STRING,
							Value:    "\"Close\"",
						}, &ast.BasicLit{
							ValuePos: 0,
							Kind:     token.INT,
							Value:    strconv.Itoa(intID),
						}})
						c.InsertBefore(newCall)
						addRecord(strconv.Itoa(intID) + ":chclose")
						needOracleRtImport = true
					}
				}
			}
		}

		if needOracleRtImport {
			inst.AddImport(iCtx.FS, iCtx.AstFile, oraclertImportName, oraclertImportPath)
		}

		return true
	}
}
