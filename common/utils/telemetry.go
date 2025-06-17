package utils

import (
	log "github.com/sirupsen/logrus"
)

func IsTelemetryInit() bool {
	return true
}

func LogTelemetry(event string, message string) error {
	log.Debugf("Telemetry Event: %s, Message: %s", event, message)
	//TODO- use opentelemetry for capturing all important events for analysis
	//TODO - for tracing and debugging use grafana loki and prometheus
	return nil
}
