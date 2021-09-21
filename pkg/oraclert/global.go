package oraclert

import (
	rtConfig "gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/selefcm"
	"io/ioutil"
	"os"
)

var (
	// oracle runtime configuration
	config *rtConfig.Config
	// select enforcement strategy
	efcmStrat selefcm.SelectCaseStrategy
	// timeout each select should wait
	selTimeout int
)

func init() {
	rtConfigFile := os.Getenv(ORACLERT_CONFIG_FILE)
	data, err := ioutil.ReadFile(rtConfigFile)
	if err == nil {
		config, err = rtConfig.Deserilize(data)
		if err == nil {
			// read oracle configuration file successfully

			// We can create different strategies according to our needs
			efcmStrat = selefcm.NewSelectCaseInOrder(config.SelEfcm.Efcms)
			selTimeout = config.SelEfcm.SelTimeout
		} else {
			println(err)
		}
	} else {
		println(err)
	}
}
