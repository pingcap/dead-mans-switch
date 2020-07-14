package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var configPath = flag.String("config", "/etc/deadmansswitch/config.yaml", "Path to config.yaml.")

func main() {
	flag.Parse()

	config, err := ParseConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	evaluateMessage := make(chan string)
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/webhook", webhook(evaluateMessage))
	http.HandleFunc("/health", health)
	s := &http.Server{
		Addr: ":8080",
	}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	pagerDuty := NewPagerDutyNotify(config.Notify.Pagerduty.Key)
	dms := NewDeadMansSwitch(evaluateMessage, config.Interval, pagerDuty.Notify)
	go dms.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	dms.Stop()
}

// webhook access alert manager alert message and evaluate alert as expected
func webhook(evaluateMessage chan<- string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		data := template.Data{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("received webhook payload: %+v\n", data)

		// TODO: evaluate alert as expected
		// if evaluateMessage is "", the heartbeat is successful,
		// if it is not empty, use evaluateMessage as notify message
		evaluateMessage <- ""

		w.WriteHeader(http.StatusOK)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok!")
}
