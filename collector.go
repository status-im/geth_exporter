package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	ipcPath string
	fs      []*regexp.Regexp
}

func newCollector(ipcPath string, rawFilters []string) *collector {
	collector := &collector{ipcPath: ipcPath}
	collector.compileFilters(rawFilters)

	return collector
}

func (c *collector) compileFilters(rawFilters []string) {
	for _, raw := range rawFilters {
		s := fmt.Sprintf("^%s", raw)
		f, err := regexp.Compile(s)
		if err != nil {
			log.Printf("error adding filter %s, %v", s, err)
			continue
		}

		c.fs = append(c.fs, f)
	}
}

//FIXME: remove method
func (c *collector) collect() (flatMetrics, error) {
	cl, err := newClient(c.ipcPath)
	if err != nil {
		return nil, err
	}

	m, err := cl.metrics()
	if err != nil {
		return nil, err
	}

	all := transformMetrics(m)

	for k := range all {
		if !c.match(k) {
			delete(all, k)
		}
	}

	return all, nil
}

var scrapeDurationDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "no_subsystem", "no_name"),
	"geth_exporter: Metrics from geth.",
	[]string{"collector"},
	nil,
)

var namespace = "geth"

// Collect implements the prometheus.Collector interface.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	m, err := c.collect()
	if err != nil {
		log.Printf("Error collecting: %+v", err)
		return
	}

	begin := time.Now()
	duration := time.Since(begin)

	for k := range m {
		if c.match(k) {
			ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.UntypedValue, duration.Seconds(), k)
		}
	}
}

// Describe implements the prometheus.Collector interface.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	// ch <- scrapeSuccessDesc
}

func (c *collector) match(key string) bool {
	if len(c.fs) == 0 {
		return true
	}

	for _, filter := range c.fs {
		if filter.MatchString(key) {
			return true
		}
	}

	return false
}
