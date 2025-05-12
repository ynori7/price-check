package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ynori7/price-check/aida/cmd/report/display"
	"github.com/ynori7/price-check/aida/cmd/report/domain"
	"github.com/ynori7/price-check/aida/config"
	aidadomain "github.com/ynori7/price-check/aida/domain"
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
		logger.Fatal("You must specify the path to the config file")
	}

	//Get the config
	data, err := os.ReadFile(config.CliConf.ConfigFile)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error reading config file")
	}

	var conf config.Config
	if err := conf.Parse(data); err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error parsing config")
	}

	if !conf.Debug {
		log.SetLevel(log.InfoLevel)
	}

	// open all the files for today's date in the output directory
	// unmarshal it into a list of domain.Result
	outputDir := fmt.Sprintf("%s/%s", conf.OutputDirectory, time.Now().Format("2006-01-02"))

	// Get a list of files in the output directory
	files, err := os.ReadDir(outputDir)
	if err != nil {
		log.Fatalf("Failed to read output directory: %v", err)
	}

	var results []aidadomain.Result

	// Iterate over the files and process those matching today's date
	for _, file := range files {
		filePath := filepath.Join(outputDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v", filePath, err)
			continue
		}

		var fileResults []aidadomain.Result
		if err := json.Unmarshal(data, &fileResults); err != nil {
			log.Printf("Failed to unmarshal file %s: %v", filePath, err)
			continue
		}

		for i := range fileResults {
			fileParts := strings.Split(file.Name(), "_")
			if len(fileParts) != 2 {
				log.Printf("Invalid file name format: %s", file.Name())
				continue
			}
			fileResults[i].TimePeriod = strings.TrimSuffix(fileParts[1], ".json")
			baseDayPrice, _ := strconv.ParseFloat(fileResults[i].BaseDayPrice, 64)
			fileResults[i].BaseDayPriceFloat = baseDayPrice
		}

		results = append(results, fileResults...)
	}

	// Log the results or process them further
	logger.Infof("Loaded %d results", len(results))

	report := domain.Report{}

	//cheapest overall and average overall
	cheapestPrice := 99999.0
	cheapestIndex := 0
	sum := 0.0
	for i, result := range results {
		if result.BaseDayPriceFloat < cheapestPrice {
			cheapestPrice = result.BaseDayPriceFloat
			cheapestIndex = i
		}
		sum += result.BaseDayPriceFloat
	}
	report.Overall.CheapestOffer = results[cheapestIndex]
	report.Overall.CheapestPrice = cheapestPrice
	report.Overall.AveragePrice = sum / float64(len(results))

	// range per trip
	tripsIndices := make(map[string]int)
	trips := make([]domain.TripReport, 0)
	i := 0
	for _, result := range results {
		if _, ok := tripsIndices[result.Name]; !ok {
			tripsIndices[result.Name] = i
			i++
			trips = append(trips, domain.TripReport{
				Name:     result.Name,
				MinPrice: result.BaseDayPriceFloat,
				MaxPrice: result.BaseDayPriceFloat,
			})
		} else {
			index := tripsIndices[result.Name]
			if result.BaseDayPriceFloat < trips[index].MinPrice {
				trips[index].MinPrice = result.BaseDayPriceFloat
			}
			if result.BaseDayPriceFloat > trips[index].MaxPrice {
				trips[index].MaxPrice = result.BaseDayPriceFloat
			}
		}
	}
	sort.Slice(trips, func(i, j int) bool {
		return trips[i].Name < trips[j].Name
	})
	report.TripReports = trips

	//cheapest per duration and average per duration
	report.Durations = getMinAndAvgForGroup(results, func(r aidadomain.Result) string {
		return r.Duration
	})
	sort.Slice(report.Durations, func(i, j int) bool {
		groupI, _ := strconv.Atoi(strings.TrimSpace(strings.Split(report.Durations[i].GroupName, " ")[0]))
		groupJ, _ := strconv.Atoi(strings.TrimSpace(strings.Split(report.Durations[j].GroupName, " ")[0]))
		return groupI < groupJ
	})

	//cheapest per scan period and average per scan period
	report.ScanPeriods = getMinAndAvgForGroup(results, func(r aidadomain.Result) string {
		return r.TimePeriod
	})
	sort.Slice(report.ScanPeriods, func(i, j int) bool {
		return report.ScanPeriods[i].GroupName < report.ScanPeriods[j].GroupName
	})

	template := display.NewHtmlTemplate(report)
	html, err := template.ExecuteHtmlTemplate()
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Fatal("Error executing template")
	}

	//Send email
	if conf.Email.Enabled {
		mailer := emailer.NewMailer(conf.Email)
		if err := mailer.SendHTMLMail("Aida price check report", html); err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("Error sending email")
		}
	} else {
		fmt.Printf("%s\n", html)
	}
}

func getMinAndAvgForGroup(results []aidadomain.Result, keyFunc func(aidadomain.Result) string) []domain.MinAndAvg {
	groupedByKey := groupBy(results, keyFunc)
	groupedResults := make([]domain.MinAndAvg, 0, len(groupedByKey))
	for _, group := range groupedByKey {
		cheapestPrice := 99999.0
		cheapestIndex := 0
		sum := 0.0
		for i, result := range group {
			if result.BaseDayPriceFloat < cheapestPrice {
				cheapestPrice = result.BaseDayPriceFloat
				cheapestIndex = i
			}
			sum += result.BaseDayPriceFloat
		}
		groupedResults = append(groupedResults, domain.MinAndAvg{
			GroupName:     keyFunc(group[0]),
			CheapestOffer: group[cheapestIndex],
			CheapestPrice: cheapestPrice,
			AveragePrice:  sum / float64(len(group)),
		})
	}

	return groupedResults
}

func groupBy(results []aidadomain.Result, keyFunc func(aidadomain.Result) string) map[string][]aidadomain.Result {
	grouped := make(map[string][]aidadomain.Result)
	for _, result := range results {
		key := keyFunc(result)
		grouped[key] = append(grouped[key], result)
	}
	return grouped
}
