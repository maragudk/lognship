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
	props.Description = "AI-first software development in the Atmosphere! Agent-friendly, you own your data."

	return Page(props.PageProps,
		Div(Class("space-y-8"),
			H1(Class("text-6xl font-bold"), Text("log’n’ship")),
			P(Class("text-xl"), Text("AI-first software development in the Atmosphere! Agent-friendly, you own your data.")),
			P(a(Class("text-xl"), Href("https://bsky.app/profile/lognship.com"), Text("Follow updates on Bluesky"))),
			P(a(Class("text-xl"), Href("https://pdsls.dev/at://did:plc:ln7l376244zkvhwu4wgr3rzz"), Text("See schemas and more on PDSls.dev"))),
		),
	)
}
