# ghx field

Manage custom fields in GitHub Projects.

## Synopsis

```bash
ghx field <command> [arguments] [flags]
```

## Description

The `field` command group provides comprehensive field management for GitHub Projects. Fields allow you to track additional metadata for project items with support for various data types.

## Field Types

| Type | Description | Example Values |
|------|-------------|----------------|
| `text` | Free-form text | "Implementation notes" |
| `number` | Numeric values | 5, 3.14, -10 |
| `date` | Date values | 2024-01-31 |
| `single_select` | Single choice from options | "High", "Medium", "Low" |
| `iteration` | Sprint/iteration cycles | "Sprint 1", "Sprint 2" |

## Commands

| Command | Description |
|---------|-------------|
| `list` | List project fields |
| `create` | Create a new field |
| `update` | Update field properties |
| `delete` | Delete a field |
| `add-option` | Add option to single select field |
| `update-option` | Update single select option |
| `delete-option` | Delete single select option |

## ghx field list

List fields in a project.

```bash
ghx field list <project-ref> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format (table, json) | table |

### Examples

```bash
# List all fields
ghx field list myorg/123

# JSON output
ghx field list myorg/123 --format json
```

## ghx field create

Create a new project field.

```bash
ghx field create <project-ref> <name> <type> [flags]
```

### Arguments

- `<project-ref>` - Project reference in format `owner/number`
- `<name>` - Field name
- `<type>` - Field type (text, number, date, single_select, iteration)

### Flags

| Flag | Description |
|------|-------------|
| `--options` | Comma-separated options (for single_select) |
| `--description` | Field description |

### Examples

```bash
# Create text field
ghx field create myorg/123 "Notes" text

# Create number field
ghx field create myorg/123 "Story Points" number

# Create date field
ghx field create myorg/123 "Due Date" date

# Create single select field
ghx field create myorg/123 "Priority" single_select --options "Critical,High,Medium,Low"

# Create iteration field
ghx field create myorg/123 "Sprint" iteration
```

## ghx field update

Update field properties.

```bash
ghx field update <field-id> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--name` | New field name |
| `--description` | New description |

### Examples

```bash
# Rename field
ghx field update PVTF_xxx --name "New Priority"

# Update description
ghx field update PVTF_xxx --description "Task priority level"
```

## ghx field delete

Delete a project field.

```bash
ghx field delete <field-id> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--force` | Skip confirmation prompt |

### Examples

```bash
# Delete with confirmation
ghx field delete PVTF_xxx

# Delete without confirmation
ghx field delete PVTF_xxx --force
```

## ghx field add-option

Add an option to a single select field.

```bash
ghx field add-option <field-id> <option-name> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--color` | Option color (e.g., red, blue, green) |
| `--description` | Option description |

### Examples

```bash
# Add basic option
ghx field add-option PVTF_xxx "Critical"

# Add option with color
ghx field add-option PVTF_xxx "Urgent" --color red

# Add option with description
ghx field add-option PVTF_xxx "Blocked" --color yellow --description "Waiting for dependencies"
```

## ghx field update-option

Update a single select field option.

```bash
ghx field update-option <field-id> <option-id> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--name` | New option name |
| `--color` | New option color |
| `--description` | New description |

### Examples

```bash
# Rename option
ghx field update-option PVTF_xxx OPT_xxx --name "Very High"

# Change color
ghx field update-option PVTF_xxx OPT_xxx --color purple
```

## ghx field delete-option

Delete a single select field option.

```bash
ghx field delete-option <field-id> <option-id> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--force` | Skip confirmation prompt |

### Examples

```bash
# Delete option
ghx field delete-option PVTF_xxx OPT_xxx --force
```

## Available Colors

For single select options, these colors are available:

| Color | Preview |
|-------|---------|
| `gray` | Default |
| `red` | High priority |
| `orange` | Warning |
| `yellow` | Attention |
| `green` | Success |
| `blue` | Info |
| `purple` | Special |
| `pink` | Highlight |
