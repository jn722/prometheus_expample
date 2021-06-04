package main

import (
	"context"
	"log"
	"prom_test/src/collector"
	"time"

	"github.com/prometheus/client_golang/api"
	promapi "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func main() {
	go func() {
		time.Sleep(1 * time.Minute)
		c, _ := api.NewClient(api.Config{
			Address:      "http://prometheus:9090",
			RoundTripper: api.DefaultRoundTripper,
		})
		api := promapi.NewAPI(c)
		timeEnd := time.Now()
		value, _, err := api.Query(context.TODO(), "test_stat[1h]", timeEnd)
		if err != nil {
			log.Println(err)
		}
		log.Println(value.Type().String())
		retV := value.(model.Matrix)
		for _, m := range retV {
			v1 := m.Metric["v1"]
			switch string(v1) {
			case "1":
				log.Println("v1:")
				for _, v := range m.Values {
					log.Printf("time:%v value:%v\n", v.Timestamp, v.Value)
				}
			case "2":
				log.Println("v2:")
				for _, v := range m.Values {
					log.Printf("time:%v value:%v\n", v.Timestamp, v.Value)
				}
			}

		}
	}()
	collector.InitProm()
}
