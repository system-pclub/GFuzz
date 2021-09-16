package fuzz

type FuzzStage string

const (
	// InitStage simply run the empty without any mutation
	InitStage FuzzStage = "init"

	// DeterStage is to create input by tweak select choice one by one
	DeterStage FuzzStage = "deter"

	// CalibStage choose an input from queue to run (prepare for rand)
	CalibStage FuzzStage = "calib"

	// RandStage randomly mutate select choice
	RandStage FuzzStage = "rand"
)
