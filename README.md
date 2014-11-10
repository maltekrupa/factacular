# factacular

Factacular is a Command Line Interface (CLI) for convinient querying of the PuppetDB.

## Repository status

Travis-ci.org: [![Build Status](https://travis-ci.org/temal-/factacular.svg?branch=master)](https://travis-ci.org/temal-/factacular)

## Install

```
go get github.com/temal-/factacular
```

## Prerequisites

`factacular` needs to know the HTTP address of your PuppetDB. You can choose between
the following options:
- `--puppetdb http://puppetdb.example.com`
- `-p http://puppetdb.example.com`
- `export PUPPETDB_HOST=http://puppetdb.example.com`

If you have `$GOROOT/bin` in your `$PATH` you can start with `factacular`.

## Help

```
$ factacular
NAME:
   factacular - Get facts and informations from PuppetDB.

USAGE:
   factacular [global options] command [command options] [arguments...]

VERSION:
   0.3.2

COMMANDS:
   list-facts, lf   List all available facts.
   list-nodes, ln   List all available nodes.
   node-facts, nf   List all facts for a specific node.
   fact, f          List fact for all nodes (which have this fact).
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --puppetdb, -p 'http://localhost:8080'   PuppetDB host. [$PUPPETDB_HOST]
   --help, -h                               show help
   --version, -v                            print the version
```

## Examples

```
$ factacular fact os
PuppetDB host: http://localhost:8080
Fact per node:
FQDN00 - os - {"family":"RedHat","name":"CentOS","release":{"full":"6.5","major":"6","minor":"5"}}
FQDN01 - os - {"family":"Debian","lsb":{"distcodename":"squeeze","distdescription":"Debian GNU/Linux 6.0.9 (squeeze)","distid":"Debian","distrelease":"6.0.9","majdistrelease":"6","minordistrelease":"0"},"name":"Debian","release":{"full":"6.0.9","major":"6","minor":"0"}}
FQDN02 - os - {"family":"Debian","lsb":{"distcodename":"squeeze","distdescription":"Debian GNU/Linux 6.0.9 (squeeze)","distid":"Debian","distrelease":"6.0.9","majdistrelease":"6","minordistrelease":"0"},"name":"Debian","release":{"full":"6.0.9","major":"6","minor":"0"}}
FQDN03 - os - {"family":"RedHat","name":"CentOS","release":{"full":"6.5","major":"6","minor":"5"}}
```
