package main

import "github.com/loadimpact/k6/stats"

var (
	PluginCounter = stats.New("pluginname.counter", stats.Counter)
	PluginGauge   = stats.New("pluginname.gauge", stats.Gauge)
	PluginRate    = stats.New("pluginname.rate", stats.Rate)
	PluginTrend   = stats.New("pluginname.trend", stats.Trend, stats.Time)
)
