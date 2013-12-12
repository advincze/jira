jira-client
===========

simple jira REST API client

	$go get github.com/advincze/jira


### configuration

jira.yaml in your application folder

	baseurl: http://my.jira.server
	auth:
  	  login: jira_user
  	  password: his_passwd


## usage

	boards := jira.FetchBoards()
	sprints := jira.FetchSprints()
	sprintDetails := jira.FetchSprintDetails()
	
	
## documentation

[http://godoc.org/github.com/advincze/jira](http://godoc.org/github.com/advincze/jira)