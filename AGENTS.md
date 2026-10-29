# Working agreement

- Before creating, editing or deleting files, show the proposed diff and wait for explicit approval.
- Apply only the approved changes. Agreement with an approach does not approve unseen code.
- If additional changes are needed, propose them separately before editing.
- Commit, push and change infrastructure only when explicitly requested.
- Prefer native Woodpecker plugins and configuration over custom scripts.
- Keep proposals, comments and documentation short: what, why, how to verify.
- Distinguish proposed, implemented and tested behavior.
- Implement CI functionality in RevoTale/ci.
- Propose extraction into RevoTale/woodpecker-toolkit only when it provides
  clear general-purpose value to users outside RevoTale.
- Do not extract code merely to share it between our repositories.

# Decisions

- Read DECISIONS.md before proposing or making changes.
- Check relevant decisions against the current implementation and sources.
  Explain conflicts; do not silently override agreed decisions.
- After the user approves a key decision, update DECISIONS.md automatically.
  Recording it needs no separate approval; implementation still does.
- Record decisions that affect future work: architecture, dependencies,
  security, interfaces and working agreements.
- Use "We do / Why". Keep wording short and concrete.
  Include the important tradeoff or rejected alternative when relevant.
- Keep one entry per decision. Update changed decisions, merge duplicates
  and remove obsolete entries. Git preserves their history.
- Record agreement, not implementation or test status.
- Record the user's reason for a decision, not your inferred justification.
- If the reason is unclear or missing, ask before recording it.
- You may shorten the wording, but must preserve the user's meaning.
