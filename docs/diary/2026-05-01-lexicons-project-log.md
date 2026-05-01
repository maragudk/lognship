# Diary: Add `com.lognship.project` and `com.lognship.log.entry` lexicons

Adding two new record lexicons alongside the existing `com.lognship.actor.profile`: a per-actor `project` (multi-record) and a `log.entry` (issue-shaped, log-stream content). Also tightening profile and project to require `createdAt` and document it. None of the lexicons are public yet, so breaking changes are still free.

## Step 1: Create `com.lognship.project` lexicon

**Author:** main

### Prompt Context

**Verbatim prompt:** "Let's create a new lexicon com.lognship.project with records displayName, description, and createdAt, mirroring the profile."
**Interpretation:** Create a new record lexicon with the same shape as `com.lognship.actor.profile` but with the three named fields.
**Inferred intent:** Establish the lognship "project" concept (analogous to a GitHub repo) as a first-class atproto record so actors can publish projects to their PDS.

### What I did

Created `/lexicons/com/lognship/project.json` with NSID `com.lognship.project`. Mirrored the profile structure: same `maxGraphemes`/`maxLength` on `displayName` and `description`, same `datetime` format on `createdAt`. Used `key: "tid"` instead of `literal:self` because an actor will have multiple projects, not one.

Verified with `make lint` -- both `goat lex parse` and `goat lex lint` reported success on all three files.

### Why

The user explicitly asked for a mirror of profile minus `website`. The only real design decision was the record key, and `tid` is the sensible default for collections.

### What worked

Reading `/lexicons/com/lognship/actor/profile.json` first gave me the exact field constraints to mirror. The `make lint` target made validation a one-liner.

### What didn't work

Nothing in this step.

### What I learned

The profile lives at `lexicons/com/lognship/actor/profile.json` (4-segment NSID), but the user chose `com.lognship.project` (3 segments). Both are valid -- 3 is the NSID minimum -- but it puts `project.json` directly under `lexicons/com/lognship/` rather than in a grouping subdirectory. Worth flagging in case a later rename to e.g. `com.lognship.project.record` is wanted, since NSID changes are breaking.

### What was tricky

None.

### What warrants review

The choice of NSID depth (`com.lognship.project` vs. e.g. `com.lognship.project.record`). I flagged this to the user; they implicitly accepted by moving on.

### Future work

None implied at this step.

## Step 2: Create `com.lognship.log.entry` lexicon

**Author:** main

### Prompt Context

**Verbatim prompt:** The user pasted a proposed lexicon JSON for `com.lognship.log.entry` with `text` (`maxLength: 3000`, `maxGraphemes: 300`) and `createdAt`, both required, and asked for it to be created.
**Interpretation:** Write the file exactly as pasted.
**Inferred intent:** Establish the log-entry record type -- the core content unit of lognship, analogous to a GitHub issue/comment.

### What I did

Wrote `/lexicons/com/lognship/log/entry.json` with the user's content verbatim. NSID `com.lognship.log.entry`, key `tid`, two required fields. `make lint` passed clean on all three lexicons.

### Why

The user provided the exact schema. No interpretation needed.

### What worked

Direct write, no surprises.

### What didn't work

Nothing.

### What I learned

Nothing new in this step.

### What was tricky

None.

### What warrants review

The `text` cap (3000 bytes / 300 graphemes) was chosen by analogy with Bluesky posts. We later revisited this -- see Step 4.

### Future work

None implied at this step.

## Step 3: Require `createdAt` on profile and project; add descriptions; add `website` to project

**Author:** main

### Prompt Context

A series of small edits across both lexicons:
- "Also add `\"required\": [\"createdAt\"]` to the profile"
- "Add this description to profile createdAt: `\"description\": \"Client-declared timestamp when this account profile was originally created.\"`"
- "Also make the same changes to createdAt as for profile" (applied to project)
- "add a website to the project"

**Interpretation:** Tighten `createdAt` to be required on both records, document it inline, and bring `website` over from profile to project.
**Inferred intent:** Bring profile and project to a shippable shape with consistent metadata. The user noted "None of this is public yet" to authorise the breaking change.

### What I did

On `/lexicons/com/lognship/actor/profile.json`:
- Added `"required": ["createdAt"]` on the record object.
- Expanded `createdAt` from `{ "type": "string", "format": "datetime" }` to a multi-line object with the user-supplied description.

On `/lexicons/com/lognship/project.json`:
- Mirrored both changes with description "Client-declared timestamp when this project was originally created."
- Added `"website": { "type": "string", "format": "uri" }` between `description` and `createdAt`, matching the profile's placement.

Re-ran `make lint` after each edit. All three lexicons stayed green.

### Why

`createdAt` being required is a real constraint -- log entries already require it, and an undated profile/project record would be a foot-gun for downstream indexers. The descriptions show up in generated client docs. The `website` field on project parallels profile and gives projects a homepage URL.

### What worked

The user's "None of this is public yet" cleared the path for breaking edits without ceremony. I noted the evolution-rule concern once after the first required-field add ("adding a required field is technically a breaking lexicon change") so it's recorded for future edits, then dropped it.

### What didn't work

Nothing.

### What I learned

The lognship project is in a phase where lexicon edits are still cheap. After publication, every one of these tweaks would have required a new NSID. Worth keeping that boundary visible -- the moment a lexicon's records start federating, the rules harden.

### What was tricky

None.

### What warrants review

Whether `website` belongs on a project record at all -- many projects will have a repo URL but no separate marketing site. Currently optional, so harmless. Reviewer should confirm the field is what they want, or rename it (e.g. `homepageUrl`).

### Future work

None implied -- the user moved on to the log-entry text cap.

## Step 4: Revisit `text` cap on `log.entry`; end with no cap

**Author:** main + sub-agent (Tangled/Leaflet research)

### Prompt Context

**Verbatim prompts:**
- "I'm worried about the max length for the log entry. Is there a max length for github issue text?"
- "Characters or graphemes?"
- "Just set maxlength for now, to 65536"
- "Research in a subagent what e.g. Leaflet or Tangled do in the atmosphere."
- "And maybe also what the linter considers large in another subagents"
- (interrupt) "Actually, just remove the limit for now"
- "Let's live with it for now."

**Interpretation:** Find a defensible cap for issue-shaped log entry bodies, then -- after research -- commit to a no-cap policy and accept the resulting `make lint` warning.
**Inferred intent:** Reach a `text`-field constraint that matches how issue-shaped content actually behaves on atproto today, rather than guessing from Bluesky's post limits.

### What I did

First answered from general knowledge: GitHub caps issue/PR/comment bodies at 65,536 characters. Distinguished `maxLength` (UTF-8 bytes) vs. `maxGraphemes` (user-perceived length) and noted GitHub's "characters" doesn't map cleanly to either.

On the user's instruction, set `maxLength: 65536` and dropped `maxGraphemes`. `make lint` failed:

```
🟡 lexicons/com/lognship/log/entry.json
   [large-string]: large max length
error: linting issues detected
```

Dispatched a sub-agent to survey Tangled and Leaflet (the two closest atmosphere precedents -- Git forge and long-form publishing). Findings:

- **Tangled** (`sh.tangled.repo.issue.body`, PR body, all comment bodies) -- bare `{"type": "string"}`, no `maxLength`, no `maxGraphemes`. Closest analogue to lognship; ships completely uncapped.
- **Leaflet** -- uses `maxGraphemes`/`maxLength` only on short metadata (title 500/5000, description 3000/30000); body content is split across `pub.leaflet.blocks.*` records, all uncapped.

Started a second sub-agent to find the goat lint threshold. The user interrupted with "Actually, just remove the limit for now" before that ran, so I removed both `maxLength` and `maxGraphemes` on `text`. `make lint` then flagged the opposite warning:

```
🟡 lexicons/com/lognship/log/entry.json
   [unlimited-string]: no max length
error: linting issues detected
```

Presented options; the user chose to live with it.

### Why

Tangled is the most direct analogue to lognship and they explicitly chose no cap, letting the PDS' record-size limit (~1 MB) bound things. Adopting their convention is the cheapest path to "matches how the atmosphere actually works for this content type." The lint warning is informational, not blocking, and `goat lex parse` still succeeds.

### What worked

The Tangled/Leaflet research was the right call -- it replaced a guess (65536 from GitHub docs) with an evidence-based "the convention here is no cap." Sub-agent returned in one shot with file-level URLs.

### What didn't work

The `goat lex lint` rule fires both ways: 65,536 trips `[large-string]`, and "no cap" trips `[unlimited-string]`. There's no obvious middle that satisfies both rules without research into the threshold. The user chose to short-circuit by accepting the warning rather than chasing the exact bound.

### What I learned

`goat lex lint` enforces a "Goldilocks" length policy that real atmosphere apps don't follow. Tangled's lexicons would presumably trip the same lint warnings, which strongly suggests the lint rule is informational and ignored in practice. The Makefile `make lint` target makes this surface as a hard CI failure, though, which is worth knowing if CI gets wired up.

### What was tricky

The lint failure on uncapped strings was unexpected after the failure on overly-large strings -- the linter complains about both extremes. Without the research detour I'd have suggested a guess like 25000 bytes, which would have been pure cargo-culting.

### What warrants review

The Makefile `lint` target now exits non-zero. If CI runs `make lint`, it will fail. Reviewer should decide whether to:
1. Filter the `[unlimited-string]` rule (need to research goat's flags).
2. Split `make lint` so `parse` is hard-required and `lint` is advisory.
3. Accept that `make lint` is currently broken until the rule is silenced.

The `text` field having no `maxGraphemes` also means downstream indexers can't easily implement "show first N graphemes" without computing it themselves. That's fine -- Tangled does the same thing.

### Future work

- Decide on a `make lint` strategy (see "What warrants review" above).
- Consider whether `log.entry` should support markdown/rich text or stays plain. Tangled and Leaflet both store plaintext at the lexicon layer; Bluesky uses `facets` for richer structure. Defer until lognship has a UI to render anything.
- The first sub-agent's research notes for Leaflet/Tangled are worth keeping accessible -- if the user wants to formalise lognship's lexicon style guide later, those URLs are the source material.
