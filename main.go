package main

import (
	"net/http"
	"os"

	"github.com/kalkspace/orno-exporter/config"
	"github.com/kalkspace/orno-exporter/exporter"
	"github.com/kalkspace/orno-exporter/orno"
	"github.com/kalkspace/orno-exporter/orno/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var ConfigFile string

func main() {
	log := logrus.StandardLogger()

	if configFileEnv, ok := os.LookupEnv("ORNO_CONFIG_FILE"); ok {
		ConfigFile = configFileEnv
	}
	config, err := config.LoadConfig("ORNO", ConfigFile)
	if err != nil {
		log.WithError(err).Error("Failed loading config")
		return
	}

	reader := orno.NewReader(log.WithField("component", "reader"), config.Serial.Address, new(models.WE516))

	exporter := exporter.NewExporter(log.WithField("component", "exporter"), reader)
	if err := prometheus.Register(exporter); err != nil {
		log.WithError(err).Error("Unable to register exporter")
		return
	}

	addr := config.Server.Host + ":" + config.Server.Port
	log.WithField("address", addr).Info("Starting http server...")
	if err := http.ListenAndServe(addr, promhttp.Handler()); err != nil {
		log.WithError(err).Error("HTTP server failed")
	}
}
