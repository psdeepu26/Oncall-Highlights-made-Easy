package wf

import (
	"fmt"
	"os"

	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/spatrayuni/tobs-oncall-highlights/controllers"

	log "github.com/sirupsen/logrus"
)

func WavefrontConnection() *wavefront.Alerts {
	client, err := wavefront.NewClient(
		&wavefront.Config{
			Address: "mon.wavefront.com",
			Token:   os.Getenv("MON_API_TOKEN"),
		},
	)

	if err != nil {
		log.Fatal("Error during Wavefront connection :", err)
	}

	alertsClient := client.Alerts()
	fmt.Println("Wavefront Connection Successfull")

	return alertsClient
}

func FindAlerts(alertsClient *wavefront.Alerts) []*wavefront.Alert {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	searchCondition := []*wavefront.SearchCondition{
		{
			Key:            "tagpath",
			Value:          "group.ops",
			MatchingMethod: "TAGPATH",
		},
		{
			Key:            "status",
			Value:          "FIRING",
			MatchingMethod: "EXACT",
		},
	}

	wavefrontAlerts, err1 := alertsClient.Find(
		searchCondition,
	)

	if err1 != nil {
		log.Error("Error while finding alerts : ", err1)
	}

	return wavefrontAlerts
}

func CheckMap(alertAffectedHosts map[string][]string, AlertName string, Host string) bool {
	for _, item := range alertAffectedHosts[AlertName] {
		if item == Host {
			return true
		}
	}
	return false
}

func AggregateAlerts() (map[string][]string, map[string][]string, map[string]string) {
	alertsClient := WavefrontConnection()
	wavefrontAlerts := FindAlerts(alertsClient)

	aggregatedAlerts := make(map[string][]string)
	alertAffectedHosts := make(map[string][]string)
	alertLink := make(map[string]string)

	for _, eachAlert := range wavefrontAlerts {

		alertLink[eachAlert.Name] = controllers.MonURL + "/alert/" + *eachAlert.ID

		for _, host := range eachAlert.FailingHostLabelPairs {
			repeated := CheckMap(alertAffectedHosts, eachAlert.Name, host.Host)
			if repeated == false {
				alertAffectedHosts[eachAlert.Name] = append(alertAffectedHosts[eachAlert.Name], host.Host)
			}
		}

		aggregatedAlerts[eachAlert.Severity] = append(aggregatedAlerts[eachAlert.Severity], eachAlert.Name)
	}

	return aggregatedAlerts, alertAffectedHosts, alertLink
}

/*
func PrintAlerts(aggregatedAlerts map[string][]string, alertAffectedHosts map[string][]string, alertLink map[string]string) {
	for Severity, alertList := range aggregatedAlerts {
		output := fmt.Sprintf("\n\n************ %s _Total Firing - %v_ ************\n", Severity, len(alertList))
		log.Info(output)
		for _, alert := range alertList {
			log.Info("* %s \n", alert)
			log.Info("  Alert link: %s\n", alertLink[alert])
			log.Info("	Affected cluster/mirrors/hosts for this alert")
			for _, host := range alertAffectedHosts[alert] {
				log.Info("		%s\n", host)
			}
		}
	}
}*/
