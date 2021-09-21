package selefcm

// SelectCaseStrategy defines the interface for guiding select case during application running
// It tries to provide answer to question 'when application reach the select at given file and line, which case it should choose'
// This functionality expected to be cooperated with instrumentation.
type SelectCaseStrategy interface {
	GetCase(selectID string) int
}

// SelEfcm, stands for select enforcement, defines which select case will be enforced during runtime
type SelEfcm struct {
	// ID is where select is, usually is filename + linenumber
	ID string `json:"id"`

	// Case to be enforced
	Case int `json:"case"`
}

type SelEfcmConfig struct {
	SelTimeout int       `json:"sel_timeout"`
	Efcms      []SelEfcm `json:"efcms"`
}
