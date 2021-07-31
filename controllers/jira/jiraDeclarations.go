package jira

type jira_epic_details struct {
	Expand     string `json:"expand"`
	Startat    int    `json:"startAt"`
	Maxresults int    `json:"maxResults"`
	Total      int    `json:"total"`
	Issues     []struct {
		Expand string `json:"expand"`
		ID     string `json:"id"`
		Self   string `json:"self"`
		Key    string `json:"key"`
		Fields struct {
			Summary  string `json:"summary"`
			Assignee struct {
				Self         string `json:"self"`
				Accountid    string `json:"accountId"`
				Emailaddress string `json:"emailAddress"`
				Avatarurls   struct {
					Four8X48  string `json:"48x48"`
					Two4X24   string `json:"24x24"`
					One6X16   string `json:"16x16"`
					Three2X32 string `json:"32x32"`
				} `json:"avatarUrls"`
				Displayname string `json:"displayName"`
				Active      bool   `json:"active"`
				Timezone    string `json:"timeZone"`
				Accounttype string `json:"accountType"`
			} `json:"assignee"`
			Status struct {
				Self           string `json:"self"`
				Description    string `json:"description"`
				Iconurl        string `json:"iconUrl"`
				Name           string `json:"name"`
				ID             string `json:"id"`
				Statuscategory struct {
					Self      string `json:"self"`
					ID        int    `json:"id"`
					Key       string `json:"key"`
					Colorname string `json:"colorName"`
					Name      string `json:"name"`
				} `json:"statusCategory"`
			} `json:"status"`
		} `json:"fields"`
	} `json:"issues"`
}
