package utils

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	//TODO - user prometheus and grafana loki for tracing, logging and metrics monitoring
	log.Debug("initializing tracing utility for prometheus and grafan loki tracing and metrics monitoring")
}

func LogTracing(guid string, message string) {
	log.Debugf("going to add tracing for guid-  %s, message - %s", guid, message)
}

func IsTracingInitialized() bool {
	return true
}
