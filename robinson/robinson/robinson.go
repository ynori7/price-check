package robinson

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/price-check/robinson/config"
)

var client = http.Client{}

type Result struct {
	Name  string
	URL   string
	Price string
}

func CheckPrice(conf config.PriceConfig) ([]Result, error) {
	logger := log.WithFields(log.Fields{"Logger": "robinson"})

	offerList, err := GetOffers(conf)
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0)
	for _, offer := range offerList.Items {
		offerDetails, err := LookupOffer(offer)
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Warn("Error looking up offer")
			continue
		}

		if offerDetails.Price <= conf.MinPrice {
			results = append(results, Result{
				Name:  offerDetails.Name,
				URL:   "https://www.robinson.com/" + offer.Offer.Button.Href,
				Price: fmt.Sprintf("%.2f", offerDetails.Price),
			})
		} else {
			logger.WithFields(log.Fields{"title": offerDetails.Name, "price": offerDetails.Price}).Info("Price is too high")
		}
	}

	return results, nil
}

func GetOffers(conf config.PriceConfig) (*OfferList, error) {
	logger := log.WithFields(log.Fields{"Logger": "robinson-offer-list"})

	params := url.Values{}
	params.Add("adults", strconv.Itoa(conf.Adults))
	params.Add("children", strconv.Itoa(conf.ChildAge))
	params.Add("departureDate", conf.StartDate)
	params.Add("returnDate", conf.EndDate)
	params.Add("duration", strconv.Itoa(conf.Duration))
	params.Add("hotelOnly", "1")
	params.Add("productCodes", conf.ProductCodes)
	params.Add("query", "")
	params.Add("rooms[0][adults]", strconv.Itoa(conf.Adults))
	params.Add("rooms[0][children]", strconv.Itoa(conf.ChildAge))

	req, _ := http.NewRequest(http.MethodPost, buildOfferListURL(), strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error making http request")
		return nil, err
	}
	defer resp.Body.Close()
	offerList := new(OfferList)
	if err := json.NewDecoder(resp.Body).Decode(&offerList); err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error unmarshalling json")
		return nil, err
	}

	return offerList, nil
}

func LookupOffer(offer Items) (*OfferDetails, error) {
	logger := log.WithFields(log.Fields{"Logger": "robinson-offer-details"})

	req, _ := http.NewRequest(http.MethodGet, buildOfferDetailsURL(offer.Offer.Button.Href), nil)
	resp, err := client.Do(req)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error making http request")
		return nil, err
	}
	defer resp.Body.Close()
	offerDetailResp := new(OfferDetailsResponse)
	if err := json.NewDecoder(resp.Body).Decode(&offerDetailResp); err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error unmarshalling json")
		return nil, err
	}

	if len(offerDetailResp.Offers) == 0 {
		logger.Warn("No offers found")
		return nil, fmt.Errorf("no offers found")
	}

	offerDetails := offerDetailResp.Offers[0]
	offerDetails.Name = offer.Club.Headline

	return &offerDetails, nil
}

func buildOfferListURL() string {
	return "https://www.robinson.com/de/de/angebot?tx_pituimodules_seljson%5Baction%5D=matches&tx_pituimodules_seljson%5Bcontroller%5D=SelJson&type=8050&cHash=b4da64d5e24b018a389d49f2d071a30f"
}

func buildOfferDetailsURL(buttonLink string) string {
	//offer/button/href
	strings.ReplaceAll(buttonLink, "/de/de/dispatch/?", "https://www.robinson.com/api/v1/offer/tcp/?hotelOnly=1&")
	return strings.ReplaceAll(buttonLink, "/de/de/dispatch/?", "https://www.robinson.com/api/v1/offer/tcp/?hotelOnly=1&")
}
