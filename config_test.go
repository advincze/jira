package jira

import (
	"io/ioutil"
	"os"

	"testing"
)

func writeTempFile(data string) *os.File {
	file, _ := ioutil.TempFile("", "config-test")
	ioutil.WriteFile(file.Name(), []byte(data), 0644)
	return file
}

func TestLoadConfig(t *testing.T) {
	data := `
baseurl: http://jira.example.com
auth:
  login: jirauser
  password: jirapass`
	file := writeTempFile(data)
	defer os.Remove(file.Name())

	config := LoadConfig(file.Name())
	if config == nil {
		t.Errorf("the config should not be empty")
	}
}
