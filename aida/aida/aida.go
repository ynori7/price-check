package aida

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/ynori7/hulksmash/anonymizer"
	hulkhttp "github.com/ynori7/hulksmash/http"
)

var client *http.Client
var reqAnonymizer anonymizer.Anonymizer

func init() {
	client = hulkhttp.NewClient()
	reqAnonymizer = anonymizer.New(int64(rand.Int()))
}

func CheckPrice(url string, minPrice float64) (string, error) {
	logger := log.WithFields(log.Fields{"Logger": "aida"})

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	reqAnonymizer.AnonymizeRequest(req)

	resp, err := client.Do(req)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error making http request")
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error parsing document")
		return "", err
	}

	priceRaw, _ := doc.Find("#cruisedetail .route-offers__cabin-price").Attr("price")
	price, err := strconv.ParseFloat(priceRaw, 64)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error parsing price")
		return "", err
	}

	priceString := fmt.Sprintf("%.0f", price)

	if price <= minPrice {
		return priceString, nil
	}

	return "", fmt.Errorf("price too high. current price:" + priceString)
}
