package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type Builder struct {
	Name       string
	Namespace  string
	Subsystem  string
	InstanceID string
	Help       string
}

func NewBuilder(name string, namespace string, subsystem string, instanceID string, help string) *Builder {
	return &Builder{Name: name, Namespace: namespace, Subsystem: subsystem, InstanceID: instanceID, Help: help}
}

func (b *Builder) BuildResponseTime() gin.HandlerFunc {
	labels := []string{"method", "pattern", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_resp_time",
		Help:      b.Help,
		ConstLabels: map[string]string{
			"instance_id": b.InstanceID,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, labels)
	prometheus.MustRegister(vector)
	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			pattern := ctx.FullPath()
			status := ctx.Writer.Status()
			vector.WithLabelValues(method, pattern, strconv.Itoa(status)).Observe(float64(duration))
		}()

		ctx.Next()
	}
}

func (b *Builder) BuildActiveRequest() gin.HandlerFunc {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_active_request",
		Help:      b.Help,
		ConstLabels: map[string]string{
			"instance_id": b.InstanceID,
		},
	})

	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		gauge.Inc()
		defer gauge.Dec()

		ctx.Next()
	}
}
