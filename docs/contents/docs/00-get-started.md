---
title: Get Started
published: true
---

## Install

- You can download the binaries from the [releases](/releases) page for your specific system and add it to a directory that is in your `PATH` variables.

## Usage

```sh
$ statico [flags]

Usage of statico:
  -c -config
        alias -config
  -config string
        Config file to use
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
