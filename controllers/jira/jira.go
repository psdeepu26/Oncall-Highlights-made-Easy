package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spatrayuni/tobs-oncall-highlights/controllers"
	"github.com/spatrayuni/tobs-oncall-highlights/controllers/slackshop"
)

func postRequest() *http.Response {
	var payload = strings.NewReader(`{
		"jql": "project = OPS And \"Epic Link\" = OPS-9551 And (status = \"IN PROGRESS\" OR status = \"BLOCKED\" OR status = \"TRIAGED\" OR status = \"NEXT\")",
		"fields": ["summary","assignee","status"]
	}`)

	var searchurl string = controllers.JiraUrl + "/rest/api/2/search"

	req, err := http.NewRequest("POST", searchurl, payload)
	var authorization string = "Basic " + os.Getenv("JIRA_TOKEN")

	if err != nil {
		log.Error("Unable to create new Http request: ", err)
	}

	log.Info("Http Post request connection successfull!")

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{authorization},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Unable to get Epic details of ", controllers.EpicId, ": ", err)
	}

	log.Info("Succesful in getting the Epic Status and Sumamry details for Epic: %s", controllers.EpicId)

	return resp
}

func convertResponsetoByte(resp *http.Response) []byte {
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Unable to get JIRA Status and Summary Details.\nError: ", err)
	}

	return out
}

func unmarshalResponse(response_body []byte) *jira_epic_details {
	var data = new(jira_epic_details)

	err := json.Unmarshal(response_body, &data)
	if err != nil {
		log.Error("Unable to unmarshal the request body.\nError: ", err)
	}
	return data
}

func GetJIRADetails(jira_details *jira_epic_details) (map[string]string, map[string]string, map[string]string) {
	inProgressJiraDetails := make(map[string]string)
	toDoJiraDetails := make(map[string]string)
	jiraAssignee := make(map[string]string)

	for _, issue := range jira_details.Issues {
		link := controllers.JiraUrl + "/browse/" + issue.Key

		if issue.Fields.Status.Statuscategory.Name == "In Progress" {
			inProgressJiraDetails[link] = issue.Fields.Summary
		} else if issue.Fields.Status.Statuscategory.Name == "To Do" {
			toDoJiraDetails[link] = issue.Fields.Summary
		}
		jiraAssignee[link] = issue.Fields.Assignee.Displayname

		log.Info("https://wavefront.atlassian.net/browse", issue.Key, " - ", issue.Fields.Status.Statuscategory.Name, " - ", issue.Fields.Summary, " - ", issue.Fields.Assignee.Displayname)
	}

	return inProgressJiraDetails, toDoJiraDetails, jiraAssignee
}

func SendJiraTitelDetails(InProgressJiraCount string, ToDoJiraCount string) {
	preText := slackshop.GetEmoji("jira") + "*Oncall Jira Issues - " + slackshop.GetEmoji("jiraepic") + controllers.JiraUrl + "/browse/" + controllers.EpicId + ":*" + slackshop.GetEmoji("jira") + "\n"

	attachment := slack.Attachment{
		Color:      controllers.SlackColor,
		Pretext:    preText,
		MarkdownIn: []string{"text", "pretext", "fields"},
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Current InProgress Jiras :",
				Value: InProgressJiraCount,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Current ToDo Jiras",
				Value: ToDoJiraCount,
				Short: true,
			},
		},
	}

	params := slack.MsgOptionAttachments(attachment)
	slackshop.SlackSendMessage(params)
}

func sendDetails(jiraText string, Assignee string, Text string) {
	attachment := slack.Attachment{
		Color:      controllers.SlackColor,
		MarkdownIn: []string{"text", "pretext", "fields"},
		Pretext:    Text,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Jira Link",
				Value: jiraText,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Jira Assignee",
				Value: Assignee,
				Short: true,
			},
		},
	}

	params := slack.MsgOptionAttachments(attachment)
	slackshop.SlackSendMessage(params)
}

func SendJiraDetails(jiraDetails map[string]string, jiraAssignee map[string]string, jiratype string) {
	Text := ""
	count := 0

	if len(jiraDetails) > 0 {
		for link, summary := range jiraDetails {
			issueSummary := strings.Replace(summary, "\n", "", -1)
			jiraText := "<" + link + "|" + issueSummary + ">" + "\n"
			if count == 0 {
				if jiratype == "In Progress" {
					Text = slackshop.GetEmoji("inprogress") + "*InProgress Jira Issues*" + slackshop.GetEmoji("inprogress") + "\n"
				} else if jiratype == "To Do" {
					Text = slackshop.GetEmoji("todo") + " *To Do Jira Issues * " + slackshop.GetEmoji("todo") + "\n"
				}
				sendDetails(jiraText, jiraAssignee[link], Text)
				count++
			} else {
				sendDetails(jiraText, jiraAssignee[link], "")
			}

		}
	}
}

func GetJiraEpicDetails() {
	response := postRequest()
	response_body := convertResponsetoByte(response)
	jira_details := unmarshalResponse(response_body)

	inProgressJiraDetails, toDoJiraDetails, jiraAssignee := GetJIRADetails(jira_details)

	SendJiraTitelDetails(fmt.Sprintf("%d", len(inProgressJiraDetails)), fmt.Sprintf("%d", len(toDoJiraDetails)))

	SendJiraDetails(inProgressJiraDetails, jiraAssignee, "In Progress")
	SendJiraDetails(toDoJiraDetails, jiraAssignee, "To Do")
}
