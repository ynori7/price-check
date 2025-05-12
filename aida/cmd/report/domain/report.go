package domain

import "github.com/ynori7/price-check/aida/domain"

type Report struct {
	Overall     MinAndAvg
	TripReports []TripReport
	Durations   []MinAndAvg
	ScanPeriods []MinAndAvg
}

type TripReport struct {
	Name     string
	MinPrice float64
	MaxPrice float64
}

type MinAndAvg struct {
	GroupName     string
	CheapestOffer domain.Result

	CheapestPrice float64
	AveragePrice  float64
}
