package jira

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := loadConfig("jira.yaml")
	if config.BaseUrl == "" || config.Login == "" || config.Password == "" {
		t.Error("config should not be empty")
	}
}

func TestWriteConfig(t *testing.T) {
	config := &Config{}
	config.BaseUrl = "http://google.com"
	config.Login = "rob"
	config.Password = "passwd1"

	config.writeConfig("jira.yaml")
}
