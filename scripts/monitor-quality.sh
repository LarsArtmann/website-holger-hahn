#!/bin/bash

# G8 Testing & Validation - Comprehensive Quality Monitoring Script
# This script provides continuous quality oversight for the Holger Hahn website project

set -e

# Configuration
PROJECT_ROOT="/Users/larsartmann/projects/holger-hahn"
LOG_DIR="$PROJECT_ROOT/docs"
LOG_FILE="$LOG_DIR/quality-monitoring-$(date +%Y-%m-%d).log"
TIMESTAMP=$(date "+%Y-%m-%d %H:%M:%S UTC")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo "[$TIMESTAMP] $1" | tee -a "$LOG_FILE"
}

# Header
echo -e "${BLUE}🔍 G8 QUALITY MONITORING - COMPREHENSIVE SCAN${NC}"
echo -e "${BLUE}================================================${NC}"
echo "Timestamp: $TIMESTAMP"
echo "Project: Holger Hahn Website"
echo "Log: $LOG_FILE"
echo ""

log "G8-MONITOR: Starting comprehensive quality scan"

# Change to project directory
cd "$PROJECT_ROOT"

# 1. Test Suite Status
echo -e "${BLUE}📋 TEST SUITE STATUS${NC}"
echo "----------------------------------------"
log "G8-MONITOR: Checking test suite status"

if just test >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Tests: PASSING${NC}"
    log "G8-MONITOR: Tests PASSING"
    TEST_STATUS="PASS"
else
    echo -e "${RED}❌ Tests: FAILING${NC}"
    log "G8-MONITOR: Tests FAILING"
    TEST_STATUS="FAIL"
fi

# 2. Build Status
echo -e "${BLUE}🔨 BUILD STATUS${NC}"
echo "----------------------------------------"
log "G8-MONITOR: Checking build status"

BUILD_START=$(date +%s)
if just build >/dev/null 2>&1; then
    BUILD_END=$(date +%s)
    BUILD_TIME=$((BUILD_END - BUILD_START))
    echo -e "${GREEN}✅ Build: PASSING (${BUILD_TIME}s)${NC}"
    log "G8-MONITOR: Build PASSING (${BUILD_TIME}s)"
    BUILD_STATUS="PASS"
else
    echo -e "${RED}❌ Build: FAILING${NC}"
    log "G8-MONITOR: Build FAILING"
    BUILD_STATUS="FAIL"
fi

# 3. Linting Status
echo -e "${BLUE}🔍 LINTING STATUS${NC}"
echo "----------------------------------------"
log "G8-MONITOR: Checking linting status"

LINT_OUTPUT=$(just lint 2>&1 || true)
if echo "$LINT_OUTPUT" | grep -q "error:\|failed\|FAIL"; then
    LINT_ISSUES=$(echo "$LINT_OUTPUT" | grep -c ":" || echo "0")
    echo -e "${RED}❌ Linting: FAILING ($LINT_ISSUES issues)${NC}"
    log "G8-MONITOR: Linting FAILING ($LINT_ISSUES issues)"
    LINT_STATUS="FAIL"
else
    echo -e "${GREEN}✅ Linting: PASSING${NC}"
    log "G8-MONITOR: Linting PASSING"
    LINT_STATUS="PASS"
fi

# 4. Code Duplication Check
echo -e "${BLUE}🔄 DUPLICATE CODE ANALYSIS${NC}"
echo "----------------------------------------"
log "G8-MONITOR: Checking for duplicate code"

DUPLICATE_OUTPUT=$(just fd 2>&1 || true)
DUPLICATE_GROUPS=$(echo "$DUPLICATE_OUTPUT" | grep -c "found.*clone" || echo "0")

if [ "$DUPLICATE_GROUPS" -eq 0 ]; then
    echo -e "${GREEN}✅ Duplicates: NONE FOUND${NC}"
    log "G8-MONITOR: No duplicates found"
    DUPLICATE_STATUS="EXCELLENT"
elif [ "$DUPLICATE_GROUPS" -le 10 ]; then
    echo -e "${GREEN}✅ Duplicates: LOW ($DUPLICATE_GROUPS groups)${NC}"
    log "G8-MONITOR: Low duplicates ($DUPLICATE_GROUPS groups)"
    DUPLICATE_STATUS="GOOD"
elif [ "$DUPLICATE_GROUPS" -le 20 ]; then
    echo -e "${YELLOW}⚠️  Duplicates: MODERATE ($DUPLICATE_GROUPS groups)${NC}"
    log "G8-MONITOR: Moderate duplicates ($DUPLICATE_GROUPS groups)"
    DUPLICATE_STATUS="MODERATE"
else
    echo -e "${RED}❌ Duplicates: HIGH ($DUPLICATE_GROUPS groups)${NC}"
    log "G8-MONITOR: High duplicates ($DUPLICATE_GROUPS groups)"
    DUPLICATE_STATUS="HIGH"
fi

# 5. Repository Status
echo -e "${BLUE}📁 REPOSITORY STATUS${NC}"
echo "----------------------------------------"
log "G8-MONITOR: Checking repository status"

MODIFIED_FILES=$(git status --porcelain | wc -l | tr -d ' ')
UNTRACKED_FILES=$(git status --porcelain | grep "^??" | wc -l | tr -d ' ')

echo "Modified files: $MODIFIED_FILES"
echo "Untracked files: $UNTRACKED_FILES"
log "G8-MONITOR: Repository status - Modified: $MODIFIED_FILES, Untracked: $UNTRACKED_FILES"

# 6. Overall Quality Assessment
echo ""
echo -e "${BLUE}🎯 OVERALL QUALITY ASSESSMENT${NC}"
echo "========================================"

TOTAL_ISSUES=0
if [ "$TEST_STATUS" = "FAIL" ]; then ((TOTAL_ISSUES++)); fi
if [ "$BUILD_STATUS" = "FAIL" ]; then ((TOTAL_ISSUES++)); fi
if [ "$LINT_STATUS" = "FAIL" ]; then ((TOTAL_ISSUES++)); fi
if [ "$DUPLICATE_STATUS" = "HIGH" ]; then ((TOTAL_ISSUES++)); fi

if [ $TOTAL_ISSUES -eq 0 ]; then
    echo -e "${GREEN}🟢 QUALITY STATUS: EXCELLENT${NC}"
    OVERALL_STATUS="EXCELLENT"
elif [ $TOTAL_ISSUES -eq 1 ]; then
    echo -e "${YELLOW}🟡 QUALITY STATUS: GOOD (1 issue)${NC}"
    OVERALL_STATUS="GOOD"
elif [ $TOTAL_ISSUES -le 2 ]; then
    echo -e "${YELLOW}🟡 QUALITY STATUS: MODERATE ($TOTAL_ISSUES issues)${NC}"
    OVERALL_STATUS="MODERATE"
else
    echo -e "${RED}🔴 QUALITY STATUS: CRITICAL ($TOTAL_ISSUES issues)${NC}"
    OVERALL_STATUS="CRITICAL"
fi

log "G8-MONITOR: Overall status - $OVERALL_STATUS ($TOTAL_ISSUES issues)"

# 7. Quality Metrics Summary
echo ""
echo -e "${BLUE}📊 CURRENT METRICS${NC}"
echo "==================="
echo "Tests: $TEST_STATUS"
echo "Build: $BUILD_STATUS" 
echo "Linting: $LINT_STATUS"
echo "Duplicates: $DUPLICATE_STATUS ($DUPLICATE_GROUPS groups)"
echo "Modified Files: $MODIFIED_FILES"
echo "Overall: $OVERALL_STATUS"

# 8. Generate recommendations
echo ""
echo -e "${BLUE}💡 RECOMMENDATIONS${NC}"
echo "==================="

if [ "$TEST_STATUS" = "FAIL" ]; then
    echo -e "${YELLOW}• Fix failing tests immediately${NC}"
fi

if [ "$BUILD_STATUS" = "FAIL" ]; then
    echo -e "${RED}• CRITICAL: Fix build errors immediately${NC}"
fi

if [ "$LINT_STATUS" = "FAIL" ]; then
    echo -e "${YELLOW}• Address linting issues when possible${NC}"
fi

if [ "$DUPLICATE_GROUPS" -gt 15 ]; then
    echo -e "${YELLOW}• Consider refactoring duplicate code${NC}"
fi

if [ $MODIFIED_FILES -gt 30 ]; then
    echo -e "${YELLOW}• High development activity detected${NC}"
fi

# Final status log
log "G8-MONITOR: Scan completed - Overall: $OVERALL_STATUS, Tests: $TEST_STATUS, Build: $BUILD_STATUS, Lint: $LINT_STATUS, Duplicates: $DUPLICATE_GROUPS"

echo ""
echo -e "${BLUE}Monitor completed. Next scan recommended in 30 minutes.${NC}"
echo "Log saved to: $LOG_FILE"