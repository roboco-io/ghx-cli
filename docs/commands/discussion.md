# ghx discussion

Manage GitHub Discussions in repositories.

## Synopsis

```bash
ghx discussion <command> [arguments] [flags]
```

## Aliases

- `ghx discussion`
- `ghx disc`
- `ghx discussions`

## Description

The `discussion` command group provides comprehensive discussion management for GitHub Discussions, including creating, viewing, commenting, and managing discussion state.

## Commands

| Command | Description |
|---------|-------------|
| `list` | List discussions |
| `view` | View discussion details |
| `create` | Create a new discussion |
| `edit` | Edit a discussion |
| `delete` | Delete a discussion |
| `close` | Close a discussion |
| `reopen` | Reopen a discussion |
| `lock` | Lock a discussion |
| `unlock` | Unlock a discussion |
| `comment` | Add a comment |
| `answer` | Mark/unmark as answer |
| `category` | Manage categories |

## ghx discussion list

List discussions in a repository.

```bash
ghx discussion list <owner/repo> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--category` | Filter by category slug | - |
| `--state` | Filter by state (open, closed, all) | all |
| `--answered` | Show only answered | false |
| `--unanswered` | Show only unanswered | false |
| `-L, --limit` | Maximum number of discussions | 20 |
| `--format` | Output format (table, json) | table |

### Examples

```bash
# List all discussions
ghx discussion list myorg/repo

# List open discussions
ghx discussion list myorg/repo --state open

# Filter by category
ghx discussion list myorg/repo --category ideas

# Show unanswered Q&A
ghx discussion list myorg/repo --category q-a --unanswered

# JSON output
ghx discussion list myorg/repo --format json
```

## ghx discussion view

View a discussion.

```bash
ghx discussion view <owner/repo> <number> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--comments` | Include comments | false |
| `--format` | Output format (table, json) | table |

### Examples

```bash
# View discussion
ghx discussion view myorg/repo 123

# View with comments
ghx discussion view myorg/repo 123 --comments

# JSON output
ghx discussion view myorg/repo 123 --format json
```

## ghx discussion create

Create a new discussion.

```bash
ghx discussion create <owner/repo> [flags]
```

### Flags

| Flag | Description | Required |
|------|-------------|----------|
| `--category` | Category slug | Yes |
| `--title` | Discussion title | Yes |
| `--body` | Discussion body | No |
| `--body-file` | Read body from file | No |

### Examples

```bash
# Create discussion
ghx discussion create myorg/repo \
  --category general \
  --title "Feature Request" \
  --body "I would like to suggest..."

# Create from file
ghx discussion create myorg/repo \
  --category ideas \
  --title "New Feature" \
  --body-file proposal.md
```

## ghx discussion edit

Edit an existing discussion.

```bash
ghx discussion edit <owner/repo> <number> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--title` | New title |
| `--body` | New body |
| `--category` | New category |

### Examples

```bash
# Update title
ghx discussion edit myorg/repo 123 --title "Updated Title"

# Update body
ghx discussion edit myorg/repo 123 --body "Updated content"

# Change category
ghx discussion edit myorg/repo 123 --category announcements
```

## ghx discussion delete

Delete a discussion.

```bash
ghx discussion delete <owner/repo> <number> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--force` | Skip confirmation prompt |

### Examples

```bash
# Delete with confirmation
ghx discussion delete myorg/repo 123

# Delete without confirmation
ghx discussion delete myorg/repo 123 --force
```

## ghx discussion close

Close a discussion.

```bash
ghx discussion close <owner/repo> <number> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--reason` | Close reason (resolved, outdated, duplicate) |

### Examples

```bash
# Close discussion
ghx discussion close myorg/repo 123

# Close with reason
ghx discussion close myorg/repo 123 --reason resolved
```

## ghx discussion reopen

Reopen a closed discussion.

```bash
ghx discussion reopen <owner/repo> <number>
```

### Examples

```bash
ghx discussion reopen myorg/repo 123
```

## ghx discussion lock

Lock a discussion to prevent new comments.

```bash
ghx discussion lock <owner/repo> <number> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--reason` | Lock reason (off_topic, too_heated, resolved, spam) |

### Examples

```bash
# Lock discussion
ghx discussion lock myorg/repo 123

# Lock with reason
ghx discussion lock myorg/repo 123 --reason resolved
```

## ghx discussion unlock

Unlock a discussion.

```bash
ghx discussion unlock <owner/repo> <number>
```

### Examples

```bash
ghx discussion unlock myorg/repo 123
```

## ghx discussion comment

Add a comment to a discussion.

```bash
ghx discussion comment <owner/repo> <number> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--body` | Comment body |
| `--body-file` | Read body from file |
| `--reply-to` | Reply to specific comment ID |

### Examples

```bash
# Add comment
ghx discussion comment myorg/repo 123 --body "Thanks for the suggestion!"

# Reply to comment
ghx discussion comment myorg/repo 123 \
  --body "I agree with this" \
  --reply-to DC_kwDOxxxxxx
```

## ghx discussion answer

Mark or unmark a comment as the answer (Q&A categories only).

```bash
ghx discussion answer <owner/repo> <number> [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--comment-id` | Comment ID to mark as answer |
| `--unmark` | Unmark the current answer |

### Examples

```bash
# Mark comment as answer
ghx discussion answer myorg/repo 123 --comment-id DC_kwDOxxxxxx

# Unmark answer
ghx discussion answer myorg/repo 123 --unmark
```

## ghx discussion category

Manage discussion categories.

```bash
ghx discussion category <action> <owner/repo> [flags]
```

### Actions

- `list` - List categories

### Examples

```bash
# List categories
ghx discussion category list myorg/repo

# JSON output
ghx discussion category list myorg/repo --format json
```

## Discussion Categories

Common category types:

| Category | Type | Description |
|----------|------|-------------|
| Announcements | Announcement | News and updates |
| General | Discussion | General conversations |
| Ideas | Discussion | Feature suggestions |
| Q&A | Question | Questions (supports answers) |
| Show and tell | Discussion | Share projects |
