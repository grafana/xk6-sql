package main

import (
	"context"
	"log"
	"time"

	"github.com/loadimpact/k6/stats"
)

type pluginName struct{}

func New() *pluginName {
	return &pluginName{}
}

func (*pluginName) Func(ctx context.Context, count int) error {
	state, err := GetState(ctx)

	if err == nil {
		// Do something here
		log.Print(count) // count input variable

		now := time.Now()
		tags := make(map[string]string)
		// tags["tag"] = "value"

		value := 1 // metrics and statistics from the plugin

		stats.PushIfNotDone(ctx, state.Samples, stats.Sample{
			Time:   now,
			Metric: PluginCounter,
			Tags:   stats.IntoSampleTags(&tags),
			Value:  float64(value),
		})

		stats.PushIfNotDone(ctx, state.Samples, stats.Sample{
			Time:   now,
			Metric: PluginGauge,
			Tags:   stats.IntoSampleTags(&tags),
			Value:  float64(value),
		})

		stats.PushIfNotDone(ctx, state.Samples, stats.Sample{
			Time:   now,
			Metric: PluginRate,
			Tags:   stats.IntoSampleTags(&tags),
			Value:  float64(value),
		})

		stats.PushIfNotDone(ctx, state.Samples, stats.Sample{
			Time:   now,
			Metric: PluginTrend,
			Tags:   stats.IntoSampleTags(&tags),
			Value:  float64(value),
		})

		return nil
	}

	return err
}
