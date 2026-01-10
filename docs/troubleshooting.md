# Troubleshooting

Common issues and solutions for ghx-cli.

## Authentication Issues

### No authentication token available

```
Error: no authentication token available
```

**Cause**: ghx-cli cannot find a valid GitHub token.

**Solutions**:

1. Login with GitHub CLI:
   ```bash
   gh auth login
   ```

2. Set environment variable:
   ```bash
   export GITHUB_TOKEN="ghp_your_token"
   ```

3. Add to config file (`~/.ghx.yaml`):
   ```yaml
   token: "ghp_your_token"
   ```

### Token is invalid or expired

```
Error: Bad credentials
```

**Cause**: The token is invalid, expired, or revoked.

**Solutions**:

1. Refresh GitHub CLI authentication:
   ```bash
   gh auth refresh
   ```

2. Generate a new token at [GitHub Settings](https://github.com/settings/tokens)

3. Check token validity:
   ```bash
   ghx auth status
   ```

### Missing required scopes

```
Error: Resource not accessible by integration
```

**Cause**: Token lacks required permissions.

**Required scopes**:
- `repo` - Repository access
- `project` - Projects access
- `write:discussion` - Discussion write access

**Solutions**:

1. Add scopes via GitHub CLI:
   ```bash
   gh auth refresh -s repo,project
   ```

2. Generate new token with required scopes

## API Errors

### Rate limit exceeded

```
Error: API rate limit exceeded
```

**Cause**: Too many API requests in a short time.

**Solutions**:

1. Wait for rate limit reset (usually 1 hour)

2. Check current rate limit:
   ```bash
   gh api rate_limit
   ```

3. Use caching to reduce requests:
   ```bash
   # Ensure caching is enabled
   ghx project list  # Cached after first request
   ```

### Resource not found

```
Error: Could not resolve to a Project
```

**Cause**: Project or resource doesn't exist or is not accessible.

**Solutions**:

1. Verify the project reference:
   ```bash
   # Correct format: owner/number
   ghx project view myorg/123
   ```

2. Check access permissions

3. Verify the project exists:
   ```bash
   gh project list --owner myorg
   ```

### GraphQL errors

```
Error: GraphQL error: Field 'xxx' doesn't exist
```

**Cause**: API schema mismatch, usually due to GitHub API changes.

**Solutions**:

1. Update ghx-cli to latest version:
   ```bash
   go install github.com/roboco-io/ghx-cli/cmd/ghx@latest
   ```

2. Report the issue on [GitHub Issues](https://github.com/roboco-io/ghx-cli/issues)

## Command Errors

### Invalid repository format

```
Error: invalid repository format: expected owner/repo
```

**Cause**: Repository reference is not in correct format.

**Solution**: Use `owner/repo` format:
```bash
# Correct
ghx discussion list myorg/my-repo

# Incorrect
ghx discussion list my-repo
ghx discussion list myorg/my-repo/issues
```

### Invalid project reference

```
Error: invalid project reference
```

**Cause**: Project reference is not in correct format.

**Solution**: Use `owner/number` or project ID:
```bash
# By owner/number
ghx project view myorg/123

# By project ID
ghx project view PVT_kwDOBcH12s4AXxyz
```

### Required flag missing

```
Error: required flag "category" not set
```

**Cause**: A required flag was not provided.

**Solution**: Add the required flag:
```bash
ghx discussion create myorg/repo --category general --title "Title"
```

## Output Issues

### No output displayed

**Cause**: Command succeeded but returned no results.

**Solutions**:

1. Check if resources exist:
   ```bash
   ghx project list --format json
   ```

2. Try different filters:
   ```bash
   ghx discussion list myorg/repo --state all
   ```

3. Enable debug mode:
   ```bash
   ghx project list --debug
   ```

### JSON parsing error

```
Error: invalid character in JSON
```

**Cause**: Output contains non-JSON content.

**Solutions**:

1. Check for debug output mixed with JSON:
   ```bash
   # Disable debug when piping
   ghx project list --format json 2>/dev/null | jq .
   ```

2. Redirect stderr:
   ```bash
   ghx project list --format json 2>&1 | jq .
   ```

## Performance Issues

### Slow commands

**Cause**: Large data sets or disabled caching.

**Solutions**:

1. Enable caching (default):
   ```bash
   # Ensure --no-cache is not set
   ghx project list
   ```

2. Use limits:
   ```bash
   ghx item list myorg/repo --limit 20
   ```

3. Use filters:
   ```bash
   ghx discussion list myorg/repo --state open --limit 10
   ```

### High memory usage

**Cause**: Processing large projects with many items.

**Solutions**:

1. Use pagination:
   ```bash
   ghx item list myorg/repo --limit 100
   ```

2. Export in chunks:
   ```bash
   ghx analytics export myorg/123 --items-only
   ```

## Build Issues

### Build fails

```
Error: cannot find package
```

**Solutions**:

1. Update Go modules:
   ```bash
   go mod tidy
   ```

2. Clear module cache:
   ```bash
   go clean -modcache
   go mod download
   ```

3. Ensure Go version 1.21+:
   ```bash
   go version
   ```

### Test failures

```
FAIL: TestXxx
```

**Solutions**:

1. Run with verbose output:
   ```bash
   go test -v ./...
   ```

2. Run specific test:
   ```bash
   go test -v -run TestXxx ./internal/service/...
   ```

3. Skip integration tests:
   ```bash
   go test -short ./...
   ```

## Getting Help

If your issue isn't covered here:

1. **Check debug output**:
   ```bash
   ghx <command> --debug
   ```

2. **Search existing issues**:
   [GitHub Issues](https://github.com/roboco-io/ghx-cli/issues)

3. **Ask in discussions**:
   [GitHub Discussions](https://github.com/roboco-io/ghx-cli/discussions)

4. **Report a bug**:
   [New Issue](https://github.com/roboco-io/ghx-cli/issues/new)

When reporting issues, include:
- ghx-cli version (`ghx --version`)
- Go version (`go version`)
- OS and architecture
- Full error message
- Steps to reproduce
