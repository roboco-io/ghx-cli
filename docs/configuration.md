# Configuration

ghx-cli can be configured through a config file, environment variables, or command-line flags.

## Configuration Precedence

Settings are applied in this order (later overrides earlier):

1. Config file (`~/.ghx.yaml`)
2. Environment variables
3. Command-line flags

## Config File

Create a config file at `~/.ghx.yaml`:

```yaml
# Authentication
token: "ghp_your_token_here"  # GitHub Personal Access Token

# Default context
org: "myorg"                  # Default organization
user: "myuser"                # Default user

# Output settings
format: "table"               # Default output format (table, json, yaml)
debug: false                  # Enable debug output
no-cache: false               # Disable caching

# Request settings
timeout: 30s                  # Request timeout
retries: 3                    # Number of retries for failed requests
```

### Custom Config Location

Use a different config file:

```bash
ghx --config /path/to/config.yaml project list
```

## Environment Variables

All settings can be configured via environment variables with the `GHX_` prefix:

| Variable | Description | Example |
|----------|-------------|---------|
| `GHX_TOKEN` | GitHub token | `ghp_xxxx` |
| `GITHUB_TOKEN` | GitHub token (fallback) | `ghp_xxxx` |
| `GH_TOKEN` | GitHub token (fallback) | `ghp_xxxx` |
| `GHX_ORG` | Default organization | `myorg` |
| `GHX_USER` | Default user | `myuser` |
| `GHX_FORMAT` | Output format | `json` |
| `GHX_DEBUG` | Enable debug | `true` |
| `GHX_NO_CACHE` | Disable cache | `true` |

### Example

```bash
export GHX_ORG="myorg"
export GHX_FORMAT="json"
export GHX_DEBUG="true"

ghx project list  # Uses org=myorg, format=json, debug=true
```

## Command-Line Flags

Override settings for a single command:

```bash
ghx project list --org myorg --format json --debug
```

### Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Config file path | `~/.ghx.yaml` |
| `--token` | GitHub token | - |
| `--org` | Organization | - |
| `--user` | User | - |
| `--format` | Output format | `table` |
| `--debug` | Debug output | `false` |
| `--no-cache` | Disable cache | `false` |

## Output Formats

### Table (Default)

Human-readable table format, best for interactive use:

```bash
ghx project list --format table
```

```
NUMBER  TITLE           ITEMS  UPDATED
1       Sprint 1        42     2024-01-15
2       Backlog         128    2024-01-14
```

### JSON

Machine-readable JSON, best for scripting:

```bash
ghx project list --format json
```

```json
[
  {"number": 1, "title": "Sprint 1", "items": 42, "updated": "2024-01-15"},
  {"number": 2, "title": "Backlog", "items": 128, "updated": "2024-01-14"}
]
```

### YAML

YAML format, useful for configuration files:

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

## Caching

ghx-cli caches API responses to improve performance. Cache settings:

```yaml
# In ~/.ghx.yaml
cache:
  enabled: true
  ttl: 5m              # Cache TTL
  directory: ~/.ghx/cache
```

### Disable Cache

```bash
# For single command
ghx project list --no-cache

# Via environment
export GHX_NO_CACHE=true
```

## Debug Mode

Enable debug mode to see detailed request/response information:

```bash
# Via flag
ghx project list --debug

# Via environment
export GHX_DEBUG=true
```

Debug output includes:
- API requests and responses
- Rate limit information
- Cache hits/misses
- Timing information

## Shell Completion

Generate shell completion scripts:

```bash
# Bash
ghx completion bash > /etc/bash_completion.d/ghx

# Zsh
ghx completion zsh > "${fpath[1]}/_ghx"

# Fish
ghx completion fish > ~/.config/fish/completions/ghx.fish

# PowerShell
ghx completion powershell > ghx.ps1
```

### Bash

Add to `~/.bashrc`:

```bash
source <(ghx completion bash)
```

### Zsh

Add to `~/.zshrc`:

```bash
source <(ghx completion zsh)
```

## Multiple Profiles

For managing multiple GitHub accounts, use different config files:

```bash
# Work profile
ghx --config ~/.ghx-work.yaml project list

# Personal profile
ghx --config ~/.ghx-personal.yaml project list
```

Or use shell aliases:

```bash
# In ~/.bashrc or ~/.zshrc
alias ghx-work='ghx --config ~/.ghx-work.yaml'
alias ghx-personal='ghx --config ~/.ghx-personal.yaml'
```

## Security Best Practices

1. **Never commit tokens**: Add `.ghx.yaml` to `.gitignore`
2. **Use environment variables**: Prefer `GITHUB_TOKEN` over config file
3. **Use GitHub CLI**: `gh auth login` is the most secure method
4. **Rotate tokens regularly**: Generate new tokens periodically
5. **Minimal scopes**: Only grant required permissions
