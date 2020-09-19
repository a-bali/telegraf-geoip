package geoip

import (
	"fmt"
	"net"

	"github.com/IncSW/geoip2"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

const sampleConfig = `
  ## country_db is the location of the MaxMind GeoIP2 Country database
  country_db = "/home/bali/GeoLite2-Country.mmdb"

  [[processors.geoip.lookup]
	# get the ip from the field "source_ip" and put the result in the field "source_country"
	field = "source_ip"
	dest = "source_country"
  `

type lookupEntry struct {
	Field string `toml:"field"`
	Dest  string `toml:"dest"`
}

type GeoIP struct {
	CountryDB string          `toml:"country_db"`
	Lookups   []lookupEntry   `toml:"lookup"`
	Log       telegraf.Logger `toml:"-"`
}

var reader *geoip2.CountryReader

func (g *GeoIP) SampleConfig() string {
	return sampleConfig
}

func (g *GeoIP) Description() string {
	return "GeoIP looks up the country code for IP addresses in the MaxMind GeoIP database"
}

func (g *GeoIP) Apply(metrics ...telegraf.Metric) []telegraf.Metric {
	for _, point := range metrics {
		for _, lookup := range g.Lookups {
			if lookup.Dest == "" {
				continue
			}
			if lookup.Field != "" {
				if value, ok := point.GetField(lookup.Field); ok {
					record, err := reader.Lookup(net.ParseIP(value.(string)))
					if err != nil {
						g.Log.Errorf("GeoIP lookup error: %v", err)
						continue
					}
					point.AddField(lookup.Dest, record.Country.ISOCode)
				}
			}
		}
	}
	return metrics
}

func (g *GeoIP) Init() error {
	r, err := geoip2.NewCountryReaderFromFile(g.CountryDB)
	if err != nil {
		return fmt.Errorf("Error opening GeoIP database: %v", err)
	} else {
		reader = r
	}
	return nil
}

func init() {
	processors.Add("geoip", func() telegraf.Processor {
		return &GeoIP{
			CountryDB: "/var/lib/GeoIP/GeoLite2-Country.mmdb",
		}
	})
}
