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
	"strings"
	"syscall"
	"time"

	"github.com/google/go-cmp/cmp"
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
	http.Handle("/webhook", webhook(evaluateMessage, config.Evaluate))
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
func webhook(evaluateMessage chan<- string, evaluate *Evaluate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		data := new(template.Data)
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("received webhook payload")

		// if evaluateMessage is "", the heartbeat is successful,
		// if it is not empty, use evaluateMessage as notify message
		if evaluate != nil {
			// only compare .Receiver .Status .Alerts as expected
			copyData := template.Data{}
			copyData.Receiver = data.Receiver
			copyData.Status = data.Status
			copyData.Alerts = make(template.Alerts, len(data.Alerts))
			for index, v := range data.Alerts {
				copyData.Alerts[index].Status = v.Status
				copyData.Alerts[index].Labels = v.Labels
			}
			diff := cmp.Diff(evaluate.Data, copyData)
			switch evaluate.Type {
			case EvaluateEqual:
				if diff != "" {
					evaluateMessage <- diff
					fmt.Fprintf(os.Stderr, "error: %s, diff: %s\n", "alert payload not euqal", diff)
					w.WriteHeader(http.StatusOK)
					return
				}
			case EvaluateInclude, "":
				if diff != "" {
					// todo: cmp package does not support get only the more or less part.
					if strings.Contains(diff, "- ") {
						evaluateMessage <- diff
						fmt.Fprintf(os.Stderr, "error: %s, diff: %s\n", "alert payload not included", diff)
						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}
		evaluateMessage <- ""

		w.WriteHeader(http.StatusOK)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok!")
}
