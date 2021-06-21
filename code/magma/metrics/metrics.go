package metrics

import (
	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type (
	// MagmaService represents struct contains all custom metrics of the Magma GRPC service.
	MagmaService struct {
		// Stores uploaded data info by session ID for each session.
		//
		// Fully-qualified name of the metric: "magma_session_data_uploaded"
		dataUploaded map[string]prometheus.Gauge

		// Stores downloaded data info by session ID for each session.
		//
		// Fully-qualified name of the metric: "magma_session_data_downloaded"
		dataDownloaded map[string]prometheus.Gauge

		// sessionStarted counts started sessions by session ID for each session.
		//
		// Fully-qualified name of the metric: "magma_session_started"
		sessionStarted prometheus.Counter

		// sessionStopped counts stopped sessions by session ID for each session.
		//
		// Fully-qualified name of the metric: "magma_session_stopped"
		sessionStopped prometheus.Counter
	}
)

const (
	// magmaNameSpace represents name space for all MagmaService metrics.
	magmaNameSpace = "magma"

	// sessionSubsystem represents subsystem name for metrics, which stores info about sessions.
	sessionSubsystem = "session"

	// uploadedName represents name for MagmaService.dataUploaded metrics.
	uploadedName = "data_uploaded"
	// downloadedName represents name for MagmaService.dataDownloaded metrics.
	downloadedName = "data_downloaded"

	// startedName represents name for MagmaService.sessionStarted metric.
	startedName = "started"

	// stoppedName represents name for MagmaService.sessionStopped metric.
	stoppedName = "stopped"

	// sessionIDKey represents key for session ID metrics label.
	sessionIDKey = "session_id"
)

// NewMagmaServiceMetrics creates initialized empty MagmaService.
//
// If an error occurs during execution, the program terminates with code 2 and the error will be written in os.Stderr.
func NewMagmaServiceMetrics() *MagmaService {
	sessStarted, err := newSessionStartedCounter()
	if err != nil {
		errors.ExitErr("error while registering session started metric", err, 2)
	}
	sessStopped, err := newSessionStoppedCounter()
	if err != nil {
		errors.ExitErr("error while registering session stopped metric", err, 2)
	}

	return &MagmaService{
		dataUploaded:   make(map[string]prometheus.Gauge),
		dataDownloaded: make(map[string]prometheus.Gauge),
		sessionStarted: sessStarted,
		sessionStopped: sessStopped,
	}
}

// newSessionStartedCounter creates new counter metric with fully-qualified name "magma_session_started".
func newSessionStartedCounter() (prometheus.Counter, error) {
	metric := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: magmaNameSpace,
		Subsystem: sessionSubsystem,
		Name:      startedName,
		Help:      "Counts started sessions.",
	})

	err := prometheus.Register(metric)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

// newSessionStoppedCounter creates new counter metric with fully-qualified name "magma_session_stopped".
func newSessionStoppedCounter() (prometheus.Counter, error) {
	metric := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: magmaNameSpace,
		Subsystem: sessionSubsystem,
		Name:      stoppedName,
		Help:      "Counts stopped sessions.",
	})

	err := prometheus.Register(metric)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

// IncSessionStarted increments counter value of session metric with fully-qualified name "magma_session_started".
func (ms *MagmaService) IncSessionStarted() {
	ms.sessionStarted.Inc()
}

// IncSessionStopped increments counter value of session metric with fully-qualified name "magma_session_stopped".
func (ms *MagmaService) IncSessionStopped() {
	ms.sessionStopped.Inc()
}

// UpdateDataUploadedMetric sets data uploaded value represented in octets
// to session metric with fully-qualified name "magma_session_data_uploaded".
//
// If metric with provided session ID is not registered, new Counter will be created and registered.
func (ms *MagmaService) UpdateDataUploadedMetric(sessionID string, dataUploaded uint64) {
	metric, ok := ms.dataUploaded[sessionID]
	if !ok {
		metric = ms.addNewDataUploadedMetric(sessionID)
	}

	metric.Set(float64(dataUploaded))
}

// addNewDataUploadedMetric creates new metric with fully-qualified name "magma_session_data_uploaded",
// stores metric value with session ID key in MagmaService.dataUploaded map and returns created metric.
//
// Resulted metric has sessionIDKey label which can be used for identifying.
func (ms *MagmaService) addNewDataUploadedMetric(sessionID string) prometheus.Gauge {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: magmaNameSpace,
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

	ms.dataUploaded[sessionID] = metric

	return metric
}

// UpdateDataDownloadedMetric sets data downloaded value represented in octets
// to session metric with fully-qualified name "magma_session_data_downloaded".
//
// If metric with provided session ID is not registered, new Counter will be created and registered.
func (ms *MagmaService) UpdateDataDownloadedMetric(sessionID string, dataDownloaded uint64) {
	metric, ok := ms.dataDownloaded[sessionID]
	if !ok {
		metric = ms.addNewDataDownloadedMetric(sessionID)
	}

	metric.Set(float64(dataDownloaded))
}

// addNewDataUploadedMetric creates new metric with fully-qualified name "magma_session_data_downloaded",
// stores metric value with session ID key in MagmaService.dataDownloaded map and returns created metric.
//
// Resulted metric has sessionIDKey label which can be used for identifying.
func (ms *MagmaService) addNewDataDownloadedMetric(sessionID string) prometheus.Gauge {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: magmaNameSpace,
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

	ms.dataDownloaded[sessionID] = metric

	return metric
}
