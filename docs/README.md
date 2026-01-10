# ghx-cli Documentation

Welcome to the ghx-cli documentation. This guide covers all features and commands of the GitHub eXtensions CLI.

## Table of Contents

### Getting Started
- [Installation & Quick Start](getting-started.md) - Install ghx-cli and run your first commands

### Command Reference
- [Command Overview](commands/README.md) - All available commands at a glance
- [project](commands/project.md) - Manage GitHub Projects v2
- [item](commands/item.md) - Manage project items (issues, PRs, drafts)
- [field](commands/field.md) - Manage custom fields
- [view](commands/view.md) - Manage project views (table, board, roadmap)
- [discussion](commands/discussion.md) - Manage GitHub Discussions
- [analytics](commands/analytics.md) - Generate reports and bulk operations
- [auth](commands/auth.md) - Manage authentication

### Guides
- [Configuration](configuration.md) - Configure ghx-cli settings
- [Examples & Cookbook](examples.md) - Common use cases and workflows
- [Troubleshooting](troubleshooting.md) - Common issues and solutions

### Reference
- [Feature Comparison](feature-comparison.md) - ghx-cli vs gh CLI capabilities
- [PRD](PRD.md) - Product Requirements Document

## Quick Links

```bash
# Check authentication status
ghx auth status

# List your projects
ghx project list

# List discussions in a repository
ghx discussion list owner/repo

# Get help for any command
ghx [command] --help
```

## Support

- [GitHub Issues](https://github.com/roboco-io/ghx-cli/issues) - Report bugs and request features
- [GitHub Discussions](https://github.com/roboco-io/ghx-cli/discussions) - Ask questions and share ideas
