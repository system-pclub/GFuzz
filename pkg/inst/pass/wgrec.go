package pass

import (
	"gfuzz/pkg/inst"
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

// WgResPass, Wait Group Record Pass.
type WgRecPass struct{}

func (p *WgRecPass) Name() string {
	return "wgrec"
}

func (p *WgRecPass) Run(iCtx *inst.InstContext) error {
	inst.AddImport(iCtx.FS, iCtx.AstFile, oraclertImportName, "gfuzz/pkg/oraclert")
	iCtx.AstFile = astutil.Apply(iCtx.AstFile, instWgOps, nil).(*ast.File)
	return nil
}

func (p *WgRecPass) Deps() []string {
	return nil
}

func instWgOps(c *astutil.Cursor) bool {
	switch concrete := c.Node().(type) {
	case *ast.ExprStmt:
		if callExpr, ok := concrete.X.(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok { // like `mu.Lock()`
				var matched bool = true
				var op string
				switch selectorExpr.Sel.Name {
				case "Add":
					op = "Add"
				case "Done":
					op = "Done"
				case "Wait":
					op = "Wait"
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
