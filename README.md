# MacroFlow CLI

A powerful command-line tool for creating directory-scoped command aliases (macros). Define shortcuts for complex commands that work within specific project directories and their subdirectories.

## Features

- **Directory-Scoped**: Macros are tied to specific directories and work in all subdirectories
- **Project-Based**: Organize macros into projects for better management
- **Parameter Support**: Use `$1`, `$2`, etc. or `$@` for dynamic command parameters
- **Import/Export**: Share macro configurations or back them up
- **Smart Lookup**: Automatically finds the right project for your current directory
- **Easy Management**: List, delete, and organize macros with simple commands

## Installation

### macOS (Homebrew)

```bash
brew tap <user>/macroflow
brew install macroflow
```

### From Source

```bash
git clone https://github.com/anishdpatel28/MacroFlowCLI.git
cd MacroFlowCLI
go build -o macro main.go
sudo cp macro /usr/local/bin/macro
```

## Quick Start

```bash
# Initialize a project
cd ~/projects/my-app
macro init

# Add macros
macro add dev "npm run dev"
macro add build "npm run build"
macro add test "npm test"

# Use them
macro dev
```

## Usage

### Commands

```bash
macro init [name]                    # Initialize project
macro add <name> <command>           # Add macro
macro <name> [args...]               # Execute macro
macro list                           # List macros
macro list --all                     # List all macros
macro list --projects                # List all projects
macro delete <id>                    # Delete by ID
macro delete-name <name>             # Delete by name
macro export --file <path>           # Export data
macro import --file <path>           # Import data
```

### Parameter Substitution

```bash
# Single parameter
macro add goto "cd $1"
macro goto src              # → cd src

# Multiple parameters
macro add copy "cp $1 $2"
macro copy file.txt backup.txt

# All parameters
macro add commit "git commit -m \"$@\""
macro commit Fix the bug    # → git commit -m "Fix the bug"
```

## Examples

### Web Development
```bash
macro add dev "npm run dev"
macro add build "npm run build"
macro add test "npm test"
macro add deploy "npm run build && npm run deploy"
```

### Git Shortcuts
```bash
macro add s "git status"
macro add c "git commit -m \"$@\""
macro add p "git push"
macro add save "git add . && git commit -m \"$@\" && git push"
```

### Docker
```bash
macro add up "docker-compose up -d"
macro add down "docker-compose down"
macro add logs "docker-compose logs -f $1"
```

## Documentation

For more examples and detailed guides, see [MacroFlowDocs](https://github.com/anishdpatel28/MacroFlowDocs).

## License

MIT License - See [LICENSE](LICENSE) file.

## Author

Built with ❤️ by <user>
