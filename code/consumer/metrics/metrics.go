package metrics

import (
	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type (
	// ProxyService represents struct contains all custom metrics of the Proxy GRPC service.
	ProxyService struct {
		// Counts accepted terms for each user.
		//
		// Fully-qualified name of the metric: "proxy_terms_accepted"
		termsAccepted map[string]prometheus.Counter

		// Counts declined terms for each user.
		//
		// Fully-qualified name of the metric: "proxy_terms_declined"
		termsDeclined map[string]prometheus.Counter
	}
)

const (
	// proxyNameSpace represents name space for all ProxyService metrics.
	proxyNameSpace = "proxy"

	// termsSubsystem represents subsystem name for metrics, which stores info about terms.
	termsSubsystem = "terms"

	// acceptedName represents name for ProxyService.termsAccepted metrics.
	acceptedName = "accepted"

	// acceptedName represents name for ProxyService.termsDeclined metrics.
	declinedName = "declined"

	// userIDKey represents key for user ID metrics label.
	userIDKey = "user_id"
)

// NewProxyServiceMetrics creates initialized empty ProxyService.
func NewProxyServiceMetrics() *ProxyService {
	return &ProxyService{
		termsAccepted: make(map[string]prometheus.Counter),
		termsDeclined: make(map[string]prometheus.Counter),
	}
}

// UpdateTermsAcceptedMetric increments counter to terms metric with fully-qualified name "proxy_terms_accepted".
//
// If metric with provided user ID is not registered, new Counter will be created and registered.
func (ps *ProxyService) UpdateTermsAcceptedMetric(userID string) {

	metric, ok := ps.termsAccepted[userID]
	if !ok {
		metric = ps.addNewTermsAcceptedMetric(userID)
	}

	metric.Inc()
}

// addNewTermsAcceptedMetric creates new metric with fully-qualified name "proxy_terms_accepted",
// stores metric value with user ID key in ProxyService.termsAccepted map and returns created metric.
//
// Resulted metric has userIDKey, label which can be used for identifying.
func (ps *ProxyService) addNewTermsAcceptedMetric(userID string) prometheus.Gauge {

	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: proxyNameSpace,
		Subsystem: termsSubsystem,
		Name:      acceptedName,
		Help:      "Counts accepted terms",
		ConstLabels: map[string]string{
			userIDKey: userID,
		},
	})
	err := prometheus.Register(metric)
	if err != nil {
		log.Logger.Error("Fail registering metric.", zap.Error(err))
	}

	ps.termsAccepted[userID] = metric

	return metric
}

// UpdateTermsDeclinedMetric increments counter to terms metric with fully-qualified name "proxy_terms_declined".
//
// If metric with provided user ID is not registered, new Counter will be created and registered.
func (ps *ProxyService) UpdateTermsDeclinedMetric(userID string) {
	metric, ok := ps.termsAccepted[userID]
	if !ok {
		metric = ps.addNewTermsDeclinedMetric(userID)
	}

	metric.Inc()
}

// addNewTermsAcceptedMetric creates new metric with fully-qualified name "proxy_terms_declined",
// stores metric value with user ID key in ProxyService.termsDeclined map and returns created metric.
//
// Resulted metric has userIDKey label which can be used for identifying.
func (ps *ProxyService) addNewTermsDeclinedMetric(userID string) prometheus.Gauge {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: proxyNameSpace,
		Subsystem: termsSubsystem,
		Name:      declinedName,
		Help:      "Counts uploaded data in octets",
		ConstLabels: map[string]string{
			userIDKey: userID,
		},
	})
	err := prometheus.Register(metric)
	if err != nil {
		log.Logger.Error("Fail registering metric.", zap.Error(err))
	}

	ps.termsAccepted[userID] = metric

	return metric
}
