package aida

type CruiseResponse struct {
	CruiseItems []CruiseItem `json:"cruiseItems"`
}

type CruiseItem struct {
	Title              string              `json:"title"`
	Duration           int                 `json:"duration"`
	StartDate          string              `json:"startDate"`
	EndDate            string              `json:"endDate"`
	FlightIncluded     bool                `json:"flightIncluded"`
	CruiseItemVariants []CruiseItemVariant `json:"cruiseItemVariant"`
}

type CruiseItemVariant struct {
	FromCity          string  `json:"fromCity"`
	ToCity            string  `json:"toCity"`
	Amount            float64 `json:"amount"`
	TariffType        string  `json:"tariffType"`
	JourneyIdentifier string  `json:"journeyIdentifier"`
	DepartureAirport  string  `json:"departureAirport"`
}

type PriceDetailsResponse struct {
	CabinItemsVariants []CabinItemsVariant `json:"cabinItemsVariant"`
}

type CabinItemsVariant struct {
	CabinCode         string            `json:"cabinCode"`
	CabinPriceDetails CabinPriceDetails `json:"ind"`
}

type CabinPriceDetails struct {
	Amount float64 `json:"amount"`
}
