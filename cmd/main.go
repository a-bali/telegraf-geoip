package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/a-bali/telegraf-geoip/plugins/processors/geoip"

	"github.com/influxdata/telegraf/plugins/common/shim"
)

var configFile = flag.String("config", "", "path to the config file for this plugin")
var err error

func main() {
	// parse command line options
	flag.Parse()

	// create the shim. This is what will run your plugins.
	s := shim.New()

	// Check for settings from a config toml file,
	err = s.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err loading input: %s\n", err)
		os.Exit(1)
	}

	// run the plugin until stdin closes or we receive a termination signal
	if err := s.Run(shim.PollIntervalDisabled); err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		os.Exit(1)
	}
}
