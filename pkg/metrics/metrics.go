package metrics

type Counter interface {
	With(labelValues ...string) Counter
	Add(delta float64)
	Inc()
}

type Gauge interface {
	With(labelValues ...string) Gauge
	Set(value float64)
	Add(delta float64)
	Inc()
}

type Histogram interface {
	With(labelValues ...string) Histogram
	Observe(value float64)
}

type Provider interface {
	NewCounter(name string, labelNames ...string) Counter
	NewGauge(name string, labelNames ...string) Gauge
	NewHistogram(name string, buckets []float64, labelNames ...string) Histogram
	Stop()
}
