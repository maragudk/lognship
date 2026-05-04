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
