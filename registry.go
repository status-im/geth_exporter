package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const namespace = "geth"

type registry struct {
	c *collector
	*prometheus.Registry
}

func newRegistry(ipcPath string, rawFilters []string) *registry {
	return &registry{
		c:        newCollector(ipcPath, rawFilters),
		Registry: prometheus.NewRegistry(),
	}
}

// Gather implements the prometheus.Registerer interface
func (r *registry) Gather() (families []*dto.MetricFamily, err error) {
	fs, err := r.Registry.Gather()
	if err != nil {
		log.Println(err)
	} else {
		families = fs
	}

	fm, err := r.c.collect()
	if err != nil {
		log.Printf("Error collecting: %+v", err)
		return
	}

	for k, v := range fm {
		mf, err := r.buildMetricFamily(k, v)
		if err != nil {
			log.Printf("error building metric: %v", err)
			continue
		}

		families = append(families, mf)

	}

	return
}

func (r *registry) buildMetricFamily(name string, stringValue string) (*dto.MetricFamily, error) {
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return nil, err
	}

	labelPairs := make([]*dto.LabelPair, 0)
	m := &dto.Metric{
		Label:   labelPairs,
		Untyped: &dto.Untyped{Value: proto.Float64(value)},
	}

	name = fmt.Sprintf("%s_%s", namespace, name)

	metricFamily := &dto.MetricFamily{}
	metricFamily.Name = proto.String(name)
	metricFamily.Help = proto.String("metric exported from geth with debug.metrics")
	metricFamily.Type = dto.MetricType_UNTYPED.Enum()
	metricFamily.Metric = []*dto.Metric{m}

	return metricFamily, nil
}
