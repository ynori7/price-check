package amazon

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/ynori7/price-check/amazon/config"
)

const (
	baseUrl = "https://www.amazon.de/gp/product/"
)

type Result struct {
	ASIN   string
	Result string
	Error  error
}

var client = http.Client{}

func CheckPrices(priceConfig []config.PriceConf) []Result {
	logger := log.WithFields(log.Fields{"Logger": "amazon"})

	results := []Result{}

	for _, priceGroup := range priceConfig {
		for _, asin := range priceGroup.ASINs {
			result := check(logger, asin, priceGroup.MinPrice)
			if result != nil {
				results = append(results, *result)
			}
		}

	}

	return results
}

func check(logger *log.Entry, asin string, minPrice float64) *Result {
	result := &Result{ASIN: asin}
	req, _ := http.NewRequest(http.MethodGet, baseUrl+asin, nil)

	resp, err := client.Do(req)
	if err != nil {
		logger.WithFields(log.Fields{"error": err, "asin": asin}).Warn("Error making http request")
		result.Error = err
		return result
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.WithFields(log.Fields{"error": err, "asin": asin}).Warn("Error parsing document")
		result.Error = err
		return result
	}

	priceRaw, _ := doc.Find("[name=displayedPrice]").Attr("value")
	price, err := strconv.ParseFloat(priceRaw, 64)
	if err != nil {
		logger.WithFields(log.Fields{"error": err, "asin": asin}).Warn("Error parsing price")
		result.Error = err
		return result
	}
	title := strings.TrimSpace(doc.Find("#productTitle").Text())
	kindleUnlimitedIcon := doc.Find(".a-icon-kindle-unlimited")
	isKindleUnlimited := kindleUnlimitedIcon != nil && kindleUnlimitedIcon.Nodes != nil

	if price <= minPrice || isKindleUnlimited {
		result.Result = title + " is "
		if isKindleUnlimited {
			result.Result += "available on Kindle Unlimited"
		} else {
			result.Result += fmt.Sprintf("available for %.2f. Threshold price was: %.2f", price, minPrice)
		}
		return result
	}

	return nil
}
