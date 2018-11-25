package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

// ArgsConfig parses the `cfg` parameter from the command line which should contain a valid configuration file
// Check config/Config for possible config values
func ArgsConfig() Config {
	f := flag.String("cfg", "", "-cfg=/path/to/config")
	flag.Parse()
	bts, err := ioutil.ReadFile(*f)
	if *f == "" || err != nil {
		log.Fatalf("Failed to read config file since it is empty or does not exist. Error: %v", err)
	}
	var cfg Config
	err = json.Unmarshal(bts, &cfg)
	if err != nil {
		log.Fatalf("Failed to parse config file since it is not a valid JSON. Error: %v", err)
	}
	if !cfg.valid() {
		log.Fatalf("Provided config is not valid since there are either empty domains and or missing credentials (including fallback) for at least one update.")
	}
	return cfg
}
