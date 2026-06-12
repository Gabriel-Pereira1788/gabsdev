---
title: "Rewriting this blog: from Next.js to Go + templ + HTMX"
date: "2026-06-11"
tags: ["go", "templ", "htmx", "nextjs", "web"]
excerpt: "Does every site really need a full-blown Next.js, or can we simplify things? I rewrote this very blog in plain Go, templ and HTMX to answer that question in practice."
readMin: 10
---

If there's one decision that has become almost automatic in modern web development, it's reaching for the same usual suspect on every project. Want a blog? Next. A portfolio? Next. A landing page? Next. But after all, does every site really need a full-blown Next.js, or can we simplify things?

This blog you're reading right now started exactly like that: a Next.js project. In this article I'll show you how I rewrote it from scratch in **plain Go, templ and HTMX**. And, more important than the "how", why this path makes sense for a site that is, at the end of the day, mostly content.

## The weight of the modern default

Before we talk about the solution, let's understand the problem. A site like this one has a simple nature: it renders a list of articles, shows an article, filters by tags. Almost none of that is *truly* interactive. It's content that needs to arrive fast and be read.

Next does render on the server, sure, nobody's saying otherwise. But it still ships React to the browser and hydrates the page on the client. And that package charges a price by default, even when you don't need it:

- A JavaScript bundle the browser has to download, parse and execute before the page becomes useful.
- Hydration: the server sends HTML, and the client re-runs everything to "bring that HTML to life".
- A build pipeline with dozens of dependencies that age quickly.

For an app full of client-side state, that cost pays off. For a blog, you're loading a truck to deliver a letter. **One important point to highlight is that** the problem isn't React. It's using the heavy interactivity tool in a place where interactivity is the exception, not the rule.

## Simplicity as a choice

It's worth separating two kinds of complexity. There's the essential one, which is the problem itself: rendering articles, searching, filtering by tag. And there's the accidental one, the kind the tool drags along without the problem asking for it. Next handles the essential part very well, but it charges a fair dose of the accidental: a runtime to hydrate, a build pipeline, a dependency graph to keep standing. For this blog, almost all of that weight was accidental.

Switching stacks was, above all, throwing pieces away. There's no `node_modules` to install, no Node toolchain to version, no frontend build to break on a dependency bump. What's left is a Go binary and a handful of `.templ` files. Fewer pieces means less surface that ages, less to update when some library decides to change its API.

And there's a gain that shows up in no bundle: the whole project fits in your head. The path of any request is always the same, Route → Handler → Service → Component, and you can follow it end to end without opening the docs of three libraries. Simplicity here isn't about doing less. **It's about removing what wasn't helping, so that what remains is easy to understand and hard to break.**

Now that it's clear why simplifying pays off, let's look at the mental model that makes all of this possible.

## The inverted pyramid: HTML on the server

The idea behind the rewrite is to flip that pyramid upside down. Instead of "JavaScript on the client by default, HTML as a detail", the rule becomes **HTML on the server by default, JavaScript only where it hurts**.

This mental model has three layers, and it's worth fixing them before looking at any code:

- **Go** is the server. It takes the request, fetches the data, decides what to render.
- **templ** is the HTML, but type-checked. Components that compile into real Go functions.
- **HTMX** is the pinpoint interactivity. Small pieces of HTML swapped on demand, with no JavaScript to write.

Now that we understand the model, let's finally see how it materializes in this project.

## Hands on: the skeleton

The first surprise for someone coming from Next is finding out you don't need any framework for the server. Go's standard library already gives you routing, and the `cmd/server/main.go` of this entire blog fits on one screen:

```go
mux := http.NewServeMux()

mux.HandleFunc("/", handlers.Home)
mux.HandleFunc("/articles", handlers.Articles)
mux.HandleFunc("/articles/search", handlers.ArticlesSearch)
mux.HandleFunc("/articles/{slug}", handlers.ArticleDetail)
mux.HandleFunc("/about", handlers.About)
mux.HandleFunc("/tags", handlers.Tags)
mux.HandleFunc("/rss.xml", handlers.RSS)

mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

log.Println("Server on http://localhost:8080")
http.ListenAndServe(":8080", mux)
```

Notice the `{slug}` on the `/articles/{slug}` line: a route with a parameter, no external dependency. Each route points to a *handler*, and the pattern that repeats across the whole project is simple: **Route → Handler → Service → Component**. The handler fetches the data from a service and asks a component to render. That's it.

## templ instead of JSX

This is the part that feels most familiar to anyone coming from React. With **templ**, you write components that look a lot like JSX, but compile into Go functions with checked types. Here's the articles page:

```go
templ ArticlesPage(pathname string, posts []domain.Post, tags []string, activeTag string) {
	@layouts.Layout("Articles", pathname) {
		<div class="page enter">
			@components.SearchBar(activeTag)
			<div class="post-list" id="post-list">
				for _, post := range posts {
					@components.PostRow(post, activeTag)
				}
			</div>
		</div>
	}
}
```

If you've written React, this feels familiar: the layout receives children (templ's `{ children... }`), components compose with `@components.PostRow(...)`, and the `for` is the `.map()` you already know, except it's a real Go `for`.

**One important point to highlight is that** this `.templ` is not interpreted at runtime. Before compiling, you run `templ generate` and each file becomes a matching `_templ.go`, which is ordinary Go code. That means two things: if you pass the wrong type to a component, the build breaks (not your user's browser); and there's no reflection or template parsing in production, just already-compiled string concatenation.

## HTMX: interactivity without a framework

"Fine, but what about search? And the tag filter? That's interactive." This is where HTMX comes in, and where the rewrite proves you can have interactivity without dragging a whole framework along.

HTMX enters through a single `<script>` in the layout (loaded from a CDN) and works by extending HTML with attributes. Here's this blog's search field:

```go
templ SearchBar(activeTag string) {
	<div class="search-wrap">
		<input
			class="search-input"
			type="text"
			name="q"
			placeholder={ searchPlaceholder(activeTag) }
			hx-get="/articles/search"
			hx-trigger="input changed delay:300ms, search"
			hx-target="#post-list"
			hx-swap="innerHTML"
			hx-vals={ `{"tag":"` + activeTag + `"}` }
		/>
	</div>
}
```

Let's read these attributes as a sentence: when the input *changes* (`hx-trigger`, with a free 300ms debounce), make a `GET` to `/articles/search` (`hx-get`), take whatever HTML comes back and swap the contents of `#post-list` with it (`hx-target` + `hx-swap`). I didn't write a single line of JavaScript for this.

And what does the server reply with? Here's the prettiest detail of the approach. The handler **does not return the whole page**. It returns only the piece that changed:

```go
func ArticlesSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	tag := r.URL.Query().Get("tag")

	posts := services.SearchPosts(query, tag)

	components.PostList(posts, tag).Render(r.Context(), w)
}
```

Compare it with the full-page handler: that one renders `pages.ArticlesPage(...)`, wrapped in `@layouts.Layout`. This one renders only `components.PostList(...)`, the same list component, without the layout around it. **Same mental route, two response sizes:** the full page on navigation, the fragment on interaction. HTMX takes care of stitching the fragment back into the screen.

Clicking a tag follows the exact same idea, with a bonus `hx-push-url` so the URL stays shareable:

```go
<span
	class={ getTagClass(tag, activeTag) }
	hx-get={ "/articles/search?tag=" + tag }
	hx-target="#post-list"
	hx-swap="innerHTML"
	hx-push-url={ "/articles?tag=" + tag }
>{ tag }</span>
```

## Content as files, not a database

There's one more decision the rewrite let go of, at least for now: the database. And here's an honest note about the project's current state. I chose, on purpose, not to spin up any database yet. The idea was the simplest possible one: write each article in a Markdown file and keep it locally, versioned alongside the code.

**One important point to highlight is that** this is a temporary choice, not a dogma. Down the road I plan to migrate all this content to a real database and build an admin dashboard, just to keep a reasonable grip on the data. But for the stage this blog is at today, a database would be more of a piece to maintain than a problem solved. So let's take the lighter path while it holds up.

In practice, a blog's content is the articles, and articles are text. Each post becomes a Markdown file in `content/posts/`, with a YAML frontmatter on top (including the very one you're reading).

The service reads the folder, splits the frontmatter from the body with a `SplitN` on `---`, and runs the Markdown through goldmark with GFM and syntax highlighting:

```go
md := goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle("nord"),
		),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)
```

`WithAutoHeadingID` even generates the heading IDs automatically, which feeds the side table of contents you see next to this article. The upside of this choice is hard to overstate: the content lives in git, versioned alongside the code, with no migration, no admin panel, no database to spin up on deploy.

## Upsides (and the honest trade-offs)

Rewriting this blog produced concrete gains over Next:

- **Less JavaScript on the client.** What reaches the browser is ready HTML. The only meaningful JS is HTMX itself, small and stable.
- **No frontend build.** The build step is `templ generate && go build`. The generated `_templ.go` files are committed, which makes deploying on Railway trivial, with no Node toolchain and no `node_modules`.
- **End-to-end type safety.** The same compiler that validates the service validates the template.
- **Simplicity that fits in your head.** Route → Handler → Service → Component. You can read the whole project in an afternoon.

But I wouldn't be honest if I sold this as magic. The Go + HTMX ecosystem is smaller than React's. You'll find fewer ready-made components and fewer answers on Stack Overflow. And genuinely rich interactions (a text editor, a canvas, complex drag-and-drop) still call for real JavaScript. HTMX shines at swapping fragments, not at complex client state. There are good solutions to cover that reactivity gap, like Alpine.js, which sits well next to HTMX for the local state that lives in the browser. Even so, the question was never "which one is better", but rather **what's the right weight for the right problem**.

## Conclusion

It was quite a journey going from Next to a server that returns ready HTML 🚀. We took a look at how Go handles routing on its own, how templ brings JSX ergonomics with type checking, and how HTMX delivers interactivity without dragging a whole framework along.

Let me be clear: this is not an anti-React manifesto. Next remains an excellent choice for real apps full of client-side state. The point is different. Before accepting the default, it's worth asking how much interactivity your project actually needs. In this blog's case, the answer was "almost none", and the stack came to reflect that.

I hope this journey makes you look at your next project and consider the simpler path before grabbing the truck. Thanks for reading this far 🙏. And, of course, this whole blog is open source: go poke around the code.
