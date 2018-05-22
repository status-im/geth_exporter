package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type registry struct {
	c *collector
	*prometheus.Registry
}

func newRegistry(ipcPath string) *registry {
	return &registry{
		c:        newCollector(ipcPath, []string{}),
		Registry: prometheus.NewRegistry(),
	}
}

func (r *registry) Gather() ([]*dto.MetricFamily, error) {
	var families []*dto.MetricFamily

	fs, err := r.Registry.Gather()
	if err != nil {
		log.Println(err)
	} else {
		families = fs
	}

	fm, err := r.c.collect()
	if err != nil {
		log.Printf("Error collecting: %+v", err)
		return families, err
	}

	for k, _ := range fm {
		labelPairs := make([]*dto.LabelPair, 0)
		m := &dto.Metric{
			Label:   labelPairs,
			Untyped: &dto.Untyped{Value: proto.Float64(100)},
		}

		name := fmt.Sprintf("geth_%s", k)

		metricFamily := &dto.MetricFamily{}
		metricFamily.Name = proto.String(name)
		metricFamily.Help = proto.String("test help")
		metricFamily.Type = dto.MetricType_UNTYPED.Enum()
		metricFamily.Metric = []*dto.Metric{m}

		families = append(families, metricFamily)
	}

	return families, nil
}
