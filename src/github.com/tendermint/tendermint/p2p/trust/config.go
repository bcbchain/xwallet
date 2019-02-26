package trust

import "time"

type TrustMetricConfig struct {
	ProportionalWeight	float64

	IntegralWeight	float64

	TrackingWindow	time.Duration

	IntervalLength	time.Duration
}

func DefaultConfig() TrustMetricConfig {
	return TrustMetricConfig{
		ProportionalWeight:	0.4,
		IntegralWeight:		0.6,
		TrackingWindow:		(time.Minute * 60 * 24) * 14,
		IntervalLength:		1 * time.Minute,
	}
}

func customConfig(tmc TrustMetricConfig) TrustMetricConfig {
	config := DefaultConfig()

	if tmc.ProportionalWeight > 0 {
		config.ProportionalWeight = tmc.ProportionalWeight
	}

	if tmc.IntegralWeight > 0 {
		config.IntegralWeight = tmc.IntegralWeight
	}

	if tmc.IntervalLength > time.Duration(0) {
		config.IntervalLength = tmc.IntervalLength
	}

	if tmc.TrackingWindow > time.Duration(0) &&
		tmc.TrackingWindow >= config.IntervalLength {
		config.TrackingWindow = tmc.TrackingWindow
	}
	return config
}
