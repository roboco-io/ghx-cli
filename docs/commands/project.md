# ghx project

Manage GitHub Projects v2.

## Synopsis

```bash
ghx project <command> [arguments] [flags]
```

## Description

The `project` command group provides comprehensive project management capabilities for GitHub Projects v2, including creating, viewing, editing, and deleting projects, as well as export/import and workflow automation.

## Commands

| Command | Description |
|---------|-------------|
| `list` | List projects |
| `view` | View project details |
| `create` | Create a new project |
| `edit` | Edit project properties |
| `delete` | Delete a project |
| `export` | Export project data |
| `import` | Import project data |
| `link` | Link project to repository |
| `template` | Manage project templates |
| `workflow` | Manage project workflows |

## ghx project list

List projects for a user or organization.

```bash
ghx project list [owner] [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--org` | List organization projects | - |
| `--user` | List user projects | - |
| `-L, --limit` | Maximum number of projects | 20 |
| `--format` | Output format (table, json, yaml) | table |

### Examples

```bash
# List your projects
ghx project list

# List organization projects
ghx project list --org myorg

# List user projects
ghx project list octocat

# JSON output
ghx project list --format json
```

## ghx project view

View project details.

```bash
ghx project view <project-ref> [flags]
```

### Arguments

- `<project-ref>` - Project reference in format `owner/number` or project ID

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format (table, json, yaml) | table |

### Examples

```bash
# View by owner/number
ghx project view myorg/123

# View by project ID
ghx project view PVT_kwDOBcH12s4AXxyz
```

## ghx project create

Create a new project.

```bash
ghx project create <title> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--org` | Create in organization | - |
| `--user` | Create for user | - |
| `--description` | Project description | - |
| `--public` | Make project public | false |

### Examples

```bash
# Create user project
ghx project create "My Project"

# Create organization project
ghx project create "Sprint 1" --org myorg

# Create with description
ghx project create "Q1 Planning" --org myorg --description "Planning for Q1 2024"

# Create public project
ghx project create "Public Roadmap" --org myorg --public
```

## ghx project edit

Edit project properties.

```bash
ghx project edit <project-ref> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--title` | New project title |
| `--description` | New description |
| `--public` | Make public (true/false) |
| `--readme` | Update README content |

### Examples

```bash
# Update title
ghx project edit myorg/123 --title "New Title"

# Update description
ghx project edit myorg/123 --description "Updated description"

# Make project public
ghx project edit myorg/123 --public true
```

## ghx project delete

Delete a project.

```bash
ghx project delete <project-ref> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--force` | Skip confirmation prompt |

### Examples

```bash
# Delete with confirmation
ghx project delete myorg/123

# Delete without confirmation
ghx project delete myorg/123 --force
```

## ghx project export

Export project data to a file.

```bash
ghx project export <project-ref> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output file path | stdout |
| `--format` | Export format (json, csv, yaml) | json |
| `--include-items` | Include project items | true |
| `--include-fields` | Include field definitions | true |
| `--include-views` | Include view configurations | true |

### Examples

```bash
# Export to JSON file
ghx project export myorg/123 -o project.json

# Export as CSV
ghx project export myorg/123 -o project.csv --format csv

# Export only structure (no items)
ghx project export myorg/123 -o template.json --include-items=false
```

## ghx project import

Import project data from a file.

```bash
ghx project import <project-ref> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-i, --input` | Input file path | - |
| `--strategy` | Import strategy (merge, replace, append) | merge |
| `--dry-run` | Preview changes without applying | false |

### Examples

```bash
# Import with merge strategy
ghx project import myorg/123 -i project.json

# Import with replace strategy
ghx project import myorg/123 -i project.json --strategy replace

# Preview import
ghx project import myorg/123 -i project.json --dry-run
```

## ghx project link

Link a project to a repository.

```bash
ghx project link <project-ref> <repository> [flags]
```

### Examples

```bash
# Link project to repository
ghx project link myorg/123 myorg/my-repo
```

## ghx project template

Manage project templates.

```bash
ghx project template <action> [flags]
```

### Actions

- `list` - List available templates
- `create` - Create template from project
- `apply` - Apply template to project

### Examples

```bash
# List templates
ghx project template list --org myorg

# Create template from project
ghx project template create myorg/123 --name "Sprint Template"

# Apply template
ghx project template apply myorg/456 --template "Sprint Template"
```

## ghx project workflow

Manage project workflows and automation.

```bash
ghx project workflow <action> [flags]
```

### Actions

- `list` - List workflows
- `create` - Create workflow
- `enable` - Enable workflow
- `disable` - Disable workflow
- `delete` - Delete workflow

### Examples

```bash
# List workflows
ghx project workflow list myorg/123

# Create auto-archive workflow
ghx project workflow create myorg/123 \
  --name "Auto Archive" \
  --trigger "status:done" \
  --action "archive"

# Disable workflow
ghx project workflow disable myorg/123 workflow-id
```
