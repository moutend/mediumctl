# mediumctl

[![GitHub release](https://img.shields.io/github/release/moutend/mediumctl.svg?style=flat-square)][release]
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![CircleCI](https://circleci.com/gh/moutend/mediumctl.svg?style=svg&circle-token=e7748578056ded93a5532904c047fc0f23db3bba)](https://circleci.com/gh/moutend/mediumctl)

[release]: https://github.com/moutend/mediumctl/releases
[license]: https://github.com/moutend/mediumctl/blob/master/LICENSE
[status]: https://circleci.com/gh/moutend/mediumctl

`mediumctl` is a CLI tool for publishing HTML or Markdown file to the Medium.

# Installation

## Windows / Linux

You can download the executable for 32 bit / 64 bit at [GitHub releases page](https://github.com/moutend/mediumctl/releases/).

## Mac

The easiest way is Homebrew.

```shell
$ brew tap moutend/homebrew-mediumctl
$ brew install mediumctl
```

## `go install`

If you have already set up Go environment, just `go install`.

```shell
$ go install github.com/moutend/mediumctl
```

# Usage

## Setting up API token with OAuth

First off, you need set up API token. Open https://medium.com/me/applications, please create a new OAuth application.
You can specify any client name and description, but you must specify local IP address assigned your machine (e.g. `192.168.1.2`) as the redirect URI.

After creating OAuth application, please run `auth` command the following flags:

- `-i` ... client ID
- `-s` ... client secret
- `-u` ... redirect URI

In the following example, it assumes that you have specified `http://192.168.1.2:4567` as the redirect URI for OAuth application.

```shell
$ mediumctl auth -u http://192.168.1.2:4567 -i CLIENT_ID -s CLIENT_SECRET
```

Then browser will be automatically opened, please check the listed grant types and press OK to continue.
If you would be asked network access permission during this step, please allow that permission.
Internally, `mediumctl` launches local web server with given host name and port number, and then extract the shortlive code from redirected HTTP request to generate API token. This is the most tricky part of `mediumctl`.

Your API token will be saved at `$HOME/.mediumctl` and the web browser will be closed automatically.
Now you can post an article to your user profile and your publications.

## Setting up API token with self-issued token

Alternatively, it's not recommended but you can set up an API token by hand.
Open https://medium.com/me/settings, please generate self-issued API token.
Then create a JSON file that contains following key-value pairs at `$HOME/.mediumctl`.

```js
{
  "ApplicationID": "",
  "ApplicationSecret": "",
  "AccessToken": "SELF_ISSUED_TOKEN",
  "ExpiresAt": 0,
}
```

Please replace `SELF_ISSUED_TOKEN` to your self-issued token. This file must be treated as the password and do NOT expose it.

## Publishing an article

Publishing an article is easy.

```shell
$ mediumctl publication example.md
```

That's it! The Markdown file `example.md` will be published at your publication.
If you have more than one publications, you can specify which publication to be published at. For more details, please read the next section.

In the same manner, you can publish an article to the your user profile.

```shell
$ mediumctl user example.html
```

The HTML file `example.html` will be published at your profile.

# Frontmatter for HTML and Markdown

You can provide additional information with frontmatter. The following table shows what property can be used.

| Property | Description | Default value |
|:--|:--|:--|
| `title` | Title of the article | `Untitled` |
| `tags` | Tags associated the article. Only three tags can be specified. | blank |
| `status` | One of `public`, `draft` and `unlisted`. `public` |
| `number` | Publication number displayed when you run `info` command. | `0` |
| `notify` | Whether notify followers that the user has published new article. | `false` |
| `license` | License of the article listed below. | `all-rights-reserved` |
| `canonicalURL` | Canonical URL for the article. | blank |

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

For example, you can create a Markdown file like this:

```markdown
---
title: The best way to learn Go
tags: golang programming
status: draft
canonicalURL: https://blog.example.com/the-best-way-to-learn-go
---

# The best way to learn Go

If you're looking for the best way to learn Go, this article might help you.

# A Tour of Go

As you know, [a tour of Go](https://golang.org) is the best way to learn go.

## Why Go?

Simple is not equal to easy, but simple made you easy.
```

# Tips for publishing HTML / Markdown

## Valid HTML tags

Some HTML tags cannot be used. For a full list of accepted HTML tags, please see [Medium API documentation](https://medium.com/@katie/a4367010924e).

## Special heading element

If you have the markdown like:

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

In the example above, first heading level 1 element `# Title of the article` will be treated as heading level 1.
However, second and third heading level 1 elements are treated as heading level 2.
This is the specification of Medium API and you cannot change this behavior 

Note that only the first heading level 1 is treated as heading level 1.
I recommend you to specify the first heading level 1 as same as title of the article.

## API limitation

If you published a large number of articles in a short time, publishing would be restricted.

```shell
$ mediumctl publication example.md
error: User has reached the rate limit for publishing today. (code:-1)
```

# LICENSE

MIT
