Complex Generic Maps for Go
===========================

[![Go Reference](https://pkg.go.dev/badge/github.com/thejerf/cm.svg)](https://pkg.go.dev/github.com/thejerf/cm)

    go get github.com/thejerf/cm

    import "github.com/thejerf/cm"

`cm` provides some generic complex maps for Go;

  * `MapMap` and `MapMapMap` provide maps based on two and
    three keys that provide some convenience functions around
    the equivalent of `map[K1]map[K2]Value` and
    `map[K1]map[K2]map[K3]Value` respectively, for easily setting and
    fetching values.
  * DualMap implements a map that can be keyed by either of two keys,
    packaging up a `map[Left]map[Right]Value` and
    `map[Right]map[Left]Value` into a single coherent package.

There's nothing particularly "special" about this implementation, no magic
sauce or anything. Just code I've had to write in several projects and
would like to get factored out and into a well-tested library, rather than
write over and over.

Status
======

version v0.2.0 is new, but headed towards production-grade. I intend to use
this in my code.

PRs
===

I am actively accepting PRs. If this almost does what you want, please by
all means file a PR to add it rather than starting a new project.

The major caveat I'd give is, I'm trying not to be redundant to anything
Go is going to provide anyhow in the future. I expect people using this
project to understand that these types are _also_ normal Go maps that can
be passed to
the
[maps package](https://pkg.go.dev/golang.org/x/exp@v0.0.0-20220307200941-a1099baf94bf/maps) (or
whatever future version of that exists).

Code Signing
============

All commits and releases will be signed with a GPG key, as verified by
GitHub. It is the ["jerf" keybase account](https://keybase.io/jerf) key.

(Bear in mind that due to the nature of how git commit signing works, there
may be runs of unverified commits. What matters is that the top one is
signed.)

Changelog
=========

cm uses semantic versioning.

At the moment, this is in pre-release, which means no guarantees whatsoever
about backwards compatibility. I would expect there's a reasonably chance
of at least some incompatible changes in the future, though I wouldn't
expect them to be major.

* 0.2.0:
    * Renamed almost everything to something shorter, but also more
      comprehensible, I hope.
    * Additional data types:
      * Set
      * MapSet
      * MapSlice
* 0.1: Initial release to GitHub. This release has not been publicized as
  it is still missing a couple of key methods:
    * Methods for getting values.
    * Methods for getting full key/value tuples.
    * KeyTree2 may need some renaming; it works in the context of a
      MultiLevelMap2 but when recursively used by the MultiLevelMap3 the
      struct names are a bit off.
