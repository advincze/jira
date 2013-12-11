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

func (c *Config) writeConfig(fileName string) {
	bytes, err := goyaml.Marshal(c)
	panicerr(err)

	err = ioutil.WriteFile(fileName, bytes, 0777)
	panicerr(err)
}
