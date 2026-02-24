# GitHub Issue Importer

Import multiple GitHub issues from a CSV file.

## Prerequisites

```bash
# Install Python 3 if needed
# Install requests
pip install requests
```

## Setup

1. Create a GitHub Personal Access Token:
   - Go to https://github.com/settings/tokens
   - Generate new token (classic)
   - Select `repo` scope

2. Set environment variable:
   ```bash
   export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
   ```

## CSV Format

Create a CSV file with these columns:

```csv
title,body,labels
"Issue title here","Issue description goes here","label1,label2"
"Another issue","More details","bug,high-priority"
```

## Usage

```bash
python import_issues.py issues.csv owner/repo
```

Example for this project:
```bash
python import_issues.py issues.csv wahyudesu/quran-api-go
```

## Files

- `import_issues.py` - Main script
- `issues.csv` - Example CSV file
- `.gitignore` - Ignores token and Python cache
