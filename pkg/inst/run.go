package inst

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
)

// runPasses executes given passes with provided instrumentation context
func runPasses(iCtx *InstContext, passes []InstPass) error {
	for _, p := range passes {
		err := RunPass(p, iCtx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Run executes passes with given a list of pass name and instrumentation context.
func Run(iCtx *InstContext, r *PassRegistry, passNames []string) error {
	var passes = make([]InstPass, 0, len(passNames))
	for _, passName := range passNames {
		pass, err := r.GetNewPassInstance(passName)
		if err != nil {
			return err
		}
		passes = append(passes, pass)
	}
	return runPasses(iCtx, passes)
}

func RunPass(p InstPass, iCtx *InstContext) error {
	p.Before(iCtx)
	iCtx.AstFile = astutil.Apply(iCtx.AstFile, p.GetPreApply(iCtx), p.GetPostApply(iCtx)).(*ast.File)
	p.After(iCtx)
	return nil
}
