Complex Generic Maps for Go
===========================

[![Go Reference](https://pkg.go.dev/badge/github.com/thejerf/cm.svg)](https://pkg.go.dev/github.com/thejerf/cm)<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat)</a>

    go get github.com/thejerf/cm

    import "github.com/thejerf/cm"

`cm` provides some generic complex maps for Go;

  * The equivalent of `map[A]map[B]C` and `map[A]map[B]map[C]D`, in two
    flavors:
    * `MapMap` and `MapMapMap` implement that, with the constraint that the
      value type is `comparable`. This allows the equvialent of the maps'
      package `.Equal` method.
    * `MapMapAny` and `MapMapMapAny` are the same, but allowing the Value
      type to be `any`. This removes the `.Equal` method but allows storing
      any value.
  * `DualMap` implements a map that can be keyed by either of two keys,
    packaging up a `map[A]map[B]C` and `map[B]map[A]C` into a single
    coherent package. 
  * `MapSet` implements a map that contains sets, like `map[K]Set[V]`.
  * To support `MapSet`, there's a ful `Set` implementation.

    As I write this, there's a [proposal for a standard-library `Set` type
    on github](https://github.com/golang/go/discussions/47331). It has
    stalled out on lacking iterator support. This `Set` sidesteps that by
    simply accepting that it is based on a map type, thus permitting a
    standard `range` iteration. While a custom Set type may theoretically
    permit higher performance for large sets, and it does make sense for
    the standard library to consider such a case, this set type will be
    acceptable for quite significantly sized sets.

    This set is also biased in the direction of mutability and
    performance. So, for instance, `.Subtract` will modify the `Set` it is
    called on. I chose this because if you have a mutation-based library,
    but want to create a new set, it is easy to
    `set.Close().Subtract(set2)`, but if you have a library that only works
    by creating new values you can't get the higher performance of direct
    mutation. And in my experience, direct mutation is a frequently common
    case, as is cloning a set once and performing many mutation operations
    on it (like subtracting several sets).

There's nothing particularly "special" about this implementation, no magic
sauce or anything. Just code I've had to write in several projects and
would like to get factored out and into a well-tested library, rather than
write over and over.

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

version v0.4.0 is new, but headed towards production-grade. I intend to use
this in my code. But it is still pretty new.

PRs
===

I am actively accepting PRs. If this almost does what you want, please by
all means file a PR to add it rather than starting a new project.

However, I'm aiming for a library that *does not replace things we already
have built-ins or implementations in the
to-be-standard
[maps package](https://pkg.go.dev/golang.org/x/exp@v0.0.0-20220307200941-a1099baf94bf/maps)*. Especially
when implementations would be wildly slower than native support.

So, for instance, I'm not interested in a `.OnValues(func (val V) ...)`
implementation. The correct spelling of that is to go ahead and range over
the relevant maps. This prevents a huge amount of function call overhead
for something that the range-based for loops are more efficient for.

Users of this library are expected to understand this library as helpful
methods that work on the relevant data types, not as a complete replacement
data type for maps.

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

* 0.5.0:
    * Trimmed off the .Get functions from the `MapMap` and `MapMapMap` as
      they are unnecessary.
    * NewDualMap is removed because it doesn't matter; the first Set will
      initialize the values and there is no reason to need them to be
      initialized before that.
    * DualMap's two maps are now called Primary and Reverse, which will
      probably be easier to remember than Left and Right.
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
