hcd
====

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

## Important Info for RC1 Testnet
RC1 Testnet is being restarted with updated parameters to prevent the issue of no blocks being made for long periods of time after ASIC miners join the network and leave once the difficulty has increased. Users who have previously connected to RC1 testnet will need to follow the instructions below to ensure they are on the most up-to-date network. When the status below is "Active" you can follow the install instructions below to connect. When the status is "Temporarily Inactive" please wait until the status is updated to "Active" before attempting to join the testnet. *If you have previously connected to RC1 testnet you will need to follow the instructions below to connect to the most recent testnet! If you cannot find peers it is because the DNS seeds are not yet active. In this case you must wait until the Status is updated to Active and follow the instructions below to grab the latest peer info. Additionally you will need to create a new wallet and set a new config file.*

### Status: *Temporarily Inactive*

#### Linux
```
$ rm -rf $HOME/.hcd
$ rm -rf $HOME/.hcwallet
$ cd $HOME/go/src/github.com/coolsnady/hcd
$ git pull
$ go install $(glide nv)
```

#### Mac/OSX
```
$ rm -rf $HOME/Library/Application\ Support/Hcwallet
$ rm -rf $HOME/Library/Application\ Support/Hcd
$ cd $HOME/go/src/github.com/coolsnady/hcd
$ git pull
$ go install $(glide nv)
```

hcd is a Hc full node implementation written in Go (golang).

This acts as a chain daemon for the Hc cryptocurrency.
hcd maintains the entire past transactional ledger of Hc and allows
 relaying of quantum resistant transactions to other Hc nodes across the world.

Note: To send or receive funds and join Proof-of-Stake mining, you will also need
[hcwallet](https://github.com/coolsnady/hcwallet).

Hc is forked from [decred](https://github.com/decred) and [btcd](https://github.com/btcsuite/btcd) which are full node implementations written in Go. Both projects are ongoing and under active development. Since hcd is synced and will merge with upstream commits from hcd and btcd, it will get the benefit of both hcd and btcd's ongoing upgrades to staking, voting, peer and connection handling, database optimization and other blockchain related technology improvements. Advances made by hcd can also be pulled back upstream to hcd and btcd including quantum resistant signature schemes and more.

## Current State
This project is currently under active development and is in a Beta state. The default branch of hcd is currently testnet1. Please make sure to use --testnet flag when running hcd and report any issues by using the integrated issue tracker. Do not yet use this software yet as a store of value!

## Requirements

[Go](http://golang.org) 1.7 or newer.

## Getting Started

- hcd (and utilities) will now be installed in either ```$GOROOT/bin``` or
  ```$GOPATH/bin``` depending on your configuration.  If you did not already
  add the bin directory to your system path during Go installation, we
  recommend you do so now.

## Installing

#### Build from Source

- **Dependencies**

  Glide is used to manage project dependencies and provide reproducible builds.
  To install:

  `go get -u github.com/Masterminds/glide`

Unfortunately, the use of `glide` prevents a handy tool such as `go get` from
automatically downloading, building, and installing the source in a single
command.  Instead, the latest project and dependency sources must be first
obtained manually with `git` and `glide`, and then `go` is used to build and
install the project.

**Getting the source**:

For a first time installation, the project and dependency sources can be
obtained manually with `git` and `glide` (create directories as needed):

```
git clone https://github.com/coolsnady/hcd $GOPATH/src/github.com/coolsnady/hcd
cd $GOPATH/src/github.com/coolsnady/hcd
glide install
go install $(glide nv)
```

To update an existing source tree, pull the latest changes and install the
matching dependencies:

```
cd $GOPATH/src/github.com/coolsnady/hcd
git pull
glide install
go install $(glide nv)
```

## Running

Make sure you are working the correct GOPATH and run the following in your terminal:

```
hcd -u YOURUNIQUERPCUSERNAME -P YOURUNIQUERPCPASSWORD --testnet
```

To use your node for mining add the miningaddr flag when running hcd:

```
hcd -u YOURUNIQUERPCUSERNAME -P YOURUNIQUERPCPASSWORD --testnet --miningaddr=YOURTESTNETADDRESS
```

To generate a testnet mining address you must install [hcwallet](https://github.com/coolsnady/hcwallet)

To begin CPU mining after hcd is already running you can run the following in your terminal:

```
hcctl -u YOURUNIQUEUSERNAME -P YOURUNIQUEPASSWORD --testnet setgenerate true
```

## Issue Tracker

The [integrated github issue tracker](https://github.com/coolsnady/hcd/issues)
is used for this project.

## Documentation

The documentation is a work-in-progress.  It is located in the [docs](https://github.com/coolsnady/hcd/tree/master/docs) folder.

## License

hcd is licensed under the [copyfree](http://copyfree.org) ISC License.
