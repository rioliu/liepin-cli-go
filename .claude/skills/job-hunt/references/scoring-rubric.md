# Job Scoring Rubric

## How to Score

Compare each JD against the candidate's **resume abilities** (extracted in Step 1). Score by actual skill match — **not** by job title.

### Dimensions (weighted)

| Dimension | Weight | What to check |
|-----------|--------|---------------|
| Skills match | 40% | Which required skills does the candidate have vs missing? Match against the structured ability list, not titles. |
| Experience fit | 25% | Years, seniority level alignment. Overqualified is a gap too (signals wrong level). |
| Industry/domain | 20% | Relevance of candidate's background to the job's domain. Hard domain requirements (e.g. "medical experience", "automotive ADAS") are critical gaps. |
| Education | 15% | Meets requirements? Degree level, major relevance. |

### Score Bands

| Score | Band | Meaning |
|-------|------|---------|
| 90-100 | Near-perfect fit | Strong apply — skills, domain, and level all align |
| 70-89 | Good fit | Apply — minor gaps only, none critical |
| 50-69 | Moderate fit | Significant gaps in one or two dimensions |
| Below 50 | Poor fit | Skip — core requirements not met |

### Critical Gap Rules

A "critical gap" is a hard requirement the candidate lacks. It caps the score:

- **Hard domain requirement** (e.g. "must have automotive experience", "medical/pharma background required") → cap at 65
- **Hard tech requirement the candidate lacks** (e.g. "must know ARM SoC", "Android framework expert") → cap at 55
- **Degree requirement not met** (e.g. "硕士及以上学历" when candidate has 本科) → cap at 60
- **Salary below candidate's floor** → note as gap but do not cap score (negotiable)

### Matching Principle

Match by **JD skill content**, not by job title. Many relevant roles have non-obvious titles:

- "Agent工程师" may match if JD requires Python + Agent architecture + CI/CD
- "系统工程师" may match if JD requires Linux + Python + automation
- "架构师" may not match if JD is about chip/hardware architecture

Always read the full JD before scoring. A title like "测试架构师" with ARM SoC requirements scores lower than "Agent工程师" with Python + CI/CD requirements.
