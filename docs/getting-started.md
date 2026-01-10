# Getting Started with ghx-cli

This guide will help you install ghx-cli and run your first commands.

## Prerequisites

- Go 1.21 or later (for building from source)
- GitHub CLI (`gh`) for authentication (recommended)
- A GitHub account with access to Projects or Discussions

## Installation

### Using Go Install (Recommended)

```bash
go install github.com/roboco-io/ghx-cli/cmd/ghx@latest
```

### Building from Source

```bash
git clone https://github.com/roboco-io/ghx-cli.git
cd ghx-cli
make build
# Binary will be at ./bin/ghx
```

### Download Binary

```bash
# macOS/Linux
curl -L https://github.com/roboco-io/ghx-cli/releases/latest/download/ghx-$(uname -s)-$(uname -m) -o ghx
chmod +x ghx
sudo mv ghx /usr/local/bin/
```

## Authentication

ghx-cli uses GitHub CLI (`gh`) for authentication. If you're already logged in with `gh`, ghx-cli will automatically use your token.

### Option 1: Use GitHub CLI (Recommended)

```bash
# Login with GitHub CLI
gh auth login

# Verify ghx-cli can authenticate
ghx auth status
```

### Option 2: Environment Variable

```bash
# Set your GitHub Personal Access Token
export GITHUB_TOKEN="ghp_your_token_here"

# Or use GHX_TOKEN
export GHX_TOKEN="ghp_your_token_here"

# Verify authentication
ghx auth status
```

### Required Scopes

Your token needs these scopes:
- `repo` - Access repositories
- `project` - Access Projects v2
- `read:discussion` - Read discussions (for discussion commands)
- `write:discussion` - Write discussions (for creating/editing)

## Quick Start Examples

### Working with Projects

```bash
# List your projects
ghx project list

# List projects for an organization
ghx project list --org myorg

# View project details
ghx project view myorg/123

# Create a new project
ghx project create "Sprint Planning" --org myorg
```

### Working with Items

```bash
# List items in a project
ghx item list myorg/123

# Add an issue to a project
ghx item add myorg/123 myorg/repo#42

# Create a draft issue
ghx item add myorg/123 --draft --title "New Task" --body "Task description"
```

### Working with Discussions

```bash
# List discussions in a repository
ghx discussion list owner/repo

# View a specific discussion
ghx discussion view owner/repo 123

# Create a new discussion
ghx discussion create owner/repo --category general --title "Question" --body "How do I...?"

# Add a comment
ghx discussion comment owner/repo 123 --body "Thanks for the help!"
```

### Managing Fields

```bash
# List fields in a project
ghx field list myorg/123

# Create a priority field
ghx field create myorg/123 "Priority" single_select --options "High,Medium,Low"

# Create a date field
ghx field create myorg/123 "Due Date" date
```

### Managing Views

```bash
# List views in a project
ghx view list myorg/123

# Create a board view
ghx view create myorg/123 "Sprint Board" board

# Create a table view with filter
ghx view create myorg/123 "Bugs" table --filter "label:bug"
```

## Output Formats

ghx-cli supports multiple output formats:

```bash
# Default table format
ghx project list

# JSON format (for scripting)
ghx project list --format json

# YAML format
ghx project list --format yaml
```

## Getting Help

```bash
# General help
ghx --help

# Command-specific help
ghx project --help
ghx project create --help

# Version information
ghx --version
```

## Next Steps

- Read the [Command Reference](commands/README.md) for detailed command documentation
- Check out [Examples & Cookbook](examples.md) for common workflows
- Configure your settings in [Configuration](configuration.md)
