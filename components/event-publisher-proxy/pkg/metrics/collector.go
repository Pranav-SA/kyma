package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Errors name of the errors metric
	Errors = "event_publish_to_messaging_server_errors_total"
	// Latency name of the latency metric
	Latency = "event_publish_to_messaging_server_latency"
	// EventTypePublishedMetricKey name of the eventType metric
	EventTypePublishedMetricKey = "event_type_published"
	//EventRequests name if the eventRequests metric
	EventRequests = "event_requests"
	// errorsHelp help for the errors metric
	errorsHelp = "The total number of errors while sending Events to the messaging server"
	// latencyHelp help for the latency metric
	latencyHelp = "The duration of sending Events to the messaging server"
	// EventTypePublishedMetricHelp help for the eventType metric
	EventTypePublishedMetricHelp = "The total number of events published for a given eventType"
	// EventRequestsHelp help for event requests metric
	EventRequestsHelp = "The total number of event requests"
	//responseCode is the name of the status code labels used by multiple metrics
	responseCode = "response_code"
	//destSvc is the name of the destination service label used by multiple metrics
	destSvc = "destination_service"
)

// Collector implements the prometheus.Collector interface
type Collector struct {
	errors    *prometheus.CounterVec
	latency   *prometheus.HistogramVec
	eventType *prometheus.CounterVec
	requests  *prometheus.CounterVec
}

// NewCollector a new instance of Collector
func NewCollector() *Collector {
	return &Collector{
		errors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: Errors,
				Help: errorsHelp,
			},
			[]string{},
		),
		latency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: Latency,
				Help: latencyHelp,
			},
			[]string{responseCode, destSvc},
		),
		eventType: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: EventTypePublishedMetricKey,
				Help: EventTypePublishedMetricHelp,
			},
			[]string{"event_type", "event_source"},
		),
		requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: EventRequests,
				Help: EventRequestsHelp,
			},
			[]string{responseCode, destSvc},
		),
	}
}

// Describe implements the prometheus.Collector interface Describe method
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.errors.Describe(ch)
	c.latency.Describe(ch)
	c.eventType.Describe(ch)
	c.requests.Describe(ch)
}

// Collect implements the prometheus.Collector interface Collect method
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.errors.Collect(ch)
	c.latency.Collect(ch)
	c.eventType.Collect(ch)
	c.requests.Collect(ch)
}

// RecordError records an error metric
func (c *Collector) RecordError() {
	c.errors.WithLabelValues().Inc()
}

// RecordLatency records a latency metric
func (c *Collector) RecordLatency(duration time.Duration, statusCode int, destSvc string) {
	c.latency.WithLabelValues(fmt.Sprint(statusCode), destSvc).Observe(float64(duration.Milliseconds()))
}

// RecordEventType records a eventType metric
func (c *Collector) RecordEventType(eventType, eventSource string) {
	c.eventType.WithLabelValues(eventType, eventSource).Inc()
}

// RecordRequests records the eventRequests metric
func (c *Collector) RecordRequests(statusCode int, destSvc string) {
	c.requests.WithLabelValues(fmt.Sprint(statusCode), destSvc).Inc()
}
