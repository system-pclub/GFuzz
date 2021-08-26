package inst

// runPasses executes given passes with provided instrumentation context
func runPasses(iCtx *InstContext, passes []InstPass) error {
	for _, p := range passes {
		err := p.Run(iCtx)
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
		pass, err := r.GetPass(passName)
		if err != nil {
			return err
		}
		passes = append(passes, pass)
	}
	return runPasses(iCtx, passes)
}

// NewInstContext creates a InstContext by given Golang source file
func NewInstContext(goSrcFile string) (*InstContext, error) {

}
