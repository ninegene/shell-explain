# Independent contributor trial

This is a facilitator plan, not evidence that a trial has happened. The repository became public on 2026-11-17 after the owner waived this trial as a gate, so no outside contributor has yet exercised the workflow. The plan remains available if a trial is wanted later. Do not create accounts, issues, forks, or pull requests just to run it without the necessary authorization.

## Before the session

- Choose a participant who did not build the project or its contribution instructions. Use their normal development environment and ask them not to inspect this facilitator plan before the trial.
- Confirm the participant has an authorized copy of the repository and can use their own disposable working directory. Record OS, architecture, Go, Node, npm, just, shell, and relevant host-tool versions. Do not record personal identifiers or machine names in the repository.
- Select a real, still-open, small task from [candidate starter issues](starter-issues.md), such as clarifying the description of one existing example. Confirm the page and proposed change are appropriate before the session; do not ask the participant to invent a defect or copy another project's prose. Supply the task and the [contribution guide](../CONTRIBUTING.md), not a step-by-step solution.
- Agree on a place for feedback and whether a draft change may be shared. Do not publish the participant's name, logs, or contact details without their agreement. Do not send an invitation on their behalf without authorization.

## Participant task

Ask the participant to work from the documented instructions with minimal help. Have them locate the relevant command page and upstream reference, make the agreed original change, run the documented checks appropriate to it, and explain the change and validation as they would in a PR. Ask them to record failures and skipped implementations. If a check is slow or a prerequisite is unavailable, record the outcome and continue; do not call that check passed.

The guide's `ps` examples assume a controlling terminal. Run `just test` in a terminal; record failures in noninteractive environments separately rather than treating them as a successful test. On Linux, follow the Playwright system-dependency guidance when needed. Ask the participant to note unclear or missing instructions before providing assistance.

## Observe and record

Capture the repository revision, date, task, OS and tool versions; the instructions and commands followed; pass/fail/skip results; time and points of confusion; and any assistance given. Check that the participant can identify the source file, preserve the content format, provide evidence for changed claims, understand the contribution licensing terms, and describe checks and limitations. Avoid recording private data or copying terminal output containing secrets.

If a participant cannot use GitHub, a patch or local branch may stand in for a PR; record that fork creation and the actual GitHub PR workflow were not tested. The repository is public, so a trial can verify the fork, clone, branch, push, PR template, review, and response workflow directly.

## Completion

The trial is complete only when an independent participant has attempted the relevant steps, the outcomes and limitations are recorded, and material blockers have been fixed or explicitly carried forward. Record the public fork/PR portion separately and keep it pending until actually exercised. Update the [publication review](publication-review.md) with evidence; do not treat this plan or the earlier container run as a completed trial.
