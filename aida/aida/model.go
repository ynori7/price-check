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

/*
{
  "resultsTotal": 19,
  "nextPage": 2,
  "totalPages": 1,
  "paging": false,
  "currency": "€",
  "success": true,
  "cruiseItems": [
    {
      "routeGroupCode": "Kurzreise nach Danzig & Bornholm",
      "routeCode": "WAR04096",
      "yieldRouteCode": "BAWAR04",
      "duration": 4,
      "title": "Kurzreise nach Danzig & Bornholm",
      "startDate": "2025-06-10",
      "endDate": "2025-06-14",
      "flightIncluded": true,
      "transCruise": false,
      "flightDirection": "two-way",
      "adults": 2,
      "children": 1,
      "juveniles": 0,
      "babies": 0,
      "notes": [
        null,
        "4",
        "inkl. Flug"
      ],
      "ports": [
        {
          "code": "DEWAR",
          "name": null,
          "countryCode": null,
          "departureDateTime": "2025-06-10T17:00:00",
          "arrivalDateTime": null
        },
        {
          "code": "SEE",
          "name": null,
          "countryCode": null,
          "departureDateTime": null,
          "arrivalDateTime": null
        },
        {
          "code": "PLGDY",
          "name": null,
          "countryCode": null,
          "departureDateTime": "2025-06-12T18:00:00",
          "arrivalDateTime": "2025-06-12T08:00:00"
        },
        {
          "code": "DKRNN",
          "name": null,
          "countryCode": null,
          "departureDateTime": "2025-06-13T18:00:00",
          "arrivalDateTime": "2025-06-13T08:00:00"
        },
        {
          "code": "DEWAR",
          "name": null,
          "countryCode": null,
          "departureDateTime": null,
          "arrivalDateTime": "2025-06-14T08:00:00"
        }
      ],
      "cruiseItemVariant": [
        {
          "fromCity": "Warnemünde",
          "toCity": "Warnemünde",
          "sortDate": "Jun 10, 2025, 12:00:00 AM",
          "startDate": "2025-06-10",
          "endDate": "2025-06-14",
          "ship": {
            "name": "AIDAmar",
            "code": "MA",
            "marketingName": "AIDAmar"
          },
          "shipVariation": "07",
          "duration": 4,
          "amount": 2998,
          "amountString": "2.998",
          "amountPerPerson": null,
          "journeyIdentifier": "MA04250610",
          "flightIncluded": true,
          "flightDirection": "two-way",
          "matches": true,
          "tariffType": "VARIO",
          "currency": "€",
          "imageUrl": "/content/dam/aida/de/medien/kreuzfahrt/routes/WAR04096/DE_AIDA_ADAPTIVE_WAR04096_de-DE.png",
          "departureAirport": "MUC",
          "notes": [
            null,
            "4",
            "inkl. Flug"
          ],
          "information": [],
          "campaigns": [
            {
              "code": "BFT 202514",
              "name": "Kurzreisen Special",
              "medium": "Aktion BF Template",
              "validity": {
                "bookingDateFrom": "2025-04-03",
                "bookingDateTo": "2025-05-07",
                "currentDate": "2025-05-07"
              }
            },
            {
              "code": "SKV 202518",
              "name": "Mai Push",
              "medium": "Sonderkampagne",
              "validity": {
                "bookingDateFrom": "2025-04-29",
                "bookingDateTo": "2025-05-19",
                "currentDate": "2025-05-07"
              }
            }
          ],
          "amenities": [],
          "transCruise": false,
          "targetPageVariant": null
        }
      ]
    }
  ]
}*/
