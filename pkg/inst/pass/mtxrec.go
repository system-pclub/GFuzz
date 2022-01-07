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

func (p *MtxRecPass) Before(iCtx *inst.InstContext) {
	iCtx.SetMetadata(MetadataKeyRequiredOrtImport, false)
}

func (p *MtxRecPass) After(iCtx *inst.InstContext) {
	needOracleRtImportItf, _ := iCtx.GetMetadata(MetadataKeyRequiredOrtImport)
	needOracleRtImport := needOracleRtImportItf.(bool)
	if needOracleRtImport {
		inst.AddImport(iCtx.FS, iCtx.AstFile, oraclertImportName, oraclertImportPath)
	}
}

func (p *MtxRecPass) Deps() []string {
	return nil
}

func (p *MtxRecPass) GetPostApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	return nil
}

func (p *MtxRecPass) GetPreApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	return func(c *astutil.Cursor) bool {
		defer func() {
			if r := recover(); r != nil {
			}
		}()
		switch concrete := c.Node().(type) {
		case *ast.ExprStmt:
			if callExpr, ok := concrete.X.(*ast.CallExpr); ok {
				if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok { // like `mu.Lock()`
					if !SelectorCallerHasTypes(iCtx, selectorExpr, true, "sync.Mutex", "*sync.Mutex", "sync.RWMutex", "*sync.RWMutex") {
						return true
					}
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
