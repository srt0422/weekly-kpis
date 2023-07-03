package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/jszwec/csvutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	analytics "google.golang.org/api/analytics/v3"
)

type Record struct {
	Who string `csv:"Who"`
	WeeklyKPIs string `csv:"Weekly KPIs"`
	Date string `csv:",omitempty"` // The field name should be the date
}

func main() {
	// read the service account key file
	key, err := ioutil.ReadFile("path/to/your/service-account.json")
	if err != nil {
		log.Fatal(err)
	}

	// get the JWT config
	conf, err := google.JWTConfigFromJSON(key, analytics.AnalyticsReadonlyScope)
	if err != nil {
		log.Fatal(err)
	}

	// create a new OAuth2 client
	client := conf.Client(oauth2.NoContext)

	// create a new Analytics service
	svc, err := analytics.New(client)
	if err != nil {
		log.Fatal(err)
	}

	// perform the API request
	data, err := svc.Data.Ga.Get("ga:YOUR_VIEW_ID", "7daysAgo", "today", "ga:sessions,ga:users").Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: process the data object to extract the values you want
	alexa_total_downloads_to_date := data.TotalsForAllResults["ga:sessions"]

	// open the CSV file for appending
	f, err := os.OpenFile("path/to/your.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// create a new CSV encoder
	enc := csvutil.NewEncoder(f)

	// write a new record to the CSV file
	err = enc.Encode(Record{
		Who: "Alexa",
		WeeklyKPIs: "Total Downloads to Date",
		Date: alexa_total_downloads_to_date,  // The field name should be the date
	})
	if err != nil {
		log.Fatal(err)
	}
}
