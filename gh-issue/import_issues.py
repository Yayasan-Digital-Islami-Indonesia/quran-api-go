#!/usr/bin/env python3
"""
GitHub Issue Importer from CSV

Usage:
1. Create a CSV file with columns: title, body, labels
2. Set GITHUB_TOKEN environment variable
3. Run: python import_issues.py <csv-file> <owner/repo>

Example:
    export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
    python import_issues.py issues.csv wahyudesu/quran-api-go
"""

import csv
import os
import sys
import requests
from typing import List, Dict


def read_csv(file_path: str) -> List[Dict[str, str]]:
    """Read issues from CSV file."""
    issues = []
    with open(file_path, "r", encoding="utf-8") as f:
        reader = csv.DictReader(f)
        for row in reader:
            issues.append({
                "title": row.get("title", "").strip(),
                "body": row.get("body", "").strip(),
                "labels": [l.strip() for l in row.get("labels", "").split(",") if l.strip()]
            })
    return issues


def create_issue(repo: str, issue: Dict[str, str], token: str) -> bool:
    """Create a single GitHub issue."""
    url = f"https://api.github.com/repos/{repo}/issues"
    headers = {
        "Authorization": f"Bearer {token}",
        "Accept": "application/vnd.github.v3+json",
        "X-GitHub-Api-Version": "2022-11-28"
    }
    payload = {
        "title": issue["title"],
        "body": issue["body"],
        "labels": issue["labels"]
    }

    response = requests.post(url, json=payload, headers=headers)

    if response.status_code in [201, 200]:
        print(f"✓ Created: {issue['title']}")
        return True
    else:
        print(f"✗ Failed: {issue['title']}")
        print(f"  Status: {response.status_code}")
        print(f"  Error: {response.text}")
        return False


def main():
    if len(sys.argv) < 3:
        print("Usage: python import_issues.py <csv-file> <owner/repo>")
        print("Example: python import_issues.py issues.csv wahyudesu/quran-api-go")
        sys.exit(1)

    csv_file = sys.argv[1]
    repo = sys.argv[2]
    token = os.getenv("GITHUB_TOKEN")

    if not token:
        print("Error: GITHUB_TOKEN environment variable not set")
        print("Get your token at: https://github.com/settings/tokens")
        sys.exit(1)

    if not os.path.exists(csv_file):
        print(f"Error: File '{csv_file}' not found")
        sys.exit(1)

    issues = read_csv(csv_file)

    if not issues:
        print(f"Error: No issues found in '{csv_file}'")
        sys.exit(1)

    print(f"Found {len(issues)} issues to import to {repo}")
    print("-" * 50)

    success = 0
    failed = 0

    for issue in issues:
        if create_issue(repo, issue, token):
            success += 1
        else:
            failed += 1

    print("-" * 50)
    print(f"Done! {success} created, {failed} failed")


if __name__ == "__main__":
    main()
