/*

Package cm contains code for dealing with some "complicated maps":
multi-level maps (maps that contain maps as values) and dual-key
maps (maps that can be accessed by two distinct keys).

This package provides no locking in the datastructures. All locking is
the responsibility of code using these maps.

Multilevel Maps

Multi-level maps are maps that have other maps as their values.
This provides convenience functions for interacting with multi-level
maps. It is intended to work harmoniously with golang.org/x/maps, and
tries not to replicate any functionality already found there. For
instance, to get the first level of keys of these maps, simply pass them
as normal maps to maps.Keys. The internal maps are exported so normal map
operations work, so redundant operations already provided by range and
such are not implemented.

It is safe to write to these maps directly, no constraints maintained
by this code will be violated. The delete methods provided by the
multi-level maps will also clean up any higher-level maps left emptied by
a delete. Directly executing deletes on the lower-level maps yourself
will not automatically clean these maps up, which may also cause spurious
keys to appear in the KeySlice and KeyTree functions, so I advise against
deleting directly.

In theory, you can drop this into any existing multilevel map you already
have, and they should continue to work, give or take any type conversions
as you pass them around. You just also have the additional methods added by
this type.

This allows setting values when the previous levels do not exist yet, and
if all values from a particular sub-level are removed, all now-empty maps
will be removed.

Unlike single level maps where a sequence of the key values is the only
sensible representation of the keys, multi-level maps have more than one
useful representation. You can either look at the set of keys as a set
of tuples for all levels, or you can look at them as a tree. Each
representation has its costs and benefits, so this package provides both.

As multilevel maps are just Go maps under the hood, they scale the same
as Go maps do in general.

Dual Keyed Maps

A dual-keyed map is a map that allows you to lookup values by either
of the two keys. Under the hood, it is simply both possible maps, and
functions for setting and deleting by both keys.

For your convenience, the two maps are left exported so you can efficiently
read from them. Bear in mind that if you write directly to them, you will
break the guarantees provided by the methods!

Values are stored as given in both maps. This means that a dual-keyed map
consumes twice the resources of a normal map. As a result you may want to
consider storing pointers in the map if it is going to be very large.
This is targeted for cases where a dual-key map is very convenient, but
not large by modern standards. As you scale up needs like this you
eventually need a database.

For dual-key maps, it is obvious how to store them, with a reasonable
penalty. As you get into needs for three or more keys, the cost of this
technique multiplies resource consumption by the number of permutations
of the keys, which by three keys is already six times a single map.
So this package stops at dual-level maps.

*/
package cm