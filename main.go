package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"
	"xray-vars-exporter/internal/pkg/xray"

	"xray-vars-exporter/internal/pkg/metrics"
	"xray-vars-exporter/internal/pkg/web"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/lo"
)

var (
	FlagXrayMetricsAddress  = flag.String("xray-metrics-address", "http://127.0.0.1:11111", "Xray metrics addsress")
	FlagXrayPollingInterval = flag.Duration("xray-polling-interval", 5*time.Second, "Xray polling interval")
	FlagListen              = flag.String("listen", "127.0.0.1:3000", "Exporter listen address")
	FlagMetricsPath         = flag.String("metrics-path", "/metrics", "Path that serves the Xray metrics")
	FlagProxyPassAddress    = flag.String("proxy-pass-metrics-before-address", "", "Append proxy result before the metrics payload (like xray-exporter address)")
)

func init() {
	flag.Parse()
}

func main() {
	ctx := context.Background()
	registry := prometheus.NewRegistry()
	m := metrics.NewMetrics()
	m.MustRegister(registry)

	client := xray.NewClient(lo.FromPtr(FlagXrayMetricsAddress))
	poller := xray.NewPoller(client, m)
	poller.AsyncStartPolling(ctx, lo.FromPtr(FlagXrayPollingInterval))

	mux := http.NewServeMux()
	api := web.NewWeb()
	api.SetProxyAddress(lo.FromPtr(FlagProxyPassAddress))
	mux.HandleFunc(lo.FromPtr(FlagMetricsPath), api.GetMetricsHandler(registry))

	server := &http.Server{
		Addr:    lo.FromPtr(FlagListen),
		Handler: web.ContextTimeoutMiddleware(10*time.Second, mux),
	}

	log.Fatal(server.ListenAndServe())
}
