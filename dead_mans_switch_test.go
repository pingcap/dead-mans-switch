package main

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestDeadMansSwitchDoesntTrigger(t *testing.T) {
	evaluateMessage := make(chan string)
	defer close(evaluateMessage)

	notify := ""
	d := NewDeadMansSwitch(evaluateMessage, 100*time.Millisecond, func(summary, detail string) error {
		notify = summary
		return nil
	})

	go d.Run()
	defer d.Stop()

	evaluateMessage <- ""

	time.Sleep(150 * time.Millisecond)
	if notify != "" {
		t.Fatal("deadman should not trigger!")
	}
}

func TestDeadMansSwitchTrigger(t *testing.T) {
	evaluateMessage := make(chan string)
	defer close(evaluateMessage)

	notify := ""
	d := NewDeadMansSwitch(evaluateMessage, 100*time.Millisecond, func(summary, detail string) error {
		notify = summary

		return nil
	})

	go d.Run()
	defer d.Stop()

	time.Sleep(150 * time.Millisecond)
	if notify != "WatchdogDown" {
		t.Fatal("deadman should trigger!")
	}
}

func TestEvaluateMessageNotNull(t *testing.T) {
	evaluateMessage := make(chan string)
	defer close(evaluateMessage)

	d := NewDeadMansSwitch(evaluateMessage, 100*time.Millisecond, func(summary, detail string) error {
		return nil
	})

	go d.Run()
	defer d.Stop()

	evaluateMessage <- "alert not as expected"

	if testutil.ToFloat64(failedEvaluatePayload) == 0 {
		t.Fatal("failedEvaluatePayload should be > 0 when evaluate failed")
	}
}
