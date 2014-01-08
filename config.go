package jira

import (
	"io/ioutil"
	"launchpad.net/goyaml"
)

type Config struct {
	BaseUrl string
	Auth    *Auth
}

func LoadConfig(fileName string) (config *Config) {

	file, err := ioutil.ReadFile(fileName)
	panicerr(err)

	err = goyaml.Unmarshal(file, &config)
	panicerr(err)
	return
}
