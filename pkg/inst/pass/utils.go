package pass

import (
	"gfuzz/pkg/inst"
	"go/ast"
)

func getSelectorCallerType(iCtx *inst.InstContext, selExpr *ast.SelectorExpr) string {
	if callerIdent, ok := selExpr.X.(*ast.Ident); ok {
		if callerIdent.Obj != nil {
			if objStmt, ok := callerIdent.Obj.Decl.(*ast.AssignStmt); ok {
				if objIdent, ok := objStmt.Lhs[0].(*ast.Ident); ok {
					if to := iCtx.Type.Defs[objIdent]; to == nil {
						return ""
					} else {
						return to.Type().String()
					}

				}
			}
		}
	}

	return ""
}

func SelectorCallerHasTypes(iCtx *inst.InstContext, selExpr *ast.SelectorExpr, tys ...string) bool {
	t := getSelectorCallerType(iCtx, selExpr)
	for _, ty := range tys {
		if ty == t {
			return true
		}
	}

	return false
}
