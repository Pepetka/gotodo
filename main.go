package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pepetka/gotodo/internal/storage"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
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
		fmt.Println(id)
	case CommandList:
		tasks, listErr := list(store, args)
		if listErr != nil {
			return listErr
		}
		b, _ := json.MarshalIndent(tasks, "", "  ")
		fmt.Println(string(b))
	case CommandDone:
		id, doneErr := done(store, args)
		if doneErr != nil {
			return doneErr
		}
		needWrite = true
		fmt.Println(id)
	case CommandRm:
		id, rmErr := rm(store, args)
		if rmErr != nil {
			return rmErr
		}
		needWrite = true
		fmt.Println(id)
	case CommandClear:
		tasks, clearErr := clear(store, args)
		if clearErr != nil {
			return clearErr
		}
		needWrite = true
		b, _ := json.MarshalIndent(tasks, "", "  ")
		fmt.Println(string(b))
	}

	if !needWrite {
		return nil
	}
	err = store.SaveStore(path)
	if err != nil {
		return err
	}
	return nil
}

func preparePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, ".gotodo", "tasks.json")

	return path, nil
}
