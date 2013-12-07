package jira

import (
	"testing"
)

func _TestLoadConfig(t *testing.T) {
	config := LoadConfig("jira.yaml")
	if config.BaseUrl == "" || config.Login == "" || config.Password == "" {
		t.Error("config should not be empty")
	}
}

func _TestWriteConfig(t *testing.T) {
	config := &Config{}
	config.BaseUrl = "http://google.com"
	config.Login = "rob"
	config.Password = "passwd1"

	config.writeConfig("jira.yaml")
}
