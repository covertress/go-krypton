## Krypton Go

Official golang implementation of the Krypton protocol

          | Linux   | OSX | ARM | Windows | Tests
----------|---------|-----|-----|---------|------
develop   | [![Build+Status](https://build.krdev.com/buildstatusimage?builder=Linux%20Go%20develop%20branch)](https://build.krdev.com/builders/Linux%20Go%20develop%20branch/builds/-1) | [![Build+Status](https://build.krdev.com/buildstatusimage?builder=Linux%20Go%20develop%20branch)](https://build.krdev.com/builders/OSX%20Go%20develop%20branch/builds/-1) | [![Build+Status](https://build.krdev.com/buildstatusimage?builder=ARM%20Go%20develop%20branch)](https://build.krdev.com/builders/ARM%20Go%20develop%20branch/builds/-1) | [![Build+Status](https://build.krdev.com/buildstatusimage?builder=Windows%20Go%20develop%20branch)](https://build.krdev.com/builders/Windows%20Go%20develop%20branch/builds/-1) | [![Buildr+Status](https://travis-ci.org/krypton/go-krypton.svg?branch=develop)](https://travis-ci.org/krypton/go-krypton) [![codecov.io](http://codecov.io/github/krypton/go-krypton/coverage.svg?branch=develop)](http://codecov.io/github/krypton/go-krypton?branch=develop)
master    | [![Build+Status](https://build.krdev.com/buildstatusimage?builder=Linux%20Go%20master%20branch)](https://build.krdev.com/builders/Linux%20Go%20master%20branch/builds/-1) | [![Build+Status](https://build.krdev.com/buildstatusimage?builder=OSX%20Go%20master%20branch)](https://build.krdev.com/builders/OSX%20Go%20master%20branch/builds/-1) | [![Build+Status](https://build.krdev.com/buildstatusimage?builder=ARM%20Go%20master%20branch)](https://build.krdev.com/builders/ARM%20Go%20master%20branch/builds/-1) | [![Build+Status](https://build.krdev.com/buildstatusimage?builder=Windows%20Go%20master%20branch)](https://build.krdev.com/builders/Windows%20Go%20master%20branch/builds/-1) | [![Buildr+Status](https://travis-ci.org/krypton/go-krypton.svg?branch=master)](https://travis-ci.org/krypton/go-krypton) [![codecov.io](http://codecov.io/github/krypton/go-krypton/coverage.svg?branch=master)](http://codecov.io/github/krypton/go-krypton?branch=master)

[![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://godoc.org/github.com/krypton/go-krypton) 
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/krypton/go-krypton?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Automated development builds

The following builds are build automatically by our build servers after each push to the [develop](https://github.com/krypton/go-krypton/tree/develop) branch.

* [Docker](https://registry.hub.docker.com/u/krypton/client-go/)
* [OS X](http://build.krdev.com/builds/OSX%20Go%20develop%20branch/Mist-OSX-latest.dmg)
* Ubuntu
  [trusty](https://build.krdev.com/builds/Linux%20Go%20develop%20deb%20i386-trusty/latest/) |
  [utopic](https://build.krdev.com/builds/Linux%20Go%20develop%20deb%20i386-utopic/latest/)
* [Windows 64-bit](https://build.krdev.com/builds/Windows%20Go%20develop%20branch/Gkr-Win64-latest.zip)
* [ARM](https://build.krdev.com/builds/ARM%20Go%20develop%20branch/gkr-ARM-latest.tar.bz2)

## Building the source

For prerequisites and detailed build instructions please read the
[Installation Instructions](https://github.com/krypton/go-krypton/wiki/Building-Krypton)
on the wiki.

Building gkr requires both a Go and a C compiler.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    make gkr

## Executables

Go Krypton comes with several wrappers/executables found in 
[the `cmd` directory](https://github.com/krypton/go-krypton/tree/develop/cmd):

 Command  |         |
----------|---------|
`gkr` | Krypton CLI (krypton command line interface client) |
`bootnode` | runs a bootstrap node for the Discovery Protocol |
`krtest` | test tool which runs with the [tests](https://github.com/krypton/tests) suite: `/path/to/test.json > krtest --test BlockTests --stdin`.
`evm` | is a generic Krypton Virtual Machine: `evm -code 60ff60ff -gas 10000 -price 0 -dump`. See `-h` for a detailed description. |
`disasm` | disassembles EVM code: `echo "6001" | disasm` |
`rlpdump` | prints RLP structures |

## Command line options

`gkr` can be configured via command line options, environment variables and config files.

To get the options available:

    gkr help

For further details on options, see the [wiki](https://github.com/krypton/go-krypton/wiki/Command-Line-Options)

## Contribution

If you'd like to contribute to go-krypton please fork, fix, commit and
send a pull request. Commits who do not comply with the coding standards
are ignored (use gofmt!). If you send pull requests make absolute sure that you
commit on the `develop` branch and that you do not merge to master.
Commits that are directly based on master are simply ignored.

See [Developers' Guide](https://github.com/krypton/go-krypton/wiki/Developers'-Guide)
for more details on configuring your environment, testing, and
dependency management.
