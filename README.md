# gscript

Genesis Scripting Engine

![Genesis Logo](http://www.city-church.org.uk/sites/default/files/styles/sidebar_series_image/public/series/image/genesis-in-the-beginning.jpg?itok=EJFz0LWt)

**WARNING**: This library is under active development. API is ***NOT*** stable and will have breaking changes for the foreseeable future.

## Description

GENESIS Scripting (gscript for short) is a technology I've developed to allow dynamic runtime execution of malware installation based on parameters determined at runtime.

Inspiration for this comes from my old AutoRune™ days and from the need for malware to basically become self aware without a bunch of duplicate overhead code.

GScript uses a JS V8 Virtual Machine to interpret your genesis script and allow it to hook into the malware initialization.

The Engine itself is referred commonly as "GSE" - Genesis Scripting Engine.

## Installation

We have created a command line SDK for gscript. You can download it from our Releases page:

<https://github.com/gen0cide/gscript/releases>

If you want to compile final binaries using gscripts compiler, you'll need the following dependencies installed and configured on your system:

 * Golang v1.10 or higher
 * jshint
 
After you have both of those, run:

```
$ go get github.com/gen0cide/gscript/cmd/gscript
```

## What is GENESIS btw?

GENESIS was created by @vyrus, @gen0cide, @emperorcow, and @ahhh for dynamically bundling multiple payloads into one dropper for faster deployment of implants for the CCDC Red Team.

For more information on this work we do every year, see my blog post outlining our toolbox:

 * <https://alexlevinson.wordpress.com/2017/05/09/know-your-opponent-my-ccdc-toolbox/>

GSE's goal is to allow intelligent deployment of those payloads.

## Variables

These variables are pre-defined and injected into the GENESIS VM at runtime for your convenience.

| Variable Name | Type     | Example                                                 | Purpose                                                                                                      |
|---------------|----------|---------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `USER_INFO`   | `object` | `{uid: 0, gid: 0, username: "root", home_dir: "/root"}` | Information about the User. Will be basically whatever is returned with https://golang.org/pkg/os/user/#User |
| `HOSTNAME`    | `string` | `example01`                                             | The hostname of the machine.                                                                                 |
| `IP_ADDRS`    | `array`  | `["127.0.0.1","192.168.1.5"]`                           | The IP addresses of the machine.                                                                             |
| `OS`          | `string` | `linux`                                                 | The operating system (basically `runtime.GOOS`)                                                              |
| `ARCH`        | `string` | `amd64`                                                 | The CPU architecture (basically `runtime.GOARCH`)                                                            |

## Function Docs

Check out `FUNCTIONS.md` in the root of this repo.

## Credits

Shoutouts to the homies:

 * vyrus
 * ahhh
 * cmccsec
 * carnal0wnage
 * indi303
 * emperorcow
 * rossja
\n\ntest
