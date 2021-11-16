package env

// All environment consumed by oracle runtime
const (
	// ORACLERT_CONFIG_FILE is an environment variable, set by fuzzer,
	// pass to oraclert, a filepath to oracle runtime configuration file.
	ORACLERT_CONFIG_FILE = "ORACLERT_CONFIG_FILE"

	// ORACLERT_DEBUG is an boolean environment variable
	ORACLERT_DEBUG = "ORACLERT_DEBUG"

	// ORACLERT_BENCHMARK is an boolean environment variable
	ORACLERT_BENCHMARK = "ORACLERT_BENCHMARK"

	// ORACLERT_STDOUT_FILE is an environment variable indicates the file in where
	// oracle runtime should write bug report (by default bug report will be directly printed
	// out to the stdout, however, in some cases, if trigger program from go test, timeout causes some stdout been ate by caller)
	ORACLERT_STDOUT_FILE = "ORACLERT_STDOUT_FILE"

	// ORACLERT_OUTPUT_FILE is an environment variable dumps selects, channels, xorLoc usage
	ORACLERT_OUTPUT_FILE = "ORACLERT_OUTPUT_FILE"
)
