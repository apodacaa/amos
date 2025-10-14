```text
Claude 4.5 — minimal safety checklist for quick prototyping

- Always ask for a minimal patch/diff for any code change. Example: "Return only a unified diff for file X."
- Keep changes tiny: 1 logical change per prompt.
- Require a human review before committing anything the model suggests.
- Include environment constraints in prompts (Python version, formatters: black/ruff).
- Label AI-assisted commits/PRs with `ai-assisted`.
- Never accept code that runs elevated OS commands or eval() without explicit approval.
```