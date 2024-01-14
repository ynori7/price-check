package aida

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/ynori7/hulksmash/anonymizer"
	hulkhttp "github.com/ynori7/hulksmash/http"
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

func CheckPrice(url string, minPrice float64) ([]Result, error) {
	logger := log.WithFields(log.Fields{"Logger": "aida"})

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	reqAnonymizer.AnonymizeRequest(req)

	resp, err := client.Do(req)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error making http request")
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warn("Error parsing document")
		return nil, err
	}

	results := make([]Result, 0)

	doc.Find(".new-card-body").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Find(".route-teaser-title").Html()
		title = strings.ReplaceAll(title, "&amp;", "&")

		url, _ := s.Find(".new-routeteaser-action a.btn").Attr("href")

		priceRaw, _ := s.Find(".new-routeteaser-price-amount").Html()
		priceRaw = strings.ReplaceAll(strings.TrimSpace(priceRaw), ".", "")
		price, err := strconv.ParseFloat(priceRaw, 64)
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Warn("Error parsing price")
			return
		}

		priceString := fmt.Sprintf("%.0f", price)

		if price > minPrice {
			logger.WithFields(log.Fields{"title": title, "price": priceString}).Info("Price is too high")
			return
		}

		results = append(results, Result{
			Name:  title,
			URL:   "https://www.aida.de/" + url,
			Price: priceString,
		})
	})

	return results, nil
}
