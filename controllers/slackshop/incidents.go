package slackshop

import (
	"strings"

	"github.com/slack-go/slack"
	"github.com/spatrayuni/tobs-oncall-highlights/controllers"
)

func GetOnGoingIncidents() (map[string]string, map[string]string) {
	api := SlackApi("SLACK_BOT_TOKEN_LEGACY")
	var param1 slack.GetConversationsForUserParameters
	param1.UserID = controllers.UserID
	param1.ExcludeArchived = false

	var linkParameters slack.PermalinkParameters
	var link string
	incidentSummary := make(map[string]string)
	incidentPinLink := make(map[string]string)

	channelslist2, _, _ := api.GetConversationsForUser(&param1)

	for _, channel := range channelslist2 {
		if strings.Contains(channel.Name, "_incident-") && int64(channel.Created) >= controllers.GetDiffEpochTime(31) {

			items, _, _ := api.ListPins(channel.ID)
			for index, item := range items {
				if index == 0 {
					linkParameters.Channel = channel.ID
					linkParameters.Ts = item.Message.Msg.Timestamp
					link, _ = api.GetPermalink(&linkParameters)
				}
			}
			channelName := "#" + channel.Name
			incidentSummary[channelName] = channel.Topic.Value
			incidentPinLink[channelName] = link
		}
	}
	return incidentSummary, incidentPinLink
}
