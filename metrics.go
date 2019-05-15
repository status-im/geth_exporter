package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type metrics map[string]interface{}

type flatMetrics map[string]string

func (fm flatMetrics) String() string {
	var s string
	keys := make([]string, 0)

	for k := range fm {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		s += fmt.Sprintf("%s %s\n", k, fm[k])
	}

	return s
}

func transformMetrics(data metrics) flatMetrics {
	return flattenMetrics(data, make(flatMetrics), "")
}

func flattenMetrics(data metrics, memo flatMetrics, prefix string) flatMetrics {
	for k, v := range data {
		key := prefix + normalizeKey(k)

		switch value := v.(type) {
		case map[string]interface{}:
			memo = flattenMetrics(value, memo, key+"_")
		case string:
			memo[key] = value
		default:
			stringValue := fmt.Sprintf("%+v", value)
			memo[key] = stringValue
		}
	}

	return memo
}

func decodeSyncData(data metrics, prefix string) flatMetrics {
	memo := make(flatMetrics)
	for k, v := range data {
		key := prefix + normalizeKey(k) + "_value"
		val, ok := v.(string)
		if !ok {
			log.Printf("error casting to string: %s -> %v", k, v)
		}
		value := decodeHexAddr(val)
		memo[key] = value
	}
	return memo
}

func decodeHexAddr(s string) string {
	// Prometheus needs a number, not a hex
	number, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return "-1"
	}
	return fmt.Sprintf("%d", number)
}

func normalizeKey(s string) string {
	r, n := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return ""
	}

	return string(unicode.ToLower(r)) + s[n:]
}
