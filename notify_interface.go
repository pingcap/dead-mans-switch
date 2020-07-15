package main

type NotifyInterface interface {
	Notify(msg string) error
}