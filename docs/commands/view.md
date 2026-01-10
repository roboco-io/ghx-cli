# ghx view

Manage project views (table, board, roadmap).

## Synopsis

```bash
ghx view <command> [arguments] [flags]
```

## Description

The `view` command group provides comprehensive view management for GitHub Projects. Views allow you to organize and visualize project items in different layouts.

## View Layouts

| Layout | Description | Use Case |
|--------|-------------|----------|
| `table` | Spreadsheet-style table | Detailed data view |
| `board` | Kanban-style board | Workflow visualization |
| `roadmap` | Timeline view | Planning and scheduling |

## Commands

| Command | Description |
|---------|-------------|
| `list` | List project views |
| `create` | Create a new view |
| `update` | Update view properties |
| `copy` | Copy an existing view |
| `delete` | Delete a view |
| `sort` | Configure view sorting |
| `group` | Configure view grouping |

## ghx view list

List views in a project.

```bash
ghx view list <project-ref> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format (table, json) | table |

### Examples

```bash
# List all views
ghx view list myorg/123

# JSON output
ghx view list myorg/123 --format json
```

## ghx view create

Create a new project view.

```bash
ghx view create <project-ref> <name> <layout> [flags]
```

### Arguments

- `<project-ref>` - Project reference in format `owner/number`
- `<name>` - View name
- `<layout>` - View layout (table, board, roadmap)

### Flags

| Flag | Description |
|------|-------------|
| `--filter` | View filter expression |
| `--group-by` | Field to group by |
| `--sort-by` | Field to sort by |

### Examples

```bash
# Create table view
ghx view create myorg/123 "All Items" table

# Create board view
ghx view create myorg/123 "Sprint Board" board

# Create roadmap view
ghx view create myorg/123 "Timeline" roadmap

# Create view with filter
ghx view create myorg/123 "Bugs" table --filter "label:bug"

# Create board grouped by status
ghx view create myorg/123 "Kanban" board --group-by Status
```

## ghx view update

Update view properties.

```bash
ghx view update <view-id> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--name` | New view name |
| `--filter` | New filter expression |

### Examples

```bash
# Rename view
ghx view update PVV_xxx --name "Sprint 2 Board"

# Update filter
ghx view update PVV_xxx --filter "status:open"
```

## ghx view copy

Create a copy of an existing view.

```bash
ghx view copy <view-id> <new-name> [flags]
```

### Examples

```bash
# Copy view
ghx view copy PVV_xxx "Sprint 2 Board"
```

## ghx view delete

Delete a project view.

```bash
ghx view delete <view-id> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--force` | Skip confirmation prompt |

### Examples

```bash
# Delete with confirmation
ghx view delete PVV_xxx

# Delete without confirmation
ghx view delete PVV_xxx --force
```

## ghx view sort

Configure view sorting.

```bash
ghx view sort <view-id> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--field` | Field to sort by | - |
| `--direction` | Sort direction (asc, desc) | asc |
| `--clear` | Clear sorting | false |

### Examples

```bash
# Sort by priority descending
ghx view sort PVV_xxx --field Priority --direction desc

# Sort by due date
ghx view sort PVV_xxx --field "Due Date" --direction asc

# Clear sorting
ghx view sort PVV_xxx --clear
```

## ghx view group

Configure view grouping.

```bash
ghx view group <view-id> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--field` | Field to group by | - |
| `--direction` | Group direction (asc, desc) | asc |
| `--clear` | Clear grouping | false |

### Examples

```bash
# Group by status
ghx view group PVV_xxx --field Status

# Group by assignee descending
ghx view group PVV_xxx --field Assignee --direction desc

# Clear grouping
ghx view group PVV_xxx --clear
```

## Filter Expressions

Views support filter expressions to show only matching items:

| Filter | Description | Example |
|--------|-------------|---------|
| `label:name` | Filter by label | `label:bug` |
| `status:value` | Filter by status | `status:open` |
| `assignee:user` | Filter by assignee | `assignee:octocat` |
| `milestone:name` | Filter by milestone | `milestone:v1.0` |
| `no:label` | Items without labels | `no:label` |
| `is:issue` | Only issues | `is:issue` |
| `is:pr` | Only pull requests | `is:pr` |

Multiple filters can be combined:

```bash
ghx view create myorg/123 "My Bugs" table --filter "label:bug assignee:@me status:open"
```
