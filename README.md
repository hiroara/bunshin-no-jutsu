# Bunshin no Jutsu

![Bunshin no Jutsu](./images/ninja.svg)

Bunshin-no-Jutsu is a CLI to make complete copies of your files.
You can use this for making backup :)

This tool works as a configuration-oriented tool.
The name of the configuration file is `.makimono.yml` :scroll:

## Usecase

### Working with Dropbox

For example, you can copy your project under a directory which syncs with Dropbox.

Dropbox is sometimes heavy especially when we have many files, and doing many writing/deleting files.
(e.g. Developing a nodejs project under Dropbox directory and build it)

In this case, we can reduce the waste of machine resources by using Bunshin-no-Jutsu.

This can be done by syncing two directories:

- `~/dev/some-project`: Working directory
- `~/Dropbox/dev/some-project`: Backup directory to sync with Dropbox

## Getting Started

### Install released version (MacOS)

```
$ sudo curl -o /usr/local/bin/bunshin -sL https://github.com/hiroara/bunshin-no-jutsu/releases/download/1.0.1/bunshin-darwin-amd64 && sudo chmod +x /usr/local/bin/bunshin
$ bunshin version
```

### Installation with build

You can build the executable binary as follows:

```
$ go get -u github.com/hiroara/bunshin-no-jutsu/bunshin
$ bunshin version
```

## How to use

### :dancers: `bunshin`

`bunshin` is a CLI of Bunshin-no-Jutsu.
When you want to sync files, all you have to do is simply type `bunshin`.

For more details, please type `bunshin --help`.

### :scroll: `.makimono.yml`

`.makimono.yml` is a configuration file for Bunshin-no-Jutsu.
Please put this file in a root directory you want to sync.

`bunshin` command detects `.makimono.yml` with going up to parent directories, and it determines the directory which contains `.makimono.yml` as a root directory to sync.

For more details, please see the example file: [./docs/.makimono.yml](./docs/.makimono.yml)

- _TODO_ : This can be generated with `bunshin init`.
