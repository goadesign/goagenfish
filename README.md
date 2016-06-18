# goagenfish

goagenfish is a fish shell completion script generator for
[goagen](https://goa.design/implement/goagen/), the generation tool of
[goa](https://gihub.com/goadesign/goa).

## Installation

Assuming a correct [Go](https://golang.org) setup:

Install `goagenfish`:

```
go get -u github.com/goadesign/goa/goagen
go get -u github.com/goadesign/goagenfish
```

And run it:

```
goagenfish
goagen.fish
```

Move the generated file to a fish completion directory, for example:

```
mkdir -p ~/.config/fish/completions
mv goagen.fish ~/.config/fish/completions
```

## Usage

Newly opened fish shell have now "smart" completion for `goagen`. For example:

```
goagen <tab>
app        client  js    schema 
bootstrap  gen     main  swagger
```

Or:

```
goagen app -d github.com/goadesign/<tab>
…oadesign/examples/types/design          (design package import path)  …oadesign/gorma/example/design                     (design package import path)
…oadesign/examples/upload/design         (design package import path)  …oadesign/oauth2/design                            (design package import path)
…oadesign/examples/websocket/design      (design package import path)  …oadesign/swagger-service/design                   (design package import path)
…oadesign/goa-cellar/design              (design package import path)  
…and 5 more rows
```

## Features

Auto-completion is provided for:

* All command names
* All flag names
* `--design` flag values (only package import paths to `design` package directories are shown)
* `--pkg-path` flag values (only package import paths are shown)
