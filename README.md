<h1 align="center">Statico</h1>
<p align="center">A Markdown Based Static Site Generator</p>

If you like any of my work, you can support me on: https://barelyhuman.dev/donate

[![](https://img.shields.io/badge/license-mit-black?style=for-the-badge)](LICENSE)

## Motivation

I use markdown for most of my writing and really like the simplicity. This started as an experiment to power my websites but is now at a decent level of generalisation to be used by others.

## Features

- Supports Templates (Go Templates)
- Watch Mode
- ⚡ Fast
- Markdown => HTML

## Install

- You can download the binaries from the [releases](/releases) page for your specific system and add it to a directory that is in your `PATH` variables.

## Usage

```sh
$ statico [flags]

Usage of statico:
  -watch Run in watch mode
  -h show this list
```

The cli has just 2 functions,

1. To convert the given folder of markdown into html based on the given templates and has no default templates, you can use [barelyhuman.dev](https://github.com/barelyhuman.dev) as a base if you don't want to write your own templates.

2. To do the same as above but with a watcher, so it can monitor changes for you.

The tool needs a `config.yml` and you can find a template in this repository [`config.template.yml`](/config.template.yml)

## Contribute

Issues and PR's are your way to go, fork the repository, create a PR and you're done. Just make sure you let the maintainer know about the issue you pick up to avoid overlaps

## License

[MIT](LICENSE) &copy; Reaper
