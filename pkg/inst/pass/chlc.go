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
	chMakeOutsideUsedInGs map[*ast.GoStmt]map[string]struct{}
	chMakeInGs            map[*ast.GoStmt]map[string]struct{}
	requiredOrtImport     bool
}

func NewChLifeCyclePass() *ChLifeCyclePass {
	return &ChLifeCyclePass{
		chMakeOutsideUsedInGs: make(map[*ast.GoStmt]map[string]struct{}),
		chMakeInGs:            make(map[*ast.GoStmt]map[string]struct{}),
	}
}
func (p *ChLifeCyclePass) Deps() []string {
	return nil
}

func (p *ChLifeCyclePass) Before(iCtx *inst.InstContext) {
}

func (p *ChLifeCyclePass) After(iCtx *inst.InstContext) {

	for g, chs := range p.chMakeOutsideUsedInGs {
		for ch := range chs {

			addRefCall := NewArgCallExpr(oraclertImportName, "CurrentGoAddCh", []ast.Expr{
				&ast.BasicLit{
					ValuePos: 0,
					Kind:     token.STRING,
					Value:    ch,
				},
			})
			if f, ok := g.Call.Fun.(*ast.FuncLit); ok {
				// append to beginning of goroutine
				f.Body.List = append([]ast.Stmt{addRefCall}, f.Body.List...)
			}
			p.requiredOrtImport = true

		}

	}

	if p.requiredOrtImport {
		inst.AddImport(iCtx.FS, iCtx.AstFile, oraclertImportName, oraclertImportPath)
	}

}

func (p *ChLifeCyclePass) GetPostApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	return nil
}

func (p *ChLifeCyclePass) GetPreApply(iCtx *inst.InstContext) func(*astutil.Cursor) bool {
	return func(c *astutil.Cursor) bool {
		defer func() {
			if r := recover(); r != nil { // This is allowed. If we insert node into nodes not in slice, we will meet a panic
				// For example, we may identified a receive in select and wanted to insert a function call before it, then this function will panic
			}
		}()

		switch concrete := c.Node().(type) {
		case *ast.GoStmt:
			p.chMakeOutsideUsedInGs[concrete] = make(map[string]struct{})
			p.chMakeInGs[concrete] = make(map[string]struct{})
			ast.Inspect(concrete, func(n ast.Node) bool {
				switch tn := n.(type) {
				case *ast.Ident:
					if tyObj, exist := iCtx.Type.Defs[tn]; exist && tyObj.Type() != nil {
						if strings.HasPrefix(tyObj.Type().String(), "chan ") {
							p.chMakeInGs[concrete][tn.Name] = struct{}{}
						}

					} else if tyObj, exist := iCtx.Type.Uses[tn]; exist && tyObj.Type() != nil {
						if strings.HasPrefix(tyObj.Type().String(), "chan ") {
							if _, exist := p.chMakeInGs[concrete][tn.Name]; !exist {
								p.chMakeOutsideUsedInGs[concrete][tn.Name] = struct{}{}
							}
						}

					}
				}

				return true
			})

		}
		return true
	}
}
