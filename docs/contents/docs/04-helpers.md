---
title: Helpers
published: true
---

There's a few helpers that are added to certain types of data

## `Format`

<small>Added in: v0.0.5</small>
<small>Prop: <code>Metadata.Date</code></small>

Returns a string date in the requested format layout (compatible with Time.Format Layouts) from go

```html
<!-- Example -->
<ul class="p-0">
  {{range .Files}}
  <li class="mb-2 list-style-none">
    <a href="/{{.Slug}}"> {{.Title}} </a>
    <p class="m-0 p-0">
      <small> {{.Date.Format "02/01/2006"}} </small>
    </p>
  </li>
  {{end}}
</ul>
```
