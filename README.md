# MacroFlow CLI

A powerful command-line tool for creating directory-scoped command aliases (macros). Define shortcuts for complex commands that work within specific project directories and their subdirectories.

> **Note**: This is the CLI application repository. For documentation, see [MacroFlowDocs](https://github.com/anish/MacroFlowDocs).

## Features

- 📁 **Directory-Scoped**: Macros are tied to specific directories and work in all subdirectories
- 🎯 **Project-Based**: Organize macros into projects for better management
- 🔄 **Parameter Support**: Use `$1`, `$2`, etc. or `$@` for dynamic command parameters
- 💾 **Import/Export**: Share macro configurations or back them up
- 🔍 **Smart Lookup**: Automatically finds the right project for your current directory
- 🗂️ **Easy Management**: List, delete, and organize macros with simple commands

## Installation

### macOS (Homebrew)

```bash
brew tap anish/macroflow
brew install macroflow
```

### From Source

```bash
git clone https://github.com/anish/macroflow.git
cd MacroFlowCLI
go build -o macroflow
sudo mv macroflow /usr/local/bin/macro
```

## Quick Start

1. **Initialize a project** in your working directory:
```bash
cd ~/projects/my-app
macro init
```

2. **Add some macros**:
```bash
macro add dev "npm run dev"
macro add build "npm run build"
macro add goto "cd $1"
macro add commit "git commit -m \"$@\""
```

3. **Use your macros**:
```bash
macro dev                    # Runs: npm run dev
macro goto src/components    # Runs: cd src/components
macro commit Fix bug         # Runs: git commit -m "Fix bug"
```

## Usage

### Project Management

#### Initialize a new project
```bash
macro init [name]
```
Creates a new macro project in the current directory. If no name is provided, uses the directory name.

#### List all projects
```bash
macro list --projects
```

#### Delete a project
```bash
macro delete <project-id> --project
```

### Macro Management

#### Add a macro
```bash
macro add <name> <command>
```

Examples:
```bash
# Simple alias
macro add dev "npm run dev"

# With single parameter
macro add goto "cd $1"

# With multiple parameters
macro add copy "cp $1 $2"

# With all parameters
macro add commit "git commit -m \"$@\""
```

#### List macros

```bash
# List macros in current project
macro list

# List all macros across all projects
macro list --all
```

#### Delete macros

```bash
# Delete by ID
macro delete <macro-id>

# Delete multiple macros
macro delete <id1> <id2> <id3>

# Delete by name (in current project)
macro delete-name <name>
```

### Import/Export

#### Export all data
```bash
macro export --file macros.json
```

#### Import data
```bash
# Replace all existing data
macro import --file macros.json

# Merge with existing data
macro import --file macros.json --merge
```

## Parameter Substitution

MacroFlow supports flexible parameter substitution:

- `$1`, `$2`, `$3`, etc. - Individual positional parameters
- `$@` - All parameters as a single string

### Examples

```bash
# Add a macro with parameters
macro add goto "cd $1"
macro add search "grep -r \"$1\" $2"
macro add commit "git commit -m \"$@\""

# Use them
macro goto src/              # → cd src/
macro search TODO ./src      # → grep -r "TODO" ./src
macro commit Fix the bug     # → git commit -m "Fix the bug"
```

## How It Works

### Directory Scoping

When you run a macro, MacroFlow:
1. Checks your current directory
2. Finds the deepest (most specific) project that matches your location
3. Uses macros from that project

This means:
- Macros work in the project directory and ALL subdirectories
- Child directories inherit macros from parent projects
- You can have nested projects with different macros

### Example Directory Structure

```
~/projects/
├── web-app/                    # Project A
│   ├── frontend/              # Inherits Project A macros
│   │   └── src/              # Still uses Project A macros
│   └── backend/               # Project B (has own macros)
│       └── api/              # Uses Project B macros
└── cli-tool/                  # Project C
```

### Data Storage

MacroFlow stores all data in `~/.macroflow/macroflow.json`:
- Projects (with IDs, names, and paths)
- Macros (with IDs, commands, and project associations)

## Command Reference

| Command | Description |
|---------|-------------|
| `macro init [name]` | Initialize a project |
| `macro add <name> <command>` | Add a macro |
| `macro <name> [args...]` | Execute a macro |
| `macro list` | List macros in current project |
| `macro list --all` | List all macros |
| `macro list --projects` | List all projects |
| `macro delete <id> [ids...]` | Delete macro(s) by ID |
| `macro delete <id> --project` | Delete a project |
| `macro delete-name <name>` | Delete macro by name |
| `macro export --file <path>` | Export data to JSON |
| `macro import --file <path>` | Import data from JSON |
| `macro import --file <path> --merge` | Merge imported data |

## Tips & Tricks

### 1. Common Development Workflows
```bash
macro add dev "npm run dev"
macro add test "npm test"
macro add build "npm run build"
macro add lint "npm run lint --fix"
```

### 2. Git Shortcuts
```bash
macro add s "git status"
macro add c "git commit -m \"$@\""
macro add p "git push"
macro add pl "git pull"
```

### 3. Directory Navigation
```bash
macro add src "cd src"
macro add back "cd .."
macro add root "cd \$(git rev-parse --show-toplevel)"
```

### 4. Docker Commands
```bash
macro add dup "docker-compose up -d"
macro add ddown "docker-compose down"
macro add dlogs "docker-compose logs -f $1"
```

### 5. Complex Workflows
```bash
macro add deploy "npm run build && npm run test && git push heroku main"
macro add fresh "rm -rf node_modules package-lock.json && npm install"
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - See LICENSE file for details

## Related Repositories

- **[MacroFlowDocs](https://github.com/anish/MacroFlowDocs)** - Official documentation
- **MacroFlowGUI** (Coming soon) - GUI application

## Author

Built with ❤️ by Anish
