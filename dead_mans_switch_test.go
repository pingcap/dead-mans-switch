package main

import (
	"testing"
	"time"
)

func TestDeadMansSwitchDoesntTrigger(t *testing.T) {
	evaluateMessage := make(chan string)
	defer close(evaluateMessage)

	notify := ""
	d := NewDeadMansSwitch(evaluateMessage, 100*time.Millisecond, func(msg string) error {
		notify = msg
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
	d := NewDeadMansSwitch(evaluateMessage, 100*time.Millisecond, func(msg string) error {
		notify = msg
		return nil
	})

	go d.Run()
	defer d.Stop()

	time.Sleep(150 * time.Millisecond)
	if notify != "DeadMansSwitchDown" {
		t.Fatal("deadman should trigger!")
	}
}

func TestEvaluateMessageNotNull(t *testing.T) {
	evaluateMessage := make(chan string)
	defer close(evaluateMessage)

	notify := ""
	d := NewDeadMansSwitch(evaluateMessage, 100*time.Millisecond, func(msg string) error {
		notify = msg
		return nil
	})

	go d.Run()
	defer d.Stop()

	evaluateMessage <- "alert not as expected"

	if notify != "alert not as expected" {
		t.Fatal("notify msg should equal with evaluate message!")
	}
}
