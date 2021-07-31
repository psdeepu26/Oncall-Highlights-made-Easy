package slackshop

import (
	"fmt"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/spatrayuni/tobs-oncall-highlights/controllers"
	"github.com/spatrayuni/tobs-oncall-highlights/controllers/pd"
	"github.com/spatrayuni/tobs-oncall-highlights/controllers/wf"

	log "github.com/sirupsen/logrus"
)

// Ongoing Highlights
func SendOngoingIncidentDetails() {
	incidentSummary, incidentPinLink := GetOnGoingIncidents()
	count := 0
	preText := ""

	for incident, summary := range incidentSummary {
		issueSummary := strings.Replace(summary, "\n", "", -1)
		Text := "*Incident No.: <" + incidentPinLink[incident] + "|" + incident + ">*"
		pinnedValue := "<" + incidentPinLink[incident] + "|Last Pinned Update>"

		if count == 0 {
			preText = GetEmoji("oncall") + " *Ongoing Incidents* " + GetEmoji("oncall") + "\n"
			onGoingCount := fmt.Sprintf("Currently Ongoing Incidents : %d \n", len(incidentSummary))
			preText = preText + onGoingCount

			OnGoingIncSlackAttachment(preText, Text, incident, issueSummary, pinnedValue)
			count++
		} else {
			preText = ""
			OnGoingIncSlackAttachment(preText, Text, incident, issueSummary, pinnedValue)
			count++
		}
	}

}

func OnGoingIncSlackAttachment(preText string, Text string, incident string, issueSummary string, pinnedValue string) {
	var fields []slack.AttachmentField

	message1 := slack.AttachmentField{
		Title: "Incident Title",
		Value: issueSummary,
		Short: false,
	}
	fields = append(fields, message1)

	message2 := slack.AttachmentField{
		Title: "Update",
		Value: pinnedValue,
		Short: false,
	}
	fields = append(fields, message2)

	attachment := slack.Attachment{
		Color:      controllers.SlackColor,
		Pretext:    preText,
		Text:       Text,
		MarkdownIn: []string{"text", "pretext", "fields"},
		Fields:     fields,
	}

	params := slack.MsgOptionAttachments(attachment)
	SlackSendMessage(params)
}

//PD Highlights
func SendPDDetails() {
	preText := GetEmoji("pd") + " *PagerDuty Details :* " + GetEmoji("pd") + "\n\n"
	currentOncall, toBeOncall, currentOncallDuration, toBeOncallDuration := pd.GetPDOncallSchedule()
	openInc, _ := pd.GetPDIncidents()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	log.Info("PagerDuty Details: \n")
	log.Info("Current Oncall Engineer : %s and Duration is till: %s\n", currentOncall, currentOncallDuration)
	log.Info("Next Oncall Engineer : %s and Duration is till: %s\n", toBeOncall, toBeOncallDuration)

	attachment := slack.Attachment{
		Color:      controllers.SlackColor,
		Pretext:    preText,
		MarkdownIn: []string{"text", "pretext", "fields"},
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Current Oncall Engineer",
				Value: currentOncall + " - (" + currentOncallDuration + ")",
				Short: true,
			},
			slack.AttachmentField{
				Title: "Next Oncall Engineer",
				Value: toBeOncall + " - (" + toBeOncallDuration + ")",
				Short: true,
			},
			slack.AttachmentField{
				Title: "Currently Firing PD Incidents:",
				Value: fmt.Sprintf("%d", openInc),
				Short: true,
			},
		},
	}
	params := slack.MsgOptionAttachments(attachment)
	SlackSendMessage(params)
}

// Alerts Highlights
func SendOncallHighlights() {
	aggregatedAlerts, alertAffectedHosts, alertLink := wf.AggregateAlerts()
	count := 0

	if len(aggregatedAlerts) == 0 {
		msg := formatNoAlertsMessage()
		SlackSendMessage(msg)
	}

	for Severity, alertList := range aggregatedAlerts {
		if Severity == "SEVERE" {
			formatAlertAttachment(alertList, Severity, alertAffectedHosts, alertLink, count)
			count++
		} else if Severity == "WARN" {
			formatAlertAttachment(alertList, Severity, alertAffectedHosts, alertLink, count)
			count++
		} else if Severity == "SMOKE" {
			formatAlertAttachment(alertList, Severity, alertAffectedHosts, alertLink, count)
			count++
		} else if Severity == "INFO" {
			formatAlertAttachment(alertList, Severity, alertAffectedHosts, alertLink, count)
			count++
		}
	}
}

// Title
func TitleSection() slack.MsgOption {
	preText := GetEmoji("wavefront") + "  *Oncall Highlights at " + controllers.GetCurrentDateTimeInPST() + "*  " + GetEmoji("wavefront")
	preTextField := slack.NewTextBlockObject("mrkdwn", preText+"\n\n", false, false)
	preTextSection := slack.NewSectionBlock(preTextField, nil, nil)
	dividerSection := slack.NewDividerBlock()

	msg := slack.MsgOptionBlocks(
		preTextSection,
		dividerSection,
	)

	return msg
}

// Divider
func Divider() slack.MsgOption {
	dividerSection1 := slack.NewDividerBlock()

	msg := slack.MsgOptionBlocks(
		dividerSection1,
	)
	return msg
}
