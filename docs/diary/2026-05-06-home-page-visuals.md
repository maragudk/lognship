# Diary: Home page visuals — logo placement, animated wave background, polish

A run of visual polish on the lognship landing page. Started with placing the new logo on the home page, snowballed into an animated SVG wave background inspired by `/banner.png`, then a series of tuning passes (transparency, layer count, motion model, timing function, logo bobbing, logo PNG transparency). All changes scoped to `/html/common.go`, `/html/home.go`, `/tailwind.css`, and the logo asset at `/public/images/logo.png`. No new dependencies.

## Step 1: Logo on the right side of the home page

**Author:** main

### Prompt Context

**Verbatim prompt:** "Add the logo in the right side on @html/home.go" — followed up with "Add it to public/", "public/images/", and "Center it horizontally when it moves above the headline".
**Interpretation:** Drop `/logo.png` into `/public/images/`, render it on the home page to the right of the existing headline/copy on `sm` and up, and stack it above the headline (centered) on mobile.
**Inferred intent:** Give the landing page some visual identity beyond the wordmark, while keeping the existing copy and links readable on phones.

### What I did

Copied `/logo.png` to `/public/images/logo.png`. The static handler at `maragu.dev/glue/http/static.go` already routes `/{images|scripts|styles}/*` to `public/`, so no routing changes needed. Wrapped the home content in `flex flex-col-reverse sm:flex-row sm:items-center gap-8`: text block (`grow max-w-2xl`) on the left at `sm+`, logo (`w-48 sm:w-64 shrink-0`) on the right. `flex-col-reverse` puts the logo above the headline on mobile, and `mx-auto sm:mx-0` centers it there. First pass used `getHashedPath("public/images/logo.png")` for cache-busting; Markus said "Drop the getHashedPath stuff, it's ok" — switched to a plain `/images/logo.png`.

### Why

`flex-col-reverse` is the simplest way to swap the visual stacking on mobile while keeping a single source of children. Hashed paths matter for CSS/JS where stale caches break the page; for a hero image, plain path is fine and reads cleaner.

### What worked

The static handler's `/{images|scripts|styles}/*` route covered the new directory without any glue changes. `flex-col-reverse` + `mx-auto sm:mx-0` did mobile centering in two utility classes.

### What didn't work

Nothing broke.

### What I learned

`maragu.dev/httph.VersionedAssets` is the middleware that strips the content hash from URLs like `/styles/app.<hash>.css` before the file server sees them — that's why `getHashedPath` works without any files actually being renamed on disk. Worth knowing for any future asset wiring.

### What was tricky

Briefly worried about whether nesting an `absolute` sibling under `relative` would paint under or over a `static` sibling — that came up in the next step, not here. For the logo step itself, the layout was straight CSS.

### What warrants review

Reviewer should resize the browser across the `sm` breakpoint and confirm the logo flips between right-of-text and centered-above-text without layout shift.

### Future work

None directly implied.

## Step 2: Animated wave background inspired by `/banner.png`

**Author:** main

### Prompt Context

**Verbatim prompt:** "Can you make some animated waves in the background, very subtle, inspired by @banner.png ?"
**Interpretation:** Add a decorative animated wave field behind the home page content, in the brand pink, low-key enough not to fight the copy.
**Inferred intent:** Match the hand-drawn pink-wave feel of `/banner.png` so the landing page has some movement and atmosphere.

### What I did

Added `waves()` in `/html/common.go` — three absolute SVG layers stacked vertically inside the main white content area, each rendered via `Raw(...)` (the `gomponents/svg` package was retired, per Markus). Each SVG is `w-[200%]` with a tileable `viewBox="0 0 2400 100"` wave path of 12 quadratic bows (period 200, 6 cycles in each half), so translating `-50%` shows the second half of the path, which is identical in shape to the first. Animations defined as Tailwind v4 `@theme` utilities: `--animate-wave-slow/medium/fast` driving `wave-drift` and `wave-drift-reverse` keyframes. The white grow container in `Page()` got `relative overflow-hidden`, the content got wrapped in a `relative` div so it paints above the absolute waves, and a `prefers-reduced-motion` rule disables the animations.

### Why

CSS `transform: translateX` on an SVG is GPU-friendly and stays cheap regardless of layer count. Doubling the SVG width and translating `-50%` is the standard seamless-loop trick for tileable patterns, no JS required. Putting waves in `Page()` rather than only `HomePage()` is intentional — the only public route right now is home, and if other pages get added later, the brand-feel should carry through.

### What worked

The static-then-positioned painting order: `<div class="absolute ...">waves</div>` followed by `<div class="relative">content</div>` — both children of a `relative` parent — paints content on top of waves without needing explicit `z-index`. Both are positioned, so DOM order wins.

`Raw(...)` for the SVG content kept the gomponents code readable and avoided the retired svg package.

### What didn't work

Nothing in this step. The "jump after a while" issue surfaced later (step 4).

### What I learned

Tailwind v4's `@source` directive scans Go files for class names, so arbitrary values like `w-[200%]` and percentage `top-[15%]` inside `Class(...)` strings get picked up the same as in templates. Means I can keep all the wave geometry in `common.go` without any manual safelisting.

### What was tricky

Confirming the wave path actually tiles. Each `Q` followed by sequential `T`s makes a parabola-then-mirrored-parabola pattern with period 400 (one up-bow + one down-bow), not 200 as I first assumed. With viewBox 2400 that's 6 full periods, and the pattern at viewBox 1200 (the seam) matches viewBox 0 in both position and slope, so the seam is invisible. (Browser sub-pixel rounding still bit us later — see step 4.)

### What warrants review

Reviewer should pull the page up in light and dark mode, confirm waves are visible but not distracting, and check that the home content stays clickable (waves have `pointer-events-none`).

### Future work

None directly implied.

## Step 3: More transparent, more layers

**Author:** main

### Prompt Context

**Verbatim prompt:** "Make them a bit more transparent, and add a few more"
**Interpretation:** Drop the existing layer opacities and bring layer count up from 3 to ~7, with a mix of speeds, colors, and stroke widths so they don't all sync.
**Inferred intent:** The first pass was too present; the desired feel is closer to a watermark than a feature.

### What I did

Refactored `waves()` into a small inner closure `wave(positionClasses, colorClass, opacityClass, animClass, strokeWidth)` so each layer is one line. Bumped layer count from 3 to 7, positions from `top-[5%]` down to `top-[85%]` in roughly even steps. Opacities dropped from 20/25/30 to 10/15. Stroke widths varied (1.5 / 2 / 2.5). Added three more animation variants in `tailwind.css` so directions and speeds don't repeat across adjacent layers: `wave-slow-reverse`, `wave-medium-forward`, `wave-fast-reverse`. The `prefers-reduced-motion` selector got the new classes too.

### Why

A closure beats six near-identical `Div(Class("absolute inset-x-0 ..."))` blocks; the differences (position, color, opacity, animation, stroke) are exactly the closure's parameters. More layers + lower opacity gives a denser-but-quieter pattern, which is more "watermark" than "feature".

### What worked

The closure compressed seven layers into seven readable lines. Mixing forward/reverse and varying durations across adjacent layers prevents the visual lockstep that would have been obvious with three or four layers all running the same animation.

### What didn't work

Nothing broke. The new `animate-wave-*` classes all showed up correctly in both the served HTML and the compiled CSS — verified by grepping `/styles/app.<hash>.css`.

### What I learned

Tailwind v4's default opacity scale supports `opacity-15` (every 5%, not every 10%), so I didn't need arbitrary values for the half-step.

### What was tricky

Nothing tricky here.

### What warrants review

Reviewer should confirm the waves still feel "very subtle" with seven layers — easy to overshoot density when adding count.

### Future work

None directly implied.

## Step 4: Pendulum motion + logo bob

**Author:** main

### Prompt Context

**Verbatim prompt:** "The animation 'jumps' after a while. Should it go back and forth instead?" — followed by "Also, make the logo on the front page move subtly up and down."
**Interpretation:** Switch the wave loop from "translate to -50% then snap back to 0" to "translate to -50% then reverse back to 0" so there's no seam. Add a subtle vertical bob to the logo.
**Inferred intent:** The loop seam was visually catching even though the path tiled mathematically — likely browser sub-pixel rounding accumulating. A pendulum motion sidesteps the seam entirely. The logo bob should hint at the boat-on-water idea without being kinetic.

### What I did

Switched all six wave animation utilities in `/tailwind.css` from `linear infinite` to `ease-in-out infinite alternate` (or `alternate-reverse`), and dropped the now-unused `wave-drift-reverse` keyframes since `alternate-reverse` covers the same effect from the same `wave-drift` keyframes. Added `--animate-bob: bob 6s ease-in-out infinite alternate` with a `bob` keyframe translating `-0.5rem` to `0.5rem`. Tagged `<img>` in `/html/home.go` with `animate-bob`. Updated the `prefers-reduced-motion` rule to include `.animate-bob`.

### Why

`alternate` removes the seam by reversing direction at each iteration boundary instead of looping — there's no instant translate jump, only a velocity reversal. The motion stays continuous in position. For the bob, a 1rem peak-to-peak swing over 6s is small enough to read as "afloat" and slow enough not to compete with the waves.

### What worked

`alternate` plus `alternate-reverse` covered all the variation I had with the two original keyframes, so I could delete `wave-drift-reverse` and consolidate to a single keyframe.

### What didn't work

Initially shipped with `ease-in-out`, which became its own problem (step 6).

### What I learned

`animation-direction: alternate-reverse` is the right tool for "start at the end, oscillate from there" — it lets you stagger phase across multiple elements without manually offsetting each animation's `animation-delay`.

### What was tricky

Decision between linear-with-jolt and ease-in-out-with-slow-ramp. I picked the latter because at slow speeds (22s–55s) instant velocity reversal at the turnaround feels wrong for a fluid motion. Markus pushed back the other way (step 6), which was correct — the slow ramp at the start of each swing was more noticeable than the velocity reversal would have been.

### What warrants review

Reviewer should let the page sit for several minutes and confirm the seam is gone — the previous bug only surfaced after a few cycles.

### Future work

None directly implied.

## Step 5: Make outer-most white in the logo transparent

**Author:** main

### Prompt Context

**Verbatim prompt:** "Can you make the outer-most white in the logo transparent?"
**Interpretation:** Flood-fill the white background outside the boat to alpha 0, while preserving the white inside the sail and the hull (those are part of the logo, not background).
**Inferred intent:** So the logo composites cleanly over the wave background and any future colored backgrounds, without a visible white square around it.

### What I did

Used ImageMagick: `magick logo.png -alpha set -bordercolor white -border 1 -fill none -fuzz 10% -floodfill +0+0 white -shave 1x1 public/images/logo.png`. The `-border 1` adds a guaranteed-white perimeter so flood-fill from `+0+0` always starts in outer white; `-shave 1x1` trims that border off after the fill. Verified by sampling pixels: corner is now `srgba(0,0,0,0)`, sail interior and hull interior remain `srgba(255,255,255,1)`. Composited the result over the brand pink as a sanity check at `/tmp/logo-on-pink.png` and confirmed the topology was correct.

### What worked

`-floodfill +0+0` (the operator) succeeded where my first attempt with `-draw "alpha 0,0 floodfill"` (the draw command) silently did nothing — corner pixel was still `srgba(255,255,255,1)` after that one. Adding the bordered perimeter made the start position robust regardless of any anti-aliased fringe touching the actual edge.

### What didn't work

First attempt:

```
magick logo.png -alpha set -bordercolor none -border 1 \
  -fill none -fuzz 10% -draw "alpha 0,0 floodfill" \
  -shave 1x1 public/images/logo.png
```

Output had no transparency. The `-draw "alpha ... floodfill"` form wasn't doing what I expected on this image. Switching to the dedicated `-floodfill` operator with `-fill none` worked.

### What I learned

In ImageMagick 7, prefer the standalone `-floodfill X,Y color` operator over the `-draw "alpha X,Y floodfill"` form when you want to write transparency into a flood-filled region — the former is more predictable.

### What was tricky

There's a faint anti-aliased pink halo at the edge of the boat where the original PNG had partially transparent pink-on-white pixels. With `-fuzz 10%` it stays. Bumping fuzz higher would risk eroding the inner white in the sail and hull (they're connected to the outer white through thin gaps in the hand-drawn outline). Left as-is — the halo isn't visible against pink-leaning backgrounds, and the alternative is more invasive cleanup.

### What warrants review

Reviewer should pull `/public/images/logo.png` into a viewer that shows alpha (or open over a colored swatch) and confirm: (1) outer is transparent, (2) sail interior is solid white, (3) hull interior is solid white, (4) the pink outline is intact and not aliased to crap. The halo is acceptable; "white square around the boat" is not.

### Future work

If we ever want a crisper edge, redraw the logo as SVG. PNG flood-fill has a ceiling.

## Step 6: Linear wave timing

**Author:** main

### Prompt Context

**Verbatim prompt:** "The waves start off very slow"
**Interpretation:** The `ease-in-out` curve was making the first ~25% of each swing feel sluggish — at 55s periods, that's ~14s of slow ramp every iteration, all visible.
**Inferred intent:** Constant drift speed throughout the swing, not pendulum-like acceleration.

### What I did

Changed all six wave animation utilities in `/tailwind.css` from `ease-in-out` to `linear`. Left the `alternate` / `alternate-reverse` directions and the durations unchanged. The `bob` animation kept `ease-in-out` — at 6s it doesn't have the same problem, and a tiny vertical motion benefits from the gentle ease.

### Why

At these durations and translation distances, the velocity reversal at each turnaround is small enough not to read as a jolt — it just looks like the waves smoothly change direction. The slow start under `ease-in-out` was much more visible than that reversal.

### What worked

Single-character-class change in six places. Verified the compiled CSS picked it up by grepping the served `app.<hash>.css` for `wave-drift NNs linear` patterns.

### What didn't work

Nothing — this one was a clean correction of step 4's choice.

### What I learned

`linear` + `alternate` is the right default for slow background drift. `ease-in-out` is the right default for short, expressive motions (like the logo bob). The break-point is roughly: under 10s, ease-in-out reads as natural; over 20s, linear reads as natural and ease-in-out reads as sluggish-then-frantic.

### What was tricky

Nothing tricky.

### What warrants review

Reviewer should confirm the velocity reversal at each wave turnaround isn't perceptible — open the page and watch one of the slower waves through a full cycle.

### Future work

None.
