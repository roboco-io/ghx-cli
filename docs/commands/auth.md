# ghx auth

Manage GitHub authentication.

## Synopsis

```bash
ghx auth <command> [flags]
```

## Description

The `auth` command group provides authentication management for ghx-cli. It integrates with GitHub CLI for seamless authentication with fallback to environment variables.

## Commands

| Command | Description |
|---------|-------------|
| `status` | Show authentication status |

## ghx auth status

Display current authentication status.

```bash
ghx auth status [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--format` | Output format (table, json) | table |

### Examples

```bash
# Check status
ghx auth status

# JSON output
ghx auth status --format json
```

### Output

The status command displays:

- **GitHub CLI Status**: Whether `gh` is installed
- **Environment Token**: Whether `GITHUB_TOKEN` or `GHX_TOKEN` is set
- **Token Availability**: Whether a valid token is available
- **Token Validity**: Whether the token is valid and not expired
- **Available Scopes**: Token permission scopes
- **Required Scopes**: Scopes needed for ghx-cli features

### Example Output

```
GitHub CLI Authentication Status
================================

Status: Ready

Details:
--------
GitHub CLI: Installed
Environment Token: Not set
Token: Available
Token Validity: Valid
Required Scopes: Available

Available Scopes: [admin:org delete_repo gist project repo workflow]
Required Scopes: [repo project]

Recommendation:
---------------
Authentication is properly configured
```

## Authentication Methods

### Method 1: GitHub CLI (Recommended)

ghx-cli automatically uses your GitHub CLI authentication:

```bash
# Login with GitHub CLI
gh auth login

# Verify
ghx auth status
```

### Method 2: Environment Variables

Set a GitHub Personal Access Token:

```bash
# Using GITHUB_TOKEN
export GITHUB_TOKEN="ghp_your_token_here"

# Or using GHX_TOKEN
export GHX_TOKEN="ghp_your_token_here"
```

### Method 3: Config File

Add token to config file (`~/.ghx.yaml`):

```yaml
token: "ghp_your_token_here"
```

### Method 4: Command Line Flag

Pass token directly (not recommended for security):

```bash
ghx project list --token "ghp_your_token_here"
```

## Token Scopes

### Required Scopes

| Scope | Description | Commands |
|-------|-------------|----------|
| `repo` | Full repository access | All project/item commands |
| `project` | Project access | All project commands |

### Optional Scopes

| Scope | Description | Commands |
|-------|-------------|----------|
| `read:discussion` | Read discussions | `discussion list`, `discussion view` |
| `write:discussion` | Write discussions | `discussion create`, `discussion edit` |
| `admin:org` | Organization admin | Org-level project management |

## Creating a Token

1. Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. Click "Generate new token (classic)"
3. Select required scopes:
   - `repo` - Full control of private repositories
   - `project` - Full control of projects
   - `write:discussion` - Read and write discussions
4. Copy the token and store it securely

## Troubleshooting

### Token Not Found

```
Error: no authentication token available
```

**Solution**: Login with GitHub CLI or set environment variable:
```bash
gh auth login
# or
export GITHUB_TOKEN="your_token"
```

### Invalid Token

```
Error: token is invalid or expired
```

**Solution**: Generate a new token or re-authenticate:
```bash
gh auth refresh
```

### Missing Scopes

```
Error: token lacks required scope: project
```

**Solution**: Generate a new token with required scopes or:
```bash
gh auth refresh -s project
```
