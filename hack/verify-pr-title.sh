# Copyright Contributors to the Open Cluster Management project
#!/bin/bash

# Original Copyright:
# Copyright 2024 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Define regex patterns
WIP_REGEX="^\W?WIP\W"
TAG_REGEX="^\[[[:alnum:]\._-]*\]"
PR_TITLE="$1"

# Trim WIP and tags from title
trimmed_title=$(echo "$PR_TITLE" | sed -E "s/$WIP_REGEX//" | sed -E "s/$TAG_REGEX//" | xargs -0)

# Normalize common emojis in text form to actual emojis
trimmed_title=$(echo "$trimmed_title" | sed -E "s/:sparkles:/✨/g")
trimmed_title=$(echo "$trimmed_title" | sed -E "s/:bug:/🐛/g")
trimmed_title=$(echo "$trimmed_title" | sed -E "s/:book:/📖/g")
trimmed_title=$(echo "$trimmed_title" | sed -E "s/:memo:/📝/g")
trimmed_title=$(echo "$trimmed_title" | sed -E "s/:warning:/⚠️/g")
trimmed_title=$(echo "$trimmed_title" | sed -E "s/:seedling:/🌱/g")
trimmed_title=$(echo "$trimmed_title" | sed -E "s/:question:/❓/g")

# Check PR type prefix
if [[ "$trimmed_title" =~ ^(⚠|✨|🐛|📖|🚀|🌱) ]]; then
    echo "PR title is valid: $trimmed_title"
else
    echo "Error: No matching PR type indicator found in title."
    echo "Please copy the appropriate `:text:` or icon to the beginning of your PR title:"
    echo ":sparkles: ✨ feature"
    echo ":bug: 🐛 bug fix"
    echo ":book: 📖 docs"
    echo ":memo: 📝 proposal"
    echo ":warning: ⚠️ breaking change"
    echo ":seedling: 🌱 other/misc"
    echo ":question: ❓ requires manual review/categorization"
    exit 1
fi

# Check that PR title does not contain Issue or PR number
if [[ "$trimmed_title" =~ \#[0-9]+ ]]; then
    echo "Error: PR title should not contain issue or PR number."
    echo "Issue numbers belong in the PR body as either \"Fixes #XYZ\" (if it closes the issue or PR), or something like \"Related to #XYZ\" (if it's just related)."
    exit 1
fi
