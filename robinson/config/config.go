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
	PriceConfig []PriceConfig  `yaml:"priceConfig"`
	Email       emailer.Config `yaml:"email"`
}

type PriceConfig struct {
	ProductCodes string  `yaml:"productCodes"`
	StartDate    string  `yaml:"startDate"`
	EndDate      string  `yaml:"endDate"`
	Duration     int     `yaml:"duration"`
	Adults       int     `yaml:"adults"`
	ChildAge     int     `yaml:"childAge"`
	MinPrice     float64 `yaml:"minPrice"`
}

/**
 * Parse the contents of the YAML file into the Config object.
 */
func (c *Config) Parse(data []byte) error {
	return yaml.Unmarshal(data, &c)
}
