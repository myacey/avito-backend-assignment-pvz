package metrics

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Tech Metrics

// requestCount - request counter with Vecotor3: handler, method, code.
var requestCount = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http.requests.total",
		Help: "Total number of HTTP requests by handler, method and status code",
	},
	[]string{"handler", "method", "code"},
)

// responseTime - hisogram of response time in seconds.
var responseTime = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http.response.time.seconds",
		Help:    "Histogram of response times for HTTP requests",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"handler", "method"},
)

// Buisness Metrics

// createPVZCount - counter of created PVZs.
var createPvzCount = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "created.pvz.total",
		Help: "Total number of created PVZ",
	},
)

func CreatePVZ() {
	createPvzCount.Inc()
}

// createdReceptionCount - counter of created receptions.
var createdReceptionCount = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "created.reception.total",
		Help: "Total number of created Receptions",
	},
)

func CreateReception() {
	createdReceptionCount.Inc()
}

// addedProductCount - counter of products, added to receptions.
var addedProductCount = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "added.product.total",
		Help: "Total number of added Products",
	},
)

func AddProduct() {
	addedProductCount.Inc()
}

func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":9000", nil); err != nil {
			log.Fatalf("metrics server failed: %v", err)
		}
	}()
}

// GetMetricsMiddleware - middleware func for
// gathering tech metrics.
func GetMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		status := c.Writer.Status()
		path := sanitizePath(c.Request.URL.Path)

		requestCount.WithLabelValues(path, c.Request.Method, getStatusCode(status)).Inc()
		responseTime.WithLabelValues(path, c.Request.Method).Observe(duration)
	}
}

func sanitizePath(p string) string {
	if strings.HasPrefix(p, "/pvz/") {
		return "/pvz/:id"
	}
	if strings.HasPrefix(p, "/reception/") {
		return "/reception/:id"
	}
	return p
}

func getStatusCode(code int) string {
	return fmt.Sprintf("%d", code)
}
