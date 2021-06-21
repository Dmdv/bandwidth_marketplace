package metrics

import (
	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type (
	// ProxyService represents struct contains all custom metrics of the Proxy GRPC service.
	ProxyService struct {
		// Stores number of verified acknowledgments.
		//
		// Fully-qualified name of the metric: "proxy_acknowledgment_verified"
		acknowledgmentVerified prometheus.Counter

		// Stores number of unverified acknowledgments.
		//
		// Fully-qualified name of the metric: "proxy_acknowledgment_unverified"
		acknowledgmentUnverified prometheus.Counter

		// Stores uploaded data info by session ID for each session.
		//
		// Fully-qualified name of the metric: "proxy_session_data_uploaded"
		dataUploaded map[string]prometheus.Gauge

		// Stores downloaded data info by session ID for each session.
		//
		// Fully-qualified name of the metric: "proxy_session_data_downloaded"
		dataDownloaded map[string]prometheus.Gauge
	}
)

const (
	// proxyNameSpace represents name space for all ProxyService metrics.
	proxyNameSpace = "proxy"

	// acknowledgmentSubsystem represents subsystem name for metrics, which stores info about acknowledgments.
	acknowledgmentSubsystem = "acknowledgment"

	// verifiedName represents name for ProxyService.acknowledgmentVerified metric.
	verifiedName = "verified"
	// unverifiedName represents name for ProxyService.acknowledgmentUnverified metric.
	unverifiedName = "unverified"

	// sessionSubsystem represents subsystem name for metrics, which stores info about sessions.
	sessionSubsystem = "session"

	// uploadedName represents name for ProxyService.dataUploaded metrics.
	uploadedName = "data_uploaded"

	// downloadedName represents name for ProxyService.dataDownloaded metrics.
	downloadedName = "data_downloaded"

	// sessionIDKey represents key for session ID metrics label.
	sessionIDKey = "session_id"
)

// NewProxyServiceMetrics creates initialized empty ProxyService.
//
// If an error occurs during execution, the program terminates with code 2 and the error will be written in os.Stderr.
func NewProxyServiceMetrics() *ProxyService {
	acknVer, err := newAcknowledgmentVerifiedCounter()
	if err != nil {
		errors.ExitErr("error while registering session started metric", err, 2)
	}

	acknUnver, err := newAcknowledgmentUnverifiedCounter()
	if err != nil {
		errors.ExitErr("error while registering session started metric", err, 2)
	}

	return &ProxyService{
		acknowledgmentVerified:   acknVer,
		acknowledgmentUnverified: acknUnver,
		dataUploaded:             make(map[string]prometheus.Gauge),
		dataDownloaded:           make(map[string]prometheus.Gauge),
	}
}

// addNewDataUploadedMetric creates and registers new metric with fully-qualified name "proxy_acknowledgment_verified"
func newAcknowledgmentVerifiedCounter() (prometheus.Gauge, error) {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: proxyNameSpace,
		Subsystem: acknowledgmentSubsystem,
		Name:      verifiedName,
		Help:      "Counts verified acknowledgments.",
	})
	err := prometheus.Register(metric)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

// addNewDataUploadedMetric creates and registers new metric with fully-qualified name "proxy_acknowledgment_unverified".
func newAcknowledgmentUnverifiedCounter() (prometheus.Gauge, error) {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: proxyNameSpace,
		Subsystem: acknowledgmentSubsystem,
		Name:      unverifiedName,
		Help:      "Counts verified acknowledgments.",
	})
	err := prometheus.Register(metric)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

// IncAcknowledgmentVerified increments counter value of acknowledgment metric
// with fully-qualified name "proxy_acknowledgment_verified".
func (ps *ProxyService) IncAcknowledgmentVerified() {
	ps.acknowledgmentVerified.Inc()
}

// IncAcknowledgmentUnverified increments counter value of acknowledgment metric
// with fully-qualified name "proxy_acknowledgment_unverified".
func (ps *ProxyService) IncAcknowledgmentUnverified() {
	ps.acknowledgmentUnverified.Inc()
}

// UpdateDataUploadedMetric sets data uploaded value represented in octets
// to session metric with fully-qualified name "proxy_session_data_uploaded".
//
// If metric with provided session ID is not registered, new Gauge will be created and registered.
func (ps *ProxyService) UpdateDataUploadedMetric(sessionID string, dataUploaded uint64) {
	metric, ok := ps.dataUploaded[sessionID]
	if !ok {
		metric = ps.addNewDataUploadedMetric(sessionID)
	}

	metric.Set(float64(dataUploaded))
}

// addNewDataUploadedMetric creates new metric with fully-qualified name "proxy_session_data_uploaded",
// stores metric value with session ID key in ProxyService.dataUploaded map and returns created metric.
//
// Resulted metric has sessionIDKey label which can be used for identifying.
func (ps *ProxyService) addNewDataUploadedMetric(sessionID string) prometheus.Gauge {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: proxyNameSpace,
		Subsystem: sessionSubsystem,
		Name:      uploadedName,
		Help:      "Counts uploaded data in octets",
		ConstLabels: map[string]string{
			sessionIDKey: sessionID,
		},
	})
	err := prometheus.Register(metric)
	if err != nil {
		log.Logger.Error("Fail registering metric.", zap.Error(err))
	}

	ps.dataUploaded[sessionID] = metric

	return metric
}

// UpdateDataDownloadedMetric sets data downloaded value represented in octets
// to session metric with fully-qualified name "proxy_session_data_downloaded".
//
// If metric with provided session ID is not registered, new Counter will be created and registered.
func (ps *ProxyService) UpdateDataDownloadedMetric(sessionID string, dataDownloaded uint64) {
	metric, ok := ps.dataDownloaded[sessionID]
	if !ok {
		metric = ps.addNewDataDownloadedMetric(sessionID)
	}

	metric.Set(float64(dataDownloaded))
}

// addNewDataUploadedMetric creates new metric with fully-qualified name "proxy_session_data_downloaded",
// stores metric value with session ID key in ProxyService.dataDownloaded map and returns created metric.
//
// Resulted metric has sessionIDKey label which can be used for identifying.
func (ps *ProxyService) addNewDataDownloadedMetric(sessionID string) prometheus.Gauge {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: proxyNameSpace,
		Subsystem: sessionSubsystem,
		Name:      downloadedName,
		Help:      "Counts downloaded data in octets",
		ConstLabels: map[string]string{
			sessionIDKey: sessionID,
		},
	})
	err := prometheus.Register(metric)
	if err != nil {
		log.Logger.Error("Fail registering metric.", zap.Error(err))
	}

	ps.dataDownloaded[sessionID] = metric

	return metric
}
