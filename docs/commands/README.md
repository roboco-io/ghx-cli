# Command Reference

ghx-cli provides commands organized into logical groups for managing GitHub Projects and Discussions.

## Command Groups

| Command | Description |
|---------|-------------|
| [ghx project](project.md) | Manage GitHub Projects v2 |
| [ghx item](item.md) | Manage project items (issues, PRs, drafts) |
| [ghx field](field.md) | Manage custom fields |
| [ghx view](view.md) | Manage project views |
| [ghx discussion](discussion.md) | Manage GitHub Discussions |
| [ghx analytics](analytics.md) | Generate reports and bulk operations |
| [ghx auth](auth.md) | Manage authentication |

## Global Flags

These flags are available for all commands:

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Config file path | `$HOME/.ghx.yaml` |
| `--token` | GitHub Personal Access Token | - |
| `--org` | GitHub organization | - |
| `--user` | GitHub user | - |
| `--format` | Output format (table, json, yaml) | `table` |
| `--debug` | Enable debug output | `false` |
| `--no-cache` | Disable caching | `false` |
| `-h, --help` | Help for command | - |
| `-v, --version` | Version information | - |

## Command Structure

Commands follow this pattern:

```
ghx <group> <action> [arguments] [flags]
```

Examples:
```bash
ghx project list                    # List projects
ghx project view myorg/123          # View specific project
ghx item add myorg/123 repo#42      # Add item to project
ghx discussion create owner/repo    # Create discussion
```

## Output Formats

### Table (Default)

Human-readable table format:

```bash
ghx project list
```

```
NUMBER  TITLE           ITEMS  UPDATED
1       Sprint 1        42     2024-01-15
2       Backlog         128    2024-01-14
```

### JSON

Machine-readable JSON format:

```bash
ghx project list --format json
```

```json
[
  {"number": 1, "title": "Sprint 1", "items": 42},
  {"number": 2, "title": "Backlog", "items": 128}
]
```

### YAML

YAML format for configuration files:

```bash
ghx project list --format yaml
```

```yaml
- number: 1
  title: Sprint 1
  items: 42
- number: 2
  title: Backlog
  items: 128
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Authentication error |
| 4 | API error |

## Getting Help

```bash
# General help
ghx --help

# Group help
ghx project --help

# Command help
ghx project create --help
```
