package main

import (
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/price-check/amazon/amazon"
	"github.com/ynori7/price-check/amazon/config"
	"github.com/ynori7/price-check/emailer"
)

const (
	baseUrl = "https://www.amazon.de/gp/product/"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	logger := log.WithFields(log.Fields{"Logger": "main"})

	config.ParseCliFlags()
	if config.CliConf.ConfigFile == "" {
		log.Fatal("You must specify the path to the config file")
	}

	//Get the config
	data, err := ioutil.ReadFile(config.CliConf.ConfigFile)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error reading config file")
	}

	var conf config.Config
	if err := conf.Parse(data); err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error parsing config")
	}

	//Check the prices
	results := amazon.CheckPrices(conf.PriceConfig)
	logger.WithFields(log.Fields{"Count": len(results)}).Info("Results found")
	if len(results) == 0 {
		return
	}

	//Build the email report
	body := strings.Builder{}
	for _, res := range results {
		if res.Error == nil {
			body.WriteString(res.Result)
			body.WriteString("\n\n")
		} else {
			body.WriteString("ASIN: " + res.ASIN)
			body.WriteString(" failed: ")
			body.WriteString(err.Error())
			body.WriteString("\n\n")
		}
	}

	//Send email
	if conf.Email.Enabled {
		mailer := emailer.NewMailer(conf.Email)
		if err := mailer.SendMail("Amazon price check results", body.String()); err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("Error sending email")
		}
	}
}
