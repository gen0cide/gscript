![Gscript Logo](https://i.imgur.com/16lZGrA.png)

# Genesis Scripting Engine (gscript)

> Security framework to rapidly implement custom droppers for all three major operating systems



[![CircleCI](https://circleci.com/gh/gen0cide/gscript/tree/master.svg?style=svg)](https://circleci.com/gh/gen0cide/gscript/tree/master)

## About

Gscript is a framework for building multi-tenant executors for several implants in a stager. The engine works by embedding runtime logic (powered by the V8 Javascript Virtual Machine) for each persistence technique. This logic gets run at deploy time on the victim machine, in parallel for every implant contained with the stager. The Gscript engine leverages the multi-platform support of Golang to produce final stage one binaries for Windows, Mac, and Linux. 

**We encourage you to read through the slides from DEFCON26:**

https://docs.google.com/presentation/d/1kHdz8DY0Zn44yn_XrZ2RVqDY1lpADThLPNPwHP-njbc/edit?usp=sharing


## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Docs](#docs)
- [Shoutouts](#shoutouts)

## Features

- Easy to learn and write - uses javascript.
- Portable - Compile droppers for OSX, Windows, or Linux from any OS.
- Robust - Script's are isolated from each other in a safe execution way.
- Fast.
- Extensible - Can link native Golang packages directly into your Javascript.

## Installation

### Docker (Easiest)

If you have docker installed, you can run:

```
$ docker pull gen0cide/gscript:v1
```

Make a local directory where you can share files between your local machine and the docker container. Replace `$LOCAL_DIR` in the following command with the path to that:

```
$ docker run -it -v $LOCAL_DIR:/root/share gen0cide/gscript:v1
```

Thats it! You're good to go.


### Local (Good for advanced)

**Local installation requires you to have a Golang compiler setup and working on your machine. If you need to do this, you can grab an installer [here](https://golang.org/dl/). Make sure `$GOPATH/bin` is in your `$PATH`.**


After that, all you need to do is run:

```
$ go get github.com/gen0cide/gscript/cmd/gscript
```

## Quick Start

Check out the tutorial docs here:

https://github.com/gen0cide/gscript/tree/master/docs/tutorials

If you want to see example scripts, we have a separate repo you can clone:

https://github.com/ahhh/gscripts

## Docs

Here's a list of docs and tutorials that might be helpful for you:

 - [Godoc for Engine and Compiler](https://godoc.org/github.com/gen0cide/gscript)
 - [Tutorials in docs/tutorials](https://github.com/gen0cide/gscript/tree/master/docs/tutorials)
 
(more to come soon)

## Shoutouts

mentors, contributors, and great friends of gscript

- @cmc
- @hecfblog
- @ccdcredteam
- @1njecti0n
- @emperorcow
- @vyrus001
- @kos
- @davehughes
- @maus
- @javuto

