---
title: Templates
published: true
---

**Statico** was built with the idea of letting the developer decide what the resulting website looks like and so there's no default template that it comes with. You decide on it and you can do so with simple HTML templates or make it a little more dynamic by using Go Templates to compile the pages as needed.

A good example of this would be the following websites

- [barelyhuman.dev](https://barelyhuman.dev)
- [reaper.im](https://reaper.im)
- [statico.reaper.im](https://statico.reaper.im)

## Predefined Variables

To work with the dynamic templates, there's going to be certain variables passed by the **statico** to help you.

### Global

The following variable and it's children will be available in all templates and can be used with the general format (`{{ .VariableName.ChildProp}}`) of go templates, refer to existing templates of the above mentioned websites for more info

```go
type Site struct {
	Name        string
	Link        string
	Description string
}
```

### `post_template` and `page_template` | Post Template and Page Template

The following are the variables that will be passed to the `post_template` and `page_template` defined go template in `config.yml`

```go
type Post struct {
	Site    Site // from the global variable
	Meta    Metadata
	Content string
}

type Metadata struct {
	Published  bool   // comes from the post's frontmatter key `published`
	Title      string // comes from the post's frontmatter key `title`
	Date       Date   // comes from the post's frontmatter key `date`
	ImageURL   string // comes from the post's frontmatter key `image_url`
	AGImageURL ImageURLGen
	Slug       string // current url path for the file
	Content    string // compiled content of the file
	OutPath    string
}
```

```html
<!-- example -->
<h1>{{.Site.Name}}</h1>
<main>
  <h2>{{.Meta.Title}}</h2>
  <article>{{.Content}</article>
</main>
```

### `post_index_template` / Post Index

The following are the variables that will be passed to the `post_index_template` defined go template in `config.yml`

```go
type IndexedFiles struct {
	Site  Site
	Files []Metadata
}

type Metadata struct {
	Published  bool   // comes from the post's frontmatter key `published`
	Title      string // comes from the post's frontmatter key `title`
	Date       Date   // comes from the post's frontmatter key `date`
	ImageURL   string // comes from the post's frontmatter key `image_url`
	AGImageURL ImageURLGen
	Slug       string // current url path for the file
	Content    string // compiled content of the file
	OutPath    string
}
```

```html
<!-- example -->
{{range .Files}}
<li class="mb-2 list-style-none">
  <a href="/{{.Slug}}"> {{.Title}} </a>
</li>
{{end}}
```

### `rss_template` / RSS Template

The following are the variables that will be passed to the `rss_template` defined go template in `config.yml`

```go
    type ATOMFeed struct {
	Site struct {
		Name        string
		Link        string
		Description string
	}
	Posts []struct {
		Slug    string
		Title   string
		Link    string
		Content string
		Date    time.Time
	}
}
```

```html
<!-- example -->
{{define "rssTemplate"}}<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
  <title>{{.Site.Name}}</title>
  <link>{{.Site.Link}}</link>
  <description>{{.Site.Description}}</description>
  <atom:link href="{{.Site.Link}}" rel="self" type="application/rss+xml" />
      {{range .Posts}}
      <item>
      <guid>{{.Slug}}</guid>
      <title>{{.Title}}</title>
      <link>{{.Link}}</link>
      <description>
        <![CDATA[{{.Content}}]]>
      </description>
      <pubDate>{{.Date}}</pubDate>
    </item>
      {{end}}
</channel>
</rss>{{end}}
```
