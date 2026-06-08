---
name: job-hunt
description: >
  Search for jobs on liepin.com, fetch JD details, AI-match against the candidate's resume,
  score and rank results, and optionally apply. Supports multi-keyword pipe-separated search,
  single URL, or job ID. Use when the user says "job hunt", "search jobs", "find jobs",
  "apply to jobs", or passes a liepin job URL/ID.
argument-hint: <keywords|URL|job-id>
allowed-tools: [Bash, Read, Agent, AskUserQuestion]
---

# Job Hunt

Search for jobs, AI-match against the candidate's resume, and apply.

## Input

$ARGUMENTS

Supported formats:
- Single keyword: `测试架构师`
- Multiple keywords (pipe-separated): `测试架构师|QA架构师|测试开发` — runs separate searches, merges results
- URL: `https://www.liepin.com/a/12345.shtml`
- Job ID: `12345`

## Steps

### 1. Fetch resume and extract preferences

Run `bin/liepin-cli resume get -o json` to get the full resume.

Extract **search filters**:
- **Location**: from 求职期望 → 期望地点 (e.g. 北京)
- **Salary floor**: from 求职期望 → 期望薪资下限 (e.g. 30000)
- **Work years**: from 开始工作年份 → calculate years of experience

Extract **resume abilities** as a structured list for matching. Group by category:
- **Languages**: e.g. Python (primary), Go, Java, Shell
- **Frameworks/Tools**: e.g. Ginkgo, Robot Framework, FastMCP, Jenkins
- **Platforms**: e.g. Kubernetes, OpenShift, Linux/RHEL, AWS
- **Domain expertise**: e.g. CI/CD pipelines, MCP server development, AI Agent integration, test automation
- **Leadership**: e.g. team management, cross-team collaboration

These become the matching basis. Do not hardcode them.

### 2. Determine mode

- If input starts with `http` → **single mode**: fetch one job detail, skip to step 4
- If input is a numeric ID → **single mode**: fetch with `--job-id`, skip to step 4
- Otherwise → **search mode**: continue to step 3

### 3. Search mode — fetch jobs

Split input by `|` into a list of keyword phrases. If no `|` is present, treat the entire input as a single keyword phrase.

For **each** keyword phrase, run a search in parallel:

```
bin/liepin-cli job search --job-name "<keyword>" --salary-floor <from-resume> --address <from-resume> --work-experience <from-resume> -o json
```

After all searches complete, **merge** the job lists and **deduplicate by jobId** (same jobId from different searches → keep one).

**Fetch JD details** for all unique jobs. Run fetches in parallel (up to 10 concurrently):

```
bin/liepin-cli job detail --job-detail-url "<jobDetailUrl>" -o json
```

If a fetch fails, skip that job and note it.

### 4. Score and rank

Read `references/scoring-rubric.md` for the full scoring methodology.

For each fetched JD, score against the resume abilities from Step 1:

1. **Skills match** (40%): Compare JD required skills against the structured ability list. Match by actual skill content, not job title.
2. **Experience fit** (25%): Years and seniority alignment. Overqualified is a gap.
3. **Industry/domain** (20%): Relevance of candidate's background. Hard domain requirements are critical gaps.
4. **Education** (15%): Meets degree/major requirements?

**Critical gap rules** — these cap the score:
- Hard domain requirement the candidate lacks → cap at 65
- Hard tech requirement the candidate lacks → cap at 55
- Degree requirement not met → cap at 60

**Key principle**: Score by JD content, not title. "Agent工程师" with Python + CI/CD may score higher than "测试架构师" with ARM SoC requirements.

### 5. Present results

**Single mode** — compact assessment:

```
## Job Match: <title> @ <company>

**Score: XX/100**

### Strengths
- ...

### Gaps
- ...

### Verdict
Apply / Skip / Borderline
```

**Search mode** — ranked table:

```
## Job Hunt Results: "<keywords>"

| # | Score | Job Title | Company | Salary | Key Match | Key Gap |
|---|-------|-----------|---------|--------|-----------|---------|
| 1 | 85    | ...       | ...     | ...    | ...       | ...     |
| 2 | 72    | ...       | ...     | ...    | ...       | ...     |

**Summary**: X jobs evaluated, Y scored 70+, Z scored below 50.
```

### 6. Niche detection

If fewer than 3 jobs score 70+ in search mode, analyze why and suggest alternatives:

- **Platform mismatch**: if the platform returns mostly unrelated roles (hardware, automotive, chip testing), suggest the user try different platforms (LinkedIn, Boss直聘, company career pages)
- **Keyword mismatch**: suggest alternative keywords that better target the candidate's niche
- **Niche too specific**: if the candidate's skillset is rare, suggest broadening keywords or targeting companies directly

### 7. Ask to apply

After showing results, ask the user:

- Single mode: "Want me to apply?"
- Search mode: "Want me to apply to jobs scoring 70+? Or pick specific ones by number?"

Apply with: `bin/liepin-cli job apply --job-id <jobId> --job-kind <jobKind> -o json`

Report results:

```
| Job | Company | Result |
|-----|---------|--------|
| ... | ...     | Applied / Failed / Already applied |
```

## Notes

- For multiple keywords, run searches **in parallel** for speed
- Deduplicate by jobId — same job found under different keywords is listed once
- If a keyword returns no results, note it and continue with other keywords
- Fetch JD details in parallel (up to 10 concurrent fetches)
- Always read the full JD before scoring — title alone is not reliable
