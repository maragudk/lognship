# Diary: Articulate the project's principles in docs and on the front page

Markus wanted the core values / design choices of lognship written down — both as a doc in `docs/` and as a section on the home page. The README and home page had been carrying a marketing-style tagline ("AI-first software development in the Atmosphere! Where you own your data."), but the actual principles behind the project had never been spelled out anywhere.

## Step 1: Brainstorm, draft, and ship the doc + home-page section

**Author:** main

### Prompt Context

**Verbatim prompt:** "I'd like to figure out the wording of the core values / design choices for this project, so we can put them in a document in docs/ as well as on the front page."

**Interpretation:** Run a brainstorm interview to extract Markus's actual principles, then write them up as a layered values doc (principle → design choices) and produce a front-page distillation.

**Inferred intent:** The project has been shipping for a couple of weeks (lexicons, basic app view, Tap replication). The conceptual identity has been implicit. Markus wants it explicit so the rest of the work — and the public-facing surface — can refer back to a stable articulation.

### What I did

Ran the `fabrik:brainstorm` skill. Asked five multiple-choice questions to converge on the shape:

1. Framing of the doc → "principles + design choices, layered".
2. Whether the existing README/home language captured the principle set → Markus chose "let me describe in my own words first" and gave a backstory.
3. Where "AI-first / agent-friendly" sits → Markus called it a top-level principle and sharpened it to *"the orchestration / knowledge / context layer between agents and humans"*.
4. How the data-ownership / open-write / open-source claims relate → collapse into one principle: "Your data, on the open network".
5. Front-page layout → keep the existing tagline, add a "Principles" section below it; keep the elaborated bullets; don't link to the doc (internal for now).

Then drafted in sections per the brainstorm skill, presenting 200-300 words at a time and waiting for "yes" before advancing: outline & ordering → thesis + Principle 1 → Principle 2 → Principle 3 → Principle 4 → front-page distillation.

The four principles, in causal order:

1. The product, not the code, is the artifact.
2. The log is the unit.
3. A layer between agents and humans.
4. Your data, on the open network.

Wrote `/docs/principles.md` with thesis paragraph, four principles (each: prose explanation + "What this means in practice" bullets).

Edited `/html/home.go` to add a Principles section below the existing flex hero block — wrapped both in a `space-y-16` div, used `<strong>` lead-in via `Raw` to match the existing tagline's style of inline emphasis.

Verified with `go build ./...`. App was not running locally (no `.env`, no `app.log`), so I didn't visually confirm in browser.

### Why

Markus runs lognship as a one-person project right now, but the principles document doubles as load-bearing context for agents. Once `docs/principles.md` exists, future agents (and Markus himself, on cold cache) have a single canonical place to ground design decisions.

The layered structure (principle → design choices) was Markus's explicit pick — it makes the doc useful as both a manifesto and an architecture rationale.

The causal ordering (artifact → log → agent layer → open network) flips the README's lead. The README leads with the agent angle, but agents-as-collaborators is *downstream* of "the product, not the code, is the artifact". Putting the artifact-shift first makes the agent claim land harder when it arrives.

### What worked

The multiple-choice brainstorm format converged fast — five questions and we had the shape. Markus's answer to Q3 ("orchestration / knowledge / context layer between agents and humans") was the sharpest single sentence in the entire conversation; that phrase carried straight into Principle 3's prose.

Causal ordering instead of "marketing-importance ordering" (i.e., not leading with the AI-first headline) made the doc cohere as an argument rather than a list of features.

Section-by-section drafting with explicit "yes to advance" gates kept each section short enough that course corrections were cheap.

### What didn't work

Initial copyediting pass was sloppy — see Step 2.

I used the term "the product's meta" in Principle 2's opening sentence as a callback to Principle 1, but I had never actually introduced "meta" as a term in Principle 1 (I'd used "reasoning" and "context" interchangeably). Caught and fixed in Step 2.

### What I learned

When the user gives backstory in their own words, the *exact phrases they use* are the most reusable raw material. "Code is cheap now", "actionable insight from real users or developers is worth much more", "an append-only log of relevant information", "lognship is just one possible app view on top of the data in the network", "the app view (and the project owners) decides whether to surface it" — every one of those landed in the doc nearly verbatim.

The home page's existing styling pattern (`Raw` with inline `<em>`/`<strong>`) is the right pattern for the principles bullets too, since each bullet is one sentence with mid-sentence emphasis. Trying to use separate `Strong()` + `Text()` nodes would have produced spacing weirdness around the period after the headline phrase.

### What was tricky

The agent-friendly / AI-first framing. It's the marketing headline (`HomePageProps.Description` says "Agent-first software development in the Atmosphere! Where you own your data."), but it isn't *the* principle — it's downstream of the artifact-shift. Asking Q3 explicitly to position it surfaced this. Markus's answer ("top-level, AND it's the orchestration / knowledge / context layer") gave it a sharper home than the existing marketing language.

Picking the doc filename. `docs/spec.md`, `docs/decisions.md`, and `docs/values.md` were all candidates. `principles.md` won because (a) "principles" was Markus's framing, (b) the fabrik `spec` skill targets `docs/spec.md` for "what the product is and why it exists" which is a different doc, and (c) the fabrik `decisions` skill targets `docs/decisions.md` for ADR-style records, also different. `principles.md` is its own thing and won't collide with either of those if they get added later.

### What warrants review

- `/docs/principles.md` — full doc. Read end-to-end; check that the four principles read as four distinct claims and not three claims plus one rephrasing.
- `/html/home.go:16-37` — the Principles section. Build was clean but I didn't visually verify; eyeball at `make watch` to make sure the `space-y-16` outer wrapper doesn't crush the layout on mobile, and that `list-none pl-0` on the `<ul>` actually suppresses the default bullet styling under the project's tailwind config.
- The `Raw("<strong>...</strong> ...")` pattern in the bullets — confirm it renders with the right spacing after the bold lead-in and that the period sits flush against the closing `</strong>`.

### Future work

- The home-page Principles section currently doesn't link to `docs/principles.md` (Markus said "internal for now"). When the doc goes public, add the link.
- The README still has the old marketing-style tagline. At some point the README should pick up the same principle list — but not yet; Markus didn't ask.
- The four principles imply some concrete UI decisions that aren't yet made: e.g., what exactly does a project's home page show instead of a file tree? That's a downstream design question this doc deliberately doesn't answer.

## Step 2: Strunk copyedit pass

**Author:** main

### Prompt Context

**Verbatim prompt:** "Did you use the writing skill?"

**Interpretation:** Markus had asked me earlier ("BTW, make sure to use your writing-clearly-and-concisely skill before writing anything") and was checking whether I'd actually done it. I had loaded the skill's overview but hadn't read the full `elements-of-style.md`, and I'd been applying the rules from memory. Markus deserved an honest answer plus a real pass.

**Inferred intent:** Quality gate. The doc is going on the home page and is meant to be the project's articulated identity — it should be properly edited, not just plausibly edited.

### What I did

Came clean: said I'd loaded the skill overview but not read the full guide, and offered two options — dispatch a subagent for a copyedit pass, or read the full guide myself and edit in place. Markus picked option 2.

Read `/Users/maragubot/Developer/fabrik/skills/writing-clearly-and-concisely/elements-of-style.md` end-to-end.

Walked the doc and home-page bullets against the rules. Made nine targeted edits:

1. `/docs/principles.md:3` — "Code is cheap to produce; agents write it." → "Code is cheap; agents write it." (Rule 13: omit "to produce".)
2. `/docs/principles.md:3` — "lognship is built on that shift." → "lognship answers that shift." (Rule 10: active voice.)
3. `/docs/principles.md:7` — "the tasks open" → "the tasks taken on". (Rule 15: parallel construction. The other three items in the list — "decisions made", "tradeoffs accepted", "bugs found" — are noun + past-participle; "tasks open" was noun + adjective.)
4. `/docs/principles.md:14` — "Anything an agent or a human needs to gather context lives as readable, addressable records, not as comments..." → "What an agent or human needs to gather context lives as readable, addressable records — not as comments...". (Rule 13: "Anything that" → "What"; em-dash for stronger antithesis per Rule 11.)
5. `/docs/principles.md:18` — "If the product's meta is the artifact, the natural shape for it is a log:" → "Once the product is the artifact, the log is the natural form for everything that surrounds it:". (Removed undefined term "meta"; Rule 12: definite, concrete language.)
6. `/docs/principles.md:26` — "Reactions and labels attach to entries, so the log carries its own social signal without leaving the format." → dropped "without leaving the format" (Rule 13: unclear words don't earn their place).
7. `/docs/principles.md:27` — "Append-only means corrections add an entry; they don't overwrite one. The reasoning trail stays intact." → "Corrections add a new entry rather than overwriting an old one. The reasoning trail stays intact." (Rule 11: positive form; "they don't overwrite" was a negation.)
8. `/docs/principles.md:38` — "Every entry is readable by both." → "Both read every entry." (Rule 10: active voice.)
9. `/docs/principles.md:44` — "Everything above only holds if" → "Everything above holds only if". (Rule 16: place modifiers next to what they modify; "only" was modifying "holds", not "everything".)

Then two on the home page:

10. `/html/home.go:31` — "An append-only stream anyone can read for context or act on." → "An append-only stream — read it for context, write to add to it." (The original had a dangling object after "act on"; the rewrite makes the action explicit and parallel.)
11. `/html/home.go:32` — "The shared surface where work coordinates." → "The shared surface for coordination." ("Where work coordinates" was an awkward intransitive use.)

Verified with `go build ./...` after the home-page edits.

### Why

The doc is going on the home page. If the wording is sloppy, every reader notices and Markus is the one who has to live with it. A real copyedit pass — with the full Strunk guide loaded — was the right gate.

Also: Markus had explicitly asked me to use the skill before writing. Not actually doing so, then claiming I had, would have been worse than the original sloppy pass.

### What worked

Reading the full guide surfaced things the overview didn't catch — specifically Rule 16's "only" placement (which is exactly the kind of thing you don't notice without the rule top-of-mind), and the "tasks open" parallel-construction violation (the overview mentions Rule 15 but doesn't show how subtle the violation can get).

The Edit tool was the right shape for this work — small, targeted swaps with full context. No risk of accidentally rewriting something Markus had already approved.

### What didn't work

The first attempt at editing the home page bullets failed with `String to replace not found in file` because I tried to combine two `Li(...)` line edits into a single Edit call with embedded tabs that didn't match the file's actual indentation. Fix: split into two separate single-line Edit calls.

```
String: 				Li(Raw("<strong>The log is the unit.</strong> ..."

<tool_use_error>String to replace not found in file.</tool_use_error>
```

### What I learned

For tab-indented files, prefer single-line Edit calls over multi-line ones. The whitespace-matching is fragile when the editor's serializer of the `old_string` doesn't roundtrip cleanly through tabs.

The "honest admission + offer two paths" response shape was useful — it surfaced the gap, gave Markus a concrete choice, and didn't waste time on apologies. He picked the path he wanted in one word ("2") and the work resumed.

When the user asks "did you use X?", the answer should be "partially, here's exactly what I did and didn't do" — not a defensive yes-or-no.

### What was tricky

Resisting the temptation to edit prose I'd already gotten "yes" on during the brainstorm. Markus had approved each section. The Strunk pass was a *narrow* mandate — only fix things that actually violate the rules — not a license to keep tinkering. I limited myself to nine + two changes that each pointed to a specific rule.

### What warrants review

The eleven specific edits — each one references the Strunk rule that motivated it. If Markus disagrees with any rule application, the change is small and reversible.

In particular: edit #5 (replacing "If the product's meta is the artifact" with "Once the product is the artifact") changes the logical hinge of Principle 2's opener. The original implied "the meta of the product is the artifact (per Principle 1), so..."; the rewrite says "the product is the artifact, so the log is what surrounds it." Slightly different claim. Worth a re-read.

### Future work

None directly. Possible follow-up: ask Markus whether the writing-clearly-and-concisely skill's overview should call out a few rules that are easy to miss without the full guide (Rule 16 "only" placement, Rule 15 parallel construction in lists). That's a `/improve-skill` candidate, not work for this diary.
