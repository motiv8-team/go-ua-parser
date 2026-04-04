// Example: Prometheus metrics for go-ua-parser using post-parse hooks.
//
// Exposes /metrics with:
//   - ua_parse_duration_seconds  — histogram of parse latencies
//   - ua_requests_total          — counter by browser, os, device_type, is_bot
//
// Run:
//
//	go run . &
//	curl -s localhost:8080/ -A "Mozilla/5.0 (Linux; Android 14) Chrome/123.0 Mobile Safari/537.36"
//	curl -s localhost:8080/metrics | grep ua_
package main

import (
	"fmt"
	"net/http"
	"time"

	uax "github.com/motiv8-team/go-ua-parser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	parseDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "ua_parse_duration_seconds",
		Help:    "Histogram of User-Agent parse latencies.",
		Buckets: prometheus.DefBuckets,
	})

	requestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ua_requests_total",
		Help: "Total requests by browser, os, device type, and bot status.",
	}, []string{"browser", "os", "device_type", "is_bot"})
)

func init() {
	prometheus.MustRegister(parseDuration, requestsTotal)
}

func main() {
	parser, err := uax.NewParser(
		uax.WithPostParseHook(func(_ uax.Input, r uax.Result, d time.Duration) {
			parseDuration.Observe(d.Seconds())

			isBot := "false"
			if r.IsBot {
				isBot = "true"
			}
			requestsTotal.WithLabelValues(
				r.ShortBrowser(),
				r.ShortOS(),
				r.DeviceClass(),
				isBot,
			).Inc()
		}),
	)
	if err != nil {
		panic(err)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		result := parser.ParseRequest(r)
		fmt.Fprintf(w, "Browser: %s\nOS: %s\nDevice: %s\nBot: %v\n",
			result.ShortBrowser(), result.ShortOS(), result.DeviceClass(), result.IsBot)
	})

	fmt.Println("Listening on :8080 (metrics at /metrics)")
	http.ListenAndServe(":8080", nil)
}
