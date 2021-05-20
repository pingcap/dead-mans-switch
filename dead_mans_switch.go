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

	failedEvaluatePayload = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dead_mans_switch_evaluate_failed",
			Help: "The timestamps of failed evaluate.",
		}, []string{"labels"},
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
	message  <-chan string
	interval time.Duration
	ticker   *time.Ticker
	closer   chan struct{}
	notifier func(summary, detail string) error
}

func NewDeadMansSwitch(message <-chan string, interval time.Duration, notifier func(summary, detail string) error) *DeadmansSwitch {
	return &DeadmansSwitch{
		message:  message,
		interval: interval,
		notifier: notifier,
		closer:   make(chan struct{}),
	}
}

func (d *DeadmansSwitch) Run() error {
	log.Println("starting dead mans switch")
	d.ticker = time.NewTicker(d.interval)

	skip := false

	for {
		select {
		case <-d.ticker.C:
			if !skip {
				d.Notify("WatchdogDown", "alerting pipeline is unhealthy")
			} else {
				log.Println("received Deadman's Switch alert, skip notify")
			}

			skip = false

		case <-d.message:
			heatbeatSuccess.Inc()
			skip = true

		case <-d.closer:
			break
		}
	}
}

// Notify send special message to notifier
func (d *DeadmansSwitch) Notify(summary, detail string) {
	if err := d.notifier(summary, detail); err != nil {
		failedNotifications.Inc()
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func (d *DeadmansSwitch) Stop() {
	if d.ticker != nil {
		d.ticker.Stop()
	}

	d.closer <- struct{}{}
}
