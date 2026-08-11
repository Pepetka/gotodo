package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/pepetka/gotodo/internal/storage"
)

func main() {
	if err := run(); err != nil {
		red := color.New(color.FgRed)
		fmt.Fprintln(os.Stderr, red.Sprint(err))
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		return errors.New("missing command: add|list|done|rm|clear")
	}
	cmd, args := args[0], args[1:]

	command, err := ParseCommand(cmd)
	if err != nil {
		return err
	}

	path, err := preparePath()
	if err != nil {
		return err
	}

	store, err := storage.LoadStore(path)
	if err != nil {
		return err
	}

	needWrite := false
	switch command {
	case CommandAdd:
		id, addErr := add(store, args)
		if addErr != nil {
			return addErr
		}
		needWrite = true
		color.Green("added #%d", id)
	case CommandList:
		tasks, listErr := list(store, args)
		if listErr != nil {
			return listErr
		}
		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return nil
		}
		format(tasks, store.LastID)
	case CommandDone:
		id, doneErr := done(store, args)
		if doneErr != nil {
			return doneErr
		}
		needWrite = true
		color.Green("done #%d", id)
	case CommandRm:
		id, rmErr := rm(store, args)
		if rmErr != nil {
			return rmErr
		}
		needWrite = true
		color.Green("removed #%d", id)
	case CommandClear:
		tasks, clearErr := clear(store, args)
		if errors.Is(clearErr, errAborted) {
			return nil
		}
		if clearErr != nil {
			return clearErr
		}
		needWrite = true
		color.Green("removed %d tasks", len(tasks))
	}

	if !needWrite {
		return nil
	}
	return store.SaveStore(path)
}

func preparePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, ".gotodo", "tasks.json")

	return path, nil
}
