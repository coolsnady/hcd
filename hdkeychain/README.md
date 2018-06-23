hdkeychain
==========

[![Build Status](http://img.shields.io/travis/coolsnady/hxd.svg)](https://travis-ci.org/coolsnady/hxd)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/coolsnady/hxd/hdkeychain)

Package hdkeychain provides an API for Decred hierarchical deterministic
extended keys (based on BIP0032).

A comprehensive suite of tests is provided to ensure proper functionality.  See
`test_coverage.txt` for the gocov coverage report.  Alternatively, if you are
running a POSIX OS, you can run the `cov_report.sh` script for a real-time
report.

## Feature Overview

- Full BIP0032 implementation
- Single type for private and public extended keys
- Convenient cryptograpically secure seed generation
- Simple creation of master nodes
- Support for multi-layer derivation
- Easy serialization and deserialization for both private and public extended
  keys
- Support for custom networks by registering them with chaincfg
- Obtaining the underlying EC pubkeys, EC privkeys, and associated coolsnady
  addresses ties in seamlessly with existing btcec and hxutil types which
  provide powerful tools for working with them to do things like sign
  transations and generate payment scripts
- Uses the btcec package which is highly optimized for secp256k1
- Code examples including:
  - Generating a cryptographically secure random seed and deriving a
    master node from it
  - Default HD wallet layout as described by BIP0032
  - Audits use case as described by BIP0032
- Comprehensive test coverage including the BIP0032 test vectors
- Benchmarks

## Installation and Updating

```bash
$ go get -u github.com/coolsnady/hxd/hdkeychain
```

## Examples

* [NewMaster Example](http://godoc.org/github.com/coolsnady/hxd/hdkeychain#example-NewMaster)
  Demonstrates how to generate a cryptographically random seed then use it to
  create a new master node (extended key).
* [Default Wallet Layout Example](http://godoc.org/github.com/coolsnady/hxd/hdkeychain#example-package--DefaultWalletLayout)
  Demonstrates the default hierarchical deterministic wallet layout as described
  in BIP0032.
* [Audits Use Case Example](http://godoc.org/github.com/coolsnady/hxd/hdkeychain#example-package--Audits)
  Demonstrates the audits use case in BIP0032.

## License

Package hdkeychain is licensed under the [copyfree](http://copyfree.org) ISC
License.
