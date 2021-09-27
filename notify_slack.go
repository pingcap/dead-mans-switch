package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

var _ NotifyInterface = new(Slack)

type Slack struct {
	WebhookURL string
}

func NewSlackNotify(webhookURL string) *Slack {
	return &Slack{WebhookURL: webhookURL}
}

// Notify send notify message to Slack
func (s *Slack) Notify(summary, detail string) error {
	log.Printf("sending notify: %s to slack\n", summary)

	rendered := fmt.Sprintf(slackTmpl, summary, detail, time.Now().Format(time.RFC3339))

	req, err := http.NewRequest(http.MethodPost, s.WebhookURL, bytes.NewBuffer([]byte(rendered)))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if resp.StatusCode >= 400 || buf.String() != "ok" {
		return fmt.Errorf("failed to send message to slack. status: %d - %s, body: %s", resp.StatusCode, resp.Status, buf.String())
	}

	return nil
}

const slackTmpl = `
{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "DeadMansSwitch - %s"
			}
		},
		{
			"type": "section",
			"fields": [
				{
					"type": "mrkdwn",
					"text": "*Details:*\n%s"
				},
				{
					"type": "mrkdwn",
					"text": "*When:*\n%s"
				}
			]
		}
	]
}`
