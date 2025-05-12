package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/price-check/aida/aida"
	"github.com/ynori7/price-check/aida/config"
	"github.com/ynori7/price-check/aida/domain"
	"github.com/ynori7/price-check/emailer"
)

type Success struct {
	Name      string
	Threshold float64
	Price     string
	URL       string
}

func main() {
	date := time.Now()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	logger := log.WithFields(log.Fields{"Logger": "main"})

	config.ParseCliFlags()
	if config.CliConf.ConfigFile == "" {
		log.Fatal("You must specify the path to the config file")
	}
	logger.Info("Starting Aida price check...")

	//Get the config
	data, err := os.ReadFile(config.CliConf.ConfigFile)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error reading config file")
	}

	var conf config.Config
	if err := conf.Parse(data); err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error parsing config")
	}
	logger = logger.WithFields(log.Fields{"config_name": conf.ConfigName})

	if !conf.Debug {
		log.SetLevel(log.InfoLevel)
	}

	tripSpecs := domain.BuildTripSpecifications(conf.PriceConfig)
	logger.WithFields(log.Fields{"count": len(tripSpecs)}).Info("Crawling trip specifications...")

	//Check the prices
	allResults := make([]domain.Result, 0)
	successes := make([]Success, 0)
	for _, trip := range tripSpecs {
		results, err := aida.CheckPriceOverview(trip)
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Info("Error checking price")
		} else {
			for _, result := range results {
				allResults = append(allResults, result)
				if result.IsPreferredBelowThreshold {
					successes = append(successes, Success{
						Name:      result.Name,
						Threshold: trip.DayPriceThreshold,
						Price:     result.TotalBasePrice,
						URL:       result.URL,
					})
				}
			}
		}
	}

	// Write the results to a JSON file for reporting later
	resultsJson, err := json.Marshal(allResults)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error marshalling results to JSON")
	}

	outDir := fmt.Sprintf("%s/%s", conf.OutputDirectory, date.Format("2006-01-02"))
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.Mkdir(outDir, os.ModePerm); err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("Error creating output directory")
			return
		}
	}

	fileName := fmt.Sprintf("%s/%s-%s_%s.json", outDir, conf.ConfigName, date.Format("2006-01-02T15-04"), config.CliConf.ScanName)
	if err := os.WriteFile(fileName, resultsJson, 0644); err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error writing results to file")
	}

	if len(successes) == 0 {
		return //nothing more to do
	}

	//Build the email report
	body := ""
	for _, s := range successes {
		body += fmt.Sprintf("The price for Aida \"%s\" has fallen below the threshold of %.2f. Current price: %s\n\n%s\n---\n", s.Name, s.Threshold, s.Price, s.URL)
		logger.Info(body)
	}

	//Send email
	if conf.Email.Enabled {
		mailer := emailer.NewMailer(conf.Email)
		if err := mailer.SendMail("Aida price check results", body); err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("Error sending email")
		}
	}
}
