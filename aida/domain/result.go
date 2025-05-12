package domain

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`

	Region     string `json:"region"`
	Duration   string `json:"duration"`
	Port       string `json:"port"`
	WithFlight bool   `json:"with_flight"`

	BaseDayPrice   string `json:"day_price"`
	TotalBasePrice string `json:"total_price"`

	IsBaseBelowThreshold      bool `json:"is_base_below_threshold"`
	IsPreferredBelowThreshold bool `json:"is_preferred_below_threshold"`

	TimePeriod        string  `json:"-"`
	BaseDayPriceFloat float64 `json:"-"`
}
