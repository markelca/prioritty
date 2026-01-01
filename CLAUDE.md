# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Prioritty is a TUI and CLI application for managing tasks and notes. Written in Go 1.23, it uses SQLite for persistence and Bubble Tea for the terminal UI.

## Build Commands

```bash
# Build the binary (outputs to bin/pt)
make build

# Run directly
go run cmd/*.go

# Run with demo data (temporary database)
./bin/pt --demo
```

## Architecture

### Layered Structure

```
CLI Layer (internal/cli/)
    ↓
Service Layer (pkg/items/service/)
    ↓
Repository Layer (pkg/items/repository/)
    ↓
SQLite Database
```

### Key Directories

- `cmd/` - Entry point (`main.go`)
- `internal/cli/` - Cobra command handlers (one file per command)
- `internal/config/` - Viper-based YAML configuration
- `internal/tui/` - Bubble Tea TUI (model, view, update, keys, styles)
- `internal/migrations/` - Embedded SQL schema and seed data
- `pkg/items/` - Domain models (Task, Note, Item interface)
- `pkg/items/service/` - Business logic layer
- `pkg/items/repository/` - SQLite data access

### Domain Models

- **Item**: Base interface implemented by Task and Note
- **Task**: Has Status (Todo, InProgress, Done, Cancelled)
- **Note**: Simple content container without status
- **Tag**: Can be assigned to items for grouping

### Configuration

Config file at `~/.config/prioritty/prioritty.yaml`:
- `database_path` - SQLite database location
- `default_command` - Command run when no subcommand given (default: "tui")
- `log_file_path` - Log file location
- `editor` - External editor for editing items (default: "nano")

Environment variables use `PRIORITTY_` prefix.

## Key Patterns

- **Repository Pattern**: `SQLiteRepository` implements `TaskRepository` and `NoteRepository` interfaces
- **Service Composition**: `Service` struct composes `TaskService` and `NoteService`
- **Embedded SQL**: Migration files use `//go:embed` directive
- **Editor Integration**: `pkg/editor/` creates temp files for external editor editing
- **Item Sorting**: Items with tags appear first, then sorted by creation date (newest first)

## TUI Keybindings

- Navigation: `↑/k`, `↓/j`, `←/h`, `→/l`
- Status: `p` (in progress), `d` (done), `t` (todo), `c` (cancelled)
- Actions: `e` (edit), `s` (show), `a` (add), `r` (remove)
- UI: `?` (help), `q` (quit)

## CLI Commands

Main commands: `task`, `note`, `list`, `show`, `edit`, `remove`, `tui`
Status commands: `done`, `todo`, `start`, `cancel`
Tag commands: `tag`, `tags`
