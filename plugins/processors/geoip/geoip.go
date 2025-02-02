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
  # db_type = "city" # city, country or asn (default: city)

  [[processors.geoip.lookup]
	# get the ip from the field "source_ip" and put the lookup results in the respective destination fields (if specified)
	field = "source_ip"
	dest_country = "source_country"
	dest_city = "source_city"
	dest_lat = "source_lat"
	dest_lon = "source_lon"
	# from the ASN database
	dest_autonomous_system_organization = "source_autonomous_system_organization"
	dest_autonomous_system_number = "source_autonomous_system_number"
	dest_network = "source_network"
  `

type lookupEntry struct {
	Field                            string `toml:"field"`
	DestCountry                      string `toml:"dest_country"`
	DestCity                         string `toml:"dest_city"`
	DestLat                          string `toml:"dest_lat"`
	DestLon                          string `toml:"dest_lon"`
	DestAutonomousSystemOrganization string `toml:"dest_autonomous_system_organization"`
	DestAutonomousSystemNumber       string `toml:"dest_autonomous_system_number"`
	DestNetwork                      string `toml:"dest_network"`
}

type GeoIP struct {
	DBPath  string          `toml:"db_path"`
	DBType  string          `toml:"db_type"`
	Lookups []lookupEntry   `toml:"lookup"`
	Log     telegraf.Logger `toml:"-"`
}

var cityReader *geoip2.CityReader
var countryReader *geoip2.CountryReader
var ASNReader *geoip2.ASNReader

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
					if g.DBType == "city" || g.DBType == "" {
						record, err := cityReader.Lookup(net.ParseIP(value.(string)))
						if err != nil {
							if err.Error() != "not found" {
								g.Log.Errorf("GeoIP lookup error: %v", err)
							}
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
					} else if g.DBType == "country" {
						record, err := countryReader.Lookup(net.ParseIP(value.(string)))
						if err != nil {
							if err.Error() != "not found" {
								g.Log.Errorf("GeoIP lookup error: %v", err)
							}
							continue
						}
						if len(lookup.DestCountry) > 0 {
							point.AddField(lookup.DestCountry, record.Country.ISOCode)
						}
					} else if g.DBType == "asn" {
						record, err := ASNReader.Lookup(net.ParseIP(value.(string)))
						if err != nil {
							if err.Error() != "not found" {
								g.Log.Errorf("GeoIP lookup error: %v", err)
							}
							continue
						}
						if len(lookup.DestAutonomousSystemNumber) > 0 {
							point.AddField(lookup.DestAutonomousSystemNumber, record.AutonomousSystemNumber)
						}
						if len(lookup.DestAutonomousSystemOrganization) > 0 {
							point.AddField(lookup.DestAutonomousSystemOrganization, record.AutonomousSystemOrganization)
						}
						if len(lookup.DestNetwork) > 0 {
							point.AddField(lookup.DestNetwork, record.Network)
						}
					} else {
						g.Log.Errorf("Invalid GeoIP database type specified: %s", g.DBType)
					}
				}
			}
		}
	}
	return metrics
}

func (g *GeoIP) Init() error {
	if g.DBType == "city" || g.DBType == "" {
		r, err := geoip2.NewCityReaderFromFile(g.DBPath)
		if err != nil {
			return fmt.Errorf("Error opening GeoIP database: %v", err)
		} else {
			cityReader = r
		}
	} else if g.DBType == "country" {
		r, err := geoip2.NewCountryReaderFromFile(g.DBPath)
		if err != nil {
			return fmt.Errorf("Error opening GeoIP database: %v", err)
		} else {
			countryReader = r
		}
	} else if g.DBType == "asn" {
		r, err := geoip2.NewASNReaderFromFile(g.DBPath)
		if err != nil {
			return fmt.Errorf("Error opening GeoIP database: %v", err)
		} else {
			ASNReader = r
		}
	} else {
		return fmt.Errorf("Invalid GeoIP database type specified: %s", g.DBType)
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
