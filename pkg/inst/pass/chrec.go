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
type ChRecPass struct{}

func (p *ChRecPass) Name() string {
	return "chrec"
}

func (p *ChRecPass) Run(iCtx *inst.InstContext) error {
	inst.AddImport(iCtx.FS, iCtx.AstFile, "gooracle", "gooracle")
	iCtx.AstFile = astutil.Apply(iCtx.AstFile, instChOps, nil).(*ast.File)
	return nil
}

func (p *ChRecPass) Deps() []string {
	return nil
}

func instChOps(c *astutil.Cursor) bool {
	switch concrete := c.Node().(type) {

	// channel send operation
	case *ast.SendStmt:
		intID := int(Uint16OpID)
		newCall := NewArgCallExpr("gooracle", "StoreOpInfo", []ast.Expr{&ast.BasicLit{
			ValuePos: 0,
			Kind:     token.STRING,
			Value:    "\"Send\"",
		}, &ast.BasicLit{
			ValuePos: 0,
			Kind:     token.INT,
			Value:    strconv.Itoa(intID),
		}})
		c.InsertBefore(newCall) // Insert the call to store this operation's type and ID into goroutine local storage
		records = append(records, strconv.Itoa(intID)+":chsend")
		Uint16OpID++

	// channel make operation
	case *ast.AssignStmt:
		if len(concrete.Rhs) == 1 {
			if callExpr, ok := concrete.Rhs[0].(*ast.CallExpr); ok { // calling on the right
				if funcIdent, ok := callExpr.Fun.(*ast.Ident); ok {
					if funcIdent.Name == "make" {
						if len(callExpr.Args) == 1 { // This is a make operation
							if _, ok := callExpr.Args[0].(*ast.ChanType); ok {
								intID := int(Uint16OpID)
								newCall := NewArgCallExpr("gooracle", "StoreChMakeInfo", []ast.Expr{
									concrete.Lhs[0],
									&ast.BasicLit{
										ValuePos: 0,
										Kind:     token.INT,
										Value:    strconv.Itoa(intID),
									}})
								c.InsertAfter(newCall)
								records = append(records, strconv.Itoa(intID)+":chmake")
								Uint16OpID++
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
				intID := int(Uint16OpID)
				newCall := NewArgCallExpr("gooracle", "StoreOpInfo", []ast.Expr{&ast.BasicLit{
					ValuePos: 0,
					Kind:     token.STRING,
					Value:    "\"Recv\"",
				}, &ast.BasicLit{
					ValuePos: 0,
					Kind:     token.INT,
					Value:    strconv.Itoa(intID),
				}})
				c.InsertBefore(newCall)
				records = append(records, strconv.Itoa(intID)+":chrecv")
				Uint16OpID++
			}
		} else if callExpr, ok := concrete.X.(*ast.CallExpr); ok { // like `close(ch)` or `mu.Lock()`
			if funcIdent, ok := callExpr.Fun.(*ast.Ident); ok { // like `close(ch)`
				// channel close operation
				if funcIdent.Name == "close" {
					intID := int(Uint16OpID)
					newCall := NewArgCallExpr("gooracle", "StoreOpInfo", []ast.Expr{&ast.BasicLit{
						ValuePos: 0,
						Kind:     token.STRING,
						Value:    "\"Close\"",
					}, &ast.BasicLit{
						ValuePos: 0,
						Kind:     token.INT,
						Value:    strconv.Itoa(intID),
					}})
					c.InsertBefore(newCall)
					records = append(records, strconv.Itoa(intID)+":chclose")
					Uint16OpID++
				}
			}
		}
	}

	return true
}
