# Diary: Project membership lexicon

The piece that was missing before lognship could start being used by more than one DID: a way to say "this person is a collaborator on this project". We brainstormed the data model, settled on owner-side membership records, wrote and published `com.lognship.project.membership`, attached the first two members (maragubot.com and markus.maragu.dev) to the lognship project, and refreshed the project's display copy now that the project actually has collaborators.

## Step 1: Brainstorm and write the lexicon

**Author:** main

### Prompt Context

**Verbatim prompt:** "I have one missing piece of the design to start using it: how collaborators are attached to the project. See the diary for recent work, then let's brainstorm."
**Interpretation:** Walk through the design space for "collaborator on a project" with multiple-choice questions, one per turn, until we have a fully formed design -- then write the lexicon JSON.
**Inferred intent:** Lock down the single missing data model before any UI/app-view work starts, so downstream tooling has a stable schema to read against.

### What I did

Read the most recent diary (`/docs/diary/2026-05-04-first-records-on-the-atmosphere.md`) and the existing lexicons under `/lexicons/com/lognship/` to ground the brainstorm in current state. Invoked the brainstorm skill and worked through the design with one question per turn:

1. **Meaning of "collaborator"**: someone whose work shows up in the app view by default. Git access is out of scope.
2. **Direction of the relationship**: owner-side only -- the project owner publishes the membership records. Markus rejected the collaborator-side option (anyone could "join" without consent) and the two-sided handshake (collaborators wouldn't necessarily care to ratify).
3. **Record shape**: one record per (project, person), not one list-record per project. Reasons in priority order: firehose-friendly, per-record metadata is free, deletion is just `record delete`, and concurrent writes don't clobber each other.
4. **Optional fields**: minimal now (`project`, `subject`, `createdAt`), grow later. No `role`, no `note`, no `invitedBy`.
5. **rkey shape**: `tid`, consistent with the rest of the project.
6. **Project reference**: plain `at-uri` (matching `repo` and `log.entry`), not `strongRef` -- a project rename shouldn't invalidate memberships.
7. **Subject type**: `format: did`, not `at-identifier` -- a handle change shouldn't silently transfer membership.
8. **Field name for the subject**: Markus didn't like any of `subject`/`actor`/`did` and proposed `member`, which reads cleanly alongside `project`.

Wrote the design summary, got "yes", then wrote `/lexicons/com/lognship/project/member.json` with `mkdir -p lexicons/com/lognship/project` first (the `project.json` file and `project/` subdirectory coexist fine on disk; goat walks the tree). Ran `make lint` -- 🟢 across all seven lexicons.

### Why

Owner-side, one-record-per-membership is the smallest design that (a) gives the project owner unilateral control over who counts, (b) is readable from the firehose without diffing, and (c) is consistent with the existing single-record-per-relation pattern of `log.label` and `log.reaction`. Pushing back on the strongRef and at-identifier alternatives kept the design robust to renames and handle changes -- both of which atproto explicitly supports.

### What worked

The "ask one multiple-choice question, save follow-ups for the next turn" rhythm. It surfaced design decisions one at a time and gave Markus room to push back without my pre-baked rationale closing off the option. The brainstorm skill's "describe the design in sections, stay until told to move on" rule didn't get exercised this round because the design was small enough to summarize in one section -- but the section-then-confirm cadence still applied.

### What didn't work

I jumped past Question 4 by reading Markus's clarifying question ("What's the difference between A and D?") as a vote and committing to D. He caught it ("Hey, you decided 4 for me"); I rolled back, restated A vs C cleanly, and let him pick. Lesson: clarifying questions are not votes.

I also offered to "let app views dedupe" as part of the rkey design language, and Markus told me to drop it ("don't mention the app view"). The lexicon is the wire format; app-view behaviour belongs in app-view docs, not in the lexicon's design rationale.

### What I learned

The atproto `at-uri` vs `strongRef` choice maps onto a semantic question: "should this reference survive an edit to the target?" Memberships *should* survive a project rename, so plain at-uri. Reactions and labels *shouldn't* (a reaction is to a specific version of an entry), so strongRef. That framing decides the field type without having to think about it as a JSON-schema problem.

### What was tricky

Naming. I was anchored on "member" as the record type, with `subject` as the field. Markus pointed out the record is the *relation*, not the person, and proposed `membership`. Renaming late was cheap (the lexicon hadn't been published yet), but it's the kind of decision that's much more expensive to flip post-publication, so it was worth the back-and-forth. Renamed file (`project/member.json` → `project/membership.json`), updated the `id` field, updated the description to "Declare a membership of a project.", simplified the `member` field description from "DID of the collaborator." to "DID of the member."

### What warrants review

`/lexicons/com/lognship/project/membership.json` -- 32 lines, one record def, three required fields. The most opinionated decisions (owner-side only, one-record-per-membership, plain at-uri, did not at-identifier, rkey is tid) are all immutable post-publish, so worth one more glance.

### Future work

Optional fields when need surfaces: `role` (for distinguishing trust levels in the app view), `note` (free-form context), `invitedBy` (delegation). All can be added later as optional properties without breaking existing records.

## Step 2: Publish to the network

**Author:** main

### Prompt Context

**Verbatim prompt:** "Perfect. Publish the schema"
**Interpretation:** Run `goat lex publish lexicons/` from the lognship session so the new schema record lands in the lognship.com repo and becomes resolvable to any consumer of the network.
**Inferred intent:** Move from "lints clean locally" to "live and authoritative on atproto", closing the same loop as the other six lexicons.

### What I did

Verified the active goat profile with `bin/goat-as --which` -- already on `lognship`. Before that, ran `goat lex check-dns` and `dig +short TXT _lexicon.project.lognship.com` and was surprised to find a TXT record already in place pointing at the lognship DID. Markus must have set it up when he added the others on 2026-05-04, or there's a wildcard in play. Either way, no DNS work needed.

Ran `goat lex publish lexicons/`. The output was six 🟠 (unchanged) and one 🟢 (`com.lognship.project.membership`). `goat lex status` afterward came back all-🟢, confirming the local copies match the network for all seven lexicons.

### Why

DNS authority + schema record in the publisher's repo is what lets any party resolve `com.lognship.project.membership` by NSID and validate records against the authoritative schema. The `goat lex publish` step is the schema-record half of that contract; the DNS half was already in place.

### What worked

Two commands, no surprises. `goat lex publish` only ships the changed schema -- no mass-resign of the unchanged lexicons -- which is good hygiene, because it means the CIDs of the six existing schema records are preserved.

### What didn't work

Nothing this step.

### What I learned

`goat lex publish` is incremental. The 🟠 markers next to the six unchanged lexicons aren't "skipped" -- they're "verified-already-current and not republished". That distinction is reassuring for anyone worried that a routine publish run might accidentally rewrite history.

### What was tricky

Nothing this step. The friction from 2026-05-04 (figuring out DNS + the publisher contract) had paid off here.

### What warrants review

`goat lex status` should keep returning all-🟢 going forward. Anything else means a local lexicon edit didn't get republished.

### Future work

Same as the previous diary's: when bsky-network PDSes start dynamically resolving lexicons, drop the `-n` flag from record writes. Still not happening.

## Step 3: Add the first two memberships

**Author:** main

### Prompt Context

**Verbatim prompt:** "As the lognship.com account, add @maragubot.com and @markus.maragu.dev as members to the lognship project"
**Interpretation:** Resolve both handles to DIDs, draft two `com.lognship.project.membership` records, validate locally, write them under the lognship session with `-n` (the bsky-network-PDS limitation hasn't gone away).
**Inferred intent:** Exercise the new lexicon end-to-end with real members, so the next thing built (an app view, a CLI, anything) has actual data to read.

### What I did

Resolved both handles in parallel: `goat resolve maragubot.com` returned `did:plc:ccrzw4q3x5cknyo5uqsti74x` (already known from the previous diary), `goat resolve markus.maragu.dev` returned `did:plc:hfjv6h7kg43jey4tqyys7uzp` (new -- markus.maragu.dev's PDS is `phellinus.us-west.host.bsky.network`, distinct from lognship's `stropharia` PDS). Ran `date -u +"%Y-%m-%dT%H:%M:%SZ"` to get a fresh `createdAt` (`2026-05-05T09:09:49Z`) -- per the auto-memory rule, never extrapolate from session anchors.

Drafted both record bodies at `/tmp/lognship-records/membership-maragubot.json` and `/tmp/lognship-records/membership-markus.json`, each pointing at `at://did:plc:ln7l376244zkvhwu4wgr3rzz/com.lognship.project/3mkzclpilfo2v` and including the right DID. Validated locally -- both `valid com.lognship.project.membership record`.

`goat record create -n` for each, in sequence:

- maragubot: `at://did:plc:ln7l376244zkvhwu4wgr3rzz/com.lognship.project.membership/3ml3tp554ah23` (CID `bafyrei...lbai`)
- markus.maragu.dev: `at://did:plc:ln7l376244zkvhwu4wgr3rzz/com.lognship.project.membership/3ml3tp5yqxu2c` (CID `bafyrei...3t5i`)

`goat record list --collection com.lognship.project.membership did:plc:ln7l376244zkvhwu4wgr3rzz` confirmed both records live under the lognship account.

### Why

The brainstorm decided owner-side: only the project owner writes membership records. Both records had to be written from the lognship session (which we were already in). The DIDs anchor the membership to identities, not handles, so a future handle migration on either account is a no-op for membership.

### What worked

The `goat lex validate -c lexicons` → `goat record create -n` → `goat record list` round trip was the same tight loop as the previous diary, applied to a new lexicon. The local validation caught nothing this time because the records were trivially correct, but it's still the right guard before any wire write.

I parallelized the two `goat resolve` calls and `date -u` in a single Bash call to save a round trip. Both DIDs and the timestamp came back in one shot.

### What didn't work

First `goat record list` attempt used positional args (`goat record list -c com.lognship.project.membership lognship.com`) and `-c` -- both wrong. The `-c` short-form is `--collections` (plural, "list collections, not records"), not `--collection`. The handle wasn't recognized either, which is odd because `goat resolve` had just returned for it. Second attempt: `--collection com.lognship.project.membership did:plc:ln7l376244zkvhwu4wgr3rzz` -- worked. Lesson: `goat record list` takes the long-form `--collection` and prefers DIDs over handles.

### What I learned

`goat record list -c` and `goat record list --collection` mean different things -- `-c` is `--collections` (the plural plural-listing flag), `--collection` is the singular filter. Two flags one letter apart in semantics. Worth remembering before the next time I reach for `-c`.

### What was tricky

Nothing material. The two records were straightforward and the tooling did what it said it would.

### What warrants review

The two membership records' AT URIs and CIDs are pinned in this diary. If `goat record get` ever returns a different CID for either, something replaced them. Most likely benign cause is a `record update`; otherwise worth a closer look.

### Future work

A third membership (lognship.com itself, as a self-reference) might be useful for app-view ergonomics -- "show me everyone who's a member, including the owner" without special-casing ownership. Out of scope for now; we'd want to hear from the app view first about what's annoying.

## Step 4: Update the project's display copy

**Author:** main

### Prompt Context

**Verbatim prompt:** `As @lognship.com , change the project description to "AI-first software development in the Atmosphere! Agent-friendly, you own your data." and displayName to "log'n'ship"`
**Interpretation:** Read the existing project record, replace `description` and `displayName` with the new values, leave `createdAt` and `website` alone, write back via `goat record update`.
**Interred intent:** Refresh the user-visible copy of the project now that there are actual collaborators on it -- the original "An AI-first GitHub alternative living in the Atmosphere." was a placeholder, this is the real pitch.

### What I did

Read the current project record at `at://did:plc:ln7l376244zkvhwu4wgr3rzz/com.lognship.project/3mkzclpilfo2v`. Wrote `/tmp/lognship-records/project-updated.json` with:

- `displayName`: `log'n'ship` (was `lognship`)
- `description`: `AI-first software development in the Atmosphere! Agent-friendly, you own your data.` (was `An AI-first GitHub alternative living in the Atmosphere.`)
- `createdAt`: `2026-05-04T08:52:00Z` (preserved verbatim)
- `website`: `https://www.lognship.com` (preserved verbatim)
- `$type`: `com.lognship.project`

Validated locally (`valid com.lognship.project record`), then `goat record update -n -r 3mkzclpilfo2v /tmp/lognship-records/project-updated.json`. New CID `bafyrei...igmi`. Read it back with `goat record get` to confirm the field values landed exactly as written.

### Why

`goat record update` rewrites a record in place at the same AT URI -- no rkey change, no consumer migration. Preserving `createdAt` from the original keeps the project's "originally created" timestamp meaningful (it's documentation of when the project was *first* declared, not when its description was last touched). The new copy reflects the project's actual positioning ("AI-first software development", "agent-friendly", "you own your data") rather than the early-stage GitHub-comparison placeholder.

### What worked

`goat record update` is a single command that took the rkey via `-r` and the JSON file as a positional. Same `-n` rule as creates -- the bsky-network PDS still doesn't dynamically resolve `com.lognship.*` lexicons. Read-back confirmed the record serialized exactly as expected.

### What didn't work

Nothing this step.

### What I learned

`goat record update` doesn't merge -- it replaces. So I needed to include `createdAt` and `website` explicitly in the new JSON, otherwise they'd have been dropped on the floor. The `$type` field is also required (PDS-side), so don't forget it. The "read existing record, modify, write back" pattern is the safe form.

### What was tricky

The old description's CID is now stale -- any consumer who pinned a strongRef to the previous version (none currently exist; this is preventative thinking) would now have a dangling reference. This is exactly why memberships use plain at-uri rather than strongRef: a project edit shouldn't invalidate every membership record. Confirmed by inspection that no membership or repo or log entry currently strongRefs the project.

### What warrants review

The new copy itself -- the description is one sentence, two clauses, exclamation included. If the marketing voice should be more sober, that's a copy edit. The new displayName uses an apostrophe (`log'n'ship`) -- worth confirming that nothing downstream chokes on the apostrophe in URLs or grapheme counts (the lexicon limit is 64 graphemes / 640 bytes; we're well under both).

### Future work

A `goat lex status`-style tool for *records* (not just lexicons) would be useful: "list all records in my repo that strongRef this AT URI, so I know who to notify before I update". Not building it today; just noting that the absence of it is what makes record updates feel slightly more dangerous than creates.
