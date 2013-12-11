package jira

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := LoadConfig("../jira.yaml")
	if config == nil || config.BaseUrl == "" || config.Auth == nil || config.Auth.Login == "" || config.Auth.Password == "" {
		t.Error("config should not be empty")
	}
}

func TestWriteConfig(t *testing.T) {
	config := &Config{
		BaseUrl: "http://google.com",
		Auth: &Auth{
			Login:    "rob",
			Password: "passwd1",
		},
	}
	config.writeConfig("jira2.yaml")

}
