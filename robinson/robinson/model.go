package robinson

type OfferList struct {
	Items []Items `json:"items"`
}

type Button struct {
	Href  string `json:"href"`
	Text  string `json:"text"`
	Title string `json:"title"`
}

type Club struct {
	Country     string `json:"country"`
	Headline    string `json:"headline"`
	Id          string `json:"id"`
	ProductCode string `json:"productCode"`
}

type Items struct {
	Club  Club  `json:"club"`
	Offer Offer `json:"offer"`
}

type Offer struct {
	Airport      string `json:"airport"`
	Board        string `json:"board"`
	Button       Button `json:"button"`
	Duration     string `json:"duration"`
	Price        Price  `json:"price"`
	PriceType    string `json:"priceType"`
	Room         string `json:"room"`
	TravelPeriod string `json:"travelPeriod"`
}
type Price struct {
	Info   string  `json:"info"`
	Prefix string  `json:"prefix"`
	Value  float64 `json:"value"`
}

type OfferDetailsResponse struct {
	Offers []OfferDetails `json:"offers"`
}

type OfferDetails struct {
	Currency       string  `json:"currency"`
	Date           string  `json:"date"`
	DiscountAmount float64 `json:"discountAmount"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	PricePerAdult  string  `json:"pricePerAdult"`
	ReturnDate     string  `json:"returnDate"`
}
