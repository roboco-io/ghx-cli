# Examples & Cookbook

Common workflows and use cases for ghx-cli.

## Project Management

### Create and Set Up a Sprint Project

```bash
# Create project
ghx project create "Sprint 1" --org myorg --description "Sprint 1: Jan 15-29"

# Add status field
ghx field create myorg/1 "Status" single_select \
  --options "Backlog,Todo,In Progress,Review,Done"

# Add priority field
ghx field create myorg/1 "Priority" single_select \
  --options "Critical,High,Medium,Low"

# Add story points field
ghx field create myorg/1 "Points" number

# Create board view
ghx view create myorg/1 "Sprint Board" board --group-by Status

# Create table view
ghx view create myorg/1 "All Tasks" table
```

### Add Issues to Project

```bash
# Add single issue
ghx item add myorg/1 myorg/repo#42

# Add multiple issues
ghx item add-bulk myorg/1 --items "repo#1,repo#2,repo#3,repo#4,repo#5"

# Add issues from search
ghx item add-bulk myorg/1 --query "is:issue label:sprint-1 state:open"

# Create draft issues
ghx item add myorg/1 --draft --title "Design review" --body "Review new designs"
ghx item add myorg/1 --draft --title "Deploy to staging"
ghx item add myorg/1 --draft --title "Update documentation"
```

### Manage Sprint Progress

```bash
# Move items to In Progress
ghx item edit myorg/1 PVTI_xxx --field Status --value "In Progress"

# Bulk update completed items
ghx analytics bulk-update myorg/1 \
  --filter "Status:Review" \
  --field Status \
  --value Done

# Archive completed items
ghx analytics bulk-archive myorg/1 --filter "Status:Done" --force
```

### Export Project for Reporting

```bash
# Export full project
ghx project export myorg/1 -o sprint1-report.json

# Export as CSV for spreadsheets
ghx analytics export myorg/1 -o sprint1.csv --format csv --items-only

# Get project overview
ghx analytics overview myorg/1 --format json > sprint1-stats.json
```

## Discussion Management

### Create Announcement

```bash
# Create announcement
ghx discussion create myorg/repo \
  --category announcements \
  --title "v2.0 Release" \
  --body "We're excited to announce v2.0! Key features:
- New dashboard
- Performance improvements
- Bug fixes

Please report any issues in this thread."
```

### Manage Q&A

```bash
# List unanswered questions
ghx discussion list myorg/repo --category q-a --unanswered

# View question with comments
ghx discussion view myorg/repo 42 --comments

# Add answer
ghx discussion comment myorg/repo 42 \
  --body "You can solve this by..."

# Mark as answered
ghx discussion answer myorg/repo 42 --comment-id DC_kwDOxxxxxx

# Close resolved question
ghx discussion close myorg/repo 42 --reason resolved
```

### Moderate Discussions

```bash
# Lock heated discussion
ghx discussion lock myorg/repo 123 --reason too_heated

# Close duplicate
ghx discussion close myorg/repo 456 --reason duplicate

# Delete spam
ghx discussion delete myorg/repo 789 --force
```

## Automation Scripts

### Daily Standup Report

```bash
#!/bin/bash
# daily-standup.sh

PROJECT="myorg/1"
DATE=$(date +%Y-%m-%d)

echo "=== Daily Standup Report - $DATE ==="

echo ""
echo "## In Progress"
ghx item list $PROJECT --format json | \
  jq -r '.[] | select(.status == "In Progress") | "- \(.title) (@\(.assignee))"'

echo ""
echo "## Completed Yesterday"
ghx item list $PROJECT --format json | \
  jq -r '.[] | select(.status == "Done" and .updatedAt >= "'$(date -v-1d +%Y-%m-%d)'") | "- \(.title)"'

echo ""
echo "## Blocked"
ghx item list $PROJECT --format json | \
  jq -r '.[] | select(.labels | contains(["blocked"])) | "- \(.title): \(.blockedReason)"'
```

### Weekly Project Backup

```bash
#!/bin/bash
# backup-projects.sh

BACKUP_DIR="$HOME/project-backups"
DATE=$(date +%Y%m%d)

mkdir -p "$BACKUP_DIR"

# List all projects
projects=$(ghx project list --org myorg --format json | jq -r '.[].number')

for project in $projects; do
  echo "Backing up project $project..."
  ghx project export myorg/$project \
    -o "$BACKUP_DIR/project-$project-$DATE.json" \
    --include-all
done

echo "Backup complete: $BACKUP_DIR"
```

### Sprint Transition

```bash
#!/bin/bash
# new-sprint.sh

OLD_PROJECT="myorg/1"
NEW_PROJECT="myorg/2"

# Create new sprint project
ghx project create "Sprint 2" --org myorg

# Copy unfinished items
incomplete=$(ghx item list $OLD_PROJECT --format json | \
  jq -r '.[] | select(.status != "Done") | .contentId')

for item in $incomplete; do
  ghx item add $NEW_PROJECT $item
done

# Archive old sprint
ghx analytics bulk-archive $OLD_PROJECT --filter "Status:Done" --force

echo "Sprint transition complete"
```

### Discussion Auto-Response

```bash
#!/bin/bash
# auto-respond.sh

REPO="myorg/repo"

# Get new unanswered questions
ghx discussion list $REPO --category q-a --unanswered --format json | \
  jq -r '.[] | select(.createdAt >= "'$(date -v-1d +%Y-%m-%dT%H:%M:%S)'") | .number' | \
while read num; do
  ghx discussion comment $REPO $num \
    --body "Thanks for your question! Our team will respond shortly.

In the meantime, please check our [FAQ](https://docs.example.com/faq) for common solutions."
done
```

## Integration Examples

### GitHub Actions Integration

```yaml
# .github/workflows/project-sync.yml
name: Sync Issues to Project

on:
  issues:
    types: [opened, labeled]

jobs:
  add-to-project:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install ghx-cli
        run: go install github.com/roboco-io/ghx-cli/cmd/ghx@latest

      - name: Add to project
        if: contains(github.event.issue.labels.*.name, 'priority-high')
        env:
          GITHUB_TOKEN: ${{ secrets.PROJECT_TOKEN }}
        run: |
          ghx item add myorg/1 ${{ github.repository }}#${{ github.event.issue.number }}
          ghx item edit myorg/1 $ITEM_ID --field Priority --value High
```

### Slack Notification Script

```bash
#!/bin/bash
# notify-slack.sh

PROJECT="myorg/1"
SLACK_WEBHOOK="https://hooks.slack.com/..."

# Get project stats
stats=$(ghx analytics overview $PROJECT --format json)

done=$(echo $stats | jq '.statusCounts.Done')
inProgress=$(echo $stats | jq '.statusCounts["In Progress"]')
total=$(echo $stats | jq '.totalItems')

# Send to Slack
curl -X POST $SLACK_WEBHOOK \
  -H 'Content-type: application/json' \
  -d "{
    \"text\": \"Project Update\",
    \"blocks\": [
      {
        \"type\": \"section\",
        \"text\": {
          \"type\": \"mrkdwn\",
          \"text\": \"*Sprint Progress*\n:white_check_mark: Done: $done\n:hourglass: In Progress: $inProgress\n:clipboard: Total: $total\"
        }
      }
    ]
  }"
```

## JSON Processing with jq

### Filter items by assignee

```bash
ghx item list myorg/repo --format json | \
  jq '.[] | select(.assignee == "octocat")'
```

### Get item counts by status

```bash
ghx item list myorg/repo --format json | \
  jq 'group_by(.status) | map({status: .[0].status, count: length})'
```

### Export to CSV

```bash
ghx item list myorg/repo --format json | \
  jq -r '["Title","Status","Assignee"], (.[] | [.title, .status, .assignee]) | @csv' > items.csv
```

### Find overdue items

```bash
TODAY=$(date +%Y-%m-%d)
ghx item list myorg/repo --format json | \
  jq --arg today "$TODAY" '.[] | select(.dueDate < $today and .status != "Done")'
```
