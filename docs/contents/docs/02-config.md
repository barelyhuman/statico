---
title: Configuration
published: true
---

## Template

```yml
# Details about the website
site:
  name: BarelyHuman
  link: https://barelyhuman.dev
  description: Indie Developer

# Names of the tempaltes to be used for each scenario
template_names:
  post_template: "blogPostHTML"
  post_index_template: "blogIndexHTML"
  page_template: "simplePageHTML"
  rss_template: "rssTemplate"

# Path definitions
content_path: "./contents"
templates_path: "./templates"
out_path: "./out"
public_folder: "./public"

# Folders that need a subsitute index.html file
# will use the `post_index_template` during generation and
# uses `post_template` for the actual page content
indexed_folders:
  - devlogs

# RSS Options, generates RSS for all the posts from indexed_folders
generate_rss: true
rss_out_path: "./out/blog.xml"

# CLI Options, change the serving port
port: 3000
```
