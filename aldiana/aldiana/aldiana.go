package aldiana

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

type Result struct {
	ASIN   string
	Result string
	Error  error
}

var client = http.Client{}

func CheckPrice(url string, minPrice float64) (string, error) {
	logger := log.WithFields(log.Fields{"Logger": "amazon"})

	req, _ := http.NewRequest(http.MethodGet, url, nil)

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

	priceRaw, _ := doc.Find(".offerSummary--container .price--amount--integer").Html()
	price, err := strconv.ParseFloat(strings.ReplaceAll(priceRaw, ".", ""), 64)
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
