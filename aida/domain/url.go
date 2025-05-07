package domain

import (
	"fmt"
	"net/url"
)

func BuildURL(
	adults, children, juveniles, babies int,
	durations []string,
	start, end string,
	port string,
	airport string,
) string {
	params := url.Values{
		"pax[adults]":    {fmt.Sprintf("%d", adults)},
		"pax[children]":  {fmt.Sprintf("%d", children)},
		"pax[juveniles]": {fmt.Sprintf("%d", juveniles)},
		"pax[babies]":    {fmt.Sprintf("%d", babies)},
		"duration":       durations,
		"from":           {start},
		"to":             {end},
		"sortCriteria":   {"Price"},
		"sortDirection":  {"Asc"},
		"size":           {"40"},
	}

	if port != "" {
		params.Add("port", port)
	}

	if airport != "" {
		params.Add("airport", airport)
	}

	return "https://aida.de/content/aida-search-and-booking/requests/search.cruise.json?" + params.Encode()
}

func BuildDetailsURL(journeyIdentifier, departureAirport string, tripSpec TripSpec) string {
	params := url.Values{
		"pax[adults]":    {fmt.Sprintf("%d", tripSpec.Adults)},
		"pax[children]":  {fmt.Sprintf("%d", tripSpec.Children)},
		"pax[juveniles]": {fmt.Sprintf("%d", tripSpec.Juveniles)},
		"pax[babies]":    {fmt.Sprintf("%d", tripSpec.Babies)},
	}

	if departureAirport != "" {
		params.Add("departureAirport", departureAirport)
	}

	return fmt.Sprintf("https://aida.de/finden/%s/VARIO?", journeyIdentifier) + params.Encode()
}
