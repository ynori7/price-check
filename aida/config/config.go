package config

import (
	"flag"

	"github.com/ynori7/price-check/emailer"
	yaml "gopkg.in/yaml.v2"
)

var CliConf CliConfig

type CliConfig struct {
	ConfigFile string
}

func ParseCliFlags() {
	configFile := flag.String("config", "", "the path to the configuration yaml")

	flag.Parse()

	CliConf.ConfigFile = *configFile
}

type Config struct {
	PriceConfig PriceConfig    `yaml:"priceConfig"`
	Email       emailer.Config `yaml:"email"`
}

type PriceConfig struct {
	Adults                         int      `yaml:"adults"`
	Children                       int      `yaml:"children"`
	Juveniles                      int      `yaml:"juveniles"`
	Babies                         int      `yaml:"babies"`
	Durations                      []string `yaml:"durations"`
	Ports                          []string `yaml:"ports"`
	Airports                       []string `yaml:"airports"`
	PreferredCabinType             string   `yaml:"preferredCabinType"`
	Trips                          []Trip   `yaml:"trips"`
	WithFlightDayPriceThreshold    float64  `yaml:"withFlightDayPriceThreshold"`
	WithoutFlightDayPriceThreshold float64  `yaml:"withoutFlightDayPriceThreshold"`
}

type Trip struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

/**
 * Parse the contents of the YAML file into the Config object.
 */
func (c *Config) Parse(data []byte) error {
	return yaml.Unmarshal(data, &c)
}
