---
title: CLI
published: true
---

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

## Flags

### `-c -config`

Self Explanatory, you pass the path to the config to be used, a template can be found on the [Configuration](/docs/02-config.html) page.

```sh
# eg:
> statico -c config.yml
# or
> statico -config ./path/to/config.yml
```

### `-s -serve`

Statico comes with it's own http server which cna be chained with other flags to give you an easier development environment. The `PORT` that it runs on can be changed in the configuration file

```sh
# eg:
> statico -s
# or
> statico -serve

# general use case
> statico -s -w
```

### `-w -watch`

Statico also comes with a file watcher which currently monitors the files and directories mentioned in the configuration, so any change or addition of file in the user configured directory will trigger a re-compile and will be served accordingly.

```sh
# eg:
> statico -w
# or
> statico -watch

# general use case
> statico -s -w
```
