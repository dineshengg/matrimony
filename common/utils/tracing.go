package utils

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yukitsune/lokirus"
)

var (
	initOnce        sync.Once
	tracingInit     bool
	requestsCounter prometheus.Counter
	logChan         chan *logrus.Entry
	logWorkerWG     sync.WaitGroup
	Logger          *logrus.Logger // Exported logger for use in other modules
)

func init() {
	initOnce.Do(func() {
		Logger = logrus.New()

		// Prometheus metrics setup
		if viper.GetBool("prometheus.enabled") {
			addr := viper.GetString("prometheus.listen_addr")
			if addr == "" {
				addr = ":2112"
			}
			requestsCounter = prometheus.NewCounter(prometheus.CounterOpts{
				Name: "app_requests_total",
				Help: "Total number of requests processed.",
			})
			prometheus.MustRegister(requestsCounter)
			go func() {
				http.Handle("/metrics", promhttp.Handler())
				_ = http.ListenAndServe(addr, nil)
			}()
		}

		// Grafana Loki setup
		if viper.GetBool("loki.enabled") {
			lokiURL := viper.GetString("loki.url")
			if lokiURL == "" {
				lokiURL = "http://localhost:3100/loki/api/v1/push"
			}
			labels := lokirus.Labels{
				"app":         viper.GetString("loki.app"),
				"environment": viper.GetString("loki.environment"),
			}
			opts := lokirus.NewLokiHookOptions().
				WithFormatter(&logrus.JSONFormatter{}).
				WithStaticLabels(labels)
			if viper.IsSet("loki.username") && viper.IsSet("loki.password") {
				opts = opts.WithBasicAuth(
					viper.GetString("loki.username"),
					viper.GetString("loki.password"),
				)
			}
			hook := lokirus.NewLokiHookWithOpts(
				lokiURL,
				opts,
				logrus.InfoLevel,
				logrus.WarnLevel,
				logrus.ErrorLevel,
				logrus.FatalLevel,
			)
			Logger.AddHook(hook)

			// Non-blocking log channel for high QPS
			logChan = make(chan *logrus.Entry, 10000)
			for i := 0; i < 8; i++ {
				logWorkerWG.Add(1)
				go lokiLogWorker(hook)
			}
		}

		tracingInit = true
		Logger.Debug("Tracing utility initialized for Prometheus and Grafana Loki")
	})
}

func lokiLogWorker(hook logrus.Hook) {
	defer logWorkerWG.Done()
	for entry := range logChan {
		logEntry := entry.WithFields(logrus.Fields{
			"app":         viper.GetString("loki.app"),
			"environment": viper.GetString("loki.environment"),
		})
		_ = hook.Fire(logEntry)
	}
}

// LogTracing logs tracing info and increments Prometheus counter.
func LogTracing(guid string, message string) {
	if requestsCounter != nil {
		requestsCounter.Inc()
	}
	entry := Logger.WithField("guid", guid)
	entry.Debugf("tracing: %s", message)
	// Non-blocking send to Loki
	if logChan != nil {
		select {
		case logChan <- entry:
		default:
			// Drop log if channel is full to avoid blocking
		}
	}
}

// IsTracingInitialized returns true if tracing has been initialized.
func IsTracingInitialized() bool {
	return tracingInit
}
