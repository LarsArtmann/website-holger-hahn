#!/bin/bash

# Quick Status Check for Holger Hahn Website
# Group G8 - Testing & Validation Baseline

echo "🔍 QUICK QUALITY STATUS CHECK"
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Test Status
echo -n "Tests: "
if just test > /dev/null 2>&1; then
    echo "✅ PASSING"
else
    echo "❌ FAILING"
fi

# Build Status  
echo -n "Build: "
if just build > /dev/null 2>&1; then
    echo "✅ PASSING"
else
    echo "❌ FAILING"
fi

# Lint Status
echo -n "Lint: "
if just lint > /dev/null 2>&1; then
    echo "✅ PASSING"
else
    echo "❌ FAILING"
fi

# Duplicate Count
DUPLICATE_COUNT=$(just fd 2>/dev/null | grep -c "found.*clones:" || echo "0")
echo "Duplicates: $DUPLICATE_COUNT groups"

# Modified Files
MODIFIED_COUNT=$(git status --porcelain | wc -l)
echo "Modified: $MODIFIED_COUNT files"

echo ""
if [ "$DUPLICATE_COUNT" -ge 10 ]; then
    echo "⚠️  High duplicate count detected"
fi

if [ "$MODIFIED_COUNT" -gt 20 ]; then
    echo "⚠️  Many modified files - consider committing progress"
fi

echo "📊 Run './scripts/monitor-quality.sh' for detailed analysis"