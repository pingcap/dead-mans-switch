package main

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Interval time.Duration
	Notify   *Notify
	Evaluate *Evaluate
}

type Notify struct {
	Pagerduty *Pagerduty
}

type Pagerduty struct {
	Key string
}

type Evaluate struct {
	// TODO: add evaluate rules
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
