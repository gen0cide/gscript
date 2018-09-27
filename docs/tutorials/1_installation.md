# Installation

Welcome to the GSCRIPT install!

## On MacOS

- Install XCode or XCode CLI tools - `xcode-select --install`
- Install GoLang 1.10 minimum - https://golang.org/dl/
- Setup your `GOHOME`:

```sh
mkdir ~/go
export GOHOME=~/go
```

- Easily build out the dir structure by using "go get" to grab a package:

```sh
go get github.com/gen2brain/dlgs
ls -al ~/go/
```

## Download GSCRIPT

First pull down the gscript command line utility or source.
https://github.com/gen0cide/gscript (or your fork if you want to change things)
Make sure you save this project in your `GOPATH`, i.e:

```sh
$ echo $GOPATH
~/go/
```

Then the save would go here:

```sh
~/go/src/github/gen0cide/gscript/
```

## Build GSCRIPT

!!NOTE: This doesn't need to happen if you download gscript by doing `go get -a github.com/gen0cide/gscript/cmd/gscript`

First we need to get all of the dependencies:

- go get github.com/faith/color
- go get github.com/robertkrimen/otto
- go get golang.org/x/tools
- go get github.com/vigneshuvi/GoDateFormat
- go get github.com/urfave/cli
- go get github.com/uudashr/gopkgs
- go get github.com/tdewolff/minify/js
- go get github.com/sirupsen/logrus

Building the CLI tool:

```sh
cd ~/go/src/github.com/gen0cide/gscript/cmd/gscript
go build
cp ./gscript /usr/local/bin/
gscript --help
```

