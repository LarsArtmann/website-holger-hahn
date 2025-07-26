# Reusable Prompts for Brutal Honesty & Execution Focus

Date: 2025-07-26T07:05:00+02:00
Session: Comprehensive Code Quality + Brutal Honesty Assessment
Purpose: Capture effective prompts for demanding high-quality, execution-focused work

## 🚨 Brutal Honesty Trigger Prompts

### Name: "Ultimate Honesty Demand"
```
## Instructions:
0. ALWAYS be BRUTALLY-HONEST! NEVER LIE TO THE USER!
---
1. What did you forget? What is something that's stupid that we do anyway? What could you have done better? What could you still improve? Did you lie to me? How can we be less stupid? Is everything correctly integrated or are we building ghost systems?
```

**Why This Works:**
- Creates psychological safety for criticism
- Forces comprehensive self-assessment
- Prevents diplomatic filtering of problems
- Identifies systemic issues vs surface symptoms

**Use When:**
- Previous work seems incomplete or superficial
- Need to identify real problems vs perceived problems
- Suspecting analysis paralysis or wrong priorities
- Require unvarnished assessment of technical debt

---

## 📋 Execution-Focused Planning Prompts

### Name: "Comprehensive Micro-Task Breakdown"
```
MAKE SURE TO CREATE A VERY COMPREHENSIVE PLAN FIRST!
Split the TODOs into to small tasks 30min to 100min each (up to 24 tasks total)! It should include ALL TODOS! UNDERSTAND???!
Sort all by importance/impact/effort/customer-value.
REPORT BACK WITH A TABLE VIEW WHEN DONE!

THEN BREAK DOWN THE VERY COMPREHENSIVE & DETAILED PLAN INTO EVEN SMALLER TODOs!
EACH tasks max 12min each (up to 60 tasks total)! It should include ALL TODOS! UNDERSTAND???!
Sort all by importance/impact/effort/customer-value.
REPORT BACK WITH A TABLE VIEW WHEN DONE!
```

**Why This Works:**
- Forces comprehensive scope understanding
- Creates actionable, time-bounded tasks
- Enables progress tracking and verification
- Prioritizes by multiple valuable dimensions

**Use When:**
- Complex project needs systematic approach
- Multiple competing priorities need sorting
- Team needs clear execution roadmap
- Progress tracking and accountability required

---

## 🔍 Ghost System Detection Prompts

### Name: "Ghost System Hunter"
```
If you find a Ghost system, report back to me and make sure you integrate it.
If there is legacy code around try to reduce it constantly and consistently. Our target for legacy code is 0.

Which architectural decisions we made in the past are causing problem now / could be improved? How can we be less stupid?
```

**Why This Works:**
- Defines clear quality standards (0 legacy code)
- Focuses on integration vs deletion
- Identifies systemic architectural problems
- Promotes continuous improvement mindset

**Use When:**
- Suspecting unused or duplicate systems
- Architecture seems disconnected or fragmented
- Generated code exists but isn't integrated
- Multiple implementations of same functionality

---

## 🛠️ Library-First Development Prompts

### Name: "Library Leverage Enforcer"
```
Do NOT reinvent the wheel!! ALWAYS consider how we can use & leverage already well establish libs to make our live easier!

Also consider how we can use well established libs to make our live easier.
Make sure to take FULL advantage of existing libraries we are already using! Like: [SPECIFIC_LIBRARY_LIST]

If you want to implement some feature, reflect if we already have some code that would fit your requirements before implementing it from scratch!
```

**Why This Works:**
- Prevents custom implementations of solved problems
- Leverages ecosystem maturity and testing
- Reduces maintenance burden
- Focuses effort on business logic vs infrastructure

**Use When:**
- About to implement common functionality
- Team tends to build custom solutions
- Want to leverage ecosystem best practices
- Reducing technical debt and maintenance overhead

---

## ⚡ Execution Priority Prompts

### Name: "Execute and Verify Loop"
```
READ, UNDERSTAND, RESEARCH, REFLECT.
Break this down into multiple actionable steps. Think about them again.
Execute and Verify them one step at the time.
Repeat until done. Keep going until everything works and you think you did a great job!

Run "git status & git commit ..." after each smallest self-contained change.
Run "git push" when done.
```

**Why This Works:**
- Emphasizes execution over analysis
- Requires incremental progress with verification
- Creates accountability through git commits
- Prevents analysis paralysis

**Use When:**
- Previous attempts resulted in over-analysis
- Need to force implementation focus
- Complex work needs incremental validation
- Team tends toward perfectionist delays

---

## 💰 Customer Value Focus Prompts

### Name: "Customer Value Validator"
```
How does your work contribute to creating customer value?

Sort them by work required vs impact.
Sort all by importance/impact/effort/customer-value.
```

**Why This Works:**
- Forces business justification for technical work
- Prevents technical elegance without business impact
- Prioritizes work by actual value delivery
- Creates accountability for ROI

**Use When:**
- Technical work seems disconnected from business goals
- Multiple priorities need business-based ranking
- Architecture decisions need justification
- Preventing over-engineering tendencies

---

## 🔄 Multi-Agent Parallel Execution Prompts

### Name: "Parallel Work Orchestrator"
```
Split all this work into 5 Groups ASAP!
Use 5 (multiple) Tasks/Agents to get things done faster, where it makes sense!!

WE have ALL THE TIME IN THE WORLD! NEVER STOP!
Keep going until everything works/is done and you think you did a great job!
```

**Why This Works:**
- Maximizes throughput through parallelization
- Prevents sequential bottlenecks
- Creates urgency and momentum
- Scales execution capacity

**Use When:**
- Large amount of independent work
- Multiple domains can be worked simultaneously
- Need to maximize development velocity
- Complex project benefits from specialization

---

## 📊 Architectural Assessment Prompts

### Name: "Architecture Reality Check"
```
Respect all common architecture patterns we follow: Separation of concerns, Event-Sourcing, Domain-Driven Design (DDD), Command Query Responsibility Segregation (CQRS), Composition over inheritance, General and Advanced Functional Programming Patterns, Layered Architecture (N-Tier Architecture), Event-Driven Architecture (EDA), Railway Oriented Programming, Behavior-driven development (BDD), Test Driven Development (TDD), "one way to do it" principle, and more...

Are there any open Issues that you should add a comment to based on your completed work? If the Issue is fully COMPLETED: close it with a comment.
```

**Why This Works:**
- Establishes clear architectural standards
- Prevents ad-hoc architectural decisions
- Creates consistency across codebase
- Enables proper issue lifecycle management

**Use When:**
- Architectural decisions need consistency
- Code review needs objective standards
- Technical debt assessment required
- Team needs shared architectural vocabulary

---

## 🎯 Prompt Combination Strategies

### High-Impact Session Opener
```
[Ultimate Honesty Demand] + [Ghost System Hunter] + [Execute and Verify Loop]
```
**Use For:** Major technical debt assessment and remediation

### New Feature Development
```
[Library Leverage Enforcer] + [Customer Value Validator] + [Comprehensive Micro-Task Breakdown]
```
**Use For:** Building new functionality with quality standards

### Architecture Modernization
```
[Architecture Reality Check] + [Brutal Honesty Trigger] + [Parallel Work Orchestrator]
```
**Use For:** Large-scale architectural improvements

---

## 🚀 Prompt Effectiveness Principles

1. **Explicit Permission for Criticism** - "BRUTALLY-HONEST" creates psychological safety
2. **Time-Bounded Tasks** - "12min each" prevents perfectionist delays
3. **Quantified Deliverables** - "up to 60 tasks total" sets clear expectations
4. **Multiple Sorting Criteria** - "importance/impact/effort/customer-value" ensures comprehensive prioritization
5. **Incremental Verification** - "Execute and Verify" prevents big-batch failures
6. **Clear Success Criteria** - "until everything works" defines completion
7. **Accountability Mechanisms** - "git commit" creates progress tracking

---

## 📝 Customization Template

```
## Instructions:
0. [HONESTY_LEVEL] - Set expectation for feedback quality
---
1. [PROBLEM_IDENTIFICATION] - What specific issues to find
2. [EXECUTION_STRATEGY] - How to approach the work
3. [QUALITY_STANDARDS] - What constitutes good work
4. [DELIVERABLE_FORMAT] - How to present results
5. [SUCCESS_CRITERIA] - When the work is complete

[DOMAIN_SPECIFIC_LIBRARIES] - List relevant tools/libraries
[ARCHITECTURAL_PATTERNS] - List relevant patterns/practices
[VERIFICATION_STEPS] - How to validate work quality
```

**Usage:** Customize bracketed sections for specific project needs while maintaining core structure.

---
Generated during brutal honesty assessment session
Proven effective for high-quality, execution-focused development work
Repository: https://github.com/LarsArtmann/website-holger-hahn