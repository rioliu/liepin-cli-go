Search for jobs, AI-match against my resume, and apply. Works in two modes:

- **Search mode**: pass keywords → search, fetch details, match all, rank, apply
- **Single mode**: pass a URL or job ID → fetch detail, match, ask to apply

## Input

$ARGUMENTS

- Keywords: "Go Kubernetes", "Python AI platform", "DevOps SRE"
- URL: `https://www.liepin.com/a/12345.shtml`
- Job ID: `12345`

## Steps

### 1. Fetch resume and extract preferences

Run `bin/liepin-cli resume get -o json` to get my full resume.

From the resume, extract:
- **Location**: from 求职期望 → 期望地点 (e.g. 北京)
- **Salary floor**: from 求职期望 → 期望薪资下限 (e.g. 30000)
- **Work years**: from 开始工作年份 → calculate years of experience

These become the search filters. Do not hardcode them.

### 2. Determine mode

- If input starts with `http` → **single mode**: fetch one job detail, skip to step 4
- If input is a numeric ID → **single mode**: fetch with `--job-id`, skip to step 4
- Otherwise → **search mode**: continue to step 3

### 3. Search mode — fetch jobs

Run search using the extracted preferences:

```
bin/liepin-cli job search --job-name "<keywords>" --salary-floor <from-resume> --address <from-resume> --work-experience <from-resume> -o json
```

Extract the job list. For each job (up to 10), fetch full JD:

```
bin/liepin-cli job detail --job-detail-url "<jobDetailUrl>" -o json
```

If a fetch fails, skip that job and note it.

### 4. AI match and rank

Compare every fetched JD against my resume. For each job, assess:

- **Skills match**: which required skills do I have vs missing?
- **Experience fit**: years, seniority level alignment
- **Industry/domain**: relevance of my background
- **Education**: meets requirements?

Assign a score 0-100:
- 90-100: near-perfect fit, strong apply
- 70-89: good fit, minor gaps
- 50-69: moderate fit, significant gaps
- Below 50: poor fit, skip

### 5. Present results

**Single mode** — show a compact assessment:

```
## Job Match: <title> @ <company>

**Score: XX/100**

### Strengths / Gaps / Verdict
```

**Search mode** — show a ranked table:

```
## Job Hunt Results: "<keywords>"

| # | Score | Job Title | Company | Salary | Key Match | Key Gap |
|---|-------|-----------|---------|--------|-----------|---------|
| 1 | 85    | ...       | ...     | ...    | ...       | ...     |
| 2 | 72    | ...       | ...     | ...    | ...       | ...     |

**Summary**: X jobs found, Y scored 70+, Z scored below 50.
```

### 6. Ask to apply

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

- If search returns no results, suggest alternative keywords
- If fewer than 3 jobs score 70+, suggest broadening the search
- Deduplicate by jobId if the same job appears in multiple searches
