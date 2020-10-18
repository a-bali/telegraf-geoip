# Telegraf GeoIP processor plugin

This processor plugin for [telegraf](https://github.com/influxdata/telegraf) looks up IP addresses in the [MaxMind GeoLite2](https://dev.maxmind.com/geoip/geoip2/geolite2/) database and adds the respective ISO country code, city name, latitude and longitude as new fields to the output.

# Installation

This module is to be used as an external plugin to telegraf, therefore first compile it using Go:

    $ git clone https://github.com/a-bali/telegraf-geoip
    $ cd telegraf-geoip
    $ go build -o geoip cmd/main.go

This will create a standalone binary named `geoip`.

# Usage

You will need to add this plugin as an external plugin to your telegraf config as follows:

    [[processors.execd]]
    command = ["/path/to/geoip_binary", "--config", "/path/to/geoip_config_file"]

# Configuration

As specified above, the plugin uses a separate configuration file, where you can specify where it can find the downloaded GeoLite2 database (you will need the City version), which field to read as input and how to name the newly created fields. For details, please see the [sample config](https://github.com/a-bali/telegraf-geoip/blob/master/plugin.conf).

# License

This software is licensed under GPL 3.0.
