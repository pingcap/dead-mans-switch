package main

import (
	"testing"
	"time"
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

	notify := ""
	notifyDetail := ""
	d := NewDeadMansSwitch(evaluateMessage, 100*time.Millisecond, func(summary, detail string) error {
		notify = summary
		notifyDetail = detail
		return nil
	})

	go d.Run()
	defer d.Stop()

	evaluateMessage <- "alert not as expected"

	if notify != "WatchdogAlertPayloadNotAsExpected" {
		t.Fatal("summary should equal with WatchdogAlertPayloadNotAsExpected")
	}

	if notifyDetail != "alert not as expected" {
		t.Fatal("notify detail should be equal with evaluate message")
	}
}
