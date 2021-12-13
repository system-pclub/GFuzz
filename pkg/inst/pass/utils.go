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
					if to := iCtx.Type.Defs[objIdent]; to == nil || to.Type() == nil {
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

func SelectorCallerHasTypes(iCtx *inst.InstContext, selExpr *ast.SelectorExpr, trueIfUnknown bool, tys ...string) bool {
	t := getSelectorCallerType(iCtx, selExpr)
	if t == "" && trueIfUnknown {
		return true
	}
	for _, ty := range tys {
		if ty == t {
			return true
		}
	}

	return false
}
