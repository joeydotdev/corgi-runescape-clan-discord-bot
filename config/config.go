package config

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Token string
}

func Load() *Configuration {
	var config *Configuration

	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	config = &Configuration{}
	err = json.Unmarshal(file, config)
	if err != nil {
		panic(err)
	}

	return config
}
