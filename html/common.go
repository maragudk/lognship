package html

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"maragu.dev/glue/html"
	. "maragu.dev/gomponents"
	data "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

var hashOnce sync.Once
var appCSSPath, appJSPath string
var datastarJSPath string

func Page(props PageProps, body ...Node) Node {
	hashOnce.Do(func() {
		appCSSPath = getHashedPath("public/styles/app.css")
		appJSPath = getHashedPath("public/scripts/app.js")

		datastarJSPath = getHashedPath("public/scripts/datastar.js")
	})

	return HTML5(HTML5Props{
		Title:       props.Title,
		Description: props.Description,
		Language:    "en",
		Head: []Node{
			Link(Rel("stylesheet"), Href(appCSSPath)),
			Script(Type("module"), Src(datastarJSPath), Defer()),
			Script(Src(appJSPath), Defer()),
			Script(Src("https://cdn.usefathom.com/script.js"), Data("site", "FIRHULCJ"), Defer()),
			html.FavIcons("lognship"),
		},
		Body: []Node{Class("bg-primary-600 text-gray-900 dark:text-white font-mono"),
			Div(Class("min-h-dvh flex flex-col justify-between"),
				header(props),
				Div(Class("grow bg-white dark:bg-gray-800 h-auto relative overflow-hidden"),
					waves(),
					Div(Class("relative"),
						container(true,
							Group(body),
						),
					),
				),
				Div(Class("bg-white dark:bg-gray-800"),
					footer(),
				),
			),
		},
	})
}

func header(_ PageProps) Node {
	return Div(Class("min-h-8 bg-primary-600"),
		container(false),
	)
}

func waves() Node {
	wavePath := `M0,60 Q100,30 200,60 T400,60 T600,60 T800,60 T1000,60 T1200,60 T1400,60 T1600,60 T1800,60 T2000,60 T2200,60 T2400,60`
	wave := func(positionClasses, colorClass, opacityClass, animClass, strokeWidth string) Node {
		return Div(Class("absolute inset-x-0 h-40 "+positionClasses+" "+colorClass+" "+opacityClass+" "+animClass),
			Raw(`<svg class="w-[200%] h-full block" viewBox="0 0 2400 100" preserveAspectRatio="none"><path d="`+wavePath+`" stroke="currentColor" stroke-width="`+strokeWidth+`" stroke-linecap="round" fill="none"/></svg>`),
		)
	}
	return Div(Class("absolute inset-0 pointer-events-none overflow-hidden"), Attr("aria-hidden", "true"),
		wave("top-[5%]", "text-primary-300", "opacity-15", "animate-wave-slow", "1.5"),
		wave("top-[18%]", "text-primary-400", "opacity-10", "animate-wave-medium", "1.5"),
		wave("top-[30%]", "text-primary-500", "opacity-15", "animate-wave-fast-reverse", "2"),
		wave("top-[44%]", "text-primary-300", "opacity-10", "animate-wave-medium-forward", "1.5"),
		wave("top-[58%]", "text-primary-400", "opacity-15", "animate-wave-slow-reverse", "2"),
		wave("top-[72%]", "text-primary-500", "opacity-10", "animate-wave-fast", "2.5"),
		wave("top-[85%]", "text-primary-400", "opacity-15", "animate-wave-medium", "1.5"),
	)
}

func footer() Node {
	return Div(
		container(false,
			Div(Class("flex items-center justify-center space-x-4 sm:space-x-8 py-2"),
				data.Init("$counter = 0"),
				a(Href("https://www.maragu.dev"), Text("by maragu")),
			),
		),
	)
}

func a(children ...Node) Node {
	return A(JoinAttrs("class", Class("text-primary-600 hover:text-primary-400 dark:text-white dark:hover:text-gray-700"), Group(children)))
}

func container(padY bool, children ...Node) Node {
	return Div(
		Classes{
			"max-w-7xl mx-auto h-full": true,
			"px-4 sm:px-6 lg:px-8":     true,
			"py-4 md:py-8":             padY,
		},
		Group(children),
	)
}

func getHashedPath(path string) string {
	externalPath := strings.TrimPrefix(path, "public")
	ext := filepath.Ext(path)
	if ext == "" {
		panic("no extension found")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("%v.x%v", strings.TrimSuffix(externalPath, ext), ext)
	}

	return fmt.Sprintf("%v.%x%v", strings.TrimSuffix(externalPath, ext), sha256.Sum256(data), ext)
}
