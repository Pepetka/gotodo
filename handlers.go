package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pepetka/gotodo/internal/storage"
	"golang.org/x/term"
)

func add(s *storage.Store, args []string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("missing task text")
	}
	text, args := args[0], args[1:]
	var rawPriority string
	var rawDue string
	fl := flag.NewFlagSet("add", flag.ContinueOnError)
	fl.StringVar(&rawPriority, "priority", "medium", "high|medium|low")
	fl.StringVar(&rawDue, "due", "", "YYYY-MM-DD")
	err := fl.Parse(args)
	if err != nil {
		return 0, err
	}
	args = fl.Args()
	if len(args) > 0 {
		return 0, errors.New("too many arguments")
	}

	task := storage.Task{
		Text:    text,
		Status:  storage.StatusActive,
		Created: time.Now(),
		Updated: time.Now(),
	}

	priority, err := ParsePriority(rawPriority)
	if err != nil {
		return 0, err
	}
	task.Priority = priority
	if rawDue != "" {
		due, err := ParseDate(rawDue)
		if err != nil {
			return 0, err
		}
		task.Due = due
	}

	id := s.AddTask(task)
	return id, nil
}

func list(s *storage.Store, args []string) ([]storage.Task, error) {
	fl := flag.NewFlagSet("list", flag.ContinueOnError)
	var rawFilter string
	var rawPriority string
	var rawBefore string
	fl.StringVar(&rawFilter, "filter", "active", "all|active|done")
	fl.StringVar(&rawPriority, "priority", "", "high|medium|low")
	fl.StringVar(&rawBefore, "before", "", "YYYY-MM-DD")
	err := fl.Parse(args)
	if err != nil {
		return nil, err
	}
	args = fl.Args()
	if len(args) > 0 {
		return nil, errors.New("too many arguments")
	}

	filters := []storage.FilterByFn{}

	filter, err := ParseFilter(rawFilter)
	if err != nil {
		return nil, err
	}
	if filter != FilterAll {
		filters = append(filters, storage.FilterByStatusFn(storage.Status(filter)))
	}
	if rawPriority != "" {
		priority, err := ParsePriority(rawPriority)
		if err != nil {
			return nil, err
		}
		filters = append(filters, storage.FilterByPriorityFn(priority))
	}
	if rawBefore != "" {
		due, err := ParseDate(rawBefore)
		if err != nil {
			return nil, err
		}
		filters = append(filters, storage.FilterByDueFn(*due))
	}

	tasks, _ := storage.FilterSlice(s.Tasks, filters...)
	return tasks, nil
}

func done(s *storage.Store, args []string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("missing task id")
	}
	if len(args) > 1 {
		return 0, errors.New("too many arguments")
	}

	taskID, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, err
	}

	ok := s.DoneTask(taskID)
	if !ok {
		return 0, errors.New("task not found")
	}
	return taskID, nil
}

func rm(s *storage.Store, args []string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("missing task id")
	}
	if len(args) > 1 {
		return 0, errors.New("too many arguments")
	}

	taskID, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, err
	}

	ok := s.RmTask(taskID)
	if !ok {
		return 0, errors.New("task not found")
	}
	return taskID, nil
}

func clear(s *storage.Store, args []string) ([]storage.Task, error) {
	var confirm bool
	fl := flag.NewFlagSet("clear", flag.ContinueOnError)
	fl.BoolVar(&confirm, "yes", false, "confirm")
	err := fl.Parse(args)
	if err != nil {
		return nil, err
	}
	args = fl.Args()
	if len(args) > 0 {
		return nil, errors.New("too many arguments")
	}

	if confirm {
		removed := s.RmTasks(storage.FilterByStatusFn(storage.StatusActive))
		return removed, nil
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, errors.New("confirm required: use --yes")
	}
	if !requestConfirm("clear all tasks?") {
		return nil, nil
	}
	removed := s.RmTasks(storage.FilterByStatusFn(storage.StatusActive))
	return removed, nil
}

func requestConfirm(prompt string) bool {
	fmt.Printf("%s [y/N] ", prompt)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(input), "y")
}
