package main

import (
	"net/http"

	"github.com/kalkspace/orno-exporter/exporter"
	"github.com/kalkspace/orno-exporter/orno"
	"github.com/kalkspace/orno-exporter/orno/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.StandardLogger()

	reader := orno.NewReader(log.WithField("component", "reader"), "/dev/cu.usbserial-AM00EBGZ", new(models.WE516))

	exporter := exporter.NewExporter(log.WithField("component", "exporter"), reader)
	if err := prometheus.Register(exporter); err != nil {
		log.WithError(err).Error("Unable to register exporter")
		return
	}

	addr := ":9090"
	log.WithField("address", addr).Info("Starting http server...")
	if err := http.ListenAndServe(addr, promhttp.Handler()); err != nil {
		log.WithError(err).Error("HTTP server failed")
	}
}
