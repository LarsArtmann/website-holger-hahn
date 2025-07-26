#!/bin/bash

# G8 Testing & Validation - Quick Status Check
# Provides rapid quality status overview without detailed analysis

set -e

# Configuration  
PROJECT_ROOT="/Users/larsartmann/projects/holger-hahn"
TIMESTAMP=$(date "+%H:%M:%S")

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}⚡ G8 QUICK STATUS CHECK${NC} [$TIMESTAMP]"
echo "=================================="

cd "$PROJECT_ROOT"

# Quick test check
if just test >/dev/null 2>&1; then
    echo -e "Tests: ${GREEN}✅ PASS${NC}"
else
    echo -e "Tests: ${RED}❌ FAIL${NC}"
fi

# Quick build check
if just build >/dev/null 2>&1; then
    echo -e "Build: ${GREEN}✅ PASS${NC}"
else
    echo -e "Build: ${RED}❌ FAIL${NC}"
fi

# Quick lint check
if just lint >/dev/null 2>&1; then
    echo -e "Lint:  ${GREEN}✅ PASS${NC}"
else
    echo -e "Lint:  ${YELLOW}⚠️  ISSUES${NC}"
fi

# Quick duplicate check
DUPLICATES=$(just fd 2>&1 | grep -c "found.*clone" || echo "0")
if [ "$DUPLICATES" -eq 0 ]; then
    echo -e "Dupes: ${GREEN}✅ NONE${NC}"
elif [ "$DUPLICATES" -le 10 ]; then
    echo -e "Dupes: ${GREEN}✅ LOW ($DUPLICATES)${NC}"
elif [ "$DUPLICATES" -le 20 ]; then
    echo -e "Dupes: ${YELLOW}⚠️  MOD ($DUPLICATES)${NC}"
else
    echo -e "Dupes: ${RED}❌ HIGH ($DUPLICATES)${NC}"
fi

# Modified files
MODIFIED=$(git status --porcelain | wc -l | tr -d ' ')
if [ "$MODIFIED" -eq 0 ]; then
    echo -e "Files: ${GREEN}✅ CLEAN${NC}"
elif [ "$MODIFIED" -le 10 ]; then
    echo -e "Files: ${GREEN}✅ ACTIVE ($MODIFIED)${NC}"
else
    echo -e "Files: ${YELLOW}⚠️  BUSY ($MODIFIED)${NC}"
fi

echo ""
echo -e "${BLUE}Use 'just monitor' for detailed analysis${NC}"