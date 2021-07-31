package main

import (
	"github.com/spatrayuni/tobs-oncall-highlights/controllers/jira"
	"github.com/spatrayuni/tobs-oncall-highlights/controllers/slackshop"
)

//Main function for Sending Highlights
func SlackMain() {

	divider := slackshop.Divider()
	preMsg := slackshop.TitleSection()
	slackshop.SlackSendMessage(preMsg)

	slackshop.SendOngoingIncidentDetails()

	slackshop.SlackSendMessage(divider)
	slackshop.SendPDDetails()

	slackshop.SlackSendMessage(divider)
	jira.GetJiraEpicDetails()

	slackshop.SlackSendMessage(divider)
	slackshop.SendOncallHighlights()

}

func main() {
	SlackMain()
}
