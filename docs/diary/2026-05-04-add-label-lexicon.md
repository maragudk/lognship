# Diary: Add `com.lognship.log.label` lexicon

The sixth lognship lexicon: actor-applied labels on log entries. Same shape as `log.reaction` (strongRef-pinned subject, short bounded text, datetime), different intent: a label is a tag the actor wants to attach to an entry, not an opinion about it. The consistency pass that fell out of this also caught a missing `createdAt` description on `log.reaction` from earlier in the session, so this diary covers both.

## Step 1: Draft the lexicon, then tighten it for codebase consistency

**Author:** main

### Prompt Context

**Verbatim prompt:** "In the meantime, let's add log entry labels. Make this lexicon consistent with the rest of the codebase.\n\n{ ... }" -- followed by a JSON draft of `com.lognship.log.label` with `subject` as `com.atproto.repo.strongRef`, `text` bounded at 32/320, and `createdAt` as a bare datetime.
**Interpretation:** Take the user's draft, sanity-check every detail against the existing five lexicons (record description style, field naming, bounds, createdAt shape, indentation, ordering), fix anything inconsistent, lint clean, ship.
**Inferred intent:** Add a fluid actor-side labelling vocabulary to the project so consumers can attach short tags to specific entry versions, without any of it slipping into per-author drift from the existing schemas.

### What I did

Wrote `/lexicons/com/lognship/log/label.json` directly from Markus's draft. Cross-checked it against `/lexicons/com/lognship/log/reaction.json` (the closest cousin) and the other four lexicons. The structural shape, indentation, required-array order, property order, and 32/320 bound on `text` all matched the reaction lexicon. Ran `make lint` -- all six green on first try.

Reported "ship as-is" to Markus. He pushed back: the record description was too verbose ("A label applied to a log entry by the actor under whose DID this record lives.") and the `createdAt` should have a description, since other lexicons' timestamps do. I checked: `actor.profile`, `project`, `repo`, `log.entry` all have `createdAt` descriptions following a consistent pattern ("Client-declared timestamp when this <thing> was originally <verb>."). `log.reaction` was the outlier -- I'd shipped it earlier this session without one. Markus's original `log.reaction` draft had also omitted it; I'd followed the draft rather than the broader convention.

Three edits:

1. `log/label.json` record description: `"A label applied to a log entry by the actor under whose DID this record lives."` -> `"A label applied to a log entry."`. The dropped clause was true but inferable from the lexicon's location (`com.lognship.log`, not a labeler context); the terser form matches `log.reaction`'s `"A reaction to a log entry."`.
2. `log/label.json` `createdAt` gained `"Client-declared timestamp when this label was originally applied."`. "Applied" mirrors `repo`'s "linked" -- specific verb when the lexicon's action verb is more precise than "created".
3. `log/reaction.json` `createdAt` gained `"Client-declared timestamp when this reaction was originally created."`. This is a non-breaking edit to an already-published lexicon; on the wire nothing changes, but `goat lex status` will go red on `log.reaction` until we republish.

Re-ran `make lint`: all six green.

### Why

A quiet consistency drift in lexicon descriptions becomes a loud one once we have a dozen lexicons -- this is the cheap moment to fix it. And catching `log.reaction`'s missing description here, while we have the file open and the convention freshly re-derived, is meaningfully cheaper than catching it later when the next-but-one lexicon also omits it and the pattern has dissolved.

### What worked

The "ship the draft, get pushback, audit every lexicon for the convention being violated" loop. The grep across `createdAt` blocks made the divergence obvious in one screenful -- no need to read each lexicon end-to-end. Once Markus flagged "verbose" and "other timestamps have descriptions", the right edits were mechanical.

The pattern-naming approach to verbs ("created" for things that come into existence, "applied" for tagging actions, "linked" for repo connections) makes future additions easy: pick the verb that matches the action, drop the boilerplate around it. No shrubbery.

### What didn't work

My initial review of the draft missed both inconsistencies. I'd looked field-by-field at `log.reaction` as the comparator and concluded "matches"; I should have looked at the broader cohort, where `log.reaction` itself was the odd one out. The lesson: when a new lexicon "matches" the most-recent one, also compare against the rest -- the recent precedent might itself be wrong.

### What I learned

`knownValues` is genuinely open -- it's a hint, not a constraint. So `log.reaction.value` accepts strings beyond `["like", "dislike"]` (a `maxLength: 320` string can be anything). That same openness means `log.label.text` doesn't *need* `knownValues`: any short string is a valid label. The two fields end up with the same shape (string + bounds) for the same reason, and we don't need to enumerate the future label vocabulary to ship a label lexicon today.

Adding a `description` to a property of a record-type lexicon is a non-breaking change. The serialized DAG-CBOR shape of records doesn't include lexicon descriptions; the schema and the data are separate. So fixing `log.reaction` mid-session is safe -- existing reaction records still validate -- it just means re-running `goat lex publish` to update the schema record on the network.

### What was tricky

Choosing the verb for `log.label`'s `createdAt` description. "Created" was the safe default and matches profile/project/entry, but labels aren't created in the same sense -- they're applied to something. "Applied" mirrors repo's "linked" specificity. The risk is bikeshedding: every future lexicon could spawn a new verb. The right rule going forward is: if there's a natural action verb for the lexicon's primary action (link, apply, react, post, comment), use it; otherwise default to "created".

The `log.reaction` patch sits in awkward territory -- it's a fix to a lexicon I shipped this morning, in the same session, by my own hand. The diary skill says don't modify old diary entries; the original `log.reaction` diary at `/docs/diary/2026-05-04-add-reaction-lexicon.md` correctly described what I shipped, so it stays. This new diary documents both edits because they're the same logical move.

### What warrants review

Two files: `/lexicons/com/lognship/log/label.json` end-to-end (this is the new lexicon) and `/lexicons/com/lognship/log/reaction.json` -- specifically the new `createdAt` description. Both should match the convention now visible across `actor.profile`, `project`, `repo`, and `log.entry`. If a future contributor wants to add a seventh lexicon, the convention is "createdAt has a description that says 'Client-declared timestamp when this <thing> was originally <verb>.', and the verb is whatever the lexicon's primary action is."

Also worth eyeballing: `log.reaction` is now divergent from its on-network schema until we re-publish. `goat lex status` will be the canary -- if we forget the republish, the next person to run it sees a red row.

### Future work

Publish both lexicons (`goat lex publish lexicons/`). Write a first label record from one of the accounts to exercise the strongRef wire format end-to-end (mirrors what we did for reaction). Both deferred -- the user wanted the diary written and the change committed before any network writes.

The same audit could be applied to record descriptions across all lexicons -- they vary in tone (`"A label applied to a log entry."` vs. `"Record containing a lognship log entry."` vs. `"A declaration of a lognship project."`). Not addressed here; the user only flagged the new file's description as too long, not the broader inconsistency. Worth a future "tone pass" if it ever bothers anyone.
