# prioritty
A Terminal User Interface (TUI) and CLI application for managing your tasks. Focused on:
- Good looks
- Performance
- Nice defaults
- Customization
- Autocompletion support

**🚧 Disclaimer: This project is still under development.**
Home screen             |  Task content
:-------------------------:|:-------------------------:
![image](https://github.com/user-attachments/assets/d9a74cb8-e64e-4d16-8b4f-43833a9d5067) | ![image](https://github.com/user-attachments/assets/8096063f-b35a-4c7e-88e3-881d7a1bd9e3)


---

## Configuration
You can configure the tool with the `config` command, or by modifying the configuration yaml file. You can also provide a filepath to use another config (`pt --config ./config.yaml`)
`config.yaml:`
```yaml
database_path: "./data/prioritty.db"
log_file_path: "./logs/prioritty.log"
default_command: "tui"
```

## Usage
### CLI
Run the `help` command to find out the usage:
```
pt help
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
  add         Adds a new task
  completion  Generate the autocompletion script for the specified shell
  config      Show current configuration
  help        Help about any command
  list        Shows all the tasks
  remove      Removes a task by ID
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
