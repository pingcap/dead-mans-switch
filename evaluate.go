package main

import (
	"strings"

	"github.com/prometheus/alertmanager/template"
)

func InitPrometheusValues(evaluate *Evaluate) {
	for _, j := range evaluate.Data.Alerts {
		labelValuesStr := strings.Join(j.Labels.Values(), " ")
		lastAlertTimestamp.WithLabelValues(labelValuesStr).Set(0)
	}
}

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
		if received {
			lastAlertTimestamp.WithLabelValues(label).SetToCurrentTime()
		}
	}
}
