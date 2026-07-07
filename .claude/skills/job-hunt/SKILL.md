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

Extract **resume abilities** as a structured list for matching. Group by category:
- **Languages**: e.g. Python (primary), Go, Java, Shell
- **Frameworks/Tools**: e.g. Ginkgo, Robot Framework, FastMCP, Jenkins
- **Platforms**: e.g. Kubernetes, OpenShift, Linux/RHEL, AWS
- **Domain expertise**: e.g. CI/CD pipelines, MCP server development, AI Agent integration, test automation
- **Leadership**: e.g. team management, cross-team collaboration

These become the matching basis. Do not hardcode them.

Determine **internet company filter** — internet companies (互联网) have implicit age restrictions. Apply this rule:

1. **Calculate age** from 生日 field (e.g. 19810301 → ~45 in 2026)
2. **Check industry background** from 当前行业 and work history industries

| Age | Industry Background | Filter |
|-----|-------------------|--------|
| < 40 | Any | **Include** internet companies |
| ≥ 40 | 互联网 career | **Include** internet companies (familiar territory) |
| ≥ 40 | Non-互联网 career (计算机软件, IT服务, 云计算, 企业软件, etc.) | **Exclude** internet companies |

Reason: internet companies favor younger candidates; non-internet candidates over 40 face significant bias there and have better odds at enterprise/foreign companies.

Note the age, industry classification, and filter decision in the results header.

### 2. Determine mode

- If input starts with `http` → **single mode**: fetch one job detail, skip to step 4
- If input is a numeric ID → **single mode**: fetch with `--job-id`, skip to step 4
- Otherwise → **search mode**: continue to step 3

### 3. Search mode — fetch jobs

Split input by `|` into a list of keyword phrases. If no `|` is present, treat the entire input as a single keyword phrase.

For **each** keyword phrase, run a search in parallel:

```
bin/liepin-cli job search --job-name "<keyword>" --salary-floor <from-resume> --address <from-resume> -o json
```

After all searches complete, **merge** the job lists and **deduplicate by jobId** (same jobId from different searches → keep one).

**Apply industry filter** (from Step 1):
- If the candidate is NOT from the internet industry → remove jobs where `industry` contains "互联网"
- Log how many jobs were filtered and why

**Fetch JD details** for remaining unique jobs. Run fetches in parallel (up to 10 concurrently):

```
bin/liepin-cli job detail --job-detail-url "<jobDetailUrl>" -o json
```

If a fetch fails, skip that job and note it.

### 4. Score and rank

Read `references/scoring-rubric.md` for the full scoring methodology.

For each fetched JD, score against the resume abilities from Step 1:

1. **Skills match** (50%): Compare JD required skills against the structured ability list. Match by actual skill content, not job title.
2. **Industry/domain** (30%): Relevance of candidate's background. Hard domain requirements are critical gaps.
3. **Education** (20%): Meets degree/major requirements?

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

## Bulk Apply Rules

Liepin enforces a **daily application limit** (~50-80 per day). The server returns HTTP 200 with `"您的投递已达上限"` in the response body when the cap is hit. HTTP 429 rate limiting is handled by the CLI automatically.

### Response handling:

| Response body contains | Meaning | Action |
|------------------------|---------|--------|
| `应聘成功` | Applied successfully | Continue |
| `您已投递过该职位` | Already applied | Skip, continue |
| `您的投递已达上限` | Daily app limit reached | **Stop immediately** — no retry |
| `rateLimited: true` (error JSON) | HTTP 429 or app rate limit (code 429xxx) — handled by CLI | Already retried automatically |
| Other error | Transient failure | Skip, continue |

### When applying to multiple jobs:

1. **Add a delay** between requests: `sleep 0.5` to avoid triggering HTTP 429
2. **Check for "已达上限"**: if found, **stop immediately** — further attempts will also fail
3. **Report the cap**: tell the user how many succeeded and how many remain

### Bulk apply strategy:

- **Score first, then apply**: use the daily quota on the best matches
- **Prioritize 70+ scored jobs**: apply to highest-scoring jobs first
