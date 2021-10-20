package pass

import (
	"gfuzz/pkg/inst"
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

// MtxResPass, Mutex (and RWMutex) Record Pass.
type MtxRecPass struct{}

func (p *MtxRecPass) Name() string {
	return "mtxrec"
}

func (p *MtxRecPass) Run(iCtx *inst.InstContext) error {
	inst.AddImport(iCtx.FS, iCtx.AstFile, oraclertImportName, "gfuzz/pkg/oraclert")
	iCtx.AstFile = astutil.Apply(iCtx.AstFile, instMtxOps, nil).(*ast.File)
	return nil
}

func (p *MtxRecPass) Deps() []string {
	return nil
}

func instMtxOps(c *astutil.Cursor) bool {
	switch concrete := c.Node().(type) {
	case *ast.ExprStmt:
		if callExpr, ok := concrete.X.(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok { // like `mu.Lock()`
				var matched bool = true
				var op string
				switch selectorExpr.Sel.Name {
				case "Lock":
					op = "Lock"
				case "RUnlock":
					op = "RUnlock"
				case "RLock":
					op = "RLock"
				case "Unlock":
					op = "Unlock"
				default:
					matched = false
				}

				if matched {
					intID := int(Uint16OpID)
					newCall := NewArgCallExpr(oraclertImportName, "StoreOpInfo", []ast.Expr{&ast.BasicLit{
						ValuePos: 0,
						Kind:     token.STRING,
						Value:    "\"" + op + "\"",
					}, &ast.BasicLit{
						ValuePos: 0,
						Kind:     token.INT,
						Value:    strconv.Itoa(intID),
					}})
					c.InsertBefore(newCall)
					records = append(records, strconv.Itoa(intID)+":trad"+op)
					Uint16OpID++
				}
			}

		}
	}
	return true
}
