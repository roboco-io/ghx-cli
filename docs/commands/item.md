# ghx item

Manage project items (issues, pull requests, and draft issues).

## Synopsis

```bash
ghx item <command> [arguments] [flags]
```

## Description

The `item` command group provides comprehensive item management for GitHub Projects. Items can be existing issues, pull requests, or draft issues created directly in the project.

## Commands

| Command | Description |
|---------|-------------|
| `list` | List items in a repository or project |
| `view` | View item details |
| `add` | Add item to project |
| `add-bulk` | Add multiple items at once |
| `edit` | Edit item field values |
| `remove` | Remove item from project |
| `update-bulk` | Update multiple items at once |

## ghx item list

List issues and pull requests.

```bash
ghx item list <owner/repo> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--state` | Filter by state (open, closed, all) | all |
| `--type` | Filter by type (issue, pr) | - |
| `--label` | Filter by label | - |
| `--assignee` | Filter by assignee | - |
| `--milestone` | Filter by milestone | - |
| `-L, --limit` | Maximum number of items | 30 |
| `--format` | Output format (table, json) | table |

### Examples

```bash
# List all items
ghx item list myorg/repo

# List open issues only
ghx item list myorg/repo --state open --type issue

# Filter by label
ghx item list myorg/repo --label bug

# Filter by assignee
ghx item list myorg/repo --assignee octocat
```

## ghx item view

View details of an issue or pull request.

```bash
ghx item view <item-ref> [flags]
```

### Arguments

- `<item-ref>` - Item reference in format `owner/repo#number`

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format (table, json) | table |
| `--comments` | Include comments | false |

### Examples

```bash
# View issue
ghx item view myorg/repo#123

# View with comments
ghx item view myorg/repo#123 --comments

# JSON output
ghx item view myorg/repo#123 --format json
```

## ghx item add

Add an item to a project.

```bash
ghx item add <project-ref> <item-ref> [flags]
ghx item add <project-ref> --draft --title <title> [flags]
```

### Arguments

- `<project-ref>` - Project reference in format `owner/number`
- `<item-ref>` - Item reference in format `owner/repo#number` (for existing items)

### Flags

| Flag | Description |
|------|-------------|
| `--draft` | Create a draft issue |
| `--title` | Draft issue title |
| `--body` | Draft issue body |

### Examples

```bash
# Add existing issue to project
ghx item add myorg/123 myorg/repo#42

# Add pull request to project
ghx item add myorg/123 myorg/repo#100

# Create draft issue in project
ghx item add myorg/123 --draft --title "New Task" --body "Task description"
```

## ghx item add-bulk

Add multiple items to a project at once.

```bash
ghx item add-bulk <project-ref> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--items` | Comma-separated list of item references |
| `--file` | File with item references (one per line) |
| `--query` | GitHub search query to find items |

### Examples

```bash
# Add multiple items
ghx item add-bulk myorg/123 --items "repo#1,repo#2,repo#3"

# Add from file
ghx item add-bulk myorg/123 --file items.txt

# Add items matching search query
ghx item add-bulk myorg/123 --query "is:issue label:priority-high"
```

## ghx item edit

Edit item field values in a project.

```bash
ghx item edit <project-ref> <item-id> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--field` | Field name to update |
| `--value` | New value for the field |

### Examples

```bash
# Set status field
ghx item edit myorg/123 PVTI_xxx --field Status --value "In Progress"

# Set priority field
ghx item edit myorg/123 PVTI_xxx --field Priority --value High

# Set date field
ghx item edit myorg/123 PVTI_xxx --field "Due Date" --value "2024-01-31"
```

## ghx item remove

Remove an item from a project.

```bash
ghx item remove <project-ref> <item-id> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--force` | Skip confirmation prompt |

### Examples

```bash
# Remove with confirmation
ghx item remove myorg/123 PVTI_xxx

# Remove without confirmation
ghx item remove myorg/123 PVTI_xxx --force
```

## ghx item update-bulk

Update multiple project items at once.

```bash
ghx item update-bulk <project-ref> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--items` | Comma-separated list of item IDs |
| `--field` | Field name to update |
| `--value` | New value for the field |
| `--filter` | Filter items to update |

### Examples

```bash
# Update multiple items
ghx item update-bulk myorg/123 --items "id1,id2,id3" --field Status --value Done

# Update items matching filter
ghx item update-bulk myorg/123 --filter "Status:Todo" --field Status --value "In Progress"
```

## Item References

ghx-cli supports multiple formats for referencing items:

| Format | Example | Description |
|--------|---------|-------------|
| `owner/repo#number` | `myorg/repo#42` | Standard reference |
| `repo#number` | `repo#42` | Repo in current org |
| `#number` | `#42` | Issue in current repo |
| GitHub URL | `https://github.com/myorg/repo/issues/42` | Full URL |
