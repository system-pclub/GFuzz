package oraclert

import (
	"fmt"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/env"
	"gfuzz/pkg/selefcm"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	// oracle runtime configuration
	ortConfig *config.Config
	// select enforcement strategy
	efcmStrat selefcm.SelectCaseStrategy
	// timeout each select should wait
	selTimeout    int
	ortDebug      bool
	ortBenchmark  bool
	ortStdoutFile string
	ortOutputFile string
)

func init() {
	ortDebug, _ = strconv.ParseBool(os.Getenv(env.ORACLERT_DEBUG))
	ortBenchmark, _ = strconv.ParseBool(os.Getenv(env.ORACLERT_BENCHMARK))
	ortStdoutFile = os.Getenv(env.ORACLERT_STDOUT_FILE)
	ortOutputFile = os.Getenv(env.ORACLERT_OUTPUT_FILE)
	rtConfigFile := os.Getenv(env.ORACLERT_CONFIG_FILE)
	data, err := ioutil.ReadFile(rtConfigFile)
	if err == nil {
		ortConfig, err = config.Deserilize(data)
		if err == nil {
			// read oracle configuration file successfully

			// We can create different strategies according to our needs
			efcmStrat = selefcm.NewSelectCaseInOrder(ortConfig.SelEfcm.Efcms)
			selTimeout = ortConfig.SelEfcm.SelTimeout
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
