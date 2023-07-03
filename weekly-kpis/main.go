package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	analytics "google.golang.org/api/analytics/v3"
)

func main() {
	err := godotenv.Load() // load .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gaViewID := os.Getenv("GA_VIEW_ID") // get the view id from .env
	csvFile := os.Getenv("CSV_FILE")

	key, err := ioutil.ReadFile(os.Getenv("GA_SERVICE_ACCOUNT"))
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(key, analytics.AnalyticsReadonlyScope)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(oauth2.NoContext)

	svc, err := analytics.New(client)
	if err != nil {
		log.Fatal(err)
	}

	data, err := svc.Data.Ga.Get("ga:"+gaViewID, "7daysAgo", "today", "ga:sessions,ga:users").Do()
	if err != nil {
		log.Fatal(err)
	}
fmt.Printf("data: %v\n", data)
return
	// TODO: process the data object to extract the values you want
	// example:
	kpis := process(data)

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		log.Fatal(err)
	}

	// ... same code as before ...

	// Assuming you have date in "mm/dd" format
	lines = buildRecords(kpis, "06/29", lines)

	// TODO: Add other KPIs...

	// Write back to file
	file, err = os.Create("path/to/your.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	writer.Flush()
}

func process(data *analytics.GaData) map[string]string {  
	total_downloads_to_date := "0"  
	visits_to_download_page := "0"  
	download_clicks := "0"  

	// Iterate through rows to extract metrics  
	for _, row := range data.Rows {  
		if row.Metrics[0].Value != nil {  
			total_downloads_to_date = *row.Metrics[0].Value  
		}  
		if row.Metrics[1].Value != nil {  
			visits_to_download_page = *row.Metrics[1].Value  
		}  
		if row.Metrics[2].Value != nil {  
			download_clicks = *row.Metrics[2].Value  
		}  
	}  

	return map[string]string{  
		"Total Downloads to Date": total_downloads_to_date,  
		"Visits to Download Page": visits_to_download_page,  
		"Download Clicks":         download_clicks,  
	}  
} 



func buildRecords(kpis map[string]string, date string, lines []string) []string {
	for kpi, value := range kpis {
		lines = updateRecord(kpi, date, value, lines)
	}
	return lines
}

// This function updates each record with the data provided.
func updateRecord(kpi, date, value string, lines []string) []string {
	for i, line := range lines {
		// split the line into columns
		columns := strings.Split(line, ",")

		// if the line's second column (index 1) matches the KPI, update the line
		if columns[1] == kpi {
			// loop over columns to find the date column and update value
			for j, col := range columns {
				if col == date {
					columns[j] = value
					break
				}
			}

			// join columns back into a line
			lines[i] = strings.Join(columns, ",")
			break
		}
	}
	return lines
}
