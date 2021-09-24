<h1 align="center">Statico</h1>
<p align="center">A static site generator for creative devs</p>

If you like any of my work, you can support me on: https://barelyhuman.dev/donate

[![](https://img.shields.io/badge/license-mit-black?style=for-the-badge)](LICENSE)

## Motivation

I use markdown for most of my writing and really like the simplicity. This started as an experiment to power my websites but is now at a decent level of generalisation to be used by others.

## Features

- Supports Templates (Go Templates)
- Watch Mode
- Local server
- ⚡ Fast
- Markdown => HTML

## Install

- You can download the binaries from the [releases](/releases) page for your specific system and add it to a directory that is in your `PATH` variables.

## Usage

```sh
$ statico [flags]

Usage of statico:
  -s -serve
        alias -serve
  -serve
        Enable file server
  -w -watch
        alias -watch
  -watch
        Start statico in watch mode
```

**statico** doesn't come with any necessary boilerplate and just injects a few variables into the provided templates using a `.config.yml` file. 
You can either use the [default template](https://github.com/barelyhuman/statico-default-template/) or make your own with inspirations and hacking the following available templates 

- [Default](https://github.com/barelyhuman/statico-default-template/)
- [reaper.im](https://github.com/barelyhuman/reaper.im)
- [barelyhuman.dev](https://github.com/barelyhuman/barelyhuman.dev)

The tool needs a `config.yml` and you can find a template in this repository [`config.template.yml`](/config.template.yml)

## Contribute

Issues and PR's are your way to go, fork the repository, create a PR and you're done. Just make sure you let the maintainer know about the issue you pick up to avoid overlaps

## License

[MIT](LICENSE) &copy; Reaper
