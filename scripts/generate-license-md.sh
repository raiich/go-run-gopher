#!/bin/bash
# Generate license markdown for a single module
# Usage: ./generate-license-md.sh <csv-line> <license-dir>
# CSV format: module_name,license_url,license_type

set -euo pipefail

CSV_LINE="$1"
LICENSE_DIR="$2"

# Parse CSV with awk
read -r MODULE_NAME LICENSE_URL LICENSE_TYPE <<< $(echo "$CSV_LINE" | awk -F',' '{printf "%s %s %s", $1, $2, $3}')

# Get LICENSE file from saved directory (try LICENSE, then LICENSE.txt)
LICENSE_FILE="$LICENSE_DIR/$MODULE_NAME/LICENSE"

if [[ ! -f "$LICENSE_FILE" ]]; then
    LICENSE_FILE="$LICENSE_DIR/$MODULE_NAME/LICENSE.txt"
fi

if [[ ! -f "$LICENSE_FILE" ]]; then
    echo "Error: LICENSE file not found at $LICENSE_DIR/$MODULE_NAME/LICENSE{,.txt}" >&2
    exit 1
fi

# Output markdown
cat << EOF
## ${MODULE_NAME}

- ${LICENSE_TYPE} license
- ${LICENSE_URL}

\`\`\`
$(cat "$LICENSE_FILE")
\`\`\`

---

EOF
