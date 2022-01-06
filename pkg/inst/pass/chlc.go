package pass

import (
	"gfuzz/pkg/inst"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// ChLifeCyclePass tries to mark channel's end
type ChLifeCyclePass struct {
}

func (p *ChLifeCyclePass) Name() string {
	return "chlc"
}

func (p *ChLifeCyclePass) Deps() []string {
	return nil
}

func (p *ChLifeCyclePass) Before(iCtx *inst.InstContext) {
	iCtx.SetMetadata(MetadataKeyRequiredOrtImport, false)
}

func (p *ChLifeCyclePass) After(iCtx *inst.InstContext) {
	needOracleRtImportItf, _ := iCtx.GetMetadata(MetadataKeyRequiredOrtImport)
	needOracleRtImport := needOracleRtImportItf.(bool)
	if needOracleRtImport {
		inst.AddImport(iCtx.FS, iCtx.AstFile, oraclertImportName, oraclertImportPath)
	}
}

func (p *ChLifeCyclePass) GetPostApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	return nil
}

func TryAddChDerefGoroutine(iCtx *inst.InstContext, ident *ast.Ident) {
	if to := iCtx.Type.Defs[ident]; to == nil || to.Type() == nil {
		return
	} else {
		println(to.Type().String())

		if to.Type().String() == "chan" {
			to.Parent()
		}
	}
}

func (p *ChLifeCyclePass) GetPreApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	return func(c *astutil.Cursor) bool {
		defer func() {
			if r := recover(); r != nil { // This is allowed. If we insert node into nodes not in slice, we will meet a panic
				// For example, we may identified a receive in select and wanted to insert a function call before it, then this function will panic
			}
		}()

		switch concrete := c.Node().(type) {

		// try to find channel assignment
		case *ast.AssignStmt:
			// if len(concrete.Lhs) == 0:
			// 	return true
			identExpr := concrete.Lhs[0]
			if ident, ok := identExpr.(*ast.Ident); ok {
				if tyObj, ok := iCtx.Type.Defs[ident]; ok {
					if tyObj == nil || tyObj.Type() == nil {
						return true
					}
					if strings.HasPrefix(tyObj.Type().String(), "chan ") {
						// confirm this x :=, x type is chan
						var blockStmtNode ast.Node
						for {
							blockStmtNode = c.Parent()
							if blockStmtNode == nil {
								break
							}
							if bs, ok := blockStmtNode.(*ast.BlockStmt); ok {
								// TODO: this one should be added when go block found
								// addRefCall := NewArgCallExpr(oraclertImportName, "AddChRefFromG", []ast.Expr{
								// 	&ast.BasicLit{
								// 		ValuePos: 0,
								// 		Kind:     token.STRING,
								// 		Value:    ident.Name,
								// 	},
								// })
								rmRefCall := NewArgCallExpr(oraclertImportName, "RemoveChRefFromG", []ast.Expr{
									&ast.BasicLit{
										ValuePos: 0,
										Kind:     token.STRING,
										Value:    ident.Name,
									},
								})
								bs.List = append(bs.List, rmRefCall)
								break
							}
						}
					}
				}

			}

		}

		return true
	}
}
