package main

import (
	"fmt"
	"time"

	"github.com/pepetka/gotodo/internal/storage"
)

type Commands string

const (
	CommandAdd   Commands = "add"
	CommandList  Commands = "list"
	CommandDone  Commands = "done"
	CommandRm    Commands = "rm"
	CommandClear Commands = "clear"
)

type Filter string

const (
	FilterAll    Filter = "all"
	FilterActive Filter = Filter(storage.StatusActive)
	FilterDone   Filter = Filter(storage.StatusDone)
)

func ParseCommand(command string) (Commands, error) {
	switch parsed := Commands(command); parsed {
	case CommandAdd, CommandList, CommandDone, CommandRm, CommandClear:
		return parsed, nil
	}
	return "", fmt.Errorf("unknown command: %s", command)
}

func ParsePriority(priority string) (storage.Priority, error) {
	switch parsed := storage.Priority(priority); parsed {
	case storage.PriorityHigh, storage.PriorityMedium, storage.PriorityLow:
		return parsed, nil
	}
	return "", fmt.Errorf("unknown priority: %s", priority)
}

func ParseDate(due string) (*time.Time, error) {
	if due == "" {
		return nil, nil
	}
	date, err := time.ParseInLocation(time.DateOnly, due, time.Local)
	return &date, err
}

func ParseFilter(filter string) (Filter, error) {
	switch parsed := Filter(filter); parsed {
	case FilterAll, FilterActive, FilterDone:
		return parsed, nil
	}
	return "", fmt.Errorf("unknown filter: %s", filter)
}
