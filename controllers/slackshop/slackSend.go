package slackshop

import (
	"log"
	"os"

	"github.com/slack-go/slack"
)

func SlackApi(token string) *slack.Client {
	return slack.New(os.Getenv(token))
}

func SlackSendMessage(msg slack.MsgOption) {
	api := SlackApi("SLACK_BOT_TOKEN")
	channelId := os.Getenv("SLACK_CHANNEL_ID")

	_, _, _, err := api.SendMessage(
		channelId,
		msg,
	)

	if err != nil {
		log.Printf("Error while connection to Slack: %s\n", err)
	}
}
