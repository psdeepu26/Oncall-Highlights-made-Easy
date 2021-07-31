package slackshop

import (
	"github.com/slack-go/slack"
	"github.com/spatrayuni/tobs-oncall-highlights/controllers"
)

func formatNoAlertsMessage() slack.MsgOption {
	msgDetailsBlock := slack.NewTextBlockObject("mrkdwn", "`No Alerts Firing`", false, false)
	msgDetailsSection := slack.NewSectionBlock(msgDetailsBlock, nil, nil)

	msg := slack.MsgOptionBlocks(
		msgDetailsSection,
	)

	return msg
}

func getAffectedhosts(alertAffectedHosts map[string][]string, alert string) string {
	affectedHosts := ""
	if len(alertAffectedHosts[alert]) > 0 {
		for _, host := range alertAffectedHosts[alert] {
			affectedHosts = host + "\n" + affectedHosts
		}
	} else {
		affectedHosts = "verify alert for more details" + "\n" + affectedHosts
	}

	return affectedHosts
}

func formatAlertAttachment(alertList []string, Severity string, alertAffectedHosts map[string][]string, alertLink map[string]string, count int) {
	var sevPreText string

	//preText := slackshop.GetEmoji("wavefront") + "*Oncall Highlights for today:*" + slackshop.GetEmoji("wavefront") + "\n\n"
	if count == 0 {
		alertSummary := " *Summary Of Current Firing Alerts - " + controllers.MonAlertsLink + " * \n"
		sevText := Severity + "_Total Firing - " + "_ "
		sevPreText = alertSummary + GetEmoji(Severity) + "*" + sevText + "*" + GetEmoji(Severity) + "\n\n"
		count++
	} else {

		sevText := Severity + "_Total Firing - " + "_ "
		sevPreText = GetEmoji(Severity) + "*" + sevText + "*" + GetEmoji(Severity) + "\n\n"
	}

	for index, alert := range alertList {
		formatAlerted := "<" + alertLink[alert] + "|" + alert + ">"
		affectedHosts := getAffectedhosts(alertAffectedHosts, alert)
		if index == 0 {
			sendAlertDetails(sevPreText, formatAlerted, affectedHosts)
		} else {
			sendAlertDetails("", formatAlerted, affectedHosts)
		}

	}
}

func sendAlertDetails(preText string, formatAlerted string, affectedHosts string) {
	attachment := slack.Attachment{
		Color:      controllers.SlackColor,
		Pretext:    preText,
		Text:       formatAlerted,
		MarkdownIn: []string{"text", "pretext", "fields"},
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Affected Objects:",
				Value: affectedHosts,
				Short: false,
			},
		},
	}

	params := slack.MsgOptionAttachments(attachment)
	SlackSendMessage(params)
}
