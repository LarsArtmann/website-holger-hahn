# Professional Website Transformation Prompt

**Name:** `PROFESSIONAL_WEBSITE_TRANSFORMATION`

**Purpose:** Transform a boring, colorful CV-style website into a compelling, professional client magnet that makes enterprise prospects want to hire the consultant immediately.

---

## Prompt Template

```
# PROFESSIONAL WEBSITE TRANSFORMATION MISSION

You are tasked with transforming a consultant's website from a boring CV dump into a compelling client acquisition machine.

## CURRENT PROBLEM
- Website looks like a colorful, unprofessional CV
- Boring bullet points that could apply to anyone
- No visual credibility or company branding
- Missing the "so what?" factor for enterprise clients
- Lacks compelling business impact stories

## TARGET TRANSFORMATION
- Compelling narratives showing real business impact
- Visual company logos building credibility
- Specific metrics and achievements that matter
- Story-driven content that builds trust
- Professional, enterprise-grade appearance

## CORE REQUIREMENTS

### 1. CONTENT TRANSFORMATION
- Transform CV bullets into "Challenge → Solution → Result" impact stories
- Add specific metrics and quantified business results
- Include "What This Means For You" client benefit sections
- Highlight regulatory compliance and risk mitigation expertise
- Show progression of technical expertise over time

### 2. VISUAL CREDIBILITY
- Research and integrate company logos from LinkedIn history
- Create professional logo display system with hover effects
- Ensure mobile-responsive presentation
- Link to company websites for credibility verification
- Maintain consistent brand colors (Julius Bär Blue #141d55)

### 3. PROFESSIONAL DESIGN SYSTEM
- Eliminate "rainbow disaster" - use monochrome with single brand accent
- Implement flat design (zero border-radius, zero box-shadows)
- Force all devicon colors to grayscale with CSS overrides
- Create consistent typography hierarchy
- Ensure mobile-first responsive design

### 4. TECHNICAL EXCELLENCE
- Use proper Go project structure with recommended libraries
- Implement gin (HTTP server), viper (config), templ (templates), sqlc (database)
- Add ginkgo/gomega for BDD testing
- Integrate uniflow for error handling
- Follow DDD, CQRS, and Event Sourcing patterns where appropriate

## EXECUTION STRATEGY

### Phase 1: Customer Value First (70% effort)
1. Transform experience section content into compelling stories
2. Add company logos and visual credibility markers
3. Create clear client benefit propositions
4. Test narrative flow and engagement

### Phase 2: Technical Foundation (30% effort)
1. Implement proper Go architecture
2. Add testing and error handling
3. Optimize for mobile and performance
4. Complete integration and QA

## SUCCESS CRITERIA

### Business Impact Metrics
- Compelling experience narratives that showcase real business value
- Visual credibility through recognized company branding
- Clear value propositions that answer "What's in it for me?"
- Professional appearance that builds immediate trust

### Technical Quality Metrics
- Proper Go project structure with recommended libraries
- Mobile-optimized responsive design
- Fast loading performance (<3s)
- Cross-browser compatibility
- Comprehensive test coverage

## COMMON PITFALLS TO AVOID

1. **Over-engineering:** Don't create complex architecture for simple websites
2. **API limits:** Don't spawn too many Task agents simultaneously
3. **Technical perfectionism:** Prioritize customer value over code elegance
4. **Ignoring mobile:** Enterprise clients use mobile devices too
5. **Boring content:** CV bullets don't sell - impact stories do

## QUALITY ASSURANCE

### Content Review Questions
- Does this make enterprise clients want to hire this person immediately?
- Are the business impact stories compelling and specific?
- Is the value proposition crystal clear?
- Does this build trust and credibility?

### Technical Review Questions
- Is the architecture appropriate for the project complexity?
- Does it work perfectly on mobile devices?
- Is the performance optimized for fast loading?
- Are all dependencies properly integrated?

## DELIVERABLES

1. **Transformed Experience Section** - Compelling business impact stories
2. **Visual Credibility System** - Professional company logo display
3. **Professional Design System** - Consistent, trustworthy appearance
4. **Technical Foundation** - Proper Go architecture with recommended libraries
5. **Documentation** - Architecture decisions and learnings recorded

---

**Expected Outcome:** Website transformation from "boring CV" to "client magnet" that makes enterprise prospects say "I need to hire this person immediately!"

**Time Estimate:** 12-16 hours for complete transformation with proper planning and execution.

**Success Metric:** Enterprise client conversion and professional credibility establishment.
```

---

## Usage Instructions

1. **Assess Current State:** Identify boring/unprofessional elements
2. **Gather LinkedIn Data:** Collect detailed professional history
3. **Define Success Criteria:** What makes clients want to hire this person?
4. **Execute in Phases:** Customer value first, technical excellence second  
5. **Test with Target Audience:** Validate with actual enterprise prospects
6. **Document Learnings:** Record what worked and what didn't

## Customization Points

- **Brand Colors:** Adjust Julius Bär Blue to client's brand colors
- **Industry Context:** Customize for specific professional service (legal, consulting, etc.)
- **Target Clients:** Adapt messaging for specific client personas
- **Technical Stack:** Modify library recommendations based on preferences
- **Compliance Requirements:** Adjust for industry-specific regulations

## Related Prompts

- `ENTERPRISE_CREDIBILITY_BUILDER` - Focus on trust and authority
- `PROFESSIONAL_CONTENT_TRANSFORMER` - CV bullets to impact stories
- `GO_ARCHITECTURE_MODERNIZER` - Technical foundation improvement
- `MOBILE_FIRST_TO_OPTIMIZER` - Mobile experience enhancement