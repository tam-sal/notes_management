package handlers

import (
	"log/slog"
	"net/http"
	"notes/pkg/response"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handlers struct {
	NoteHandler       *NoteHandler
	CategoryHandler   *CategoryHandler
	UserHandler       *UserHandler
	Logger            *slog.Logger
	HttpErrs          *HttpErrors
	HttpRequestsTotal *prometheus.CounterVec
	HttpDuration      *prometheus.HistogramVec
}

func New(nh *NoteHandler, ch *CategoryHandler, uh *UserHandler, logger *slog.Logger, httpErrs *HttpErrors) *Handlers {
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "path"},
	)

	// Register metrics
	prometheus.MustRegister(requestsTotal, requestDuration)
	return &Handlers{
		NoteHandler:       nh,
		CategoryHandler:   ch,
		UserHandler:       uh,
		Logger:            logger,
		HttpErrs:          httpErrs,
		HttpRequestsTotal: requestsTotal,
		HttpDuration:      requestDuration,
	}
}

func (h *Handlers) RegisterSwaggerHandler(mux *mux.Router) {
	swaggerRoute := mux.PathPrefix("/swagger").Subrouter()
	swaggerRoute.PathPrefix("/").Handler(httpSwagger.WrapHandler)
}

// StatusHandler checks the health of the API.
// @Summary Check API status
// @Description Returns the status of the API to confirm it's running correctly.
// @Tags status
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse "Status OK"
// @Failure 500 {object} map[CustomErrKey]string "Internal server error"
// @Router /status [get]
func (h *Handlers) StatusHandler(w http.ResponseWriter, r *http.Request) {
	res := map[string]string{"status": "OK"}
	if err := response.JSON(w, http.StatusOK, res); err != nil {
		h.HttpErrs.CheckErrType(w, r, err)
		return
	}
}

// MetricsHandler serves Prometheus metrics
// @Summary Get application metrics
// @Description Returns Prometheus metrics for monitoring
// @Tags metrics
// @Produce text/plain
// @Success 200 {string} string "Metrics data"
// @Router /metrics [get]
func (h *Handlers) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
