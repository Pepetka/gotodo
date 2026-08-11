package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pepetka/gotodo/internal/storage"
)

var red = color.New(color.FgRed).SprintFunc()
var bold = color.New(color.Bold).SprintFunc()
var faint = color.New(color.Faint).SprintFunc()

func format(tasks []storage.Task, lastID int) {
	idLen := strconv.Itoa(len(strconv.Itoa(lastID)))
	dueLen := "3"
	for _, task := range tasks {
		if task.Due != nil {
			dueLen = "10"
			break
		}
	}
	fID := "%-" + idLen + "s"
	fStatus := "%-6s"
	fPriority := "%-8s"
	fDue := "%-" + dueLen + "s"
	fText := "%s"
	fSlice := []string{fID, fStatus, fPriority, fDue, fText}
	f := strings.Join(fSlice, "  ")
	fmt.Printf(f, "#", "STATUS", "PRIORITY", "DUE", "TEXT")
	fmt.Println()

	n := time.Now()
	today := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local)

	for _, task := range tasks {
		isDone := task.Status == storage.StatusDone
		isOverdue := task.Due != nil && task.Due.Before(today)
		isHighPriority := task.Priority == storage.PriorityHigh
		isLowPriority := task.Priority == storage.PriorityLow

		status := "[ ]"
		if isDone {
			status = "[x]"
		}
		due := "-"
		if task.Due != nil {
			due = task.Due.Format(time.DateOnly)
		}
		formattedID := fmt.Sprintf(fID, strconv.Itoa(task.ID))
		formattedStatus := fmt.Sprintf(fStatus, status)
		formattedPriority := fmt.Sprintf(fPriority, task.Priority)
		formattedDue := fmt.Sprintf(fDue, due)
		formattedText := fmt.Sprintf(fText, task.Text)

		if !isDone {
			if isOverdue {
				formattedDue = red(formattedDue)
			}
			if isHighPriority {
				formattedPriority = bold(formattedPriority)
			}
			if isLowPriority {
				formattedPriority = faint(formattedPriority)
			}
		}

		formattedSlice := []string{formattedID, formattedStatus, formattedPriority, formattedDue, formattedText}
		output := strings.Join(formattedSlice, "  ")

		if isDone {
			output = faint(output)
		}

		fmt.Println(output)
	}
}
