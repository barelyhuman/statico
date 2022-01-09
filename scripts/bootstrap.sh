#!/bin/bash
set -e

# create directories 
mkdir -p public 
mkdir -p content
mkdir -p templates

# create needed files 
touch .gitignore
touch config.yml
touch content/index.md
touch public/styles.css
touch templates/meta.html
touch templates/page.html

# to ignore in version control 
echo "out" >> .gitignore

# templates
config=$(cat <<EOF
site:
  name: notes
  link: https://notes.reaper.im
  description: ""

template_names:
  page_template: "simplePageHTML"

content_path: "./content"
templates_path: "./templates"
out_path: "./out"
public_folder: "./public"
generate_rss: false
EOF
)

pageTemplate=$(cat << EOF
{{define "simplePageHTML"}}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    {{template "metaHTML" .}}
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <link rel="icon" type="image/png" href="/assets/favicon.png" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.Site.Name}} | {{.Meta.Title}}</title>
    <link rel="stylesheet" href="/styles.css" />
  </head>
  <body>
    <main>{{.Content}}</main>
  </body>
</html>
{{end}}
EOF
)

metaTemplate=$(cat << EOF
{{define "metaHTML"}}
<!-- Primary Meta Tags -->
<title>{{.Site.Name}}</title>
<meta name="title" content="{{.Site.Name}}" />
<meta name="description" content="{{.Site.Description}}" />

<!-- Open Graph / Facebook -->
<meta property="og:type" content="website" />
<meta property="og:url" content="{{.Site.Link}}" />
<meta property="og:title" content="{{.Site.Name}}" />
<meta property="og:description" content="{{.Site.Description}}" />
<meta property="og:image" content="{{.Site.Link}}/assets/logo.png" />

<!-- Twitter -->
<meta property="twitter:card" content="summary_large_image" />
<meta property="twitter:url" content="{{.Site.Link}}/" />
<meta property="twitter:title" content="{{.Site.Name}} " />
<meta
  property="twitter:description"
  content="Minimalist | Developer | Designer. 
I develop open source software and tools to ease development and life overall. Dreams of making a living off of open source software."
/>
<meta property="twitter:image" content="{{.Site.Link}}/assets/logo.png" />
{{end}}
EOF
)

indexTemplate=$(cat <<EOF
# Statico
Minimal Boostrapped Template
EOF
)

styleTemplate=$(cat <<EOF
:root {
  --color-black-50: #fafafa;
  --color-black-100: #f4f4f5;
  --color-black-200: #e4e4e7;
  --color-black-300: #d4d4d8;
  --color-black-400: #a1a1aa;
  --color-black-500: #71717a;
  --color-black-600: #52525b;
  --color-black-700: #3f3f46;
  --color-black-800: #27272a;
  --color-black-900: #18181b;

  --space-0: 0;
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 16px;
  --space-4: 32px;

  --font-base: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
    Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
  --text-base: 0.95rem;
  --background: var(--color-black-100);
  --foreground: var(--color-black-600);
  --bright-foregound: var(--color-black-900);
}

html,
body {
  width: 100%;
  height: 100%;
}

body {
  max-width: 900px;
  margin: 0 auto;
  padding: var(--space-2);
  background: var(--background);
  color: var(--foreground);
  font-family: var(--font-base);
}

a {
  color: inherit;
}

a:hover {
  color: var(--bright-foregound);
}

p {
  font-size: var(--text-base);
  line-height: calc(var(--text-base) * 1.8);
}

code,
strong,
em {
  color: var(--bright-foregound);
}
EOF
)

# write to template
echo "$config" > config.yml
echo "$pageTemplate" > templates/page.html
echo "$metaTemplate" > templates/meta.html
echo "$indexTemplate" > content/index.md
echo "$styleTemplate" > public/styles.css

















