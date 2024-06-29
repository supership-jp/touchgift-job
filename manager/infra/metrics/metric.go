package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		counters:   make(map[string]*prometheus.CounterVec),
		histograms: make(map[string]*prometheus.HistogramVec),
	}
}

func (m *Metrics) GetCounter(name string) *prometheus.CounterVec {
	return m.counters[name]
}

func (m *Metrics) AddCounter(name string, description string, labels []string) {
	_, ok := m.counters[name]
	if !ok {
		vec := prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: name, Help: description},
			labels,
		)
		prometheus.MustRegister(vec)
		m.counters[name] = vec
	}
}

func (m *Metrics) GetHistogram(name string) *prometheus.HistogramVec {
	return m.histograms[name]
}

func (m *Metrics) AddHistogram(name string, description string, labels []string, buckets []float64) {
	_, ok := m.histograms[name]
	if !ok {
		vec := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: name, Help: description, Buckets: buckets},
			labels,
		)
		prometheus.MustRegister(vec)
		m.histograms[name] = vec
	}
}
