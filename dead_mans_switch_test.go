package main

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testNotifier struct {
	Summary string
	Detail  string
}

func (tn *testNotifier) Notify(summary, detail string) error {
	tn.Summary = summary
	tn.Detail = detail
	return nil
}

func TestDeadMansSwitchDoesntTrigger(t *testing.T) {
	evaluateMessage := make(chan string)
	defer close(evaluateMessage)

	tn := testNotifier{}
	notifiers := []NotifyInterface{&tn}

	d := NewDeadMansSwitch(evaluateMessage, 100*time.Millisecond, notifiers)

	go d.Run()
	defer d.Stop()

	evaluateMessage <- ""

	time.Sleep(150 * time.Millisecond)
	if tn.Summary != "" {
		t.Fatal("deadman should not trigger!")
	}
}

func TestDeadMansSwitchTrigger(t *testing.T) {
	evaluateMessage := make(chan string)
	defer close(evaluateMessage)

	tn := testNotifier{}
	notifiers := []NotifyInterface{&tn}

	d := NewDeadMansSwitch(evaluateMessage, 100*time.Millisecond, notifiers)

	go d.Run()
	defer d.Stop()

	time.Sleep(150 * time.Millisecond)
	if tn.Summary != "WatchdogDown" {
		t.Fatal("deadman should trigger!")
	}
}

func TestEvaluateMessageNotNull(t *testing.T) {
	evaluateMessage := make(chan string)
	defer close(evaluateMessage)

	tn := testNotifier{}
	notifiers := []NotifyInterface{&tn}

	d := NewDeadMansSwitch(evaluateMessage, 100*time.Millisecond, notifiers)

	go d.Run()
	defer d.Stop()

	evaluateMessage <- "alert not as expected"

	if testutil.ToFloat64(failedEvaluatePayload) == 0 {
		t.Fatal("failedEvaluatePayload should be > 0 when evaluate failed")
	}
}
