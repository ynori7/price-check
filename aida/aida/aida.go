package aida

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/hulksmash/anonymizer"
	hulkhttp "github.com/ynori7/hulksmash/http"
	"github.com/ynori7/price-check/aida/domain"
)

type Result struct {
	Name  string
	URL   string
	Price string
}

var client *http.Client
var reqAnonymizer anonymizer.Anonymizer

func init() {
	client = hulkhttp.NewClient()
	reqAnonymizer = anonymizer.New(int64(rand.Int()))
}

func CheckPrice(tripSpec domain.TripSpec) ([]Result, error) {
	url := tripSpec.URL
	priceThreshold := tripSpec.DayPriceThreshold

	logger := log.WithFields(log.Fields{"Logger": "aida"})
	logger.WithFields(log.Fields{"url": url}).Info("Scanning trip: " + tripSpec.Title)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	reqAnonymizer.AnonymizeRequest(req)

	resp, err := client.Do(req)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error making http request")
		return nil, err
	}
	defer resp.Body.Close()

	cruiseResp := new(CruiseResponse)
	err = json.NewDecoder(resp.Body).Decode(cruiseResp)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error decoding response")
		return nil, err
	}

	results := make([]Result, 0)

	for _, cruiseItem := range cruiseResp.CruiseItems {
		if len(cruiseItem.CruiseItemVariants) == 0 {
			logger.WithFields(log.Fields{"title": cruiseItem.Title}).Debug("No cruise item variants found")
			continue
		}

		pricePerDay := cruiseItem.CruiseItemVariants[0].Amount / float64(cruiseItem.Duration)

		logger.WithFields(log.Fields{"title": cruiseItem.Title, "duration": cruiseItem.Duration, "price_per_day": pricePerDay}).Debug("Checking cruise item")

		if pricePerDay <= priceThreshold {
			results = append(results, Result{
				Name: cruiseItem.Title,
				URL: domain.BuildDetailsURL(
					cruiseItem.CruiseItemVariants[0].JourneyIdentifier,
					cruiseItem.CruiseItemVariants[0].DepartureAirport,
					tripSpec,
				),
				Price: fmt.Sprintf("%.02f", cruiseItem.CruiseItemVariants[0].Amount),
			})
		}
	}

	return results, nil
}
