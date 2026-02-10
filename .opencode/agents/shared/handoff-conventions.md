# BrewPipes Inter-Agent Handoff Conventions

This file defines how agents communicate with each other — the formats for passing work, reporting results, and flagging issues. All agents should follow these conventions.

## API designer → developers

When the API designer produces a contract, it becomes the **source of truth** for implementation. Developers must treat the contract as authoritative. If a developer believes the contract is wrong, they must flag it — not silently deviate.

The API designer's output must include:

1. **Endpoint specification:** URL, HTTP method, authentication requirement.
2. **Request schema:** All fields with types, required/optional, validation rules.
3. **Response schema:** All fields with types, nullable indicators.
4. **Error responses:** Status codes with error message formats.
5. **TypeScript types:** For frontend consumption.
6. **Go DTO types:** For backend consumption, including `Validate()` methods.
7. **Example request/response:** At least one happy-path example.

## Tech lead review reports

When a tech lead (frontend or backend) reviews code, findings must be reported using severity levels:

### Severity levels

- **`[BLOCKER]`** — Must be fixed before the change can be accepted. Examples: data integrity risk, security vulnerability, broken existing functionality, missing error handling on critical paths.
- **`[ISSUE]`** — Should be fixed but doesn't block acceptance if there's a strong reason to defer. Examples: missing test coverage, inconsistent patterns, suboptimal query, accessibility gap.
- **`[NIT]`** — Minor style or preference item. Fix if convenient. Examples: naming, comment wording, import ordering.

### Report format

```
## Review: [what was reviewed]

### Summary
[1-2 sentence overall assessment. Lead with problems found, not affirmations.]

### Findings

[BLOCKER] file/path:line — Description of the problem and why it matters.
  Suggested fix: [concrete suggestion]

[ISSUE] file/path:line — Description of the problem.
  Suggested fix: [concrete suggestion]

[NIT] file/path:line — Description.

### Verdict
[ACCEPT | ACCEPT WITH CHANGES | REJECT]
- ACCEPT: No blockers, no significant issues.
- ACCEPT WITH CHANGES: No blockers, but issues should be addressed.
- REJECT: Blockers present; must be fixed and re-reviewed.
```

## QA engineer test reports

When the QA engineer runs tests or validates a feature, results must include:

```
## QA Report: [feature or area tested]

### Test Results
- Tests added: [count]
- Tests modified: [count]
- Tests passing: [count] / [total]
- Tests failing: [count] (list each with reason)

### Coverage
- Critical paths covered: [yes/no, list gaps]
- Error cases covered: [yes/no, list gaps]
- Mobile tested: [yes/no, breakpoints checked]

### Bugs Found
[BLOCKER] Description — Steps to reproduce. Expected vs actual behavior.
[ISSUE] Description — Steps to reproduce. Expected vs actual behavior.

### Recommendations
- [Any additional testing needed]
- [Any areas of concern]
```

## Brewing industry expert research briefs

When the brewing expert provides domain research, output must be structured as:

```
## Research Brief: [topic]

### Question
[What was asked]

### Findings
[Factual answer with sources or reasoning]

### Implications for BrewPipes
[How this affects data models, UI, calculations, or workflows]

### Recommendations
[Specific, actionable recommendations for the requesting agent]

### Confidence
[HIGH | MEDIUM | LOW] — and what would increase confidence
```

## UX designer deliverables

When the UX designer produces designs, output must include:

1. **User story:** Who and what they're trying to accomplish.
2. **User flow:** Entry → action → completion (numbered steps).
3. **Layout specification:** Wireframe or detailed description for each breakpoint.
4. **Component mapping:** Which Vuetify components to use.
5. **States:** Loading, empty, error, success.
6. **Mobile considerations:** Specific adaptations for small screens.

## Developer done criteria

When a developer (frontend or backend) completes a task, they must verify and report:

### Backend developer

1. `go build ./...` succeeds with no errors.
2. `go vet ./...` produces no warnings.
3. `go test ./...` passes (all existing + new tests).
4. New tests added for new/changed behavior.
5. Migrations have matching up/down files.
6. No unparameterized SQL.
7. No hardcoded secrets.

### Frontend developer

1. `pnpm build` succeeds with no errors (from `service/www/`).
2. `pnpm test run` passes (all existing + new tests, from `service/www/`).
3. New tests added for new/changed behavior.
4. No TypeScript `any` types without justification.
5. Mobile responsiveness verified at 375px.
6. No `console.log` in production code.

## General handoff principles

- **Be explicit about what you're handing off.** Don't assume the receiving agent has context from a previous conversation.
- **Include file paths.** Reference specific files with `file/path:line` format.
- **State what's done and what's not.** If you completed 3 of 4 subtasks, say so.
- **Flag risks and assumptions.** If you made a judgment call, document it.
- **Keep it concise.** The receiving agent's context window is limited.
