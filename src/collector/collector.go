package collector

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type stat_t struct {
	v1 uint32
	v2 uint32
}

type testCollector struct {
	v1    *prometheus.Desc
	value *stat_t
	mu    sync.Mutex
}

func newTestCollector() *testCollector {
	return &testCollector{
		v1: prometheus.NewDesc("test_stat", "test stat", []string{"v1"}, nil),
	}
}

func (c *testCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.v1
}
func (c *testCollector) Collect(ch chan<- prometheus.Metric) {
	var vc *stat_t
	c.mu.Lock()
	if c.value != nil {
		vc = c.value
		c.value = nil
	}
	c.mu.Unlock()
	if vc != nil {
		log.Println(vc)
		ch <- prometheus.MustNewConstMetric(c.v1, prometheus.GaugeValue, float64(vc.v1), "1")
		ch <- prometheus.MustNewConstMetric(c.v1, prometheus.GaugeValue, float64(vc.v2), "2")
	}
}

var defaultCollector = newTestCollector()

func init() {
	prometheus.MustRegister(defaultCollector)
}

func InitProm() {
	go doCollectProcess()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
func (c *testCollector) Update() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.value == nil {
		c.value = &stat_t{}
	}
	c.value.v1 = uint32(rand.Intn(100))
	c.value.v2 = uint32(rand.Intn(200))
}

func doCollectProcess() {
	for {
		defaultCollector.Update()
		// log.Println(defaultCollector.value)
		time.Sleep(20 * time.Second)
	}
}
