package main

import (
	"log"
	"regexp"
)

type collector struct {
	ipcPath    string
	fs         []*regexp.Regexp
	syncFails  bool
	debugFails bool
}

func newCollector(ipcPath string, rawFilters []string) *collector {
	collector := &collector{ipcPath: ipcPath, syncFails: false}
	collector.compileFilters(rawFilters)

	return collector
}

func (c *collector) compileFilters(rawFilters []string) {
	for _, raw := range rawFilters {
		f, err := regexp.Compile(raw)
		if err != nil {
			log.Printf("error adding filter %s, %v", raw, err)
			continue
		}

		c.fs = append(c.fs, f)
	}
}

func (c *collector) collect() (flatMetrics, error) {
	cl, err := newClient(c.ipcPath)
	if err != nil {
		return nil, err
	}

	defer cl.close()

	var all flatMetrics

	if !c.debugFails {
		m, err := cl.metrics()
		if err != nil {
			/* not all geth nodes have debug enabled, s might be nil */
			log.Printf("failed to get metrics: %v", err)
			c.debugFails = true
		} else {
			/* can handle m being nil, will just return an empty flatMetrics */
			all = transformMetrics(m)
		}
	}

	if !c.syncFails {
		s, err := cl.syncingMetrics()
		if err != nil {
			/* not all geth nodes have eth enabled, s might be nil */
			log.Printf("failed to get syncing stats: %v", err)
			c.syncFails = true
		} else {
			sync := decodeSyncData(s, "sync_")
			for k, v := range sync {
				all[k] = v
			}
		}
	}

	for k := range all {
		if !c.matchAllFilters(k) {
			delete(all, k)
		}
	}

	return all, nil
}

func (c *collector) matchAllFilters(key string) bool {
	for _, filter := range c.fs {
		if !filter.MatchString(key) {
			return false
		}
	}

	return true
}
