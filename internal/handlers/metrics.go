package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"handler", "method", "code"},
	)
	UploadFileSize = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "upload_file_size_bytes",
			Help:    "Size of uploaded files",
			Buckets: prometheus.ExponentialBuckets(1024, 2, 10), // 1KB to ~1MB
		},
	)
	UploadFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "upload_failures_total",
			Help: "Total number of failed uploads",
		},
		[]string{"reason"},
	)
	UploadMIMEs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "upload_mime_total",
			Help: "Total number of uploads by MIME type",
		},
		[]string{"mime"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(HTTPRequestsTotal, UploadFileSize, UploadFailures, UploadMIMEs)
}
