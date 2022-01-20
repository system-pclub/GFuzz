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
func (p *WgRecPass) Before(iCtx *inst.InstContext) {
	iCtx.SetMetadata(MetadataKeyRequiredOrtImport, false)
}

func (p *WgRecPass) After(iCtx *inst.InstContext) {
	needOracleRtImportItf, _ := iCtx.GetMetadata(MetadataKeyRequiredOrtImport)
	needOracleRtImport := needOracleRtImportItf.(bool)
	if needOracleRtImport {
		inst.AddImport(iCtx.FS, iCtx.AstFile, oraclertImportName, oraclertImportPath)
	}
}

func (p *WgRecPass) Deps() []string {
	return nil
}

func (p *WgRecPass) GetPostApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	return nil
}

func (p *WgRecPass) GetPreApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {

	return func(c *astutil.Cursor) bool {
		defer func() {
			if r := recover(); r != nil {
			}
		}()
		switch concrete := c.Node().(type) {
		case *ast.AssignStmt:

		case *ast.ExprStmt:

			if callExpr, ok := concrete.X.(*ast.CallExpr); ok {
				if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok { // like `mu.Lock()`
					if !SelectorCallerHasTypes(iCtx, selectorExpr, true, "sync.WaitGroup", "*sync.WaitGroup") {
						return true
					}

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
						intID := int(getNewOpID())
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
						addRecord(strconv.Itoa(intID) + ":trad" + op)
						iCtx.SetMetadata(MetadataKeyRequiredOrtImport, true)
					}
				}

			}
		}

		return true

	}

}
