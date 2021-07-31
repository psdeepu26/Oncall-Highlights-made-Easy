package slackshop

func GetEmoji(emojiName string) string {
	if emojiName == "SEVERE" {
		return " :fire:"
	} else if emojiName == "WARN" {
		return " :warning:"
	} else if emojiName == "SMOKE" {
		return " :mag:"
	} else if emojiName == "pd" {
		return " :pd:"
	} else if emojiName == "jira" {
		return " :jira:"
	} else if emojiName == "jiraepic" {
		return " :jira_epic:"
	} else if emojiName == "inprogress" {
		return " :in_progress-5659:"
	} else if emojiName == "todo" {
		return " :todo:"
	} else if emojiName == "oncall" {
		return " :notify_support_oncall:"
	} else if emojiName == "wavefront" {
		return " :wavefront:"
	} else {
		return " :information_source:"
	}
}
