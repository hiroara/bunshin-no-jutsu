# Bunshin no Jutsu

![Bunshin no Jutsu](./images/ninja.svg)

Bunshin-no-Jutsu is a CLI to duplicate your files.
This tool works as configuration-oriented.

And the name of the configuration file is `.makimono.yml`.
This can be used for backup.

For example, you can copy your project into a directory which syncs with Dropbox.

## Getting Started

### Installation with build

You can build the executable binary as follows:

`go install github.com/hiroara/bunshin-no-jutsu/bunshin`

## How to use

### `.makimono.yml`

`.makimono.yml` is a configuration file for this tool.
Please put this file in a root of a directory you want to sync.

`bunshin` detects `.makimono.yml` with going up to parent directories, and it determines the directory which contains `.makimono.yml` as a root directory to sync.

- _TODO_ : This can be generated with `bunshin init`.
- _TODO_ : About how to write this file, please refer the example.
