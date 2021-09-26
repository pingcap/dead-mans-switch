package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	heatbeatSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dead_mans_switch_heatbeat_success",
			Help: "The number of heatbeat receive from webhook.",
		},
	)

	failedNotifications = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dead_mans_switch_notifications_failed",
			Help: "The number of failed notifications.",
		},
	)

	failedEvaluatePayload = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dead_mans_switch_evaluate_failed",
			Help: "The timestamps of failed evaluate.",
		},
	)
)

func init() {
	prometheus.MustRegister(
		heatbeatSuccess,
		failedNotifications,
		failedEvaluatePayload,
	)
}

type DeadmansSwitch struct {
	message   <-chan string
	interval  time.Duration
	ticker    *time.Ticker
	closer    chan struct{}
	notifiers []NotifyInterface
}

func NewDeadMansSwitch(message <-chan string, interval time.Duration, notifiers []NotifyInterface) *DeadmansSwitch {
	return &DeadmansSwitch{
		message:   message,
		interval:  interval,
		notifiers: notifiers,
		closer:    make(chan struct{}),
	}
}

func (d *DeadmansSwitch) Run() error {
	log.Println("starting dead mans switch")
	d.ticker = time.NewTicker(d.interval)

	skip := false
outer:
	for {
		select {
		case <-d.ticker.C:
			if !skip {
				d.Notify("WatchdogDown", "alerting pipeline is unhealthy")
			} else {
				log.Println("received Deadman's Switch alert, skip notify")
			}

			skip = false

		case msg := <-d.message:
			if msg != "" {
				failedEvaluatePayload.SetToCurrentTime()
			} else {
				// message is null, heatbeat success, just skip current check
				failedEvaluatePayload.Set(0)
				heatbeatSuccess.Inc()
				skip = true
			}

		case <-d.closer:
			break outer
		}
	}
	return nil
}

// Notify send special message to notifier
func (d *DeadmansSwitch) Notify(summary, detail string) {
	for _, notifier := range d.notifiers {
		if err := notifier.Notify(summary, detail); err != nil {
			failedNotifications.Inc()
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
	}
}

func (d *DeadmansSwitch) Stop() {
	if d.ticker != nil {
		d.ticker.Stop()
	}

	d.closer <- struct{}{}
}
