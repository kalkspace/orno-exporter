package exporter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var _ prometheus.Collector = (*Exporter)(nil)

const (
	namespace = "orno"
)

type StateSource interface {
	Fetch() (map[string]interface{}, error)
	Metrics() []string
}

type Exporter struct {
	log logrus.FieldLogger

	source StateSource

	metrics      map[string]*prometheus.Desc
	up           prometheus.Gauge
	totalScrapes prometheus.Counter
}

var slugRegex = regexp.MustCompile("[^a-z0-9]+")

func toMetricName(name string) string {
	name = strings.ToLower(name)
	slug := slugRegex.ReplaceAllString(name, "_")
	return fmt.Sprintf("%s_%s", namespace, slug)
}

func prepareMetrics(keys []string) map[string]*prometheus.Desc {
	metrics := make(map[string]*prometheus.Desc)
	for _, label := range keys {
		name := toMetricName(label)
		metrics[name] = prometheus.NewDesc(name, label, nil, nil)
	}
	return metrics
}

func NewExporter(log logrus.FieldLogger, source StateSource) *Exporter {
	return &Exporter{
		log:     log,
		source:  source,
		metrics: prepareMetrics(source.Metrics()),

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Whether last read from source was successful",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_total_scrapes",
			Help:      "Current total exporter scrapes",
		}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range e.metrics {
		ch <- desc
	}
	ch <- e.up.Desc()
	ch <- e.totalScrapes.Desc()
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.totalScrapes.Add(1)
	e.collectFromSource(ch)
	ch <- e.up
	ch <- e.totalScrapes
}

func (e *Exporter) collectFromSource(ch chan<- prometheus.Metric) {
	values, err := e.source.Fetch()
	if err != nil {
		e.log.WithError(err).Warn("Fetching failed")
		e.up.Set(0)
		return
	}
	e.up.Set(1)

	for name, value := range values {
		name = toMetricName(name)
		log := e.log.WithField("metric", name)
		desc, ok := e.metrics[name]
		if !ok {
			log.Warn("Metric not predefined. Skipping...")
			continue
		}
		var val float64
		switch t := value.(type) {
		case float64:
			val = t
		case float32:
			val = float64(t)
		case uint16:
			val = float64(t)
		case uint32:
			val = float64(t)
		case uint64:
			val = float64(t)
		default:
			log.WithField("type", fmt.Sprintf("%T", t)).Warn("Unknown metric type. Skipping...")
			continue
		}
		m, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, val)
		if err != nil {
			log.WithError(err).Warn("Failed creating const metric. Skipping...")
			continue
		}
		ch <- m
	}
}
