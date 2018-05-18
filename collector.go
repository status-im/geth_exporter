package main

import (
	"fmt"
	"log"
	"regexp"
)

type collector struct {
	cl *client
	fs []*regexp.Regexp
}

func newCollector(ipcPath string, rawFilters []string) (*collector, error) {
	cl, err := newClient(ipcPath)
	if err != nil {
		return nil, err
	}

	collector := &collector{cl: cl}
	collector.compileFilters(rawFilters)

	return collector, nil
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

func (c *collector) collect() (string, error) {
	m, err := c.cl.metrics()
	if err != nil {
		return "", err
	}

	all := transformMetrics(m)

	for k := range all {
		if !c.match(k) {
			delete(all, k)
		}
	}

	return all.String(), nil
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
