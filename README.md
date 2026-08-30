# gabsdev-go

Personal blog rewritten in Go — no framework, no unnecessary dependencies.

## About

A retro-pixel aesthetic blog. Posts written in Markdown with YAML frontmatter, rendered server-side with syntax highlighting. No JS framework, no database — just files.

## Stack

| Technology | Role |
|---|---|
| [Go](https://go.dev) | Core language — HTTP server using stdlib `net/http` |
| [templ](https://templ.guide) | Type-safe HTML templating in Go |
| [Goldmark](https://github.com/yuin/goldmark) | Markdown parser and renderer |
| [Goldmark Highlighting](https://github.com/yuin/goldmark-highlighting) | Syntax highlighting via Chroma |
| [Air](https://github.com/air-verse/air) | Live reload for development |

## Structure

```
gabsdev-go/
├── cmd/server/         # Entrypoint — routes and HTTP server
├── content/posts/      # Blog posts as Markdown files
├── internal/
│   ├── domain/         # Domain structs (Post, PostMeta)
│   ├── handlers/       # HTTP handlers (home, articles, about, tags)
│   ├── services/posts/ # Post reading and search logic
│   ├── utils/          # Markdown, TOC, date formatting, icons
│   └── views/          # templ templates (layouts, pages, components)
└── static/             # Static assets (CSS, images, favicon)
```

## Pages

- `/` — Home
- `/articles` — Post list with search
- `/articles/{slug}` — Single post
- `/tags` — Posts grouped by tag
- `/about` — About

## Running locally

**Requirements:** Go 1.25+, [templ CLI](https://templ.guide/quick-start/installation), [air](https://github.com/air-verse/air)

```bash
make dev      # live reload (air)
make run      # generate + run once, no reload
make build    # compile binary to ./bin/server
make generate # regenerate *_templ.go from .templ sources
make fmt      # templ fmt + gofmt
make test     # go test ./...
```

Without `make`: `air`, or `templ generate && go run ./cmd/server/main.go`.

Server starts at `http://localhost:8080`.

## Writing posts

Drop a `.md` file in `content/posts/` with YAML frontmatter:

```markdown
---
title: "Post title"
date: "2025-06-05"
tags: ["go", "web"]
excerpt: "Short description shown in the post list."
readMin: 5
slug: "post-title"
---

Markdown content here.
```
