package web

import (
	"io"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (w *Web) GetMetricsHandler(registry *prometheus.Registry) func(writer http.ResponseWriter, reader *http.Request) {
	return func(writer http.ResponseWriter, reader *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if len(w.proxyAddress) > 0 {
			w.proxyBeforeMetrics(writer, reader)
		}

		promhttp.
			HandlerFor(registry, promhttp.HandlerOpts{
				DisableCompression: true,
			}).
			ServeHTTP(NewFakeResponseWriter(writer), reader)
	}
}

func (w *Web) proxyBeforeMetrics(writer http.ResponseWriter, reader *http.Request) {
	resp, err := w.proxy.R().
		WithContext(reader.Context()).
		SetResponseDoNotParse(true).
		Get(w.proxyAddress)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	writer.WriteHeader(resp.StatusCode())

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
