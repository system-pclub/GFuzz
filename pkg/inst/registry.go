package inst

func NewPassRegistry() *PassRegistry {
	return &PassRegistry{
		n2p: make(map[string]InstPassConstructor),
	}
}

// AddPass adds a unique pass into registry
func (r *PassRegistry) Register(name string, passc InstPassConstructor) error {
	_, exist := r.n2p[name]
	if exist {
		return &PassExistedError{Name: name}
	}
	r.n2p[name] = passc
	return nil
}

// GetPass returns the pass with given name
func (r *PassRegistry) GetNewPassInstance(name string) (InstPass, error) {
	c, exist := r.n2p[name]
	if exist {
		return c(), nil
	}

	return nil, &NoPassError{Name: name}
}

func (r *PassRegistry) ListOfPassNames() []string {
	passes := make([]string, 0, len(r.n2p))

	for n, _ := range r.n2p {
		passes = append(passes, n)
	}
	return passes
}

// HasPass return true if pass registered, false otherwise
func (r *PassRegistry) HasPass(name string) bool {
	_, exist := r.n2p[name]
	return exist
}
