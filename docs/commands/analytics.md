# ghx analytics

Generate analytics, reports, and perform bulk operations on GitHub Projects.

## Synopsis

```bash
ghx analytics <command> [arguments] [flags]
```

## Description

The `analytics` command group provides reporting, data export/import, and bulk operation capabilities for GitHub Projects.

## Commands

| Command | Description |
|---------|-------------|
| `overview` | Generate project overview |
| `velocity` | Generate velocity metrics |
| `timeline` | Generate timeline analytics |
| `distribution` | Generate item distribution |
| `export` | Export project data |
| `import` | Import project data |
| `bulk-update` | Update multiple items |
| `bulk-delete` | Delete multiple items |
| `bulk-archive` | Archive multiple items |
| `operation-status` | Check bulk operation status |

## ghx analytics overview

Generate project overview statistics.

```bash
ghx analytics overview <project-ref> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format (table, json) | table |

### Examples

```bash
# Get project overview
ghx analytics overview myorg/123

# JSON output
ghx analytics overview myorg/123 --format json
```

### Output

- Total items count
- Items by status
- Items by assignee
- Recent activity
- Completion percentage

## ghx analytics velocity

Generate velocity metrics over time.

```bash
ghx analytics velocity <project-ref> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--period` | Time period (daily, weekly, monthly) | weekly |
| `--range` | Date range (e.g., 30d, 3m) | 30d |
| `--format` | Output format (table, json) | table |

### Examples

```bash
# Weekly velocity
ghx analytics velocity myorg/123 --period weekly

# Monthly velocity for 3 months
ghx analytics velocity myorg/123 --period monthly --range 3m
```

## ghx analytics timeline

Generate timeline analytics.

```bash
ghx analytics timeline <project-ref> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--range` | Date range | 30d |
| `--format` | Output format (table, json) | table |

### Examples

```bash
ghx analytics timeline myorg/123 --range 90d
```

## ghx analytics distribution

Generate item distribution analytics.

```bash
ghx analytics distribution <project-ref> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--by` | Group by (status, assignee, label, priority) | status |
| `--format` | Output format (table, json) | table |

### Examples

```bash
# Distribution by status
ghx analytics distribution myorg/123 --by status

# Distribution by assignee
ghx analytics distribution myorg/123 --by assignee
```

## ghx analytics export

Export project data.

```bash
ghx analytics export <project-ref> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output file | stdout |
| `--format` | Export format (json, csv, xml) | json |
| `--include-all` | Include all data | false |
| `--items-only` | Export items only | false |

### Examples

```bash
# Export to JSON
ghx analytics export myorg/123 -o project.json

# Export as CSV
ghx analytics export myorg/123 -o items.csv --format csv --items-only

# Export all data
ghx analytics export myorg/123 -o full-backup.json --include-all
```

## ghx analytics import

Import project data.

```bash
ghx analytics import <project-ref> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-i, --input` | Input file | - |
| `--strategy` | Import strategy | merge |
| `--dry-run` | Preview without changes | false |

### Import Strategies

| Strategy | Description |
|----------|-------------|
| `merge` | Merge with existing data |
| `replace` | Replace existing items |
| `append` | Add without modifying existing |
| `skip_conflicts` | Skip conflicting items |

### Examples

```bash
# Import with merge
ghx analytics import myorg/123 -i data.json --strategy merge

# Preview import
ghx analytics import myorg/123 -i data.json --dry-run

# Replace existing
ghx analytics import myorg/123 -i data.json --strategy replace
```

## ghx analytics bulk-update

Update multiple project items at once.

```bash
ghx analytics bulk-update <project-ref> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--items` | Comma-separated item IDs |
| `--filter` | Filter expression |
| `--field` | Field to update |
| `--value` | New value |

### Examples

```bash
# Update specific items
ghx analytics bulk-update myorg/123 \
  --items "PVTI_a,PVTI_b,PVTI_c" \
  --field Status \
  --value Done

# Update items matching filter
ghx analytics bulk-update myorg/123 \
  --filter "Status:In Progress" \
  --field Status \
  --value Done
```

## ghx analytics bulk-delete

Delete multiple project items.

```bash
ghx analytics bulk-delete <project-ref> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--items` | Comma-separated item IDs |
| `--filter` | Filter expression |
| `--force` | Skip confirmation |

### Examples

```bash
# Delete specific items
ghx analytics bulk-delete myorg/123 --items "PVTI_a,PVTI_b" --force

# Delete items matching filter
ghx analytics bulk-delete myorg/123 --filter "label:wontfix" --force
```

## ghx analytics bulk-archive

Archive multiple project items.

```bash
ghx analytics bulk-archive <project-ref> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--items` | Comma-separated item IDs |
| `--filter` | Filter expression |
| `--force` | Skip confirmation |

### Examples

```bash
# Archive completed items
ghx analytics bulk-archive myorg/123 --filter "Status:Done" --force
```

## ghx analytics operation-status

Check the status of a bulk operation.

```bash
ghx analytics operation-status <operation-id>
```

### Examples

```bash
ghx analytics operation-status OP_kwDOxxxxxx
```

### Output

- Operation type
- Status (pending, in_progress, completed, failed)
- Items processed
- Errors if any
