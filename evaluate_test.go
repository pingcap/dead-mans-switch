package main

import (
	"testing"

	"github.com/prometheus/alertmanager/template"
)

func TestIncludeWhenDiffOrder(t *testing.T) {
	expected := template.Data{
		Status: "FIRING",
		Alerts: []template.Alert{
			{
				Status: "s1",
			},
			{
				Status: "s2",
			},
		},
	}

	got := template.Data{
		Status: "FIRING",
		Alerts: []template.Alert{
			{
				Status: "s2",
			},
			{
				Status: "s1",
			},
		},
	}

	diff := Include(expected, got)
	if diff != "" {
		t.Fatal("should be include when alerts different order")
	}
}

func TestIncludeWhenMissItems(t *testing.T) {
	expected := template.Data{
		Status: "FIRING",
		Alerts: []template.Alert{
			{
				Status: "s1",
			},
		},
	}

	got := template.Data{
		Status: "FIRING",
		Alerts: []template.Alert{
			{
				Status: "s2",
			},
			{
				Status: "s1",
			},
		},
	}

	diff := Include(expected, got)
	if diff != "" {
		t.Fatal("should be include when alerts different order")
	}
}
