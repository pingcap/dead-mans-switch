package main

type NotifyInterface interface {
	Notify(summary, detail string) error
}