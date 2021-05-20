package main

import (
	"strings"

	"github.com/prometheus/alertmanager/template"
)

func Include(i, j template.Data) {
	exportedReceiveAlerts := make(map[string]bool)
	for _, k := range i.Alerts {
		labelValuesStr := strings.Join(k.Labels.Values(), " ")
		exportedReceiveAlerts[labelValuesStr] = false
	}

	for _, l := range j.Alerts {
		labelValuesStr := strings.Join(l.Labels.Values(), " ")
		exportedReceiveAlerts[labelValuesStr] = true
	}
	for label, received := range exportedReceiveAlerts {
		v := 1.0
		if received {
			v = 0
		}
		failedEvaluatePayload.WithLabelValues(label).Set(v)
	}
}
