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

var client *http.Client
var reqAnonymizer anonymizer.Anonymizer

func init() {
	client = hulkhttp.NewClient()
	reqAnonymizer = anonymizer.New(int64(rand.Int()))
}

func CheckPriceOverview(tripSpec domain.TripSpec) ([]domain.Result, error) {
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

	results := make([]domain.Result, 0)

	for _, cruiseItem := range cruiseResp.CruiseItems {
		if len(cruiseItem.CruiseItemVariants) == 0 {
			logger.WithFields(log.Fields{"title": cruiseItem.Title}).Debug("No cruise item variants found")
			continue
		}

		pricePerDay := cruiseItem.CruiseItemVariants[0].Amount / float64(cruiseItem.Duration)

		res := domain.Result{
			Name: cruiseItem.Title,
			URL: domain.BuildDetailsURL(
				cruiseItem.CruiseItemVariants[0].JourneyIdentifier,
				cruiseItem.CruiseItemVariants[0].DepartureAirport,
				tripSpec,
			),
			Region:     "", //TODO: where to find this?
			Duration:   fmt.Sprintf("%d days", cruiseItem.Duration),
			Port:       cruiseItem.CruiseItemVariants[0].FromCity,
			WithFlight: cruiseItem.FlightIncluded,
			BaseDayPrice:   fmt.Sprintf("%.02f", pricePerDay),
			TotalBasePrice: fmt.Sprintf("%.02f", cruiseItem.CruiseItemVariants[0].Amount),
			
		}

		logger.WithFields(log.Fields{"title": cruiseItem.Title, "duration": cruiseItem.Duration, "price_per_day": pricePerDay}).Debug("Checking cruise item")

		if pricePerDay <= priceThreshold {
			res.IsBaseBelowThreshold = true

			// if it's cheap enough, also check if the preferred cabin is available and cheap enough
			if preferredCabinCheapEnough, err := CheckIfPreferredCabinCheapEnough(
				cruiseItem.CruiseItemVariants[0].JourneyIdentifier,
				cruiseItem.CruiseItemVariants[0].TariffType,
				cruiseItem.Duration,
				tripSpec,
			); err == nil && preferredCabinCheapEnough {
				res.IsPreferredBelowThreshold = true
			}
		}
		
		results = append(results, res)
	}

	return results, nil
}

func CheckIfPreferredCabinCheapEnough(journeyIdentifier, tariffType string, duration int, tripSpec domain.TripSpec) (bool, error) {
	url := domain.BuildDetailsAPIUrl(journeyIdentifier, tariffType, tripSpec)
	logger := log.WithFields(log.Fields{"Logger": "aida"})
	logger.WithFields(log.Fields{"url": url}).Info("Getting additional price details")

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	reqAnonymizer.AnonymizeRequest(req)

	resp, err := client.Do(req)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error making http request")
		return false, err
	}
	defer resp.Body.Close()

	priceDetailsResp := new(PriceDetailsResponse)
	if err = json.NewDecoder(resp.Body).Decode(priceDetailsResp); err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error reading response body")
		return false, err
	}

	for _, cabinItem := range priceDetailsResp.CabinItemsVariants {
		if cabinItem.CabinCode == tripSpec.PreferredCabinType {
			pricePerDay := cabinItem.CabinPriceDetails.Amount / float64(duration)
			if pricePerDay <= tripSpec.DayPriceThreshold {
				logger.WithFields(log.Fields{"cabin_code": cabinItem.CabinCode, "price_per_day": pricePerDay}).Info("Preferred cabin price is below threshold")
				return true, nil
			}
			logger.WithFields(log.Fields{"cabin_code": cabinItem.CabinCode, "price_per_day": pricePerDay}).Debug("Preferred cabin price is above threshold")
			return false, nil
		}
	}

	logger.Info("Preferred cabin not available")
	return false, nil
}
