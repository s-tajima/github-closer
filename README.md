github-closer
---
[![Build Status](https://travis-ci.org/s-tajima/github-closer.svg?branch=master)](https://travis-ci.org/s-tajima/github-closer)

A useful tool for closing GitHub Issue/Pull Requests that be left for a long time.

```
$ ./github-closer --debug
2016/12/15 22:54:18 #10 https://github.com/s-tajima/github-closer/issues/10 left for a long time issue is closed.
2016/12/15 22:54:19 #999 https://github.com/s-tajima/github-closer/issues/999 recently created issue was updated recently. skipped.
```

## Index

* [Concepts](#concepts)
* [Requirements](#requirements)
* [Installation](#installation)
* [Configure](#configure)
* [Usage](#usage)       
* [License](#license)    

## Concepts

* It would be better to keep status that only Issues/Pull Requests that actually needed is opened.
* So, let's automate to closing Issues/Pull Requests that be left for a long time.
* No worries, Issues/Pull Requests is easily reopened if you desired.

## Requirements

github-closer requires the following to run:

* Golang

## Installation

```
$ go get github.com/s-tajima/github-closer
```

## Configure

Set your configuration as Environment Variables.
```
export GITHUB_ACCESS_TOKEN=
export GITHUB_ORGANIZATION=
export GITHUB_REPO=
export GITHUB_TEAM=
```
You can use .env file as well.


## Usage

```
Usage:
  github-closer [OPTIONS]

Application Options:
  -q, --query=    Query strings. For search Issues/Pull Requests. (default: is:issue is:open)
  -d, --duration= Duration. Issues would be closed if left over this duration. (days) (default: 7)
  -c, --comment=  Comment. Would be posted before an Issue is closed. (default: :alarm_clock: this Issue was left for a long time.)
  -n, --dry-run   If true, show target Issues without closing.
  -o, --run-once  If true, close only one Issue.
  -l, --limit=    A maximum number of closed Issues. (default: 0)
      --debug

Help Options:
  -h, --help      Show this help message
```

## License

[MIT](./LICENSE)

## Author

[Satoshi Tajima](https://github.com/s-tajima)
