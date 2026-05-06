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
		Div(Class("flex flex-col-reverse sm:flex-row sm:items-center gap-8"),
			Div(Class("space-y-8 max-w-2xl grow"),
				H1(Class("text-3xl sm:text-6xl font-bold tracking-tighter"), Raw(`log’<span class="text-primary-600">n</span>’ship`)),
				P(Class("sm:text-xl"), Raw("Agent-first software development in the Atmosphere! Where <em>you</em> own your data.")),
				P(a(Class("sm:text-xl"), Href("https://bsky.app/profile/lognship.com"), Text("Follow updates on Bluesky"))),
				P(a(Class("sm:text-xl"), Href("https://pdsls.dev/at://did:plc:ln7l376244zkvhwu4wgr3rzz"), Text("See schemas and more on PDSls.dev"))),
			),
			Img(Src("/images/logo.png"), Alt("log’n’ship logo"), Class("w-48 sm:w-64 shrink-0 mx-auto sm:mx-0")),
		),
	)
}
