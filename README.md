mediumctl
=========

[![GitHub release](https://img.shields.io/github/release/moutend/mediumctl.svg?style=flat-square)][release]
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![CircleCI](https://circleci.com/gh/moutend/mediumctl.svg?style=svg&circle-token=e7748578056ded93a5532904c047fc0f23db3bba)](https://circleci.com/gh/moutend/mediumctl)

[release]: https://github.com/moutend/mediumctl/releases
[license]: https://github.com/moutend/mediumctl/blob/master/LICENSE
[status]: https://circleci.com/gh/moutend/mediumctl

`mediumctl` is a CLI tool for publishing HTML / Markdown to the Medium.

## Installation

### Windows / Linux

You can download the executable for 32 bit / 64 bit at [GitHub releases page](https://github.com/moutend/mediumctl/releases/).

### Mac

The easiest way is Homebrew.

```console
brew tap moutend/homebrew-mediumctl
brew install mediumctl
```

### `go build`

If you have already set up Go environment, just `go build`.

Note: Go v1.12 required.

```console
git clone https://github.com/moutend/mediumctl
$ go build
```

## Setup API token with OAuth

First, you need create an OAuth application. Visit https://medium.com/me/settings. You'll find the Developer section, please create an OAuth application for `mediumctl`.

You can set the redirect URL to the localhost. However, the host name `https://localhost` is not accepted, use the IP address `http://127.0.0.1`.

After creating the OAuth application, run `auth` command with the following flags:

- `-u` ... redirect URI
- `-i` ... client ID
- `-s` ... client secret

For example:

```console
mediumctl auth -u http://127.0.0.1:4000 -i xxxxxxxx -s xxxxxxxx
```

Because the `mediumctl` launches the local web server to obtain the response from the API server, allow the network access when the permission is required.

Now you can post an article to your user profile and your publications.

### Setup API token with self-issued token

Alternatively, it's not recommended but you can set up an API token by hand.

Go to https://medium.com/me/settings, please generate self-issued API token.

Then create a JSON file that contains following key-value pairs at `$HOME/.mediumctl`.

```javascript
{
  "AccessToken": "SELF_ISSUED_TOKEN"
}
```

`SELF_ISSUED_TOKEN` is your self-issued token. This file must be treated as the password and do NOT expose it.

## Publish your article

```console
mediumctl publication ./example.md
```

That's it! The Markdown file `example.md` will be published at your publication.

If you have more than one publications, you can specify which publication to be published at. For more details, please read the next section.

You can also publish an article to the your user profile.

```console
mediumctl user ./example.html
```

The HTML file `example.html` will be published at your profile.

## Get information about the user and its publications

To get information about the user and its publications, use `info` command.

```console
mediumctl info
```

## HTML and Markdown format

You can provide additional information with frontmatter. The following table shows what property can be used.

| Property | Description | Default value |
|:--|:--|:--|
| `title` | Title of the article | `Untitled` |
| `tags` | Tags associated the article. Only three tags can be specified. | blank |
| `publishedAt` | The date that the article published at. (Use RFC3339 date and time format) | current time |
| `status` | One of `public`, `draft` and `unlisted`. `public` |
| `number` | Publication number displayed when you run `info` command. | `0` |
| `notify` | Whether notify followers that the user has published new article. | `false` |
| `license` | License of the article listed below. | `all-rights-reserved` |
| `curl` | Canonical URL for the article. | blank |

note: You can't specify the future date as the publish date. Also, you can't specify the older date before Jan 1st, 1970.

### Valid licenses for `license`

Valid values for `licence` are:

- all-rights-reserved
- cc-40-by
- cc-40-by-sa
- cc-40-by-nd
- cc-40-by-nc
- cc-40-by-nc-nd
- cc-40-by-nc-sa
- cc-40-zero
- public-domain

### Example

```markdown
---
title: The best way to learn Go
tags: golang programming
status: draft
curl: https://blog.example.com/the-best-way-to-learn-golang
---

# The best way to learn Go

If you're looking for the best way to learn Go, this article might help you.

# A Tour of Go

As you know, [a tour of Go](https://golang.org) is the best way to learn go.
```

### Tips for publishing HTML / Markdown

#### Valid HTML tags

You can use the accepted HTML tags. For a full list of accepted HTML tags, please see [Medium API documentation](https://medium.com/@katie/a4367010924e).

#### `h1` conversion

Note that heading elements are automatically converted according to the following rules.

| Before | After |
|:--|:--|
| The first `h1` | `h1` (title of the article) |
| The second and subsequent `h1` | `h3` |
| The first `h2` | `h2` (Subtitle of the article) |
| The second and subsequent `h2` | `h4` |
| `h3` | `h4` |
| `h4` | `h4` |
| `h5` and `h6` | `p` (Normal paragraph) |

Note that the only first h1 and h2 are treated as title and sub title of the article.

For example, if you have the markdown file like this:

```markdown
---
title: Title of the article
---

# Title of the article

# Heading level 1

first paragraph ...

# Heading level 1

second paragraph ...
```

The first `h1` treated as `h1`. The second and third `h1` are treated as `h3`.

#### API limitation

If you did publish the many articles in a short time, you may restricted temporary.

```console
mediumctl publication example.md
error: User has reached the rate limit for publishing today. (code:-1)
```

If you got the rate and limit error, you must wait 24 hours. Then you'll be activated automatically.

## LICENSE

MIT

## Author

[Yoshiyuki Koyanagi <moutend@gmail.com>](https://github.com/moutend)
