/*

Package cm contains generic "complicated maps": multi-level maps, dual-key
maps, and maps containing sets.

This package provides no locking in the datastructures. All locking is
the responsibility of code using these maps.

This code panics analogously to normal map behaviors. When there is no
existing map behavior to guide, it tries to match the same logic Go
normally uses. This is justified because these are just wrappers around
maps, rather than independent data structures. Most or all of the places
where this library panics is places where the code was going to panic
anyhow; the panic calls in the code simply offer more direct guidance on
the problem rather than panicking deep in library code.

Similarly, none of these data structures are thread-safe on their own,
just like conventional Go maps.

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
keys to appear in the KeySlice method. Otherwise it is safe too.

In theory, you can drop this into any existing multilevel map you already
have, and they should continue to work, give or take any type conversions
as you pass them around. You just also have the additional methods added by
this type.

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
read from them. Writing directly to them will violate the guarantees
provided by this implementation and should generally not be done.

Values are stored as given in both maps. This means that a dual-keyed map
consumes twice the resources of a normal map. This is targeted for cases
where a dual-key map is very convenient, but not large by modern
standards. As you scale up needs like this you eventually need a database.

Because this simply stores the maps in both directions, you may want
to double-check before using pointer types for either type. It is legal
in Go to use pointers to key maps, but it may not have the desired
or expected result, as it will result in one of the two directions keying
off of object identity rather than value. This has its uses too, though.

For dual-key maps, it is obvious how to store them, with a reasonable
penalty. As you get into needs for three or more keys, the cost of this
technique multiplies resource consumption by the number of permutations
of the keys, which by three keys is already six times a single map.
So this package stops at dual-level maps.

Map Sets

A MapSet is a map, whose value is a set. Several convenience functions can
be implemented for manipulating such values.

As a consequence of offering this functionality, this package also provides
a Set implementation.

Key Trees And Key Slices

Each of these structures implements the ability to get data structures
representing the set of all keys, or keys and values in the set, as a
single static data structure.

It is an anti-pattern to use them as such:

    keySlice := someMapMap.KeySlice()
    for _, keys := range keySlice {
        // ...
    }

This causes the needless instantiation of a data structure in memory.
This should be written as

    for key1, submap := range someMapMap {
        for key2, val := range submap {
            // do work here
        }
    }

Normal use of KeySlice and KeyTree would be sorting it somehow before
iterating, or possibly serializing them somewhere.

*/
package cm
