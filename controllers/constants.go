package controllers

import "os"

const (
	PdServiceID  string = "PdServiceID"
	UserID       string = "USERID"
	PdUrl        string = "https://<>.pagerduty.com"
	PdScheduleId string = "PdScheduleId"

	EpicId  string = "OPS-0000"
	JiraUrl string = "https://id.atlassian.net"

	//SlackColor    string = "#00FF00"
	MonAlertsLink = "https://<>.wavefront.com/u/<>"
	MonURL        = "https:/<>.wavefront.com"
)

var SlackColor string = os.Getenv("SLACK_COLOR")
