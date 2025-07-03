# tada - Terminal Todo Application

A modern, fast terminal-based todo application written in Go that helps you manage tasks efficiently from the command line.

## Features

- ✅ **CRUD Operations**: Create, read, update, and delete todos
- 🎯 **Priority Levels**: Low, medium, and high priority tasks
- ⏱️ **Task Age Tracking**: See how long tasks have been open or completed
- 🔍 **Flexible Filtering**: Filter by status and priority
- 📊 **Rich Display**: Colorful output with emojis and detailed information
- 💾 **SQLite Storage**: Fast, reliable local database storage
- 🚀 **Command Aliases**: Convenient shortcuts for common operations

## Installation

`tada` can be installed via Homebrew.

First, add the Homebrew tap (you only need to do that once):

```sh
brew tap tada/tap git@github.com:negadras/tada.git
```

Then, install it via:

```sh
brew install tada
```

Or install the latest HEAD via:

```sh
brew install --HEAD tada
```

To verify it, try running the `help` command:

```sh
tada help
```

## Usage

### Basic Commands

```bash
# Show all available commands
tada --help

# Add a new todo
tada add "Your task description"

# List todos (shows open tasks by default)
tada list

# Update a todo
tada update [id] --status done

# Delete a todo
tada delete [id]

# Quick shortcuts
tada done [id]    # Mark as done
tada open [id]    # Mark as open
```

### Adding Todos

```bash
# Add a todo with medium priority (default)
tada add "Complete the project"

# Add a todo with specific priority
tada add "Fix critical bug" --priority high
tada add "Clean up code" --priority low

# Using short flags
tada add "Review pull request" -p high
```

### Listing and Filtering Todos

```bash
# List all open todos (default)
tada list

# List all todos regardless of status
tada list --status all

# List only completed todos
tada list --status done

# List high priority tasks
tada list --priority high

# List all high priority tasks (open and done)
tada list --status all --priority high

# Using short flags
tada list -s done -p high
```

### Updating Todos

```bash
# Mark a todo as done
tada update 5 --status done

# Change priority
tada update 5 --priority high

# Update description
tada update 5 --description "Updated task description"

# Update multiple properties at once
tada update 5 --status done --priority low

# Using short flags
tada update 5 -s done -p high -d "New description"
```

### Command Aliases

```bash
# List todos
tada ls           # Same as 'tada list'

# Delete todos
tada rm 5         # Same as 'tada delete 5'
tada del 5        # Same as 'tada delete 5'

# Quick status updates
tada done 5       # Same as 'tada update 5 --status done'
tada open 5       # Same as 'tada update 5 --status open'

# Show all aliases
tada aliases
```

## Priority Levels

- **🟢 Low**: Regular tasks with no urgency
- **🟡 Medium**: Standard priority tasks (default)
- **🔴 High**: Important tasks that need attention

## Status Types

- **Open**: New or pending tasks (default)
- **Done**: Completed tasks

## Data Storage

Your todos are automatically saved to `~/.tada/todos.db` in your home directory using SQLite. The database is created
automatically when you add your first todo.

## Examples

Here are some practical examples:

```bash
# Personal task management
tada add "Doctor appointment" --priority high
tada add "Call mom" --priority medium
tada add "Read book" --priority low

# Work task management
tada add "Sprint planning meeting" --priority high
tada add "Code review" --priority medium
tada add "Update documentation" --priority low

# View work progress
tada list --status all    # See all tasks
tada list --priority high # Focus on high priority

# Mark tasks complete
tada done 3               # Quick completion
tada update 5 --status done --priority low  # Update multiple fields

# Clean up completed tasks
tada list --status done   # Review completed work
tada delete 7             # Remove unnecessary tasks
```

## Command Reference

| Command  | Description                        | Example                           |
|----------|------------------------------------|-----------------------------------|
| `add`    | Create a new todo                  | `tada add "Task" --priority high` |
| `list`   | Show todos with optional filtering | `tada list --status done`         |
| `update` | Modify an existing todo            | `tada update 1 --status done`     |
| `delete` | Remove a todo                      | `tada delete 1`                   |
| `done`   | Mark todo as completed             | `tada done 1`                     |
| `open`   | Mark todo as open                  | `tada open 1`                     |
| `ls`     | Alias for list                     | `tada ls`                         |
| `rm`     | Alias for delete                   | `tada rm 1`                       |
| `del`    | Alias for delete                   | `tada del 1`                      |

## Tips

1. **Use priorities** to organize your tasks by importance
2. **Filter by status** to focus on open work or review completed tasks
3. **Use aliases** for faster command entry (`ls`, `rm`, `done`)
4. **Check task age** to see how long items have been pending
5. **Update multiple fields** at once for efficiency
6. **Use the done command** as a quick shortcut to mark tasks complete

## Contributing

If you are interested in contributing to `tada`, have ideas for features that could be a useful addition or just want to
ask a question, please reach out or open an issue on the GitHub repository.
