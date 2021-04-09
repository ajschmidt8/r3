package shared

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func ReadConfig() (config ConfigInterface) {
	ymlBytes, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("cannot read config.yaml file: %v", err)
	}
	err = yaml.Unmarshal(ymlBytes, &config)
	if err != nil {
		log.Fatalf("cannot decode config.yaml: %v", err)
	}
	return
}
