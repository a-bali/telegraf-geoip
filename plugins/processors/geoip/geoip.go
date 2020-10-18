package geoip

import (
	"fmt"
	"net"

	"github.com/IncSW/geoip2"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

const sampleConfig = `
  ## db_path is the location of the MaxMind GeoIP2 City database
  db_path = "/var/lib/GeoIP/GeoLite2-City.mmdb"

  [[processors.geoip.lookup]
	# get the ip from the field "source_ip" and put the lookup results in the respective destination fields (if specified)
	field = "source_ip"
	dest_country = "source_country"
	dest_city = "source_city"
	dest_lat = "source_lat"
	dest_lon = "source_lon"
  `

type lookupEntry struct {
	Field       string `toml:"field"`
	DestCountry string `toml:"dest_country"`
	DestCity    string `toml:"dest_city"`
	DestLat     string `toml:"dest_lat"`
	DestLon     string `toml:"dest_lon"`
}

type GeoIP struct {
	DBPath  string          `toml:"db_path"`
	Lookups []lookupEntry   `toml:"lookup"`
	Log     telegraf.Logger `toml:"-"`
}

var reader *geoip2.CityReader

func (g *GeoIP) SampleConfig() string {
	return sampleConfig
}

func (g *GeoIP) Description() string {
	return "GeoIP looks up the country code, city name and latitude/longitude for IP addresses in the MaxMind GeoIP database"
}

func (g *GeoIP) Apply(metrics ...telegraf.Metric) []telegraf.Metric {
	for _, point := range metrics {
		for _, lookup := range g.Lookups {
			if lookup.Field != "" {
				if value, ok := point.GetField(lookup.Field); ok {
					record, err := reader.Lookup(net.ParseIP(value.(string)))
					if err != nil {
						g.Log.Errorf("GeoIP lookup error: %v", err)
						continue
					}
					if len(lookup.DestCountry) > 0 {
						point.AddField(lookup.DestCountry, record.Country.ISOCode)
					}
					if len(lookup.DestCity) > 0 {
						point.AddField(lookup.DestCity, record.City.Names["en"])
					}
					if len(lookup.DestLat) > 0 {
						point.AddField(lookup.DestLat, record.Location.Latitude)
					}
					if len(lookup.DestLon) > 0 {
						point.AddField(lookup.DestLon, record.Location.Longitude)
					}

				}
			}
		}
	}
	return metrics
}

func (g *GeoIP) Init() error {
	r, err := geoip2.NewCityReaderFromFile(g.DBPath)
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
			DBPath: "/var/lib/GeoIP/GeoLite2-Country.mmdb",
		}
	})
}
