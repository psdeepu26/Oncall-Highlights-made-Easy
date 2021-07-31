package pd

import (
	"os"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"
	"github.com/spatrayuni/tobs-oncall-highlights/controllers"
)

func pagerDutyClient() *pagerduty.Client {
	return pagerduty.NewClient(os.Getenv("PAGERDUTY_TOKEN"))
}

func GetPDOncallSchedule() (string, string, string, string) {
	var scheduleOpts pagerduty.GetScheduleOptions
	client := pagerDutyClient()

	// Building query options
	scheduleOpts.Since = time.Now().String()
	scheduleOpts.Until = time.Now().Add(time.Hour * 24).String()

	// API call
	sched, err := client.GetSchedule(controllers.PdScheduleId, scheduleOpts)
	if err != nil {
		log.Println(err)
	}
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	log.Info("PagerDuty Schedule collection is successfull!")

	currentOncall := sched.FinalSchedule.RenderedScheduleEntries[0].User.Summary
	toBeOncall := sched.FinalSchedule.RenderedScheduleEntries[1].User.Summary
	curentShiftEndTime := sched.FinalSchedule.RenderedScheduleEntries[0].End
	nextShiftEndTime := sched.FinalSchedule.RenderedScheduleEntries[1].End

	currentOncallDuration := controllers.GetDiffTime(controllers.GetParsedTime(curentShiftEndTime), time.Now())
	toBeOncallDuration := controllers.GetDiffTime(controllers.GetParsedTime(nextShiftEndTime), controllers.GetParsedTime(curentShiftEndTime))

	return currentOncall, toBeOncall, controllers.ShortDur(currentOncallDuration), controllers.ShortDur(toBeOncallDuration)
}

func GetPDIncidents() (int, map[string]string) {
	client := pagerDutyClient()
	var incidentOpts pagerduty.ListIncidentsOptions
	incidentList := make(map[string]string)

	incidentOpts.Statuses = []string{"triggered", "acknowledged"}
	incidentOpts.ServiceIDs = []string{controllers.PdServiceID}

	openIncidents, _ := client.ListIncidents(incidentOpts)
	openInc := len(openIncidents.Incidents)

	log.Printf("No. of Open Incidents : %d\n", openInc)

	if openInc > 0 {
		for _, incident := range openIncidents.Incidents {
			incidentLink := controllers.PdUrl + "/incidents/" + incident.Id
			incidentList[incident.Title] = incidentLink
		}
	}

	return openInc, incidentList
}
