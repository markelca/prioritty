# prioritty
A Terminal User Interface (TUI) and CLI application for managing your tasks. Focused on:
- Good looks
- Performance
- Nice defaults
- Customization
- Autocompletion support

**🚧 Disclaimer: This project is still under development.**
![Peek 2025-06-09 18-58](https://github.com/user-attachments/assets/24c2bd12-a714-4d69-bc01-ba28c34b8f32)


---

## Configuration
You can configure the tool with the `config` command, or by modifying the configuration yaml file. You can also provide a filepath to use another config (`pt --config ./config.yaml`)
`config.yaml:`
```yaml
database_path: "./data/prioritty.db"
log_file_path: "./logs/prioritty.log"
default_command: "tui"
editor: vim
```

## Usage
### CLI
Run the `help` command to find out the usage:
```
A Terminal User Interface (TUI) and CLI application for managing your tasks. Focused on:
        - Good looks
        - Performance
        - Nice defaults
        - Customization
        - Autocompletion support

🚧 Disclaimer: This project is still under development.

Usage:
  pt [flags]
  pt [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Show current configuration
  help        Help about any command
  list        Shows all the tasks
  note        Adds a new note
  remove      Removes a task by ID
  tag         Sets the tag for a task
  task        Adds a new task
  tui         Launch the interactive TUI
  version     Print the version number of Hugo

Flags:
      --config string   config file (default is $HOME/.cobra.yaml)
      --demo            Use a temporal demo database with predefined values
  -h, --help            help for pt

Use "pt [command] --help" for more information about a command.
```
### TUI
You can also press the `?` key to toggle the full help in TUI mode:
![image](https://github.com/user-attachments/assets/bcc53f9c-8250-45e8-bb2d-edaaeebdbf95)


---
Inspired by [taskbook](https://github.com/klaudiosinani/taskbook)
