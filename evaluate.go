package main

import (
	"fmt"
	"reflect"

	"github.com/prometheus/alertmanager/template"
)

func Include(i, j template.Data) (diff string) {
	if i.Status != j.Status {
		return fmt.Sprintf("status different, expected: %v, got: %v", i.Status, j.Status)
	}

	if i.Receiver != j.Receiver {
		return fmt.Sprintf("receiver different, expected: %v, got: %v", i.Status, j.Status)
	}

	for index, k := range i.Alerts {
		include := false
		for _, l := range j.Alerts {
			if reflect.DeepEqual(k, l) {
				include = true
				break
			}
		}
		if include {
			continue
		} else {
			return fmt.Sprintf(".Alerts[%v] not included, got: %v", index, j.Alerts)
		}
	}
	return ""
}
