package selectstrtg

// SelectCaseController defines the interface for guiding select case during application running
// It tries to provide answer to question 'when application reach the select at given file and line, which case it should choose'
// This functionality expected to be cooperated with instrumentation.
type SelectCaseStrategy interface {
	GetCase(filename string, line, numOfCases int) int
}
