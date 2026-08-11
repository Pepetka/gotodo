# gotodo

A minimal command-line todo manager written in Go. Tasks are stored locally
in a single JSON file — no database, no daemon, no accounts.

## Features

- Add tasks with priority (`high`, `medium`, `low`) and an optional due date
- List tasks with status, priority, and due-date filters
- Mark tasks as done, remove them, or clear all active tasks
- Colorized output: overdue dates in red, high priority in bold, done tasks dimmed
- Atomic saves to `~/.gotodo/tasks.json`

## Installation

```sh
go install github.com/pepetka/gotodo@latest
```

Or build from source:

```sh
git clone https://github.com/pepetka/gotodo.git
cd gotodo
go build -o gotodo
```

## Usage

```sh
gotodo <command> [arguments] [flags]
```

### add

Add a new task:

```sh
gotodo add "Buy milk"
gotodo add "Submit report" --priority high --due 2026-08-15
```

| Flag         | Default  | Description                     |
| ------------ | -------- | ------------------------------- |
| `--priority` | `medium` | `high`, `medium`, or `low`      |
| `--due`      | —        | Due date in `YYYY-MM-DD` format |

### list

List tasks (active by default):

```sh
gotodo list
gotodo list --filter all
gotodo list --filter done --priority high --before 2026-09-01
```

| Flag         | Default  | Description                             |
| ------------ | -------- | --------------------------------------- |
| `--filter`   | `active` | `all`, `active`, or `done`              |
| `--priority` | —        | Show only tasks with this priority      |
| `--before`   | —        | Show only tasks due before `YYYY-MM-DD` |

### done

Mark a task as done:

```sh
gotodo done 3
```

### rm

Remove a task permanently:

```sh
gotodo rm 3
```

### clear

Remove all active tasks at once. Asks for confirmation interactively;
use `--yes` to skip the prompt (required in non-interactive shells):

```sh
gotodo clear
gotodo clear --yes
```

## Storage

All tasks live in `~/.gotodo/tasks.json`. The file is written atomically
(write to a temp file, then rename), so an interrupted save won't corrupt
your data.

## License

[MIT](LICENSE)
