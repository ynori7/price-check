package domain

import (
	"fmt"

	"github.com/ynori7/price-check/aida/config"
)

type TripSpec struct {
	URL               string //URL to the trip
	DayPriceThreshold float64

	//details about travelers. Used also to build the details URL
	Adults    int
	Children  int
	Juveniles int
	Babies    int

	Title string
}

func BuildTripSpecifications(conf config.PriceConfig) []TripSpec {
	fmt.Printf("%+v\n", conf)
	trips := make([]TripSpec, 0)
	for _, trip := range conf.Trips {
		// Build URLs for each trip starting directly from a port
		for _, port := range conf.Ports {
			trips = append(trips, TripSpec{
				URL: BuildURL(
					conf.Adults,
					conf.Children,
					conf.Juveniles,
					conf.Babies,
					conf.Durations,
					trip.Start,
					trip.End,
					port,
					"",
				),
				DayPriceThreshold: conf.WithoutFlightDayPriceThreshold,
				Adults:            conf.Adults,
				Children:          conf.Children,
				Juveniles:         conf.Juveniles,
				Babies:            conf.Babies,
				Title:             fmt.Sprintf("%s to %s from port %s", trip.Start, trip.End, port),
			})
		}

		// Build URLs for each trip including a flight from an airport
		for _, airport := range conf.Airports {
			trips = append(trips, TripSpec{
				URL: BuildURL(
					conf.Adults,
					conf.Children,
					conf.Juveniles,
					conf.Babies,
					conf.Durations,
					trip.Start,
					trip.End,
					"",
					airport,
				),
				DayPriceThreshold: conf.WithFlightDayPriceThreshold,
				Adults:            conf.Adults,
				Children:          conf.Children,
				Juveniles:         conf.Juveniles,
				Babies:            conf.Babies,
				Title:             fmt.Sprintf("%s to %s with flight from %s", trip.Start, trip.End, airport),
			})
		}
	}
	return trips
}
