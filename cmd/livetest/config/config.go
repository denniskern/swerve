package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	envDynoPW = "SWERVE_DYNO_DEFAULT_PW"
)

func Get() LiveConfig {
	cfg := LiveConfig{}
	fh, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(fh, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	pw := os.Getenv(envDynoPW)
	if pw == "" {
		log.Fatal(fmt.Errorf("provide the password for the dyno user %s as env %s", cfg.DynoPw, envDynoPW))
	}
	cfg.DynoPw = pw

	return cfg

}
