# mediumctl

[![GitHub release](https://img.shields.io/github/release/moutend/mediumctl.svg?style=flat-square)][release]
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![CircleCI](https://circleci.com/gh/moutend/mediumctl/tree/master.svg?style=svg&circle-token=7555891ce81c70dfcfd4279e73d9043a53e5129d)][status]

[release]: https://github.com/moutend/mediumctl/releases
[license]: https://github.com/moutend/mediumctl/blob/master/LICENSE
[status]: https://circleci.com/gh/moutend/mediumctl

`mediumctl` is CLI tool for posting an article to Medium.

# Installation

## Windows / Linux

You can download the executable for 32 bit / 64 bit at [GitHub releases page](https://github.com/moutend/gip/releases/).

## Mac

The easiest way is Homebrew.

```shell
$ brew tap moutend/homebrew-mediumctl
$ brew install mediumctl
```

# Usage

## Before starting

First off, you need login to Medium with `auth` subcommand.
Go to https://medium.com/me/applications, and please create a new OAuth application.
You can specify any client name and description, but you must specify your local IP address (e.g. `192.168.1.2:4000`) as the redirect URI.

Internally, `mediumctl auth` launches local web server with given flags  and obtain the shortlive code for generating API token. This is the most tricky part.


```shell
$ mediumctl auth -h http://192.168.1.2 -p 4000 -i YOUR_CLIENT_ID -s YOUR_CLIENT_SECRET
```

And then the web browser will be automatically opened, please click OK to continue.
After clicking OK, the browser will be automatically closed and your API token will be stored if the authorization was successfully done.

```shell
```

Now you can post an article to your user page an your publications.

## Post an article to the publication

You can post an HTML file or a Markdown file with special syntax.

**example.md**

```markdown
---
title: The best way to learn Go
tags: golang programming
status: draft
---

# A Tour of Go

If you are looking for the way to learn go, a tour of Go is the best way to learn go.

## Why Go?

Simple is not equal to easy, but simple made me easy.
```shell

To post an article above, please run the command below:

```shell
$ mediumctl post example.md
```

If the article was successfully posted, you can see the message like:

```shell
```

| property | Description | Default value |
|:--|:--|:--|
| title | Title of the article | `Untitled` |
| tags | Tags associated with the article. You can not specify more than 3 tags. | empty |
| status | You can specify one of `public`, `draft` and `editing`. `draft` |
| canonicalURL | Canonical URL for the article. | empty |
| publication | URL or publication number. | empty |

# LICENSE

MIT
