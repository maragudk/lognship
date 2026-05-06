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

Drafted the worked example from the brainstorm (an entry with a `#mention` of `@markus.maragu.dev` resolved to `did:plc:hfjv6h7kg43jey4tqyys7uzp`, plus an `#atLink` to another lognship entry) at `/tmp/lognship-records/log-entry-with-facets-example.json`. Hand-computed byte offsets — 15/33 for the mention, 44/118 for the AT URI. Ran `goat lex validate` and got `valid com.lognship.log.entry record`. *(See Step 4 for why this initial validation was misleading: without `--catalog lexicons`, the validator couldn't load the unpublished facet lexicon and silently passed the entire `facets` array unchecked. Step 4 re-runs validation with `-c lexicons` and confirms the design actually holds.)*

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

## Step 4: Code review, triage, and tighten

**Author:** main

### Prompt Context

**Verbatim prompt:** "Use your code review skill" / "/fabrik:address-code-review"
**Interpretation:** Run two competing reviewer subagents over the branch, then walk through their findings one at a time per the address-code-review protocol — triage all comments before any code changes.
**Inferred intent:** Catch what the implementation diary missed before publishing. Lexicons are write-once-evolve-only, so getting the wire format right pre-publication is much cheaper than after.

### What I did

Dispatched two parallel reviewers (`general-purpose` agent type) with a self-contained brief covering the design context, settled decisions to skip re-litigating, and the seven file/area review prompts. They returned independent reports; cross-referenced findings via the code-review skill's "consensus + serious single-reviewer" rule, surfacing six items.

Triaged each comment one at a time per the address-code-review skill — strict no-edits-during-triage rule. Decisions:

- **Skip:** documenting `byteEnd >= byteStart` invariant in the `byteSlice` description; capping the `facets` array with `maxLength`. Both can be added later if friction appears; not load-bearing now.
- **Apply:** add `minLength: 1` to the `features` array (matches the `main` description's "one or more semantic features" and rejects meaningless empty-features facets); move the `facets` field to sit between `text` and `createdAt` in `com.lognship.log.entry` to match the project-wide convention of keeping `createdAt` last (verified across `label`, `reaction`, `project`, `actor/profile`, `project/membership`).
- **Verify, then update:** the diary's "design holds end-to-end" claim, plus a contradicting reviewer claim about whether `format: "at-uri"` actually rejects malformed AT URIs.

For the verification: built four probe records under `/tmp/lognship-records/strictness-test/` (garbage with spaces, `https://` shoved into `#atLink`, valid DID with malformed NSID, well-formed). Ran each through `goat lex validate` twice — without `--catalog` and with `-c lexicons`. Without the catalog, *all four passed* (validator silently skipped the unloadable union). With `-c lexicons`, the three malformed cases failed with precise errors (`AT-URI syntax didn't validate via regex`, `AT-URI first path segment not an NSID: badnsid`); only the well-formed AT URI passed. Reviewer B's reading was correct; the strict-validation justification for splitting `#link`/`#atLink` is real, not aspirational.

After applying both edits and verifying empty-features now fails (`array length out of bounds: 0`), re-ran `goat lex parse`/`lint` over all eight lexicons (clean) and re-validated the worked example with `-c lexicons` (still valid).

### Why

The reviewers' job was to find what the implementation diary missed. Two competing reviewers + the consensus rule keeps the signal high without a single critic dominating. The `--catalog` finding was particularly load-bearing: the original Step 3's "design holds end-to-end" was an honest mistake (the validator failed silent rather than loud), and discovering it pre-publish means we never publish a lexicon whose strictness claims are unverified.

### What worked

The address-code-review skill's no-edits-during-triage rule prevented eager partial fixes that would have muddled the discussion. Walking through six items one at a time, with my recommendation paired against each one, kept Markus in the driver's seat — he skipped two I'd been inclined to apply and confirmed the other four cleanly. Cheap, fast, no surprises.

The empirical reconciliation of the reviewer disagreement (probe records + before/after `--catalog` runs) settled the dispute in five minutes and produced concrete evidence either way. Better than re-reading the spec.

### What didn't work

The original Step 3 ran `goat lex validate` without `--catalog` and called the result "design holds end-to-end". That phrasing was wrong — the validator silently fell back to schema-less validation of the unpublished `com.lognship.richtext.facet` union, so reversed byte ranges, malformed AT URIs, invalid DIDs, and empty `features` arrays all returned `valid` cleanly. Recovering required the code review surfacing the gap; otherwise we'd have published lexicons whose strictness was untested.

### What I learned

`goat lex validate` without `--catalog` is misleadingly permissive when local-only lexicons are involved. It only checks against schemas it can resolve (network-fetchable, published) — the moment a record's `$type` (or any nested union) references an unpublished local lexicon, the validator silently skips that part rather than erroring. **Always pass `-c lexicons` for local lexicon work; treat anything else as smoke-testing only.**

Format-level validation in `format: "at-uri"` is real — given the catalog, it enforces both AT URI grammar and NSID grammar on the path. The Bluesky-style strict format strings carry their weight.

### What was tricky

The reviewers disagreed on whether `format: "at-uri"` rejects malformed input even with `--catalog`. Reading their reports, both seemed plausible. Resolution required running the experiment myself — neither secondhand account would have been authoritative.

Triage discipline. The temptation to "just fix the obvious one quickly" while discussing the next is real. The skill's strict separation of triage and implementation is annoying in the moment and clearly correct in retrospect — partial fixes during discussion would have left half-edited files when Markus then said "actually skip that one".

### What warrants review

The diary's Step 3 inline pointer to Step 4 (the parenthetical caveat) — keeps the historical narrative honest without rewriting the original entry. Acceptable but worth confirming the form is right.

The `minLength: 1` constraint on `features`. With it, an unfortunate future feature-type addition that's the *only* feature on a facet would still validate fine — unions allow any single recognised type. The constraint catches "no features at all", which is the only meaningless case.

### Future work

Carry the `--catalog lexicons` habit into the publication runbook: before `goat lex publish`, validate any draft records against `-c lexicons` rather than the unfiltered `goat lex validate`. Add to the project's eventual contributor docs once they exist.

## Step 5: Can we enforce `byteEnd > byteStart` in the schema?

**Author:** main

### Prompt Context

**Verbatim prompt:** "Can we enforce byteEnd > byteStart in the schema?"
**Interpretation:** Re-open Reviewer A's original concern (the one we skipped during triage) — is there a lexicon mechanism we missed?
**Inferred intent:** Get a clear answer on the limit of atproto's lexicon expressiveness so future "can we enforce X" questions don't keep re-litigating this.

### What I did

Surveyed atproto's lexicon expressiveness. Single-field constraints supported: `minimum`/`maximum` on integers, `minLength`/`maxLength` on arrays, `enum`, `const`, `default`, plus the type-specific `format` strings (`did`, `nsid`, `uri`, `at-uri`, etc.). What's *not* supported: any cross-field relationship — no `dependentSchemas`, no `if/then/else`, no expressions over multiple fields. atproto's data model is a deliberately small subset of JSON Schema designed to keep validation portable.

Surfaced two workarounds: (A) reshape `byteSlice` to `{byteStart, byteLength}` with `byteLength: minimum: 1`, which converts the cross-field invariant into a single-field constraint; (B) document the invariant in `byteSlice.description` — the option Reviewer A originally raised, soft-contract only. Markus took option C: "no action".

### Why

The code review had surfaced the gap; the question deserved a definitive answer rather than a half-remembered "I think it can't be done". Now if it comes up again, the answer is recorded.

### What worked

Direct answer with both workarounds plus my recommendation (no action) let the question close in two turns. Option A is the only schema-side fix and trades away Bluesky-shape parity, which we deliberately preserved during the brainstorm — not a swap I wanted to propose seriously.

### What didn't work

Nothing — small research question, clean answer.

### What I learned

atproto lexicons enforce only single-field constraints. Anything that depends on one field's value relative to another (consistency, derived equality, ordering) lives in client-and-renderer convention, not the schema. Worth noting as a recurring pattern: every time the question "can we enforce X" comes up, if X involves two or more fields, the answer is no by construction. The escape hatch is always to reshape the data so the invariant becomes a single-field property — sometimes worth it, often not.

### What was tricky

Resisting the temptation to propose option A as a clever fix. Bluesky-shape parity was load-bearing during the brainstorm; it shouldn't get sacrificed for one invariant after the fact.

### What warrants review

Nothing changed in code; the diary entry is the only artefact.

### Future work

If a future feature type wants more invariants the schema can't express, the answer will be the same — schema gives single-field bounds, relationships live in client/renderer convention. Worth recording in contributor docs once those exist.
