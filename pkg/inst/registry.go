package inst

func NewPassRegistry() *PassRegistry {
	return &PassRegistry{
		n2p: make(map[string]InstPass),
	}
}

// AddPass adds a unique pass into registry
func (r *PassRegistry) AddPass(pass InstPass) error {
	_, exist := r.n2p[pass.Name()]
	if exist {
		return &PassExistedError{Name: pass.Name()}
	}
	r.n2p[pass.Name()] = pass
	return nil
}

// GetPass returns the pass with given name
func (r *PassRegistry) GetPass(name string) (InstPass, error) {
	p, exist := r.n2p[name]
	if exist {
		return p, nil
	}

	return nil, &NoPassError{Name: name}
}

// ListOfPasses return a list of registered passes
func (r *PassRegistry) ListOfPasses() []InstPass {
	passes := make([]InstPass, 0, len(r.n2p))

	for _, p := range r.n2p {
		passes = append(passes, p)
	}
	return passes
}

func (r *PassRegistry) ListOfPassNames() []string {
	passes := make([]string, 0, len(r.n2p))

	for _, p := range r.n2p {
		passes = append(passes, p.Name())
	}
	return passes
}

// HasPass return true if pass registered, false otherwise
func (r *PassRegistry) HasPass(name string) bool {
	_, err := r.GetPass(name)
	return err == nil
}
