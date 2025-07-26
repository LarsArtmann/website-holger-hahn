#!/bin/bash

# Quality Monitoring Script for Holger Hahn Website
# Group G8 - Testing & Validation Baseline
# Usage: ./scripts/monitor-quality.sh

set -e

TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
REPORT_FILE="docs/quality-monitoring-$(date '+%Y-%m-%d').log"

echo "=== QUALITY MONITORING REPORT ===" | tee -a "$REPORT_FILE"
echo "Timestamp: $TIMESTAMP" | tee -a "$REPORT_FILE"
echo "" | tee -a "$REPORT_FILE"

# Function to check and log status
check_status() {
    local test_name="$1"
    local command="$2"
    
    echo "Checking: $test_name" | tee -a "$REPORT_FILE"
    
    if eval "$command" > /tmp/quality_check.log 2>&1; then
        echo "✅ PASS: $test_name" | tee -a "$REPORT_FILE"
        return 0
    else
        echo "❌ FAIL: $test_name" | tee -a "$REPORT_FILE"
        echo "Error details:" | tee -a "$REPORT_FILE"
        cat /tmp/quality_check.log | head -10 | sed 's/^/  /' | tee -a "$REPORT_FILE"
        return 1
    fi
}

# Initialize counters
PASS_COUNT=0
FAIL_COUNT=0

# Test Suite
if check_status "Test Suite" "just test"; then
    ((PASS_COUNT++))
else
    ((FAIL_COUNT++))
fi

# Build System
if check_status "Build System" "just build"; then
    ((PASS_COUNT++))
else
    ((FAIL_COUNT++))
fi

# Linting
if check_status "Linting" "just lint"; then
    ((PASS_COUNT++))
else
    ((FAIL_COUNT++))
fi

# Duplicate Analysis
echo "Checking: Duplicate Code Analysis" | tee -a "$REPORT_FILE"
DUPLICATE_COUNT=$(just fd 2>/dev/null | grep -c "found.*clones:" || echo "0")
echo "Found $DUPLICATE_COUNT duplicate groups" | tee -a "$REPORT_FILE"

if [ "$DUPLICATE_COUNT" -lt 10 ]; then
    echo "✅ PASS: Duplicate Code Analysis (<10 groups)" | tee -a "$REPORT_FILE"
    ((PASS_COUNT++))
else
    echo "❌ FAIL: Duplicate Code Analysis (≥10 groups)" | tee -a "$REPORT_FILE"
    ((FAIL_COUNT++))
fi

# Git Status
MODIFIED_FILES=$(git status --porcelain | wc -l)
echo "Modified files: $MODIFIED_FILES" | tee -a "$REPORT_FILE"

# Summary
echo "" | tee -a "$REPORT_FILE"
echo "=== SUMMARY ===" | tee -a "$REPORT_FILE"
echo "Passed: $PASS_COUNT" | tee -a "$REPORT_FILE"
echo "Failed: $FAIL_COUNT" | tee -a "$REPORT_FILE"
echo "Duplicate Groups: $DUPLICATE_COUNT" | tee -a "$REPORT_FILE"
echo "Modified Files: $MODIFIED_FILES" | tee -a "$REPORT_FILE"

if [ "$FAIL_COUNT" -gt 0 ]; then
    echo "🚨 ALERT: Quality regression detected!" | tee -a "$REPORT_FILE"
    exit 1
else
    echo "✅ All quality checks passing" | tee -a "$REPORT_FILE"
    exit 0
fi