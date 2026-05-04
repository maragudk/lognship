# Diary: Add `com.lognship.repo` lexicon

Adding a new record lexicon that links one or more git repositories to a `com.lognship.project`. This is the first lognship lexicon to use a strong reference to another lognship record.

## Step 1: Create `com.lognship.repo` lexicon

**Author:** main

### Prompt Context

**Verbatim prompt:** "I'd like to add a com.lognship.repo lexicon next. See the diary for recent work. The lexicon links git repos to a project. It has: - strong-ref project - URL (uri?) - createdAt"
**Interpretation:** Create a record lexicon at `com.lognship.repo` with three required fields: a strong reference to a `com.lognship.project`, a repo URL (using the `uri` string format), and a `createdAt` datetime. The user flagged `uri` with a question mark, asking for confirmation that it's the right format.
**Inferred intent:** Establish the project-to-repo linkage so an actor can declare which git repositories belong to one of their lognship projects, paving the way for repo-scoped log entries.

### What I did

Read the prior diary at `/docs/diary/2026-05-01-lexicons-project-log.md` to pick up the conventions established for `project` and `log.entry`: 3-segment NSIDs directly under `com.lognship`, `key: "tid"` for multi-record collections, descriptions on every field, `createdAt` required and described.

Wrote `/lexicons/com/lognship/repo.json` with NSID `com.lognship.repo`, `key: "tid"`, and three required fields:

- `project` -- `ref` to `com.atproto.repo.strongRef` (the canonical uri+cid pair).
- `url` -- string with `format: "uri"`, which accepts `https://`, `git://`, `ssh://`, etc.
- `createdAt` -- datetime with description matching the project/profile pattern.

`make lint` passed: the new file is 🟢, and the only remaining error is the pre-existing `[unlimited-string]` warning on `log/entry.json` that the user explicitly chose to live with in Step 4 of the previous diary.

Confirmed two design questions with the user before committing:
1. NSID depth: `com.lognship.repo` (3 segments) stays parallel with `com.lognship.project`. If a `com.lognship.repo.*` namespace is later wanted (e.g. `repo.commit`, `repo.branch`), this is the cheap moment to rename. User: "we're good".
2. `url` strictness: `format: "uri"` accepts anything URI-shaped. A regex or a `cloneUrl`/`webUrl` split would be over-engineering for now. User: "we're good".

Committed with no issue references (user confirmed none to ref).

### Why

A strong reference to the project record is the right primitive: it pins both the AT URI and the CID, so a repo record can't silently drift to point at a different version of the project. `format: "uri"` matches the precedent set by `website` on `project.json`. `key: "tid"` matches `project` and `log.entry` -- one actor, many repos.

### What worked

The prior diary at `/docs/diary/2026-05-01-lexicons-project-log.md` was the right thing to read first. It established the `key: "tid"`, "every field gets a description", and required-`createdAt`-with-description patterns, so the new lexicon could mirror them without re-deriving the conventions. `make lint` returned the answer in one shot.

### What didn't work

Nothing.

### What I learned

`com.atproto.repo.strongRef` is the canonical way to link records across lexicons in atproto. The lexicon side is just `{ "type": "ref", "ref": "com.atproto.repo.strongRef" }` -- the indexer is responsible for resolving the uri+cid pair on read. This is the first place lognship's lexicons cross-reference each other, so this pattern will recur (e.g. when log entries gain a repo or project reference).

### What was tricky

None.

### What warrants review

Whether `url` should be the only repo address. A repo typically has both a clone URL (`git@host:owner/repo.git` or `https://host/owner/repo.git`) and a web URL (`https://host/owner/repo`). Currently we collapse them to one field. If the indexer needs to render both -- "view on web" and "clone command" -- this becomes a breaking change later. User accepted the tradeoff for now.

### Future work

- Decide whether `log.entry` should optionally carry a strong-ref to a `com.lognship.repo` (or `com.lognship.project`) so log streams can be repo-scoped. The strongRef pattern from this step is reusable.
- If `com.lognship.repo.*` ends up needing siblings (commit, branch, tag), revisit the NSID and rename before publication.

## Step 2: Switch cross-record references from strongRef to bare AT URI; add `project` to `log.entry`

**Author:** main

### Prompt Context

**Verbatim prompts:**
- "Should the log entry maybe also have a strong-ref to a project?"
- "Are there any alternatives to strong-ref?"
- "Yes." (to switching `log.entry` to bare AT URI)
- "I don't want the repo link in the log entry"

**Interpretation:** Reconsider the linkage primitive. Bluesky's pattern is strongRef for "pin to a version" (likes, reposts, replies, quote posts) and bare reference for "follow updates" (follows). For lognship, project metadata will evolve (rename, description edits) and entries should re-render with the live state, so bare AT URI is the right call -- on both `log.entry` and `repo`. Add a required `project` AT URI to `log.entry`; no repo link on entries.
**Inferred intent:** Avoid the foot-gun of pinning entries to a stale CID of the project record, and keep the log.entry shape minimal -- entries are project-scoped, not repo-scoped, at this layer.

### What I did

On `/lexicons/com/lognship/log/entry.json`:
- Added required `project` field with `format: "at-uri"` and a description.
- Promoted `required` from `["text", "createdAt"]` to `["project", "text", "createdAt"]`.

On `/lexicons/com/lognship/repo.json`:
- Replaced the `project` strongRef (`{ "type": "ref", "ref": "com.atproto.repo.strongRef" }`) with a bare AT URI (`{ "type": "string", "format": "at-uri" }`), keeping the field name and required status.

Briefly added an optional `repo` AT URI to `log.entry` as part of the same edit, then removed it on user instruction ("I don't want the repo link in the log entry"). The `repo` -> `project` linkage in `repo.json` is enough to derive a project's repos; logs are kept project-scoped only.

`make lint` after each edit: all four lexicons parse and lint clean. Only the pre-existing `[unlimited-string]` warning on `log/entry.json` remains.

### Why

- **Bare AT URI over strongRef**: lognship is an issue-tracker shape, not a snapshot-attachment shape. If a project gets renamed, the log entry should re-render with the new name -- strongRef would freeze the project at its CID at entry time, which is the wrong default. Bluesky reserves strongRef for content-anchored attachments (likes, reposts, replies, quote posts); follows use a bare DID. Our case maps to the follow side of that split.
- **Required `project` on `log.entry`**: without a parent reference, a log entry is just an undirected note. Indexers can't query "entries for project X" without it.
- **No `repo` link on `log.entry`**: keeps the shape minimal. A log entry is project-level for now; if we later want repo-scoped logs, that's an additive optional field (still a non-breaking edit if we add it before publication).

### What worked

The Bluesky precedent (`like.subject`/`repost.subject`/reply pins are strongRef; `follow.subject` is a bare DID) gave a clean rule of thumb to apply: pin when content-anchored, bare when you want updates to flow through. That rule decided both edits in one shot.

### What didn't work

I jumped ahead and added an optional `repo` field to `log.entry` along with the `project` field, on the assumption that "project + optional repo" matched the design conversation. The user interrupted and removed it. Lesson: even with auto-mode, when a design conversation ends with a confirmed scope ("yes, switch to bare AT URI"), don't expand scope mid-edit.

### What I learned

The strongRef-vs-bare-AT-URI choice is more load-bearing than it looks. It's not about precision -- it's about whether the reference is a snapshot or a live pointer. Once published, the choice is locked (bare AT URI -> strongRef would be a breaking type change, and vice versa). Worth getting right before any of these lexicons see the firehose.

### What was tricky

Recognising that the "is this strong enough?" framing was actually a "do you want updates to flow through?" question. The first instinct is to reach for strongRef because it sounds safer; the actual question is about update semantics, not integrity.

### What warrants review

Whether `log.entry.project` should stay required. Right now it forces every log entry to be project-scoped, which is what we want today, but it forecloses the option of standalone or actor-level log entries. If that ever matters, this becomes a breaking change. Currently safe: nothing is published yet.

### Future work

- If repo-scoped logs become a real ask, add an optional `repo` (`at-uri`) to `log.entry`. Optional addition is non-breaking even after publication.
- Consider whether `repo.project` should also be required to be the AT URI of a `com.lognship.project` specifically. The lexicon can't enforce the target's NSID -- the `at-uri` format only checks shape -- so this is an indexer-side validation either way.
