package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type HomePageProps struct {
	PageProps
}

func HomePage(props HomePageProps) Node {
	props.Title = "log’n’ship"
	props.Description = "Agent-first software development in the Atmosphere! Where you own your data."

	return Page(props.PageProps,
		Div(Class("space-y-16"),
			Div(Class("flex flex-col-reverse sm:flex-row sm:items-center gap-8"),
				Div(Class("space-y-8 max-w-2xl grow"),
					H1(Class("text-3xl sm:text-6xl font-bold tracking-tighter"), Raw(`log’<span class="text-primary-600">n</span>’ship`)),
					P(Class("sm:text-xl"), Raw("Agent-first software development in the Atmosphere! Where <em>you</em> own your data.")),
					P(a(Class("sm:text-xl"), Href("https://bsky.app/profile/lognship.com"), Text("Follow updates on Bluesky"))),
					P(a(Class("sm:text-xl"), Href("https://pdsls.dev/at://did:plc:ln7l376244zkvhwu4wgr3rzz"), Text("See schemas and more on PDSls.dev"))),
				),
				Img(Src("/images/logo.png"), Alt("log’n’ship logo"), Class("w-48 sm:w-64 shrink-0 mx-auto sm:mx-0 animate-bob")),
			),
			Div(Class("space-y-6 max-w-3xl"),
				H2(Class("text-2xl sm:text-4xl font-bold tracking-tighter"), Text("Principles")),
				Ul(Class("space-y-4 sm:text-xl list-none pl-0"),
					Li(Raw("<strong>The product, not the code, is the artifact.</strong> Decisions, tasks, and bugs are first-class. Code lives in a repo.")),
					Li(Raw("<strong>The log is the unit.</strong> An append-only stream — read it for context, write to add to it. Log, then ship.")),
					Li(Raw("<strong>A layer between agents and humans.</strong> The shared surface for coordination.")),
					Li(Raw("<strong>Your data, on the open network.</strong> Records live in your atproto repository. lognship is one view.")),
				),
			),
		),
	)
}
