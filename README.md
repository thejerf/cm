Complex Generic Maps for Go
===========================

[![Go Reference](https://pkg.go.dev/badge/github.com/thejerf/cm.svg)](https://pkg.go.dev/github.com/thejerf/cm)<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat)</a>

    go get github.com/thejerf/cm

    import "github.com/thejerf/cm"

`cm` provides some generic complex maps for Go;

  * `MapMap`, `MapMapMap`, `MapMapAny`, and `MapMapMapAny` provide maps
    based on two and three keys that provide some convenience functions
    around  the equivalent of `map[K1]map[K2]Value` and
    `map[K1]map[K2]map[K3]Value` respectively, for easily setting and
    fetching values.
  * `DualMap` implements a map that can be keyed by either of two keys,
    packaging up a `map[Left]map[Right]Value` and
    `map[Right]map[Left]Value` into a single coherent package.
  * `MapSet` implements a map that contains sets, like `map[K]Set[V]`.

There's nothing particularly "special" about this implementation, no magic
sauce or anything. Just code I've had to write in several projects and
would like to get factored out and into a well-tested library, rather than
write over and over.

MapMap(Map) vs. MapMap(Map)Any
==============================

I took inspiration from the
current
[maps package](https://pkg.go.dev/golang.org/x/exp@v0.0.0-20220307200941-a1099baf94bf/maps) and
tried to implement all the functionality present there on the multi-level
maps.

In order to implement `.Equal`, the values of the map must be `comparable`.
The `MapMap` and `MapMapMap` types implement this restriction, because
many things stored in maps are indeed `comparable`. All methods on
`MapMap(Map)Any` are available for those maps, and they also have a
`.Equal`.

If you want to store a non-`comparable` value in the map, `MapMap(Map)Any`
is available. It drops the `comparable` restriction on the Value, and loses
the `.Equal` method as a result. Switching the type of the map should not
be very complicated in general, as the method sets are almost identical.

Performance
===========

All of these methods are extremely strong candidates for inlining. Some
brief checks with `-gcflags="-m"` on some test programs suggest that
they are indeed all inlined. Consequently, this library should generally
be zero-performance-impact versus having directly written the code.
(There are a few places where nil is checked for where you might not have,
but a highly-predictable branch should be lost in the noise compared to
what even one map lookup requires.)

Status
======

version v0.3.0 is new, but headed towards production-grade. I intend to use
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
about backwards compatibility. Change is still happening frequently as I
hone in on the best solutions.

* 0.4.0:
    * Removed MapSlice because upon an even more careful review of the
      behavior of nil slices, it adds nothing. The `.Append` method I wrote
      can also just be written as `m[key] = append(m[key], vals...)`
      without loss. The remaining `.Set` method that removes the key
      entirely for empty slices doesn't justify a full datatype.
* 0.3.0:
    * Worked out a way to recover the `.Equal` method without a lot
      of nasty casting on values.
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
