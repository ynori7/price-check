package main

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/price-check/aldiana/aldiana"
	"github.com/ynori7/price-check/aldiana/config"
	"github.com/ynori7/price-check/emailer"
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
	successes := make([]string, 0)
	for _, priceConf := range conf.PriceConfig {
		price, err := aldiana.CheckPrice(priceConf.URL, priceConf.MinPrice)
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Info("Error checking price")
		} else {
			successes = append(successes, price+" \n\n"+priceConf.URL)
		}
	}

	if len(successes) == 0 {
		return //nothing more to do
	}

	//Build the email report
	body := ""
	for _, s := range successes {
		body += "The price for Aldiana has fallen below the threshold. Current price: " + s + "\n---\n"
		logger.Info(body)
	}

	//Send email
	if conf.Email.Enabled {
		mailer := emailer.NewMailer(conf.Email)
		if err := mailer.SendMail("Amazon price check results", body); err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("Error sending email")
		}
	}
}
