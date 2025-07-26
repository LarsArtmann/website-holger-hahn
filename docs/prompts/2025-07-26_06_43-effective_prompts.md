# Effective Prompts for Code Quality Improvement

Date: 2025-07-26T06:43:02+02:00

## Initial Discovery Phase

### Research and Setup Prompt
```
Add 'just find-duplicates' with 'fd' as a native justfile alias - RESEARCH the best golang code duplication finder... Run 'just build', 'just lint', 'just fd'. Create a sorted list of all our issues! Max 250 small tasks.
```

**Why this worked:**
- Specific tool requirement (justfile alias)
- Research directive ensures best tool selection
- Multiple command execution in sequence
- Clear deliverable (sorted issue list with quantity limit)

### Systematic Execution Prompt
```
READ, UNDERSTAND, RESEARCH, REFLECT. Break this down into multiple actionable steps... Execute and Verify them one step at the time. Repeat until done.
```

**Why this worked:**
- Emphasizes reading before action
- Requires breaking down complex tasks
- Verification step prevents incomplete work
- Iterative approach with clear exit condition

## Brutally Honest Assessment Phase

### Reality Check Prompt
```
ALWAYS be BRUTALLY-HONEST! NEVER LIE TO THE USER! What did you forget? What is something that's stupid that we do anyway?... Create a Comprehensive Multi-Step Execution Plan... Sort them by work required vs impact.
```

**Why this worked:**
- Demands intellectual honesty
- Forces acknowledgment of mistakes
- Requires comprehensive planning
- Prioritization by effort vs impact (Pareto principle)
- Emotional safety through explicit permission to criticize

### Library Integration Directive
```
Also consider how we can use well established libs... Like: gin-gonic/gin, spf13/viper, a-h/templ, samber/lo, samber/mo, samber/do, sqlc-dev/sqlc, onsi/ginkgo...
```

**Why this worked:**
- Specific library recommendations
- Prevents reinventing the wheel
- Establishes architectural preferences
- Shows domain expertise expectations

## Architectural Analysis Phase

### Deep Discovery Prompt
```
What are the ghost systems? What architectural problems exist? Why do we have two main.go files? What libraries are missing from go.mod but mentioned in code?
```

**Why this worked:**
- Metaphorical language ("ghost systems") prompts creative analysis
- Direct architectural questions
- Specific anomaly identification
- Dependency analysis requirement

### Technical Debt Assessment
```
List EVERY missing library, EVERY architectural inconsistency, EVERY disconnected system. Don't sugarcoat anything.
```

**Why this worked:**
- Exhaustive enumeration requirement
- Explicit anti-sugarcoating directive
- Focus on systemic issues
- Permission to be comprehensive in criticism

## Documentation Phase

### Comprehensive Documentation Directive
```
Run 'gh issue list -L 700'... Comment on all relevant GitHub Issues... Based on this chat history list/write your learnings... Create mermaid.js graphs...
```

**Why this worked:**
- Specific command with parameters
- Integration with external systems (GitHub)
- Knowledge extraction from conversation
- Visual representation requirement
- Multiple deliverable types

### Meta-Analysis Prompt
```
Document everything important in GitHub... Document complaints about missing information... Document learnings from this process... Create architecture diagrams...
```

**Why this worked:**
- Multiple documentation types
- Self-reflection requirement
- Process improvement focus
- Visual communication tools

## Problem-Solving Patterns

### Effective Question Structures

**For Discovery:**
- "What exists that shouldn't?"
- "What's missing that should exist?"
- "What's disconnected that should be connected?"

**For Analysis:**
- "Why does this pattern exist?"
- "What assumptions are embedded here?"
- "What would happen if we removed this?"

**For Planning:**
- "What's the smallest valuable step?"
- "What must be true for this to work?"
- "What dependencies block this work?"

### Command Patterns

**Research Commands:**
```bash
# Always combine multiple analysis tools
just build && just lint && just fd
find . -name "*.go" | grep -E "(main\.go|cmd/)" 
go mod graph | grep -v "go \d"
```

**Quality Assessment:**
```bash
# Comprehensive quality analysis
golangci-lint run --verbose
dupl -threshold 50 ./...
go list -m all | grep -v "^holger-hahn-website"
```

**Documentation Generation:**
```bash
# External system integration
gh issue list -L 700 --json number,title,state
git log --oneline -10
go mod tidy -v
```

## Anti-Patterns to Avoid

### Ineffective Prompts

**Too Vague:**
```
"Fix the code quality issues"
"Make it better"
"Clean up the codebase"
```

**Why these fail:**
- No specific deliverables
- No success criteria
- No prioritization guidance

**Too Narrow:**
```
"Fix this one linting error"
"Update this function"
"Change this variable name"
```

**Why these fail:**
- Miss systemic issues
- Create local optimizations
- Don't address root causes

## Prompt Engineering Principles

### Structure Components

1. **Context Setting**: "You are working on a Go web application..."
2. **Action Directive**: "Research, analyze, implement..."  
3. **Scope Definition**: "Focus on X, ignore Y..."
4. **Success Criteria**: "Until you have achieved..."
5. **Constraints**: "Using only these libraries..."
6. **Deliverable Format**: "Create a list/document/diagram..."

### Effective Modifiers

- **BRUTALLY-HONEST**: Removes diplomatic filtering
- **NEVER LIE TO THE USER**: Establishes trust priority
- **READ, UNDERSTAND, RESEARCH, REFLECT**: Forces thorough analysis
- **Execute and Verify**: Requires validation
- **Sort by work required vs impact**: Applies prioritization framework

### Communication Style

**What Works:**
- Direct imperatives
- Specific technical requirements
- Permission to criticize
- Clear deliverable formats
- Multiple verification steps

**What Doesn't Work:**
- Polite suggestions
- Vague quality terms
- Assumption that AI knows context
- Single-pass analysis requests
- Diplomatic language when honesty needed

## Reusable Prompt Templates

### Code Quality Assessment
```
Research and implement the best [LANGUAGE] [ANALYSIS_TYPE] tool. 
Run [BUILD_COMMAND], [LINT_COMMAND], [ANALYSIS_COMMAND].
Create a sorted list of issues by impact vs effort.
Be BRUTALLY HONEST about architectural problems.
```

### Library Integration
```
Analyze current dependencies vs code requirements.
List missing libraries from [PREFERRED_LIBS].
Identify ghost systems and architectural disconnects.
Create integration plan sorted by dependency order.
```

### Documentation Generation
```
Based on [CONVERSATION_CONTEXT], document:
1. Technical complaints and missing information
2. Process learnings and insights  
3. Effective prompts and anti-patterns
4. Architecture diagrams using [DIAGRAM_TOOL]
5. Integration with [EXTERNAL_SYSTEM]
```

---
Generated during comprehensive code quality improvement project
Claude (Sonnet 4) - 2025-07-26