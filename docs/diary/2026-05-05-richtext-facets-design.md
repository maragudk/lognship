# Diary: Design rich-text facets for log entries

A `com.lognship.log.entry` carries free-form `text`. Until now, references inside the prose to other log entries, atproto accounts, or external URLs were just text. This task settled how to add structured references — borrowing Bluesky's facet pattern — and shipped the lexicon changes locally (without publishing them).

## Step 1: Settle the high-level pattern

**Author:** main

### Prompt Context

**Verbatim prompt:** "If I'd like to reference a log entry (or several) from another log entry, what would be the best approach? Just in the text? Or make it structured?"
**Interpretation:** Three options to weigh — bare URIs in prose; a top-level structured `references` field on the entry record; a separate side-record (like `com.lognship.log.label`) whose subject points at the related entry.
**Inferred intent:** Markus has settled what an entry is, and wants to settle how entries link to each other before any reader UI gets built.

### What I did

Surveyed the three options. The deciding question was who declares the link and when: author-at-write-time (field on the entry) vs. anyone-later (side record). Recommended a field on the entry, since for a developer log "this entry builds on entry N-3" is the author's narrative, not a third-party assertion.

Then sent an Explore subagent to read up on Bluesky's pattern. The brief came back compact and clear: `app.bsky.richtext.facet` annotates byte ranges in `text` with typed features (`#link`, `#mention`, `#tag`); clients parse prose, compute UTF-8 byte offsets, and attach facets at write time. The text stays syntax-free (no markdown, no `<a>` tags), and unknown feature types degrade to plain text.

That research pivoted the design. The "structured top-level field" framing was the wrong shape: Bluesky uses *both* inline facets *and* top-level `app.bsky.embed.record` for different intents (mid-prose mention vs. quote-tweet). Lognship's question — "how to reference other entries from prose" — is the inline-facet half. Locked the high-level: text + facets, no top-level references field, no relationship taxonomy.

### Why

Inline references mirror how GitHub renders `#123` cross-references and how Bluesky handles links/mentions: declaratively, with byte ranges pre-computed so renderers don't re-parse prose forever. Meaning lives in the words around the link; the schema just keeps the link mechanically reliable.

### What worked

The subagent investigation pulled together a coherent ~400-word brief in one round-trip — Bluesky's design rationale (syntax-free text, no parsing burden, easy length validation), the byte-range structure, the client responsibilities, and the relationship to embed records. That brief was load-bearing for the rest of the brainstorm.

### What didn't work

Initial framing was a false binary: I posed it as "structured field on the entry vs. side record". Bluesky's actual pattern is both, for different intents. Pivoting after the research cost one exchange and produced a sounder framing.

### What I learned

Bluesky's facet system has a subtlety I'd missed: the `features` array per byte range allows multiple semantic features to overlay (e.g. a span that's both a mention and a link), and the union is open so unknown future feature types degrade gracefully without breaking older readers.

### What was tricky

Markus's mid-step reset — "I don't think they're obvious. Let's throw out assumptions and start from scratch" — was uncomfortable but right. I was inheriting Bluesky's set unexamined. Re-deriving from lognship's needs (mention project members, reference URLs and AT URIs) and *then* discovering the answer is close to Bluesky's anyway made the choice principled rather than imitative.

### What warrants review

The decision *not* to add a top-level "references" field with relationship types (`supersedes`, `continues`, etc.). If lognship later wants a typed graph of entry-to-entry relationships (for project tracking, dependency graphs, etc.), this model would shift. For now, prose-around-facet conveys meaning, which is enough.

### Future work

A git-ref feature type — `{repo: at-uri, sha: string}` — could land in `com.lognship.richtext.facet` later, pointing at `com.lognship.repo` records. Deliberately deferred from v1; the schema dependency (`com.lognship.repo`) is already in place.

## Step 2: Brainstorm the facet types and naming

**Author:** main

### Prompt Context

**Verbatim prompt:** "I'd like to brainstorm the types of facets first"
**Interpretation:** Don't lock the section yet — survey candidate feature types, debate which earn their keep.
**Inferred intent:** Markus wants the lexicon to ship with intentional choices, not a copy of Bluesky's set.

### What I did

Surveyed candidates beyond the obvious `#link`/`#mention`: developer-log-specific options like `#commit`, `#file`, `#pr`, plus a formatting tier (bold/italic/code). Markus narrowed to "@mentions and URLs/AT URIs" — back to two-ish types, but with a real follow-up: are URLs and AT URIs the same kind of thing?

Argued first for a single `#link` covering both URI flavours (app view branches on scheme at render time). Markus pushed back — "isn't `format: at-uri` kind of nice?" — surfacing the strict-validation argument I'd dismissed too quickly. Walked the recommendation back: split into `#link` (loose `format: "uri"`) and `#atLink` (strict `format: "at-uri"`). Markus picked the asymmetric naming `#link` + `#atLink`: keep "link" as the unifying word, qualify the AT one.

Confirmed handle-to-DID normalisation at write time for both `#mention` and `#atLink` (because AT URIs can use handle authorities and the same migration durability concern applies). Confirmed: skip `#tag` (labels do that job better — third-party-applicable, independently created/deleted), skip the formatting tier (Bluesky's restraint is wise; keep facets semantic, not visual).

Also surfaced the precedence rule for the client encoder: every AT URI is technically a valid `#link`, but the client picks the most specific applicable type (`#atLink` for `at://`, `#link` otherwise).

### Why

The strict-validation argument for splitting was load-bearing. Silently broken references in a developer log are worse than the marginal cost of one extra feature type. Catching `at://garbage with spaces/foo` at write time via lexicon enforcement beats finding it months later when an app view tries to render.

### What worked

Markus's "throw out assumptions" reset pulled the design back to first principles. The end set (`#mention` + `#link` + `#atLink`) is identical in structure to Bluesky's minus `#tag`, but reached via lognship-specific reasoning. Same shape, principled choice.

### What didn't work

Initially argued single `#link` is sufficient and only split when Markus proposed it. Walking that back honestly was better than digging in.

### What I learned

The lexicon's `format: "uri"` accepts any RFC 3986 URI including `at://` URIs, but doesn't enforce internal AT URI structure (NSID + rkey grammar). `format: "at-uri"` does. That distinction is the entire reason `#link` and `#atLink` justify being separate types — same field shape, different validators.

Also: `at://...` URIs and `https://...` URLs share generic URI syntax but have completely different internal structure (authority-as-DID-or-handle vs. authority-as-DNS-host, restricted path grammar vs. arbitrary hierarchical path). The shared `format: "uri"` validation is the loose floor; per-scheme parsers do the real work.

### What was tricky

Naming. "RecordRef" felt jargon-heavy ("they're both links" Markus pointed out); "atLink" emerged as the cleanest qualifier of `#link`. Also tricky: explaining the precedence rule. Every `#atLink` is also semantically a valid `#link`, but clients pick the more specific. The rule lives in client convention and lexicon descriptions, not in lexicon enforcement (which would need a custom format validator).

### What warrants review

The `#link` description nudges authors toward `#atLink` for AT URIs but doesn't enforce it. A malformed AT URI smuggled into `#link` will validate at the lexicon level; only the renderer would notice it doesn't parse properly as an AT URI. Acceptable graceful degradation — but worth knowing if "broken atLinks slipped through as #link" ever appears in a debugging session later.

### Future work

When a client lands, it should implement the precedence rule (AT URI → `#atLink`) in its facet builder. Documenting the rule in lexicon descriptions is the only enforcement mechanism short of writing a custom format validator.

## Step 3: Implement the lexicons (no publish)

**Author:** main

### Prompt Context

**Verbatim prompt:** "we're good"
**Interpretation:** Brainstorm complete; ship the design as lexicon files. Markus had explicitly said "Don't publish anything yet" earlier, so the work stops short of `goat lex publish` and DNS setup.
**Inferred intent:** Convert the design into reviewable artefacts so Markus can see the JSON and validate it himself before deciding when to publish.

### What I did

Created `lexicons/com/lognship/richtext/facet.json` with five defs:

- `main` — the facet object: `index` (a `#byteSlice` ref) + `features` (an array of unions over `#mention`, `#link`, `#atLink`).
- `byteSlice` — `{byteStart, byteEnd}`, both non-negative integers (UTF-8 byte offsets, `byteStart` inclusive, `byteEnd` exclusive).
- `mention` — `{did}` with `format: "did"`.
- `link` — `{uri}` with `format: "uri"`. Description nudges authors toward `#atLink` for AT URIs.
- `atLink` — `{uri}` with `format: "at-uri"`. Description notes the authority-normalisation-to-DID convention.

Updated `lexicons/com/lognship/log/entry.json` to add an optional `facets` property of `array of ref to com.lognship.richtext.facet`. No version bump (additive optional change), no changes to required fields (`project`, `text`, `createdAt`), no changes to `text` size limits.

Validated everything in two passes: `find lexicons -name '*.json' -print0 | xargs -0 goat lex parse` and the same with `goat lex lint`. Both came back clean across all eight lexicon files.

Drafted the worked example from the brainstorm (an entry with a `#mention` of `@markus.maragu.dev` resolved to `did:plc:hfjv6h7kg43jey4tqyys7uzp`, plus an `#atLink` to another lognship entry) at `/tmp/lognship-records/log-entry-with-facets-example.json`. Hand-computed byte offsets — 15/33 for the mention, 44/118 for the AT URI. Ran `goat lex validate` and got `valid com.lognship.log.entry record`. Design holds end-to-end.

### Why

Validation against a sample record proves the lexicon JSON, the facet union, and the entry's array-of-ref typing wire together correctly. Without that test, the lexicon could be silently wrong in ways that surface only at first publish.

### What worked

Mirroring `app.bsky.richtext.facet`'s structure made the JSON straightforward and the standard atproto formats (`did`, `uri`, `at-uri`) did the validation heavy lifting. `goat lex parse`, `lint`, and `validate` all came back clean on the first try.

### What didn't work

My first instinct was `goat lex validate-record /path/to/file` — there's no such subcommand, the correct one is `goat lex validate`. Quick recovery via `goat lex --help`.

### What I learned

`goat lex validate` validates a record against the local lexicons in `lexicons/`. It's the right tool for "does this draft record satisfy its declared `$type`?" — useful before any `record create` write, since validation errors at the PDS are less informative than locally.

### What was tricky

Computing byte offsets by hand. The example was all ASCII, so byte == char, but the process is still error-prone — a real client must not delegate byte-offset math to humans.

### What warrants review

The two new lexicon files. The `description` strings encode design conventions: the precedence rule, the authority-normalisation convention, why facets stay separate from text. If Markus wants those to read differently for clarity or tone, this is the moment to adjust before publication.

The `_lexicon.richtext.lognship.com` DNS TXT record doesn't exist yet — it's needed before `goat lex publish` can publish the new schema. Out of scope for this task per Markus's "don't publish" instruction.

### Future work

When ready: add the DNS TXT record at `_lexicon.richtext.lognship.com` pointing at the lognship DID, run `goat lex publish lexicons/`, confirm with `goat lex status`. After publication, the next chunk is a client that builds facets programmatically (parses prose, resolves handles, computes byte offsets per the section 4 spec from the brainstorm) and an app view that consumes them per the section 5 spec.
