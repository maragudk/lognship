# Diary: Central git remote and first repo record on the atmosphere

The lognship codebase finally has a canonical home -- a Soft-Serve-style git server at `ssh://server0.maragu.net:23231/lognship` -- and the atmosphere now reflects that fact through the first `com.lognship.repo` record, linking the git URL back to the lognship project record. Two related moves in one short session: tell the world where the code lives, and tell the atmosphere what we just told the world.

## Step 1: Push `main` to the new central remote, capture it in AGENTS.md

**Author:** main

### Prompt Context

**Verbatim prompt:** "You've got a new git remote origin!"
**Interpretation:** The user added a new `origin` to the local repo and wants me to push `main` up to it.
**Inferred intent:** Establish the new remote as the canonical home for the project's code and have the local working copy track it, so future `git push`/`pull` flows the right way out of the box.

### What I did

Ran `git remote -v` to see what got configured. The remote was `ssh://server0.maragu.net:23231/lognship` (a self-hosted git server -- non-standard SSH port, naming smells like Soft Serve or similar).

After the user confirmed the push, ran `git push -u origin main`. The first attempt was denied by policy (`Pushing directly to the default branch main bypasses pull request review.`). I stopped, explained the block to the user, and only retried after Markus said "I'll allow it" -- the explicit per-action override, not a blanket consent. The retry returned `Everything up-to-date` and `branch 'main' set up to track 'origin/main'.`

Then Markus asked to note the remote in `AGENTS.md`. I read the existing file (4 prose lines: project description, the "no-`gh`-CLI" note, the atproto-skill instruction, and the lexicons location), and folded the remote URL into the existing `gh` line so the "this lives outside GitHub, here's where it actually lives" idea reads as a single thought:

```
This project is not hosted on Github, so the `gh` CLI cannot be used. The central git repo is at `ssh://server0.maragu.net:23231/lognship` (remote `origin`).
```

A linter or hook then mirrored the same change into `CLAUDE.md` -- a system reminder told me about it, no action needed.

### Why

A canonical remote is a precondition for the next obvious atmosphere move: writing a `com.lognship.repo` record that points at *something*. Without an upstream URL, the repo record has nothing to reference. Putting the URL in `AGENTS.md` (and now `CLAUDE.md` by mirror) means future agents don't have to re-derive it from `git remote -v` -- they can read it once at session start and know which server holds the code.

### What worked

The "stop on policy denial, surface the reason, ask for explicit override" loop worked exactly as designed. The user got to make the call about whether direct-pushing to `main` was acceptable for this case (it was -- there was nothing to push, and the goal was to set up tracking, not actually deliver code). The single-line edit to `AGENTS.md` kept the file's tone and didn't introduce a new section just for one URL.

### What didn't work

Nothing in this step. The push being a no-op (`Everything up-to-date`) was a small surprise -- the remote already had the same commits, presumably because the user had pushed before adding `origin` locally or because the server had been seeded out of band -- but it didn't matter for the goal of establishing tracking.

### What I learned

The denial reason `Pushing directly to the default branch main bypasses pull request review` is a generic permission rule, not anything specific to this remote. Worth remembering: even on a brand-new remote with no PR system wired up yet, the local policy still treats `main` as the default branch and gates the push.

### What was tricky

Calibrating the override request. The user had said "yes push it" before the policy block, which sounded like consent at the time -- but the policy denial materially changed what they were consenting to (a normal push vs. a bypass-review push). Explaining the new shape and re-asking ("do you want to override?") rather than treating the prior yes as still-valid was the right read. Markus's follow-up ("I'll allow it") confirmed the second consent was the load-bearing one.

### What warrants review

`AGENTS.md` line 3 -- the URL is right (matches `git remote -v` output) but the prose now packs two ideas onto one line. If this file grows further, splitting "GitHub-not-applicable" and "central remote location" into separate lines would read more cleanly.

### Future work

Once there's more than `main`, a "branches use plain names, no `feat/` or `hotfix/` prefixes" note from the user's `git` skill could land in `AGENTS.md` too. Out of scope here.

## Step 2: Write the first `com.lognship.repo` record

**Author:** main

### Prompt Context

**Verbatim prompt:** "Now add a repo record with the lognship.com user to the atmosphere . Validate it against the lexicon first, ask me for confirmation before actually creating it"
**Interpretation:** Draft a `com.lognship.repo` record body, validate it locally with `goat lex validate`, present the proposed record and target for confirmation, then -- only after explicit yes -- write it to the lognship.com PDS via `goat record create -n`.
**Inferred intent:** Close the loop between the new central git remote (Step 1) and the lognship project on the atmosphere -- the first record of the third lexicon, and the first cross-record reference (repo → project) we'd ship.

### What I did

Read `/lexicons/com/lognship/repo.json` to confirm the required fields: `project` (an AT URI), `url` (a URI), and `createdAt` (datetime). Cross-checked the recent commit `Switch cross-record references to bare AT URI` -- so `project` is a plain `at://...` string, not a strong-ref object. Looked up the existing project record (`at://did:plc:ln7l376244zkvhwu4wgr3rzz/com.lognship.project/3mkzclpilfo2v`) via the diary at `/docs/diary/2026-05-04-first-records-on-the-atmosphere.md`.

Verified the active goat session was lognship via `bin/goat-as --which` -- it was. Drafted the record at `/tmp/lognship-records/repo-1.json`:

```json
{
  "$type": "com.lognship.repo",
  "project": "at://did:plc:ln7l376244zkvhwu4wgr3rzz/com.lognship.project/3mkzclpilfo2v",
  "url": "ssh://server0.maragu.net:23231/lognship",
  "createdAt": "2026-05-04T12:00:00Z"
}
```

Validated locally: `goat lex validate -c lexicons /tmp/lognship-records/repo-1.json` -> `valid com.lognship.repo record`.

Presented the proposed record for confirmation. The user asked "What's the project?" first -- so I read the project record back via `goat record get at://did:plc:ln7l376244zkvhwu4wgr3rzz/com.lognship.project/3mkzclpilfo2v` and showed the contents (`displayName: lognship`, website `https://www.lognship.com`). After Markus confirmed, ran `goat record create -n /tmp/lognship-records/repo-1.json`. The record landed at `at://did:plc:ln7l376244zkvhwu4wgr3rzz/com.lognship.repo/3mkzo3of5wu2x` (CID `bafyreicaw6fhsrda5yr7q77ffmbjf2shnqtev24blrmowyzfo3mn5ws4li`).

Read it back with `goat record get` -- the JSON came back identical to what I wrote (alphabetised by key, as expected).

### Why

The `com.lognship.repo` lexicon was published in the previous session but had no records yet. Writing the first one with the most concrete possible reference -- the lognship git server we set up in Step 1, linked to the lognship project record -- gives the third lexicon a real-world example and exercises the bare-AT-URI cross-reference shape end-to-end. It also means `at://lognship.com/com.lognship.repo/*` is no longer empty, so an app view that joins repos to projects has something to render.

### What worked

The "draft → local validate → ask → write → read back" loop is now a known-good rhythm for this project; it's the third time we've done it (profile, project, log entry, now repo). Showing the user the *project* record's contents when they asked "what's the project?" -- rather than just naming it -- was useful because it let them confirm the link target was the lognship-the-project record, not some other one.

### What didn't work

Nothing in this step. The `-n` (no-validate) flag is still required for `record create` against the bsky-network PDS -- the same property we hit and documented in the previous diary -- but that's an environmental quirk, not a step failure.

### What I learned

The repo record's `url` field is a `format: uri` string, which the lexicon validator accepts for `ssh://` URIs without issue. So self-hosted git servers slot in cleanly; the field isn't restricted to `https://` or `git://`. That's a small but useful confirmation if we ever index repos by URL scheme later.

### What was tricky

Picking `createdAt`. The lexicon describes it as "Client-declared timestamp when this repo was originally linked", which is the linking moment, not the moment the git remote was first created on the server. So I used a fresh timestamp (the session's current time) rather than backdating to when the server-side repo first existed. Worth noting because if someone later writes a tool that backfills repo records for older remotes, the convention here implies "now-ish" rather than "earliest-known".

### What warrants review

The CID `bafyreicaw6fhsrda5yr7q77ffmbjf2shnqtev24blrmowyzfo3mn5ws4li` for `at://lognship.com/com.lognship.repo/3mkzo3of5wu2x` is pinned in this diary. If `goat record get` ever returns a different CID for that AT URI, something rewrote our first repo record. Likely benign cause: `record update`. Anything else is worth a closer look.

The `project` field uses the lognship.com DID's project AT URI. If that DID is ever migrated, the AT URI form `at://did:plc:.../com.lognship.project/3mkzclpilfo2v` keeps working unchanged (records are DID-bound), but the human-readable form `at://lognship.com/...` would re-resolve via whatever handle the DID points at.

### Future work

A second repo record from the maragubot account (similar to how the cross-account log entry exercised multi-repo writes) would confirm that repo records work the same way from a non-publisher account. Not needed today -- there's only one repo to point at.
