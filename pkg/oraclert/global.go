package oraclert

import (
	"fmt"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/env"
	"gfuzz/pkg/selefcm"
	"io/ioutil"
	"os"
	"strconv"
	"time"
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
	timeDivStr := os.Getenv(env.ORACLERT_TIME_DIVIDE)
	if timeDivStr == "" {
		timeDivStr = "1"
	}
	var err error
	time.DurDivideBy, err = strconv.Atoi(timeDivStr)
	if err != nil {
		time.DurDivideBy = 1
		fmt.Println("Failed to set time.DurDivideBy. time.DurDivideBy is set to 1. Err:", err)
	} 

	if rtConfigFile == "" {
		return
	}
	data, err := ioutil.ReadFile(rtConfigFile)
	if err == nil {
		ortConfig, err = config.Deserilize(data)
		if err == nil {
			// read oracle configuration file successfully

			// We can create different strategies according to our needs
			efcmStrat = selefcm.NewSelectCaseInOrder(ortConfig.SelEfcm.Efcms)
			selTimeout = ortConfig.SelEfcm.SelTimeout
			fmt.Printf("[oraclert] selefcm timeout: %d", selTimeout)

		} else {
			fmt.Printf("OracleRt deserilize config: %s", err)
		}
	} else {
		fmt.Printf("OracleRt read config: %s", err)
	}
}
