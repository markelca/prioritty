# prioritty
A Terminal User Interface (TUI) and CLI application for managing your tasks. Focused on:
- Good looks
- Performance
- Nice defaults
- Customization
- Autocompletion support

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
repository_type: sqlite  # or "obsidian"
```

### Repository Types

Prioritty supports two storage backends:

- **sqlite** (default): Traditional SQLite database storage
- **obsidian**: Store items as markdown files in an Obsidian vault

To use an Obsidian vault:
```yaml
repository_type: obsidian
database_path: /path/to/your/obsidian-vault
```

Or via environment variables:
```bash
export PRIORITTY_REPOSITORY_TYPE=obsidian
export PRIORITTY_DATABASE_PATH=/path/to/your/obsidian-vault
```

When using the Obsidian backend, each task/note is stored as a markdown file with YAML frontmatter in the vault root directory.

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

Usage:
  pt [flags]
  pt [command]

Available Commands:
  cancel      Mark tasks as cancelled
  completion  Generate the autocompletion script for the specified shell
  config      Show current configuration
  done        Mark tasks as done
  edit        Edit a task or note by index
  help        Help about any command
  list        Shows all the tasks
  note        Adds a new note
  remove      Removes one or more tasks by ID
  show        Show task or note details by index
  start       Mark tasks as in progress
  tag         Sets the tag for one or more tasks
  tags        Lists all available tags
  task        Adds a new task
  todo        Mark tasks as todo
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

### Frontmatter Syntax

When creating or editing items (via `pt task`, `pt note`, `pt edit`, or pressing `e`/`a` in the TUI), the editor opens with YAML frontmatter format:

**Task example:**
```yaml
---
title: Complete project report
type: task
status: in-progress
tag: work
---
Optional body/description here.
Can span multiple lines.
```

**Note example:**
```yaml
---
title: Meeting notes
type: note
tag: work
---
Note content here.
```

**Supported fields:**
| Field | Description | Values |
|-------|-------------|--------|
| `title` | Item title (required) | Any text |
| `type` | Item type | `task` or `note` |
| `status` | Task status (tasks only) | `todo`, `in-progress`, `done`, `cancelled` |
| `tag` | Single tag name | Any text |

You can view an item's raw frontmatter with `pt show <index> --raw`.

### Autocompletion

To enable shell autocompletion for `pt`, add the appropriate line to your shell configuration:

#### Bash
```bash
# Add to your ~/.bashrc
source <(pt completion bash)
```

#### Zsh
```bash
# Add to your ~/.zshrc
source <(pt completion zsh)
```

#### Fish
```bash
# Add to your fish config
source (pt completion fish | psub)
```

#### PowerShell
```powershell
# Add to your PowerShell profile
Invoke-Expression (pt completion powershell | Out-String)
```


---
Inspired by [taskbook](https://github.com/klaudiosinani/taskbook)
