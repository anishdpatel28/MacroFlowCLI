# Changelog

All notable changes to MacroFlow will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of MacroFlow CLI
- Project management (`macro init`, `macro list --projects`)
- Macro management (`macro add`, `macro list`, `macro delete`)
- Macro execution with parameter substitution (`$1`, `$2`, `$@`)
- Directory-scoped project detection
- Import/Export functionality (`macro export`, `macro import`)
- Delete by name command (`macro delete-name`)
- Comprehensive documentation
- Homebrew formula template
- Cross-platform build support

### Features
- 📁 Directory-scoped macros
- 🎯 Project-based organization
- 🔄 Parameter support for dynamic commands
- 💾 JSON-based storage
- 🔍 Smart project lookup
- 📤 Import/Export capabilities

## [1.0.0] - TBD

Initial release.

### Added
- Core CLI functionality
- Storage layer with JSON persistence
- All basic commands
- Documentation and examples
- Installation guides

---

## Version History

### Versioning Strategy

- **Major** (X.0.0): Breaking changes
- **Minor** (0.X.0): New features, no breaking changes
- **Patch** (0.0.X): Bug fixes and minor improvements

### Planned Features

Future versions may include:
- Shell completion support
- Macro templates/snippets
- Remote macro sharing
- Macro groups/categories
- Environment variable support in macros
- Conditional macro execution
- Macro chaining
- Interactive mode
- GUI application (separate project)

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on our development process.
